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
	workflowOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_repository_operations_total",
			Help: "The total number of workflow repository operations",
		},
		[]string{"operation", "status"},
	)
	workflowLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "workflow_repository_operation_duration_seconds",
			Help:    "Duration of workflow repository operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// WorkflowRepo implements the WorkflowRepo interface using GORM for database operations
// and River for workflow queue management.
type WorkflowRepo struct {
	db *gorm.DB
}

// CreateWorkflow creates a new workflow in the database and enqueues it for processing.
// It returns the created workflow with its assigned ID or an error if the operation fails.
func (s *WorkflowRepo) CreateWorkflow(ctx context.Context, workflow models.Workflow) (models.Workflow, error) {
	timer := prometheus.NewTimer(workflowLatency.WithLabelValues("create"))
	defer timer.ObserveDuration()

	result := s.db.Create(&workflow)
	if result.Error != nil {
		workflowOperations.WithLabelValues("create", "error").Inc()
		return models.Workflow{}, fmt.Errorf("failed to create workflow: %w", result.Error)
	}

	if workflow.ID == 0 {
		workflowOperations.WithLabelValues("create", "error").Inc()
		return models.Workflow{}, fmt.Errorf("failed to get workflow ID after creation")
	}

	workflowOperations.WithLabelValues("create", "success").Inc()
	return workflow, nil
}

// GetWorkflow retrieves a workflow from the database by its ID.
// It returns a pointer to the workflow if found, or an error if the workflow doesn't exist or if the operation fails.
func (s *WorkflowRepo) GetWorkflow(ctx context.Context, workflowID uint) (*models.Workflow, error) {
	timer := prometheus.NewTimer(workflowLatency.WithLabelValues("get"))
	defer timer.ObserveDuration()

	var workflow models.Workflow
	if err := s.db.First(&workflow, workflowID).Error; err != nil {
		workflowOperations.WithLabelValues("get", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve workflow by ID: %w", err)
	}
	workflowOperations.WithLabelValues("get", "success").Inc()
	return &workflow, nil
}

// ListWorkflows retrieves a paginated list of workflows from the database, filtered by status and type.
// The 'limit' parameter specifies the maximum number of workflows to return,
// 'offset' determines the starting point for pagination,
// 'status' allows filtering by workflow status, and 'workflowType' allows filtering by workflow type.
// It returns a slice of workflows and an error if the operation fails.
func (s *WorkflowRepo) ListWorkflow(ctx context.Context) ([]models.Workflow, error) {
	timer := prometheus.NewTimer(workflowLatency.WithLabelValues("list"))
	defer timer.ObserveDuration()

	var workflows []models.Workflow

	// Execute the query
	if err := s.db.Find(&workflows).Error; err != nil {
		workflowOperations.WithLabelValues("list", "error").Inc()
		return nil, fmt.Errorf("failed to retrieve workflows: %w", err)
	}

	workflowOperations.WithLabelValues("list", "success").Inc()
	return workflows, nil
}

// NewWorkflowRepo creates and returns a new instance of WorkflowRepo.
// It requires a GORM database connection and a River client for workflow queue management.
func NewWorkflowRepo(db *gorm.DB) interfaces.WorkflowRepo {
	return &WorkflowRepo{
		db: db,
	}
}
