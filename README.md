# LabOps

A hands-on DevOps lab: a small multi-tier application used to practice
containers, Kubernetes, CI/CD, and beyond.

```
Browser ──┬── bunny-page (Go static server, cute meadow page)
          │
          └── Squirrel ──> postgres      (Go REST API + embedded UI)
```

The `bunny-api` REST service lives in `app/bunny-api` and is used by the
Kubernetes labs (05+); it is not needed to run the pages.

## Components

| Path               | What it is                                          |
| ------------------ | --------------------------------------------------- |
| `app/BunnyPage`    | Static: Go file server for a cute animated page     |
| `app/squirrel`     | Squirrel: Go note app — REST API + embedded UI + DB |
| `app/bunny-api`    | Backend: Go REST API backed by PostgreSQL (labs)    |
| `kubernetes/`      | Labs 01–08, self-contained with their own READMEs   |
| `.github/workflows`| CI: image build/push and manifest validation        |
| `docker-compose.yml`| Local dev stack (pages + squirrel + postgres)      |

## Local development

```bash
cp .env.example .env
docker compose up --build
# bunny page -> http://localhost:8080
# bunny api  -> http://localhost:8081/api/bunnies
# Squirrel   -> http://localhost:8082
```

## Build and push images

Images are built and pushed to GitHub Container Registry by
`.github/workflows/docker-build-push.yml`:

- `ghcr.io/<owner>/<repo>/labops-bunny`
- `ghcr.io/<owner>/<repo>/labops-bunny-api`
- `ghcr.io/<owner>/<repo>/labops-squirrel`

Update the `ghcr.io/your-org/...` placeholders in `kubernetes/**` to the images
pushed to your registry.

## Kubernetes

Start with `kubernetes/README.md`. Labs build on each other:

```bash
kubectl apply -f kubernetes/namespace.yaml
# then follow kubernetes/01-pod-service, 02-deployment, ...
```

## Roadmap

See `TODO.md`.

## Notes

- Backend uses `github.com/jackc/pgx/v5` and Postgres 16.
- Manifests target a managed cluster (EKS/GKE/AKS): default StorageClass,
  `ingressClassName: nginx` (switch to `alb`/`gce`/`azure` as needed).
