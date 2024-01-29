# Litmus Edge Device Plugin for Kubernetes

## Overview

Pods could request to be allocated devices using the Kubernetes Pod `resources` field:
```yaml
resources:
  limits:
    litmus.io/edge: 10
```

## Getting Started

To install the plugin, choose what devices should be discovered and deploy the following DaemonSet:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: edge-device-plugin
  namespace: kube-system
  labels:
    app.kubernetes.io/name: edge-device-plugin
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: edge-device-plugin
  template:
    metadata:
      labels:
        app.kubernetes.io/name: edge-device-plugin
    spec:
      priorityClassName: system-node-critical
      tolerations:
      - operator: "Exists"
        effect: "NoExecute"
      - operator: "Exists"
        effect: "NoSchedule"
      containers:
      - image: quay.io/eochieng/edge-deviceplugin:latest
        name: edge-device-plugin
        resources:
          requests:
            cpu: 50m
            memory: 10Mi
          limits:
            cpu: 50m
            memory: 20Mi
        ports:
        - containerPort: 8080
          name: http
        securityContext:
          privileged: true
        volumeMounts:
        - name: device-plugin
          mountPath: /var/lib/kubelet/device-plugins
        - name: dev
          mountPath: /dev
      volumes:
      - name: device-plugin
        hostPath:
          path: /var/lib/kubelet/device-plugins
      - name: dev
        hostPath:
          path: /dev
  updateStrategy:
    type: RollingUpdate
```
