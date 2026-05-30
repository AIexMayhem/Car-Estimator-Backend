# Local Development

## Start Environment

```bash
docker compose up -d --build
```

## Check Containers

```bash
docker compose ps
docker compose logs --tail=100
```

## Rebuild One Service

```bash
docker compose build gateway
docker compose up -d gateway
```

Replace `gateway` with another compose service name when working on a different service.

## Reset Local Data

```bash
docker compose down -v
docker compose up -d --build
```

This removes database and Redis volumes, so use it only when local data can be recreated.
