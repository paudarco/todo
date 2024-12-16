# Строим приложение на основе официального образа Golang
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Сборка исполняемого файла
RUN CGO_ENABLED=0 GOOS=linux go build -o todo ./cmd/main.go

# Создаем минимальный образ для запуска
FROM ubuntu:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исполняемый файл и папку web из предыдущего шага
COPY --from=builder /app/todo .
COPY --from=builder /app/web ./web

# Копируем .env файл в контейнер
COPY .env .env

# Указываем порт приложения
EXPOSE 7540

# Команда запуска
CMD ["./todo"]