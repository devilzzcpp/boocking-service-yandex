# Бронирование

## О проекте

HTTP-сервис для управления бронированиями мест на PostgreSQL.

Стек проекта:

- Go
- Gin
- GORM
- PostgreSQL
- Viper
- Zap
- Docker Compose (для локальной БД)

Примечание: проект реализован в рамках Яндекс контестера в быстром темпе, примерно за 50 минут, с фокусом на выполнение требований ТЗ и прохождение проверок.

---

## Техническое задание (ТЗ)

### Условие

Вам дана база данных PostgreSQL с таблицей бронирований bookings.
Необходимо реализовать HTTP-сервис для управления бронированиями мест.

### Схема базы данных

```sql
CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    place_id INTEGER NOT NULL,
    time_from TIMESTAMP NOT NULL,
    time_to TIMESTAMP NOT NULL
);
```

### Запуск

Во всех случаях ожидается, что собранный бинарь будет принимать аргументы командной строки в формате:

```bash
--port <port>
```

где port — порт, на котором будет развернут сервис.

### API

#### GET /ping

Health check.
Должен вернуть HTTP 200 с телом:

```json
{"status":"ok"}
```

#### POST /book

Создание нового бронирования.

Query параметры:

- place_id (integer) — ID места
- user_id (integer) — ID пользователя
- from (string) — начало бронирования в формате RFC3339
- to (string) — конец бронирования в формате RFC3339

Ответ:

- HTTP 200 — бронирование успешно создано
- HTTP 409 — конфликт: запрашиваемый интервал пересекается с существующим бронированием для этого места

Два полуинтервала [from1, to1) и [from2, to2) (from1 < from2) считаются пересекающимися, если from2 < to1.

#### GET /booklist

Получение списка бронирований.

Query параметры (один из двух):

- user_id (integer) — вернуть все бронирования пользователя
- place_id (integer) — вернуть все бронирования места

Ответ HTTP 200 с телом:

```json
{
  "bookings": [
    {
      "id": 1,
      "user_id": 10,
      "place_id": 5,
      "from": "2024-01-01T10:00:00Z",
      "to": "2024-01-01T12:00:00Z"
    }
  ]
}
```

Бронирования должны быть отсортированы по (from, id) по возрастанию.

### Подключение к базе данных

Параметры подключения передаются через переменные окружения:

- DB_HOST — хост (по умолчанию localhost)
- DB_PORT — порт (по умолчанию 5432)
- DB_USER — пользователь (по умолчанию postgres)
- DB_PASSWORD — пароль (по умолчанию postgres)
- DB_NAME — имя базы данных (по умолчанию contest)

### Требования

Сервис должен корректно обрабатывать последовательные блокирующие запросы:
тестирующая система будет посылать следующий запрос только после получения ответа на предыдущий.

---

## Структура проекта

- main.go — входная точка приложения
- config/config.go — загрузка конфигурации из env
- utils/db.go — подключение и закрытие БД
- utils/logger.go — инициализация логгера
- internal/app/router.go — роутинг
- internal/app/middleware.go — middleware логирования
- internal/models/booking.go — модель данных
- internal/entity/booking/repository.go — слой репозитория
- internal/entity/booking/service.go — слой бизнес-логики
- internal/entity/booking/handler.go — HTTP-хендлеры
- migrations/ — SQL-миграции
- docker-compose.yml — локальный PostgreSQL для запуска/тестов

---

## Запуск проекта

### 1. Поднять PostgreSQL

```bash
docker compose up -d
```

### 2. Собрать бинарник

```bash
go build -o server .
```

### 3. Запустить сервис

Linux/macOS:

```bash
./server --port 8080
```

Windows:

```bash
go run main.go --port 8080
```

или

```bash
go build -o server.exe .
./server.exe --port 8080
```

---

## Проверка эндпоинтов

### Ping

```http
GET http://localhost:8080/ping
```

### Создание бронирования

```http
POST http://localhost:8080/book?place_id=1&user_id=10&from=2026-04-03T10:00:00Z&to=2026-04-03T11:00:00Z
```

### Получение списка

```http
GET http://localhost:8080/booklist?place_id=1
```

или

```http
GET http://localhost:8080/booklist?user_id=10
```