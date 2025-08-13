#!/bin/bash

# ScopeAPI Local Development Script
# Manages Go microservices for local development (direct processes)
# Requires infrastructure to be running (use docker-infrastructure.sh first)
# Combines start, stop, and status functionality into a single script

set -e

# Color definitions for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Project configuration
PROJECT_NAME="ScopeAPI"
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="$PROJECT_ROOT/logs/scopeapi.pid"
LOGS_DIR="$PROJECT_ROOT/logs"

# Load environment variables from .env file if it exists
    # Check infrastructure dependencies
    check_infrastructure
    
load_env_vars() {

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
    if [ -f ".env" ]; then

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        print_info "Loading environment variables from .env file"

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        # Use source to properly handle variables with spaces

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        set -a

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        source .env

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        set +a

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
    else

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        print_error ".env file not found. Please create one based on env.example"

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        print_info "Run: ./docker-infrastructure.sh setup-env or ./scripts/generate-passwords.sh"

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
        exit 1

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
    fi

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}
}

# Function to check if infrastructure services are running
check_infrastructure() {
    print_info "Checking infrastructure dependencies..."
    
    # Check PostgreSQL
    if ! docker ps | grep -q "scopeapi-postgres"; then
        print_error "PostgreSQL is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Kafka
    if ! docker ps | grep -q "scopeapi-kafka"; then
        print_error "Kafka is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    # Check Redis
    if ! docker ps | grep -q "scopeapi-redis"; then
        print_error "Redis is not running. Please start infrastructure first:"
        print_info "  ./docker-infrastructure.sh start"
        exit 1
    fi
    
    print_success "Infrastructure services are running"
}

# Service configuration
DATA_INGESTION_PORT=8081
API_DISCOVERY_PORT=8082
THREAT_DETECTION_PORT=8083
GATEWAY_INTEGRATION_PORT=8084
ATTACK_BLOCKING_PORT=8085
DATA_PROTECTION_PORT=8086
ADMIN_CONSOLE_PORT=4200

# Global PID variables
DATA_INGESTION_PID=""
API_DISCOVERY_PID=""
THREAT_DETECTION_PID=""
GATEWAY_INTEGRATION_PID=""
ATTACK_BLOCKING_PID=""
DATA_PROTECTION_PID=""
ADMIN_CONSOLE_PID=""

# Print functions
print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_header() {
    echo -e "${PURPLE}==========================================${NC}"
    echo -e "${PURPLE}$1${NC}"
    echo -e "${PURPLE}==========================================${NC}"
}

# Function to check if a service is running
check_service_status() {
    local service_name="$1"
    local service_pattern="$2"
    local port="$3"
    local health_endpoint="$4"
    
    echo -e "${CYAN}$service_name:${NC}"
    
    # Check if process is running
    if pgrep -f "$service_pattern" >/dev/null; then
        local pid=$(pgrep -f "$service_pattern" | head -1)
        echo -e "  Process: ${GREEN}Running${NC} (PID: $pid)"
        
        # Check port
        if netstat -tlnp 2>/dev/null | grep -q ":$port "; then
            echo -e "  Port: ${GREEN}Active${NC} ($port)"
        else
            echo -e "  Port: ${RED}Not listening${NC} ($port)"
        fi
        
        # Check health endpoint if provided
        if [[ -n "$health_endpoint" ]]; then
            if curl -s "$health_endpoint" >/dev/null 2>&1; then
                echo -e "  Health: ${GREEN}Healthy${NC}"
            else
                echo -e "  Health: ${YELLOW}Unhealthy${NC}"
            fi
        fi
    else
        echo -e "  Process: ${RED}Not running${NC}"
        echo -e "  Port: ${RED}Not active${NC}"
        if [[ -n "$health_endpoint" ]]; then
            echo -e "  Health: ${RED}N/A${NC}"
        fi
    fi
    echo
}

# Function to check PID file
check_pid_file() {
    echo -e "${CYAN}PID File Status:${NC}"
    if [[ -f "$PID_FILE" ]]; then
        echo -e "  File: ${GREEN}Exists${NC} ($PID_FILE)"
        echo -e "  Contents:"
        while IFS= read -r line; do
            if [[ -n "$line" ]]; then
                local service=$(echo "$line" | cut -d'=' -f1)
                local pid=$(echo "$line" | cut -d'=' -f2)
                if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
                    echo -e "    $service: ${GREEN}$pid${NC} (Running)"
                else
                    echo -e "    $service: ${RED}$pid${NC} (Not running)"
                fi
            fi
        done < "$PID_FILE"
    else
        echo -e "  File: ${RED}Not found${NC} ($PID_FILE)"
    fi
    echo
}

# Function to check system resources
check_system_resources() {
    echo -e "${CYAN}System Resources:${NC}"
    
    # Memory Usage
    echo -e "  Memory Usage:"
    free -h | grep -E "^(Mem|Swap)" | while read -r line; do
        local type=$(echo "$line" | awk '{print $1}')
        local total=$(echo "$line" | awk '{print $2}')
        local used=$(echo "$line" | awk '{print $3}')
        local percentage=$(echo "$line" | awk '{if($2 != "0B") printf "%.1f", ($3/$2)*100; else print "0.0"}')
        
        if [[ "$type" == "Mem" ]]; then
            echo -e "    $used/$total used ($percentage%) - Main Memory (RAM)"
        elif [[ "$type" == "Swap" ]]; then
            echo -e "    $used/$total used ($percentage%) - Swap Memory (Virtual)"
        fi
    done
    
    # ScopeAPI Process Resource Usage
    echo -e "  ScopeAPI Process Resource Usage:"
    if [[ -f "$PID_FILE" ]]; then
        while IFS= read -r line; do
            if [[ -n "$line" ]]; then
                local service=$(echo "$line" | cut -d'=' -f1)
                local pid=$(echo "$line" | cut -d'=' -f2)
                if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
                    local mem_usage=$(ps -o rss= -p "$pid" 2>/dev/null | awk '{printf "%.1f", $1/1024}')
                    local cpu_usage=$(ps -o %cpu= -p "$pid" 2>/dev/null | awk '{printf "%.1f", $1}')
                    echo -e "    $service: ${GREEN}${mem_usage}MB RAM${NC}, ${GREEN}${cpu_usage}% CPU${NC}"
                fi
            fi
        done < "$PID_FILE"
    else
        echo -e "    ${RED}No PID file found${NC}"
    fi
    echo
}

# Function to check dependencies
check_dependencies() {
    echo -e "${CYAN}Dependencies:${NC}"
    
    # PostgreSQL
    if command -v pg_isready >/dev/null 2>&1; then
        if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
            echo -e "  PostgreSQL: ${GREEN}Available${NC} (localhost:5432)"
        else
            echo -e "  PostgreSQL: ${RED}Not accessible${NC} (localhost:5432)"
        fi
    else
        echo -e "  PostgreSQL: ${RED}pg_isready not found${NC}"
    fi
    
    # Kafka
    if nc -z localhost 9092 2>/dev/null; then
        echo -e "  Kafka: ${GREEN}Available${NC} (localhost:9092)"
    else
        echo -e "  Kafka: ${RED}Not accessible${NC} (localhost:9092)"
    fi
    
    echo
}

# Function to start data ingestion service
start_data_ingestion() {
    print_info "Starting Data Ingestion Service..."
    
    if pgrep -f "data-ingestion" >/dev/null; then
        print_warning "Data Ingestion Service is already running"
        return
    fi
    
    cd "$PROJECT_ROOT/backend/services/data-ingestion"
    
    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER="${POSTGRES_USER:-scopeapi}"
    export DB_PASSWORD="${POSTGRES_PASSWORD}"
    export DB_NAME="${POSTGRES_DB:-scopeapi}"
    export KAFKA_BROKERS=localhost:9092
    export SERVER_PORT=$DATA_INGESTION_PORT
    
    # Start the service
    nohup ./data-ingestion > "$LOGS_DIR/data-ingestion.log" 2>&1 &
    DATA_INGESTION_PID=$!
    
    print_success "Data Ingestion Service started with PID: $DATA_INGESTION_PID"
}

# Function to start API discovery service
start_api_discovery() {
    print_info "Starting API Discovery Service..."
    
    if pgrep -f "api-discovery" >/dev/null; then
        print_warning "API Discovery Service is already running"
        return
    fi
    
    cd "$PROJECT_ROOT/backend/services/api-discovery"
    
    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER="${POSTGRES_USER:-scopeapi}"
    export DB_PASSWORD="${POSTGRES_PASSWORD}"
    export DB_NAME="${POSTGRES_DB:-scopeapi}"
    export KAFKA_BROKERS=localhost:9092
    export SERVER_PORT=$API_DISCOVERY_PORT
    
    # Start the service
    nohup go run cmd/main.go > "$LOGS_DIR/api-discovery.log" 2>&1 &
    API_DISCOVERY_PID=$!
    
    print_success "API Discovery Service started with PID: $API_DISCOVERY_PID"
}

# Function to start threat detection service
start_threat_detection() {
    print_info "Starting Threat Detection Service..."
    
    # Load environment variables if not already loaded
    # Check infrastructure dependencies
    check_infrastructure
    
    if [[ -z "$DB_PASSWORD" ]]; then
        load_env_vars
    fi
    
    if pgrep -f "threat-detection" >/dev/null; then
        print_warning "Threat Detection Service is already running"
        return
    fi
    
    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER="${POSTGRES_USER:-scopeapi}"
    export DB_PASSWORD="${POSTGRES_PASSWORD}"
    export DB_NAME="${POSTGRES_DB:-scopeapi}"
    export KAFKA_BROKERS=localhost:9092
    export SERVER_PORT=$THREAT_DETECTION_PORT
    
    # Start the service
    nohup go run cmd/main.go > "$LOGS_DIR/threat-detection.log" 2>&1 &
    THREAT_DETECTION_PID=$!
    
    print_success "Threat Detection Service started with PID: $THREAT_DETECTION_PID"
}

# Function to start gateway integration service
start_gateway_integration() {
    print_info "Starting Gateway Integration Service..."
    
    if pgrep -f "gateway-integration" >/dev/null; then
        print_warning "Gateway Integration Service is already running"
        return
    fi
    
    cd "$PROJECT_ROOT/backend/services/gateway-integration"
    
    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER="${POSTGRES_USER:-scopeapi}"
    export DB_PASSWORD="${POSTGRES_PASSWORD}"
    export DB_NAME="${POSTGRES_DB:-scopeapi}"
    export KAFKA_BROKERS=localhost:9092
    export SERVER_PORT=$GATEWAY_INTEGRATION_PORT
    
    # Start the service
    nohup go run cmd/main.go > "$LOGS_DIR/gateway-integration.log" 2>&1 &
    GATEWAY_INTEGRATION_PID=$!
    
    print_success "Gateway Integration Service started with PID: $GATEWAY_INTEGRATION_PID"
}

# Function to start attack blocking service
start_attack_blocking() {
    print_info "Starting Attack Blocking Service..."
    
    if pgrep -f "attack-blocking" >/dev/null; then
        print_warning "Attack Blocking Service is already running"
        return
    fi
    
    cd "$PROJECT_ROOT/backend/services/attack-blocking"
    
    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER="${POSTGRES_USER:-scopeapi}"
    export DB_PASSWORD="${POSTGRES_PASSWORD}"
    export DB_NAME="${POSTGRES_DB:-scopeapi}"
    export KAFKA_BROKERS=localhost:9092
    export SERVER_PORT=$ATTACK_BLOCKING_PORT
    
    # Start the service
    nohup go run cmd/main.go > "$LOGS_DIR/attack-blocking.log" 2>&1 &
    ATTACK_BLOCKING_PID=$!
    
    print_success "Attack Blocking Service started with PID: $ATTACK_BLOCKING_PID"
}

# Function to start data protection service
start_data_protection() {
    print_info "Starting Data Protection Service..."
    
    if pgrep -f "data-protection" >/dev/null; then
        print_warning "Data Protection Service is already running"
        return
    fi
    
    cd "$PROJECT_ROOT/backend/services/data-protection"
    
    # Set environment variables
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER="${POSTGRES_USER:-scopeapi}"
    export DB_PASSWORD="${POSTGRES_PASSWORD}"
    export DB_NAME="${POSTGRES_DB:-scopeapi}"
    export KAFKA_BROKERS=localhost:9092
    export SERVER_PORT=$DATA_PROTECTION_PORT
    
    # Start the service
    nohup go run cmd/main.go > "$LOGS_DIR/data-protection.log" 2>&1 &
    DATA_PROTECTION_PID=$!
    
    print_success "Data Protection Service started with PID: $DATA_PROTECTION_PID"
}

# Function to build gateway integration service
build_gateway_integration() {
    print_info "Building Gateway Integration Service..."
    
    cd "$PROJECT_ROOT/backend/services/gateway-integration"
    
    # Check if Makefile exists and use it
    if [[ -f "Makefile" ]]; then
        print_info "Using Makefile to build service..."
        make build
    else
        print_info "Building with go build..."
        go build -o gateway-integration ./cmd/main.go
    fi
    
    print_success "Gateway Integration Service built successfully!"
}

# Function to build all services
build_all_services() {
    print_header "Building ScopeAPI Services"
    
    print_info "Building services..."
    
    # Build Gateway Integration Service (since it's complete)
    build_gateway_integration
    
    print_success "All services built successfully!"
}

# Function to start admin console
start_admin_console() {
    print_info "Starting Admin Console..."
    
    if pgrep -f "ng serve" >/dev/null; then
        print_warning "Admin Console is already running"
        return
    fi
    
    cd "$PROJECT_ROOT/adminConsole"
    
    # Check if node_modules exists
    if [[ ! -d "node_modules" ]]; then
        print_info "Installing dependencies..."
        npm install
    fi
    
    # Start the Angular development server
    nohup npm start > "$LOGS_DIR/admin-console.log" 2>&1 &
    ADMIN_CONSOLE_PID=$!
    
    print_success "Admin Console started with PID: $ADMIN_CONSOLE_PID"
}

# Function to start all services
start_all_services() {
    print_header "Starting ScopeAPI Services"
    
    # Load environment variables
    # Check infrastructure dependencies
    check_infrastructure
    
    load_env_vars
    
    # Create logs directory if it doesn't exist
    mkdir -p "$LOGS_DIR"
    
    print_info "Starting services in order..."
    
    # Start services
    start_data_ingestion
    start_api_discovery
    start_threat_detection
    start_gateway_integration
    start_attack_blocking
    start_data_protection
    start_admin_console
    
    # Wait a moment for services to start
    sleep 3
    
    # Write PIDs to file
    print_info "Writing PIDs to $PID_FILE..."
    echo "DATA_INGESTION_PID=$DATA_INGESTION_PID" > "$PID_FILE"
    echo "API_DISCOVERY_PID=$API_DISCOVERY_PID" >> "$PID_FILE"
    echo "THREAT_DETECTION_PID=$THREAT_DETECTION_PID" >> "$PID_FILE"
    echo "GATEWAY_INTEGRATION_PID=$GATEWAY_INTEGRATION_PID" >> "$PID_FILE"
    echo "ATTACK_BLOCKING_PID=$ATTACK_BLOCKING_PID" >> "$PID_FILE"
    echo "DATA_PROTECTION_PID=$DATA_PROTECTION_PID" >> "$PID_FILE"
    echo "ADMIN_CONSOLE_PID=$ADMIN_CONSOLE_PID" >> "$PID_FILE"
    
    print_success "All services started successfully!"
    print_info "PID file created at: $PID_FILE"
    
    # Show status
    echo
    show_status
}

# Function to stop all services
stop_all_services() {
    print_header "Stopping ScopeAPI Services"
    
    if [[ ! -f "$PID_FILE" ]]; then
        print_error "PID file not found: $PID_FILE"
        print_info "Attempting to stop services by process name..."
        
        # Stop by process name if PID file doesn't exist
        pkill -f "data-ingestion" 2>/dev/null || true
        pkill -f "api-discovery" 2>/dev/null || true
        pkill -f "threat-detection" 2>/dev/null || true
        pkill -f "gateway-integration" 2>/dev/null || true
        pkill -f "attack-blocking" 2>/dev/null || true
        pkill -f "data-protection" 2>/dev/null || true
        pkill -f "ng serve" 2>/dev/null || true
        
        print_success "Services stopped by process name"
        return
    fi
    
    print_info "Reading PIDs from $PID_FILE..."
    
    # Read and stop each service
    while IFS= read -r line; do
        if [[ -n "$line" ]]; then
            local service=$(echo "$line" | cut -d'=' -f1)
            local pid=$(echo "$line" | cut -d'=' -f2)
            
            if [[ -n "$pid" ]]; then
                print_info "Stopping $service (PID: $pid)..."
                
                # Try graceful shutdown first
                if kill -TERM "$pid" 2>/dev/null; then
                    sleep 2
                    # Force kill if still running
                    if kill -0 "$pid" 2>/dev/null; then
                        print_warning "Force killing $service..."
                        kill -KILL "$pid" 2>/dev/null || true
                    fi
                    print_success "$service stopped"
                else
                    print_warning "$service (PID: $pid) not found"
                fi
            fi
        fi
    done < "$PID_FILE"
    
    # Remove PID file
    rm -f "$PID_FILE"
    print_success "All services stopped and PID file removed"
}

# Function to show status
show_status() {
    print_header "ScopeAPI Services Status"
    
    # Check PID file
    check_pid_file
    
    # Check individual services
    check_service_status "Data Ingestion Service" "data-ingestion" "$DATA_INGESTION_PORT" "http://localhost:$DATA_INGESTION_PORT/health"
    check_service_status "API Discovery Service" "api-discovery" "$API_DISCOVERY_PORT" "http://localhost:$API_DISCOVERY_PORT/health"
    check_service_status "Threat Detection Service" "threat-detection" "$THREAT_DETECTION_PORT" "http://localhost:$THREAT_DETECTION_PORT/health"
    check_service_status "Gateway Integration Service" "gateway-integration" "$GATEWAY_INTEGRATION_PORT" "http://localhost:$GATEWAY_INTEGRATION_PORT/health"
    check_service_status "Attack Blocking Service" "attack-blocking" "$ATTACK_BLOCKING_PORT" "http://localhost:$ATTACK_BLOCKING_PORT/health"
    check_service_status "Data Protection Service" "data-protection" "$DATA_PROTECTION_PORT" "http://localhost:$DATA_PROTECTION_PORT/health"
    check_service_status "Admin Console" "ng serve" "$ADMIN_CONSOLE_PORT" "http://localhost:$ADMIN_CONSOLE_PORT"
    
    # Check system resources
    check_system_resources
    
    # Check dependencies
    check_dependencies
}

# Function to show help
show_help() {
    echo "ScopeAPI Local Development Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
echo "  start           Start all ScopeAPI services"
echo "  stop            Stop all ScopeAPI services"
echo "  status          Show status of all services"
echo "  build           Build all ScopeAPI services"
echo "  help            Show this help message"
    echo ""
    echo "Examples:"
echo "  $0 start                    # Start all services"
echo "  $0 stop                     # Stop all services"
echo "  $0 status                   # Show service status"
echo "  $0 build                    # Build all services"
echo "  $0 help                     # Show this help"
    echo ""
    echo "Services managed:"
echo "  - Data Ingestion Service (port $DATA_INGESTION_PORT)"
echo "  - API Discovery Service (port $API_DISCOVERY_PORT)"
echo "  - Threat Detection Service (port $THREAT_DETECTION_PORT)"
echo "  - Gateway Integration Service (port $GATEWAY_INTEGRATION_PORT)"
echo "  - Attack Blocking Service (port $ATTACK_BLOCKING_PORT)"
echo "  - Data Protection Service (port $DATA_PROTECTION_PORT)"
echo "  - Admin Console (port $ADMIN_CONSOLE_PORT)"
    echo ""
    echo "Files:"
    echo ""
    echo "Dependencies:"
    echo "  - Infrastructure must be running (use: ./docker-infrastructure.sh start)"
    echo "  - PostgreSQL, Kafka, Redis, Elasticsearch, Kibana"
    echo "  PID File: $PID_FILE"
    echo "  Logs: $LOGS_DIR/"
}

# Main function
main() {
    case "${1:-}" in
        start)
            start_all_services
            ;;
        stop)
            stop_all_services
            ;;
        status)
            show_status
            ;;
        build)
            build_all_services
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@" 