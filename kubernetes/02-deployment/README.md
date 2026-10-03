# Lab 02 — Deployment, scaling, rollouts

A Deployment manages a ReplicaSet, which keeps a desired number of identical
Pods running and performs declarative rolling updates.

## Concepts

- **replicas**: desired Pod count; the ReplicaSet reconciles reality to match.
- **RollingUpdate**: replaces Pods gradually (`maxSurge`, `maxUnavailable`).
- **Revisions**: each image/config change creates a new revision you can roll
  back to.
- **resources**: `requests` affect scheduling; `limits` cap usage.

## Apply

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

## Verify

```bash
kubectl get deploy,rs,pods -n labops -l app=bunny-page
kubectl rollout status deployment/bunny-page -n labops
kubectl get pods -n labops -l app=bunny-page -o custom-columns=NAME:.metadata.name,NODE:.spec.nodeName
```

## Exercises

1. **Scale** to 5 and back:
   ```bash
   kubectl scale deployment/bunny-page -n labops --replicas=5
   kubectl scale deployment/bunny-page -n labops --replicas=2
   ```
2. **Rolling update**: change the image tag (e.g. `:v2`), re-apply, then watch:
   ```bash
   kubectl rollout status deployment/bunny-page -n labops
   kubectl rollout history deployment/bunny-page -n labops
   ```
3. **Rollback**:
   ```bash
   kubectl rollout undo deployment/bunny-page -n labops
   ```
4. **Self-healing**: delete a Pod and watch the ReplicaSet recreate it:
   ```bash
   kubectl delete pod -n labops -l app=bunny-page
   ```
5. Inspect a Pod's assigned resources:
   ```bash
   kubectl describe pod -n labops -l app=bunny-page
   ```

## Clean up

```bash
kubectl delete -f deployment.yaml -f service.yaml
```
