# Observability setup
helm_repo(name='prometheus-community', url='https://prometheus-community.github.io/helm-charts')
helm_repo(name='grafana', url='https://grafana.github.io/helm-charts')

helm_chart(
    name='prometheus',
    chart='prometheus-community/prometheus',
    namespace='monitoring',
    create_namespace=True,
    values='k8s/observability/prometheus-values.yaml'
)

helm_chart(
    name='loki',
    chart='grafana/loki',
    namespace='monitoring',
    values='k8s/observability/loki-values.yaml'
)

helm_chart(
    name='grafana',
    chart='grafana/grafana',
    namespace='monitoring',
    values='k8s/observability/grafana-values.yaml'
)

k8s_yaml('k8s/observability/ingress.yaml')

### Observability setup complete ###

# Demo Nginx deploy and service
k8s_yaml('k8s/apps/nginx-hello-world/deploy.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/service.yaml')
k8s_yaml('k8s/apps/nginx-hello-world/ingress.yaml')

