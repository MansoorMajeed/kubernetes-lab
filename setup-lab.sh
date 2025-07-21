#!/bin/bash

set -o errexit

# Change these to your liking
CLUSTER_NAME="kubernetes-lab"
REGISTRY_NAME="kubernetes-lab-registry"
REGISTRY_PORT="5001"
# Using a dedicated kubeconfig for the lab to prevent conflicts with the default kubeconfig
LAB_KUBECONFIG="$HOME/.kube/config-kubernetes-lab"

# Function to display usage
usage() {
  echo "Usage: $0 [OPTIONS]"
  echo ""
  echo "Options:"
  echo "  --stop, -s      Stop the lab environment (cluster and registry remain for quick restart)"
  echo "  --delete, -d    Delete and remove the lab environment completely (cluster, registry, kubeconfig)"
  echo "  --help, -h      Show this help message"
  echo ""
  echo "Without options: Start/create the lab environment"
}

# Function to stop the lab environment (but keep it for restart)
stop_lab() {
  echo "⏸️  Stopping Kubernetes Lab Environment..."
  echo ""
  
  # Stop the cluster (but don't delete it)
  if k3d cluster list | grep -q "^${CLUSTER_NAME}[[:space:]]"; then
    RUNNING_SERVERS=$(k3d cluster list | grep "^${CLUSTER_NAME}[[:space:]]" | awk '{print $2}' | cut -d'/' -f1)
    if [ "$RUNNING_SERVERS" != "0" ]; then
      echo "Stopping k3d cluster '${CLUSTER_NAME}'..."
      k3d cluster stop ${CLUSTER_NAME}
      echo "✅ Cluster '${CLUSTER_NAME}' stopped."
    else
      echo "ℹ️  Cluster '${CLUSTER_NAME}' is already stopped."
    fi
  else
    echo "ℹ️  Cluster '${CLUSTER_NAME}' does not exist."
  fi
  
  # Stop the registry (but don't delete it)
  if k3d registry list | grep -q "^k3d-${REGISTRY_NAME}[[:space:]]"; then
    REGISTRY_STATUS=$(k3d registry list | grep "^k3d-${REGISTRY_NAME}[[:space:]]" | awk '{print $3}')
    if [ "$REGISTRY_STATUS" != "exited" ]; then
      echo "Stopping k3d registry '${REGISTRY_NAME}'..."
      docker stop k3d-${REGISTRY_NAME}
      echo "✅ Registry '${REGISTRY_NAME}' stopped."
    else
      echo "ℹ️  Registry '${REGISTRY_NAME}' is already stopped."
    fi
  else
    echo "ℹ️  Registry '${REGISTRY_NAME}' does not exist."
  fi
  
  echo ""
  echo "⏸️  Lab environment stopped!"
  echo ""
  echo "To restart the lab environment, run:"
  echo "  ./setup-lab.sh"
  echo ""
  echo "To completely remove the lab environment, run:"
  echo "  ./setup-lab.sh --delete"
  echo ""
}

# Function to stop and clean up the lab environment
cleanup_lab() {
  echo "🧹 Deleting Kubernetes Lab Environment..."
  echo ""
  
  # Stop and delete the cluster
  if k3d cluster list | grep -q "^${CLUSTER_NAME}[[:space:]]"; then
    echo "Stopping and deleting k3d cluster '${CLUSTER_NAME}'..."
    k3d cluster delete ${CLUSTER_NAME}
    echo "✅ Cluster '${CLUSTER_NAME}' deleted."
  else
    echo "ℹ️  Cluster '${CLUSTER_NAME}' does not exist."
  fi
  
  # Stop and delete the registry
  if k3d registry list | grep -q "^k3d-${REGISTRY_NAME}[[:space:]]"; then
    echo "Stopping and deleting k3d registry '${REGISTRY_NAME}'..."
    k3d registry delete ${REGISTRY_NAME}
    echo "✅ Registry '${REGISTRY_NAME}' deleted."
  else
    echo "ℹ️  Registry '${REGISTRY_NAME}' does not exist."
  fi
  
  # Remove the kubeconfig file
  if [ -f "${LAB_KUBECONFIG}" ]; then
    echo "Removing lab kubeconfig file..."
    rm "${LAB_KUBECONFIG}"
    echo "✅ Kubeconfig file removed: ${LAB_KUBECONFIG}"
  else
    echo "ℹ️  Kubeconfig file does not exist: ${LAB_KUBECONFIG}"
  fi
  
  echo ""
  echo "🗑️  Lab environment completely removed!"
  echo ""
  echo "To recreate the lab environment, run:"
  echo "  ./setup-lab.sh"
  echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --stop|-s)
      stop_lab
      exit 0
      ;;
    --delete|-d)
      cleanup_lab
      exit 0
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      echo ""
      usage
      exit 1
      ;;
  esac
done


### check for all required commands
for cmd in k3d kubectl tilt docker helm; do
  if ! command -v $cmd &> /dev/null; then
    echo "Error: $cmd is not installed. Please install it before running this script." >&2
    exit 1
  else
    echo "$cmd is installed."
  fi 
done

# Create .kube directory if it doesn't exist
mkdir -p "$HOME/.kube"

# 1. Create a registry if it doesn't exist, or start it if it's stopped
if [ -z "$(k3d registry list | grep "^k3d-${REGISTRY_NAME}[[:space:]]")" ]; then
  echo "Creating local k3d registry '${REGISTRY_NAME}'..."
  k3d registry create ${REGISTRY_NAME} --port ${REGISTRY_PORT}
else
  # Registry exists, check if it's running
  REGISTRY_STATUS=$(k3d registry list | grep "^k3d-${REGISTRY_NAME}[[:space:]]" | awk '{print $3}')
  if [ "$REGISTRY_STATUS" = "exited" ]; then
    echo "Registry '${REGISTRY_NAME}' exists but is stopped. Starting it..."
    docker start k3d-${REGISTRY_NAME}
  else
    echo "Registry '${REGISTRY_NAME}' already exists and is running."
  fi
fi

# Verify registry is accessible
echo "Verifying registry accessibility..."
if curl -s "http://localhost:${REGISTRY_PORT}/v2/" > /dev/null 2>&1; then
  echo "Registry is accessible at localhost:${REGISTRY_PORT}"
else
  echo "Warning: Registry may not be fully ready yet. This might cause Docker push timeouts." >&2
  echo "If you encounter push timeouts, run: docker start k3d-${REGISTRY_NAME}" >&2
fi

# 2. Create a cluster if it doesn't exist, or start it if it's stopped
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
  # Cluster exists, check if it's running
  RUNNING_SERVERS=$(k3d cluster list | grep "^${CLUSTER_NAME}[[:space:]]" | awk '{print $2}' | cut -d'/' -f1)
  if [ "$RUNNING_SERVERS" = "0" ]; then
    echo "Cluster '${CLUSTER_NAME}' exists but is stopped. Starting it..."
    k3d cluster start ${CLUSTER_NAME}
  else
    echo "Cluster '${CLUSTER_NAME}' already exists and is running."
  fi
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
echo "📋 To use Tilt with this lab environment:"
echo "  ./tilt-lab up    # Using wrapper script (recommended)"
echo "  OR"
echo "  export KUBECONFIG='${LAB_KUBECONFIG}' && tilt up"
echo ""
echo "🛑 To manage the lab environment:"
echo "  ./setup-lab.sh --stop    # Stop (but keep for quick restart)"
echo "  ./setup-lab.sh --delete  # Completely remove everything"
echo ""
echo "Your cluster is accessible at: http://localhost:8081"
echo "Registry is running on: localhost:${REGISTRY_PORT}"
echo ""