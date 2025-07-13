#!/bin/bash

# Catalog Service Traffic Simulator
# Simulates realistic e-commerce traffic patterns

set -e

# Configuration
CATALOG_HOST="${CATALOG_HOST:-catalog.kubelab.lan:8081}"
BASE_URL="http://${CATALOG_HOST}"
DURATION=${DURATION:-300}  # Default 5 minutes
REQUEST_INTERVAL=${REQUEST_INTERVAL:-2}  # Seconds between requests
VERBOSE=${VERBOSE:-false}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Global variables
CREATED_PRODUCT_IDS=()
TOTAL_REQUESTS=0
SUCCESSFUL_REQUESTS=0
FAILED_REQUESTS=0

# Print colored output
log() {
    local color=$1
    local message=$2
    echo -e "${color}[$(date '+%H:%M:%S')] ${message}${NC}"
}

# Check if service is available
check_service() {
    log $BLUE "🔍 Checking if catalog service is available at ${BASE_URL}"
    
    if ! curl -sf "${BASE_URL}/health" > /dev/null; then
        log $RED "❌ Catalog service is not available at ${BASE_URL}/health"
        log $YELLOW "💡 Make sure the service is running and accessible via ingress"
        log $YELLOW "💡 Ensure you have '127.0.0.1 catalog.kubelab.lan' in /etc/hosts"
        exit 1
    fi
    
    log $GREEN "✅ Catalog service is healthy"
}

# Seed initial product data
seed_products() {
    log $BLUE "🌱 Seeding initial product data..."
    
    # E-commerce product catalog
    local products=(
        '{"name":"MacBook Pro 14\"","description":"Apple MacBook Pro 14-inch with M3 chip, 16GB RAM, 512GB SSD","price":1999.99,"stock_quantity":25}'
        '{"name":"iPhone 15 Pro","description":"Latest iPhone with titanium design and USB-C","price":999.99,"stock_quantity":50}'
        '{"name":"AirPods Pro","description":"Wireless earbuds with active noise cancellation","price":249.99,"stock_quantity":100}'
        '{"name":"iPad Air","description":"Powerful tablet with M2 chip and 10.9-inch display","price":599.99,"stock_quantity":75}'
        '{"name":"Apple Watch Series 9","description":"Advanced smartwatch with health monitoring","price":399.99,"stock_quantity":40}'
        '{"name":"Magic Keyboard","description":"Wireless keyboard for Mac with Touch ID","price":179.99,"stock_quantity":60}'
        '{"name":"Magic Mouse","description":"Multi-touch wireless mouse for Mac","price":79.99,"stock_quantity":80}'
        '{"name":"Studio Display","description":"27-inch 5K Retina display for creative professionals","price":1599.99,"stock_quantity":15}'
        '{"name":"HomePod mini","description":"Smart speaker with amazing sound and Siri","price":99.99,"stock_quantity":90}'
        '{"name":"AirTag 4-pack","description":"Find your things with precision tracking","price":99.99,"stock_quantity":120}'
        '{"name":"USB-C to Lightning Cable","description":"Fast charging cable for iPhone","price":19.99,"stock_quantity":200}'
        '{"name":"MagSafe Charger","description":"Wireless charger for iPhone with magnetic alignment","price":39.99,"stock_quantity":150}'
        '{"name":"Apple TV 4K","description":"Stream movies, shows, and games in 4K HDR","price":179.99,"stock_quantity":35}'
        '{"name":"Mac Studio","description":"High-performance desktop with M2 Ultra chip","price":3999.99,"stock_quantity":8}'
        '{"name":"Pro Display XDR","description":"32-inch Retina 6K display for professionals","price":4999.99,"stock_quantity":5}'
    )
    
    for product in "${products[@]}"; do
        local response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/v1/products" \
            -H "Content-Type: application/json" \
            -d "$product")
        
        local http_code=$(echo "$response" | tail -n1)
        local body=$(echo "$response" | head -n -1)
        
        TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
        
        if [[ "$http_code" == "201" ]]; then
            local product_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
            CREATED_PRODUCT_IDS+=($product_id)
            SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
            
            local product_name=$(echo "$product" | grep -o '"name":"[^"]*"' | cut -d':' -f2 | tr -d '"')
            log $GREEN "✅ Created: $product_name (ID: $product_id)"
        else
            FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
            log $RED "❌ Failed to create product: HTTP $http_code"
            [[ "$VERBOSE" == "true" ]] && echo "$body"
        fi
        
        sleep 0.5  # Small delay between seeding requests
    done
    
    log $GREEN "🌱 Seeded ${#CREATED_PRODUCT_IDS[@]} products successfully"
}

# Get random product ID from created products
get_random_product_id() {
    if [[ ${#CREATED_PRODUCT_IDS[@]} -eq 0 ]]; then
        echo "1"  # Fallback
    else
        echo "${CREATED_PRODUCT_IDS[$RANDOM % ${#CREATED_PRODUCT_IDS[@]}]}"
    fi
}

# Simulate browsing products (GET all)
simulate_browse_products() {
    local page=$((RANDOM % 3 + 1))  # Random page 1-3
    local limit=$((RANDOM % 15 + 5))  # Random limit 5-20
    
    local response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/api/v1/products?page=${page}&limit=${limit}")
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $GREEN "👀 Browsed products (page: $page, limit: $limit)"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ Browse failed: HTTP $http_code"
    fi
}

# Simulate viewing specific product (GET by ID)
simulate_view_product() {
    local product_id=$(get_random_product_id)
    
    local response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}/api/v1/products/${product_id}")
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $GREEN "👁️  Viewed product ID: $product_id"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $YELLOW "⚠️  Product not found: ID $product_id"
    fi
}

# Simulate stock update (PUT)
simulate_stock_update() {
    local product_id=$(get_random_product_id)
    local new_stock=$((RANDOM % 200 + 10))  # Random stock 10-210
    
    local update_data="{\"stock_quantity\":${new_stock}}"
    
    local response=$(curl -s -w "\n%{http_code}" -X PUT "${BASE_URL}/api/v1/products/${product_id}" \
        -H "Content-Type: application/json" \
        -d "$update_data")
    
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $BLUE "📦 Updated stock for product $product_id: $new_stock units"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ Stock update failed: HTTP $http_code"
    fi
}

# Simulate price update (PUT)
simulate_price_update() {
    local product_id=$(get_random_product_id)
    local price_change=$((RANDOM % 200 - 100))  # Random change -100 to +100
    local new_price=$(echo "scale=2; 99.99 + $price_change" | bc)
    
    # Ensure price is positive
    if (( $(echo "$new_price < 10" | bc -l) )); then
        new_price="19.99"
    fi
    
    local update_data="{\"price\":${new_price}}"
    
    local response=$(curl -s -w "\n%{http_code}" -X PUT "${BASE_URL}/api/v1/products/${product_id}" \
        -H "Content-Type: application/json" \
        -d "$update_data")
    
    local http_code=$(echo "$response" | tail -n1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "200" ]]; then
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $BLUE "💰 Updated price for product $product_id: \$${new_price}"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ Price update failed: HTTP $http_code"
    fi
}

# Simulate adding new product (POST)
simulate_add_product() {
    local adjectives=("Premium" "Pro" "Advanced" "Ultimate" "Elite" "Special" "Limited" "New")
    local products=("Accessory" "Device" "Gadget" "Tool" "Case" "Stand" "Adapter" "Kit")
    
    local adj="${adjectives[$RANDOM % ${#adjectives[@]}]}"
    local prod="${products[$RANDOM % ${#products[@]}]}"
    local name="$adj $prod $(date +%s)"
    local price=$(echo "scale=2; $RANDOM % 500 + 50" | bc)
    local stock=$((RANDOM % 100 + 10))
    
    local product_data="{\"name\":\"$name\",\"description\":\"Dynamically generated product for testing\",\"price\":$price,\"stock_quantity\":$stock}"
    
    local response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}/api/v1/products" \
        -H "Content-Type: application/json" \
        -d "$product_data")
    
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n -1)
    
    TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
    
    if [[ "$http_code" == "201" ]]; then
        local product_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
        CREATED_PRODUCT_IDS+=($product_id)
        SUCCESSFUL_REQUESTS=$((SUCCESSFUL_REQUESTS + 1))
        log $GREEN "➕ Added new product: $name (ID: $product_id)"
    else
        FAILED_REQUESTS=$((FAILED_REQUESTS + 1))
        log $RED "❌ Failed to add product: HTTP $http_code"
    fi
}

# Simulate realistic traffic pattern
simulate_traffic() {
    local end_time=$(($(date +%s) + DURATION))
    
    log $BLUE "🚀 Starting traffic simulation for ${DURATION} seconds..."
    log $BLUE "📊 Request interval: ${REQUEST_INTERVAL} seconds"
    
    while [[ $(date +%s) -lt $end_time ]]; do
        # Weighted random selection (realistic e-commerce patterns)
        local action=$((RANDOM % 100))
        
        if [[ $action -lt 45 ]]; then
            # 45% - Browse products (most common)
            simulate_browse_products
        elif [[ $action -lt 75 ]]; then
            # 30% - View specific product
            simulate_view_product
        elif [[ $action -lt 90 ]]; then
            # 15% - Stock updates (inventory management)
            simulate_stock_update
        elif [[ $action -lt 95 ]]; then
            # 5% - Price updates
            simulate_price_update
        else
            # 5% - Add new products
            simulate_add_product
        fi
        
        sleep $REQUEST_INTERVAL
    done
}

# Print statistics
print_stats() {
    log $BLUE "📈 Traffic Simulation Complete!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📊 Statistics:"
    echo "   Total Requests: $TOTAL_REQUESTS"
    echo "   Successful:     $SUCCESSFUL_REQUESTS"
    echo "   Failed:         $FAILED_REQUESTS"
    echo "   Success Rate:   $(echo "scale=2; $SUCCESSFUL_REQUESTS * 100 / $TOTAL_REQUESTS" | bc)%"
    echo "   Products Created: ${#CREATED_PRODUCT_IDS[@]}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Show usage
show_usage() {
    echo "Catalog Service Traffic Simulator"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --duration SECONDS    Duration to run simulation (default: 300)"
    echo "  --interval SECONDS    Interval between requests (default: 2)"
    echo "  --host HOST:PORT      Catalog service host (default: catalog.kubelab.lan:8081)"
    echo "  --verbose             Show detailed output"
    echo "  --seed-only           Only seed data, don't simulate traffic"
    echo "  --no-seed             Skip seeding, only simulate traffic"
    echo "  --help                Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Run with defaults"
    echo "  $0 --duration 600 --interval 1       # 10 minutes, 1 second intervals"
    echo "  $0 --seed-only                       # Only create initial products"
    echo "  $0 --no-seed --duration 120          # Skip seeding, 2 minutes traffic"
    echo ""
    echo "Environment Variables:"
    echo "  CATALOG_HOST          Override catalog service host"
    echo "  DURATION              Override simulation duration"
    echo "  REQUEST_INTERVAL      Override request interval"
    echo "  VERBOSE               Set to 'true' for verbose output"
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
        --host)
            CATALOG_HOST="$2"
            BASE_URL="http://${CATALOG_HOST}"
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
    log $BLUE "🎯 Catalog Service Traffic Simulator"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    check_service
    
    if [[ "$NO_SEED" != "true" ]]; then
        seed_products
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