---
title: Initialization
menuTitle: Initialization
weight: 1
---

The Service Mesh Hub is easy to install with `kubectl`.

```bash
kubectl apply -f https://github.com/solo-io/service-mesh-hub/install/service-mesh-hub.yaml
```

This will create the `sm-marketplace` namespace and install the necessary resources there.


It may take up to a minute for the Service Mesh Hub to be ready. You can check it status with:
```bash
kubectl get pods -n sm-marketplace -w
```
The Service Mesh Hub will be ready for use as soon as all the listed pods are ready.

Access the Service Mesh Hub with port-forwarding:
```bash
kubectl port-forward -n sm-marketplace deployment/smm-apiserver 8080
```

You should now be able to visit the Service Mesh Hub in your browser at http://localhost:8080.


