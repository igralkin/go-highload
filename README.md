# go-highload — микросервис на Go

Учебный высоконагруженный микросервис на Go, реализующий CRUD по пользователям, конкурентное логирование и уведомления, интеграцию с MinIO, метрики Prometheus и rate limiting. Сервис контейнеризован в Docker и может быть протестирован через `wrk`.

## Функциональность

- CRUD API над сущностью `User`
- Конкурентное логирование действий (goroutines + каналы)
- Асинхронные уведомления о событиях
- Rate limiting (1000 rps, burst 5000)
- Метрики Prometheus (`/metrics`)
- Интеграция с MinIO (сохранение дампа пользователей, вывод списка объектов)
- Полная контейнеризация через Docker + docker-compose
- Нагрузочное тестирование wrk

## Структура проекта

```
go-highload/
  main.go
  models/
  handlers/
  services/
  utils/
  metrics/
  Dockerfile
  docker-compose.yml
```

### Описание директорий

| Папка | Содержание |
|-------|------------|
| `models/` | доменные модели (`User`) |
| `services/` | бизнес-логика: CRUD, MinIO, уведомления |
| `handlers/` | HTTP-эндпоинты (CRUD, интеграция) |
| `utils/` | rate limiter, логгер, вспомогательные функции |
| `metrics/` | Prometheus метрики и middleware |

## Запуск через Docker

### Сборка и запуск

```bash
docker compose build
docker compose up -d
```

### Проверка работы

```bash
curl http://localhost:8080/api/users
```

### Панель MinIO

```
http://localhost:9001
login:    minioadmin
password: minioadmin
```

## CRUD API (curl-примеры)

### Создать пользователя
```bash
curl -X POST http://localhost:8080/api/users   -H "Content-Type: application/json"   -d '{"name":"Alice","email":"alice@example.com"}'
```

### Получить всех
```bash
curl http://localhost:8080/api/users
```

### Получить по ID
```bash
curl http://localhost:8080/api/users/1
```

### Обновить
```bash
curl -X PUT http://localhost:8080/api/users/1   -H "Content-Type: application/json"   -d '{"name":"Alice Updated","email":"alice2@example.com"}'
```

### Удалить
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## Интеграция с MinIO

### Сохранить пользователей в MinIO

```bash
curl -X POST http://localhost:8080/api/integration/save-users
```

### Получить список объектов

```bash
curl http://localhost:8080/api/integration/list-objects
```

## Метрики Prometheus

```
GET /metrics
```

Пример:

```bash
curl http://localhost:8080/metrics | head
```

## Rate Limiting

- limit: **1000 req/s**
- burst: **5000**

При превышении — HTTP 429.

## Нагрузочное тестирование wrk

```bash
wrk -t12 -c500 -d60s http://localhost:8080/api/users
```

Примерный результат:

```
Requests/sec: ~29298
Latency avg: ~17ms
```

## Переменные окружения

```
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=go-highload
MINIO_USE_SSL=false
```

## Локальный запуск

```bash
go mod tidy
go run .
```

## Репозиторий

https://github.com/igralkin/go-highload
