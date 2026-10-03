# Lab 01 — Pod and Service

A Pod is the smallest deployable unit. A Service gives a stable virtual IP and
DNS name to a set of Pods selected by labels.

## Concepts

- **Pod**: one or more containers sharing a network namespace and storage.
- **Label selector**: how a Service finds its Pods (`spec.selector` must match
  the Pod's `metadata.labels`).
- **ClusterIP**: internal-only service (default).
- **NodePort**: exposes the service on every node's IP at a static port
  (`30000-32767`) — good for learning, not production.

## Apply

```bash
kubectl apply -f pod.yaml
kubectl apply -f service.yaml
```

## Verify

```bash
kubectl get pod -n labops -o wide
kubectl get svc -n labops
kubectl describe svc bunny-page -n labops
```

Visit the app:

- Local/bare-metal: `http://<node-ip>:30080`
- minikube: `minikube service bunny-page -n labops`
- Managed cloud without a routable node: use `kubectl port-forward`:

```bash
kubectl port-forward -n labops svc/bunny-page 8080:80
# then open http://localhost:8080
```

## Exercises

1. Change `nodePort` to `30081` and re-apply. What happens?
2. Delete the Pod and re-apply only the Service. Why does the Service have no
   endpoints? (`kubectl get endpoints bunny-page -n labops`)
3. Break the selector label and observe the empty endpoints, then fix it.

## Clean up

```bash
kubectl delete -f pod.yaml -f service.yaml
```
