# Container API

Container API is a tool to create different instances of the same image on kubernetes.  
It's suitable for small rental container services. For example, minecraft server or code-server.

# To start using Container API

## Require System

- kubernetes
- longhorn
- cilium
- metallb

## Deploy Container API

### Minecraft BE

```bash
#!/bin/bash

NAMESPACE=ctrapi

kubectl create namespace $NAMESPACE
kubectl -n $NAMESPACE apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ctrapi
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: ctrapi
rules:
- apiGroups:
  - ""
  - apps
  - extentions
  resources:
  - pods
  - services
  - persistentvolumeclaims
  - persistentvolumeclaims/status
  - deployments
  verbs:
  - get
  - list
  - create
  - delete
  - patch
  - update
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ctrapi
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: ctrapi
subjects:
- kind: ServiceAccount
  name: ctrapi
  namespace: ctrapi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: ctrapi
  name: ctrapi
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ctrapi
  template:
    metadata:
      labels:
        app: ctrapi
    spec:
      serviceAccountName: ctrapi
      automountServiceAccountToken: true
      containers:
      - name: ctrapi
        image: docheio/git-runner-go:1.0
        imagePullPolicy: Always
        env:
        - name: REPO
          value: https://github.com/docheio/container-api.git
        - name : DIR
          value: ./
        - name: BUILD
          value: "true"
        - name: BUILD_START_COMMAND
          value: "go build -o ./ctrapi && ./ctrapi --namespace ctrapi --key mcbe --image \"docheio/minecraft-be\""
        ports:
        - name: tcp8080
          containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: ctrapi
spec:
  selector:
    app: ctrapi
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
  type: LoadBalancer
# dont use LoadBalancer when deploying
EOF
```

## Operate

```bash
curl -s -X GET    127.0.0.1:8081/v1/
curl -s -X DELETE 127.0.0.1:8081/v1/INSTANCE_ID
curl -s -X POST   -d '{"cpu":1000,"mem":1000,"ports":[{"protocol":"UDP","internal":19132}],"pvcs":[{"mount":"/root/minecraft","size":5}]}' 127.0.0.1:8081/v1/
curl -s -X PUT    -d  '{"cpu":2000,"mem":2000,"ports":[{"protocol":"UDP","internal":19132}],"pvcs":[{"id":"","mount":"/root/minecraft","size":6}]}' 127.0.0.1:8081/v1/
```
