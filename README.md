# ikakbolit - Medication Scheduling App 🏥🦁

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.23.8-blue)
![Docker](https://img.shields.io/badge/Docker-Supported-blue)

**ikakbolit** — это бэкенд-сервис для удобного планирования приёма лекарств. Он позволяет пользователям управлять расписанием приёма медикаментов, отслеживать ближайшие приёмы и получать данные о назначенных лекарствах. В качестве БД используется **[PostgreSQL](https://www.postgresql.org)**
# P.S.
Привет! Спасибо что заглянули в проект :)
Ветка `via_dockerfile` содержит вариант запуска сервиса полностью через Docker. Ветка `main` использует `go run` для запуска приложения и контейнер для PostgreSQL.  
Тесты работают в обеих версиях, но в `via_dockerfile` пока некорректно отображается покрытие — проблема не решена.
# 📥 Установка и запуск

### 1️⃣ Клонирование репозитория
```sh
git clone --depth=1 https://github.com/infinity-ocean/ikakbolit.git
cd ikakbolit
```

### 2️⃣ Настройка переменных окружения
Создай `.env` файл и укажи необходимые настройки для подключения к базе данных (см. `.env.example`)
### 3️⃣ Установка зависимостей
```sh
go mod download
make install-deps
```
### 4️⃣ Запуск PostgreSQL в Docker контейнере
```sh
make start-infra
```
### 5️⃣ Применение миграций
```sh
make goose-up
```
### 6️⃣ Запуск приложения из корневой папки (для корректной работы с .env файлом)
```sh
make run
```

## Запуск тестов
```sh
make test
make unit-test # Юнит тесты 
```
## Линтер
```sh
make lint
```
Приложение будет доступно по адресу: http://localhost:8080 (См. ```HTTP_PORT``` в .env)

Swagger-документация API: http://localhost:8080/swagger/index.html (Неактуальная)

## 🛠 Полезные команды (Makefile)

| Команда             | Описание                                                                                   |
| ------------------- | ------------------------------------------------------------------------------------------ |
| `make install-deps` | Устанавливает CLI-инструменты: Goose, Swagger, protoc-gen-go, protoc-gen-go-grpc и Nilaway |
| `make run`          | Запускает приложение: выполняет `go run cmd/ikakbolit/main.go`                             |
| `make test`         | Полный прогон тестов с перекаткой миграций:                                                |
|                     |                                                                                            |
| `make unit-test`    | Запускает быстрые юнит-тесты (`go test -short ./...`)                                      |
| `make start-infra`  | Поднимает инфраструктуру: запускает контейнер БД через `docker-compose up -d`              |
| `make lint`         | Запускает `golangci-lint run --config .golangci.yml`                                       |
| `make stop-infra`   | Останавливает инфраструктуру: останавливает контейнеры через `docker-compose down`         |
| `make goose-up`     | Применяет миграции                                                                         |
| `make goose-down`   | Откатывает миграции                                                                        |
| `make swagger-gen`  | Генерирует серверный код по OpenAPI-спецификации:                                          |
| `make proto-gen`    | Генерирует Go и gRPC код по `.proto` файлу                                                 |

# Заготовленные запросы для Postman <img src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/postman/postman-original.svg" alt="Postman" width="40" height="40"/>
- **1️⃣ POST http://localhost:8080/schedule** 
*Прикрепляем подобный JSON в Body*
```json
{
  "user_id": 1,
  "cure_name": "Arbidol",
  "doses_per_day": 7,
  "duration_days": 5 // допускается 0, обозначает постоянный приём препарата
}
```
- **2️⃣ GET http://localhost:8080/schedules?user_id=5**
- **3️⃣ GET http://localhost:8080/schedule?user_id=6&schedule_id=6**
- **4️⃣ GET http://localhost:8080/next_takings?user_id=5**

### Генерация OpenAPI-документации
```sh
swag init -g cmd/ikakbolit/main.go
```

## 🚀 Стек технологий

**Основные библиотеки и инструменты:**
- **[go-chi/chi](https://github.com/go-chi/chi)** — маршрутизатор для удобной работы с REST API.
- **[pgx](https://github.com/jackc/pgx)** — драйвер для PostgreSQL, оптимизированный для высокой производительности.
- **[godotenv](https://github.com/joho/godotenv)** — загрузка конфигурации из `.env`.
- **[swaggo/swag](https://github.com/swaggo/swag)** + **[http-swagger](https://github.com/swaggo/http-swagger)** — генерация и отображение OpenAPI-документации.
- **[goose](https://github.com/pressly/goose)** — инструмент для управления миграциями базы данных.
- **[Docker](https://www.docker.com/)** — упаковка базы данных в контейнер.

## 🏛️ Архитектура
Проект основан на трёхслойной архитектуре.
### 1. Presentation Layer (Controller)
### 2. Business Logic Layer (Service)
### 3. Data Access Layer (Repo)

### Поток данных:
1. **HTTP-запрос** → Controller (валидация) → Service (логика) → Repo (БД)
2. **Ответ**: Repo → Service → Controller → HTTP-клиент

### Принципы:
- **Разделение ответственности**: Каждый слой решает свою задачу.
- **Инверсия зависимостей**: Слои зависят от интерфейсов (не от конкретной реализации).

## Информация по эндпоинтам:
- **1️⃣ POST /schedule**: Добавляет новое расписание приема лекарств для пользователя. Ожидает JSON с данными о приеме и сохраняет их в базе данных.
Ожидаемый JSON в теле запроса:
```json
{
  "user_id": 1,
  "cure_name": "Arbidol",
  "doses_per_day": 7,
  "duration_days": 5 // допускается 0, обозначает постоянный приём препарата
}
```
Поле ```"duration_days"``` отображает длительность приёма, указывается в днях, ```"doses_per_day"``` количество ежедневных приёмов.
Пример ответа:
```json
{
  "schedule_id": 42,
}
```
- **2️⃣ GET /schedules?user_id=**: Возвращает список ID расписаний для указанного пользователя. Требует параметр user_id в query.
Пример ответа:
```json
{
  "schedules": [1, 2, 3]
}
```
- **3️⃣ GET /schedule?user_id=&schedule_id=**: Получает расписание приема лекарства по ID пользователя и ID расписания. Требует параметры user_id и schedule_id в query. Пример ответа:
```json
{
    "ID": 4,
    "UserID": 4,
    "CureName": "Amoxicillin",
    "DosesPerDay": 3,
    "DurationDays": 12,
    "CreatedAt": "2025-03-12T17:48:42.645442Z",
    "Intakes": [
        "08:00",
        "14:52",
        "21:45"
    ]
}
```
- **4️⃣ GET /next_takings?user_id**: Возвращает список ближайших приемов лекарств для пользователя в период, определяемый в конфигурации переменной ```"CURE_SCHEDULE_WINDOW_MIN"```. Требует параметр user_id в query.
Пример ответа:
```json
{
  "schedules": [
    {
      "ID": 8,
      "user_id": "8",
      "cure_name": "Atorvastatin",
      "doses_per_day": 2,
      "duration_days": 3,
      "created_at": "2025-03-07T23:46:58.292596Z",
      "intakes": [
        "08:00",
        "21:44"
      ]
    }
  ]
}
```
