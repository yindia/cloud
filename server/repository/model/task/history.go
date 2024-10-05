package task

import (
	"time"

	"gorm.io/gorm"
)

// TaskHistory represents the history of a task with status changes and additional details.
type TaskHistory struct {
	gorm.Model
	TaskID    uint      `json:"task_id" gorm:"not null"`         // Foreign key for Task
	Status    int       `json:"status" gorm:"type:int;not null"` // Refers to the custom PostgreSQL enum
	Details   string    `json:"details" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

// TableName returns the custom table name for the TaskHistory model.
func (*TaskHistory) TableName() string {
	return "task_histories"
}

// BeforeCreate sets the CreatedAt field to the current time before creating a new TaskHistory record.
func (th *TaskHistory) BeforeCreate(tx *gorm.DB) (err error) {
	th.CreatedAt = time.Now()
	return nil
}
