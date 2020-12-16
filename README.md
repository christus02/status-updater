# Controller to Update the LoadBalancer Status of a Service

Status Updater is a controller that updates the `Service.Status.LoadBalancer.Ingress` field of a Kubernetes Service. We can choose to annotate a Kubernetes Service where we specify the IP or the Hostname it needs to be updated with.

When a Kubernetes Service is annotated, this controller parses the annotation and then updates the LoadBalancer Status respectively for that Kubernetes Service.

This Controller supports both `hostname` and `ip` updates in `Service.Status.LoadBalancer.Ingress` field.

The annotation for providing the value of the Hostname and IP can also be modified based on user's needs. If not specified, it falls back to the default annotation

| ENV Variable | Defaults (if not specified) | Explanation | 
|:------------:| :--------------------------:|:-----------:|
| `LOADBALANCER_IP_ANNOTATION` | `status.service.com/loadbalancer-ip` | Argument to Specify a custom annotation for IP address to be updated |
| `LOADBALANCER_HOSTNAME_ANNOTATION` | `status.service.com/loadbalancer-hostname` | Argument to Specify a custom annotation for Hostname address to be updated | 

## Status Updater Controller Deployment

```
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: status-updater
rules:
  - apiGroups: [""]
    resources: ["services/status"]
    verbs: ["update", "patch"]
  - apiGroups: [""]
    resources: ["services"]
    verbs: ["get", "list", "watch", "patch"]
---

kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: status-updater
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: status-updater
subjects:
- kind: ServiceAccount
  name: status-updater
  namespace: default
apiVersion: rbac.authorization.k8s.io/v1

---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: status-updater
  namespace: default

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: status-updater
spec:
  selector:
    matchLabels:
      app: status-updater
  replicas: 1
  template:
    metadata:
      name: status-updater
      labels:
        app: status-updater
      annotations:
    spec:
      serviceAccountName: status-updater
      containers:
        - name: status-updater
          image: "raghulc/status-patcher:v2.0"
          # Specify ENV variables if a custom annotation is needed
          # See previous section for explanation
          imagePullPolicy: Always
```

### Sample Service example

```
apiVersion: v1
kind: Service
metadata:
  name: cpx-service
  annotations:
    status.service.com/loadbalancer-hostname: "cpx.status-update.com"
  labels:
    app: cpx-service
spec:
  type: LoadBalancer
  externalTrafficPolicy: Local
  ports:
  - port: 80
    protocol: TCP
    name: http
  - port: 443
    protocol: TCP
    name: https
  selector:
    app: cpx-ingress
```

### Annotate a Service using `kubectl`

`kubectl annotate service cpx-service status.service.com/loadbalancer-hostname="new-hostname.abc.com" --overwrite`

**Note:** To delete an existing IP or Hostname, just clear the respective annotation.

`kubectl annotate service cpx-service status.service.com/loadbalancer-hostname= --overwrite`
