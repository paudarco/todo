package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo/internal/errors"
	"github.com/paudarco/todo/internal/models"
)

// GetAllTasks возвращает в JSON все задачи из базы данных.
func (h *Handler) GetAllTasks(c *gin.Context) {
	// Так как эндпоинт для получения всех задач и поиска задачи один и тот же,
	// первым делом после запроса проверяем поле search
	// и если оно не пустое - вызываем соответствующий обработчик и останавливаем работу этого.
	search := c.Query("search")
	if len(search) > 0 {
		c.Set("search", search)
		h.SearchTasks(c)
		return
	}

	// Получаем массив с данными
	tasks, err := h.service.GetAllTasks(c)
	if err != nil {
		// Если ошибка существует в мапе сохраненных ошибок - отправляем ответ с кодом плохого запроса
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Если массив равен nil - отправляем ответ с пустым массивом
	if len(tasks) == 0 || tasks == nil {
		c.JSON(http.StatusOK, gin.H{
			"tasks": make([]map[string]string, 0),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

// GetTaskById получает одну задачу по идентификатору.
func (h *Handler) GetTaskById(c *gin.Context) {
	idQuery := c.Query("id")
	if len(idQuery) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrMissingId.Error())
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	task, err := h.service.GetTaskById(c, id)
	if err != nil {
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

// SearchTasks возвращает задачи по заданному параметру.
func (h *Handler) SearchTasks(c *gin.Context) {
	search := c.GetString("search")

	tasks, err := h.service.SearchTasks(c, search)
	if err != nil {
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

// CreateTask создает новую задачу в базе данных.
func (h *Handler) CreateTask(c *gin.Context) {
	var task models.TodoItem
	if err := c.ShouldBindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	if len(task.Title) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	if err := h.service.CreateTask(c, &task); err != nil {
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": task.Id,
	})
}

// UpdateTask обновляет задачу в базе данных.
func (h *Handler) UpdateTask(c *gin.Context) {
	// task - переменная для парсинга тела запроса.
	var task models.TodoItemResponse

	if err := c.ShouldBindJSON(&task); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	if len(task.Title) == 0 || len(task.Id) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidDataFormat.Error())
		return
	}

	if err := h.service.UpdateTask(c, task); err != nil {
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DoneTask помечает запись как выполненную.
func (h *Handler) DoneTask(c *gin.Context) {
	// Получаем id из url запроса и проверяем его корректность.
	idQuery := c.Query("id")
	if len(idQuery) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrMissingId.Error())
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidId.Error())
		return
	}

	if err = h.service.DoneTask(c, id); err != nil {
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteTask удаляет запись из базы данных.
func (h *Handler) DeleteTask(c *gin.Context) {
	// Получаем id из url запроса и проверяем его корректность.
	idQuery := c.Query("id")
	if len(idQuery) == 0 {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrMissingId.Error())
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidId.Error())
		return
	}

	if err = h.service.DeleteTask(c, id); err != nil {
		if _, ok := errors.ErrList[err]; ok {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
