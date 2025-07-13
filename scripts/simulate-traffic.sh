#!/usr/bin/env bash

# LocalMart Microservices Traffic Simulator
# Simulates realistic e-commerce traffic patterns across multiple services


# This script is a vibe code product. Not every flag might work. I have not tested everything

set -e

# Check dependencies
if ! command -v bc &> /dev/null; then
    echo "Error: 'bc' command is required but not installed."
    echo "Install with: brew install bc (macOS) or apt-get install bc (Ubuntu)"
    exit 1
fi

# Configuration
DOMAIN="${DOMAIN:-kubelab.lan}"
PORT="${PORT:-8081}"
DURATION=${DURATION:-300}  # Default 5 minutes
REQUEST_INTERVAL=${REQUEST_INTERVAL:-2}  # Seconds between requests
VERBOSE=${VERBOSE:-false}

# Service Configuration
declare -A SERVICES=(
    ["catalog"]="catalog.${DOMAIN}:${PORT}"
    ["cart"]="cart.${DOMAIN}:${PORT}"
    ["orders"]="orders.${DOMAIN}:${PORT}"
    ["users"]="users.${DOMAIN}:${PORT}"
)

# Default services to test (can be overridden)
TARGET_SERVICES=()

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Global variables
CREATED_PRODUCT_IDS=()
CREATED_USER_IDS=()
CREATED_CART_IDS=()
CREATED_ORDER_IDS=()
TOTAL_REQUESTS=0
SUCCESSFUL_REQUESTS=0
FAILED_REQUESTS=0
AVAILABLE_SERVICES=()

# Print colored output
log() {
    local color=$1
    local message=$2
    echo -e "${color}[$(date '+%H:%M:%S')] ${message}${NC}"
}

# Check if a service is available
check_service() {
    local service_name=$1
    local service_host=$2
    local base_url="http://${service_host}"
    
    log $BLUE "🔍 Checking ${service_name} service at ${base_url}"
    
    if curl -sf "${base_url}/health" > /dev/null 2>&1; then
        log $GREEN "✅ ${service_name} service is healthy"
        AVAILABLE_SERVICES+=("$service_name")
        return 0
    else
        log $YELLOW "⚠️  ${service_name} service is not available"
        return 1
    fi
}

# Check all services
check_all_services() {
    log $BLUE "🔍 Checking service availability..."
    
    for service_name in "${!SERVICES[@]}"; do
        if [[ ${#TARGET_SERVICES[@]} -eq 0 ]] || [[ " ${TARGET_SERVICES[@]} " =~ " ${service_name} " ]]; then
            check_service "$service_name" "${SERVICES[$service_name]}" || true
        fi
    done
    
    if [[ ${#AVAILABLE_SERVICES[@]} -eq 0 ]]; then
        log $RED "❌ No services are available!"
        log $YELLOW "💡 Make sure services are running and accessible via ingress"
        log $YELLOW "💡 Add these entries to /etc/hosts:"
        for service_name in "${!SERVICES[@]}"; do
            echo "   127.0.0.1 ${service_name}.${DOMAIN}"
        done
        exit 1
    fi
    
    log $GREEN "✅ Available services: ${AVAILABLE_SERVICES[*]}"
}

# Seed catalog products
seed_catalog_products() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " catalog " ]]; then
        log $YELLOW "⚠️  Skipping catalog seeding - service not available"
        return
    fi
    
    local base_url="http://${SERVICES[catalog]}"
    log $BLUE "🌱 Seeding catalog products..."
    
    # E-commerce product catalog
    local products=(
        '{"name":"MacBook Pro 14\"","description":"Apple MacBook Pro 14-inch with M3 chip, 16GB RAM, 512GB SSD","price":1999.99,"stock_quantity":25}'
        '{"name":"iPhone 15 Pro","description":"Latest iPhone with titanium design and USB-C","price":999.99,"stock_quantity":50}'
        '{"name":"AirPods Pro","description":"Wireless earbuds with active noise cancellation","price":249.99,"stock_quantity":100}'
        '{"name":"iPad Air","description":"Powerful tablet with M2 chip and 10.9-inch display","price":599.99,"stock_quantity":75}'
        '{"name":"Apple Watch Series 9","description":"Advanced smartwatch with health monitoring","price":399.99,"stock_quantity":40}'
        '{"name":"Samsung Galaxy S24 Ultra","description":"Premium Android smartphone with S Pen and 200MP camera","price":1199.99,"stock_quantity":30}'
        '{"name":"Sony WH-1000XM5","description":"Industry-leading noise canceling wireless headphones","price":399.99,"stock_quantity":60}'
        '{"name":"Dell XPS 13","description":"13-inch ultrabook with Intel Core i7 and 16GB RAM","price":1299.99,"stock_quantity":20}'
        '{"name":"Nintendo Switch OLED","description":"Gaming console with vibrant OLED screen","price":349.99,"stock_quantity":45}'
        '{"name":"Kindle Paperwhite","description":"Waterproof e-reader with 6.8-inch display","price":139.99,"stock_quantity":80}'
        '{"name":"Instant Pot Duo 7-in-1","description":"Multi-use pressure cooker and slow cooker","price":99.99,"stock_quantity":65}'
        '{"name":"Dyson V15 Detect","description":"Cordless vacuum with laser dust detection","price":749.99,"stock_quantity":15}'
        '{"name":"Nespresso Vertuo Next","description":"Coffee and espresso machine with centrifusion technology","price":199.99,"stock_quantity":35}'
        '{"name":"KitchenAid Stand Mixer","description":"Professional 5-quart stand mixer in multiple colors","price":429.99,"stock_quantity":25}'
        '{"name":"Philips Hue Smart Bulbs 4-Pack","description":"Color-changing LED smart bulbs","price":199.99,"stock_quantity":90}'
        '{"name":"Peloton Bike+","description":"Indoor exercise bike with rotating HD touchscreen","price":2495.00,"stock_quantity":8}'
        '{"name":"Fitbit Charge 5","description":"Advanced fitness tracker with built-in GPS","price":179.99,"stock_quantity":55}'
        '{"name":"Yeti Rambler 30oz","description":"Stainless steel tumbler with MagSlider lid","price":39.99,"stock_quantity":120}'
        '{"name":"Nike Air Max 270","description":"Mens running shoes with Max Air heel unit","price":149.99,"stock_quantity":75}'
        '{"name":"Hydro Flask 32oz","description":"Insulated stainless steel water bottle","price":44.99,"stock_quantity":85}'        
        '{"name":"The Thursday Murder Club","description":"Bestselling mystery novel by Richard Osman","price":16.99,"stock_quantity":150}'
        '{"name":"Atomic Habits","description":"Life-changing guide to building good habits","price":18.99,"stock_quantity":200}'
        '{"name":"PlayStation 5","description":"Next-generation gaming console","price":499.99,"stock_quantity":12}'
        '{"name":"Meta Quest 3","description":"Mixed reality VR headset with 128GB storage","price":499.99,"stock_quantity":18}'
        '{"name":"Levis 501 Original Jeans","description":"Classic straight-leg denim jeans","price":89.99,"stock_quantity":110}'
        '{"name":"Ray-Ban Aviator Sunglasses","description":"Classic gold-frame aviator sunglasses","price":154.99,"stock_quantity":40}'
        '{"name":"Patagonia Better Sweater","description":"Recycled fleece pullover jacket","price":139.99,"stock_quantity":60}'
        '{"name":"Allbirds Tree Runners","description":"Sustainable running shoes made from eucalyptus","price":98.99,"stock_quantity":70}'
        '{"name":"Casper Original Mattress Queen","description":"Premium memory foam mattress with zoned support","price":1095.00,"stock_quantity":10}'
        '{"name":"Roomba j7+","description":"Self-emptying robot vacuum with smart mapping","price":649.99,"stock_quantity":22}'
        '{"name":"Ring Video Doorbell Pro","description":"1080p HD wireless doorbell with two-way talk","price":249.99,"stock_quantity":45}'
        '{"name":"Nest Learning Thermostat","description":"Smart thermostat that learns your schedule","price":249.99,"stock_quantity":35}'
    )
    
    for product in "${products[@]}"; do
        local response=$(curl -s -w "\n%{http_code}" -X POST "${base_url}/api/v1/products" \
            -H "Content-Type: application/json" \
            -d "$product" 2>/dev/null)
        
        local http_code=$(echo "$response" | tail -n1)
        local body=$(echo "$response" | head -n -1)
        
        TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
        
        if [[ "$http_code" == "201" ]]; then
            local product_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
            CREATED_PRODUCT_IDS+=($product_id)
            SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
            
            local product_name=$(echo "$product" | grep -o '"name":"[^"]*"' | cut -d':' -f2 | tr -d '"')
            log $GREEN "📦 Created product: $product_name (ID: $product_id)"
        else
            FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
            log $RED "❌ Failed to create product: HTTP $http_code"
        fi
        
        sleep 0.2
    done
    
    log $GREEN "🌱 Seeded ${#CREATED_PRODUCT_IDS[@]} products"
}

# Seed user accounts
seed_user_accounts() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " users " ]]; then
        log $YELLOW "⚠️  Skipping user seeding - service not available"
        return
    fi
    
    local base_url="http://${SERVICES[users]}"
    log $BLUE "🌱 Seeding user accounts..."
    
    local users=(
        '{"name":"Alice Johnson","email":"alice@example.com","password":"password123"}'
        '{"name":"Bob Smith","email":"bob@example.com","password":"password123"}'
        '{"name":"Carol Davis","email":"carol@example.com","password":"password123"}'
        '{"name":"David Wilson","email":"david@example.com","password":"password123"}'
        '{"name":"Eve Brown","email":"eve@example.com","password":"password123"}'
    )
    
    for user in "${users[@]}"; do
        local response=$(curl -s -w "\n%{http_code}" -X POST "${base_url}/api/v1/users" \
            -H "Content-Type: application/json" \
            -d "$user" 2>/dev/null)
        
        local http_code=$(echo "$response" | tail -n1)
        local body=$(echo "$response" | head -n -1)
        
        TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
        
        if [[ "$http_code" == "201" ]]; then
            local user_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
            CREATED_USER_IDS+=($user_id)
            SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
            
            local user_name=$(echo "$user" | grep -o '"name":"[^"]*"' | cut -d':' -f2 | tr -d '"')
            log $GREEN "👤 Created user: $user_name (ID: $user_id)"
        else
            FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
            log $RED "❌ Failed to create user: HTTP $http_code"
        fi
        
        sleep 0.2
    done
    
    log $GREEN "🌱 Seeded ${#CREATED_USER_IDS[@]} users"
}

# Get random product ID
get_random_product_id() {
    if [[ ${#CREATED_PRODUCT_IDS[@]} -eq 0 ]]; then
        echo "1"
    else
        echo "${CREATED_PRODUCT_IDS[$RANDOM % ${#CREATED_PRODUCT_IDS[@]}]}"
    fi
}

# Get random user ID
get_random_user_id() {
    if [[ ${#CREATED_USER_IDS[@]} -eq 0 ]]; then
        echo "1"
    else
        echo "${CREATED_USER_IDS[$RANDOM % ${#CREATED_USER_IDS[@]}]}"
    fi
}

# Catalog Service Simulations
simulate_catalog_browse() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " catalog " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[catalog]}"
    local page=$((RANDOM % 3 + 1))
    local limit=$((RANDOM % 15 + 5))
    
    local response=$(curl -s -w "\n%{http_code}" -X GET "${base_url}/api/v1/products?page=${page}&limit=${limit}" 2>/dev/null)
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $GREEN "👀 [CATALOG] Browsed products (page: $page, limit: $limit)"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ [CATALOG] Browse failed: HTTP $http_code"
    fi
}

simulate_catalog_view() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " catalog " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[catalog]}"
    local product_id=$(get_random_product_id)
    
    local response=$(curl -s -w "\n%{http_code}" -X GET "${base_url}/api/v1/products/${product_id}" 2>/dev/null)
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $GREEN "👁️  [CATALOG] Viewed product ID: $product_id"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $YELLOW "⚠️  [CATALOG] Product not found: ID $product_id"
    fi
}

simulate_catalog_stock_update() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " catalog " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[catalog]}"
    local product_id=$(get_random_product_id)
    local new_stock=$((RANDOM % 200 + 10))
    
    local update_data="{\"stock_quantity\":${new_stock}}"
    
    local response=$(curl -s -w "\n%{http_code}" -X PUT "${base_url}/api/v1/products/${product_id}" \
        -H "Content-Type: application/json" \
        -d "$update_data" 2>/dev/null)
    
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $BLUE "📦 [CATALOG] Updated stock for product $product_id: $new_stock units"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ [CATALOG] Stock update failed: HTTP $http_code"
    fi
}

# Cart Service Simulations
simulate_cart_add_item() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " cart " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[cart]}"
    local user_id=$(get_random_user_id)
    local product_id=$(get_random_product_id)
    local quantity=$((RANDOM % 5 + 1))
    
    local cart_data="{\"user_id\":${user_id},\"product_id\":${product_id},\"quantity\":${quantity}}"
    
    local response=$(curl -s -w "\n%{http_code}" -X POST "${base_url}/api/v1/cart/items" \
        -H "Content-Type: application/json" \
        -d "$cart_data" 2>/dev/null)
    
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "201" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $PURPLE "🛒 [CART] Added product $product_id to cart (user: $user_id, qty: $quantity)"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ [CART] Add to cart failed: HTTP $http_code"
    fi
}

simulate_cart_view() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " cart " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[cart]}"
    local user_id=$(get_random_user_id)
    
    local response=$(curl -s -w "\n%{http_code}" -X GET "${base_url}/api/v1/cart/${user_id}" 2>/dev/null)
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $PURPLE "👀 [CART] Viewed cart for user $user_id"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $YELLOW "⚠️  [CART] Cart not found for user $user_id"
    fi
}

# Orders Service Simulations
simulate_order_create() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " orders " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[orders]}"
    local user_id=$(get_random_user_id)
    
    local order_data="{\"user_id\":${user_id},\"shipping_address\":\"123 Test St, Test City, TC 12345\"}"
    
    local response=$(curl -s -w "\n%{http_code}" -X POST "${base_url}/api/v1/orders" \
        -H "Content-Type: application/json" \
        -d "$order_data" 2>/dev/null)
    
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n -1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "201" ]]; then
        local order_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
        CREATED_ORDER_IDS+=($order_id)
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $CYAN "📋 [ORDERS] Created order $order_id for user $user_id"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ [ORDERS] Order creation failed: HTTP $http_code"
    fi
}

simulate_order_view() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " orders " ]] || [[ ${#CREATED_ORDER_IDS[@]} -eq 0 ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[orders]}"
    local order_id="${CREATED_ORDER_IDS[$RANDOM % ${#CREATED_ORDER_IDS[@]}]}"
    
    local response=$(curl -s -w "\n%{http_code}" -X GET "${base_url}/api/v1/orders/${order_id}" 2>/dev/null)
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $CYAN "👁️  [ORDERS] Viewed order $order_id"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $YELLOW "⚠️  [ORDERS] Order not found: ID $order_id"
    fi
}

# Users Service Simulations
simulate_user_login() {
    if [[ ! " ${AVAILABLE_SERVICES[@]} " =~ " users " ]]; then
        return
    fi
    
    local base_url="http://${SERVICES[users]}"
    local emails=("alice@example.com" "bob@example.com" "carol@example.com" "david@example.com" "eve@example.com")
    local email="${emails[$RANDOM % ${#emails[@]}]}"
    
    local login_data="{\"email\":\"${email}\",\"password\":\"password123\"}"
    
    local response=$(curl -s -w "\n%{http_code}" -X POST "${base_url}/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "$login_data" 2>/dev/null)
    
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $BLUE "🔐 [USERS] User logged in: $email"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ [USERS] Login failed for $email: HTTP $http_code"
    fi
}

# Cross-service workflow: Complete shopping journey
simulate_shopping_journey() {
    if [[ " ${AVAILABLE_SERVICES[@]} " =~ " catalog " ]] && [[ " ${AVAILABLE_SERVICES[@]} " =~ " cart " ]]; then
        # 1. Browse products
        simulate_catalog_browse
        sleep 1
        
        # 2. View specific product
        simulate_catalog_view
        sleep 1
        
        # 3. Add to cart
        simulate_cart_add_item
        sleep 1
        
        # 4. View cart
        simulate_cart_view
        
        # 5. Maybe create order if orders service is available
        if [[ " ${AVAILABLE_SERVICES[@]} " =~ " orders " ]]; then
            sleep 1
            simulate_order_create
        fi
        
        log $GREEN "🛍️  [JOURNEY] Completed shopping journey"
    fi
}

# Main traffic simulation
simulate_traffic() {
    local end_time=$(($(date +%s) + DURATION))
    
    log $BLUE "🚀 Starting traffic simulation for ${DURATION} seconds..."
    log $BLUE "📊 Request interval: ${REQUEST_INTERVAL} seconds"
    log $BLUE "🎯 Target services: ${AVAILABLE_SERVICES[*]}"
    
    while [[ $(date +%s) -lt $end_time ]]; do
        local action=$((RANDOM % 100))
        
        if [[ $action -lt 30 ]]; then
            # 30% - Browse catalog
            simulate_catalog_browse
        elif [[ $action -lt 50 ]]; then
            # 20% - View specific product
            simulate_catalog_view
        elif [[ $action -lt 65 ]]; then
            # 15% - Add to cart
            simulate_cart_add_item
        elif [[ $action -lt 75 ]]; then
            # 10% - View cart
            simulate_cart_view
        elif [[ $action -lt 85 ]]; then
            # 10% - Stock updates
            simulate_catalog_stock_update
        elif [[ $action -lt 90 ]]; then
            # 5% - Create order
            simulate_order_create
        elif [[ $action -lt 95 ]]; then
            # 5% - User login
            simulate_user_login
        elif [[ $action -lt 98 ]]; then
            # 3% - View order
            simulate_order_view
        else
            # 2% - Complete shopping journey
            simulate_shopping_journey
        fi
        
        sleep $REQUEST_INTERVAL
    done
}

# Seed all services
seed_all_services() {
    log $BLUE "🌱 Seeding all available services..."
    seed_catalog_products
    seed_user_accounts
    log $GREEN "🌱 Seeding complete!"
}

# Print statistics
print_stats() {
    log $BLUE "📈 Traffic Simulation Complete!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📊 Statistics:"
    echo "   Services Tested: ${#AVAILABLE_SERVICES[@]} (${AVAILABLE_SERVICES[*]})"
    echo "   Total Requests: $TOTAL_REQUESTS"
    echo "   Successful:     $SUCCESSFUL_REQUESTS"
    echo "   Failed:         $FAILED_REQUESTS"
    if [[ $TOTAL_REQUESTS -gt 0 ]]; then
        echo "   Success Rate:   $(echo "scale=2; $SUCCESSFUL_REQUESTS * 100 / $TOTAL_REQUESTS" | bc 2>/dev/null || echo "N/A")%"
    else
        echo "   Success Rate:   N/A%"
    fi
    echo "   Products Created: ${#CREATED_PRODUCT_IDS[@]}"
    echo "   Users Created:    ${#CREATED_USER_IDS[@]}"
    echo "   Orders Created:   ${#CREATED_ORDER_IDS[@]}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Show usage
show_usage() {
    echo "LocalMart Microservices Traffic Simulator"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --duration SECONDS        Duration to run simulation (default: 300)"
    echo "  --interval SECONDS        Interval between requests (default: 2)"
    echo "  --services SERVICE1,SERVICE2  Specific services to test (default: all available)"
    echo "  --domain DOMAIN           Domain suffix (default: kubelab.lan)"
    echo "  --port PORT               Port for services (default: 8081)"
    echo "  --verbose                 Show detailed output"
    echo "  --seed-only               Only seed data, don't simulate traffic"
    echo "  --no-seed                 Skip seeding, only simulate traffic"
    echo "  --help                    Show this help message"
    echo ""
    echo "Available Services:"
    echo "  catalog    - Product catalog service (Go)"
    echo "  cart       - Shopping cart service (Python/FastAPI)"
    echo "  orders     - Order management service (Java/Spring Boot)"
    echo "  users      - User management service (Node.js/TypeScript)"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Test all available services"
    echo "  $0 --services catalog,cart           # Test only catalog and cart"
    echo "  $0 --duration 600 --interval 1       # 10 minutes, 1 second intervals"
    echo "  $0 --seed-only                       # Only seed data"
    echo "  $0 --no-seed --duration 120          # Skip seeding, 2 minutes traffic"
    echo ""
    echo "Prerequisites:"
    echo "  Add these entries to /etc/hosts:"
    echo "    127.0.0.1 catalog.kubelab.lan"
    echo "    127.0.0.1 cart.kubelab.lan"
    echo "    127.0.0.1 orders.kubelab.lan"
    echo "    127.0.0.1 users.kubelab.lan"
}

# Parse command line arguments
SEED_ONLY=false
NO_SEED=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --duration)
            DURATION="$2"
            shift 2
            ;;
        --interval)
            REQUEST_INTERVAL="$2"
            shift 2
            ;;
        --services)
            IFS=',' read -ra TARGET_SERVICES <<< "$2"
            shift 2
            ;;
        --domain)
            DOMAIN="$2"
            # Update services with new domain
            SERVICES["catalog"]="catalog.${DOMAIN}:${PORT}"
            SERVICES["cart"]="cart.${DOMAIN}:${PORT}"
            SERVICES["orders"]="orders.${DOMAIN}:${PORT}"
            SERVICES["users"]="users.${DOMAIN}:${PORT}"
            shift 2
            ;;
        --port)
            PORT="$2"
            # Update services with new port
            SERVICES["catalog"]="catalog.${DOMAIN}:${PORT}"
            SERVICES["cart"]="cart.${DOMAIN}:${PORT}"
            SERVICES["orders"]="orders.${DOMAIN}:${PORT}"
            SERVICES["users"]="users.${DOMAIN}:${PORT}"
            shift 2
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --seed-only)
            SEED_ONLY=true
            shift
            ;;
        --no-seed)
            NO_SEED=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main execution
main() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log $BLUE "🎯 LocalMart Microservices Traffic Simulator"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    check_all_services
    
    if [[ "$NO_SEED" != "true" ]]; then
        seed_all_services
    fi
    
    if [[ "$SEED_ONLY" != "true" ]]; then
        simulate_traffic
    fi
    
    print_stats
}

# Trap to handle interruption
trap 'log $YELLOW "🛑 Simulation interrupted"; print_stats; exit 0' INT

# Run main function
main 