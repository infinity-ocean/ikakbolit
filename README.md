# ikakbolit - Medication Scheduling App

## Описание
`ikakbolit` — это backend-сервис для планирования приема лекарств. Приложение использует PostgreSQL для хранения данных, `gorilla/mux` для маршрутизации API-запросов и `swaggo/swag` для генерации OpenAPI-документации.

## Стек технологий

### ⚡ Основные библиотеки и инструменты:
- **[gorilla/mux](https://github.com/gorilla/mux)** — маршрутизатор для обработки HTTP-запросов.
- **[PostgreSQL](https://www.postgresql.org)** — реляционная база данных для хранения информации о пользователях и расписаниях приема лекарств.
- **[pgx](https://github.com/jackc/pgx)** — драйвер для работы с PostgreSQL, обеспечивающий эффективную работу с БД.
- **[godotenv](https://github.com/joho/godotenv)** — загрузка переменных окружения из `.env`.
- **[swaggo/swag](https://github.com/swaggo/swag)** — генерация OpenAPI-документации.
- **[swaggo/http-swagger](https://github.com/swaggo/http-swagger)** — интеграция Swagger UI для удобного тестирования API.
- **[Docker](https://www.docker.com/)** Запуск сервиса и требуемой им инфраструктуры проводится в docker контейнерах.

## Инструкция по установке и запуску сервиса с помощью Docker🐋

### 1. Клонирование репозитория
```sh
git clone https://github.com/infinity-ocean/ikakbolit.git
cd ikakbolit
```

### 2. Настройка переменных окружения
Создай `.env` файл и укажи необходимые настройки для подключения к базе данных:
```ini
POSTGRES_HOST=localhost       # Хост базы данных
POSTGRES_PORT=5432            # Порт подключения к PostgreSQL
POSTGRES_USER=postgres        # Имя пользователя БД
POSTGRES_PASSWORD=12345       # Пароль пользователя БД
POSTGRES_DB=postgres          # Имя базы данных
POSTGRES_SSL=disable          # Режим SSL default=disable

DAY_START=08:00               # Время начала дня для расписания
DAY_FINISH=21:45              # Время окончания дня для расписания
CURE_SCHEDULE_WINDOW_MIN=120  # Период отображения расписаний лекарств (в минутах)
```

### 3. Установка зависимостей
```sh
go mod tidy
make install-deps
```

### 4. Запуск базы данных PostgreSQL в Docker
```sh
make start-infra
```

### 5. Применение миграций
```sh
make migration-up
```

### 6. Запуск приложения
```sh
go run internal/cmd/main.go
```

## API Документация
Для генерации OpenAPI-документации используй `swag`:
```sh
swag init -g internal/cmd/main.go
```
После запуска сервиса документация будет доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

## Кратко про существующие эндпоинты:
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

- **GET /schedules?user_id=**: Получает список ID расписаний для указанного пользователя. Требует параметр user_id в query.
- **GET /schedule?user_id=&schedule_id=**: Получает расписание приема лекарства по ID пользователя и ID расписания. Требует параметры user_id и schedule_id в query.
- **GET /next_takings?user_id**: Возвращает список ближайших приемов лекарств для пользователя. Требует параметр user_id в query.
