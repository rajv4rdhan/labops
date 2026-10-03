# Kubernetes Labs

Hands-on labs for the LabOps stack (BunnyPage frontend, bunny-api backend,
PostgreSQL). Each folder is self-contained: a `README.md` explaining the concept
and the manifests to apply.

Labs are written for a managed cluster (EKS/GKE/AKS). Storage uses the cluster's
default StorageClass and Ingress uses the nginx `ingressClassName` by default.

## Before you start

```bash
kubectl apply -f namespace.yaml
```

Replace the placeholder image references
(`ghcr.io/your-org/labops-bunny*`) with images pushed to your registry. See the
root `README`/CI workflow for building and pushing.

## Lab order

| # | Folder | Topic |
| - | ------ | ----- |
| 01 | `01-pod-service` | Pod and Service (ClusterIP/NodePort) |
| 02 | `02-deployment` | Deployments, scaling, rollout/rollback |
| 03 | `03-configmap-secret` | ConfigMaps and Secrets |
| 04 | `04-postgres` | Postgres Deployment, PVC, Service |
| 05 | `05-bunny-api` | Backend wired to Postgres |
| 06 | `06-ingress` | Ingress routing `/` and `/api` |
| 07 | `07-probes` | Startup/readiness/liveness probes, resources |
| 08 | `08-hpa` | Horizontal Pod Autoscaler |

## Clean up

```bash
kubectl delete namespace labops
```
