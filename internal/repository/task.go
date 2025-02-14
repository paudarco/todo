package repository

import (
	"context"
	"database/sql"

	"github.com/paudarco/todo/internal/errors"
	"github.com/paudarco/todo/internal/models"
)

// Структура для репозитория, работающего с делами из списка дел.
type TodoItemRepository struct {
	db *sql.DB
}

func NewTodoItemRepository(db *sql.DB) *TodoItemRepository {
	return &TodoItemRepository{db: db}
}

// GetAllTasks возващает слайс со всеми задачами.
func (repo *TodoItemRepository) GetAllTasks(ctx context.Context) ([]models.TodoItem, error) {
	query := `SELECT * FROM scheduler`

	// получаем все строки с результатом.
	rows, error := repo.db.QueryContext(ctx, query)
	if error != nil {
		return nil, error
	}
	defer rows.Close()

	// Результирующий слайс.
	var tasks []models.TodoItem

	for rows.Next() {
		var task models.TodoItem
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskById возвращает одну задачу по идентификатору.
func (repo *TodoItemRepository) GetTaskById(ctx context.Context, id int) (models.TodoItem, error) {
	query := `SELECT * FROM scheduler WHERE id=?`

	var task models.TodoItem

	// получаем одну строку с результатом.
	err := repo.db.QueryRowContext(ctx, query, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	// Если задача не найдена по идентификатору - возвращаем ошибку "Задача не найдена".
	if err == sql.ErrNoRows {
		return task, errors.ErrTaskNotFound
	} else if err != nil {
		return task, err
	}

	return task, nil
}

// SearchTasksByDate возвращает все задачи по указанной дате.
func (repo *TodoItemRepository) SearchTasksByDate(ctx context.Context, date string) ([]models.TodoItem, error) {
	query := `SELECT * FROM scheduler WHERE DATE = ?`

	rows, err := repo.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.TodoItem
	for rows.Next() {
		var t models.TodoItem
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// SearchTasksByString возвращает все задачи, включающие в себя указанную подстроку.
func (repo *TodoItemRepository) SearchTasksByString(ctx context.Context, search string) ([]models.TodoItem, error) {
	query := `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search`

	rows, err := repo.db.QueryContext(ctx, query, sql.Named("search", "%"+search+"%"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.TodoItem
	for rows.Next() {
		var t models.TodoItem
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Create создает задачу в базе данных.
func (repo *TodoItemRepository) CreateTask(ctx context.Context, task *models.TodoItem) error {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4) RETURNING id`

	return repo.db.QueryRowContext(ctx, query, task.Date, task.Title, task.Comment, task.Repeat).Scan(&task.Id)
}

// UpdateTask изменяет задачу в базе данных.
func (repo *TodoItemRepository) UpdateTask(ctx context.Context, task models.TodoItem) error {
	query := `UPDATE scheduler SET date=$1, title=$2, comment=$3, repeat=$4 WHERE id=$5 `

	result, err := repo.db.ExecContext(ctx, query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return err
	}

	// Проверяем, изменилась ли запись.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.ErrTaskNotFound
	}

	return nil
}

// DeleteTask удаляет задачу из базы данных по идентификатору.
func (repo *TodoItemRepository) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM scheduler WHERE id=$1`

	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Проверяем, удалилась ли запись.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.ErrTaskNotFound
	}

	return nil
}
