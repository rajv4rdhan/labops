# Squirrel 🐿️

A minimal note-taking app. A single Go binary serves both a small REST API
and an embedded web UI, backed by PostgreSQL. There is **no login** — every
note is shared.

```
Browser ──> squirrel (Go: REST API + embedded UI) ──> PostgreSQL
```

## Features

- Create, read, update and delete notes
- Embedded, dependency-free frontend (no build step, no CDN at runtime)
- Auto-migrating schema on startup (seeds one welcome note)
- Health (`/health`) and readiness (`/ready`) endpoints
- Graceful shutdown

## API

| Method | Path             | Description        |
| ------ | ---------------- | ------------------ |
| GET    | `/api/notes`     | List notes         |
| POST   | `/api/notes`     | Create a note      |
| GET    | `/api/notes/{id}`| Get one note       |
| PUT    | `/api/notes/{id}`| Update a note      |
| DELETE | `/api/notes/{id}`| Delete a note      |
| GET    | `/health`        | Liveness           |
| GET    | `/ready`         | Readiness (DB ping)|

Note payload: `{"title": "...", "body": "..."}`.

## Run locally

Needs a PostgreSQL instance and `DATABASE_URL`:

```bash
export DATABASE_URL="postgres://squirrel:squirrel@localhost:5432/squirrel?sslmode=disable"
go run .
# open http://localhost:8080
```

Or use the root `docker-compose.yml` to bring up Squirrel and Postgres
together.

## Configuration

| Env            | Default | Description                |
| -------------- | ------- | -------------------------- |
| `PORT`         | `8080`  | HTTP listen port           |
| `DATABASE_URL` | —       | PostgreSQL connection string (required) |

## Project structure

```text
squirrel/
├── main.go        # startup, graceful shutdown
├── handlers.go    # HTTP router + REST handlers, embeds web/
├── store.go       # pgx pool, migrations, note CRUD
├── web/           # index.html, styles.css, app.js (go:embed)
├── Dockerfile
└── README.md
```
