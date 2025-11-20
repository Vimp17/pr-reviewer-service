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

- **Go** - основной язык программирования
- **Gin** - веб-фреймворк
- **PostgreSQL** - база данных
- **Docker** - контейнеризация
- **Goose** - миграции базы данных

## Быстрый старт

### Запуск с Docker Compose

```bash
# Клонируйте репозиторий
git clone https://github.com/Vimp17/pr-reviewer-service.git
cd pr-reviewer-service

# Запустите сервис
docker-compose up --build
```


## API Endpoints

### Команды (Teams)

| Метод | Endpoint | Описание |
|-------|----------|-----------|
| `POST` | `/team/add` | Создать новую команду с участниками |
| `GET` | `/team/get?team_name={name}` | Получить информацию о команде |

### Пользователи (Users)

| Метод | Endpoint | Описание |
|-------|----------|-----------|
| `POST` | `/users/setIsActive` | Изменить статус активности пользователя |
| `GET` | `/users/getReview?user_id={id}` | Получить список PR для ревьювера |

### Pull Requests

| Метод | Endpoint | Описание |
|-------|----------|-----------|
| `POST` | `/pullRequest/create` | Создать новый Pull Request |
| `POST` | `/pullRequest/merge` | Отметить PR как слитый |
| `POST` | `/pullRequest/reassign` | Перераспределить ревьювера |

### Системные (System)

| Метод | Endpoint | Описание |
|-------|----------|-----------|
| `GET` | `/health` | Проверка работоспособности сервиса |
| `GET` | `/stats` | Статистика по назначениям ревьюверов |

## Примеры использования

### Создание команды

```bash
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
