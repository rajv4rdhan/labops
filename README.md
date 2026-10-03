# LabOps

A hands-on DevOps lab: a small multi-tier application used to practice
containers, Kubernetes, CI/CD, and beyond.

```
Browser ──Ingress──┬── /       ──> bunny-page (frontend, Go static server)
                   └── /api/*  ──> bunny-api  (backend, Go REST) ──> postgres
```

## Components

| Path               | What it is                                          |
| ------------------ | --------------------------------------------------- |
| `app/BunnyPage`    | Frontend: Go static server + animated HTML + `/api` |
| `app/bunny-api`    | Backend: Go REST API backed by PostgreSQL           |
| `kubernetes/`      | Labs 01–08, self-contained with their own READMEs   |
| `.github/workflows`| CI: image build/push and manifest validation        |
| `docker-compose.yml`| Local dev stack (frontend + api + postgres)        |

## Local development

```bash
cp .env.example .env
docker compose up --build
# frontend  -> http://localhost:8080
# api       -> http://localhost:8081/api/bunnies
```

## Build and push images

Images are built and pushed to GitHub Container Registry by
`.github/workflows/docker-build-push.yml`:

- `ghcr.io/<owner>/<repo>/labops-bunny`
- `ghcr.io/<owner>/<repo>/labops-bunny-api`

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
