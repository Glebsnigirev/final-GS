package db

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"` // Опциональный комментарий
	Repeat  string `json:"repeat,omitempty"`  // Опциональное правило повторения
}
