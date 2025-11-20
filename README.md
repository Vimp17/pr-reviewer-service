# PR Reviewer Service

Микросервис для автоматического распределения ревьюверов для Pull Requests в командах разработки.

## Функциональность

- ✅ Создание команд с участниками
- ✅ Управление активностью пользователей  
- ✅ Создание Pull Requests с автоматическим назначением ревьюверов
- ✅ Слияние PR
- ✅ Перераспределение ревьюверов
- ✅ Получение статистики по назначениям
- ✅ Получение PR для конкретного ревьювера

## Технологии

- **Go ** - основной язык программирования
- **Gin** - веб-фреймворк
- **PostgreSQL** - база данных
- **Docker** - контейнеризация
- **Goose** - миграции базы данных

## Быстрый старт

### Запуск с Docker Compose

```bash
# Клонируйте репозиторий
git clone <your-repo-url>
cd pr-reviewer-service

# Запустите сервис
docker-compose up --build



API Endpoints
Команды
POST /team/add - Создать команду

GET /team/get?team_name=name - Получить информацию о команде

Пользователи
POST /users/setIsActive - Изменить активность пользователя

GET /users/getReview?user_id=id - Получить PR для ревьювера

Pull Requests
POST /pullRequest/create - Создать PR

POST /pullRequest/merge - Слить PR

POST /pullRequest/reassign - Перераспределить ревьювера

Системные
GET /health - Проверка здоровья

GET /stats - Статистика назначений

Примеры использования
Создание команды
bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend-team",
    "members": [
      {
        "user_id": "1",
        "username": "developer1",
        "is_active": true
      },
      {
        "user_id": "2", 
        "username": "developer2",
        "is_active": true
      }
    ]
  }'
Создание PR
bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1",
    "pull_request_name": "Add new feature",
    "author_id": "1"
  }'
Получение статистики
bash
curl http://localhost:8080/stats
Структура проекта
text
pr-reviewer-service/
├── cmd/
│   └── pr-reviewer/
│       └── main.go                 # Точка входа
├── internal/
│   ├── handlers/                   # HTTP обработчики
│   │   ├── handlers.go
│   │   ├── team_handlers.go
│   │   ├── user_handlers.go
│   │   └── pr_handlers.go
│   ├── models/                     # Модели данных
│   │   └── models.go
│   ├── services/                   # Бизнес-логика
│   │   ├── team_service.go
│   │   ├── user_service.go
│   │   └── pr_service.go
│   └── storage/                    # Работа с БД
│       ├── postgres/
│       │   ├── storage.go
│       │   ├── team_storage.go
│       │   ├── user_storage.go
│       │   └── pr_storage.go
├── migrations/                     # Миграции БД
│   ├── 000001_init_schema.up.sql
│   ├── 000002_create_users_table.up.sql
│   └── 000003_create_pull_requests_table.up.sql
├── config/                        # Конфигурация
├── Dockerfile
├── docker-compose.yml
└── README.md