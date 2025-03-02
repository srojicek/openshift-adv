#!/bin/bash

SOURCE_IMAGE="default-route-openshift-image-registry.apps-crc.testing/default/secret-operator:latest"
TARGET_NAMESPACES=("dev" "test" "prod")

echo "🔄 Pulling image from default namespace..."
podman pull --tls-verify=false "$SOURCE_IMAGE"

if [ $? -ne 0 ]; then
    echo "❌ Error: Failed to pull the image from default namespace!"
    exit 1
fi

echo "🔄 Logging into OpenShift internal registry..."
podman login -u kubeadmin -p "$(oc whoami -t)" --tls-verify=false default-route-openshift-image-registry.apps-crc.testing

if [ $? -ne 0 ]; then
    echo "❌ Error: Failed to log into OpenShift internal registry!"
    exit 1
fi

for NAMESPACE in "${TARGET_NAMESPACES[@]}"; do
    TARGET_IMAGE="default-route-openshift-image-registry.apps-crc.testing/$NAMESPACE/secret-operator:latest"

    echo "🔄 Retagging image for namespace: $NAMESPACE..."
    podman tag "$SOURCE_IMAGE" "$TARGET_IMAGE"

    echo "🚀 Pushing image to $NAMESPACE namespace..."
    podman push --tls-verify=false "$TARGET_IMAGE"

    if [ $? -ne 0 ]; then
        echo "❌ Error: Failed to push the image to $NAMESPACE namespace!"
    else
        echo "✅ Successfully pushed image to $NAMESPACE namespace!"
    fi
done