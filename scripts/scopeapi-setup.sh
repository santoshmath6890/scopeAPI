#!/bin/bash

# ScopeAPI Unified Setup Script
# Purpose: Orchestrate complete infrastructure and database setup
# Usage: ./scopeapi-setup.sh [--full|--infrastructure|--database|--validate]
# Features: Infrastructure startup, database setup, validation, test data

set -e

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
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --help, -h           Show this help message"
    echo "  --verbose, -v        Enable verbose output"
    echo "  --infrastructure     Start infrastructure services only"
    echo "  --database           Setup database only"
    echo "  --test-data          Create sample test data"
    echo "  --validate           Run validation tests"
    echo "  --full               Complete setup (infrastructure + database + validation)"
    echo ""
    echo "Examples:"
    echo "  $0 --full            # Complete setup with validation"
    echo "  $0 --infrastructure  # Start infrastructure only"
    echo "  $0 --database        # Setup database only"
    echo "  $0 --test-data       # Setup + create test data"
    echo "  $0 --validate        # Setup + run validation tests"
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
        if [ -f "env.example" ]; then
            cp env.example .env
            print_success "Created .env file from env.example"
            print_warning "Please edit .env file and set your secure passwords"
            print_warning "Then run this script again"
            exit 1
        else
            print_error "env.example file not found"
            print_info "Please create a .env file manually with required environment variables"
            exit 1
        fi
    fi
    print_success ".env file found"
}

# Function to start infrastructure
start_infrastructure() {
    print_header "Starting Infrastructure Services"
    
    print_status "Starting infrastructure services..."
    
    # Start infrastructure using docker-compose
    if docker-compose up -d zookeeper kafka postgres redis elasticsearch kibana; then
        print_success "Infrastructure services started successfully"
        
        # Wait for services to be ready
        print_status "Waiting for services to be ready..."
        sleep 15
        
        # Check service status
        if docker-compose ps | grep -q "Up"; then
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
    
    # Run the database setup script
    if ./scripts/setup-database.sh; then
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
    
    # Run database setup with test data flag
    if ./scripts/setup-database.sh --test-data; then
        print_success "Test data created successfully"
    else
        print_warning "Test data creation failed (may already exist)"
    fi
}

# Function to validate setup
validate_setup() {
    print_header "Validating Setup"
    
    print_status "Running validation tests..."
    
    # Run database setup with validation flag
    if ./scripts/setup-database.sh --validate; then
        print_success "Validation completed successfully"
    else
        print_error "Validation failed"
        exit 1
    fi
    
    # Check if all services are running
    print_status "Checking service status..."
    if docker-compose ps | grep -q "Up"; then
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
    echo "  1. Start microservices: ./scripts/dev.sh start all"
    echo "  2. Or start specific service: ./scripts/dev.sh start api-discovery"
    echo "  3. For debugging: ./scripts/debug.sh start api-discovery"
    echo ""
    print_status "Service URLs:"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Kafka: localhost:9092"
    echo "  - Redis: localhost:6379"
    echo "  - Elasticsearch: localhost:9200"
    echo "  - Kibana: localhost:5601"
    echo ""
    print_status "Useful commands:"
    echo "  - View logs: ./scripts/dev.sh logs [service]"
    echo "  - Check status: ./scripts/dev.sh status"
    echo "  - Stop services: ./scripts/dev.sh stop"
}

# Main execution
main() {
    local infrastructure_flag=false
    local database_flag=false
    local test_data_flag=false
    local validate_flag=false
    local full_flag=false
    
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
