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


