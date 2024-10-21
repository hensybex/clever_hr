# Clever HR App

Этот репозиторий содержит Clever HR App, который включает Flutter веб-приложение и Go API, а также базы данных PostgreSQL и Milvus.## Предварительные требования- Docker
- Docker Compose

## Установка

1. **Клонируйте репозиторий:**

   ```bash
   git clone https://github.com/your-username/clever_hr.git
   cd clever_hr
   ```

2. **Создайте файл `.env`:**

   Скопируйте файл `.env.example` в `.env` и обновите переменные окружения по необходимости.

   ```bash
   cp .env.example .env
   ```

3. **Соберите и запустите приложение:**

   Используйте Docker Compose для сборки и запуска приложения.

   ```bash
   docker-compose -f docker-compose.dev.yml up --build
   ```

## Переменные окружения

Файл `.env` должен содержать следующие переменные окружения:

```env
POSTGRES_DB=\"clever_hr_db\"
POSTGRES_USER=\"postgres\"
POSTGRES_PASSWORD=\"postgres_password\"
POSTGRES_HOST=\"localhost\"
POSTGRES_PORT=\"5432\"
MISTRAL_API_KEY=\"mistral_api_key\"
POSTGRES_SSL_MODE=\"disable\"
API_BASE_URL=http://api:8080/api
```

## Docker Compose файлы

### `docker-compose.dev.yml`

Этот файл используется для разработки. Он включает следующие сервисы:

- База данных PostgreSQL
- База данных Milvus (с etcd и minio)
- Go API (собрана из `Dockerfile.dev`)

### `docker-compose.yml`

Этот файл предназначен для использования в продакшене. Он включает следующие сервисы:

- База данных PostgreSQL
- База данных Milvus (с etcd и minio)
- Go API (собрана из `Dockerfile`)
