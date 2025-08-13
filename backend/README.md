# 🚀 **ScopeAPI Backend Services**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](../LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/scopeapi/backend)](https://goreportcard.com/report/github.com/scopeapi/backend)

**Backend microservices** for the ScopeAPI platform, built with **Go 1.21+** and following microservices architecture principles.

## 🏗️ **Architecture Overview**

```
┌──────────────────────────────────────────────────────────────────┐
│                        ScopeAPI Backend                          │
├──────────────────────────────────────────────────────────────────┤
│                    Microservices Layer                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │API Discovery│ │Threat Detect│ │Data Protect │ │Attack Block │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐                 │
│  │Gateway Integ│ │Data Ingest  │ │Admin Console│                 │
│  └─────────────┘ └─────────────┘ └─────────────┘                 │
├──────────────────────────────────────────────────────────────────┤
│                    Shared Libraries                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │   Auth      │ │  Database   │ │  Logging    │ │  Monitoring │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

## 📁 **Project Structure**

```
backend/
├── 📁 services/                    # Individual microservices
│   ├── 📁 api-discovery/           # API discovery & cataloging
│   │   ├── 📁 cmd/                 # Service entry points
│   │   ├── 📁 internal/            # Private application code
│   │   │   ├── 📁 handlers/        # HTTP request handlers
│   │   │   ├── 📁 models/          # Data structures
│   │   │   ├── 📁 repository/      # Data access layer
│   │   │   └── 📁 services/        # Business logic
│   │   ├── 📄 Dockerfile           # Container configuration
│   │   └── 📄 go.mod               # Go module dependencies
│   │
│   ├── 📁 threat-detection/        # Security threat detection
│   ├── 📁 data-protection/         # PII & compliance
│   ├── 📁 attack-blocking/         # Real-time threat blocking
│   ├── 📁 gateway-integration/     # API gateway management
│   ├── 📁 data-ingestion/          # Traffic data processing
│   └── 📁 admin-console/           # Backend for admin UI
│
├── 📁 shared/                      # Shared libraries & utilities
│   ├── 📁 auth/                    # JWT authentication
│   ├── 📁 database/                # Database connections
│   ├── 📁 logging/                 # Structured logging
│   ├── 📁 messaging/               # Kafka integration
│   ├── 📁 monitoring/              # Health checks & metrics
│   └── 📁 utils/                   # Common utilities
│
├── 📁 bin/                         # Compiled binaries (gitignored)
├── 📄 go.mod                       # Root module dependencies
├── 📄 go.work                      # Go workspace configuration
└── 📄 Makefile                     # Build automation
```

## 🔧 **Technology Stack**

### **Core Technologies**
- **Language**: [Go 1.21+](https://golang.org) - High-performance, concurrent programming
- **Architecture**: Microservices with event-driven communication
- **Communication**: RESTful HTTP APIs + Apache Kafka messaging
- **Containerization**: Docker with multi-stage builds

### **Data Layer**
- **Primary Database**: PostgreSQL 15+ (relational data)
- **Caching**: Redis 7+ (sessions, rate limiting)
- **Search**: Elasticsearch (logging, analytics)
- **Message Queue**: Apache Kafka 3.4+ (inter-service communication)

### **Infrastructure**
- **Service Discovery**: Built-in health checks and monitoring
- **Configuration**: Environment-based configuration management
- **Logging**: Structured logging with correlation IDs
- **Monitoring**: Prometheus metrics + health endpoints

## 🚀 **Quick Start**

### **Prerequisites**
```bash
# Install Go 1.21+
go version

# Install Docker & Docker Compose
docker --version
docker-compose --version
```

### **Development Setup**
```bash
# Clone and navigate
git clone https://github.com/your-org/scopeapi.git
cd scopeapi/backend

# Install dependencies
go mod download

# Build all services
make all

# Or build individual service
make api-discovery
make threat-detection
make data-protection
make attack-blocking
make gateway-integration
make data-ingestion
make admin-console
```

### **Running Services**
```bash
# Start infrastructure (PostgreSQL, Kafka, Redis)
../scripts/docker-infrastructure.sh

# Start all microservices
../scripts/scopeapi-services.sh start all

# Start specific service
../scripts/scopeapi-services.sh start api-discovery

# View logs
../scripts/scopeapi-services.sh logs api-discovery
```

## 📊 **Service Details**

### **🔍 API Discovery Service**
- **Port**: 8080
- **Purpose**: Automatically discover and catalog API endpoints
- **Features**: Endpoint crawling, change detection, metadata management
- **Dependencies**: PostgreSQL, Kafka

### **🛡️ Threat Detection Service**
- **Port**: 8081
- **Purpose**: Real-time security threat identification
- **Features**: ML-based detection, behavioral analysis, threat intelligence
- **Dependencies**: PostgreSQL, Elasticsearch, Kafka

### **🔒 Data Protection Service**
- **Port**: 8082
- **Purpose**: Sensitive data detection and compliance
- **Features**: PII detection, data classification, compliance monitoring
- **Dependencies**: PostgreSQL, Elasticsearch

### **⚡ Attack Blocking Service**
- **Port**: 8083
- **Purpose**: Real-time threat prevention and blocking
- **Features**: Request filtering, rate limiting, IP blocking
- **Dependencies**: PostgreSQL, Redis, Kafka

### **🌐 Gateway Integration Service**
- **Port**: 8084
- **Purpose**: Multi-gateway configuration management
- **Features**: Kong, Envoy, HAProxy, Nginx, Traefik support
- **Dependencies**: PostgreSQL, Kafka

### **📥 Data Ingestion Service**
- **Port**: 8085
- **Purpose**: Traffic data processing and normalization
- **Features**: Data parsing, normalization, streaming
- **Dependencies**: PostgreSQL, Kafka

### **🖥️ Admin Console Service**
- **Port**: 8086
- **Purpose**: Backend for admin interface
- **Features**: User management, system configuration, monitoring
- **Dependencies**: PostgreSQL, Redis, Kafka

## 🧪 **Development Workflow**

### **Building Services**
```bash
# Build all services
make all

# Build specific service
make api-discovery

# Clean binaries
make clean

# Show help
make help
```

### **Testing**
```bash
# Run all tests
go test ./...

# Test specific service
cd services/api-discovery
go test ./...

# Run with coverage
go test -cover ./...
```

### **Debugging**
```bash
# Start service in debug mode
../scripts/scopeapi-debug.sh start api-discovery

# Connect IDE to localhost:2345
# Set breakpoints and debug
```

## 📚 **Documentation**

- **[Main Project README](../README.md)** - Project overview and quick start
- **[Architecture Guide](../docs/ARCHITECTURE.md)** - Detailed system design
- **[Development Guide](../docs/DEVELOPMENT.md)** - Development setup and workflow
- **[API Reference](../docs/API.md)** - Service APIs and endpoints
- **[Docker Setup](../docs/DOCKER_SETUP.md)** - Containerization guide

## 🤝 **Contributing**

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### **Code Standards**
- Follow Go best practices and [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful variable and function names
- Add tests for new functionality
- Update documentation for API changes

## 📄 **License**

This project is licensed under the **MIT License** - see the [LICENSE](../LICENSE) file for details.

## 🔗 **Related Repositories**

- **[Frontend Admin Console](../adminConsole/)** - Angular-based admin interface
- **[Infrastructure Scripts](../scripts/README.md)** - Development and deployment automation
- **[Documentation](../docs/)** - Comprehensive project documentation

---

**🚀 Ready to build secure APIs?**
- **Star** this repository if you find it useful
- **Fork** to contribute or customize
- **Share** with your team and community

**Happy coding! 🎉✨**
