#!/bin/bash

# Kubernetes Lab - Phase Management Script
# This script helps manage the phased learning approach with git tags and releases

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Show usage
show_usage() {
    echo "Kubernetes Lab - Phase Management Script"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  list-phases              List all available phase tags"
    echo "  checkout <phase>         Checkout a specific phase"
    echo "  create-tag <tag> <msg>   Create a new phase tag"
    echo "  prepare-release <phase>  Prepare files for a new phase release"
    echo "  current-phase           Show current phase (if on a tag)"
    echo "  phase-diff <phase1> <phase2>  Show differences between phases"
    echo "  create-retrospective-tags     Create tags for phases that should have existed"
    echo ""
    echo "Examples:"
    echo "  $0 list-phases"
    echo "  $0 checkout v1.0.0-monitoring-foundation"
    echo "  $0 create-tag v1.1.0-loki-integration 'Add Loki for log aggregation'"
    echo "  $0 prepare-release v2.0.0-basic-service"
    echo "  $0 create-retrospective-tags"
}

# List all phase tags
list_phases() {
    print_info "Available phases:"
    echo ""
    
    # Get all tags that match our phase naming convention
    local tags=$(git tag | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+' | sort -V)
    
    if [ -z "$tags" ]; then
        print_warning "No phase tags found. Use 'create-retrospective-tags' to create initial tags."
        return 0
    fi
    
    local current_phase=""
    local current_tag=$(git describe --exact-match HEAD 2>/dev/null || echo "")
    
    echo "Phase 1: Observability Stack Foundation"
    echo "----------------------------------------"
    for tag in $tags; do
        if [[ $tag == v1.* ]]; then
            local marker=" "
            if [ "$tag" = "$current_tag" ]; then
                marker="*"
                current_phase="$tag"
            fi
            echo "  $marker $tag"
        fi
    done
    
    echo ""
    echo "Phase 2: First Service - Application Observability"  
    echo "------------------------------------------------"
    for tag in $tags; do
        if [[ $tag == v2.* ]]; then
            local marker=" "
            if [ "$tag" = "$current_tag" ]; then
                marker="*"
                current_phase="$tag"
            fi
            echo "  $marker $tag"
        fi
    done
    
    echo ""
    echo "Phase 3: Microservices (Future)"
    echo "------------------------------"
    for tag in $tags; do
        if [[ $tag == v3.* ]]; then
            local marker=" "
            if [ "$tag" = "$current_tag" ]; then
                marker="*"
                current_phase="$tag"
            fi
            echo "  $marker $tag"
        fi
    done
    
    if [ -n "$current_phase" ]; then
        echo ""
        print_info "Currently on phase: $current_phase"
    else
        echo ""
        print_info "Currently on: $(git branch --show-current)"
    fi
}

# Checkout a specific phase
checkout_phase() {
    local phase="$1"
    
    if [ -z "$phase" ]; then
        print_error "Please specify a phase to checkout"
        echo "Usage: $0 checkout <phase>"
        echo "Example: $0 checkout v1.0.0-monitoring-foundation"
        return 1
    fi
    
    # Check if tag exists
    if ! git tag -l | grep -q "^$phase$"; then
        print_error "Phase '$phase' not found"
        echo ""
        print_info "Available phases:"
        list_phases
        return 1
    fi
    
    print_info "Checking out phase: $phase"
    git checkout "$phase"
    
    print_success "Successfully checked out phase: $phase"
    echo ""
    print_info "To start the lab environment:"
    echo "  ./start-lab.sh"
    echo ""
    print_info "To read phase documentation:"
    echo "  cat phases/phase-*/README.md"
}

# Create a new phase tag
create_tag() {
    local tag="$1"
    local message="$2"
    
    if [ -z "$tag" ] || [ -z "$message" ]; then
        print_error "Please provide both tag name and message"
        echo "Usage: $0 create-tag <tag> <message>"
        echo "Example: $0 create-tag v1.1.0-loki-integration 'Add Loki for log aggregation'"
        return 1
    fi
    
    # Validate tag format
    if ! echo "$tag" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then
        print_error "Tag must follow semantic versioning: v<major>.<minor>.<patch>-<description>"
        return 1
    fi
    
    # Check if tag already exists
    if git tag -l | grep -q "^$tag$"; then
        print_error "Tag '$tag' already exists"
        return 1
    fi
    
    # Ensure we're on main branch
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        print_warning "Currently on branch '$current_branch', not 'main'"
        read -p "Do you want to continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Aborted"
            return 1
        fi
    fi
    
    print_info "Creating tag: $tag"
    git tag -a "$tag" -m "$message"
    
    print_success "Tag created successfully: $tag"
    echo ""
    print_info "To push the tag to remote:"
    echo "  git push origin $tag"
    echo ""
    print_info "To create a GitHub release:"
    echo "  1. Go to https://github.com/mansoormajeed/kubernetes-lab/releases"
    echo "  2. Click 'Create a new release'"
    echo "  3. Select tag '$tag'"
    echo "  4. Add release notes with learning objectives"
}

# Show current phase
current_phase() {
    local current_tag=$(git describe --exact-match HEAD 2>/dev/null || echo "")
    
    if [ -n "$current_tag" ]; then
        print_info "Currently on phase: $current_tag"
        
        # Try to find phase documentation
        local phase_num=""
        if [[ $current_tag == v1.* ]]; then
            phase_num="1"
        elif [[ $current_tag == v2.* ]]; then
            phase_num="2"
        elif [[ $current_tag == v3.* ]]; then
            phase_num="3"
        elif [[ $current_tag == v4.* ]]; then
            phase_num="4"
        fi
        
        if [ -n "$phase_num" ] && [ -f "phases/phase-$phase_num/README.md" ]; then
            echo ""
            print_info "Phase documentation available at: phases/phase-$phase_num/README.md"
        fi
    else
        local branch=$(git branch --show-current)
        print_info "Currently on branch: $branch (not on a phase tag)"
    fi
}

# Create retrospective tags for phases that should have existed
create_retrospective_tags() {
    print_info "Creating retrospective tags for missing phases..."
    echo ""
    
    print_warning "This will create tags for phases that should have existed historically."
    print_warning "This is useful for setting up the initial phase structure."
    echo ""
    
    # Check if we're on main and have a clean working tree
    local current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        print_error "Must be on main branch to create retrospective tags"
        return 1
    fi
    
    if ! git diff --quiet; then
        print_error "Working tree is not clean. Please commit or stash changes."
        return 1
    fi
    
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Aborted"
        return 1
    fi
    
    # Strategy: Create tags at specific commits that represent each phase
    # Since we already have the go service, we need to create earlier tags
    # at commits before the service was added
    
    print_info "Analyzing git history to find appropriate commits for retrospective tags..."
    
    # Find commits that would represent each phase
    local commits=$(git log --oneline --reverse)
    
    # For now, let's create tags at current HEAD with different names
    # This is a starting point - in real usage, you'd create these at appropriate historical commits
    
    print_info "Creating Phase 1 foundation tags..."
    
    # Create a temporary branch to work with historical state
    git checkout -b temp-phase-creation
    
    # Reset to a point before the Go service was added (you'd need to identify this commit)
    # For demonstration, we'll use HEAD~10 or find the appropriate commit
    
    # Find the commit before Go service was added
    local pre_service_commit=$(git log --oneline --grep="add.*service\|Add.*service" | tail -n 1 | cut -d' ' -f1)
    
    if [ -z "$pre_service_commit" ]; then
        # If we can't find a specific commit, use a reasonable point in history
        print_warning "Could not find commit before service was added. Using HEAD~5"
        pre_service_commit="HEAD~5"
    fi
    
    # Create tags at the appropriate historical points
    print_info "Creating v1.0.0-monitoring-foundation tag..."
    git tag -a v1.0.0-monitoring-foundation "$pre_service_commit" -m "Phase 1.0.0: Basic monitoring stack with Prometheus and Grafana"
    
    print_info "Creating v1.1.0-loki-integration tag..."
    git tag -a v1.1.0-loki-integration "$pre_service_commit" -m "Phase 1.1.0: Add Loki for log aggregation"
    
    print_info "Creating v1.2.0-alloy-collection tag..."
    git tag -a v1.2.0-alloy-collection "$pre_service_commit" -m "Phase 1.2.0: Add Alloy for advanced log collection"
    
    # Now tag the current state as v2.0.0
    git checkout main
    git branch -D temp-phase-creation
    
    print_info "Creating v2.0.0-basic-service tag..."
    git tag -a v2.0.0-basic-service HEAD -m "Phase 2.0.0: Deploy catalog service with basic observability"
    
    print_success "Retrospective tags created successfully!"
    echo ""
    print_info "Created tags:"
    echo "  v1.0.0-monitoring-foundation - Basic monitoring stack"
    echo "  v1.1.0-loki-integration - Add Loki for logs"
    echo "  v1.2.0-alloy-collection - Add Alloy for log collection"
    echo "  v2.0.0-basic-service - Deploy catalog service"
    echo ""
    print_info "To push all tags to remote:"
    echo "  git push origin --tags"
    echo ""
    print_info "Next steps:"
    echo "  1. Test each phase: $0 checkout <phase-tag>"
    echo "  2. Create GitHub releases for each tag"
    echo "  3. Continue with v2.1.0 for Prometheus metrics"
}

# Show differences between phases
phase_diff() {
    local phase1="$1"
    local phase2="$2"
    
    if [ -z "$phase1" ] || [ -z "$phase2" ]; then
        print_error "Please specify two phases to compare"
        echo "Usage: $0 phase-diff <phase1> <phase2>"
        echo "Example: $0 phase-diff v1.0.0-monitoring-foundation v2.0.0-basic-service"
        return 1
    fi
    
    # Check if tags exist
    if ! git tag -l | grep -q "^$phase1$"; then
        print_error "Phase '$phase1' not found"
        return 1
    fi
    
    if ! git tag -l | grep -q "^$phase2$"; then
        print_error "Phase '$phase2' not found"
        return 1
    fi
    
    print_info "Differences between $phase1 and $phase2:"
    echo ""
    
    # Show file differences
    git diff --name-status "$phase1" "$phase2"
    
    echo ""
    print_info "To see detailed diff:"
    echo "  git diff $phase1 $phase2"
}

# Main script logic
case "$1" in
    "list-phases")
        list_phases
        ;;
    "checkout")
        checkout_phase "$2"
        ;;
    "create-tag")
        create_tag "$2" "$3"
        ;;
    "current-phase")
        current_phase
        ;;
    "phase-diff")
        phase_diff "$2" "$3"
        ;;
    "create-retrospective-tags")
        create_retrospective_tags
        ;;
    *)
        show_usage
        ;;
esac 