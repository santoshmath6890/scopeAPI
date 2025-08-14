#!/bin/bash

# ScopeAPI Local Development Manager
# Manages local Go microservices processes for development
# 
# Usage: ./scopeapi-local.sh [command] [service]
# Commands: start, stop, restart, status, logs, build, clean
# Services: all, api-discovery, gateway-integration, data-ingestion, threat-detection, data-protection, attack-blocking

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"
SERVICES_DIR="$BACKEND_DIR/services"
BIN_DIR="$BACKEND_DIR/bin"
LOGS_DIR="$PROJECT_ROOT/logs"

# Service definitions
declare -A SERVICES=(
    ["api-discovery"]="8080"
    ["gateway-integration"]="8081"
    ["data-ingestion"]="8082"
    ["threat-detection"]="8083"
    ["data-protection"]="8084"
    ["attack-blocking"]="8085"
)

# PID file directory
PID_DIR="$LOGS_DIR/pids"
mkdir -p "$PID_DIR"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if service is running
is_service_running() {
    local service=$1
    local pid_file="$PID_DIR/${service}.pid"
    
    if [[ -f "$pid_file" ]]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        else
            rm -f "$pid_file"
            return 1
        fi
    fi
    return 1
}

# Get service status
get_service_status() {
    local service=$1
    if is_service_running "$service"; then
        echo "running"
    else
        echo "stopped"
    fi
}

# Build a service
build_service() {
    local service=$1
    local service_dir="$SERVICES_DIR/$service"
    
    if [[ ! -d "$service_dir" ]]; then
        log_error "Service directory not found: $service_dir"
        return 1
    fi
    
    log_info "Building $service..."
    cd "$service_dir"
    
    if go build -o "$BIN_DIR/$service" ./cmd/main.go; then
        log_success "$service built successfully"
    else
        log_error "Failed to build $service"
        return 1
    fi
}

# Start a service
start_service() {
    local service=$1
    local service_dir="$SERVICES_DIR/$service"
    local bin_path="$BIN_DIR/$service"
    local pid_file="$PID_DIR/${service}.pid"
    local log_file="$LOGS_DIR/${service}.log"
    
    if is_service_running "$service"; then
        log_warning "$service is already running"
        return 0
    fi
    
    if [[ ! -f "$bin_path" ]]; then
        log_info "Binary not found, building $service..."
        build_service "$service"
    fi
    
    log_info "Starting $service on port ${SERVICES[$service]}..."
    
    # Create logs directory if it doesn't exist
    mkdir -p "$LOGS_DIR"
    
    # Start service in background
    nohup "$bin_path" > "$log_file" 2>&1 &
    local pid=$!
    
    # Save PID
    echo "$pid" > "$pid_file"
    
    # Wait a moment to check if it started successfully
    sleep 2
    if kill -0 "$pid" 2>/dev/null; then
        log_success "$service started successfully (PID: $pid)"
    else
        log_error "Failed to start $service"
        rm -f "$pid_file"
        return 1
    fi
}

# Stop a service
stop_service() {
    local service=$1
    local pid_file="$PID_DIR/${service}.pid"
    
    if [[ ! -f "$pid_file" ]]; then
        log_warning "$service is not running"
        return 0
    fi
    
    local pid=$(cat "$pid_file")
    log_info "Stopping $service (PID: $pid)..."
    
    if kill "$pid" 2>/dev/null; then
        log_success "$service stopped successfully"
    else
        log_warning "Failed to stop $service gracefully, force killing..."
        kill -9 "$pid" 2>/dev/null || true
    fi
    
    rm -f "$pid_file"
}

# Restart a service
restart_service() {
    local service=$1
    log_info "Restarting $service..."
    stop_service "$service"
    sleep 1
    start_service "$service"
}

# Show service logs
show_logs() {
    local service=$1
    local log_file="$LOGS_DIR/${service}.log"
    
    if [[ ! -f "$log_file" ]]; then
        log_error "No log file found for $service"
        return 1
    fi
    
    log_info "Showing logs for $service:"
    echo "----------------------------------------"
    tail -f "$log_file"
}

# Show all services status
show_status() {
    echo -e "\n${BLUE}ScopeAPI Local Services Status:${NC}"
    echo "=========================================="
    
    for service in "${!SERVICES[@]}"; do
        local status=$(get_service_status "$service")
        local port="${SERVICES[$service]}"
        
        if [[ "$status" == "running" ]]; then
            local pid_file="$PID_DIR/${service}.pid"
            local pid=$(cat "$pid_file" 2>/dev/null || echo "N/A")
            echo -e "${GREEN}✓${NC} $service (port: $port, PID: $pid)"
        else
            echo -e "${RED}✗${NC} $service (port: $port, status: stopped)"
        fi
    done
    
    echo ""
}

# Clean up all services and binaries
clean_all() {
    log_warning "This will stop all services and remove binaries. Continue? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        log_info "Cleaning up..."
        
        # Stop all services
        for service in "${!SERVICES[@]}"; do
            if is_service_running "$service"; then
                stop_service "$service"
            fi
        done
        
        # Remove binaries
        if [[ -d "$BIN_DIR" ]]; then
            rm -rf "$BIN_DIR"
            log_success "Binaries removed"
        fi
        
        # Remove PID files
        if [[ -d "$PID_DIR" ]]; then
            rm -rf "$PID_DIR"
            log_success "PID files removed"
        fi
        
        log_success "Cleanup completed"
    else
        log_info "Cleanup cancelled"
    fi
}

# Main command handler
case "${1:-help}" in
    "start")
        if [[ -z "$2" || "$2" == "all" ]]; then
            log_info "Starting all services..."
            for service in "${!SERVICES[@]}"; do
                start_service "$service"
                sleep 1
            done
            log_success "All services started"
        else
            start_service "$2"
        fi
        ;;
    
    "stop")
        if [[ -z "$2" || "$2" == "all" ]]; then
            log_info "Stopping all services..."
            for service in "${!SERVICES[@]}"; do
                stop_service "$service"
            done
            log_success "All services stopped"
        else
            stop_service "$2"
        fi
        ;;
    
    "restart")
        if [[ -z "$2" || "$2" == "all" ]]; then
            log_info "Restarting all services..."
            for service in "${!SERVICES[@]}"; do
                restart_service "$service"
                sleep 1
            done
            log_success "All services restarted"
        else
            restart_service "$2"
        fi
        ;;
    
    "status")
        show_status
        ;;
    
    "logs")
        if [[ -z "$2" ]]; then
            log_error "Please specify a service name"
            exit 1
        fi
        show_logs "$2"
        ;;
    
    "build")
        if [[ -z "$2" || "$2" == "all" ]]; then
            log_info "Building all services..."
            for service in "${!SERVICES[@]}"; do
                build_service "$service"
            done
            log_success "All services built"
        else
            build_service "$2"
        fi
        ;;
    
    "clean")
        clean_all
        ;;
    
    "help"|*)
        echo -e "${BLUE}ScopeAPI Local Development Manager${NC}"
        echo "======================================"
        echo ""
        echo "Usage: $0 [command] [service]"
        echo ""
        echo "Commands:"
        echo "  start [service|all]    Start service(s) (default: all)"
        echo "  stop [service|all]     Stop service(s) (default: all)"
        echo "  restart [service|all]  Restart service(s) (default: all)"
        echo "  status                 Show status of all services"
        echo "  logs <service>         Show logs for specific service"
        echo "  build [service|all]    Build service(s) (default: all)"
        echo "  clean                  Clean up all services and binaries"
        echo "  help                   Show this help message"
        echo ""
        echo "Services:"
        for service in "${!SERVICES[@]}"; do
            echo "  $service (port: ${SERVICES[$service]})"
        done
        echo "  all                    All services"
        echo ""
        echo "Examples:"
        echo "  $0 start                    # Start all services"
        echo "  $0 start api-discovery      # Start specific service"
        echo "  $0 stop all                 # Stop all services"
        echo "  $0 status                   # Show status"
        echo "  $0 logs api-discovery       # Show logs"
        echo "  $0 build all                # Build all services"
        echo ""
        ;;
esac
