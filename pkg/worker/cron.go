package worker

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	cloudv1 "task/pkg/gen/cloud/v1"

	"task/pkg/x"
	"task/server/repository/model/task"
	"time"

	"connectrpc.com/connect"
	"github.com/riverqueue/river"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ReconcileTaskWorkerArgs contains the arguments for the ReconcileTaskWorker.
type ReconcileTaskWorkerArgs struct {
	Status int `json:"status"`
	URL    string
}

// Kind returns the kind of the task argument.
func (ReconcileTaskWorkerArgs) Kind() string { return "reconcile_tasks" }

// InsertOpts returns the insertion options for the task.
func (ReconcileTaskWorkerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 1}
}

// ReconcileTaskWorker is the worker implementation for processing and reconciling tasks.
type ReconcileTaskWorker struct {
	river.WorkerDefaults[ReconcileTaskWorkerArgs]
	Logger *slog.Logger
	db     *gorm.DB
}

// Work processes a single reconciliation job for tasks.
// It finds tasks that have been in a specific status for too long and updates them to QUEUED status.
func (w *ReconcileTaskWorker) Work(ctx context.Context, job *river.Job[ReconcileTaskWorkerArgs]) error {
	w.Logger = slog.Default().With("worker", "ReconcileTaskWorker")
	w.Logger.Info("Starting task reconciliation")

	sqlDB, err := sql.Open("pgx", job.Args.URL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer sqlDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize GORM: %w", err)
	}

	w.db = db

	runningTasks, err := w.fetchRunningTasks(ctx, job.Args.Status)
	if err != nil {
		return fmt.Errorf("failed to fetch running tasks: %w", err)
	}

	w.Logger.Info("Found running tasks", "count", len(runningTasks))

	updatedCount, err := w.updateAndQueueTasks(ctx, runningTasks)
	if err != nil {
		return fmt.Errorf("failed to update and queue tasks: %w", err)
	}

	w.Logger.Info("Finished processing tasks", "updated_count", updatedCount)
	return nil
}

func (w *ReconcileTaskWorker) fetchRunningTasks(ctx context.Context, status int) ([]task.Task, error) {
	var runningTasks []task.Task
	twentyMinutesAgo := time.Now().Add(-time.Duration(x.CRON_TIME) * time.Minute)

	err := w.db.WithContext(ctx).
		Where("status = ? AND created_at <= ?", status, twentyMinutesAgo).
		Find(&runningTasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query running tasks: %w", err)
	}

	return runningTasks, nil
}

func (w *ReconcileTaskWorker) updateAndQueueTasks(ctx context.Context, tasks []task.Task) (int, error) {
	updatedCount := 0
	for _, t := range tasks {
		if err := w.updateTaskStatus(ctx, t); err != nil {
			w.Logger.Error("Failed to update task status", "task_id", t.ID, "error", err)
			continue
		}
		updatedCount++
	}
	return updatedCount, nil
}

func (w *ReconcileTaskWorker) updateTaskStatus(ctx context.Context, t task.Task) error {
	cloud, err := x.CreateClient(os.Getenv("SERVER_ENDPOINT"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	req := &cloudv1.UpdateTaskStatusRequest{
		Id:      int32(t.ID),
		Status:  cloudv1.TaskStatusEnum_QUEUED,
		Message: "Task has been queued again",
	}

	_, err = cloud.UpdateTaskStatus(ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	w.Logger.Info("Updated task status to QUEUED", "task_id", t.ID)
	return nil
}
