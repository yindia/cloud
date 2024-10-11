package gormimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"

	interfaces "task/server/repository/interface"
	models "task/server/repository/model/task"
)

var (
	taskOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_repository_operations_total",
			Help: "The total number of task repository operations",
		},
		[]string{"operation", "status"},
	)
	taskLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "task_repository_operation_duration_seconds",
			Help:    "Duration of task repository operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// TaskRepo implements the TaskRepo interface using GORM for database operations
// and River for task queue management.
type TaskRepo struct {
	db *gorm.DB
}

// CreateTask creates a new task in the database and enqueues it for processing.
// It returns the created task with its assigned ID or an error if the operation fails.
func (s *TaskRepo) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	timer := prometheus.NewTimer(taskLatency.WithLabelValues("create"))
	defer timer.ObserveDuration()

	result := s.db.Create(&task)
	if result.Error != nil {
		taskOperations.WithLabelValues("create", "error").Inc()
		return models.Task{}, fmt.Errorf("failed to create task: %w", result.Error)
	}

	if task.ID == 0 {
		taskOperations.WithLabelValues("create", "error").Inc()
		return models.Task{}, fmt.Errorf("failed to get task ID after creation")
	}

	taskOperations.WithLabelValues("create", "success").Inc()
	return task, nil
}

// GetTaskByID retrieves a task from the database by its ID.
// It returns a pointer to the task if found, or an error if the task doesn't exist or if the operation fails.
func (s *TaskRepo) GetTaskByID(ctx context.Context, taskID uint) (*models.Task, error) {
	timer := prometheus.NewTimer(taskLatency.WithLabelValues("get"))
	defer timer.ObserveDuration()

	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		taskOperations.WithLabelValues("get", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve task by ID: %w", err)
	}
	taskOperations.WithLabelValues("get", "success").Inc()
	return &task, nil
}

// UpdateTaskStatus updates the status of a task identified by its ID.
// It returns an error if the update operation fails.
func (s *TaskRepo) UpdateTaskStatus(ctx context.Context, taskID uint, status int) error {
	timer := prometheus.NewTimer(taskLatency.WithLabelValues("update_status"))
	defer timer.ObserveDuration()

	if err := s.db.Model(&models.Task{}).Where("id = ?", taskID).Update("status", status).Error; err != nil {
		taskOperations.WithLabelValues("update_status", "error").Inc()
		return fmt.Errorf("failed to update task status: %w", err)
	}

	taskOperations.WithLabelValues("update_status", "success").Inc()
	return nil
}

// ListTasks retrieves a paginated list of tasks from the database, filtered by status and type.
// The 'limit' parameter specifies the maximum number of tasks to return,
// 'offset' determines the starting point for pagination,
// 'status' allows filtering by task status, and 'taskType' allows filtering by task type.
// It returns a slice of tasks and an error if the operation fails.
func (s *TaskRepo) ListTasks(ctx context.Context, limit, offset int, status int, taskType string) ([]models.Task, error) {
	timer := prometheus.NewTimer(taskLatency.WithLabelValues("list"))
	defer timer.ObserveDuration()

	var tasks []models.Task
	query := s.db.Limit(limit).Offset(offset)

	// Apply filters if they are provided
	if status != 5 {
		query = query.Where("status = ?", status)
	}
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}

	// Execute the query
	if err := query.Find(&tasks).Error; err != nil {
		taskOperations.WithLabelValues("list", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	taskOperations.WithLabelValues("list", "success").Inc()
	return tasks, nil
}

// GetTaskStatusCounts retrieves the count of tasks for each status.
// It returns a map where the key is the status code and the value is the count of tasks with that status.
// An error is returned if the operation fails.
func (s *TaskRepo) GetTaskStatusCounts(ctx context.Context) (map[int]int64, error) {
	timer := prometheus.NewTimer(taskLatency.WithLabelValues("status_counts"))
	defer timer.ObserveDuration()

	var results []struct {
		Status int
		Count  int64
	}

	if err := s.db.Model(&models.Task{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&results).Error; err != nil {
		taskOperations.WithLabelValues("status_counts", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve task status counts: %w", err)
	}

	counts := make(map[int]int64)
	for _, result := range results {
		counts[result.Status] = result.Count
	}

	taskOperations.WithLabelValues("status_counts", "success").Inc()
	return counts, nil
}

// GetStalledTasks retrieves tasks with status "unknown" or "queue" that have been in that state for more than 10 seconds.
// It returns a slice of tasks and an error if the operation fails.
func (s *TaskRepo) GetStalledTasks(ctx context.Context) ([]models.Task, error) {
	timer := prometheus.NewTimer(taskLatency.WithLabelValues("get_stalled"))
	defer timer.ObserveDuration()

	var tasks []models.Task
	tenSecondsAgo := time.Now().Add(-30 * time.Second)

	err := s.db.Where("(status = ?) AND updated_at < ?",
		4, tenSecondsAgo).
		Find(&tasks).Error

	if err != nil {
		taskOperations.WithLabelValues("get_stalled", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve stalled tasks: %w", err)
	}

	taskOperations.WithLabelValues("get_stalled", "success").Inc()
	return tasks, nil
}

// NewTaskRepo creates and returns a new instance of TaskRepo.
// It requires a GORM database connection and a River client for task queue management.
func NewTaskRepo(db *gorm.DB) interfaces.TaskRepo {
	return &TaskRepo{
		db: db,
	}
}
