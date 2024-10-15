/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	v1 "task/controller/api/v1"
	cloudv1 "task/pkg/gen/cloud/v1"
	cloudv1connect "task/pkg/gen/cloud/v1/cloudv1connect"
	"task/pkg/plugins"

	"connectrpc.com/connect"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TaskReconciler reconciles a Task object
type TaskReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	cloudClient cloudv1connect.TaskManagementServiceClient
}

// +kubebuilder:rbac:groups=task.io,resources=tasks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=task.io,resources=tasks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=task.io,resources=tasks/finalizers,verbs=update

// Reconcile is part of the main Kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// This function compares the state specified by the Task object against the
// actual cluster state, and then performs operations to make the cluster
// state reflect the state specified by the user.
func (r *TaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	task := &v1.Task{}
	err := r.Get(ctx, req.NamespacedName, task)
	if err != nil {
		log.FromContext(ctx).Error(err, "Failed to get task")
		return ctrl.Result{}, err
	}

	maxAttempts := 3
	initialBackoff := 1 * time.Second

	var finalStatus cloudv1.TaskStatusEnum
	var finalMessage string

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Update status to Running for each attempt
		runningMessage := fmt.Sprintf("Running attempt %d of %d", attempt, maxAttempts)
		if err := r.updateTaskStatus(ctx, int64(task.Spec.ID), cloudv1.TaskStatusEnum_RUNNING, runningMessage); err != nil {
			log.FromContext(ctx).Error(err, "Failed to update task status to Running")
			return ctrl.Result{}, err
		}

		_, message, err := processWorkflowUpdate(ctx, task)

		if err != nil {
			failedMessage := fmt.Sprintf("Attempt %d failed: %v", attempt, err)
			if err := r.updateTaskStatus(ctx, int64(task.Spec.ID), cloudv1.TaskStatusEnum_FAILED, failedMessage); err != nil {
				log.FromContext(ctx).Error(err, "Failed to update task status to Failed")
				return ctrl.Result{}, err
			}

			if attempt == maxAttempts {
				finalStatus = cloudv1.TaskStatusEnum_FAILED
				finalMessage = fmt.Sprintf("All %d attempts failed. Last error: %v", maxAttempts, err)
				log.FromContext(ctx).Error(fmt.Errorf(finalMessage), "Final failure after max attempts")
			} else {
				// Wait before the next attempt
				select {
				case <-ctx.Done():
					return ctrl.Result{}, ctx.Err()
				case <-time.After(initialBackoff * time.Duration(1<<uint(attempt-1))):
				}
				continue
			}
		} else {
			finalStatus = cloudv1.TaskStatusEnum_SUCCEEDED
			finalMessage = fmt.Sprintf("Task completed successfully on attempt %d: %s", attempt, message)
			log.FromContext(ctx).Info(finalMessage)
			break
		}
	}

	// Send final status update
	if err := r.updateTaskStatus(ctx, int64(task.Spec.ID), finalStatus, finalMessage); err != nil {
		log.FromContext(ctx).Error(err, "Failed to send final task status update")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Task{}).
		Complete(r)
}

// updateTaskStatus updates the status of a task using the Task Management Service.
func (r *TaskReconciler) updateTaskStatus(ctx context.Context, taskID int64, status cloudv1.TaskStatusEnum, message string) error {
	_, err := r.cloudClient.UpdateTaskStatus(ctx, connect.NewRequest(&cloudv1.UpdateTaskStatusRequest{
		Id:      int32(taskID),
		Status:  status,
		Message: message,
	}))
	if err != nil {
		return fmt.Errorf("failed to update task %d status: %w", taskID, err)
	}

	return nil
}

// processWorkflowUpdate handles different types of responses and returns the workflow state.
func processWorkflowUpdate(ctx context.Context, task *v1.Task) (cloudv1.TaskStatusEnum, string, error) {
	response := task

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.FromContext(ctx).Info(fmt.Sprintf("Workflow update processed in %s", duration))
	}()

	defer func() {
		if r := recover(); r != nil {
			// Return FAILED status in case of panic
			log.FromContext(ctx).Error(fmt.Errorf("Task panicked: %v", r), "Panic occurred during workflow update")
			panic(fmt.Sprintf("Task panicked: %v", r))
		}
	}()

	plugin, err := plugins.NewPlugin(response.Spec.Type)
	if err != nil {
		return cloudv1.TaskStatusEnum_FAILED, fmt.Sprintf("Failed to create plugin: %v", err), err
	}

	// Add retry logic for running the task
	runErr := plugin.Run(response.Spec.Payload.Parameters)
	if runErr != nil {
		return cloudv1.TaskStatusEnum_FAILED, fmt.Sprintf("Error running task: %v", runErr), runErr
	}

	return cloudv1.TaskStatusEnum_SUCCEEDED, "Task completed successfully", nil
}
