# Observability Stack

Complete monitoring setup for Kubernetes with Prometheus, Loki, Grafana, and Alloy.

## Components

- **Prometheus**: Metrics collection and storage
- **Loki**: Log aggregation and storage  
- **Grafana**: Visualization and dashboards
- **Alloy**: Unified telemetry collector (logs via Kubernetes API)

## Access

- **Alloy UI**: http://localhost:12345 (port-forward alloy service)
- **Prometheus**: http://localhost:9090 (port-forward prometheus service)

## Configuration

### Loki
- Stores logs from all Kubernetes pods
- Retention: Check `loki-values.yaml` for settings
- Gateway: `loki-gateway.monitoring.svc.cluster.local`

### Alloy
- Uses Kubernetes API for log collection (no filesystem access)
- Filters system pods automatically
- Deployed as regular deployment (not DaemonSet)

### Grafana
- Pre-configured with Prometheus and Loki datasources. Check the `grafana-values.yaml`
- Dashboard discovery via ConfigMap labels
- Default credentials: admin/password
