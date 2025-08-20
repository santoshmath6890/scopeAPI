# 🚀 ScopeAPI

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Angular Version](https://img.shields.io/badge/Angular-16.2+-red.svg)](https://angular.io)
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
- **Docker** 24.0+ with Docker Compose (for local development)
- **Kubernetes** cluster (for staging/production)
- **Go** 1.21+ (for backend development)
- **Node.js** 18+ (for admin console)

### **Environment Strategy**
- **Local Development**: Use `.env.local` file (your machine only)
- **Staging/Production**: Use Kubernetes Secrets (secure, encrypted)

### **1. Clone the Repository**
```bash
git clone https://github.com/your-org/scopeapi.git
cd scopeapi
```

### **2. Complete Setup**
```bash
# Complete setup with validation
cd scripts
./scopeapi.sh setup --full

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
cd scripts
./dev.sh start all

# Or start specific service
./dev.sh start api-discovery
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

For detailed development workflows, debugging, and testing instructions, see:
- **[💻 Development Guide](docs/DEVELOPMENT.md)** - Complete development setup and workflows
- **[🛠️ Scripts Usage](scripts/README.md)** - Development script commands and examples

## 📚 **Documentation**

- **[📖 Documentation Index](docs/README.md)** - Complete documentation navigation
- **[🏗️ Architecture Guide](docs/ARCHITECTURE.md)** - System design and technical details
- **[💻 Development Guide](docs/DEVELOPMENT.md)** - Development setup and workflows
- **[🚀 Deployment Guide](docs/DEPLOYMENT.md)** - Environment strategy, security, and deployment
- **[🛠️ Scripts Usage](scripts/README.md)** - Development scripts guide
- **[☸️ Kubernetes Config](k8s/README.md)** - Kubernetes deployment configurations

## 🤝 **Contributing**

We welcome contributions from the community! Please see our **[Contributing Guide](docs/CONTRIBUTING.md)** for complete details on contributing to ScopeAPI.

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
├── 📁 scripts/                     # Project automation and management scripts
│   ├── 🎯 scopeapi.sh              # Main orchestrator (setup, services, status)
│   ├── 🏗️ infrastructure.sh        # Infrastructure management
│   ├── 🚀 deploy.sh                # Deployment (Docker + K8s)
│   ├── 💻 dev.sh                   # Development workflows
│   └── 🔧 setup-database.sh        # Database setup utilities
├── 📁 k8s/                         # Kubernetes deployment configurations
│   ├── 📁 deployments/             # All microservices + admin console
│   ├── 📁 services/                # Network services configuration
│   ├── 📁 ingress/                 # Traffic routing and load balancing
│   ├── 📁 secrets/                 # Environment variables and secrets
│   └── 📁 configmaps/              # Application configuration
├── 📁 docs/                        # Comprehensive documentation
└── 📄 README.md                    # This file
```

## 🚀 **Deployment**

For comprehensive deployment instructions, environment strategy, and security guidelines, see:
- **[🚀 Deployment Guide](docs/DEPLOYMENT.md)** - Complete deployment guide with environment strategy

### **Kubernetes Configuration**
The `k8s/` directory contains all Kubernetes deployment configurations for staging and production environments, including:
- **Deployments** for all microservices and the admin console
- **Services** for network communication
- **Ingress** for traffic routing and load balancing
- **Secrets** for secure environment variable management
- **ConfigMaps** for application configuration

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

## 🖥️ **Admin Console (Frontend)**

**Modern Angular 16+ web application** providing a comprehensive interface for managing all ScopeAPI services. Features responsive design, lazy-loaded modules, real-time updates, and role-based access control.

**📱 Key Features**: Dashboard, API Discovery, Threat Detection, Data Protection, Attack Protection, Gateway Integration, and Authentication modules.

**🚀 Quick Start**: `cd adminConsole && npm install && npm start`

**📚 [Detailed Documentation →](docs/README.md#admin-console)**

For complete technology stack details, see **[🏗️ Architecture Guide](docs/ARCHITECTURE.md)**

## 📈 **Performance & Scalability**

For detailed performance characteristics, scalability features, and security architecture, see **[🏗️ Architecture Guide](docs/ARCHITECTURE.md)**

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

**We welcome contributions from the community!** For complete information on getting involved, see our **[Contributing Guide](docs/CONTRIBUTING.md)**.

## 📄 **License**

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙏 **Acknowledgments**

- **Open Source Contributors** - For making this project possible

---

**🎯 Ready to secure your APIs?**
- **Star** this repository if you find it useful
- **Fork** to contribute or customize
- **Share** with your team and community
- **Contribute** to make it even better

**Happy coding! 🚀✨** 
