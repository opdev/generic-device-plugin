# Generic Edge Device Plugin for Kubernetes

## Overview

Pods could request to be allocated devices using the Kubernetes Pod `resources` field:
```yaml
resources:
  limits:
    vendor.io/device: 10
```

## Getting Started

To install the plugin, choose what devices should be discovered and deploy the following DaemonSet:

```
oc apply -f deploy/daemonset.yaml
```
