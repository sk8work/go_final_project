# Этап сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./main.go

# Финальный этап
FROM alpine:latest

# Создаем пользователя для безопасности
RUN addgroup -g 1000 todo && \
    adduser -D -u 1000 -G todo todo

# Создаем директорию для данных
RUN mkdir -p /data && chown todo:todo /data

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/todo-app .

# Копируем статические файлы
COPY --from=builder /app/web ./web

# Копируем пример конфигурации
COPY --from=builder /app/.env.example .env.example

# Меняем владельца
RUN chown -R todo:todo /app

# Переключаемся на непривилегированного пользователя
USER todo

# Открываем порт
EXPOSE 7540

# Переменные окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

# Точка входа
ENTRYPOINT ["./todo-app"]