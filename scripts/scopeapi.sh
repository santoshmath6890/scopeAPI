#!/bin/bash

# ScopeAPI Main Orchestrator Script
# Purpose: Unified script for all ScopeAPI operations
# Usage: ./scopeapi.sh [COMMAND] [OPTIONS]
# Features: Setup, services, status, and main operations

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

# Function to show help
show_help() {
    echo "ScopeAPI Main Orchestrator"
    echo ""
    echo "Usage: ./scopeapi.sh [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  setup [OPTIONS]     - Complete setup and validation"
    echo "  start [SERVICES...] - Start infrastructure + services"
    echo "  stop                - Stop all services"
    echo "  restart [SERVICES...] - Restart services"
    echo "  status              - Show comprehensive status"
    echo "  comprehensive-status - Show detailed infrastructure + microservices status"
    echo "  logs [SERVICE]      - Show logs for service(s)"
    echo "  build [SERVICES...] - Build specified services"
    echo "  clean               - Clean all containers and volumes"
    echo "  shell [SERVICE]     - Open shell in service container"
    echo "  exec [SERVICE] [CMD] - Execute command in container"
    echo "  debug [SERVICE]     - Start service in debug mode"
    echo "  help                - Show this help message"
    echo ""
    echo "Setup Options:"
    echo "  --full              - Complete setup (infrastructure + database + validation)"
    echo "  --infrastructure    - Start infrastructure services only"
    echo "  --database          - Setup database only"
    echo "  --test-data         - Create sample test data"
    echo "  --validate          - Run validation tests"
    echo "  --cleanup           - Stop and remove services"
    echo "  --cleanup-full      - Remove everything (containers, volumes, networks)"
    echo ""
    echo "Services:"
    echo "  infrastructure      - Start only infrastructure (postgres, kafka, redis, etc.)"
    echo "  api-discovery       - API Discovery service"
    echo "  gateway-integration - Gateway Integration service"
    echo "  data-ingestion     - Data Ingestion service"
    echo "  threat-detection   - Threat Detection service"
    echo "  data-protection    - Data Protection service"
    echo "  attack-blocking    - Attack Blocking service"
    echo "  admin-console      - Admin Console service"
    echo "  all                - All services"
    echo ""
    echo "Examples:"
    echo "  ./scopeapi.sh setup --full                    # Complete setup"
    echo "  ./scopeapi.sh start infrastructure            # Start infrastructure only"
    echo "  ./scopeapi.sh start api-discovery             # Start infrastructure + API Discovery"
    echo "  ./scopeapi.sh start all                       # Start everything"
    echo "  ./scopeapi.sh stop                            # Stop all services"
    echo "  ./scopeapi.sh logs api-discovery              # Show logs for API Discovery"
    echo "  ./scopeapi.sh debug api-discovery             # Debug API Discovery"
    echo "  ./scopeapi.sh shell api-discovery             # Shell into API Discovery"
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
    
    # Check if .env.local file exists (only for local development)
    if [ ! -f ".env.local" ]; then
        print_warning "No .env.local file found"
        print_status "Creating .env.local file from template for LOCAL DEVELOPMENT ONLY..."
        if [ -f "$SCRIPT_DIR/env.example" ]; then
            cp "$SCRIPT_DIR/env.example" .env.local
            print_success "Created .env.local file from env.example"
            print_warning "⚠️  IMPORTANT: .env.local is for LOCAL DEVELOPMENT ONLY!"
            print_warning "Please edit .env.local file and set your local passwords"
            print_warning "Then run this script again"
            print_info "For staging/production, use: ./deploy.sh -e staging -p k8s"
            exit 1
        else
            print_error "env.example file not found"
            print_status "Please create a .env.local file manually with required environment variables"
            exit 1
        fi
    fi
    
    # Use .env.local for local development
    ENV_FILE=".env.local"
    print_success ".env.local file found (LOCAL DEVELOPMENT MODE)"
    print_info "⚠️  Remember: .env.local is for your local machine only!"
    
    # Load environment variables
    print_status "Loading environment variables from $ENV_FILE..."
    set -a  # automatically export all variables
    source "$ENV_FILE"
    set +a  # stop automatically exporting
    print_success "Environment variables loaded"
}

# Function to setup infrastructure
setup_infrastructure() {
    print_header "Setting Up Infrastructure Services"
    
    print_status "Starting infrastructure services..."
    
    # Start infrastructure using docker-compose
    if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" up -d zookeeper kafka postgres redis elasticsearch kibana; then
        print_success "Infrastructure services started successfully"
        
        # Wait for services to be ready
        print_status "Waiting for services to be ready..."
        sleep 15
        
        # Check service status
        if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" ps | grep -q "Up"; then
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
    if docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" ps | grep -q "Up"; then
        print_success "All services are running"
    else
        print_warning "Some services may not be running"
    fi
}

# Function to start services
start_services() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified. Use 'all' or specify individual services."
        exit 1
    fi
    
    print_status "Starting services: ${services[*]}"
    
    # Always start infrastructure first
    print_status "Starting infrastructure services..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" up -d zookeeper kafka postgres redis elasticsearch kibana
    
    # Wait for infrastructure to be ready
    print_status "Waiting for infrastructure to be ready..."
    sleep 10
    
    # Start specified services
    for service in "${services[@]}"; do
        if [ "$service" = "infrastructure" ]; then
            print_success "Infrastructure services started"
            continue
        fi
        
        if [ "$service" = "all" ]; then
            print_status "Starting all microservices..."
            docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" up -d api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
            break
        fi
        
        print_status "Starting $service..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" up -d "$service"
    done
    
    print_success "Services started successfully!"
    print_status "Use './scopeapi.sh status' to check service status"
    print_status "Use './scopeapi.sh logs [service]' to view logs"
}

# Function to stop services
stop_services() {
    print_status "Stopping all services..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" down
    print_success "All services stopped"
}

# Function to restart services
restart_services() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified for restart."
        exit 1
    fi
    
    print_status "Restarting services: ${services[*]}"
    
    for service in "${services[@]}"; do
        print_status "Restarting $service..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" restart "$service"
    done
    
    print_success "Services restarted successfully!"
}

# Function to show logs
show_logs() {
    local service="$1"
    
    if [ -z "$service" ]; then
        print_status "Showing logs for all services..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" logs -f
    else
        print_status "Showing logs for $service..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" logs -f "$service"
    fi
}

# Function to show comprehensive status
show_status() {
    print_header "ScopeAPI Services Status"
    
    echo ""
    print_status "=== INFRASTRUCTURE STATUS ==="
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" ps zookeeper kafka postgres redis elasticsearch kibana
    
    echo ""
    print_status "=== MICROSERVICES STATUS ==="
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" ps api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
    
    echo ""
    print_status "=== SYSTEM RESOURCES ==="
    docker system df
    
    echo ""
    print_status "=== NETWORK STATUS ==="
    docker network ls | grep scopeapi
}

# Function to build services
build_services() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified for building."
        exit 1
    fi
    
    print_status "Building services: ${services[*]}"
    
    for service in "${services[@]}"; do
        print_status "Building $service..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" build "$service"
    done
    
    print_success "Services built successfully!"
}

# Function to clean everything
clean_all() {
    print_warning "This will remove ALL containers, volumes, and images. Are you sure? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Cleaning all containers, volumes, and images..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" down -v --rmi all
        docker system prune -af
        print_success "Cleanup completed!"
    else
        print_status "Cleanup cancelled"
    fi
}

# Function to open shell in service container
open_shell() {
    local service="$1"
    
    if [ -z "$service" ]; then
        print_error "No service specified for shell access."
        exit 1
    fi
    
    print_status "Opening shell in $service container..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" exec "$service" sh
}

# Function to execute command in service container
execute_command() {
    local service="$1"
    local command="$2"
    
    if [ -z "$service" ] || [ -z "$command" ]; then
        print_error "Usage: ./scopeapi.sh exec <service> <command>"
        print_error "Example: $0 exec api-discovery 'ps aux'"
        exit 1
    fi
    
    print_status "Executing '$command' in $service container..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" exec "$service" sh -c "$command"
}

# Function to start debug mode
start_debug() {
    local service="$1"
    
    if [ -z "$service" ]; then
        print_error "No service specified for debug mode."
        exit 1
    fi
    
    print_status "Starting $service in debug mode..."
    
    # Stop the service first
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" stop "$service"
    
    # Start in debug mode (assuming debug configuration exists)
    print_status "Starting $service with debug configuration..."
    docker-compose -f "$SCRIPT_DIR/docker-compose.debug.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" up -d "$service"
    
    print_success "$service started in debug mode!"
    print_info "Connect your debugger to localhost:2345"
}

# Function to show comprehensive status
show_comprehensive_status() {
    print_header "Comprehensive Status"
    
    print_status "Infrastructure Services:"
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
    
    echo ""
    print_status "Network Information:"
    docker network ls | grep scopeapi || echo "No scopeapi networks found"
    
    echo ""
    print_status "Volume Information:"
    docker volume ls | grep scopeapi || echo "No scopeapi volumes found"
    
    echo ""
    print_status "Resource Usage:"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
}

# Function to cleanup services
cleanup_services() {
    print_header "Cleaning Up Services"
    print_status "Stopping and removing services..."
    
    docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" down
    print_success "Services cleaned up successfully!"
}

# Function to cleanup everything
cleanup_full() {
    print_header "Full Cleanup"
    print_warning "This will remove ALL containers, volumes, and networks!"
    
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Removing all containers, volumes, and networks..."
        docker-compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$PROJECT_ROOT/$ENV_FILE" down -v --remove-orphans
        docker system prune -f
        print_success "Full cleanup completed!"
    else
        print_status "Cleanup cancelled"
    fi
}

# Function to show final status after setup
show_final_status() {
    print_header "Setup Complete!"
    
    print_success "ScopeAPI setup has been completed successfully!"
    echo ""
    print_status "Next steps:"
    echo "  1. Start microservices: ./scopeapi.sh start all"
    echo "  2. Or start specific service: ./scopeapi.sh start api-discovery"
    echo "  3. For debugging: ./scopeapi.sh debug api-discovery"
    echo ""
    print_status "Service URLs:"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Kafka: localhost:9092"
    echo "  - Redis: localhost:6379"
    echo "  - Elasticsearch: localhost:9200"
    echo "  - Kibana: localhost:5601"
    echo ""
    print_status "Useful commands:"
    echo "  - View logs: ./scopeapi.sh logs [service]"
    echo "  - Check status: ./scopeapi.sh status"
    echo "  - Stop services: ./scopeapi.sh stop"
}

# Main execution
main() {
    local command="$1"
    shift
    
    case "$command" in
        setup)
            local setup_type=""
            local setup_flags=()
            
            # Parse setup options
            while [[ $# -gt 0 ]]; do
                case $1 in
                    --full)
                        setup_type="full"
                        shift
                        ;;
                    --infrastructure)
                        setup_type="infrastructure"
                        shift
                        ;;
                    --database)
                        setup_type="database"
                        shift
                        ;;
                    --test-data)
                        setup_type="test-data"
                        shift
                        ;;
                    --validate)
                        setup_type="validate"
                        shift
                        ;;
                    --cleanup)
                        setup_type="cleanup"
                        shift
                        ;;
                    --cleanup-full)
                        setup_type="cleanup-full"
                        shift
                        ;;
                    *)
                        setup_flags+=("$1")
                        shift
                        ;;
                esac
            done
            
            # Default to full setup if no specific type specified
            if [ -z "$setup_type" ]; then
                setup_type="full"
            fi
            
            print_header "ScopeAPI Setup"
            print_status "Starting setup process..."
            
            # Check prerequisites
            check_prerequisites
            
            # Perform setup based on type
            case "$setup_type" in
                full)
                    setup_infrastructure
                    setup_database
                    create_test_data
                    validate_setup
                    show_final_status
                    ;;
                infrastructure)
                    setup_infrastructure
                    ;;
                database)
                    setup_database
                    ;;
                test-data)
                    setup_infrastructure
                    setup_database
                    create_test_data
                    ;;
                validate)
                    setup_infrastructure
                    setup_database
                    validate_setup
                    ;;
                cleanup)
                    cleanup_services
                    ;;
                cleanup-full)
                    cleanup_full
                    ;;
            esac
            ;;
        start)
            check_prerequisites
            start_services "$@"
            ;;
        stop)
            check_prerequisites
            stop_services
            ;;
        restart)
            check_prerequisites
            restart_services "$@"
            ;;
        logs)
            check_prerequisites
            show_logs "$@"
            ;;
        status)
            check_prerequisites
            show_status
            ;;
        comprehensive-status)
            check_prerequisites
            show_comprehensive_status
            ;;
        build)
            check_prerequisites
            build_services "$@"
            ;;
        clean)
            check_prerequisites
            clean_all
            ;;
        shell)
            check_prerequisites
            open_shell "$@"
            ;;
        exec)
            check_prerequisites
            execute_command "$@"
            ;;
        debug)
            check_prerequisites
            start_debug "$@"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
