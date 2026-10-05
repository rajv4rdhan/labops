# Squirrel 🐿️

A minimal, three-tier note-taking app. Each tier is deployed independently:

```
Browser ──> squirrel-frontend ──> squirrel-backend ──> PostgreSQL
            (static UI + proxy)   (Go REST API)       (official image)
```

There is **no login** — every note is shared.

## Tiers

| Tier     | Path                    | Image                              | Description                                   |
| -------- | ----------------------- | ---------------------------------- | --------------------------------------------- |
| Frontend | `app/squirrel/frontend` | `labops-squirrel-frontend`         | Static UI + optional `/api/` reverse proxy     |
| Backend  | `app/squirrel/backend`  | `labops-squirrel-backend`          | Go REST API, auto-migrating Postgres store     |
| Database | —                       | `postgres:16-alpine` (official)    | Not built here; pulled when you run the stack  |

The frontend is just static files served by a tiny Go web server. It proxies
`/api/` to the backend when `API_URL` is set (local/compose). Under an
Ingress you can instead route `/api` straight to the backend and leave
`API_URL` unset.

## API (backend)

| Method | Path              | Description         |
| ------ | ----------------- | ------------------- |
| GET    | `/api/notes`      | List notes          |
| POST   | `/api/notes`      | Create a note       |
| GET    | `/api/notes/{id}` | Get one note        |
| PUT    | `/api/notes/{id}` | Update a note       |
| DELETE | `/api/notes/{id}` | Delete a note       |
| GET    | `/health`         | Liveness            |
| GET    | `/ready`          | Readiness (DB ping) |

Note payload: `{"title": "...", "body": "..."}`.

## Run locally

Use the root `docker-compose.yml`, which starts the two app images plus the
official Postgres image:

```bash
docker compose up --build
# UI -> http://localhost:8082
```

Or run the tiers by hand against any Postgres:

```bash
# backend
cd backend
export DATABASE_URL="postgres://squirrel:squirrel@localhost:5433/squirrel?sslmode=disable"
go run .

# frontend (separate terminal)
cd frontend
export API_URL="http://localhost:8080"
go run .
```

## Configuration

Backend:

| Env            | Default | Description                                   |
| -------------- | ------- | --------------------------------------------- |
| `PORT`         | `8080`  | HTTP listen port                              |
| `DATABASE_URL` | —       | PostgreSQL connection string (required)       |

Frontend:

| Env       | Default | Description                                          |
| --------- | ------- | ---------------------------------------------------- |
| `PORT`    | `8080`  | HTTP listen port                                     |
| `API_URL` | —       | Optional backend base URL to proxy `/api/` to        |
