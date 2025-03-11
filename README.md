# ikakbolit - Medication Scheduling App 🏥

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.23.3-blue)
![Docker](https://img.shields.io/badge/Docker-Supported-blue)

**ikakbolit** — это бэкенд-сервис для удобного планирования приёма лекарств. Он позволяет пользователям управлять расписанием приёма медикаментов, отслеживать ближайшие приёмы и получать структурированные данные о назначенных лекарствах. В качестве БД используется **[PostgreSQL](https://www.postgresql.org)**.


## 🚀 Стек технологий

**Основные библиотеки и инструменты:**
- **[gorilla/mux](https://github.com/gorilla/mux)** — маршрутизатор для удобной работы с REST API.
- **[pgx](https://github.com/jackc/pgx)** — драйвер для PostgreSQL, оптимизированный для высокой производительности.
- **[godotenv](https://github.com/joho/godotenv)** — загрузка конфигурации из `.env`.
- **[swaggo/swag](https://github.com/swaggo/swag)** + **[http-swagger](https://github.com/swaggo/http-swagger)** — генерация и отображение OpenAPI-документации.
- **[Docker](https://www.docker.com/)** — упаковка сервиса и базы данных в контейнеры.

## 📥 Установка и запуск

### 1️⃣ Клонирование репозитория
```sh
git clone --depth=1 https://github.com/infinity-ocean/ikakbolit.git
cd ikakbolit
```

### 2️⃣ Настройка переменных окружения
Создай `.env` файл и укажи необходимые настройки для подключения к базе данных:

```ini
POSTGRES_HOST=localhost       # Хост базы данных
POSTGRES_PORT=5432            # Порт подключения к PostgreSQL
POSTGRES_USER=postgres        # Имя пользователя БД
POSTGRES_PASSWORD=12345       # Пароль пользователя БД
POSTGRES_DB=postgres          # Имя базы данных
POSTGRES_SSL=disable          # Отключение SSL (для локальной разработки)

DAY_START=08:00               # Время начала дня (ЧЧ:ММ, 24-часовой формат)
DAY_FINISH=21:45              # Время окончания дня
CURE_SCHEDULE_WINDOW_MIN=120  # Окно отображения расписаний (в минутах)
```

### 3️⃣ Установка зависимостей
```sh
go mod tidy
make install-deps
```

4️⃣ Запуск PostgreSQL в Docker контейнере
```sh
make start-infra
```

5️⃣ Применение миграций
```sh
make migration-up
```

6️⃣ Генерация OpenAPI-документации
```sh
swag init -g internal/cmd/main.go
```

7️⃣ Запуск приложения
```sh
go run internal/cmd/main.go
```

Приложение будет доступно по адресу: http://localhost:8080
Swagger-документация API: http://localhost:8080/swagger/index.html

## Основные эндпоинты:
- **POST /schedule**: Добавляет новое расписание приема лекарств для пользователя. Ожидает JSON с данными о приеме и сохраняет их в базе данных.
Ожидаемый JSON в теле запроса:
```json
{
  "user_id": 1,
  "cure_name": "Paracetamol",
  "doses_per_day": 3,
  "duration": 864000000000000,
}
```
Поле "duration" отображает длительность приёма, указывается в наносекундах. 1 день = 86400000000000 наносекунд.
Пример ответа:
```json
{
  "schedule_id": 42,
}
```
- **GET /schedules?user_id=**: Возвращает список ID расписаний для указанного пользователя. Требует параметр user_id в query.
Пример ответа:
```json
{
  "schedules": [1, 2, 3]
}
```
- **GET /schedule?user_id=&schedule_id=**: Получает расписание приема лекарства по ID пользователя и ID расписания. Требует параметры user_id и schedule_id в query.
- **GET /next_takings?user_id**: Возвращает список ближайших приемов лекарств для пользователя. Требует параметр user_id в query.
Пример ответа:
```json
{
  "schedules": [
    {
      "ID": 8,
      "user_id": "8",
      "cure_name": "Atorvastatin",
      "doses_per_day": 24,
      "duration": 31536000000000000,
      "created_at": "2025-03-07T23:46:58.292596Z",
      "intakes": [
        "08:00",
        "21:44"
      ]
    }
  ]
}
```

## 🛠 Полезные команды (Makefile)

| Команда              | Описание                                         |
|----------------------|--------------------------------------------------|
| `make install-deps`  | Устанавливает зависимости                        |
| `make start-infra`   | Запускает PostgreSQL в контейнере                |
| `make stop-infra`    | Останавливает контейнер PostgreSQL               |
| `make migration-up`  | Применяет миграции базы данных                   |
| `make migration-down`| Откатывает миграции                              |

### 🛠 Возможные улучшения
* [ ] Добавить авторизацию пользователей (JWT / OAuth) 
* [ ] Улучшить логирование запросов

## ❌🔧 Ошибка, с которой столкнулся
Сначала я хотел автоматизировать все команды используя Makefile, однако при автоматизации команды запуска приложения ```go run internal/cmd/main.go``` (через ```make run```) возникают ошибки с обработкой значений и .env файла. Поэтому для корректной работы мы вручную пишем ```go run internal/cmd/main.go```

Примеры ошибок:
```
{failed to calculate intake times: failed to parse DAY_START: parsing time "\"8:00\"" as "15:04": cannot parse "\"8:00\"" as "15"}
```
и
```
{failed to parse CURE_SCHEDULE_WINDOW_MIN: strconv.Atoi: parsing "\"120\"": invalid syntax}
```