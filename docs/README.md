# 📚 ScopeAPI Documentation

Welcome to the ScopeAPI documentation! This is a comprehensive guide to understanding, developing, and contributing to the ScopeAPI project.

> **📝 Note**: Documentation has been consolidated for better organization. The [Deployment Guide](DEPLOYMENT.md) now includes environment strategy, security best practices, and Kubernetes migration in one comprehensive file.

## 🚀 **Quick Start**

- **[Project Overview](../README.md)** - What is ScopeAPI and why it matters
- **[Development Guide](DEVELOPMENT.md)** - Complete development workflow
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment and operations

## 📋 **Documentation Sections**

### **🏗️ Architecture & Design**
- **[Architecture Overview](ARCHITECTURE.md)** - System design and technical architecture
- **[Project Structure](ARCHITECTURE.md#project-structure)** - Codebase organization
- **[Data Flow](ARCHITECTURE.md#data-flow)** - How data moves through the system

### **💻 Development**
- **[Development Guide](DEVELOPMENT.md)** - Complete development setup, workflows, and script usage
- **[API Documentation](API.md)** - Service APIs and endpoints
- **[Testing Guide](DEVELOPMENT.md#testing)** - How to test your changes

### **🖥️ Frontend & UI**
- **[Admin Console](ADMIN_CONSOLE.md)** - Angular frontend application guide
- **[UI Components](ADMIN_CONSOLE.md#key-features)** - Available UI components and features
- **[Frontend Development](ADMIN_CONSOLE.md#development-commands)** - Frontend development workflow

### **🚀 Deployment & Operations**
- **[Production Deployment](DEPLOYMENT.md)** - Deploy to production (includes environment strategy, security, and Kubernetes migration)
- **[Monitoring & Logging](DEPLOYMENT.md#monitoring)** - Observability

### **🤝 Contributing**
- **[Contribution Guide](CONTRIBUTING.md)** - How to contribute to ScopeAPI
- **[Code Standards](CONTRIBUTING.md#code-standards)** - Coding conventions
- **[Pull Request Process](CONTRIBUTING.md#pull-requests)** - Submitting changes

## 🎯 **What is ScopeAPI?**

ScopeAPI is a comprehensive API security and management platform that provides:

- **🔍 API Discovery** - Automatically discover and catalog APIs
- **🛡️ Threat Detection** - Identify security threats and vulnerabilities
- **🔒 Data Protection** - Protect sensitive data and ensure compliance
- **⚡ Attack Blocking** - Real-time threat prevention and blocking
- **🌐 Gateway Integration** - Integrate with popular API gateways
- **📊 Admin Console** - Centralized management and monitoring

## 🏗️ **Architecture Overview**

ScopeAPI follows a **microservices architecture** with:

- **7 Core Microservices** - Each handling a specific security domain
- **Event-Driven Communication** - Kafka-based message queuing
- **Polyglot Persistence** - PostgreSQL, Redis, Elasticsearch
- **Containerized Deployment** - Docker and Kubernetes
- **RESTful APIs** - Standard HTTP interfaces for all services

## 🚀 **Getting Started**

### **Prerequisites**
- Docker and Docker Compose (for local development)
- Kubernetes cluster (for staging/production)
- Go 1.21+ (for backend development)
- Node.js 18+ (for admin console)

### **Quick Start**
```bash
# Clone the repository
git clone https://github.com/your-org/scopeapi.git
cd scopeapi

# Complete setup
cd scripts
./scopeapi.sh setup --full

# Start development
./dev.sh start all
```

## 🔧 **Development Workflows**

### **Daily Development**
```bash
# Start services
cd scripts
./dev.sh start all

# Make changes and test
# View logs
./dev.sh logs api-discovery

# Stop when done
./dev.sh stop
```

### **Debugging**
```bash
# Start in debug mode
cd scripts
./dev.sh debug api-discovery

# Connect IDE to localhost:2345
# Set breakpoints and debug
```

## 📚 **Detailed Documentation**

### **For Developers**
- **[Development Guide](DEVELOPMENT.md)** - Complete development workflow
- **[API Reference](API.md)** - All service APIs and endpoints
- **[Testing Guide](DEVELOPMENT.md#testing)** - Testing strategies

### **For DevOps/Operations**
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment (includes environment strategy, security, and Kubernetes migration)
- **[Monitoring](DEPLOYMENT.md#monitoring)** - Observability and alerting

### **For Contributors**
- **[Contribution Guide](CONTRIBUTING.md)** - How to contribute
- **[Code Standards](CONTRIBUTING.md#code-standards)** - Coding conventions
- **[Architecture Decisions](ARCHITECTURE.md)** - Design rationale

## 🔍 **Need Help?**

- **📖 Documentation Issues**: Open an issue in the docs repository
- **🐛 Bug Reports**: Use the main repository issue tracker
- **💡 Feature Requests**: Submit through the main repository
- **❓ Questions**: Check existing issues or create a new one

## 📖 **Documentation Structure**

```
docs/
├── README.md                    # This file - Documentation index
├── ARCHITECTURE.md             # Technical architecture and design
├── DEVELOPMENT.md              # Development setup, workflows, and Docker setup
├── DEPLOYMENT.md               # Production deployment, environment strategy, security, and Kubernetes migration
├── API.md                      # API documentation and examples
├── ADMIN_CONSOLE.md            # Frontend application guide
└── CONTRIBUTING.md             # Contribution guidelines

k8s/
├── README.md                   # Kubernetes deployment configurations
├── deployments/                # All microservices + admin console
├── services/                   # Network services configuration
├── ingress/                    # Traffic routing and load balancing
├── secrets/                    # Environment variables and secrets
└── configmaps/                 # Application configuration
```

---

**🎯 This documentation is designed to help you:**
- **Understand** the ScopeAPI architecture and design
- **Develop** new features and improvements
- **Deploy** and operate ScopeAPI in production
- **Contribute** to the open-source project

**Happy coding! 🚀✨**
