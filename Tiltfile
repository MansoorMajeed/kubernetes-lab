# Restrict Tilt to only use the k3d lab context for safety
allow_k8s_contexts('k3d-kubernetes-lab')

# Print helpful message about using the correct kubeconfig
print("📋 Make sure you're using the lab kubeconfig:")
print("   export KUBECONFIG=~/.kube/config-kubernetes-lab")
print("   OR use: ./tilt-lab up")

k8s_yaml('k8s/observability/namespace.yaml') # Create the monitoring namespace
k8s_yaml('k8s/apps/nginx-hello-world/namespace.yaml') # Create the monitoring namespace
k8s_yaml('k8s/observability/loki-dashboard-configmap.yaml') # Load the Loki dashboard
k8s_yaml('k8s/observability/dashboards/catalog-dashboard-configmap.yaml') # Load the Catalog Service dashboard

# Observability setup
load('ext://helm_resource', 'helm_resource', 'helm_repo')
helm_repo('prometheus-community', 'https://prometheus-community.github.io/helm-charts')
helm_repo('grafana-repo', 'https://grafana.github.io/helm-charts')

# Define the charts to install using helm_resource
helm_resource(
    'prometheus', # Release name
    'prometheus-community/prometheus', # Chart path (repo/chart)
    namespace='monitoring',
    resource_deps=['prometheus-community'], # Dependency on the helm_repo resource
    flags=['--values', 'k8s/observability/prometheus-values.yaml']
)

helm_resource(
    'loki', # Release name
    'grafana-repo/loki', # Chart path (repo/chart)
    namespace='monitoring',
    resource_deps=['grafana-repo'],
    flags=['--values', 'k8s/observability/loki-values.yaml']
)

helm_resource(
    'grafana', # Release name
    'grafana-repo/grafana', # Chart path (repo/chart)
    namespace='monitoring',
    resource_deps=['grafana-repo'],
    flags=['--values', 'k8s/observability/grafana-values.yaml']
)

helm_resource(
    'alloy',
    'grafana-repo/alloy',
    namespace='monitoring',
    resource_deps=['grafana-repo'],
    flags=['--values', 'k8s/observability/alloy-values.yaml']
)

# Add OTLP service for traces
k8s_yaml('k8s/observability/alloy-otlp-service.yaml')

helm_resource(
    'tempo',
    'grafana-repo/tempo',
    namespace='monitoring',
    resource_deps=['grafana-repo'],
    flags=['--values', 'k8s/observability/tempo-values.yaml']
)

k8s_yaml('k8s/observability/ingress.yaml')
k8s_yaml('k8s/observability/tempo-ingress.yaml')

### Observability setup complete

# Demo Nginx deploy and service

k8s_yaml('k8s/apps/nginx-hello-world/nginx-configmap.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/deploy.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/service.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/ingress.yaml')

### Catalog Service Setup

# Create catalog namespace
k8s_yaml('k8s/apps/catalog/namespace.yaml')

# Deploy PostgreSQL database
k8s_yaml('k8s/apps/catalog/postgres.yaml')

# Build catalog service Docker image
docker_build(
    'catalog-service:latest',
    './services/catalog',
    dockerfile='./services/catalog/Dockerfile'
)

# Deploy catalog service
k8s_yaml('k8s/apps/catalog/deployment.yaml')
k8s_yaml('k8s/apps/catalog/service.yaml')
k8s_yaml('k8s/apps/catalog/ingress.yaml')

