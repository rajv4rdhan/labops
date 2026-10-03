# Lab 07 — Health probes and resources

Kubernetes uses probes to decide when a container is alive and ready to receive
traffic. This lab adds all three probe types to `bunny-api`.

## Concepts

- **startupProbe**: "has the app finished booting?" While it runs, liveness and
  readiness are suppressed. Good for slow-starting apps.
- **readinessProbe**: "can it serve traffic?" Failing pods are removed from
  Service endpoints but not restarted. `/ready` checks the DB.
- **livenessProbe**: "is it wedged?" Failing pods are restarted. `/health` is a
  cheap process check.
- **resources**: `requests` drive scheduling; `limits` cap usage and protect
  neighbors. CPU limits cause throttling, memory limits cause OOMKill.

## Apply

```bash
kubectl apply -f deployment.yaml
```

Requires the `bunny-db` Secret (Lab 04/05) and a running Postgres.

## Verify

```bash
kubectl get pods -n labops -l app=bunny-api
kubectl describe pod -n labops -l app=bunny-api | Select-String -Pattern 'Liveness|Readiness|Startup|State'
```

On Windows PowerShell the pipe above works; on bash use:

```bash
kubectl describe pod -n labops -l app=bunny-api | grep -E 'Liveness|Readiness|Startup|State'
```

Watch endpoints change as readiness flips:

```bash
kubectl get endpoints bunny-api -n labops -w
```

## Exercises

1. Stop Postgres (`kubectl scale deploy/postgres -n labops --replicas=0`). The
   pods stay `Running` but become `NotReady`; `bunny-api` endpoints empty.
   Restart Postgres and they rejoin.
2. Set an aggressive liveness probe (`periodSeconds: 1, failureThreshold: 1`)
   and watch a healthy pod restart.
3. Set a tiny memory limit and observe `OOMKilled` status.

## Clean up

```bash
kubectl delete -f deployment.yaml
```
