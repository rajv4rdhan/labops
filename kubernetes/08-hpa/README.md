# Lab 08 — Horizontal Pod Autoscaler

The HPA adds or removes Pods based on observed metrics (CPU by default). It
targets a Deployment and adjusts `replicas` between `minReplicas` and
`maxReplicas`.

## Concepts

- **metrics-server** must be installed for resource metrics:
  ```bash
  kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
  kubectl top pods -n labops
  ```
- The HPA computes `desiredReplicas = ceil(currentReplicas * currentMetric / targetMetric)`.
- Containers must set `resources.requests.cpu` — utilization is relative to the
  request.

## Apply

```bash
kubectl apply -f deployment.yaml
```

## Verify

```bash
kubectl get hpa -n labops -w
```

Initially `TARGETS` shows `<unknown>/60%` until metrics-server has data, then a
current value.

## Load test

Generate CPU load and watch the replica count climb:

```bash
kubectl run loadgen -n labops --rm -it --image=busybox:1.36 --restart=Never -- \
  sh -c 'while true; do wget -q -O- http://bunny-api.labops.svc.cluster.local/health >/dev/null; done'
```

In another terminal:

```bash
kubectl get hpa,deploy -n labops -l app=bunny-api -w
```

Stop the load generator (Ctrl+C) and the HPA scales back down after the
cooldown window.

## Notes

- For non-CPU metrics (requests/sec, queue depth), use a custom/external metrics
  adapter.
- Consider `behavior` rules to tune scale-up/scale-down rates.

## Clean up

```bash
kubectl delete -f deployment.yaml
```
