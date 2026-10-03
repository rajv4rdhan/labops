# Lab 04 — PostgreSQL with a PersistentVolumeClaim

Stateful workloads need durable storage. Here Postgres runs as a single-replica
Deployment backed by a `PersistentVolumeClaim` (PVC). The stable DNS name
`postgres` comes from a headless-free ClusterIP Service.

## Concepts

- **PVC → PV**: a claim requests storage; the cluster binds it to a volume.
  With no `storageClassName`, the cluster's **default** StorageClass is used
  (EKS `gp3`/`gp2`, GKE `standard-rwo`, AKS `managed-csi`).
- **Recreate strategy**: the old Pod must terminate before the new one starts
  because `ReadWriteOnce` allows only one node to mount the volume.
- **StatefulSet**: for databases you often want a StatefulSet (stable identity,
  ordered rollout). A Deployment is used here to keep the lab simple.

## Apply

```bash
kubectl apply -f secret.yaml
kubectl apply -f pvc.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

## Verify

```bash
kubectl get pvc -n labops
kubectl get deploy,svc,pods -n labops -l app=postgres
kubectl exec -n labops deploy/postgres -- psql -U bunny -d bunnies -c '\l'
```

## Data persistence exercise

```bash
kubectl exec -n labops deploy/postgres -- psql -U bunny -d bunnies \
  -c 'CREATE TABLE t(id int); INSERT INTO t VALUES (1);'
kubectl delete pod -n labops -l app=postgres      # PVC survives
kubectl exec -n labops deploy/postgres -- psql -U bunny -d bunnies -c 'SELECT * FROM t;'
```

If the row is still there, the data lived on the volume, not the Pod.

## Notes for production

- Use `StatefulSet` + volumeClaimTemplates, backups, and a managed DB
  (RDS/Cloud SQL/Azure Database) when possible.
- Set `storageClassName` explicitly to avoid surprises.

## Clean up

```bash
kubectl delete -f deployment.yaml -f service.yaml
kubectl delete -f pvc.yaml          # deletes the data
kubectl delete -f secret.yaml
```
