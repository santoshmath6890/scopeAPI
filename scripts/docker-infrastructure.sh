#!/bin/bash

# ScopeAPI Infrastructure Management Script (Infrastructure Only)
# This script manages ONLY infrastructure services (PostgreSQL, Kafka, Redis, etc.)
# Note: For microservices management, use scopeapi-services.sh instead
# This script is focused on infrastructure setup and troubleshooting


set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project configuration
PROJECT_NAME="scopeapi"
NETWORK_NAME="${PROJECT_NAME}-network"

# Load environment variables from .env file if it exists
load_env_vars() {
    if [ -f ".env" ]; then
        print_info "Loading environment variables from .env file"
        # Use source to properly handle variables with spaces
        set -a
        source .env
        set +a
    else
        print_warning ".env file not found. Please create one based on env.example"
        print_info "Some services may fail to start without required environment variables"
    fi
}

# Service configurations
ZOOKEEPER_PORT=2181
KAFKA_PORT=9092
POSTGRES_PORT=5432
POSTGRES_DB="${POSTGRES_DB:-scopeapi}"
POSTGRES_USER="${POSTGRES_USER:-scopeapi}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD}"
REDIS_PORT=6379
REDIS_PASSWORD="${REDIS_PASSWORD}"
ELASTICSEARCH_PORT=9200
KIBANA_PORT=5601

# Function to print colored output
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

print_header() {
    echo -e "${BLUE}==========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}==========================================${NC}"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    print_success "Docker is running"
}

# Function to create Docker network
create_network() {
    if ! docker network ls | grep -q "$NETWORK_NAME"; then
        print_info "Creating Docker network: $NETWORK_NAME"
        docker network create "$NETWORK_NAME"
        print_success "Network created: $NETWORK_NAME"
    else
        print_info "Network already exists: $NETWORK_NAME"
    fi
}

# Function to start ZooKeeper
start_zookeeper() {
    print_info "Starting ZooKeeper..."
    
    if docker ps | grep -q "scopeapi-zookeeper"; then
        print_warning "ZooKeeper is already running"
        return
    fi
    
    docker run -d \
        --name scopeapi-zookeeper \
        --network "$NETWORK_NAME" \
        -p "$ZOOKEEPER_PORT:$ZOOKEEPER_PORT" \
        -e ZOOKEEPER_CLIENT_PORT="$ZOOKEEPER_PORT" \
        -e ZOOKEEPER_TICK_TIME=2000 \
        -e ZOOKEEPER_INIT_LIMIT=5 \
        -e ZOOKEEPER_SYNC_LIMIT=2 \
        confluentinc/cp-zookeeper:7.4.0
    
    print_success "ZooKeeper started on port $ZOOKEEPER_PORT"
}

# Function to start Kafka
start_kafka() {
    print_info "Starting Kafka..."
    
    if docker ps | grep -q "scopeapi-kafka"; then
        print_warning "Kafka is already running"
        return
    fi
    
    # Wait for ZooKeeper to be ready
    print_info "Waiting for ZooKeeper to be ready..."
    sleep 10
    
    docker run -d \
        --name scopeapi-kafka \
        --network "$NETWORK_NAME" \
        -p "$KAFKA_PORT:$KAFKA_PORT" \
        -e KAFKA_BROKER_ID=1 \
        -e KAFKA_ZOOKEEPER_CONNECT="scopeapi-zookeeper:$ZOOKEEPER_PORT" \
        -e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://localhost:$KAFKA_PORT" \
        -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT \
        -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
        -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
        -e KAFKA_AUTO_CREATE_TOPICS_ENABLE=true \
        -e KAFKA_DELETE_TOPIC_ENABLE=true \
        confluentinc/cp-kafka:7.4.0
    
    print_success "Kafka started on port $KAFKA_PORT"
}

# Function to start PostgreSQL
start_postgres() {
    print_info "Starting PostgreSQL..."
    
    if [ -z "$POSTGRES_PASSWORD" ]; then
        print_error "POSTGRES_PASSWORD environment variable is not set. Please set it in your .env file."
        exit 1
    fi
    
    # Check if PostgreSQL is already running on the system
    if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
        print_warning "PostgreSQL is already running on the system (port 5432)"
        print_info "Skipping Docker PostgreSQL container"
        return
    fi
    
    if docker ps | grep -q "scopeapi-postgres"; then
        print_warning "PostgreSQL is already running in Docker"
        return
    fi
    
    docker run -d \
        --name scopeapi-postgres \
        --network "$NETWORK_NAME" \
        -p "$POSTGRES_PORT:$POSTGRES_PORT" \
        -e POSTGRES_DB="$POSTGRES_DB" \
        -e POSTGRES_USER="$POSTGRES_USER" \
        -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
        -e POSTGRES_INITDB_ARGS="--encoding=UTF-8 --lc-collate=C --lc-ctype=C" \
        -v scopeapi-postgres-data:/var/lib/postgresql/data \
        postgres:15
    
    print_success "PostgreSQL started on port $POSTGRES_PORT"
    print_info "Database: $POSTGRES_DB, User: $POSTGRES_USER, Password: [HIDDEN]"
}

# Function to start Redis
start_redis() {
    print_info "Starting Redis..."
    
    if docker ps | grep -q "scopeapi-redis"; then
        print_warning "Redis is already running"
        return
    fi
    
    if [ -z "$REDIS_PASSWORD" ]; then
        print_error "REDIS_PASSWORD environment variable is not set. Please set it in your .env file."
        exit 1
    fi
    
    docker run -d \
        --name scopeapi-redis \
        --network "$NETWORK_NAME" \
        -p "$REDIS_PORT:$REDIS_PORT" \
        -e REDIS_PASSWORD="$REDIS_PASSWORD" \
        redis:7-alpine redis-server --requirepass "$REDIS_PASSWORD"
    
    print_success "Redis started on port $REDIS_PORT"
    print_info "Password: [HIDDEN]"
}

# Function to start Elasticsearch
start_elasticsearch() {
    print_info "Starting Elasticsearch..."
    
    if docker ps | grep -q "scopeapi-elasticsearch"; then
        print_warning "Elasticsearch is already running"
        return
    fi
    
    docker run -d \
        --name scopeapi-elasticsearch \
        --network "$NETWORK_NAME" \
        -p "$ELASTICSEARCH_PORT:$ELASTICSEARCH_PORT" \
        -e "discovery.type=single-node" \
        -e "xpack.security.enabled=false" \
        -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
        -v scopeapi-elasticsearch-data:/usr/share/elasticsearch/data \
        elasticsearch:8.11.0
    
    print_success "Elasticsearch started on port $ELASTICSEARCH_PORT"
}

# Function to start Kibana
start_kibana() {
    print_info "Starting Kibana..."
    
    if docker ps | grep -q "scopeapi-kibana"; then
        print_warning "Kibana is already running"
        return
    fi
    
    # Wait for Elasticsearch to be ready
    print_info "Waiting for Elasticsearch to be ready..."
    sleep 15
    
    docker run -d \
        --name scopeapi-kibana \
        --network "$NETWORK_NAME" \
        -p "$KIBANA_PORT:$KIBANA_PORT" \
        -e "ELASTICSEARCH_HOSTS=http://scopeapi-elasticsearch:$ELASTICSEARCH_PORT" \
        -e "ELASTICSEARCH_URL=http://scopeapi-elasticsearch:$ELASTICSEARCH_PORT" \
        kibana:8.11.0
    
    print_success "Kibana started on port $KIBANA_PORT"
}

# Function to start all services
start_all() {
    print_header "Starting ScopeAPI Infrastructure Services"
    
    check_docker
    load_env_vars
    create_network
    
    print_info "Starting services in order..."
    
    start_zookeeper
    start_kafka
    start_postgres
    start_redis
    start_elasticsearch
    start_kibana
    
    print_success "All infrastructure services started!"
    print_info "Waiting for services to be ready..."
    sleep 10
    
    show_status
}

# Function to stop all services
stop_all() {
    print_header "Stopping ScopeAPI Infrastructure Services"
    
    local containers=(
        "scopeapi-kibana"
        "scopeapi-elasticsearch"
        "scopeapi-redis"
        "scopeapi-postgres"
        "scopeapi-kafka"
        "scopeapi-zookeeper"
    )
    
    for container in "${containers[@]}"; do
        if docker ps | grep -q "$container"; then
            print_info "Stopping $container..."
            docker stop "$container"
            print_success "$container stopped"
        else
            print_info "$container is not running"
        fi
    done
    
    print_success "All infrastructure services stopped!"
}

# Function to restart all services
restart_all() {
    print_header "Restarting ScopeAPI Infrastructure Services"
    
    stop_all
    sleep 5
    start_all
}

# Function to check system PostgreSQL status
check_system_postgres() {
    if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
        print_success "System PostgreSQL: Running on port 5432"
        return 0
    else
        print_warning "System PostgreSQL: Not accessible on port 5432"
        return 1
    fi
}

# Function to show status of all services
show_status() {
    print_header "ScopeAPI Infrastructure Services Status"
    
    echo ""
    print_info "Container Status:"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep scopeapi || echo "No ScopeAPI containers running"
    
    echo ""
    print_info "System Services:"
    check_system_postgres
    
    echo ""
    print_info "Network Information:"
    if docker network ls | grep -q "$NETWORK_NAME"; then
        docker network inspect "$NETWORK_NAME" --format "{{.Name}}: {{.IPAM.Config}}" | head -1
    else
        print_warning "Network $NETWORK_NAME not found"
    fi
    
    echo ""
    print_info "Service Endpoints:"
    echo "  ZooKeeper:     localhost:$ZOOKEEPER_PORT"
    echo "  Kafka:         localhost:$KAFKA_PORT"
    echo "  PostgreSQL:    localhost:$POSTGRES_PORT (System)"
    echo "  Redis:         localhost:$REDIS_PORT"
    echo "  Elasticsearch: localhost:$ELASTICSEARCH_PORT"
    echo "  Kibana:        localhost:$KIBANA_PORT"
    
    echo ""
    print_info "Connection Details:"
    echo "  PostgreSQL: psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB"
    echo "  Redis:      redis-cli -h localhost -p $REDIS_PORT -a [PASSWORD]"
    echo "  Kafka:      kafka-console-consumer --bootstrap-server localhost:$KAFKA_PORT --topic test"
}

# Function to clean up (remove containers and volumes)
cleanup() {
    print_header "Cleaning Up ScopeAPI Infrastructure"
    
    print_warning "This will remove ALL containers and volumes. Are you sure? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_info "Removing containers..."
        docker rm -f scopeapi-kibana scopeapi-elasticsearch scopeapi-redis scopeapi-postgres scopeapi-kafka scopeapi-zookeeper 2>/dev/null || true
        
        print_info "Removing volumes..."
        docker volume rm scopeapi-postgres-data scopeapi-elasticsearch-data 2>/dev/null || true
        
        print_info "Removing network..."
        docker network rm "$NETWORK_NAME" 2>/dev/null || true
        
        print_success "Cleanup completed!"
    else
        print_info "Cleanup cancelled"
    fi
}

# Function to show logs
show_logs() {
    local service="$1"
    
    if [[ -z "$service" ]]; then
        print_error "Please specify a service name"
        echo "Available services: zookeeper, kafka, postgres, redis, elasticsearch, kibana"
        exit 1
    fi
    
    local container_name="scopeapi-$service"
    
    if docker ps | grep -q "$container_name"; then
        print_info "Showing logs for $container_name..."
        docker logs -f "$container_name"
    else
        print_error "Container $container_name is not running"
        exit 1
    fi
}

# Function to fix Docker permissions
fix_docker_permissions() {
    print_header "Fixing Docker Permissions"
    
    # Check if user is root
    if [ "$EUID" -eq 0 ]; then
        print_error "This script should not be run as root"
        print_warning "Please run as a regular user"
        exit 1
    fi
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        print_warning "Please install Docker first: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    # Check if user is already in docker group
    if groups $USER | grep -q docker; then
        print_success "User is already in docker group"
        print_warning "You may need to log out and log back in, or run: newgrp docker"
    else
        print_warning "Adding user to docker group..."
        
        # Add user to docker group
        if sudo usermod -aG docker $USER; then
            print_success "User added to docker group"
            print_warning "You need to log out and log back in for changes to take effect"
            print_warning "Or run: newgrp docker"
        else
            print_error "Failed to add user to docker group"
            exit 1
        fi
    fi
    
    # Test Docker access
    echo ""
    print_info "Testing Docker access..."
    if docker ps &> /dev/null; then
        print_success "Docker access is working correctly"
    else
        print_warning "Docker access still not working"
        print_warning "Try logging out and logging back in, or run: newgrp docker"
        print_warning "If the issue persists, restart your system"
    fi
    
    print_success "Docker permissions fix completed!"
}

# Function to setup environment file
setup_env() {
    print_header "Setting up Environment Configuration"
    
    if [ -f ".env" ]; then
        print_warning ".env file already exists"
        read -p "Do you want to overwrite it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Environment setup cancelled"
            return
        fi
    fi
    
    if [ -f "env.example" ]; then
        cp env.example .env
        print_success "Created .env file from env.example"
        print_info "Please edit .env file and set your secure passwords"
        print_info "Then run: $0 start"
    else
        print_error "env.example file not found"
        print_info "Please create a .env file manually with the required environment variables"
    fi
}

# Function to show help
show_help() {
    echo "ScopeAPI Infrastructure Management Script"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  start           Start all infrastructure services"
    echo "  stop            Stop all infrastructure services"
    echo "  restart         Restart all infrastructure services"
    echo "  status          Show status of all services"
    echo "  logs <service>  Show logs for a specific service"
    echo "  cleanup         Remove all containers and volumes"
    echo "  fix-permissions Fix Docker permission issues"
    echo "  setup-env       Set up environment file from template"
    echo "  help            Show this help message"
    echo ""
    echo "Services:"
    echo "  - ZooKeeper (port 2181)"
    echo "  - Kafka (port 9092)"
    echo "  - PostgreSQL (port 5432)"
    echo "  - Redis (port 6379)"
    echo "  - Elasticsearch (port 9200)"
    echo "  - Kibana (port 5601)"
    echo ""
    echo "Examples:"
    echo "  $0 start                    # Start all services"
    echo "  $0 logs kafka              # Show Kafka logs"
    echo "  $0 status                  # Show service status"
    echo "  $0 fix-permissions         # Fix Docker permissions"
    echo "  $0 setup-env               # Set up environment file"
    echo "  $0 cleanup                 # Remove everything"
}

# Main script logic
main() {
    case "${1:-}" in
        start)
            start_all
            ;;
        stop)
            stop_all
            ;;
        restart)
            restart_all
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs "$2"
            ;;
        cleanup)
            cleanup
            ;;
        fix-permissions)
            fix_docker_permissions
            ;;
        setup-env)
            setup_env
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