
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
cd $PATH/openshift-adv/Day1/2.crd/

podman build -t default-route-openshift-image-registry.apps-crc.testing/default/secret-operator:latest .
```

### Push to OpenShift Internal Registry
```sh
oc login -u kubeadmin -p nobleprog123.
oc new-project default
podman login -u kubeadmin -p $(oc whoami -t) default-route-openshift-image-registry.apps-crc.testing
podman push default-route-openshift-image-registry.apps-crc.testing/default/secret-operator:latest 
```

### Deploy the Operator
```sh
oc apply -f rbac.yaml
oc apply -f crd.yaml
oc apply -f cronjob.yaml
```

### Create a Custom Resource (CR)
```sh
oc apply -f cr1.yaml
```

## Verify Deployment
### Check if the CronJob is created
```sh
oc get cronjobs
```

### List Secrets created by the Operator
```sh
oc get secrets
```

### Describe a Secret to confirm content
```sh
oc describe secret <secret-name>
```

### Create a second Custom Resource (CR)
```sh
oc apply -f cr2.yaml
```

## Cleanup
```sh
oc delete -f cr1.yaml
oc delete -f cr2.yaml
oc delete -f cronjob.yaml
oc delete -f crd.yaml
oc delete -f rbac.yaml
```

---
**Copyright (c) 2025 by Alexander Kolin. All rights reserved.**