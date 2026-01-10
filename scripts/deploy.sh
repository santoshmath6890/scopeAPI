#!/bin/bash

# ScopeAPI Deployment Script
# Purpose: Unified deployment for Docker (local) and Kubernetes (staging/production)
# Usage: ./deploy.sh [OPTIONS]
# Features: Environment-specific deployment, secrets management, validation

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to project root directory
cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
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

print_header() {
    echo -e "${BLUE}==========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}==========================================${NC}"
}

# Default values
ENVIRONMENT="dev"
PLATFORM="docker"
VERBOSE=false

# Function to show usage
show_usage() {
    echo "ScopeAPI Deployment Script"
    echo ""
    echo "Usage: ./deploy.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -e, --environment ENV    Deployment environment (dev|staging|prod) [default: dev]"
    echo "  -p, --platform PLATFORM  Deployment platform (docker|k8s) [default: docker]"
    echo "  -v, --verbose            Enable verbose output"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Environment-specific behavior:"
    echo "  dev: Uses .env.local file (LOCAL DEVELOPMENT ONLY)"
    echo "  staging: Uses Kubernetes Secrets"
    echo "  prod: Uses Kubernetes Secrets"
    echo ""
    echo "Examples:"
    echo "  ./deploy.sh                           # Deploy to Docker (dev)"
    echo "  ./deploy.sh -e dev -p docker          # Deploy to Docker (dev)"
    echo "  ./deploy.sh -e staging -p k8s         # Deploy to Kubernetes (staging)"
    echo "  ./deploy.sh -e prod -p k8s            # Deploy to Kubernetes (production)"
    echo ""
    echo "Note: Docker deployment is only allowed for LOCAL DEVELOPMENT (dev environment)"
}

# Function to parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--environment)
                ENVIRONMENT="$2"
                shift 2
                ;;
            -p|--platform)
                PLATFORM="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Function to validate environment and platform combination
validate_deployment() {
    print_status "Validating deployment configuration..."
    
    # Validate environment
    case "$ENVIRONMENT" in
        dev|staging|prod)
            print_success "Environment: $ENVIRONMENT"
            ;;
        *)
            print_error "Invalid environment: $ENVIRONMENT"
            print_error "Valid environments: dev, staging, prod"
            exit 1
            ;;
    esac
    
    # Validate platform
    case "$PLATFORM" in
        docker|k8s)
            print_success "Platform: $PLATFORM"
            ;;
        *)
            print_error "Invalid platform: $PLATFORM"
            print_error "Valid platforms: docker, k8s"
            exit 1
            ;;
    esac
    
    # Validate platform/environment combination
    if [[ "$ENVIRONMENT" != "dev" && "$PLATFORM" == "docker" ]]; then
        print_error "❌ Docker deployment is ONLY allowed for LOCAL DEVELOPMENT (dev environment)!"
        print_error "For staging and production, you MUST use Kubernetes:"
        print_error "  Staging: ./deploy.sh -e staging -p k8s"
        print_error "  Production: ./deploy.sh -e prod -p k8s"
        exit 1
    fi

    if [[ "$ENVIRONMENT" == "prod" && "$PLATFORM" == "docker" ]]; then
        print_error "❌ Production environment with Docker is NOT ALLOWED!"
        print_error "Production MUST use Kubernetes for security and scalability."
        exit 1
    fi
    
    print_success "Deployment configuration validated!"
}

# Function to deploy to Docker (LOCAL DEVELOPMENT ONLY)
deploy_docker() {
    # Docker deployment is only for local development
    if [[ "$ENVIRONMENT" != "dev" ]]; then
        print_error "Docker deployment is only allowed for LOCAL DEVELOPMENT (dev environment)!"
        print_error "For staging/production, use Kubernetes: ./deploy.sh -e $ENVIRONMENT -p k8s"
        exit 1
    fi
    
    # Check for .env.local file
    if [[ ! -f ".env.local" ]]; then
        print_error "No .env.local file found for local development."
        print_info "Run: cp env.example .env.local && nano .env.local"
        exit 1
    fi
    
    local env_file=".env.local"
    print_info "⚠️  DEPLOYING TO DOCKER FOR LOCAL DEVELOPMENT ONLY!"
    print_info "⚠️  .env.local file will be used (your local machine only)"
    
    print_info "Deploying to Docker using $env_file..."
    
    # Start infrastructure first
    print_status "Starting infrastructure services..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$env_file" up -d zookeeper kafka postgres redis elasticsearch kibana
    
    # Wait for infrastructure to be ready
    print_status "Waiting for infrastructure to be ready..."
    sleep 15
    
    # Start all microservices
    print_status "Starting microservices..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$env_file" up -d api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
    
    # Check status
    print_status "Checking deployment status..."
    if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$env_file" ps | grep -q "Up"; then
        print_success "✅ Docker deployment completed successfully!"
        print_info "Services are now running on your local machine"
        print_info "Use './scopeapi.sh status' to check service status"
    else
        print_error "❌ Docker deployment failed"
        exit 1
    fi
}

# Function to deploy to Kubernetes
deploy_kubernetes() {
    print_header "Deploying to Kubernetes ($ENVIRONMENT environment)"
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed. Please install it first."
        exit 1
    fi
    
    # Check if we can connect to the cluster
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    print_success "Connected to Kubernetes cluster"
    
    # Check if secrets file exists
    if [[ ! -f "k8s/secrets.yaml" ]]; then
        print_warning "No secrets file found. Creating template..."
        if [[ -f "scripts/generate-secrets.sh" ]]; then
            print_status "Generating secrets template..."
            ./scripts/generate-secrets.sh
            print_warning "⚠️  IMPORTANT: Edit k8s/secrets.yaml with your actual secrets before deploying!"
            print_warning "⚠️  NEVER commit the secrets.yaml file with real values!"
            exit 1
        else
            print_error "Secrets generation script not found"
            exit 1
        fi
    fi
    
    # Apply Kubernetes configurations
    print_status "Applying Kubernetes configurations..."
    
    # Create namespace first
    if kubectl apply -f k8s/namespace.yaml; then
        print_success "Namespace created/updated"
    else
        print_error "Failed to create namespace"
        exit 1
    fi
    
    # Apply RBAC
    if kubectl apply -f k8s/rbac/; then
        print_success "RBAC applied"
    else
        print_error "Failed to apply RBAC"
        exit 1
    fi
    
    # Apply ConfigMap
    if kubectl apply -f k8s/configmap.yaml; then
        print_success "ConfigMap applied"
    else
        print_error "Failed to apply ConfigMap"
        exit 1
    fi
    
    # Apply Secrets
    if kubectl apply -f k8s/secrets.yaml; then
        print_success "Secrets applied"
    else
        print_error "Failed to apply Secrets"
        exit 1
    fi
    
    # Apply Deployments
    if kubectl apply -f k8s/deployments/; then
        print_success "Deployments applied"
    else
        print_error "Failed to apply Deployments"
        exit 1
    fi
    
    # Apply Services
    if kubectl apply -f k8s/services/; then
        print_success "Services applied"
    else
        print_error "Failed to apply Services"
        exit 1
    fi
    
    # Apply Ingress (if exists)
    if [[ -f "k8s/ingress/ingress.yaml" ]]; then
        if kubectl apply -f k8s/ingress/; then
            print_success "Ingress applied"
        else
            print_warning "Failed to apply Ingress (may not be configured)"
        fi
    fi
    
    # Wait for deployments to be ready
    print_status "Waiting for deployments to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment -l app=scopeapi -n scopeapi
    
    # Show deployment status
    print_status "Deployment Status:"
    kubectl get pods -n scopeapi
    kubectl get services -n scopeapi
    
    print_success "✅ Kubernetes deployment completed successfully!"
    print_info "Use 'kubectl get pods -n scopeapi' to check pod status"
    print_info "Use 'kubectl logs -n scopeapi <pod-name>' to view logs"
}

# Function to show deployment summary
show_summary() {
    print_header "Deployment Summary"
    
    echo "Environment: $ENVIRONMENT"
    echo "Platform: $PLATFORM"
    echo "Timestamp: $(date)"
    
    case "$PLATFORM" in
        docker)
            echo "Deployment Type: Local Development (Docker)"
            echo "Environment File: .env.local"
            echo "Services: Running on localhost"
            ;;
        k8s)
            echo "Deployment Type: Kubernetes ($ENVIRONMENT)"
            echo "Secrets: Kubernetes Secrets"
            echo "Namespace: scopeapi"
            echo "Services: Cluster-internal"
            ;;
    esac
    
    echo ""
    print_status "Next Steps:"
    case "$PLATFORM" in
        docker)
            echo "  - Check status: ./scopeapi.sh status"
            echo "  - View logs: ./scopeapi.sh logs [service]"
            echo "  - Stop services: ./scopeapi.sh stop"
            ;;
        k8s)
            echo "  - Check pods: kubectl get pods -n scopeapi"
            echo "  - View logs: kubectl logs -n scopeapi <pod-name>"
            echo "  - Access services: kubectl port-forward -n scopeapi service/<service-name> <local-port>:<service-port>"
            ;;
    esac
}

# Main execution
main() {
    print_header "ScopeAPI Deployment"
    
    # Parse arguments
    parse_arguments "$@"
    
    # Validate deployment configuration
    validate_deployment
    
    # Perform deployment
    case "$PLATFORM" in
        docker)
            deploy_docker
            ;;
        k8s)
            deploy_kubernetes
            ;;
        *)
            print_error "Unknown platform: $PLATFORM"
            exit 1
            ;;
    esac
    
    # Show summary
    show_summary
}

# Run main function with all arguments
main "$@"
