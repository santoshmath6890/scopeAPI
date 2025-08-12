# API Discovery Service

## Overview

The API Discovery Service is a **core component** of the ScopeAPI platform that provides automated API endpoint discovery, inventory management, and metadata analysis capabilities. It enables organizations to automatically discover, catalog, and analyze APIs across their infrastructure.

## 🎯 **Service Status: IMPLEMENTED** 🔄

This service has core functionality implemented and is ready for development and testing.

## Architecture Integration

This service is part of the **Core Services Layer** in the ScopeAPI architecture:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Core Services Layer                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  Endpoint       │  │  Threat         │  │  Attack         │             │
│  │  Discovery      │  │  Detection      │  │  Blocking       │             │
│  │  Service        │  │  Engine         │  │  Engine         │             │
│  │  🔄 IMPLEMENTED │  │                 │  │                 │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  Sensitive      │  │  Security       │  │  Gateway        │             │
│  │  Data Scanner   │  │  Testing        │  │  Integration    │             │
│  │                 │  │  Engine         │  │  Service        │             │
│  │ • PII detection │  │ • Automated     │  │                 │             │
│  │ • Data classify │  │ • Vuln scanning │  │ • Kong/NGINX    │             │
│  │ • Compliance    │  │ • Pen testing   │  │ • Traefik/Envoy │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### **Integration Points**

- **API Gateway Layer**: Discovers APIs from Kong, NGINX, Traefik, Envoy, HAProxy
- **Data Storage Layer**: Uses PostgreSQL for API inventory and metadata storage
- **Frontend**: Provides Angular components for API discovery UI
- **Monitoring**: Health checks and readiness probes

## 🚀 **Quick Start**

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional)

### Option 1: Local Development

```bash
# Clone the repository
cd backend/services/api-discovery

# Install dependencies
make deps

# Build the service
make build

# Run the service
make run

# Or run in development mode
make run-dev
```

### Option 2: Docker

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run
```

### Option 3: Using Root Makefile

```bash
# Show all available commands
cd backend
make help

# Build all services
make all

# Build specific service
make api-discovery

# Clean build artifacts
make clean
```

## 🏗️ **Service Architecture**

### **Service Structure** 🔄
```
api-discovery/
├── cmd/
│   └── main.go                 # ✅ Service entry point (IMPLEMENTED)
├── config/
│   └── api-discovery.yaml      # ✅ Configuration file (IMPLEMENTED)
├── internal/
│   ├── handlers/               # ✅ HTTP request handlers (IMPLEMENTED)
│   │   ├── discovery_handler.go # ✅ Discovery operations
│   │   ├── inventory_handler.go # ✅ Inventory management
│   │   └── endpoint_handler.go  # ✅ Endpoint analysis
│   ├── models/                 # ✅ Data models (IMPLEMENTED)
│   │   ├── api_spec.go         # ✅ API specification models
│   │   ├── endpoint.go         # ✅ Endpoint models
│   │   └── metadata.go         # ✅ Metadata models
│   ├── repository/             # ✅ Database operations (IMPLEMENTED)
│   │   ├── discovery_repository.go # ✅ Discovery operations
│   │   ├── inventory_repository.go # ✅ Inventory operations
│   │   └── repository.go       # ✅ Base repository interface
│   └── services/               # ✅ Business logic (IMPLEMENTED)
│       ├── discovery_service.go # ✅ Discovery logic
│       ├── inventory_service.go # ✅ Inventory management
│       └── metadata_service.go # ✅ Metadata analysis
├── Dockerfile                  # ✅ Container configuration (NEW)
├── go.mod                      # ✅ Go module dependencies
└── README.md                   # ✅ This file (NEW)
```

## 🔍 **Core Functionality** 🔄

### **API Discovery**
- 🔄 **Automated scanning** of API endpoints
- 🔄 **Endpoint detection** and cataloging
- 🔄 **API specification** parsing and analysis
- 🔄 **Real-time discovery** status tracking

### **Inventory Management**
- 🔄 **API catalog** maintenance
- 🔄 **Endpoint metadata** storage
- 🔄 **Version tracking** for APIs
- 🔄 **Change detection** and monitoring

### **Metadata Analysis**
- 🔄 **Endpoint analysis** and classification
- 🔄 **API documentation** extraction
- 🔄 **Security assessment** data collection
- 🔄 **Compliance information** gathering

## 🔗 **API Endpoints** 🔄

### **Discovery Operations**
```
POST   /api/v1/discovery/scan           # Start API discovery scan
GET    /api/v1/discovery/status/:id     # Get discovery scan status
```

### **Inventory Management**
```
GET    /api/v1/inventory/apis           # List all discovered APIs
GET    /api/v1/inventory/apis/:id       # Get specific API details
```

### **Endpoint Analysis**
```
POST   /api/v1/endpoints/analyze        # Analyze endpoint metadata
GET    /api/v1/endpoints/:id/metadata   # Get endpoint metadata
```

### **Health & Monitoring**
```
GET    /health                           # Service health status
GET    /ready                            # Service readiness check
```

## 🔒 **Security Features** 🔄

- 🔄 **Input validation** and sanitization
- 🔄 **Database connection** security
- 🔄 **Error handling** without information leakage
- 🔄 **Secure configuration** management

## 📊 **Monitoring & Observability** 🔄

### **Health Checks**
```
GET /health                    # Service health status
GET /ready                     # Service readiness
```

### **Logging**
- 🔄 **Structured logging** with JSON format
- 🔄 **Log levels** (debug, info, warn, error)
- 🔄 **Request/response logging**
- 🔄 **Error tracking** and stack traces

## 🗄️ **Database Schema** 🔄

### **Core Tables**
- 🔄 **`api_endpoints`** - Discovered API endpoints
- 🔄 **`api_metadata`** - Endpoint metadata and analysis
- 🔄 **`discovery_scans`** - Discovery scan history
- 🔄 **`api_inventory`** - API catalog and inventory

## 🚀 **Deployment** 🔄

### **Docker Deployment**
```bash
# Build image
docker build -t scopeapi/api-discovery .

# Run container
docker run -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e DB_HOST=your-postgres-host \
  -e DB_PORT=5432 \
  -e DB_USER=scopeapi \
  -e DB_PASSWORD=your_secure_password \
  -e DB_NAME=scopeapi \
  scopeapi/api-discovery
```

### **Kubernetes Deployment**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-discovery
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api-discovery
  template:
    metadata:
      labels:
        app: api-discovery
    spec:
      containers:
      - name: api-discovery
        image: scopeapi/api-discovery:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_PORT
          value: "8080"
        - name: DB_HOST
          value: "postgres-service"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          value: "scopeapi"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        - name: DB_NAME
          value: "scopeapi"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## 🧪 **Testing** 🔄

### **Test Coverage**
- 🔄 **Unit tests** for services and handlers
- 🔄 **Integration tests** for database operations
- 🔄 **API tests** for all endpoints

### **Running Tests**
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmark tests
make test-bench
```

## 🔧 **Configuration** 🔄

### **Environment Variables**
```bash
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=scopeapi
DB_SSL_MODE=disable
```

### **Configuration File**
The service uses `config/api-discovery.yaml` for configuration management with support for:
- 🔄 **Server settings** configuration
- 🔄 **Database connection** parameters
- 🔄 **Environment-specific** configurations

## 📈 **Performance & Scalability** 🔄

### **Performance Features**
- 🔄 **Efficient database queries** with proper indexing
- 🔄 **Async processing** for discovery operations
- 🔄 **Connection pooling** for database connections

### **Scalability Features**
- 🔄 **Stateless design** for horizontal scaling
- 🔄 **Health checks** for load balancer integration
- 🔄 **Readiness probes** for Kubernetes deployment

## 🚨 **Error Handling & Resilience** 🔄

### **Error Handling**
- 🔄 **Comprehensive error types** and messages
- 🔄 **Graceful degradation** for partial failures
- 🔄 **Detailed error logging** for debugging

### **Resilience Features**
- 🔄 **Health check endpoints** for monitoring
- 🔄 **Graceful shutdown** handling
- 🔄 **Connection retry logic** for databases
- 🔄 **Resource cleanup** on failures

## 🔄 **API Versioning** 🔄

- 🔄 **RESTful API design** following best practices
- 🔄 **Versioned endpoints** (`/api/v1/`)
- 🔄 **Backward compatibility** support

## 📚 **Documentation** 🔄

- 🔄 **API endpoint documentation**
- 🔄 **Configuration examples**
- 🔄 **Deployment guides**
- 🔄 **Development setup instructions**

## 🤝 **Contributing** 🔄

### **Development Setup**
```bash
# Clone and setup
git clone <repository>
cd backend/services/api-discovery

# Install development tools
make install-tools

# Setup development environment
make dev-setup

# Run tests
make test

# Build service
make build
```

### **Code Quality**
- 🔄 **Go linting** with golangci-lint
- 🔄 **Code formatting** with go fmt
- 🔄 **Vet checks** with go vet
- 🔄 **Test coverage** requirements

## 📞 **Support & Maintenance** 🔄

### **Monitoring**
- 🔄 **Health check endpoints** for load balancers
- 🔄 **Structured logging** for log aggregation
- 🔄 **Error tracking** and alerting

### **Maintenance**
- 🔄 **Database migration** support
- 🔄 **Configuration backup** and restore
- 🔄 **Version upgrade** procedures

## 🎯 **Service Status Summary**

The API Discovery Service currently has:

- ✅ **Core functionality** implemented (discovery, inventory, metadata)
- ✅ **HTTP handlers** for all endpoints
- ✅ **Database models** and repositories
- ✅ **Business logic services**
- ✅ **Configuration management**
- ✅ **Health checks** and monitoring
- ✅ **Dockerfile** (NEW)
- ✅ **Comprehensive documentation** (NEW)

## 🚀 **Next Steps**

The service is ready for:
1. **Integration testing** with real API endpoints
2. **Performance testing** under load
3. **Security auditing** and penetration testing
4. **User acceptance testing** with the frontend
5. **Production deployment** preparation

## 📞 **Contact & Support**

For questions about the service:
- **Documentation**: This README and API documentation
- **Issues**: GitHub issue tracker
- **Contributions**: Pull request guidelines
- **Support**: Development team contacts

---

**🎯 Status: IMPLEMENTED** 🔄
**📅 Last Updated**: $(date)
**🔄 Version**: Development
**👥 Maintainers**: ScopeAPI Development Team

