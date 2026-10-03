# Lab 05 — bunny-api wired to PostgreSQL

Deploy the backend and connect it to the database from Lab 04 using a ConfigMap
for non-secrets and a Secret for the connection string.

## Concepts

- The backend reads `DATABASE_URL` from the environment (injected from a Secret)
  and `PORT` from a ConfigMap.
- On startup it creates its `bunnies` table and seeds rows if empty.
- The Service name `bunny-api` becomes the in-cluster DNS name other workloads
  use to reach it.

## Prerequisites

Lab 04 must be running (Postgres service `postgres`).

## Apply

```bash
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

## Verify

```bash
kubectl get deploy,svc,pods -n labops -l app=bunny-api
kubectl port-forward -n labops svc/bunny-api 8081:8080
```

In another terminal:

```bash
curl http://localhost:8081/health
curl http://localhost:8081/ready
curl http://localhost:8081/api/bunnies
curl -X POST http://localhost:8081/api/bunnies \
  -H 'content-type: application/json' \
  -d '{"name":"Nibbles","color":"peach"}'
curl http://localhost:8081/api/bunnies
```

Expected: `/health` → `{"status":"ok"}`, `/ready` → `{"status":"ready"}`, and
the list contains the seeded bunnies (`Clover`, `Pip`, `Mochi`) plus any you add.

## Troubleshooting

```bash
kubectl logs -n labops -l app=bunny-api
kubectl describe pod -n labops -l app=bunny-api
# Can the pod resolve the DB service?
kubectl run dnsutils -n labops --rm -it --image=busybox:1.36 -- nslookup postgres
```

## Clean up

```bash
kubectl delete -f deployment.yaml -f service.yaml -f configmap.yaml -f secret.yaml
```
