package service

import (
	"context"
	"time"

	"github.com/paudarco/todo/internal/config"
	"github.com/paudarco/todo/internal/models"
	"github.com/paudarco/todo/internal/repository"
)

// Абстракция для работы с элементом из списка дел.
type TodoItem interface {
	GetAllTasks(ctx context.Context) ([]models.TodoItemResponse, error)
	GetTaskById(ctx context.Context, id int) (models.TodoItemResponse, error)
	SearchTasks(ctx context.Context, search string) ([]models.TodoItemResponse, error)
	CreateTask(ctx context.Context, task *models.TodoItem) error
	UpdateTask(ctx context.Context, task models.TodoItemResponse) error
	DoneTask(ctx context.Context, id int) error
	DeleteTask(ctx context.Context, id int) error
	NextDate(now time.Time, date string, repeat string) (string, error)
}

// Абстракция для работы с аутентификацией.
type Auth interface {
	SignIn(user models.User) (models.AuthResponse, error)
	ValidateToken(token string) error
}

// Кумулятивная структура для работы со всеми возможными абстракциями на уровне сервисов.
type Service struct {
	TodoItem
	Auth
}

func NewService(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		TodoItem: NewTodoItemService(repo),
		Auth:     NewAuthService(cfg),
	}
}
