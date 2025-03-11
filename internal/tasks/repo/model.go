package repo

import (
	tc "task-manager/internal/tasks_categories/repo"
	"time"
)

type Task struct {
	ID           int             `json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	IsCompleted  bool            `json:"is_completed"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	TaskCategory tc.TaskCategory `json:"task_category"`
}
