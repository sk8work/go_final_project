# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# TODO Planner - Веб-приложение для планирования задач

Веб-приложение для управления задачами с поддержкой повторяющихся событий.

## Возможности

- Добавление, редактирование и удаление задач
- Поддержка повторяющихся задач (ежедневно, ежегодно)
- Поиск задач по заголовку, комментарию или дате
- Автоматический расчет следующих дат выполнения
- Веб-интерфейс
- Аутентификация (опционально)
- Docker-контейнеризация

## Задания со звёздочкой

- Поддержка переменных окружения для порта и пути к БД
- Поиск задач по заголовку/комментарию и дате
- Аутентификация через JWT
- Docker-образ

## Запуск локально

### Требования
- Go 1.24.9+
- SQLite

### Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd go_final_project
```

2. Установка зависимостей:
```bash
go mod download
```

3. Запустите сервер:
```bash
go run main.go
```

4. Откройте в браузере: http://localhost:7540


### Настройка аутентификации
Для включения аутентификации установите переменную окружения:
```bash
export TODO_PASSWORD=ваш_пароль
# или на Windows:
# set TODO_PASSWORD=ваш_пароль
```
После этого откройте: http://localhost:7540/login.html

### Запуск тестов
1. Убедитесь, что сервер не запущен на порту 7540
2. Запустите тесты:
```bash
# Все тесты
go test ./tests

# Конкретные тесты
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
```

## Docker

### Сборка образа
```bash
docker build -t todo-planner .
```

### Запуск контейнера
```bash
# Без аутентификации
docker run -p 7540:7540 -v $(pwd)/data:/data todo-planner

# С аутентификацией
docker run -p 7540:7540 -v $(pwd)/data:/data -e TODO_PASSWORD=secret123 todo-planner
```

### Docker Compose
```bash
# Запуск
docker-compose up -d

# Остановка
docker-compose down

# Просмотр логов
docker-compose logs -f
```

## API Endpoints
- GET /api/nextdate - расчет следующей даты
- POST /api/signin - аутентификация (если включена)
- GET /api/task - получение задачи
- POST /api/task - добавление задачи
- PUT /api/task - обновление задачи
- DELETE /api/task - удаление задачи
- POST /api/task/done - завершение задачи
- GET /api/tasks - список задач

## Переменные окружения
- TODO_PORT - порт сервера (по умолчанию 7540)
- TODO_DBFILE - путь к файлу БД (по умолчанию scheduler.db)
- TODO_PASSWORD - пароль для аутентификации (если пусто - аутентификация отключена)

## Лицензия
MIT
