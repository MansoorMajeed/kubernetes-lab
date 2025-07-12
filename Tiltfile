

k8s_yaml('k8s/observability/namespace.yaml') # Create the monitoring namespace

# Observability setup
load('ext://helm_resource', 'helm_resource', 'helm_repo')
helm_repo('prometheus-community', 'https://prometheus-community.github.io/helm-charts')
helm_repo('grafana', 'https://grafana.github.io/helm-charts')

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
    'grafana/loki', # Chart path (repo/chart)
    namespace='monitoring',
    resource_deps=['grafana'],
    flags=['--values', 'k8s/observability/loki-values.yaml']
)

helm_resource(
    'grafana-release', # Release name
    'grafana/grafana', # Chart path (repo/chart)
    namespace='monitoring',
    resource_deps=['grafana'],
    flags=['--values', 'k8s/observability/grafana-values.yaml']
)

k8s_yaml('k8s/observability/ingress.yaml')

### Observability setup complete

# Demo Nginx deploy and service

k8s_yaml('k8s/apps/nginx-hello-world/deploy.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/service.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/ingress.yaml')

