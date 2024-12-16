package repository

import (
	"context"
	"database/sql"

	"github.com/paudarco/todo/internal/models"
)

// Абстракция для работы с элементом из списка дел.
type TodoItem interface {
	GetAllTasks(ctx context.Context) ([]models.TodoItem, error)
	GetTaskById(ctx context.Context, id int) (models.TodoItem, error)
	SearchTasksByDate(ctx context.Context, date string) ([]models.TodoItem, error)
	SearchTasksByString(ctx context.Context, search string) ([]models.TodoItem, error)
	CreateTask(ctx context.Context, task *models.TodoItem) error
	UpdateTask(ctx context.Context, task models.TodoItem) error
	DeleteTask(ctx context.Context, id int) error
}

// Кумулятивная структура для работы со всеми возможными абстракциями на уровне бд.
type Repository struct {
	TodoItem
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		TodoItem: NewTodoItemRepository(db),
	}
}
