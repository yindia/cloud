package gormimpl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/riverqueue/river"
	"gorm.io/gorm"

	interfaces "task/server/repository/interface"
	models "task/server/repository/model/task"
)

var (
	taskHistoryOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_history_repository_operations_total",
			Help: "The total number of task history repository operations",
		},
		[]string{"operation", "status"},
	)
)

// TaskHistoryRepo handles database operations for task history entries.
type TaskHistoryRepo struct {
	db          *gorm.DB
	riverClient *river.Client[*sql.Tx]
}

// CreateTaskHistory creates a new history entry for a task.
// It returns the created TaskHistory object and any error encountered.
func (s *TaskHistoryRepo) CreateTaskHistory(ctx context.Context, history models.TaskHistory) (models.TaskHistory, error) {
	if err := s.db.Create(&history).Error; err != nil {
		taskHistoryOperations.WithLabelValues("create", "error").Inc()
		return models.TaskHistory{}, fmt.Errorf("failed to create task history: %w", err)
	}
	taskHistoryOperations.WithLabelValues("create", "success").Inc()
	return history, nil
}

// GetTaskHistory retrieves all history entries for a task by its ID.
// It returns a slice of TaskHistory objects and any error encountered.
func (s *TaskHistoryRepo) GetTaskHistory(ctx context.Context, taskID uint) ([]models.TaskHistory, error) {
	var histories []models.TaskHistory
	if err := s.db.Where("task_id = ?", taskID).Find(&histories).Error; err != nil {
		taskHistoryOperations.WithLabelValues("get", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve task history by task ID: %w", err)
	}
	taskHistoryOperations.WithLabelValues("get", "success").Inc()
	return histories, nil
}

// ListTaskHistories retrieves all history entries for a given task, sorted by ID in ascending order.
// It returns a slice of TaskHistory objects and any error encountered.
func (s *TaskHistoryRepo) ListTaskHistories(ctx context.Context, taskID uint) ([]models.TaskHistory, error) {
	var histories []models.TaskHistory

	// Added Order clause to sort by ID in descending order (latest first)
	if err := s.db.Where("task_id = ?", taskID).Order("id ASC").Find(&histories).Error; err != nil {
		taskHistoryOperations.WithLabelValues("list", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve task histories: %w", err)
	}

	taskHistoryOperations.WithLabelValues("list", "success").Inc()
	return histories, nil
}

// NewTaskHistoryRepo creates and returns a new instance of TaskHistoryRepo.
// It takes a GORM database connection and a River client as parameters.
func NewTaskHistoryRepo(db *gorm.DB, riverClient *river.Client[*sql.Tx]) interfaces.TaskHistoryRepo {
	return &TaskHistoryRepo{
		db:          db,
		riverClient: riverClient,
	}
}
