# Catalog Service

**Technology**: Go  
**Communication**: HTTP/REST API  
**Database**: PostgreSQL  
**Purpose**: Product catalog management with search and inventory tracking

## 📋 Service Responsibilities

### Core Functions (Phase 1)
- **Product Management**: Basic CRUD operations for products
- **Product Data**: Name, description, price, SKU, stock quantity
- **Health Monitoring**: Service health and readiness checks

## 🌐 API Design

### Endpoints

**Products (Phase 1)**
```
GET    /api/v1/products           # List all products (basic pagination)
GET    /api/v1/products/:id       # Get product by ID
POST   /api/v1/products           # Create new product
PUT    /api/v1/products/:id       # Update product
DELETE /api/v1/products/:id       # Delete product
```

**Health & Metrics**
```
GET    /health                    # Health check endpoint
GET    /metrics                   # Prometheus metrics
```

**Future Endpoints (Learning Opportunities)**
```
GET    /api/v1/products/search    # Search products with filters
GET    /api/v1/categories         # Category management
POST   /api/v1/categories         # Category creation
```

## 📊 Data Models

### Product (Phase 1)
```json
{
  "id": 1,
  "name": "Product Name",
  "description": "Product description",
  "price": 29.99,
  "stock_quantity": 100
}
```


## 🗄️ Database Schema

### Phase 1 Schema (Minimal)
```sql
-- Products table (bare minimum)
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INTEGER DEFAULT 0
);
```

## Curl commands for the API

### Health Check
```bash
# Check service health
curl -X GET http://catalog.kubelab.lan:8081/health
```

### Create Products
```bash
# Create a new product
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "Apple MacBook Pro 14-inch with M3 chip",
    "price": 1999.99,
    "stock_quantity": 50
  }'

# Create another product
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone with USB-C",
    "price": 799.99,
    "stock_quantity": 100
  }'

# Create a third product
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "AirPods Pro",
    "description": "Wireless earbuds with active noise cancellation",
    "price": 249.99,
    "stock_quantity": 75
  }'
```

### Read Products
```bash
# Get all products (default pagination)
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products

# Get all products with pagination
curl -X GET "http://catalog.kubelab.lan:8081/api/v1/products?page=1&limit=5"

# Get a specific product by ID
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products/1

# Get a specific product by ID (non-existent)
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products/999
```

### Update Products
```bash
# Update a product completely
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro M3",
    "description": "Apple MacBook Pro 14-inch with M3 chip - Updated",
    "price": 1899.99,
    "stock_quantity": 45
  }'

# Partial update (only price and stock)
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/2 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 749.99,
    "stock_quantity": 120
  }'

# Update non-existent product
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/999 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Non-existent Product"
  }'
```

### Delete Products
```bash
# Delete a product
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/3

# Try to delete the same product again (should fail)
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/3

# Delete non-existent product
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/999
```

### Test Invalid Requests
```bash
# Create product with missing required fields
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Product without name or price"
  }'

# Create product with invalid price
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Product",
    "price": -10.00
  }'

# Get product with invalid ID
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products/invalid
```

### Pretty Print JSON Responses
```bash
# Use jq for pretty printing (if available)
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products | jq '.'

# Or use python for pretty printing
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products | python -m json.tool
```

### Test Complete CRUD Flow
```bash
# 1. Create a product
echo "1. Creating product..."
curl -X POST http://catalog.kubelab.lan:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A product for testing",
    "price": 99.99,
    "stock_quantity": 10
  }'

# 2. Get all products
echo -e "\n2. Getting all products..."
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products

# 3. Get specific product (assuming ID 1)
echo -e "\n3. Getting product with ID 1..."
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products/1

# 4. Update the product
echo -e "\n4. Updating product..."
curl -X PUT http://catalog.kubelab.lan:8081/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test Product",
    "price": 89.99
  }'

# 5. Verify update
echo -e "\n5. Verifying update..."
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products/1

# 6. Delete the product
echo -e "\n6. Deleting product..."
curl -X DELETE http://catalog.kubelab.lan:8081/api/v1/products/1

# 7. Verify deletion
echo -e "\n7. Verifying deletion..."
curl -X GET http://catalog.kubelab.lan:8081/api/v1/products/1
```

### Setup Instructions
Before testing, make sure to add the following to your `/etc/hosts` file:
```bash
echo "127.0.0.1 catalog.kubelab.lan" | sudo tee -a /etc/hosts
```

Or manually add this line to `/etc/hosts`:
```
127.0.0.1 catalog.kubelab.lan
```

## 🚀 Traffic Simulation

For automated testing and realistic traffic patterns, use the traffic simulation script:

```bash
# Quick data seeding and 5 minutes of realistic traffic
./scripts/simulate-traffic.sh

# Just seed the catalog with 15 realistic products
./scripts/simulate-traffic.sh --seed-only

# Heavy traffic simulation for observability testing
./scripts/simulate-traffic.sh --duration 600 --interval 1 --verbose
```

**Features:**
- Seeds 15 realistic Apple products (MacBook, iPhone, iPad, etc.)
- Simulates authentic e-commerce traffic patterns (75% reads, 25% writes)
- Real-time colored console output
- Perfect for testing observability stack (Loki, Prometheus, Grafana)
- Configurable duration, request intervals, and endpoints

See [`scripts/README.md`](scripts/README.md) for complete documentation and usage examples.

## 🎯 Phase 1 Scope

**Must Have (MVP):**
- [ ] Basic product CRUD operations (create, read, update, delete)
- [ ] PostgreSQL integration with simple schema
- [ ] Structured JSON logging (compatible with Loki)
- [ ] Health check endpoint (`/health`)
- [ ] Basic Prometheus metrics


**Phase 1 Success Criteria:**
- [ ] Service starts and connects to PostgreSQL
- [ ] Can create, read, update, and delete products
- [ ] All operations logged as structured JSON
- [ ] Health endpoint returns service status
- [ ] Prometheus metrics available at `/metrics`
- [ ] Tracing


