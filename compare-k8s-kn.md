# Your Task

Deploy the
[canonical](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)
Kubernetes PHP Web App Autoscaling and minimal `Lines-of-YAML`.

## The Kubernetes Way
### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: php-apache
spec:
  selector:
    matchLabels:
      run: php-apache
  replicas: 1
  template:
    metadata:
      labels:
        run: php-apache
    spec:
      containers:
        - name: php-apache
          image: k8s.gcr.io/hpa-example
          ports:
            - containerPort: 80
          resources:
            limits:
              cpu: 1
              memory: 100Mi
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: php-apache
  labels:
    run: php-apache
spec:
  ports:
    - port: 80
  selector:
    run: php-apache
```

### Autoscaling (HPA)

```yaml
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: php-apache
spec:
  maxReplicas: 10
  minReplicas: 1
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: php-apache
  targetCPUUtilizationPercentage: 50
```

### ⚠️ Out of Scope

Ingress, custom metrics, DNS/TLS, traffic splitting and rollout management, etc.

## The Knative Way
### KService

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: php-apache-kn
spec:
  template:
    spec:
      containers:
      - image: k8s.gcr.io/hpa-example
        ports:
        - containerPort: 80 # we only need this because port != 8080
```

Yup, that's really it...and you even get `scale-to-zero` which HPA as of today won't give you 😯