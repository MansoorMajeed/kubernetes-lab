#!/bin/bash

set -o errexit


# Change these to your liking
CLUSTER_NAME="kubernetes-lab"
REGISTRY_NAME="kubernetes-lab-registry"
REGISTRY_PORT="5001"
# Usinga  dedicated kubeconfig for the lab to prevent conflicts with the default kubeconfig
LAB_KUBECONFIG="$HOME/.kube/config-kubernetes-lab"



### check for all required commands
for cmd in k3d kubectl tilt docker; do
  if ! command -v $cmd &> /dev/null; then
    echo "Error: $cmd is not installed. Please install it before running this script." >&2
    exit 1
  else
    echo "$cmd is installed."
  fi 
done

# Create .kube directory if it doesn't exist
mkdir -p "$HOME/.kube"

# 1. Create a registry if it doesn't exist
if [ -z "$(k3d registry list | grep "^k3d-${REGISTRY_NAME}[[:space:]]")" ]; then
  echo "Creating local k3d registry '${REGISTRY_NAME}'..."
  k3d registry create ${REGISTRY_NAME} --port ${REGISTRY_PORT}
else
  echo "Registry '${REGISTRY_NAME}' already exists."
fi


# 2. Create a cluster if it doesn't exist
# 8081:80@loadbalancer  - expose port 80 of the load balancer to port 8081 on the host
if [ -z "$(k3d cluster list | grep "^${CLUSTER_NAME}[[:space:]]")" ]; then
  echo "Creating k3d cluster '${CLUSTER_NAME}'..."
  k3d cluster create ${CLUSTER_NAME} \
    --api-port 6443 \
    -p "8081:80@loadbalancer" \
    --registry-use k3d-${REGISTRY_NAME}:${REGISTRY_PORT} \
    --kubeconfig-update-default=false \
    --kubeconfig-switch-context=false
else
  echo "Cluster '${CLUSTER_NAME}' already exists."
fi

# Verify cluster creation
if k3d cluster list | grep -q "^${CLUSTER_NAME}[[:space:]]"; then
  echo "Cluster '${CLUSTER_NAME}' is running."
else
  echo "Error: Cluster '${CLUSTER_NAME}' was not created successfully." >&2
  exit 1
fi

# Get the kubeconfig for this cluster and save it to a dedicated file
echo "Setting up dedicated kubeconfig for the lab..."
k3d kubeconfig get ${CLUSTER_NAME} > "${LAB_KUBECONFIG}"

# Verify that the cluster is accessible using the dedicated kubeconfig
if kubectl --kubeconfig="${LAB_KUBECONFIG}" cluster-info > /dev/null 2>&1; then
  echo "Cluster '${CLUSTER_NAME}' is accessible via kubectl."
else
  echo "Error: Unable to access cluster '${CLUSTER_NAME}' via kubectl." >&2
  exit 1
fi

echo ""
echo "🎉 Kubernetes Lab Environment is ready!"
echo ""
echo "To use this lab environment, run one of the following:"
echo ""
echo "Option 1 - Use the wrapper script (easiest):"
echo "  ./kubectl-lab get nodes"
echo ""
echo "Option 2 - Set environment variable (recommended for longer sessions):"
echo "  export KUBECONFIG='${LAB_KUBECONFIG}'"
echo ""
echo "Option 3 - Use --kubeconfig flag:"
echo "  kubectl --kubeconfig='${LAB_KUBECONFIG}' get nodes"
echo ""
echo "Option 4 - Create an alias:"
echo "  alias kubectl-lab='kubectl --kubeconfig=${LAB_KUBECONFIG}'"
echo ""
echo "Your cluster is accessible at: http://localhost:8081"
echo "Registry is running on: localhost:${REGISTRY_PORT}"
echo ""