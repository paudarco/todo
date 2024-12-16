package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/paudarco/todo/internal/config"
	"github.com/paudarco/todo/internal/handler"
	"github.com/paudarco/todo/internal/repository"
	"github.com/paudarco/todo/internal/server"
	"github.com/paudarco/todo/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загружаем конфигурацию из .env файла
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Error initializing config: %v", err)
	}

	// Инициализируем базу данных SQLite
	db, err := repository.NewSQLiteDB(cfg.DB)
	if err != nil {
		logrus.Fatalf("Error initializing database: %v", err)
	}

	// Создаем слои репозитория, сервисов и обработчика
	repository := repository.NewRepository(db)
	service := service.NewService(repository, cfg)
	handler := handler.NewHandler(service, cfg)

	srv := new(server.Server)

	// Запускаем сервер
	go func() {
		if err = srv.Run(cfg.Server, handler.InitRouters()); err != nil {
			logrus.Fatalf("Error while running server: %v", err)
		}
	}()

	logrus.Println("Todo api started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Stopping todo api...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Printf("error while shutting down: %s\n", err.Error())
	}

	err = db.Close()
	if err != nil {
		logrus.Printf("error occured on db connection close: %s\n", err.Error())
	}
}
