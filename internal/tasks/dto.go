package tasks

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	CategoryID  int    `json:"category_id"`
}
