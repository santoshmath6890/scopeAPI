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
./scripts/scopeapi-setup.sh --full

# Verify setup
./scripts/scopeapi-setup.sh --validate
```

### **2. Environment Configuration**
```bash
# Copy environment template
cp env.example .env

# Edit with your configuration
nano .env

# Key environment variables:
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_secure_password
KAFKA_BROKER_ID=1
```

### **3. Database Setup**
```bash
# Setup database with test data
./scripts/setup-database.sh --test-data

# Verify database
./scripts/setup-database.sh --validate
```

## 🔧 **Development Workflows**

### **Daily Development Workflow**
```bash
# 1. Start infrastructure
./scripts/scopeapi-setup.sh --infrastructure

# 2. Start services for development
./scripts/./scripts/scopeapi-services.sh start all

# 3. Make code changes
# 4. View logs if needed
./scripts/./scripts/scopeapi-services.sh logs api-discovery

# 5. Test changes
# 6. Stop when done
./scripts/./scripts/scopeapi-services.sh stop
```

### **Service-Specific Development**
```bash
# Start only specific service
./scripts/./scripts/scopeapi-services.sh start api-discovery

# Start multiple services
./scripts/./scripts/scopeapi-services.sh start api-discovery gateway-integration

# View specific service logs
./scripts/./scripts/scopeapi-services.sh logs api-discovery

# Open shell in service container
./scripts/./scripts/scopeapi-services.sh shell api-discovery
```

### **Debugging Workflow**
```bash
# 1. Start service in debug mode
./scripts/scopeapi-debug.sh start api-discovery

# 2. Connect IDE to debug port (2345 for api-discovery)
# 3. Set breakpoints in your Go code
# 4. Debug and step through code
# 5. Stop debug session
./scripts/scopeapi-debug.sh stop
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
./scripts/./scripts/scopeapi-services.sh start all
# Verify all services are healthy
./scripts/./scripts/scopeapi-services.sh status
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
./scripts/scopeapi-debug.sh start api-discovery

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
./scripts/./scripts/scopeapi-services.sh status
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
./scripts/scopeapi-setup.sh --validate

# Check service logs
./scripts/./scripts/scopeapi-services.sh logs [service-name]

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
./scripts/./scripts/scopeapi-services.sh status

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
| **Initial setup** | `scopeapi-setup.sh` | - |
| **Local development (Go services)** | `./scripts/scopeapi-local.sh` | `./scripts/scopeapi-services.sh` |
| **Container development** | `./scripts/scopeapi-services.sh` | - |
| **Infrastructure issues** | `./scripts/docker-infrastructure.sh` | - |
| **Debugging** | `scopeapi-debug.sh` | - |
| **Production testing** | `./scripts/scopeapi-services.sh` | - |
