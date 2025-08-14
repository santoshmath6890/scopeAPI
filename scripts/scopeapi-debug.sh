#!/bin/bash

# ScopeAPI Debug Script
# Purpose: Start microservices in debug mode with Delve debugger
# Usage: ./scopeapi-debug.sh start [service]
# Features: Debug ports, Delve integration, interactive debugging

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

# Function to show help
show_help() {
    echo "ScopeAPI Debug Script"
echo ""
echo "Usage: ./scopeapi-debug.sh [COMMAND] [SERVICES...]"
echo ""
echo "Commands:"
echo "  start [services...]  - Start services in debug & development mode"
echo "  stop                 - Stop debug services"
echo "  restart [services...] - Restart debug services"
echo "  logs [service]       - Show debug service logs"
echo "  status               - Show debug service status"
echo "  build [services...]  - Build debug images"
echo "  clean                - Clean debug containers and images"
echo "  help                 - Show this help message"
    echo ""
    echo "Services:"
    echo "  api-discovery        - API Discovery service (port 2345)"
    echo "  gateway-integration  - Gateway Integration service (port 2346)"
    echo "  data-ingestion      - Data Ingestion service (port 2347)"
    echo "  threat-detection    - Threat Detection service (port 2348)"
    echo "  data-protection     - Data Protection service (port 2349)"
    echo "  attack-blocking     - Attack Blocking service (port 2350)"
    echo "  admin-console       - Admin Console service (port 2351)"
    echo "  all                 - All services"
    echo ""
    echo "Examples:"
    echo "  $0 start api-discovery           # Start API Discovery in debug mode"
    echo "  $0 start api-discovery gateway-integration  # Start 2 services in debug mode"
    echo "  $0 start all                     # Start all services in debug mode"
    echo "  $0 stop                          # Stop all debug services"
    echo "  $0 logs api-discovery            # Show debug logs for API Discovery"
    echo "  $0 build api-discovery           # Build debug image for API Discovery"
}

# Function to check prerequisites
check_prerequisites() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose is not installed"
        exit 1
    fi
    
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running"
        exit 1
    fi
}

# Function to start debug services
start_debug_services() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified for debugging"
        exit 1
    fi
    
    print_status "Starting services in debug mode: ${services[*]}"
    
    # Start infrastructure first if not running
    if ! docker-compose ps | grep -q "postgres.*Up"; then
        print_status "Starting infrastructure services..."
        docker-compose up -d zookeeper kafka postgres redis elasticsearch kibana
        sleep 10
    fi
    
    # Start specified debug services
    for service in "${services[@]}"; do
        if [ "$service" = "all" ]; then
            print_status "Starting all services in debug mode..."
            docker-compose -f docker-compose.yml -f docker-compose.debug.yml up -d
            break
        else
            print_status "Starting $service in debug mode..."
            docker-compose -f docker-compose.yml -f docker-compose.debug.yml up -d "$service"
        fi
    done
    
    print_success "Debug services started successfully!"
    print_info "Debug ports available:"
    print_info "  API Discovery: localhost:2345"
    print_info "  Gateway Integration: localhost:2346"
    print_info "  Data Ingestion: localhost:2347"
    print_info "  Threat Detection: localhost:2348"
    print_info "  Data Protection: localhost:2349"
    print_info "  Attack Blocking: localhost:2350"
    print_info "  Admin Console: localhost:2351"
    print_info ""
    print_info "Connect your IDE to the appropriate port for debugging"
}

# Function to stop debug services
stop_debug_services() {
    print_status "Stopping debug services..."
    docker-compose -f docker-compose.yml -f docker-compose.debug.yml down
    print_success "Debug services stopped"
}

# Function to restart debug services
restart_debug_services() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified for restart"
        exit 1
    fi
    
    print_status "Restarting debug services: ${services[*]}"
    
    for service in "${services[@]}"; do
        print_status "Restarting $service..."
        docker-compose -f docker-compose.yml -f docker-compose.debug.yml restart "$service"
    done
    
    print_success "Debug services restarted successfully!"
}

# Function to show debug service logs
show_debug_logs() {
    local service="$1"
    
    if [ -z "$service" ]; then
        print_status "Showing logs for all debug services..."
        docker-compose -f docker-compose.yml -f docker-compose.debug.yml logs -f
    else
        print_status "Showing logs for $service..."
        docker-compose -f docker-compose.yml -f docker-compose.debug.yml logs -f "$service"
    fi
}

# Function to show debug service status
show_debug_status() {
    print_status "Debug service status:"
    docker-compose -f docker-compose.yml -f docker-compose.debug.yml ps
}

# Function to build debug images
build_debug_images() {
    local services=("$@")
    
    if [ ${#services[@]} -eq 0 ]; then
        print_error "No services specified for building"
        exit 1
    fi
    
    print_status "Building debug images: ${services[*]}"
    
    for service in "${services[@]}"; do
        print_status "Building debug image for $service..."
        docker-compose -f docker-compose.yml -f docker-compose.debug.yml build "$service"
    done
    
    print_success "Debug images built successfully!"
}

# Function to clean debug environment
clean_debug_environment() {
    print_warning "This will remove ALL debug containers and images. Are you sure? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Cleaning debug environment..."
        docker-compose -f docker-compose.yml -f docker-compose.debug.yml down -v --rmi all
        docker system prune -af
        print_success "Debug environment cleaned successfully!"
    else
        print_status "Cleanup cancelled"
    fi
}

# Main function
main() {
    local command="$1"
    shift
    
    # Check prerequisites
    check_prerequisites
    
    case "$command" in
        start)
            start_debug_services "$@"
            ;;
        stop)
            stop_debug_services
            ;;
        restart)
            restart_debug_services "$@"
            ;;
        logs)
            show_debug_logs "$@"
            ;;
        status)
            show_debug_status
            ;;
        build)
            build_debug_images "$@"
            ;;
        clean)
            clean_debug_environment
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
