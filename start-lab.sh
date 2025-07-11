#!/bin/bash

set -o errexit


# Change these to your liking
CLUSTER_NAME="kubernetes-lab"
REGISTRY_NAME="kubernetes-lab-registry"
REGISTRY_PORT="5001"



### check for all required commands
for cmd in k3d kubectl tilt docker; do
  if ! command -v $cmd &> /dev/null; then
    echo "Error: $cmd is not installed. Please install it before running this script." >&2
    exit 1
  else
    echo "$cmd is installed."
  fi 
done


# 1. Create a registry if it doesn't exist
if [ -z "$(k3d registry list | grep ${REGISTRY_NAME})" ]; then
  echo "Creating local k3d registry '${REGISTRY_NAME}'..."
  k3d registry create ${REGISTRY_NAME} --port ${REGISTRY_PORT}
else
  echo "Registry '${REGISTRY_NAME}' already exists."
fi


# 2. Create a cluster if it doesn't exist
# 8081:80@loadbalancer  - expose port 80 of the load balancer to port 8081 on the host
if [ -z "$(k3d cluster list | grep ${CLUSTER_NAME})" ]; then
  echo "Creating k3d cluster '${CLUSTER_NAME}'..."
  k3d cluster create ${CLUSTER_NAME} \
    --api-port 6443 \
    -p "8081:80@loadbalancer" \
    --registry-use k3d-${REGISTRY_NAME}:${REGISTRY_PORT}
else
  echo "Cluster '${CLUSTER_NAME}' already exists."
fi

# Verify cluster creation
if k3d cluster list | grep -q ${CLUSTER_NAME}; then
  echo "Cluster '${CLUSTER_NAME}' is running."
else
  echo "Error: Cluster '${CLUSTER_NAME}' was not created successfully." >&2
  exit 1
fi

# Verify that the cluster is accessible
if kubectl cluster-info --context k3d-${CLUSTER_NAME} > /dev/null 2>&1; then
  echo "Cluster '${CLUSTER_NAME}' is accessible via kubectl."
else
  echo "Error: Unable to access cluster '${CLUSTER_NAME}' via kubectl." >&2
  exit 1
fi