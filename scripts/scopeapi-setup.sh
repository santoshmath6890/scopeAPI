#!/bin/bash

# ScopeAPI Unified Setup Script
# Purpose: Orchestrate complete infrastructure and database setup
# Usage: ./scopeapi-setup.sh [--full|--infrastructure|--database|--validate]
# Features: Infrastructure startup, database setup, validation, test data

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Get the project root directory (parent of scripts directory)
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

# Function to show help
show_help() {
    echo "ScopeAPI Setup Script"
    echo ""
    echo "Usage: ./scopeapi-setup.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --help, -h           Show this help message"
    echo "  --verbose, -v        Enable verbose output"
    echo "  --infrastructure     Start infrastructure services only"
    echo "  --database           Setup database only"
    echo "  --test-data          Create sample test data"
    echo "  --validate           Run validation tests"
    echo "  --full               Complete setup (infrastructure + database + validation)"
    echo "  --cleanup            Stop and remove all infrastructure services"
    echo "  --cleanup-full       Stop services, remove containers, volumes, and networks"
    echo ""
        echo "Examples:"
    echo "  ./scopeapi-setup.sh --full            # Complete setup with validation"
    echo "  ./scopeapi-setup.sh --infrastructure  # Start infrastructure only"
    echo "  ./scopeapi-setup.sh --database        # Setup database only"
    echo "  ./scopeapi-setup.sh --test-data       # Setup + create test data"
    echo "  ./scopeapi-setup.sh --validate        # Setup + run validation tests"
    echo "  ./scopeapi-setup.sh --cleanup         # Stop and remove services"
    echo "  ./scopeapi-setup.sh --cleanup-full    # Remove everything (containers, volumes, networks)"
    echo ""
    echo "This script will:"
    echo "  1. Start infrastructure services (ZooKeeper, Kafka, PostgreSQL, Redis, Elasticsearch, Kibana)"
    echo "  2. Setup PostgreSQL database with migrations"
    echo "  3. Create sample test data (optional)"
    echo "  4. Validate the complete setup (optional)"
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    print_success "Docker is running"
    
    # Check if docker-compose is available
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose is not installed. Please install it first."
        exit 1
    fi
    print_success "docker-compose is available"
    
    # Check if .env file exists
    if [ ! -f ".env" ]; then
        print_warning ".env file not found"
        print_status "Creating .env file from template..."
        if [ -f "$SCRIPT_DIR/env.example" ]; then
            cp "$SCRIPT_DIR/env.example" .env
            print_success "Created .env file from env.example"
            print_warning "Please edit .env file and set your secure passwords"
            print_warning "Then run this script again"
            exit 1
        else
            print_error "env.example file not found"
            print_status "Please create a .env file manually with required environment variables"
            exit 1
        fi
    fi
    print_success ".env file found"
    
    # Load environment variables
    print_status "Loading environment variables..."
    set -a  # automatically export all variables
    source .env
    set +a  # stop automatically exporting
    print_success "Environment variables loaded"
}

# Function to start infrastructure
start_infrastructure() {
    print_header "Starting Infrastructure Services"
    
    print_status "Starting infrastructure services..."
    
    # Start infrastructure using docker-compose
    if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/.env" up -d zookeeper kafka postgres redis elasticsearch kibana; then
        print_success "Infrastructure services started successfully"
        
        # Wait for services to be ready
        print_status "Waiting for services to be ready..."
        sleep 15
        
        # Check service status
        if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/.env" ps | grep -q "Up"; then
            print_success "All infrastructure services are running"
        else
            print_warning "Some services may not be fully ready yet"
        fi
    else
        print_error "Failed to start infrastructure services"
        exit 1
    fi
}

# Function to setup database
setup_database() {
    print_header "Setting Up Database"
    
    print_status "Setting up PostgreSQL database..."
    
    # Run the database setup script with environment variables (basic setup only)
    if DB_HOST="$DB_HOST" DB_PORT="$DB_PORT" DB_USER="$DB_USER" DB_PASSWORD="$DB_PASSWORD" DB_NAME="$DB_NAME" DB_SSL_MODE="$DB_SSL_MODE" "$SCRIPT_DIR/setup-database.sh" --basic; then
        print_success "Database setup completed successfully"
    else
        print_error "Database setup failed"
        exit 1
    fi
}

# Function to create test data
create_test_data() {
    print_header "Creating Test Data"
    
    print_status "Creating sample test data..."
    
    # Run database setup with test data flag (basic mode)
    if "$SCRIPT_DIR/setup-database.sh" --basic --test-data; then
        print_success "Test data created successfully"
    else
        print_warning "Test data creation failed (may already exist)"
    fi
}

# Function to validate setup
validate_setup() {
    print_header "Validating Setup"
    
    print_status "Running validation tests..."
    
    # Run database setup with validation flag (basic mode)
    if "$SCRIPT_DIR/setup-database.sh" --basic --validate; then
        print_success "Validation completed successfully"
    else
        print_error "Validation failed"
        exit 1
    fi
    
    # Check if all services are running
    print_status "Checking service status..."
    if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/.env" ps | grep -q "Up"; then
        print_success "All services are running"
    else
        print_warning "Some services may not be running"
    fi
}

# Function to show final status
show_final_status() {
    print_header "Setup Complete!"
    
    print_success "ScopeAPI setup has been completed successfully!"
    echo ""
    print_status "Next steps:"
    echo "  1. Start microservices: ./scripts/scopeapi-local.sh start all"
    echo "  2. Or start specific service: ./scripts/scopeapi-local.sh start api-discovery"
    echo "  3. For debugging: ./scripts/scopeapi-debug.sh start api-discovery"
    echo ""
    print_status "Service URLs:"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Kafka: localhost:9092"
    echo "  - Redis: localhost:6379"
    echo "  - Elasticsearch: localhost:9200"
    echo "  - Kibana: localhost:5601"
    echo ""
    print_status "Useful commands:"
    echo "  - View logs: ./scripts/scopeapi-local.sh logs [service]"
    echo "  - Check status: ./scripts/scopeapi-local.sh status"
    echo "  - Stop services: ./scripts/scopeapi-local.sh stop"
}

# Function to cleanup infrastructure services
cleanup_infrastructure() {
    print_header "Cleaning Up Infrastructure Services"
    
    print_status "Stopping and removing infrastructure services..."
    
    # Stop and remove services
    if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/.env" down; then
        print_success "Infrastructure services stopped and removed successfully"
    else
        print_error "Failed to stop infrastructure services"
        return 1
    fi
}

# Function to cleanup everything (containers, volumes, networks)
cleanup_full() {
    print_header "Full Cleanup - Removing Everything"
    
    print_warning "This will remove ALL containers, volumes, and networks!"
    print_warning "This action cannot be undone."
    echo ""
    print_status "Are you sure you want to continue? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Performing full cleanup..."
        
        # Stop and remove everything
        if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/.env" down -v --remove-orphans; then
            print_success "All containers, volumes, and networks removed successfully"
        else
            print_error "Failed to perform full cleanup"
            return 1
        fi
        
        # Remove any orphaned containers
        print_status "Removing orphaned containers..."
        docker container prune -f >/dev/null 2>&1
        
        # Remove any orphaned networks
        print_status "Removing orphaned networks..."
        docker network prune -f >/dev/null 2>&1
        
        print_success "Full cleanup completed successfully"
    else
        print_info "Cleanup cancelled"
    fi
}

# Main execution
main() {
    local infrastructure_flag=false
    local database_flag=false
    local test_data_flag=false
    local validate_flag=false
    local full_flag=false
    local cleanup_flag=false
    local cleanup_full_flag=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --infrastructure)
                infrastructure_flag=true
                shift
                ;;
            --database)
                database_flag=true
                shift
                ;;
            --test-data)
                test_data_flag=true
                shift
                ;;
            --validate)
                validate_flag=true
                shift
                ;;
            --full)
                full_flag=true
                shift
                ;;
            --cleanup)
                cleanup_flag=true
                shift
                ;;
            --cleanup-full)
                cleanup_full_flag=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            --verbose|-v)
                set -x
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # If no specific flags, default to full setup
    if [ "$infrastructure_flag" = false ] && [ "$database_flag" = false ] && [ "$test_data_flag" = false ] && [ "$validate_flag" = false ]; then
        full_flag=true
    fi
    
    print_header "ScopeAPI Setup"
    print_status "Starting setup process..."
    
    # Handle cleanup operations first
    if [ "$cleanup_flag" = true ]; then
        cleanup_infrastructure
        exit 0
    fi
    
    if [ "$cleanup_full_flag" = true ]; then
        cleanup_full
        exit 0
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Start infrastructure if requested
    if [ "$infrastructure_flag" = true ] || [ "$full_flag" = true ]; then
        start_infrastructure
    fi
    
    # Setup database if requested
    if [ "$database_flag" = true ] || [ "$full_flag" = true ]; then
        setup_database
    fi
    
    # Create test data if requested
    if [ "$test_data_flag" = true ] || [ "$full_flag" = true ]; then
        create_test_data
    fi
    
    # Validate setup if requested
    if [ "$validate_flag" = true ] || [ "$full_flag" = true ]; then
        validate_setup
    fi
    
    # Show final status
    show_final_status
}

# Run main function with all arguments
main "$@"
