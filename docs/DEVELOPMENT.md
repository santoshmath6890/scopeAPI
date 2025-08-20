# 💻 ScopeAPI Development Guide

This guide covers everything you need to know to develop, test, and contribute to the ScopeAPI project.

## 📋 **Table of Contents**

- [Development Environment](#development-environment)
- [Project Setup](#project-setup)
- [Development Workflows](#development-workflows)
- [Testing Strategies](#testing-strategies)
- [Code Standards](#code-standards)
- [Debugging](#debugging)
- [Performance Optimization](#performance-optimization)
- [Troubleshooting](#troubleshooting)

## 🚀 **Development Environment**

### **Prerequisites**
- **Docker** 24.0+ with Docker Compose
- **Go** 1.21+ for backend development
- **Node.js** 18+ for frontend development
- **Git** for version control
- **IDE/Editor** (VS Code, GoLand, etc.)

### **Recommended Tools**
- **VS Code Extensions**:
  - Go extension
  - Docker extension
  - Angular Language Service
  - GitLens
- **Go Tools**:
  - `golangci-lint` for linting
  - `goimports` for code formatting
  - `delve` for debugging
- **Frontend Tools**:
  - Angular CLI
  - ESLint and Prettier
  - Cypress for E2E testing

## 🏗️ **Project Setup**

### **1. Clone and Setup**
```bash
# Clone the repository
git clone https://github.com/your-org/scopeapi.git
cd scopeapi

# Complete setup with validation
cd scripts
./scopeapi.sh setup --full

# Verify setup
./scopeapi.sh setup --validate
```

### **2. Environment Configuration**

**⚠️ IMPORTANT: For production, use Kubernetes Secrets instead of .env files!**

**Development/Staging**: Create a `.env.local` file for local development:
```bash
# Copy environment template
cp env.example .env.local

# Edit with your configuration
nano .env.local

# Key environment variables:
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_secure_password
KAFKA_BROKER_ID=1
```

**Production**: Use Kubernetes Secrets (see `k8s/secrets.yaml`):
```bash
# Generate base64 encoded secrets
./scripts/generate-secrets.sh

# Deploy to Kubernetes
./scripts/deploy.sh -e staging -p k8s
```

### **3. Database Setup**
```bash
# Setup database with test data
./scripts/setup-database.sh --test-data

# Verify database
./scripts/setup-database.sh --validate
```

## 🐳 **Docker Setup**

### **Three-File Architecture**

This project uses a **clean three-file approach** for Docker Compose:

#### **1. `scripts/docker-compose.yml` (Main)**
- **Purpose**: Production-ready orchestration
- **Content**: Infrastructure + microservices with basic settings
- **Usage**: `docker-compose -f scripts/docker-compose.yml up` (production-like)

#### **2. `scripts/docker-compose.override.yml` (Development)**
- **Purpose**: Development environment enhancements
- **Content**: Service ports, development environment variables, source mounting
- **Usage**: Automatically loaded for development

#### **3. `scripts/docker-compose.debug.yml` (Debug)**
- **Purpose**: Debugging capabilities only
- **Content**: Debug ports, Delve debugger, debug Dockerfiles
- **Usage**: `docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.debug.yml up`

### **Docker Workflows**

#### **Complete Setup (Recommended for First Time)**
```bash
# Complete setup with validation
cd scripts
./scopeapi.sh setup --full

# Or step by step:
./scopeapi.sh setup --infrastructure  # Start infrastructure
./scopeapi.sh setup --database        # Setup database
./scopeapi.sh setup --validate        # Validate setup
```

#### **Development Mode (After Setup)**
```bash
# Start infrastructure + services with development overrides
cd scripts
./dev.sh start api-discovery

# Or use docker-compose directly (automatically loads override)
docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.override.yml up api-discovery
```

#### **Debug Mode**
```bash
# Start services in debug mode
./dev.sh debug api-discovery

# Or use docker-compose with debug config
docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.debug.yml up api-discovery
```

#### **Production Mode**
```bash
# Start services without development overrides
docker-compose -f scripts/docker-compose.yml up api-discovery
```

### **Docker Compose File Structure**

#### **`scripts/docker-compose.yml` (Main Configuration)**
```yaml
# Infrastructure services
zookeeper:
  image: confluentinc/cp-zookeeper:7.4.0
  environment:
    ZOOKEEPER_CLIENT_PORT: 2181
    ZOOKEEPER_TICK_TIME: 2000

kafka:
  image: confluentinc/cp-kafka:7.4.0
  depends_on:
    - zookeeper
  environment:
    KAFKA_BROKER_ID: 1
    KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:9092
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
    KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1

postgres:
  image: postgres:15
  environment:
    POSTGRES_DB: scopeapi
    POSTGRES_USER: scopeapi
    POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
  volumes:
    - postgres_data:/var/lib/postgresql/data

redis:
  image: redis:7-alpine
  command: redis-server --requirepass ${REDIS_PASSWORD}
  volumes:
    - redis_data:/data

# Microservices
api-discovery:
  build:
    context: ../backend/services/api-discovery
    dockerfile: Dockerfile
  depends_on:
    - postgres
    - redis
    - kafka
  environment:
    - POSTGRES_HOST=postgres
    - REDIS_HOST=redis
    - KAFKA_BROKERS=kafka:9092
  env_file:
    - .env.local

gateway-integration:
  build:
    context: ../backend/services/gateway-integration
    dockerfile: Dockerfile
  depends_on:
    - postgres
    - redis
    - kafka
  environment:
    - POSTGRES_HOST=postgres
    - REDIS_HOST=redis
    - KAFKA_BROKERS=kafka:9092
  env_file:
    - .env.local

# Volumes
volumes:
  postgres_data:
  redis_data:
```

#### **`scripts/docker-compose.override.yml` (Development Overrides)**
```yaml
# Development-specific overrides
services:
  api-discovery:
    ports:
      - "8080:8080"
    volumes:
      - ../backend/services/api-discovery:/app
      - /app/vendor
    environment:
      - GO_ENV=development
      - DEBUG=true

  gateway-integration:
    ports:
      - "8081:8081"
    volumes:
      - ../backend/services/api-discovery:/app
      - /app/vendor
    environment:
      - GO_ENV=development
      - DEBUG=true

  postgres:
    ports:
      - "5432:5432"

  redis:
    ports:
      - "6379:6379"

  kafka:
    ports:
      - "9092:9092"
```

#### **`scripts/docker-compose.debug.yml` (Debug Configuration)**
```yaml
# Debug-specific overrides
services:
  api-discovery:
    build:
      context: ../backend/services/api-discovery
      dockerfile: Dockerfile.debug
    ports:
      - "2345:2345"  # Delve debugger
    environment:
      - DEBUG=true
      - DELVE=true
    command: ["dlv", "exec", "--listen=:2345", "--headless=true", "--continue", "--accept-multiclient", "/app/api-discovery"]

  gateway-integration:
    build:
      context: ../backend/services/gateway-integration
      dockerfile: Dockerfile.debug
    ports:
      - "2346:2345"  # Delve debugger
    environment:
      - DEBUG=true
      - DELVE=true
    command: ["dlv", "exec", "--listen=:2345", "--headless=true", "--continue", "--accept-multiclient", "/app/gateway-integration"]
```

### **Docker Compose File Structure**

#### **`scripts/docker-compose.yml` (Main Configuration)**
```yaml
# Infrastructure services
zookeeper:
  image: confluentinc/cp-zookeeper:7.4.0
  environment:
    ZOOKEEPER_CLIENT_PORT: 2181
    ZOOKEEPER_TICK_TIME: 2000

kafka:
  image: confluentinc/cp-kafka:7.4.0
  depends_on:
    - zookeeper
  environment:
    KAFKA_BROKER_ID: 1
    KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:9092
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
    KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1

postgres:
  image: postgres:15
  environment:
    POSTGRES_DB: scopeapi
    POSTGRES_USER: scopeapi
    POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
  volumes:
    - postgres_data:/var/lib/postgresql/data

redis:
  image: redis:7-alpine
  command: redis-server --requirepass ${REDIS_PASSWORD}
  volumes:
    - redis_data:/data

# Microservices
api-discovery:
  build:
    context: ../backend/services/api-discovery
    dockerfile: Dockerfile
  depends_on:
    - postgres
    - redis
    - kafka
  environment:
    - POSTGRES_HOST=postgres
    - REDIS_HOST=redis
    - KAFKA_BROKERS=kafka:9092
  env_file:
    - .env.local

gateway-integration:
  build:
    context: ../backend/services/gateway-integration
    dockerfile: Dockerfile
  depends_on:
    - postgres
    - redis
    - kafka
  environment:
    - POSTGRES_HOST=postgres
    - REDIS_HOST=redis
    - KAFKA_BROKERS=kafka:9092
  env_file:
    - .env.local

# Volumes
volumes:
  postgres_data:
  redis_data:
```

#### **`scripts/docker-compose.override.yml` (Development Overrides)**
```yaml
# Development-specific overrides
services:
  api-discovery:
    ports:
      - "8080:8080"
    volumes:
      - ../backend/services/api-discovery:/app
      - /app/vendor
    environment:
      - GO_ENV=development
      - DEBUG=true

  gateway-integration:
    ports:
      - "8081:8081"
    volumes:
      - ../backend/services/gateway-integration:/app
      - /app/vendor
    environment:
      - GO_ENV=development
      - DEBUG=true

  postgres:
    ports:
      - "5432:5432"

  redis:
    ports:
      - "6379:6379"

  kafka:
    ports:
      - "9092:9092"
```

#### **`scripts/docker-compose.debug.yml` (Debug Configuration)**
```yaml
# Debug-specific overrides
services:
  api-discovery:
    build:
      context: ../backend/services/api-discovery
      dockerfile: Dockerfile.debug
    ports:
      - "2345:2345"  # Delve debugger
    environment:
      - DEBUG=true
      - DELVE=true
    command: ["dlv", "exec", "--listen=:2345", "--headless=true", "--continue", "--accept-multiclient", "/app/api-discovery"]

  gateway-integration:
    build:
      context: ../backend/services/gateway-integration
      dockerfile: Dockerfile.debug
    ports:
      - "2346:2345"  # Delve debugger
    environment:
      - DEBUG=true
      - DELVE=true
    command: ["dlv", "exec", "--listen=:2345", "--headless=true", "--continue", "--accept-multiclient", "/app/gateway-integration"]
```

## 🔧 **Development Workflows**

### **Daily Development Workflow**
```bash
# 1. Start infrastructure
cd scripts
./scopeapi.sh setup --infrastructure

# 2. Start services for development
./dev.sh start all

# 3. Make code changes
# 4. View logs if needed
./dev.sh logs api-discovery

# 5. Test changes
# 6. Stop when done
./dev.sh stop
```

### **Service-Specific Development**
```bash
# Start only specific service
./dev.sh start api-discovery

# Start multiple services
./dev.sh start api-discovery gateway-integration

# View specific service logs
./dev.sh logs api-discovery

# Open shell in service container
./dev.sh shell api-discovery
```

### **Debugging Workflow**
```bash
# 1. Start service in debug mode
cd scripts
./dev.sh debug api-discovery

# 2. Connect IDE to debug port (2345 for api-discovery)
# 3. Set breakpoints in your Go code
# 4. Debug and step through code
# 5. Stop debug session
./dev.sh stop
```

## 🧪 **Testing Strategies**

### **Backend Testing (Go)**
```bash
# Run all tests
cd backend
go test ./...

# Run specific service tests
go test ./services/api-discovery/...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

### **Frontend Testing (Angular)**
```bash
# Unit tests
cd adminConsole
ng test

# E2E tests
ng e2e

# Build for production
ng build --prod
```

### **Integration Testing**
```bash
# Test database integration
./scripts/setup-database.sh --validate

# Test service communication
cd scripts
./dev.sh start all
# Verify all services are healthy
./scopeapi.sh status
```

### **Performance Testing**
```bash
# Load testing with Apache Bench
ab -n 1000 -c 10 http://localhost:8080/health

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

## 📝 **Code Standards**

### **Go Code Standards**
- **Formatting**: Use `go fmt` and `goimports`
- **Linting**: Use `golangci-lint`
- **Documentation**: Follow Go documentation conventions
- **Testing**: Aim for 80%+ test coverage
- **Error Handling**: Use proper error wrapping and context

```go
// Good: Proper error handling
func processData(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("data cannot be empty")
    }
    
    result, err := validateData(data)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    return nil
}

// Good: Comprehensive testing
func TestProcessData(t *testing.T) {
    tests := []struct {
        name    string
        data    []byte
        wantErr bool
    }{
        {"empty data", []byte{}, true},
        {"valid data", []byte("test"), false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := processData(tt.data)
            if (err != nil) != tt.wantErr {
                t.Errorf("processData() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### **Angular Code Standards**
- **TypeScript**: Use strict mode and proper typing
- **Components**: Follow Angular style guide
- **Services**: Use dependency injection properly
- **Testing**: Unit test all components and services
- **Styling**: Use SCSS with BEM methodology

```typescript
// Good: Proper typing and interfaces
interface ApiEndpoint {
  id: string;
  url: string;
  method: HttpMethod;
  parameters?: Record<string, any>;
}

@Component({
  selector: 'app-api-endpoint',
  templateUrl: './api-endpoint.component.html',
  styleUrls: ['./api-endpoint.component.scss']
})
export class ApiEndpointComponent implements OnInit {
  @Input() endpoint!: ApiEndpoint;
  @Output() endpointSelected = new EventEmitter<ApiEndpoint>();

  constructor(private endpointService: EndpointService) {}

  ngOnInit(): void {
    // Component initialization
  }

  onSelect(): void {
    this.endpointSelected.emit(this.endpoint);
  }
}
```

## 🐛 **Debugging**

### **Go Debugging with Delve**
```bash
# Start service in debug mode
cd scripts
./dev.sh debug api-discovery

# VS Code launch.json configuration
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Attach to API Discovery",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app/src",
      "port": 2345,
      "host": "127.0.0.1",
      "showLog": true
    }
  ]
}
```

### **Frontend Debugging**
```bash
# Start admin console in development mode
cd adminConsole
ng serve

# Open browser dev tools
# Use Angular DevTools extension
# Monitor network requests and console logs
```

### **Database Debugging**
```bash
# Connect to PostgreSQL
docker exec -it scopeapi-postgres psql -U scopeapi -d scopeapi

# View logs
docker logs scopeapi-postgres

# Check service health
cd scripts
./scopeapi.sh status
```

## ⚡ **Performance Optimization**

### **Backend Optimization**
- **Database**: Use proper indexes and query optimization
- **Caching**: Implement Redis caching for frequently accessed data
- **Connection pooling**: Use connection pools for database connections
- **Goroutines**: Use goroutines for concurrent operations
- **Memory**: Profile and optimize memory usage

```go
// Good: Efficient database queries
func getEndpoints(ctx context.Context, limit int) ([]Endpoint, error) {
    query := `
        SELECT id, url, method, created_at 
        FROM api_discovery.endpoints 
        WHERE is_active = true 
        ORDER BY created_at DESC 
        LIMIT $1
    `
    
    rows, err := db.QueryContext(ctx, query, limit)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    var endpoints []Endpoint
    for rows.Next() {
        var endpoint Endpoint
        if err := rows.Scan(&endpoint.ID, &endpoint.URL, &endpoint.Method, &endpoint.CreatedAt); err != nil {
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        endpoints = append(endpoints, endpoint)
    }
    
    return endpoints, nil
}
```

### **Frontend Optimization**
- **Lazy loading**: Use Angular lazy loading for modules
- **Change detection**: Optimize change detection strategies
- **Bundle size**: Monitor and reduce bundle size
- **Caching**: Implement service worker for caching
- **Performance monitoring**: Use Angular performance tools

## 🔍 **Troubleshooting**

### **Common Issues and Solutions**

#### **Service Won't Start**
```bash
# Check if infrastructure is running
./scripts/scopeapi.sh setup --validate

# Check service logs
cd scripts
./dev.sh logs [service-name]

# Check Docker status
docker ps
docker logs [container-name]
```

#### **Database Connection Issues**
```bash
# Verify PostgreSQL is running
docker exec -it scopeapi-postgres pg_isready

# Check database connection
./scripts/setup-database.sh --validate

# Reset database if needed
docker-compose down -v
docker-compose up -d postgres
./scripts/setup-database.sh
```

#### **Port Conflicts**
```bash
# Check what's using a port
sudo lsof -i :8080

# Kill process if needed
sudo kill -9 [PID]

# Or use different ports in docker-compose.override.yml
```

#### **Memory Issues**
```bash
# Check container memory usage
docker stats

# Increase memory limits in docker-compose.yml
services:
  api-discovery:
    deploy:
      resources:
        limits:
          memory: 1G
```

### **Debugging Network Issues**
```bash
# Check network connectivity
docker network ls
docker network inspect scopeapi-network

# Test service communication
docker exec -it scopeapi-api-discovery wget -qO- http://scopeapi-postgres:5432
```

### **Performance Issues**
```bash
# Monitor resource usage
docker stats

# Check service health
./scopeapi.sh status

# View performance metrics
curl http://localhost:8080/metrics
```

## 📚 **Additional Resources**

### **Documentation**
- **[Architecture Guide](ARCHITECTURE.md)** - System design and technical details
- **[API Reference](API.md)** - Service APIs and endpoints
- **[Docker Setup](DOCKER_SETUP.md)** - Container and deployment setup

### **External Resources**
- **[Go Documentation](https://golang.org/doc/)** - Official Go language docs
- **[Angular Documentation](https://angular.io/docs)** - Official Angular docs
- **[Docker Documentation](https://docs.docker.com/)** - Docker and Docker Compose
- **[PostgreSQL Documentation](https://www.postgresql.org/docs/)** - Database reference

### **Community and Support**
- **GitHub Issues**: Report bugs and request features
- **Discussions**: Ask questions and share ideas
- **Contributing Guide**: Learn how to contribute
- **Code of Conduct**: Community guidelines

---

**🎯 This development guide helps you:**
- **Set up** your development environment quickly
- **Follow** best practices and coding standards
- **Debug** issues effectively
- **Optimize** performance and code quality
- **Contribute** to the open-source project

**Happy coding! 🚀✨**

## 🏗️ **Script Architecture & Usage**

This project provides specialized scripts for different development approaches:

### **🔄 Local Development (Process-Based)**
- **Best for**: Direct Go binary development, debugging, performance testing
- **Workflow**: Infrastructure → Process management → Direct binary execution
- **Benefits**: Faster startup, direct process control, easier debugging
- **Use when**: Developing Go services, testing performance, debugging issues

### **🐳 Container-Based Development ()**
- **Best for**: Production-like development, container management, deployment testing
- **Workflow**: Complete orchestration → Container management → Service lifecycle
- **Benefits**: Production parity, container isolation, easy scaling
- **Use when**: Testing deployment, container debugging, production preparation

### **🔧 Infrastructure Management ()**
- **Best for**: Infrastructure setup, troubleshooting, environment management
- **Workflow**: Environment setup → Infrastructure management → Troubleshooting
- **Benefits**: Focused infrastructure control, permission fixing, environment setup
- **Use when**: Setting up environment, fixing infrastructure issues, managing dependencies

### **🎯 Development Workflow Selection Guide:**

| Development Phase | Recommended Script | Alternative |
|------------------|-------------------|-------------|
| **Initial setup** | `./scopeapi.sh setup --full` | - |
| **Local development (Go services)** | `./dev.sh start all` | `./scopeapi.sh start all` |
| **Container development** | `./dev.sh start all` | - |
| **Infrastructure issues** | `./infrastructure.sh start` | - |
| **Debugging** | `./dev.sh debug [service]` | - |
| **Production testing** | `./deploy.sh -e staging -p k8s` | - |
