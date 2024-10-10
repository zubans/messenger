# Используем официальный образ Golang для сборки
FROM golang:1.19-alpine as builder

WORKDIR /app

# Установка необходимых инструментов
RUN apk --no-cache add git

# Копируем go.mod и инициализируем зависимости
COPY go.mod .
RUN go mod tidy

# Копируем весь исходный код и загружаем зависимости
COPY . .
RUN go get ./...        # Загрузка и установка всех зависимостей
RUN go build -o video-conference ./cmd/server

# Используем второй stage для минимального размера финального образа
# FROM alpine:latest

# WORKDIR /app
# COPY --from=builder /app/video-conference .

# ENTRYPOINT ["/app/video-conference"]