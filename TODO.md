# LabOps Roadmap

A DevOps lab: a small multi-tier app used to practice containers, Kubernetes,
CI/CD, and (later) gRPC, observability, and GitOps.

## Done

- [x] BunnyPage frontend — simplified to a static cute page, `scratch` image
- [x] Squirrel note app — Go REST API + embedded UI + PostgreSQL, no login
- [x] bunny-api backend — Go REST API + PostgreSQL
- [x] Local dev with docker-compose (pages + squirrel + postgres)
- [x] Kubernetes labs 01–08 (pod/service, deployment, configmap/secret,
      postgres, api, ingress, probes, HPA)
- [x] CI — build & push images, validate Kubernetes manifests

## Next

- [ ] gRPC: add `bunny-api-grpc`, headless Service (lab 09)
- [ ] WebSockets: streaming endpoint over gRPC (lab 10)
- [ ] Helm chart packaging all resources (lab 11)
- [ ] Multi-microservice: split inventory/orders services (lab 12)
- [ ] Observability: Prometheus + Grafana, structured logging (lab 13)
- [ ] GitOps with ArgoCD (lab 14)
- [ ] CI/CD: GitHub Actions -> registry -> ArgoCD (lab 15)
- [ ] Progressive delivery: Argo Rollouts / canary (lab 16)
