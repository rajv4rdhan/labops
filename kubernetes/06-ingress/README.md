# Lab 06 — Ingress routing

Ingress exposes HTTP(S) routes from outside the cluster to Services. One host
serves the frontend at `/` and the API at `/api`, so the browser calls a
same-origin `/api/bunnies`.

```
Client ──> Ingress (lab.bunny.local)
              ├── /      ──> Service bunny-page:80   ──> Pods (frontend)
              └── /api   ──> Service bunny-api:8080  ──> Pods (backend)
```

## Prerequisites

- An Ingress controller. This manifest uses the nginx IngressClass. On managed
  clusters you may instead have ALB (EKS), GCE (GKE), or Application Gateway
  (AKS) — change `ingressClassName` and annotations accordingly.
- Labs 04 (Postgres) and 05 (bunny-api) running.

Install nginx ingress quickly if you don't have one:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.3/deploy/static/provider/cloud/deploy.yaml
```

## Apply

```bash
kubectl apply -f frontend.yaml
kubectl apply -f ingress.yaml
```

## Verify

```bash
kubectl get ingress -n labops
kubectl describe ingress bunny -n labops
```

Find the external address:

```bash
kubectl get ingress bunny -n labops -o jsonpath='{.status.loadBalancer.ingress[0].ip}{"\n"}'
```

Then map the host locally (or use the IP directly with a Host header):

```bash
# /etc/hosts (or C:\Windows\System32\drivers\etc\hosts)
<ADDRESS>  lab.bunny.local
# or
curl -H 'Host: lab.bunny.local' http://<ADDRESS>/api/bunnies
```

Open `http://lab.bunny.local/` — the bunny panel should list residents fetched
from `/api/bunnies` through the Ingress.

## Variants

- **Path vs host routing**: try a second rule with `host: api.bunny.local` and
  path `/`.
- **TLS**: add a `tls:` section with a Secret of type `kubernetes.io/tls`, or use
  cert-manager for automatic certificates.
- **ALB (EKS)**: set `ingressClassName: alb` and add
  `alb.ingress.kubernetes.io/scheme: internet-facing`,
  `alb.ingress.kubernetes.io/target-type: ip`.

## Clean up

```bash
kubectl delete -f ingress.yaml
kubectl delete -f frontend.yaml
```
