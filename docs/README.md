# 📚 ScopeAPI Documentation

Welcome to the ScopeAPI documentation! This is a comprehensive guide to understanding, developing, and contributing to the ScopeAPI project.

## 🚀 **Quick Start**

- **[Project Overview](README.md)** - What is ScopeAPI and why it matters
- **[Docker Setup](DOCKER_SETUP.md)** - Get up and running quickly
- **[Development Guide](DEVELOPMENT.md)** - Start developing with ScopeAPI

## 📋 **Documentation Sections**

### **🏗️ Architecture & Design**
- **[Architecture Overview](ARCHITECTURE.md)** - System design and technical architecture
- **[Project Structure](ARCHITECTURE.md#project-structure)** - Codebase organization
- **[Data Flow](ARCHITECTURE.md#data-flow)** - How data moves through the system

### **💻 Development**

### **🖥️ Frontend & UI**
- **[Admin Console](ADMIN_CONSOLE.md)** - Angular frontend application guide
- **[UI Components](ADMIN_CONSOLE.md#key-features)** - Available UI components and features
- **[Frontend Development](ADMIN_CONSOLE.md#development-commands)** - Frontend development workflow
- **[Development Setup](DEVELOPMENT.md)** - Local development environment
- **[Scripts Guide](../scripts/README.md)** - Available development scripts
- **[API Documentation](API.md)** - Service APIs and endpoints
- **[Testing Guide](DEVELOPMENT.md#testing)** - How to test your changes

### **🚀 Deployment & Operations**
- **[Production Deployment](DEPLOYMENT.md)** - Deploy to production
- **[Docker Orchestration](DOCKER_SETUP.md)** - Container management
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
- **Containerized Deployment** - Docker and Docker Compose
- **RESTful APIs** - Standard HTTP interfaces for all services

## 🚀 **Getting Started**

### **Prerequisites**
- Docker and Docker Compose
- Go 1.21+ (for backend development)
- Node.js 18+ (for admin console)

### **Quick Start**
```bash
# Clone the repository
git clone https://github.com/your-org/scopeapi.git
cd scopeapi

# Complete setup
./scripts/scopeapi-setup.sh --full

# Start development
./scripts/scopeapi-services.sh start all
```

## 🔧 **Development Workflows**

### **Daily Development**
```bash
# Start services
./scripts/scopeapi-services.sh start all

# Make changes and test
# View logs
./scripts/scopeapi-services.sh logs api-discovery

# Stop when done
./scripts/scopeapi-services.sh stop
```

### **Debugging**
```bash
# Start in debug mode
./scripts/scopeapi-debug.sh start api-discovery

# Connect IDE to localhost:2345
# Set breakpoints and debug
```

## 📚 **Detailed Documentation**

### **For Developers**
- **[Development Guide](DEVELOPMENT.md)** - Complete development workflow
- **[API Reference](API.md)** - All service APIs and endpoints
- **[Testing Guide](DEVELOPMENT.md#testing)** - Testing strategies

### **For DevOps/Operations**
- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment
- **[Docker Setup](DOCKER_SETUP.md)** - Container orchestration
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
├── README.md               # This file - Documentation index
├── ARCHITECTURE.md        # Technical architecture and design
├── DEVELOPMENT.md         # Development setup and workflows
├── API.md                 # API documentation and examples
├── CONTRIBUTING.md        # Contribution guidelines
├── DEPLOYMENT.md          # Production deployment
# All documentation is now consolidated into focused, comprehensive guides
```

---

**🎯 This documentation is designed to help you:**
- **Understand** the ScopeAPI architecture and design
- **Develop** new features and improvements
- **Deploy** and operate ScopeAPI in production
- **Contribute** to the open-source project

**Happy coding! 🚀✨**

## 🏗️ **Script Architecture Overview**

This project provides a comprehensive set of specialized scripts:

### **🔄 Local Development Management**
- **`scopeapi-local.sh`** - Local development (process-based management)
- **`docker-infrastructure.sh`** - Infrastructure services only

### **🚀 Container-Based Management**
- **`scopeapi-services.sh`** - Complete microservices orchestration
- **`scopeapi-debug.sh`** - Debug mode management

### **🔧 Setup & Validation**
- **`scopeapi-setup.sh`** - Complete project setup

### **📖 Usage Guide**
- **`README.md`** - Comprehensive script documentation

## 🎯 **Quick Script Selection:**

- **First time**: `./scripts/scopeapi-setup.sh --full`
- **Daily development (containers)**: `./scripts/scopeapi-services.sh start all`
- **Local development (processes)**: `./scripts/scopeapi-local.sh start`
- **Infrastructure only**: `./scripts/docker-infrastructure.sh start`
- **Debugging**: `./scripts/scopeapi-debug.sh start [service]`
