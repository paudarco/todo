package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/paudarco/todo/internal/config"
	"github.com/paudarco/todo/internal/service"
)

// Handler - структура для работы с сервисами через обработчики.
type Handler struct {
	service *service.Service
	cfg     *config.Config
}

func NewHandler(service *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

// Настраиваем обработчик
func (h *Handler) InitRouters() *gin.Engine {
	r := gin.Default()

	// Настраиваем файл-сервер.
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.StaticFile("/favicon.ico", "./web/favicon.ico")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/login", "./web/login.html")

	// Сайт работает с вышеперечисленными путями,
	// но тесты почему-то принимают html только в виде таких двух строк.
	r.StaticFile("/index.html", "./web/index.html")
	r.StaticFile("/login.html", "./web/login.html")

	// Публичные эндпоинты, не требующие аутентификации.
	public := r.Group("/api")
	{
		public.POST("/signin", h.SignIn)
		public.GET("/nextdate", h.NextDate)
	}

	// Защищенные эндпоинты, требующие аутентификации.
	protected := r.Group("/api", h.AuthMiddleware)
	{
		tasks := protected.Group("/tasks")
		{
			tasks.GET("/", h.GetAllTasks)
			tasks.GET("/search", h.SearchTasks)
		}

		task := protected.Group("/task")
		{
			task.GET("/", h.GetTaskById)
			task.POST("/", h.CreateTask)
			task.PUT("/", h.UpdateTask)
			task.POST("/done", h.DoneTask)
			task.DELETE("/", h.DeleteTask)
		}
	}

	return r
}
