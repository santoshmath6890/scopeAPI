#!/bin/bash

# ScopeAPI Development Workflow Script
# Purpose: Daily development tasks, debugging, testing, and development utilities
# Usage: ./dev.sh [COMMAND] [OPTIONS]
# Features: Development workflows, debugging, testing, code quality

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
    echo "ScopeAPI Development Workflow Script"
    echo ""
    echo "Usage: ./dev.sh [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  start [SERVICES...]  - Start development environment"
    echo "  stop                 - Stop development environment"
    echo "  restart [SERVICES...] - Restart services"
    echo "  logs [SERVICE]       - Show logs for service(s)"
    echo "  status               - Show development environment status"
    echo "  debug [SERVICE]      - Start service in debug mode"
    echo "  test                 - Run all tests"
    echo "  test-backend         - Run backend tests only"
    echo "  test-frontend        - Run frontend tests only"
    echo "  build                - Build all services"
    echo "  build-backend        - Build backend services only"
    echo "  build-frontend       - Build frontend only"
    echo "  clean                - Clean development environment"
    echo "  shell [SERVICE]      - Open shell in service container"
    echo "  exec [SERVICE] [CMD] - Execute command in container"
    echo "  lint                 - Run linting and code quality checks"
    echo "  format               - Format code"
    echo "  help                 - Show this help message"
    echo ""
    echo "Services:"
    echo "  infrastructure       - Start only infrastructure (postgres, kafka, redis, etc.)"
    echo "  api-discovery        - API Discovery service"
    echo "  gateway-integration  - Gateway Integration service"
    echo "  data-ingestion      - Data Ingestion service"
    echo "  threat-detection    - Threat Detection service"
    echo "  data-protection     - Data Protection service"
    echo "  attack-blocking     - Attack Blocking service"
    echo "  admin-console       - Admin Console service"
    echo "  all                 - All services"
    echo ""
    echo "Examples:"
    echo "  ./dev.sh start all                    # Start complete development environment"
    echo "  ./dev.sh start infrastructure          # Start infrastructure only"
    echo "  ./dev.sh debug api-discovery          # Debug API Discovery service"
    echo "  ./dev.sh test                         # Run all tests"
    echo "  ./dev.sh logs api-discovery           # Show API Discovery logs"
    echo "  ./dev.sh shell api-discovery          # Shell into API Discovery"
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking development prerequisites..."
    
    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    print_success "Docker is running"
    
    # Check if docker compose is available
    if ! docker compose version &> /dev/null; then
        print_error "docker compose plugin is not installed. Please install it first."
        exit 1
    fi
    print_success "docker compose is available"
    
    # Check if .env.local file exists
    if [ ! -f ".env.local" ]; then
        print_warning "No .env.local file found"
        print_status "Creating .env.local file from template for LOCAL DEVELOPMENT ONLY..."
        if [ -f "$SCRIPT_DIR/env.example" ]; then
            cp "$SCRIPT_DIR/env.example" .env.local
            print_success "Created .env.local file from env.example"
            print_warning "⚠️  IMPORTANT: .env.local is for LOCAL DEVELOPMENT ONLY!"
            print_warning "Please edit .env.local file and set your local passwords"
            print_warning "Then run this script again"
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
    
    # Load environment variables
    print_status "Loading environment variables from $ENV_FILE..."
    set -a  # automatically export all variables
    source "$ENV_FILE"
    set +a  # stop automatically exporting
    print_success "Environment variables loaded"
}

# Function to start development environment
start_dev_environment() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified. Use 'all' or specify individual services."
        exit 1
    fi
    
    print_status "Starting development environment: ${services[*]}"
    
    # Always start infrastructure first
    print_status "Starting infrastructure services..."
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" up -d zookeeper kafka postgres redis elasticsearch kibana
    
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
            docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" up -d api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
            break
        fi
        
        print_status "Starting $service..."
        docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" up -d "$service"
    done
    
    print_success "Development environment started successfully!"
    print_info "Use './dev.sh status' to check service status"
    print_info "Use './dev.sh logs [service]' to view logs"
}

# Function to stop development environment
stop_dev_environment() {
    print_status "Stopping development environment..."
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" down
    print_success "Development environment stopped"
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
        docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" restart "$service"
    done
    
    print_success "Services restarted successfully!"
}

# Function to show logs
show_logs() {
    local service="$1"
    
    if [ -z "$service" ]; then
        print_status "Showing logs for all services..."
        docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" logs -f
    else
        print_status "Showing logs for $service..."
        docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" logs -f "$service"
    fi
}

# Function to show status
show_status() {
    print_header "Development Environment Status"
    
    echo ""
    print_status "=== INFRASTRUCTURE STATUS ==="
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" ps zookeeper kafka postgres redis elasticsearch kibana
    
    echo ""
    print_status "=== MICROSERVICES STATUS ==="
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" ps api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
    
    echo ""
    print_status "=== SYSTEM RESOURCES ==="
    docker system df
    
    echo ""
    print_status "=== NETWORK STATUS ==="
    docker network ls | grep scopeapi
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
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" stop "$service"
    
    # Start in debug mode (assuming debug configuration exists)
    print_status "Starting $service with debug configuration..."
    docker compose -f "$SCRIPT_DIR/docker-compose.debug.yml" --env-file "$ENV_FILE" up -d "$service"
    
    print_success "$service started in debug mode!"
    print_info "Connect your debugger to localhost:2345"
}

# Function to run tests
run_tests() {
    local test_type="$1"
    
    case "$test_type" in
        backend)
            print_header "Running Backend Tests"
            print_status "Running Go tests..."
            cd backend
            go test ./...
            cd ..
            print_success "Backend tests completed"
            ;;
        frontend)
            print_header "Running Frontend Tests"
            print_status "Running Angular tests..."
            cd adminConsole
            npm test
            cd ..
            print_success "Frontend tests completed"
            ;;
        *)
            print_header "Running All Tests"
            print_status "Running backend tests..."
            cd backend
            go test ./...
            cd ..
            
            print_status "Running frontend tests..."
            cd adminConsole
            npm test
            cd ..
            
            print_success "All tests completed"
            ;;
    esac
}

# Function to build services
build_services() {
    local build_type="$1"
    
    case "$build_type" in
        backend)
            print_header "Building Backend Services"
            print_status "Building Go services..."
            cd backend
            go build ./...
            cd ..
            print_success "Backend services built"
            ;;
        frontend)
            print_header "Building Frontend"
            print_status "Building Angular application..."
            cd adminConsole
            npm run build
            cd ..
            print_success "Frontend built"
            ;;
        *)
            print_header "Building All Services"
            print_status "Building backend services..."
            cd backend
            go build ./...
            cd ..
            
            print_status "Building frontend..."
            cd adminConsole
            npm run build
            cd ..
            
            print_success "All services built"
            ;;
    esac
}

# Function to clean development environment
clean_dev_environment() {
    print_warning "This will remove ALL containers, volumes, and images. Are you sure? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Cleaning development environment..."
        docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" down -v --rmi all
        docker system prune -af
        print_success "Development environment cleaned"
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
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" exec "$service" sh
}

# Function to execute command in service container
execute_command() {
    local service="$1"
    local command="$2"
    
    if [ -z "$service" ] || [ -z "$command" ]; then
        print_error "Usage: ./dev.sh exec <service> <command>"
        print_error "Example: $0 exec api-discovery 'ps aux'"
        exit 1
    fi
    
    print_status "Executing '$command' in $service container..."
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" --env-file "$ENV_FILE" exec "$service" sh -c "$command"
}

# Function to run linting and code quality checks
run_linting() {
    print_header "Running Code Quality Checks"
    
    print_status "Running Go linting..."
    if command -v golangci-lint &> /dev/null; then
        cd backend
        golangci-lint run
        cd ..
        print_success "Go linting completed"
    else
        print_warning "golangci-lint not found, skipping Go linting"
    fi
    
    print_status "Running Angular linting..."
    cd adminConsole
    npm run lint
    cd ..
    print_success "Angular linting completed"
}

# Function to format code
format_code() {
    print_header "Formatting Code"
    
    print_status "Formatting Go code..."
    cd backend
    go fmt ./...
    cd ..
    print_success "Go code formatted"
    
    print_status "Formatting Angular code..."
    cd adminConsole
    npm run format
    cd ..
    print_success "Angular code formatted"
}

# Main execution
main() {
    local command="$1"
    shift
    
    case "$command" in
        start)
            check_prerequisites
            start_dev_environment "$@"
            ;;
        stop)
            check_prerequisites
            stop_dev_environment
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
        debug)
            check_prerequisites
            start_debug "$@"
            ;;
        test)
            run_tests "$@"
            ;;
        test-backend)
            run_tests "backend"
            ;;
        test-frontend)
            run_tests "frontend"
            ;;
        build)
            build_services "$@"
            ;;
        build-backend)
            build_services "backend"
            ;;
        build-frontend)
            build_services "frontend"
            ;;
        clean)
            check_prerequisites
            clean_dev_environment
            ;;
        shell)
            check_prerequisites
            open_shell "$@"
            ;;
        exec)
            check_prerequisites
            execute_command "$@"
            ;;
        lint)
            run_linting
            ;;
        format)
            format_code
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
