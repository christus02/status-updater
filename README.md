# Controller to Update the LoadBalancer Status of a Service

This controller updates the `Service.Status.LoadBalancer.Ingress` field based on the Inputs passed to the controller.

The controller accepts inputs through ENV varibles and below are the explanations for the ENV variables.

| ENV Variables          		| Explanation           								       | Possible Values |
| :-----------------------------------:	| :------------------------------------------------------------------------------------------: | :-------------: |
| `SERVICE_NAME` 	 		| Name of the Service to watch for                                                             |
| `SERVICE_NAMESPACE`    		| Namespace of the Service                                                                     |
| `EXTERNAL_ENDPOINT_TYPE_ANNOTATION`   | Annotation to look for to get the Endpoint type. Example: `status.service.com/endpoint-type` | Values can be `ip`/`hostname`
| `ENDPOINT_ANNOTATION` 		| Annotation to look for to get the Endpoint. Example: `status.service.com/endpoint`           | Values can be any valid IP or any valid Hostname depending on the `EXTERNAL_ENDPOINT_TYPE_ANNOTATION` | 

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
          image: "raghulc/status-patcher:v1.0"
          env:
          - name: "SERVICE_NAME"
            value: "cpx-service"
          - name: "SERVICE_NAMESPACE"
            value: "default"
          - name: "EXTERNAL_ENDPOINT_TYPE_ANNOTATION"
            value: "status.service.com/endpoint-type"
          - name: "ENDPOINT_ANNOTATION"
            value: "status.service.com/endpoint"
          imagePullPolicy: Always
```

### Sample Service example

```
apiVersion: v1
kind: Service
metadata:
  name: cpx-service
  annotations:
    status.service.com/endpoint-type: hostname
    status.service.com/endpoint: "cpx.status-update.com"
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

