#!/bin/bash

# Documentation Sync Script
# Analyzes changes and suggests specific documentation updates
# Designed for token efficiency - tells AI exactly what to check

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_section() {
    echo -e "\n${BLUE}## $1${NC}"
}

print_action() {
    echo -e "${YELLOW}→ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${RED}⚠️  $1${NC}"
}

# Get the comparison point (last tag or specific commit)
get_comparison_point() {
    local comparison_point=""
    
    if [ "$1" = "--since-tag" ]; then
        # Compare since last git tag
        comparison_point=$(git describe --tags --abbrev=0 2>/dev/null || echo "HEAD~10")
    elif [ -n "$1" ]; then
        # Use provided commit/branch
        comparison_point="$1"
    else
        # Compare since last 5 commits (safe default)
        comparison_point="HEAD~5"
    fi
    
    echo "$comparison_point"
}

# Analyze what files have changed
analyze_changes() {
    local since="$1"
    print_section "Change Analysis (since $since)"
    
    # Get changed files
    local changed_files=$(git diff --name-only "$since" 2>/dev/null || git diff --name-only HEAD~1)
    
    if [ -z "$changed_files" ]; then
        print_success "No changes detected"
        return 0
    fi
    
    echo "Changed files:"
    echo "$changed_files" | sed 's/^/  - /'
    
    # Analyze by file type and suggest documentation updates
    suggest_doc_updates "$changed_files"
}

# Smart mapping: what changes affect which docs
suggest_doc_updates() {
    local changed_files="$1"
    local suggestions=()
    
    print_section "Documentation Update Suggestions"
    
    # Service code changes
    if echo "$changed_files" | grep -q "services/.*/.*\.go"; then
        suggestions+=("🔧 Service code changed → Check services/catalog/README.md#api-reference")
        suggestions+=("🔧 Service code changed → Verify services/catalog/README.md#testing-examples")
        if echo "$changed_files" | grep -q "tracing\|metrics\|logger"; then
            suggestions+=("📊 Observability code changed → Update services/catalog/README.md#observability-features-deep-dive")
        fi
    fi
    
    # Kubernetes manifests
    if echo "$changed_files" | grep -q "k8s/.*\.yaml"; then
        suggestions+=("☸️  K8s manifests changed → Check README.md#architecture diagrams")
        if echo "$changed_files" | grep -q "observability"; then
            suggestions+=("📊 Observability stack changed → Update README.md#detailed-setup")
        fi
        if echo "$changed_files" | grep -q "ingress"; then
            suggestions+=("🌐 Ingress changed → Verify README.md#host-configuration")
        fi
    fi
    
    # Helm values
    if echo "$changed_files" | grep -q "values\.yaml"; then
        suggestions+=("⚙️  Helm values changed → Check README.md#detailed-setup")
        suggestions+=("⚙️  Configuration changed → Verify .cursorrules safety rules")
    fi
    
    # Tiltfile changes
    if echo "$changed_files" | grep -q "Tiltfile"; then
        suggestions+=("🔄 Tiltfile changed → Update README.md#launch-the-lab")
    fi
    
    # Scripts
    if echo "$changed_files" | grep -q "scripts/"; then
        suggestions+=("📜 Scripts changed → Check scripts/README.md")
        if echo "$changed_files" | grep -q "simulate-traffic"; then
            suggestions+=("🚦 Traffic simulation changed → Update README.md#traffic-simulation")
        fi
    fi
    
    # Main documentation files
    if echo "$changed_files" | grep -q "README\.md"; then
        suggestions+=("📚 Main README changed → Verify .cursorrules pointers are accurate")
    fi
    
    # New dependencies
    if echo "$changed_files" | grep -q "go\.mod\|go\.sum"; then
        suggestions+=("📦 Dependencies changed → Check .cursorrules safety rules about versions")
    fi
    
    # Print suggestions
    if [ ${#suggestions[@]} -eq 0 ]; then
        print_success "No documentation updates needed for these changes"
    else
        for suggestion in "${suggestions[@]}"; do
            print_action "$suggestion"
        done
    fi
}

# Quick validation of key integration points
validate_setup() {
    print_section "Quick Validation"
    
    # Check if key documentation files exist and have expected sections
    local validation_errors=0
    
    # Check .cursorrules has required sections
    if [ -f ".cursorrules" ]; then
        if grep -q "Documentation Roadmap" .cursorrules && grep -q "Development Decision Tree" .cursorrules; then
            print_success ".cursorrules structure looks good"
        else
            print_warning ".cursorrules missing expected sections"
            validation_errors=$((validation_errors + 1))
        fi
    else
        print_warning ".cursorrules file not found"
        validation_errors=$((validation_errors + 1))
    fi
    
    # Check main README has key sections
    if [ -f "README.md" ]; then
        if grep -q "Architecture" README.md && grep -q "Quick Start" README.md; then
            print_success "README.md structure looks good"
        else
            print_warning "README.md missing expected sections"
            validation_errors=$((validation_errors + 1))
        fi
    else
        print_warning "README.md not found"
        validation_errors=$((validation_errors + 1))
    fi
    
    # Check catalog service README
    if [ -f "services/catalog/README.md" ]; then
        if grep -q "Testing & Examples" services/catalog/README.md && grep -q "Observability" services/catalog/README.md; then
            print_success "Catalog service README structure looks good"
        else
            print_warning "services/catalog/README.md missing expected sections"
            validation_errors=$((validation_errors + 1))
        fi
    else
        print_warning "services/catalog/README.md not found"
        validation_errors=$((validation_errors + 1))
    fi
    
    # Quick functional test if possible
    if command -v curl >/dev/null 2>&1; then
        if curl -s --connect-timeout 2 http://catalog.kubelab.lan:8081/health >/dev/null 2>&1; then
            print_success "Catalog service API responding"
        else
            print_action "Catalog service not responding (might be expected if not running)"
        fi
    fi
    
    return $validation_errors
}

# Comprehensive audit of current state (not based on changes)
audit_current_state() {
    print_section "Comprehensive Documentation Audit"
    
    local audit_issues=()
    local missing_docs=()
    local outdated_docs=()
    
    # 1. Scan for services that should have documentation
    print_action "Scanning for services..."
    if [ -d "services" ]; then
        for service_dir in services/*/; do
            if [ -d "$service_dir" ]; then
                local service_name=$(basename "$service_dir")
                if [ ! -f "$service_dir/README.md" ]; then
                    missing_docs+=("📄 Missing: $service_dir/README.md")
                else
                    # Check if service README has required sections
                    if ! grep -q "API" "$service_dir/README.md"; then
                        outdated_docs+=("📝 $service_dir/README.md missing API documentation")
                    fi
                    if [ -d "$service_dir/internal" ] && ! grep -q -i "observability\|tracing\|metrics" "$service_dir/README.md"; then
                        outdated_docs+=("📊 $service_dir/README.md missing observability documentation")
                    fi
                fi
            fi
        done
    fi
    
    # 2. Check if documented services still exist
    print_action "Checking documented services exist..."
    if [ -f "README.md" ]; then
        # Look for service references in main README
        while IFS= read -r line; do
            if [[ "$line" =~ services/([^/]+)/ ]]; then
                local referenced_service="${BASH_REMATCH[1]}"
                if [ ! -d "services/$referenced_service" ]; then
                    outdated_docs+=("🗑️  README.md references non-existent service: $referenced_service")
                fi
            fi
        done < README.md
    fi
    
    # 3. Check Kubernetes manifests vs documentation
    print_action "Checking Kubernetes documentation alignment..."
    local k8s_apps=$(find k8s/apps -maxdepth 1 -type d ! -name apps 2>/dev/null | wc -l)
    local documented_services=0
    if [ -f "README.md" ] && grep -q "Architecture" README.md; then
        # Count service references in architecture section
        documented_services=$(grep -A 20 -i "architecture" README.md | grep -c "Service" || echo "0")
    fi
    
    if [ "$k8s_apps" -gt 0 ] && [ "$documented_services" -eq 0 ]; then
        outdated_docs+=("🏗️  K8s apps exist but no services documented in README.md architecture")
    fi
    
    # 4. Check observability stack documentation
    print_action "Checking observability documentation..."
    local observability_components=()
    if [ -d "k8s/observability" ]; then
        # Check what's actually deployed
        for file in k8s/observability/*values.yaml; do
            if [ -f "$file" ]; then
                local component=$(basename "$file" | sed 's/-values.yaml//')
                observability_components+=("$component")
            fi
        done
        
        # Check if README documents these components
        for component in "${observability_components[@]}"; do
            if [ -f "README.md" ] && ! grep -q -i "$component" README.md; then
                outdated_docs+=("📊 README.md missing documentation for $component observability component")
            fi
        done
    fi
    
    # 5. Check scripts documentation
    print_action "Checking scripts documentation..."
    if [ -d "scripts" ]; then
        for script in scripts/*.sh; do
            if [ -f "$script" ] && [ ! -f "scripts/README.md" ]; then
                missing_docs+=("📜 Missing: scripts/README.md for script documentation")
                break
            fi
        done
    fi
    
    # 6. Check for undocumented environment variables/configuration
    print_action "Checking configuration documentation..."
    if find . -name "*.yaml" -path "*/k8s/*" -exec grep -l "env:" {} \; | head -1 >/dev/null 2>&1; then
        if [ -f "README.md" ] && ! grep -q -i "environment\|configuration\|host.*configuration" README.md; then
            outdated_docs+=("⚙️  K8s manifests have environment configs but README.md lacks configuration section")
        fi
    fi
    
    # 7. Check API endpoints vs documentation
    print_action "Checking API documentation alignment..."
    for service_dir in services/*/; do
        if [ -d "$service_dir" ]; then
            local service_name=$(basename "$service_dir")
            # Look for HTTP handlers
            if find "$service_dir" -name "*.go" -exec grep -l "router\|http\|Handler" {} \; | head -1 >/dev/null 2>&1; then
                if [ ! -f "$service_dir/README.md" ] || ! grep -q -i "api\|endpoint\|curl" "$service_dir/README.md"; then
                    outdated_docs+=("🔌 $service_name has HTTP handlers but lacks API documentation")
                fi
            fi
        fi
    done
    
    # Print results
    print_section "Audit Results"
    
    if [ ${#missing_docs[@]} -eq 0 ] && [ ${#outdated_docs[@]} -eq 0 ]; then
        print_success "Documentation appears comprehensive and up-to-date!"
        return 0
    fi
    
    if [ ${#missing_docs[@]} -gt 0 ]; then
        echo -e "\n${RED}Missing Documentation:${NC}"
        for item in "${missing_docs[@]}"; do
            echo "  $item"
        done
    fi
    
    if [ ${#outdated_docs[@]} -gt 0 ]; then
        echo -e "\n${YELLOW}Potentially Outdated/Incomplete Documentation:${NC}"
        for item in "${outdated_docs[@]}"; do
            echo "  $item"
        done
    fi
    
    print_action "Run with specific suggestions: ./scripts/sync-docs.sh audit --suggest-fixes"
    
    return 1
}

# Provide specific fix suggestions for audit issues
suggest_audit_fixes() {
    print_section "Suggested Fixes"
    
    # Re-run audit logic but with fix suggestions
    print_action "Based on audit findings, consider these actions:"
    
    # Check for missing service READMEs
    if [ -d "services" ]; then
        for service_dir in services/*/; do
            if [ -d "$service_dir" ] && [ ! -f "$service_dir/README.md" ]; then
                local service_name=$(basename "$service_dir")
                echo "  📄 Create $service_dir/README.md using services/catalog/README.md as template"
            fi
        done
    fi
    
    # Check main README completeness
    if [ -f "README.md" ]; then
        if ! grep -q -i "architecture.*diagram" README.md; then
            echo "  🏗️  Add architecture diagram to README.md (use Mermaid syntax)"
        fi
        if ! grep -q -i "host.*configuration\|etc/hosts" README.md; then
            echo "  🌐 Add host configuration section to README.md"
        fi
    fi
    
    # Check observability documentation
    if [ -d "k8s/observability" ] && [ -f "README.md" ]; then
        local components=$(find k8s/observability -name "*values.yaml" | wc -l)
        local documented=$(grep -c -i "prometheus\|grafana\|loki\|tempo" README.md || echo "0")
        if [ "$components" -gt "$documented" ]; then
            echo "  📊 Update README.md observability section to document all components"
        fi
    fi
    
    print_action "For detailed guidance, see scripts/WORKFLOW.md"
}

# Update .cursorrules if new patterns detected
update_cursorrules() {
    local changed_files="$1"
    print_section "Cursor Rules Update Check"
    
    local needs_update=false
    
    # Check if new services were added
    if echo "$changed_files" | grep -q "services/[^/]*/README\.md" && ! echo "$changed_files" | grep -q "services/catalog/"; then
        print_action "New service detected → Consider updating .cursorrules documentation roadmap"
        needs_update=true
    fi
    
    # Check if new safety patterns emerged
    if echo "$changed_files" | grep -q "\.sh$" && echo "$changed_files" | grep -q "kubectl\|tilt"; then
        print_action "New scripts with kubectl/tilt → Verify .cursorrules safety rules"
        needs_update=true
    fi
    
    # Check if new observability patterns
    if echo "$changed_files" | grep -q "tracing\|metrics\|logging" && echo "$changed_files" | grep -q "internal/"; then
        print_action "New observability patterns → Consider updating .cursorrules decision tree"
        needs_update=true
    fi
    
    if [ "$needs_update" = false ]; then
        print_success "No .cursorrules updates needed"
    fi
}

# Print usage
usage() {
    echo "Documentation Sync Script"
    echo ""
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  analyze [since]     Analyze changes and suggest doc updates"
    echo "                      since: git ref to compare against (default: HEAD~5)"
    echo "  analyze --since-tag Compare against last git tag"
    echo "  audit               Comprehensive audit of current documentation state"
    echo "  audit --suggest-fixes  Audit with specific fix suggestions"
    echo "  validate           Quick validation of documentation structure"
    echo "  full               Run both analyze and validate"
    echo ""
    echo "Examples:"
    echo "  $0 analyze                    # Check last 5 commits"
    echo "  $0 analyze --since-tag        # Check since last tag"
    echo "  $0 analyze HEAD~10            # Check last 10 commits"
    echo "  $0 audit                      # Comprehensive state audit"
    echo "  $0 audit --suggest-fixes      # Audit with fix suggestions"
    echo "  $0 validate                   # Just validate structure"
    echo "  $0 full                       # Full analysis + validation"
}

# Main function
main() {
    local command="${1:-analyze}"
    
    case "$command" in
        "analyze")
            local comparison_point=$(get_comparison_point "$2")
            analyze_changes "$comparison_point"
            update_cursorrules "$(git diff --name-only "$comparison_point" 2>/dev/null || git diff --name-only HEAD~1)"
            ;;
        "audit")
            if [ "$2" = "--suggest-fixes" ]; then
                audit_current_state
                suggest_audit_fixes
            else
                audit_current_state
            fi
            ;;
        "validate")
            validate_setup
            ;;
        "full")
            local comparison_point=$(get_comparison_point "$2")
            local changed_files=$(git diff --name-only "$comparison_point" 2>/dev/null || git diff --name-only HEAD~1)
            analyze_changes "$comparison_point"
            update_cursorrules "$changed_files"
            validate_setup
            ;;
        "--help"|"-h"|"help")
            usage
            ;;
        *)
            echo "Unknown command: $command"
            usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 