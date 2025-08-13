# 🚀 ScopeAPI

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Angular Version](https://img.shields.io/badge/Angular-16+-red.svg)](https://angular.io)
[![Docker Version](https://img.shields.io/badge/Docker-24+-blue.svg)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**ScopeAPI** is a comprehensive **API security and management platform** designed to protect, monitor, and manage APIs in modern distributed systems. It provides a unified approach to API security across multiple domains with enterprise-grade capabilities.

## 🎯 **What is ScopeAPI?**

ScopeAPI is an **open-source platform** that helps organizations secure their APIs through:

- **🔍 API Discovery & Cataloging** - Automatically discover and catalog API endpoints
- **🛡️ Threat Detection & Prevention** - Real-time security threat identification and blocking
- **🔒 Data Protection & Compliance** - Sensitive data detection and regulatory compliance
- **⚡ Attack Blocking** - Real-time threat prevention and blocking
- **🌐 Gateway Integration** - Seamless integration with popular API gateways
- **📊 Centralized Management** - Unified admin console for all security operations

## 🏗️ **Architecture**

ScopeAPI follows a **microservices architecture** with:

- **7 Core Microservices** - Each handling a specific security domain
- **Event-Driven Communication** - Kafka-based message queuing
- **Polyglot Persistence** - PostgreSQL, Redis, Elasticsearch
- **Containerized Deployment** - Docker and Docker Compose
- **RESTful APIs** - Standard HTTP interfaces for all services

```
┌──────────────────────────────────────────────────────────────────┐
│                        ScopeAPI Platform                         │
├──────────────────────────────────────────────────────────────────┤
│  Admin Console (Angular)  │  API Gateway (Kong/Envoy/Nginx)      │
├──────────────────────────────────────────────────────────────────┤
│                    Microservices Layer                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │API Discovery│ │Threat Detect│ │Data Protect │ │Attack Block │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐                 │
│  │Gateway Integ│ │Data Ingest  │ │Admin Console│                 │
│  └─────────────┘ └─────────────┘ └─────────────┘                 │
├──────────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │  PostgreSQL │ │    Kafka    │ │    Redis    │ │Elasticsearch│ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
└──────────────────────────────────────────────────────────────────┘
```

## 🚀 **Quick Start**

### **Prerequisites**
- **Docker** 24.0+ with Docker Compose
- **Go** 1.21+ (for backend development)
- **Node.js** 18+ (for admin console)

### **1. Clone the Repository**
```bash
git clone https://github.com/your-org/scopeapi.git
cd scopeapi
```

### **2. Complete Setup**
```bash
# Complete setup with validation
./scripts/scopeapi-setup.sh --full

# This will:
# - Start infrastructure services
# - Setup PostgreSQL database
# - Run migrations
# - Create test data
# - Validate everything
```

### **3. Start Development**
```bash
# Start all services for development
./scripts/scopeapi-services.sh start all

# Or start specific service
./scripts/scopeapi-services.sh start api-discovery
```

### **4. Access Services**
- **Admin Console**: http://localhost:8086
- **API Discovery**: http://localhost:8080
- **Gateway Integration**: http://localhost:8081
- **Data Ingestion**: http://localhost:8082
- **Threat Detection**: http://localhost:8083
- **Data Protection**: http://localhost:8084
- **Attack Blocking**: http://localhost:8085

## 🔧 **Development Workflows**

### **Daily Development**
```bash
# Start services
./scripts/scopeapi-services.sh start all

# Make code changes
# View logs if needed
./scripts/scopeapi-services.sh logs api-discovery

# Stop when done
./scripts/scopeapi-services.sh stop
```

### **Debugging**
```bash
# Start service in debug mode
./scripts/scopeapi-debug.sh start api-discovery

# Connect IDE to localhost:2345
# Set breakpoints and debug
```

### **Testing**
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd adminConsole && ng test

# Integration tests
./scripts/setup-database.sh --validate
```

## 📚 **Documentation**

- **[📖 Documentation Index](docs/INDEX.md)** - Complete documentation navigation
- **[🏗️ Architecture Guide](docs/ARCHITECTURE.md)** - System design and technical details
- **[💻 Development Guide](docs/DEVELOPMENT.md)** - Development setup and workflows
- **[🐳 Docker Setup](docs/DOCKER_SETUP.md)** - Container and deployment setup
- **[🛠️ Scripts Usage](scripts/USAGE.md)** - Development scripts guide

## 🤝 **Contributing**

We welcome contributions from the community! Please see our **[Contributing Guide](docs/CONTRIBUTING.md)** for details on:

- **Code Standards** - Coding conventions and best practices
- **Development Setup** - How to set up your development environment
- **Pull Request Process** - How to submit your changes
- **Testing Guidelines** - How to test your contributions

### **Quick Contribution Start**
```bash
# Fork and clone
git clone https://github.com/your-username/scopeapi.git
cd scopeapi

# Setup development environment
./scripts/scopeapi-setup.sh --full

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
./scripts/scopeapi-services.sh start all

# Commit and push
git commit -m "Add amazing feature"
git push origin feature/amazing-feature

# Create Pull Request
```

## 🏗️ **Project Structure**

```
scopeapi/
├── 📁 backend/                     # Go microservices
│   ├── 📁 services/                # Individual microservices
│   │   ├── 📁 api-discovery/       # API discovery service
│   │   ├── 📁 threat-detection/    # Threat detection service
│   │   ├── 📁 data-protection/     # Data protection service
│   │   ├── 📁 attack-blocking/     # Attack blocking service
│   │   ├── 📁 gateway-integration/ # Gateway integration service
│   │   ├── 📁 data-ingestion/      # Data ingestion service
│   │   └── 📁 admin-console/       # Admin console backend service
│   └── 📁 shared/                  # Shared libraries and utilities
├── 📁 adminConsole/                # Angular frontend application
├── �� scripts/                     # Project automation and management scripts
│   ├── 🔄 scopeapi-local.sh        # Local development (process-based management)
│   ├── 🐳 docker-infrastructure.sh # Infrastructure management
│   ├── 🚀 scopeapi-services.sh     # Container-based microservices orchestration
│   ├── 🔧 scopeapi-setup.sh        # Complete setup and validation
│   └── 🐛 scopeapi-debug.sh        # Debug mode management
├── 📁 docs/                        # Comprehensive documentation
└── 📁 README.md                    # This file
```

## 🚀 **Deployment**

### **Local Development**
```bash
./scripts/scopeapi-setup.sh --full
./scripts/scopeapi-services.sh start all
```

### **Production**
```bash
# Deploy with Docker Compose
docker-compose -f scripts/docker-compose.yml up -d

# Or deploy to Kubernetes
kubectl apply -f k8s/
```

## 📊 **Features**

### **🔍 API Discovery**
- **Automatic Endpoint Discovery** - Crawl and catalog API endpoints
- **Change Detection** - Monitor API changes and versioning
- **Documentation Generation** - Auto-generate API documentation
- **Metadata Management** - Rich metadata and tagging

### **🛡️ Threat Detection**
- **Real-time Analysis** - Continuous security monitoring
- **Machine Learning** - AI-powered threat detection
- **Behavioral Analysis** - User and API behavior monitoring
- **Threat Intelligence** - Integration with threat feeds

### **🔒 Data Protection**
- **PII Detection** - Automatic sensitive data identification
- **Data Classification** - Intelligent data categorization
- **Compliance Monitoring** - Regulatory requirement tracking
- **Audit Logging** - Comprehensive audit trails

### **⚡ Attack Blocking**
- **Real-time Filtering** - Request validation and filtering
- **Rate Limiting** - Adaptive rate limiting and throttling
- **IP Blocking** - Geographic and reputation-based blocking
- **Pattern Recognition** - Attack pattern identification

### **🌐 Gateway Integration**
- **Multi-Gateway Support** - Kong, Envoy, HAProxy, Nginx, Traefik
- **Policy Management** - Centralized policy configuration
- **Health Monitoring** - Gateway health and performance
- **Configuration Sync** - Automated policy deployment

## 🛠️ **Technology Stack**

### **Backend**
- **Language**: Go 1.21+
- **Framework**: Standard library + custom middleware
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Message Queue**: Apache Kafka 3.4+

### **Frontend**
- **Framework**: Angular 16+
- **Language**: TypeScript 5+
- **Styling**: SCSS with modern CSS features
- **Build Tool**: Angular CLI with Webpack

### **Infrastructure**
- **Containerization**: Docker 24+
- **Orchestration**: Docker Compose, Kubernetes
- **Monitoring**: Prometheus, Grafana, ELK Stack
- **CI/CD**: GitHub Actions, GitLab CI

## 📈 **Performance & Scalability**

- **Horizontal Scaling** - All services scale independently
- **Event-Driven Architecture** - Asynchronous processing
- **Caching Strategies** - Multi-layer caching for performance
- **Load Balancing** - Intelligent request distribution
- **99.9% Uptime** - High availability and reliability

## 🔒 **Security Features**

- **Zero-Trust Architecture** - No implicit trust between services
- **Multi-Factor Authentication** - Enhanced access security
- **Role-Based Access Control** - Granular permission management
- **Encryption** - Data encryption at rest and in transit
- **Audit Logging** - Comprehensive security event tracking

## 🌟 **Why Choose ScopeAPI?**

### **✅ Open Source**
- **Transparent** - Full source code visibility
- **Community Driven** - Active community contributions
- **No Vendor Lock-in** - Complete control over your deployment

### **✅ Enterprise Ready**
- **Production Grade** - Built for enterprise environments
- **Scalable** - Handles growth and increased load
- **Secure** - Security-first design principles
- **Compliant** - Built-in compliance and audit features

### **✅ Developer Friendly**
- **Easy Setup** - Simple development environment setup
- **Comprehensive Tooling** - Scripts, debugging, and monitoring
- **Clear Documentation** - Well-documented APIs and workflows
- **Testing Support** - Built-in testing and validation tools

## 🤝 **Community & Support**

### **Getting Help**
- **📖 Documentation**: Comprehensive guides and references
- **🐛 Issues**: Report bugs and request features
- **💬 Discussions**: Ask questions and share ideas
- **📧 Email**: Contact the maintainers directly

### **Contributing**
- **Code Contributions** - Bug fixes, features, and improvements
- **Documentation** - Help improve guides and references
- **Testing** - Report bugs and test new features
- **Community** - Help other users and contributors

## 📄 **License**

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙏 **Acknowledgments**

- **Go Community** - For the excellent Go ecosystem
- **Angular Team** - For the powerful frontend framework
- **Docker Community** - For containerization tools
- **Open Source Contributors** - For making this project possible

---

**🎯 Ready to secure your APIs?**
- **Star** this repository if you find it useful
- **Fork** to contribute or customize
- **Share** with your team and community
- **Contribute** to make it even better

**Happy coding! 🚀✨** 
