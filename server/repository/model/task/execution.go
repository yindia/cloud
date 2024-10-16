package task

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Execution represents a task with its attributes.
type Execution struct {
	gorm.Model
	TaskID    uint      `json:"task_id" gorm:"not null"`            // Foreign key for Task
	Status    int       `json:"status" gorm:"type:int;not null"`    // Refers to the custom PostgreSQL enum
	Payload   string    `json:"payload" gorm:"type:jsonb;not null"` // Storing JSON as a string in PostgreSQL
	Retries   int       `json:"retries" gorm:"default:0;check:retries >= 0 AND retries <= 10"`
	Priority  int       `json:"priority" gorm:"default:0;check:priority >= 0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime; not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime; not null"`
}

// TableName returns the custom table name for the Task model.
func (*Execution) TableName() string {
	return "executions"
}

// BeforeCreate ensures that the TaskType and TaskStatusEnum are valid before the task is created.
func (t *Execution) BeforeCreate(tx *gorm.DB) (err error) {
	// Set CreatedAt to current time
	t.CreatedAt = time.Now()

	// Ensure task status is valid
	if t.Status > 4 {
		return errors.New("invalid task status")
	}

	return nil
}

// BeforeCreate ensures that the TaskType and TaskStatusEnum are valid before the task is created.
func (t *Execution) BeforeUpdate(tx *gorm.DB) (err error) {
	// Set CreatedAt to current time
	t.UpdatedAt = time.Now()

	// Ensure task status is valid
	if t.Status > 4 {
		return errors.New("invalid task status")
	}

	return nil
}
