# Car Estimator Backend

Единый backend-каталог для микросервисов Car Estimator.

Фронтенд находится отдельно в соседней папке `Car-Estimator` и не входит в этот compose.

## Состав

Активные сервисы в `docker-compose.yml`:

- `gateway` - HTTP API Gateway.
- `authorization` - gRPC-сервис авторизации и профиля.
- `feed` - gRPC-сервис объявлений.
- `predictor` - gRPC-сервис предсказания цены.
- `auth_db` - PostgreSQL для `authorization`.
- `feed_db` - PostgreSQL для `feed`.
- `redis` - хранилище сессий для `authorization`.

В каталоге также лежат `CarEstimator_Chat` и `CarEstimator_Purchase`, но они не подключены к compose, потому что сейчас в них нет Dockerfile и исполняемого сервиса.

## Структура

```text
car_estimator/
  docker-compose.yml
  .env.example
  car_estimator_api_gateway/
  car_estimator_authorization/
  car_estimator_api_contracts/
  CarEstimator_Feed/
  CarEstimator_PredictionPrice/
  CarEstimator_Chat/
  CarEstimator_Purchase/
```

## Быстрый запуск

```bash
cd /Users/alexmiami/Documents/GitHub/car_estimator
docker compose up -d
```

Compose использует дефолтные значения переменных, поэтому `.env` не обязателен для локального запуска.

## Настройка окружения

Если нужно поменять порты, пароли или имена баз, создай `.env` из примера:

```bash
cp .env.example .env
```

После этого отредактируй `.env`.

Основные переменные:

```env
GATEWAY_HTTP_PORT=4242
AUTH_GRPC_PORT=4444
PREDICTION_GRPC_PORT=50051
FEED_GRPC_PORT=50052

AUTH_POSTGRES_PORT=5432
FEED_POSTGRES_PORT=5433
REDIS_PORT=6379

SECRET_KEY=change-me
```

Для локальной разработки обязательно замени `SECRET_KEY`, если планируется тестировать авторизацию с реальными токенами.

## Порты

| Сервис | Внешний порт | Внутренний порт |
| --- | ---: | ---: |
| API Gateway | `4242` | `4242` |
| Authorization gRPC | `4444` | `4444` |
| Prediction gRPC | `50051` | `50051` |
| Feed gRPC | `50052` | `50052` |
| Auth PostgreSQL | `5432` | `5432` |
| Feed PostgreSQL | `5433` | `5432` |
| Redis | `6379` | `6379` |

## Команды

Собрать образы:

```bash
docker compose build
```

Запустить сервисы:

```bash
docker compose up -d
```

Посмотреть состояние:

```bash
docker compose ps
```

Посмотреть логи:

```bash
docker compose logs -f
```

Посмотреть логи конкретного сервиса:

```bash
docker compose logs -f gateway
docker compose logs -f authorization
docker compose logs -f feed
docker compose logs -f predictor
```

Остановить сервисы:

```bash
docker compose down
```

Остановить сервисы и удалить volumes с данными БД/Redis:

```bash
docker compose down -v
```

## API Gateway

Gateway доступен по адресу:

```text
http://localhost:4242
```

Он обращается к внутренним gRPC-сервисам по именам compose-сервисов:

```env
PROFILE_SERVICE_ADDR=authorization:4444
PREDICTION_SERVICE_ADDR=predictor:50051
FEED_SERVICE_ADDR=feed:50052
```

Эти адреса задаются в `docker-compose.yml` и не требуют ручной настройки при обычном локальном запуске.

## Данные

Данные хранятся в Docker volumes:

- `auth_pg_data` - данные PostgreSQL авторизации.
- `feed_pg_data` - данные PostgreSQL feed-сервиса.
- `redis_data` - данные Redis.

Чтобы полностью пересоздать окружение с чистыми базами:

```bash
docker compose down -v
docker compose up -d --build
```

## Проверка compose

Проверить итоговую конфигурацию:

```bash
docker compose config
```

Короткая проверка без вывода полной конфигурации:

```bash
docker compose config --quiet
```
