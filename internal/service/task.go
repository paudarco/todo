package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/paudarco/todo/internal/errors"
	"github.com/paudarco/todo/internal/models"
	"github.com/paudarco/todo/internal/repository"
)

// Структура для сервиса, работающего с делами из списка
type TodoItemService struct {
	repository *repository.Repository
}

func NewTodoItemService(repo *repository.Repository) *TodoItemService {
	return &TodoItemService{repository: repo}
}

// GetAllTasks получает все задачи из репозитория и возвращает их в нужном формате
func (s *TodoItemService) GetAllTasks(ctx context.Context) ([]models.TodoItemResponse, error) {
	// Получаем все задачи из репозитория
	tasks, err := s.repository.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	// Преобразуем результаты в нужный формат
	tasksParsed := make([]models.TodoItemResponse, 0)
	now := time.Now().Format(models.Format)

	// Тесты требуют id в формате строки, а не числа,
	// поэтому меняем формат вывода
	for _, task := range tasks {
		if task.Date >= now {
			tasksParsed = append(tasksParsed, models.TodoItemResponse{
				Id:      fmt.Sprintf("%d", task.Id),
				Date:    task.Date,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			})
		}
	}

	// Сортируем задачи по дате выполнения
	sort.Slice(tasksParsed, func(i, j int) bool {
		return tasksParsed[i].Date < tasksParsed[j].Date
	})

	return tasksParsed, nil
}

// GetTaskById получает задачу по её идентификатору и возвращает её в нужном формате
func (s *TodoItemService) GetTaskById(ctx context.Context, id int) (models.TodoItemResponse, error) {
	var taskParsed models.TodoItemResponse

	// Получаем задачу из репозитория
	task, err := s.repository.GetTaskById(ctx, id)
	if err != nil {
		return taskParsed, err
	}

	// Преобразуем результаты в нужный формат
	taskParsed = models.TodoItemResponse{
		Id:      fmt.Sprintf("%d", task.Id),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	return taskParsed, nil
}

// SearchTasks возвращает список задач по поисковому запросу.
func (s *TodoItemService) SearchTasks(ctx context.Context, search string) ([]models.TodoItemResponse, error) {
	var tasksResponse []models.TodoItemResponse
	var tasks []models.TodoItem

	// Проверяем, является ли строка поиска датой.
	date, err := time.Parse("02.01.2006", search)
	if err == nil {
		// Если да, ищем задачи с указанной датой выполнения.
		dateFormatted := date.Format(models.Format)
		tasks, err = s.repository.SearchTasksByDate(ctx, dateFormatted)
		if err != nil {
			return nil, err
		}
	} else {
		// Если нет, ищем задачи содержащие указанный текст в заголовке или комментарии.
		tasks, err = s.repository.SearchTasksByString(ctx, search)
		if err != nil {
			return nil, err
		}
	}

	// Преобразуем результаты в нужный формат
	for _, task := range tasks {
		tasksResponse = append(tasksResponse, models.TodoItemResponse{
			Id:      fmt.Sprintf("%d", task.Id),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}

	return tasksResponse, nil
}

// Create создает новую запись в базе данных
func (s *TodoItemService) CreateTask(ctx context.Context, task *models.TodoItem) error {
	// Сначала проверяем корректность указанного правила повторения
	if len(task.Repeat) != 0 {
		_, _, err := parseRule(task.Repeat)
		if err != nil {
			return err
		}
	}

	// date - дата выполнения задачи
	var date string
	now := time.Now()

	if len(task.Date) != 0 {
		// Проверяем, является ли введенная дата корректной
		_, err := time.Parse(models.Format, task.Date)
		if err != nil {
			return errors.ErrInvalidDate
		}

		// Если дата выполнения позднее текущей, получаем следующую дату повторения
		// согласно правилу, если оно указано
		if task.Date < now.Format(models.Format) {
			if len(task.Repeat) != 0 {
				date, err = s.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return err
				}
			} else {
				date = now.Format(models.Format)
			}
		} else {
			date = task.Date
		}
	} else {
		date = now.Format(models.Format)
	}

	task.Date = date

	// Записываем новую задачу в базу данных
	err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTask изменяет данные задачи в базе данных
func (s *TodoItemService) UpdateTask(ctx context.Context, task models.TodoItemResponse) error {
	id, err := strconv.Atoi(task.Id)
	if err != nil {
		return errors.ErrInvalidId
	}

	var date string
	now := time.Now()

	if len(task.Date) != 0 {
		// Проверяем, является ли введенная дата корректной
		_, err := time.Parse(models.Format, task.Date)
		if err != nil {
			return errors.ErrInvalidDate
		}

		// Если дата выполнения позднее текущей, получаем следующую дату повторения
		// согласно правилу, если оно указано
		if task.Date < now.Format(models.Format) {
			if len(task.Repeat) != 0 {
				date, err = s.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return err
				}
			} else {
				date = now.Format(models.Format)
			}
		} else {
			date = task.Date
		}
	} else {
		date = now.Format(models.Format)
	}

	// Создаем новую модель задачи с измененными данными  в нужном формате
	taskParsed := models.TodoItem{
		Id:      id,
		Date:    date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	// Обновляем задачу в базе данных.
	err = s.repository.UpdateTask(ctx, taskParsed)
	if err != nil {
		return err
	}

	return nil
}

// DoneTask отмечает задачу выполненной и, если не указано правило повторения,
// удаляет ее из базы данных, в противном случае - меняет дату события на следующую.
func (s *TodoItemService) DoneTask(ctx context.Context, id int) error {
	// Получаем задачу из бд для проверки правила повторения
	task, err := s.repository.GetTaskById(ctx, id)
	if err != nil {
		return err
	}

	// Если правило повторения указано, меняем дату выполнения на следующую.
	if len(task.Repeat) != 0 {
		date, err := s.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = date

		// Обновляем задачу с новой датой выполнения.
		err = s.repository.UpdateTask(ctx, task)
		if err != nil {
			return err
		}

		return nil
	}

	// Если правило повторения не указано, удаляем задачу из базы данных.
	err = s.repository.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask удаляет задачу из базы данных.
func (s *TodoItemService) DeleteTask(ctx context.Context, id int) error {
	return s.repository.DeleteTask(ctx, id)
}
