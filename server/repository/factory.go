package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	cloudv1 "task/pkg/gen/cloud/v1"
	"task/pkg/worker"
	"task/pkg/x"
	interfaces "task/server/repository/interface"
	tasks "task/server/repository/model/task"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetRepository(url string, workerCount int, maxConns int) (interfaces.TaskManagmentInterface, error) {
	// Open database connection

	sqlDB, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections
	sqlDB.SetMaxOpenConns(maxConns)
	// Set the maximum number of open connections
	sqlDB.SetMaxOpenConns(maxConns)

	// Set the maximum number of idle connections
	sqlDB.SetMaxIdleConns(maxConns / 2)

	// Set the maximum lifetime of a connection
	sqlDB.SetConnMaxLifetime(time.Hour)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	migrator, err := rivermigrate.New(riverdatabasesql.New(sqlDB), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create river migrator: %w", err)
	}

	_, err = migrator.Migrate(context.Background(), rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{})
	if err != nil {
		panic(err)
	}

	// Set up River workers and client
	workers := river.NewWorkers()
	if err := river.AddWorkerSafely(workers, &worker.TaskWorker{}); err != nil {
		return nil, fmt.Errorf("failed to add TaskWorker: %w", err)
	}

	if err := river.AddWorkerSafely(workers, &worker.ReconcileTaskWorker{}); err != nil {
		return nil, fmt.Errorf("failed to add TaskWorker: %w", err)
	}

	// TODO: Add comprehensive documentation
	// We added periodic reconciliation jobs to handle stuck tasks. These jobs will:
	// 1. Get a list of all stuck tasks
	// 2. Change their status to "queued"
	// 3. Enqueue them for processing
	// Currently, we're making direct DB changes, but ideally, the server should
	// implement an API to handle this logic. In the future, these scheduled jobs
	// should just call the API instead of modifying the database directly.

	var reconcileTasks = []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(x.CRON_TIME)*time.Second),
			func() (river.JobArgs, *river.InsertOpts) {
				return worker.ReconcileTaskWorkerArgs{

					Status: int(cloudv1.TaskStatusEnum_RUNNING),
					URL:    url,
				}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(x.CRON_TIME)*time.Second),
			func() (river.JobArgs, *river.InsertOpts) {
				return worker.ReconcileTaskWorkerArgs{
					Status: int(cloudv1.TaskStatusEnum_QUEUED),
					URL:    url,
				}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	}

	riverClient, err := river.NewClient(riverdatabasesql.New(sqlDB), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: workerCount},
		},
		Workers:      workers,
		PeriodicJobs: reconcileTasks,
		ErrorHandler: &worker.CustomErrorHandler{},
		Logger:       slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create River client: %w", err)
	}

	// Start River client
	if err := riverClient.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start River client: %w", err)
	}

	// Perform database migrations
	if err = db.AutoMigrate(&tasks.Task{}, &tasks.TaskHistory{}); err != nil {
		return nil, fmt.Errorf("failed to run auto migrations: %w", err)
	}

	// Create necessary indexes
	indexes := []struct {
		name string
		sql  string
	}{
		{"idx_task_id_created_at", "CREATE INDEX IF NOT EXISTS idx_task_id_created_at ON task_histories (task_id, created_at DESC)"},
		{"idx_type_status", "CREATE INDEX IF NOT EXISTS idx_type_status ON tasks (type, status)"},
		{"idx_created_at", "CREATE INDEX IF NOT EXISTS idx_created_at ON tasks (created_at)"},
		{"idx_status_created_at", "CREATE INDEX IF NOT EXISTS idx_status_created_at ON tasks (status, created_at)"},
	}

	for _, idx := range indexes {
		if err := db.Exec(idx.sql).Error; err != nil {
			return nil, fmt.Errorf("failed to create index %s: %w", idx.name, err)
		}
		slog.Info("Created index", "name", idx.name)
	}

	return NewPostgresRepo(db, riverClient), nil
}
