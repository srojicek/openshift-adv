
# Secret Operator for OpenShift

## Overview
This operator generates or updates Kubernetes Secrets based on Custom Resources (CRs). A CronJob runs the operator periodically to ensure Secrets are refreshed.

## Prerequisites
- OpenShift CRC
- Podman
- oc CLI tools

## Build and Deploy

### Build the Operator Image

```sh
cd $PATH/openshift-adv/Day3/3.capstone/app/docker

podman build -t default-route-openshift-image-registry.apps-crc.testing/dev/capstone-app:latest .
podman build -t default-route-openshift-image-registry.apps-crc.testing/test/capstone-app:latest .
podman build -t default-route-openshift-image-registry.apps-crc.testing/prod/capstone-app:latest .
```

### Push to OpenShift Internal Registry
```sh
oc login -u kubeadmin https://api.crc.testing:6443
podman login -u kubeadmin -p $(oc whoami -t) default-route-openshift-image-registry.apps-crc.testing
podman push default-route-openshift-image-registry.apps-crc.testing/dev/capstone-app:latest
podman push default-route-openshift-image-registry.apps-crc.testing/test/capstone-app:latest
podman push default-route-openshift-image-registry.apps-crc.testing/prod/capstone-app:latest
```


## Deployment with Kustomize
```sh
cd ..
oc apply -k kustomize/overlays/dev
oc apply -k kustomize/overlays/test
oc apply -k kustomize/overlays/prod
```


---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**