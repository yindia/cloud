package repositories

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	interfaces "task/server/repository/interface"
	tasks "task/server/repository/model/task"

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

	return NewPostgresRepo(db), nil
}
