#!/bin/bash

# ScopeAPI Microservices Orchestration Script
# Purpose: Complete container-based microservices management
# Usage: ./scopeapi-services.sh start [service]
# Features: Infrastructure + microservices orchestration, container management, debugging

set -e

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"
ENV_FILE="$PROJECT_ROOT/.env"

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

# Function to show help
show_help() {
    echo "ScopeAPI Services Script"
    echo ""
    echo "Usage: ./scopeapi-services.sh [COMMAND] [SERVICES...]"
    echo ""
    echo "Commands:"
    echo "  start [services...]  - Start infrastructure + specified services"
    echo "  stop                 - Stop all services"
    echo "  restart [services...] - Restart services"
    echo "  logs [service]       - Show logs for service(s)"
    echo "  status               - Show status of all services"
    echo "  build [services...]  - Build specified services"
    echo "  clean                - Clean all containers and volumes"
    echo "  infrastructure        - Start only infrastructure services"
    echo "  comprehensive-status - Show detailed infrastructure + microservices status"

    
    echo "  shell [service]      - Open shell in service container"
    echo "  exec [service] [cmd] - Execute command in service container"
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
    echo "  $0 start infrastructure           # Start only infrastructure"
    echo "  $0 start api-discovery            # Start infrastructure + API Discovery"
    echo "  $0 start api-discovery gateway-integration  # Start infrastructure + 2 services"
    echo "  $0 start all                     # Start everything"
    echo "  $0 stop                          # Stop all services"
    echo "  $0 logs api-discovery            # Show logs for API Discovery"
    echo "  $0 build api-discovery           # Build API Discovery service"
    echo "  $0 infrastructure                  # Start only infrastructure"
    echo "  $0 comprehensive-status             # Show detailed status"

    
    echo "  $0 shell api-discovery           # Open shell in API Discovery container"
    echo "  $0 exec api-discovery ps aux     # Execute command in container"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Function to check if docker-compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose is not installed. Please install it and try again."
        exit 1
    fi
}

# Function to check configuration and load environment
check_configuration() {
    # Check if compose file exists
    if [ ! -f "$COMPOSE_FILE" ]; then
        print_error "Docker Compose file not found: $COMPOSE_FILE"
        exit 1
    fi
    
    # Check if .env file exists
    if [ ! -f "$ENV_FILE" ]; then
        print_error "Environment file not found: $ENV_FILE"
        exit 1
    fi
    
    # Load environment variables
    print_status "Loading environment variables from $ENV_FILE"
    set -a  # automatically export all variables
    source "$ENV_FILE"
    set +a  # stop automatically exporting
    print_success "Environment variables loaded"
}

# Function to start services
# Function to start infrastructure only
start_infrastructure_only() {
    print_status "Starting infrastructure services only..."
    
    # Start core infrastructure
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d zookeeper kafka postgres redis elasticsearch kibana
    
    print_status "Waiting for infrastructure to be ready..."
    sleep 15
    
    # Verify infrastructure health
    print_status "Verifying infrastructure health..."
    
    # Check PostgreSQL
    if docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres pg_isready -U scopeapi > /dev/null 2>&1; then
        print_success "PostgreSQL is ready"
    else
        print_warning "PostgreSQL may still be starting up..."
    fi
    
    # Check Kafka
    if docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
        print_success "Kafka is ready"
    else
        print_warning "Kafka may still be starting up..."
    fi
    
    # Check Redis
    if docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T redis redis-cli ping > /dev/null 2>&1; then
        print_success "Redis is ready"
    else
        print_warning "Redis may still be starting up..."
    fi
    
    print_success "Infrastructure services started successfully!"
    print_status "Use 047$0 status047 to check service status"
}

start_services() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified. Use 'all' or specify individual services."
        exit 1
    fi
    
    print_status "Starting services: ${services[*]}"
    
    # Always start infrastructure first
    print_status "Starting infrastructure services..."
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d zookeeper kafka postgres redis elasticsearch kibana
    
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
            docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
            break
        fi
        
        print_status "Starting $service..."
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d "$service"
    done
    
    print_success "Services started successfully!"
    print_status "Use '$0 status' to check service status"
    print_status "Use '$0 logs [service]' to view logs"
}

# Function to stop services
stop_services() {
    print_status "Stopping all services..."
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down
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
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" restart "$service"
    done
    
    print_success "Services restarted successfully!"
}

# Function to show logs
show_logs() {
    local service="$1"
    
    if [ -z "$service" ]; then
        print_status "Showing logs for all services..."
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs -f
    else
        print_status "Showing logs for $service..."
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs -f "$service"
    fi
}

# Function to show status
# Function to show comprehensive status
show_comprehensive_status() {
    print_status "=== INFRASTRUCTURE STATUS ==="
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps zookeeper kafka postgres redis elasticsearch kibana
    
    echo ""
    print_status "=== MICROSERVICES STATUS ==="
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps api-discovery gateway-integration data-ingestion threat-detection data-protection attack-blocking admin-console
    
    echo ""
    print_status "=== SYSTEM RESOURCES ==="
    docker system df
    
    echo ""
    print_status "=== NETWORK STATUS ==="
    docker network ls | grep scopeapi
}

show_status() {
    print_status "Service status:"
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps
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
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build "$service"
    done
    
    print_success "Services built successfully!"
}

# Function to clean everything
clean_all() {
    print_warning "This will remove ALL containers, volumes, and images. Are you sure? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Cleaning all containers, volumes, and images..."
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down -v --rmi all
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
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec "$service" sh
}

# Function to execute command in service container
execute_command() {
    local service="$1"
    local command="$2"
    
    if [ -z "$service" ] || [ -z "$command" ]; then
        print_error "Usage: ./scopeapi-services.sh exec <service> <command>"
        print_error "Example: $0 exec api-discovery 'ps aux'"
        exit 1
    fi
    
    print_status "Executing '$command' in $service container..."
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec "$service" sh -c "$command"
}

# Main script logic
main() {
    local command="$1"
    shift
    
    # Check prerequisites
    check_docker
    check_docker_compose
    check_configuration
    
    case "$command" in
        start)
            start_services "$@"
            ;;
        infrastructure)
            start_infrastructure_only
            ;;
        comprehensive-status)
            show_comprehensive_status
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services "$@"
            ;;
        logs)
            show_logs "$@"
            ;;
        status)
            show_status
            ;;
        build)
            build_services "$@"
            ;;
        clean)
            clean_all
            ;;
        shell)
            open_shell "$@"
            ;;
        exec)
            execute_command "$@"
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
