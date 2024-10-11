package task

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Task represents a task with its attributes.
type Task struct {
	gorm.Model
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Type        string    `json:"type" gorm:"type:varchar(255);not null"` // Refers to the custom PostgreSQL enum
	Status      int       `json:"status" gorm:"type:int;not null"`        // Refers to the custom PostgreSQL enum
	Payload     string    `json:"payload" gorm:"type:jsonb;not null"`     // Storing JSON as a string in PostgreSQL
	Retries     int       `json:"retries" gorm:"default:0;check:retries >= 0 AND retries <= 10"`
	Priority    int       `json:"priority" gorm:"default:0;check:priority >= 0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime; not null"`
}

// TableName returns the custom table name for the Task model.
func (*Task) TableName() string {
	return "tasks"
}

// BeforeCreate ensures that the TaskType and TaskStatusEnum are valid before the task is created.
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	// Set CreatedAt to current time
	t.CreatedAt = time.Now()

	// Ensure task type is valid
	if t.Type != "send_email" && t.Type != "run_query" {
		return errors.New("invalid task type")
	}

	// Ensure task status is valid
	if t.Status > 4 {
		return errors.New("invalid task status")
	}

	// Ensure payload is valid JSON
	var js json.RawMessage
	if err := json.Unmarshal([]byte(t.Payload), &js); err != nil {
		return errors.New("invalid JSON payload")
	}

	return nil
}
