package task

import (
	"time"

	"gorm.io/gorm"
)

// Workflow represents a task with its attributes.
type Workflow struct {
	gorm.Model
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Payload     string    `json:"payload" gorm:"type:jsonb;not null"` // Storing JSON as a string in PostgreSQL
	Spec        []byte    `json:"spec" gorm:"type:bytea;not null"`    // Storing binary data in PostgreSQL
	Retries     int       `json:"retries" gorm:"default:0;check:retries >= 0 AND retries <= 10"`
	Priority    int       `json:"priority" gorm:"default:0;check:priority >= 0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime; not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime; not null"`
}

// TableName returns the custom table name for the Task model.
func (*Workflow) TableName() string {
	return "workflows"
}

// BeforeCreate ensures that the TaskType and TaskStatusEnum are valid before the task is created.
func (t *Workflow) BeforeCreate(tx *gorm.DB) (err error) {
	// Set CreatedAt to current time
	t.CreatedAt = time.Now()

	return nil
}
