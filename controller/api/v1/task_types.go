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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TaskSpec defines the desired state of Task
type TaskSpec struct {
	// ID is the unique identifier for the task.
	ID int32 `json:"id,omitempty"`

	// Name is the name of the task.
	Name string `json:"name,omitempty"`

	// Type is the type of the task.
	Type string `json:"type,omitempty"`

	// Status is the current status of the task.
	Status string `json:"status,omitempty"`

	// Retries is the number of retries attempted for this task.
	Retries int32 `json:"retries,omitempty"`

	// Priority is the priority level of the task.
	Priority int32 `json:"priority,omitempty"`

	// CreatedAt is the timestamp of when the task was created.
	CreatedAt string `json:"created_at,omitempty"`

	// Payload contains task parameters.
	Payload Payload `json:"payload,omitempty"`

	// Description is a description of the task.
	Description string `json:"description,omitempty"`
}

// Payload defines the parameters for the task.
type Payload struct {
	// Parameters are dynamic key-value pairs for task parameters.
	Parameters map[string]string `json:"parameters,omitempty"`
}

// TaskStatus defines the observed state of Task
type TaskStatus struct {
	// Status is the current status of the task.
	Status string `json:"status,omitempty"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Task is the Schema for the tasks API
type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskSpec   `json:"spec,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskList contains a list of Task
type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Task `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Task{}, &TaskList{})
}
