# ReplicaSet
Inspect the ReplicaSet and Pods
```
kubectl get replicasets
kubectl get pods -l app=nginx
kubectl describe rs nginx-replicaset
```

Scale the ReplicaSet
`kubectl scale replicaset nginx-replicaset --replicas=5`

Verify
`kubectl get pods -l app=nginx`

Delete One Pod and Watch It Recreate
`kubectl delete pod <one-of-the-pod-names>`

Delete the ReplicaSet
`kubectl delete replicaset name-replicaset`

# Deployment – Scaling + Rolling Updates + Rollbacks

Apply the Deployment
`kubectl apply -f deployment.yaml`
Verify
```
kubectl get deployments
kubectl get pods -l app=bunny-page
```
Check the rollout status
`kubectl rollout status deployment bunny-page`
`kubectl rollout history deployment bunny-page`