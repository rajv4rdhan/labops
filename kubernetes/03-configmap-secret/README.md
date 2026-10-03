# Lab 03 — ConfigMaps and Secrets

Separate configuration (ConfigMap) and sensitive data (Secret) from the image
so the same image runs in every environment.

## Concepts

- **ConfigMap**: non-sensitive key/value data.
- **Secret**: sensitive data. `stringData` is written as plain text in your
  manifest and the API server base64-encodes it into `data`. It is **not**
  encrypted at rest by default — enable encryption/Secrets Manager in production.
- **Consumption**: as env vars (`envFrom`/`valueFrom`) or mounted as files.

## Apply

```bash
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f consumer.yaml
```

## Verify

```bash
kubectl get configmap,secret -n labops
kubectl logs -n labops config-demo
kubectl exec -n labops config-demo -- env | grep -E 'PORT|DB_USER|DATABASE_URL'
```

Change a value and re-run the demo Pod to see the new value:

```bash
kubectl delete pod -n labops config-demo
kubectl apply -f consumer.yaml
```

## Exercises

1. Declaratively create from files/CLI:
   ```bash
   kubectl create configmap app-flags -n labops --from-literal=FEATURE_X=on
   kubectl create secret generic api-token -n labops --from-literal=TOKEN=s3cr3t
   ```
2. Mount a ConfigMap as a file instead of env:
   ```yaml
   volumeMounts:
     - name: cfg
       mountPath: /etc/bunny
   volumes:
     - name: cfg
       configMap:
         name: bunny-config
   ```
3. Why is `stringData` preferred in GitOps? (Reviewable, no manual base64.)

## Clean up

```bash
kubectl delete -f configmap.yaml -f secret.yaml -f consumer.yaml
```
