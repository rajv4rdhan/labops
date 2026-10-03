# bunny-api

REST backend for the LabOps bunny frontend. It stores bunnies in PostgreSQL and
exposes a small JSON API.

## Endpoints

| Method | Path              | Description                     |
| ------ | ----------------- | ------------------------------- |
| GET    | `/health`         | Liveness probe                  |
| GET    | `/ready`          | Readiness probe (pings the DB)  |
| GET    | `/api/bunnies`    | List bunnies                    |
| POST   | `/api/bunnies`    | Create a bunny                  |
| GET    | `/api/bunnies/{id}` | Fetch a single bunny          |

## Configuration

| Env            | Default | Description                                  |
| -------------- | ------- | -------------------------------------------- |
| `PORT`         | `8080`  | HTTP listen port                             |
| `DATABASE_URL` | —       | PostgreSQL DSN (required)                    |

## Run locally

```bash
export DATABASE_URL="postgres://bunny:bunny@localhost:5432/bunnies?sslmode=disable"
go run .
```

## Build image

```bash
docker build -t bunny-api .
docker run -e DATABASE_URL="postgres://bunny:bunny@host:5432/bunnies?sslmode=disable" -p 8080:8080 bunny-api
```
