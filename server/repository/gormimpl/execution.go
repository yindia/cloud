package gormimpl

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"

	interfaces "task/server/repository/interface"
	models "task/server/repository/model/task"
)

var (
	executionOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "execution_repository_operations_total",
			Help: "The total number of execution repository operations",
		},
		[]string{"operation", "status"},
	)
	executionLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "execution_repository_operation_duration_seconds",
			Help:    "Duration of execution repository operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// ExecutionRepo implements the ExecutionRepo interface using GORM for database operations
// and River for execution queue management.
type ExecutionRepo struct {
	db *gorm.DB
}

// CreateExecution creates a new execution in the database and enqueues it for processing.
// It returns the created execution with its assigned ID or an error if the operation fails.
func (s *ExecutionRepo) CreateExecution(ctx context.Context, execution models.Execution) (models.Execution, error) {
	timer := prometheus.NewTimer(executionLatency.WithLabelValues("create"))
	defer timer.ObserveDuration()

	result := s.db.Create(&execution)
	if result.Error != nil {
		executionOperations.WithLabelValues("create", "error").Inc()
		return models.Execution{}, fmt.Errorf("failed to create execution: %w", result.Error)
	}

	if execution.ID == 0 {
		executionOperations.WithLabelValues("create", "error").Inc()
		return models.Execution{}, fmt.Errorf("failed to get execution ID after creation")
	}

	executionOperations.WithLabelValues("create", "success").Inc()
	return execution, nil
}

// GetExecution retrieves a execution from the database by its ID.
// It returns a pointer to the execution if found, or an error if the execution doesn't exist or if the operation fails.
func (s *ExecutionRepo) GetExecution(ctx context.Context, executionID uint) (*models.Execution, error) {
	timer := prometheus.NewTimer(executionLatency.WithLabelValues("get"))
	defer timer.ObserveDuration()

	var execution models.Execution
	if err := s.db.First(&execution, executionID).Error; err != nil {
		executionOperations.WithLabelValues("get", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve execution by ID: %w", err)
	}
	executionOperations.WithLabelValues("get", "success").Inc()
	return &execution, nil
}

// ListExecutions retrieves a paginated list of executions from the database, filtered by status and type.
// The 'limit' parameter specifies the maximum number of executions to return,
// 'offset' determines the starting point for pagination,
// 'status' allows filtering by execution status, and 'executionType' allows filtering by execution type.
// It returns a slice of executions and an error if the operation fails.
func (s *ExecutionRepo) ListExecution(ctx context.Context) ([]models.Execution, error) {
	timer := prometheus.NewTimer(executionLatency.WithLabelValues("list"))
	defer timer.ObserveDuration()

	var executions []models.Execution

	// Execute the query
	if err := s.db.Find(&executions).Error; err != nil {
		executionOperations.WithLabelValues("list", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve executions: %w", err)
	}

	executionOperations.WithLabelValues("list", "success").Inc()
	return executions, nil
}

// NewExecutionRepo creates and returns a new instance of ExecutionRepo.
// It requires a GORM database connection and a River client for execution queue management.
func NewExecutionRepo(db *gorm.DB) interfaces.ExecutionRepo {
	return &ExecutionRepo{
		db: db,
	}
}
