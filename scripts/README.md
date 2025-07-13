# Catalog Service Scripts

This directory contains scripts for testing and simulating traffic against the catalog service.

## 🚀 Traffic Simulator (`simulate-traffic.sh`)

A comprehensive script that seeds realistic product data and simulates authentic e-commerce traffic patterns.

### Features

- **Data Seeding**: Creates 15 realistic Apple products with proper pricing and stock levels
- **Traffic Simulation**: Mimics real user behavior with weighted request patterns:
  - 45% Browse products (listing with pagination)
  - 30% View specific products
  - 15% Stock updates (inventory management)
  - 5% Price updates
  - 5% Add new products
- **Real-time Monitoring**: Colored console output showing all operations
- **Statistics**: Complete success/failure metrics
- **Configurable**: Duration, request intervals, and service endpoints

### Quick Start

```bash
# Basic usage - 5 minutes of simulation
./scripts/simulate-traffic.sh

# Just seed data without traffic simulation
./scripts/simulate-traffic.sh --seed-only

# Heavy traffic simulation - 10 minutes, 1 second intervals
./scripts/simulate-traffic.sh --duration 600 --interval 1

# Light traffic simulation - 2 minutes, 5 second intervals  
./scripts/simulate-traffic.sh --duration 120 --interval 5
```

### Prerequisites

Before running the script, ensure:

1. **Service is running** - The catalog service must be deployed and accessible
2. **DNS resolution** - Add to `/etc/hosts`:
   ```bash
   echo "127.0.0.1 catalog.kubelab.lan" | sudo tee -a /etc/hosts
   ```
3. **Dependencies** - Requires `curl` and `bc` (basic calculator)

### Usage Examples

#### Development Testing
```bash
# Quick data seeding for development
./scripts/simulate-traffic.sh --seed-only

# Light traffic for development testing
./scripts/simulate-traffic.sh --duration 60 --interval 3
```

#### Observability Testing
```bash
# Generate sustained traffic for monitoring dashboards
./scripts/simulate-traffic.sh --duration 900 --interval 2 --verbose

# High-frequency requests for stress testing logs/metrics
./scripts/simulate-traffic.sh --duration 300 --interval 0.5
```

#### Demo Scenarios
```bash
# Seed data then run realistic demo traffic
./scripts/simulate-traffic.sh --duration 300 --interval 2

# Continuous background traffic (restart when finished)
while true; do
    ./scripts/simulate-traffic.sh --no-seed --duration 300 --interval 3
    sleep 10
done
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `--duration SECONDS` | How long to run simulation | 300 (5 min) |
| `--interval SECONDS` | Seconds between requests | 2 |
| `--host HOST:PORT` | Catalog service endpoint | `catalog.kubelab.lan:8081` |
| `--verbose` | Show detailed request/response info | false |
| `--seed-only` | Only create products, no traffic | false |
| `--no-seed` | Skip seeding, only simulate traffic | false |
| `--help` | Show usage information | - |

### Environment Variables

You can also configure via environment variables:

```bash
export CATALOG_HOST="localhost:8080"
export DURATION=600
export REQUEST_INTERVAL=1
export VERBOSE=true

./scripts/simulate-traffic.sh
```

### Sample Output

```bash
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 Catalog Service Traffic Simulator
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[14:30:15] 🔍 Checking if catalog service is available at http://catalog.kubelab.lan:8081
[14:30:15] ✅ Catalog service is healthy
[14:30:15] 🌱 Seeding initial product data...
[14:30:16] ✅ Created: MacBook Pro 14" (ID: 1)
[14:30:17] ✅ Created: iPhone 15 Pro (ID: 2)
[14:30:18] ✅ Created: AirPods Pro (ID: 3)
...
[14:30:30] 🌱 Seeded 15 products successfully
[14:30:30] 🚀 Starting traffic simulation for 300 seconds...
[14:30:30] 📊 Request interval: 2 seconds
[14:30:30] 👀 Browsed products (page: 2, limit: 8)
[14:30:32] 👁️ Viewed product ID: 5
[14:30:34] 📦 Updated stock for product 3: 147 units
[14:30:36] 👀 Browsed products (page: 1, limit: 12)
...
```

### Troubleshooting

**Service not available?**
- Check if catalog service is running: `kubectl-lab get pods -n catalog`
- Verify ingress: `kubectl-lab get ingress -n catalog`
- Test DNS: `ping catalog.kubelab.lan`

**Permission denied?**
```bash
chmod +x scripts/simulate-traffic.sh
```

**Missing dependencies?**
```bash
# On macOS
brew install bc

# On Ubuntu/Debian
sudo apt-get install bc
```

### Integration with Observability

This script is perfect for testing your observability stack:

- **Loki**: Generates structured logs for log aggregation testing
- **Prometheus**: Creates metrics for monitoring dashboard validation
- **Grafana**: Provides realistic data for dashboard development
- **Alloy**: Tests log collection and forwarding

Run the simulator while developing dashboards to see real traffic patterns in your monitoring tools! 