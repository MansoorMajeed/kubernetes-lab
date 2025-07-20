# Catalog Service

A **production-ready microservice** built for learning cloud-native observability patterns. This service demonstrates complete **distributed tracing**, **structured logging**, and **metrics collection** in a Kubernetes environment.

## 🚀 What This Service Does

The Catalog Service is a **product catalog API** that manages an inventory of products with full **CRUD operations**. It's designed to showcase modern observability practices in a microservices architecture.

### Key Features
- 🛍️ **Product Management**: Complete REST API for product operations
- 🔍 **Advanced Analysis**: Rich tracing demonstration endpoint with multiple spans
- 📊 **Full Observability**: Traces, logs, and metrics integrated
- 🌐 **Distributed Tracing**: HTTP requests → Database queries with OpenTelemetry
- 📝 **Structured Logging**: JSON logs with trace correlation
- 📈 **Prometheus Metrics**: Request rates, latency, error rates
- 🏥 **Health Checks**: Database connectivity monitoring
- 🐳 **Container-Ready**: Optimized for Kubernetes deployment

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │────│  Catalog API    │────│   PostgreSQL    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  Observability  │
                    │                 │
                    │ • Tempo (Traces)│
                    │ • Loki (Logs)   │
                    │ • Prometheus    │
                    └─────────────────┘
```

### Request Flow
1. **HTTP Request** arrives → Gin router
2. **Tracing Middleware** creates root span
3. **Logging Middleware** captures request details  
4. **Metrics Middleware** records performance data
5. **Handler** processes business logic
6. **ProductService** executes database operations with child spans
7. **Response** returned with full observability context

## 📁 Code Organization Guide

### 🗂️ Directory Structure
```
services/catalog/
├── main.go                 # 🚪 Entry point - start here
├── internal/               # 📦 Internal packages
│   ├── server/            # 🌐 HTTP server & middleware
│   ├── handlers/          # 🎯 Request handlers (API endpoints)
│   ├── services/          # 🧠 Business logic & analysis operations
│   ├── models/            # 💾 Data access & CRUD operations
│   ├── db/                # 🗄️ Database connection & schema
│   ├── metrics/           # 📊 Prometheus metrics
│   ├── tracing/           # 🔍 OpenTelemetry setup
│   └── logger/            # 📝 Structured logging
├── go.mod                 # 📋 Dependencies
└── Dockerfile             # 🐳 Container image
```

### 🔍 Where to Find What

| **Want to understand...** | **Look at...** | **Key files** |
|---------------------------|----------------|---------------|
| 🚪 **Application startup** | `main.go` | Entry point, initialization order |
| 🌐 **HTTP routing & middleware** | `internal/server/` | `server.go` - middleware stack |
| 🎯 **API endpoints** | `internal/handlers/` | `products.go`, `health.go` |
| 🧠 **Business logic & analysis** | `internal/services/` | `analysis.go` - complex operations |
| 💾 **Data access & CRUD** | `internal/models/` | `product.go` - database operations |
| 🗄️ **Database setup** | `internal/db/` | `connection.go` - DB configuration |
| 🔍 **Tracing implementation** | `internal/tracing/` | `tracing.go` - OpenTelemetry config |
| 📊 **Metrics collection** | `internal/metrics/` | `metrics.go` - Prometheus metrics |
| 📝 **Logging setup** | `internal/logger/` | `logger.go` - Structured logging |

## 🔍 Observability Features Deep Dive

### 📊 Distributed Tracing
**Implementation**: OpenTelemetry with automatic + manual instrumentation

**What to explore**:
- `internal/tracing/tracing.go` - OTLP exporter configuration
- `internal/server/server.go:52` - Automatic HTTP tracing with `otelgin`
- `internal/services/analysis.go` - Complex multi-span operations
- `internal/models/product.go` - Database span creation

**Rich Tracing with the Analyze Endpoint**:

The `/api/v1/products/analyze` endpoint demonstrates comprehensive distributed tracing:

![Analyze Endpoint Trace](../../img/analyze-trace-screenshot.png)

**Trace hierarchy you'll see**:
```
🌐 GET /api/v1/products/analyze     [HTTP span - 1518ms]
  └── 📊 product.analyze            [Analysis span - 1518ms]
      ├── 🧮 compute.analysis       [Compute span - 105ms]
      │   ├── compute.matrix_operations       [52ms]
      │   ├── compute.statistical_analysis    [31ms]
      │   └── compute.complexity_scoring      [23ms]
      ├── 🗄️ database.analysis      [Database span - 85ms]
      │   ├── db.count_products               [80ms]
      │   └── db.product_lookup               [120ms] (if ?id= provided)
      └── 🌐 external.api_call       [HTTP span - 1327ms]
```

**Key span attributes**:
- 🧮 **Compute**: calculations=3000, memory_bytes=8000, complexity_score
- 🗄️ **Database**: result_count=32, queries_executed=1, avg_latency_ms=85
- 🌐 **External**: status_code=200, response_time_ms=1327, service=httpbin

### 📝 Structured Logging
**Implementation**: JSON logs with trace correlation

**What to explore**:
- `internal/logger/logger.go` - Global logger setup
- `internal/server/server.go:76` - Request logging with trace IDs
- `internal/models/product.go` - Business logic logging

**Log format you'll see**:
```json
{
  "timestamp": "2025-01-19T10:30:45Z",
  "level": "info",
  "message": "HTTP request processed",
  "component": "http",
  "method": "GET",
  "path": "/api/v1/products/1",
  "status_code": 200,
  "duration_ms": 12,
  "trace_id": "abc123...",
  "span_id": "def456..."
}
```

### 📈 Prometheus Metrics
**Implementation**: Custom metrics with automatic collection

**What to explore**:
- `internal/metrics/metrics.go` - Metric definitions
- `internal/server/server.go:103` - Metrics middleware
- `GET /metrics` endpoint - Prometheus scraping endpoint

**Metrics you'll find**:
- `catalog_http_requests_total` - Request count by method/path/status
- `catalog_http_request_duration_seconds` - Request latency histograms
- `catalog_http_requests_in_flight` - Current active requests

## 🛠️ API Reference

### Product Endpoints
```http
GET    /api/v1/products          # List products (paginated)
POST   /api/v1/products          # Create product
GET    /api/v1/products/analyze  # Analyze products (rich tracing demo)
GET    /api/v1/products/:id      # Get specific product
PUT    /api/v1/products/:id      # Update product
DELETE /api/v1/products/:id      # Delete product
```

### System Endpoints
```http
GET    /health                   # Health check
GET    /metrics                  # Prometheus metrics
```

### Response Format
All endpoints use consistent JSON structure:
```json
{
  "data": { /* actual data */ },
  "page": 1,      // for paginated endpoints
  "limit": 50,    // for paginated endpoints
  "count": 25     // for paginated endpoints
}
```

## 🎓 Learning Guide: Exploring the Code

### 1. **Start with the Big Picture** 📖
1. Read `main.go` - understand initialization order
2. Look at `internal/server/server.go` - see middleware stack
3. Check `internal/handlers/products.go` - understand request flow

### 2. **Dive into Observability** 🔍
1. **Tracing**: Follow a request from `otelgin` middleware to database spans
2. **Logging**: See how trace IDs connect logs across components  
3. **Metrics**: Understand how request data becomes Prometheus metrics

### 3. **Understand Data Flow** 🔄
```
HTTP Request → Handler → ProductService → Database → Response
     ↓             ↓           ↓            ↓          ↓
  Tracing     Validation   Bus Logic   DB Query   Metrics
  Logging     Parsing      Tracing     Tracing    Logging
```

### 4. **Key Patterns to Notice** 💡
- **Dependency Injection**: Database passed to services
- **Middleware Layering**: Tracing → Logging → Metrics → Business Logic
- **Error Handling**: Consistent error responses with logging
- **Context Propagation**: Trace context flows through all layers

## 🔧 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DB_HOST` | `localhost` | PostgreSQL hostname |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `catalog_user` | Database username |
| `DB_PASSWORD` | `catalog_pass` | Database password |
| `DB_NAME` | `localmart` | Database name |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | (required) | OpenTelemetry collector endpoint |
| `OTEL_SERVICE_NAME` | `catalog-service` | Service name for tracing |
| `LOG_LEVEL` | `info` | Logging level |

## 📊 Observability in Action

### Viewing Traces
1. **Grafana Explore** → **Tempo** → Query: `{service.name="catalog-service"}`
2. See complete request → database trace hierarchy
3. Notice span attributes (product.id, db.table, etc.)

### Viewing Logs  
1. **Grafana Explore** → **Loki** → Query: `{app="catalog"}`
2. See structured JSON logs with trace correlation
3. Filter by trace_id to see all logs for a specific request

### Viewing Metrics
1. **Grafana Explore** → **Prometheus** 
2. Try queries like:
   - `sum(rate(catalog_http_requests_total[5m]))` - Request rate
   - `histogram_quantile(0.95, rate(catalog_http_request_duration_seconds_bucket[5m]))` - 95th percentile latency

## 🚀 Getting Started

1. **Deploy**: Service automatically deploys with Tilt in the lab environment
2. **Test**: `curl http://catalog.kubelab.lan:8081/api/v1/products`
3. **Observe**: Check traces, logs, and metrics in Grafana
4. **Explore**: Start with `main.go` and follow the request flow

## 🧪 Testing & Examples

### Prerequisites
Make sure the service is accessible by adding this to your `/etc/hosts`:
```bash
echo "127.0.0.1 catalog.kubelab.lan" | sudo tee -a /etc/hosts
```

### 🏥 Health Check
Start with the basics - verify the service is running:
```bash
# Check service health and database connectivity
curl -s http://catalog.kubelab.lan:8081/health | jq
```

Expected response:
```json
{
  "data": {
    "status": "healthy",
    "database": "connected", 
    "service": "catalog-service"
  }
}
```

### 📦 Product CRUD Operations

#### 🔍 **READ Operations**
```bash
# Get all products (shows pagination)
curl -s http://catalog.kubelab.lan:8081/api/v1/products | jq

# Get products with custom pagination  
curl -s "http://catalog.kubelab.lan:8081/api/v1/products?page=1&limit=5" | jq

# Get specific product by ID
curl -s http://catalog.kubelab.lan:8081/api/v1/products/1 | jq

# Try to get non-existent product (404 error)
curl -s http://catalog.kubelab.lan:8081/api/v1/products/999 | jq
```

#### ➕ **CREATE Operations**
```bash
# Create a MacBook Pro
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro 14\"",
    "description": "Apple MacBook Pro 14-inch with M3 chip",
    "price": 1999.99,
    "stock_quantity": 50
  }' | jq

# Create an iPhone
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Latest iPhone with titanium design and USB-C",
    "price": 999.99,
    "stock_quantity": 100
  }' | jq

# Create AirPods
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "AirPods Pro",
    "description": "Wireless earbuds with active noise cancellation",
    "price": 249.99,
    "stock_quantity": 75
  }' | jq
```

#### ✏️ **UPDATE Operations**
```bash
# Complete update of a product
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro 14\" M3",
    "description": "Apple MacBook Pro 14-inch with M3 chip - Updated model",
    "price": 1899.99,
    "stock_quantity": 45
  }' | jq

# Partial update (only price and stock)
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/2 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 949.99,
    "stock_quantity": 120
  }' | jq

# Try to update non-existent product (404 error)
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/999 \
  -H "Content-Type: application/json" \
  -d '{"name": "Non-existent Product"}' | jq
```

#### 🗑️ **DELETE Operations**
```bash
# Delete a product
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/3 | jq

# Try to delete the same product again (404 error)
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/3 | jq

# Delete non-existent product (404 error)
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/999 | jq
```

### 🔍 **ANALYZE Operations (Rich Tracing Demo)**

The analyze endpoint demonstrates complex distributed tracing with multiple spans:

```bash
# Analyze all products (creates 6+ spans)
curl -s http://catalog.kubelab.lan:8081/api/v1/products/analyze | jq

# Analyze specific product (creates 7+ spans with product lookup)
curl -s "http://catalog.kubelab.lan:8081/api/v1/products/analyze?id=1" | jq

# Example response structure:
curl -s http://catalog.kubelab.lan:8081/api/v1/products/analyze | jq '.data | keys'
# Expected: ["compute_stats", "database_stats", "external_data", "metadata", "timestamp", "total_duration_ms"]
```

**What this endpoint demonstrates:**
- 🧮 **Computation spans**: Matrix operations, statistical analysis, complexity scoring
- 🗄️ **Database spans**: Product counting, optional product lookup
- 🌐 **External API spans**: HTTP call to httpbin.org with 1s delay
- 📊 **Rich span attributes**: Calculations performed, query results, response times
- 🔗 **Span hierarchy**: Parent-child relationships across service layers

**Perfect for learning:**
- How business logic creates multiple spans
- Cross-cutting concerns (compute, database, external calls)
- Span attributes and metadata
- Performance bottleneck identification

### 🚨 Error Scenarios & Validation

Test the API's error handling and validation:

```bash
# Missing required fields
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Product without name or price"
  }' | jq

# Invalid price (negative)
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Product",
    "price": -10.00
  }' | jq

# Invalid product ID format
curl -s http://catalog.kubelab.lan:8081/api/v1/products/invalid | jq

# Malformed JSON
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Broken JSON"' | jq
```

### 🔄 Complete CRUD Workflow

Test the full lifecycle while observing traces and logs:

```bash
# 1. Create a test product
echo "🛍️ Creating product..."
PRODUCT_ID=$(curl -s -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A product for testing the complete workflow",
    "price": 99.99,
    "stock_quantity": 10
  }' | jq -r '.data.id')

echo "✅ Created product with ID: $PRODUCT_ID"

# 2. Read the product back
echo "📖 Reading product..."
curl -s http://catalog.kubelab.lan:8081/api/v1/products/$PRODUCT_ID | jq

# 3. Update the product
echo "✏️ Updating product..."
curl -s -X PUT http://catalog.kubelab.lan:8081/api/v1/products/$PRODUCT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test Product",
    "price": 89.99,
    "stock_quantity": 15
  }' | jq

# 4. Verify the update
echo "🔍 Verifying update..."
curl -s http://catalog.kubelab.lan:8081/api/v1/products/$PRODUCT_ID | jq

# 5. List all products to see our changes
echo "📋 Listing all products..."
curl -s http://catalog.kubelab.lan:8081/api/v1/products | jq '.data | length'

# 6. Delete the test product
echo "🗑️ Deleting product..."
curl -s -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/$PRODUCT_ID | jq

# 7. Verify deletion
echo "❌ Verifying deletion..."
curl -s http://catalog.kubelab.lan:8081/api/v1/products/$PRODUCT_ID | jq
```

### 📊 Observability Testing

Generate traffic to create rich observability data:

```bash
# Generate realistic traffic patterns (run in separate terminal)
for i in {1..20}; do
  # Mix of reads and writes
  curl -s http://catalog.kubelab.lan:8081/api/v1/products > /dev/null
  curl -s http://catalog.kubelab.lan:8081/api/v1/products/1 > /dev/null
  
  # Generate rich traces with analyze endpoint
  curl -s http://catalog.kubelab.lan:8081/api/v1/products/analyze > /dev/null
  curl -s "http://catalog.kubelab.lan:8081/api/v1/products/analyze?id=$((i % 5 + 1))" > /dev/null
  
  # Occasionally create/update products
  if [ $((i % 5)) -eq 0 ]; then
    curl -s -X POST http://catalog.kubelab.lan:8081/api/v1/products \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"Load Test Product $i\",
        \"price\": $((RANDOM % 100 + 50)).99,
        \"stock_quantity\": $((RANDOM % 100))
      }" > /dev/null
  fi
  
  sleep 2
done
```

**Pro tip**: Use the analyze endpoint to create the most interesting traces!
```bash
# Create 10 rich traces quickly
for i in {1..10}; do
  curl -s "http://catalog.kubelab.lan:8081/api/v1/products/analyze?id=$i" > /dev/null &
done
wait
```

### 🎯 Metrics Endpoint

Check Prometheus metrics:
```bash
# View all metrics
curl -s http://catalog.kubelab.lan:8081/metrics

# Filter for catalog-specific metrics
curl -s http://catalog.kubelab.lan:8081/metrics | grep catalog_

# Check request totals
curl -s http://catalog.kubelab.lan:8081/metrics | grep catalog_http_requests_total

# Check response time buckets
curl -s http://catalog.kubelab.lan:8081/metrics | grep catalog_http_request_duration_seconds_bucket
```

### 🚀 Traffic Simulation

For comprehensive testing, use the traffic simulation script:

```bash
# Quick test with realistic e-commerce data
./scripts/simulate-traffic.sh --duration 300 --interval 2

# Heavy load testing for observability
./scripts/simulate-traffic.sh --duration 600 --interval 1 --verbose

# Just seed the database with sample products
./scripts/simulate-traffic.sh --seed-only
```

**What this generates:**
- 📊 **Rich trace data** (HTTP → Database spans)
- 📝 **Correlated logs** with trace IDs
- 📈 **Realistic metrics** (request rates, latency, errors)
- 🛍️ **Sample product data** for testing

### 🔍 Observability Verification

After generating traffic, verify the observability stack:

1. **Traces in Grafana**: 
   - Go to Grafana → Explore → Tempo
   - Query: `{service.name="catalog-service"}`
   - See HTTP → Database trace hierarchy

2. **Logs in Grafana**:
   - Go to Grafana → Explore → Loki  
   - Query: `{app="catalog"} | json`
   - See structured logs with trace correlation

3. **Metrics in Grafana**:
   - Go to Grafana → Explore → Prometheus
   - Query: `sum(rate(catalog_http_requests_total[5m]))`
   - See request rates and latency percentiles

### 💡 Pro Tips

- **Use `jq`** for pretty JSON formatting
- **Set `-s` flag** for clean curl output
- **Run traffic simulation** before checking observability
- **Check different time ranges** in Grafana to see patterns
- **Filter logs by trace_id** to see complete request flows
- **Monitor metrics** during load testing to see real-time changes

## 💡 Next Steps for Learning

This service provides a foundation for exploring:
- **Service mesh** patterns (Istio)
- **Circuit breakers** for resilience  
- **Rate limiting** and authentication
- **Database migrations** and schema management
- **API versioning** strategies
- **Integration testing** with observability

---

**Happy exploring!** 🎉 This service demonstrates production-ready patterns in a learning-friendly codebase. Each component is designed to be educational while being functionally complete.


