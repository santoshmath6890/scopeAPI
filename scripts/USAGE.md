# 🛠️ ScopeAPI Scripts & Configuration

This directory contains all the scripts, Docker Compose files, and configuration files for managing ScopeAPI development, debugging, and setup.

## 📋 **Script Overview**

### **1. `scopeapi-setup.sh` - Complete Setup Orchestrator**
- **Purpose**: Orchestrates the complete ScopeAPI setup process
- **When to use**: First time setup, infrastructure management, database setup
- **Features**: 
  - Infrastructure startup (ZooKeeper, Kafka, PostgreSQL, Redis, Elasticsearch, Kibana)
  - Database setup and migrations
  - Test data creation
  - Setup validation

**Usage Examples:**
```bash
# Complete setup with validation (recommended for first time)
./scripts/scopeapi-setup.sh --full

# Start infrastructure only
./scripts/scopeapi-setup.sh --infrastructure

# Setup database only
./scripts/scopeapi-setup.sh --database

# Create test data
./scripts/scopeapi-setup.sh --test-data

# Validate existing setup
./scripts/scopeapi-setup.sh --validate
```

---

### **2. `scopeapi-services.sh` - Services Lifecycle Manager**
- **Purpose**: Manages microservices lifecycle and orchestration
- **When to use**: Daily development work, starting/stopping services, service management
- **Features**:
  - Development ports (8080-8086)
  - Source code mounting for live development
  - Development environment variables
  - Service management (start/stop/restart/logs/status)
  - Container access (shell, exec commands)

**Usage Examples:**
```bash
# Start all services
./scripts/scopeapi-services.sh start all

# Start specific service
./scripts/scopeapi-services.sh start api-discovery

# Start multiple services
./scripts/scopeapi-services.sh start api-discovery gateway-integration

# View logs
./scripts/scopeapi-services.sh logs api-discovery

# Check status
./scripts/scopeapi-services.sh status

# Stop all services
./scripts/scopeapi-services.sh stop
```

---

### **3. `scopeapi-debug.sh` - Debug Mode Manager**
- **Purpose**: Manages microservices in debug mode with Delve debugger
- **When to use**: When you need to debug services, set breakpoints
- **Features**:
  - Debug ports (2345-2351) for Delve debugger
  - Debug Dockerfiles with Delve included
  - Interactive debugging support
  - Debug environment variables

**Usage Examples:**
```bash
# Start service in debug mode
./scripts/scopeapi-debug.sh start api-discovery

# Start multiple services in debug mode
./scripts/scopeapi-debug.sh start api-discovery gateway-integration

# View debug logs
./scripts/scopeapi-debug.sh logs api-discovery

# Check debug status
./scripts/scopeapi-debug.sh status

# Stop debug services
./scripts/scopeapi-debug.sh stop
```

---

## 🔄 **Workflow Progression**

### **First Time Setup:**
```bash
# 1. Complete setup
./scripts/scopeapi-setup.sh --full

# 2. Verify everything is working
./scripts/scopeapi-setup.sh --validate
```

### **Daily Development:**
```bash
# 1. Start services for development
./scripts/scopeapi-services.sh start all

# 2. Make code changes
# 3. View logs if needed
./scripts/scopeapi-services.sh logs api-discovery

# 4. Stop when done
./scripts/scopeapi-services.sh stop
```

### **When Debugging is Needed:**
```bash
# 1. Start service in debug mode
./scripts/scopeapi-debug.sh start api-discovery

# 2. Connect IDE to localhost:2345
# 3. Set breakpoints and debug
# 4. Stop debug service when done
./scripts/scopeapi-debug.sh stop
```

---

## 🎯 **Script Naming Convention**

All scripts follow the `scopeapi-{purpose}.sh` naming convention:

- **`scopeapi-setup.sh`** - Setup and infrastructure management
- **`scopeapi-services.sh`** - Services lifecycle management  
- **`scopeapi-debug.sh`** - Debug mode management

This makes it clear that these are ScopeAPI-specific scripts and what each one does.

---

## 🚨 **Important Notes**

1. **Always run setup first**: Use `scopeapi-setup.sh` before using other scripts
2. **Development vs Debug**: Use `scopeapi-dev.sh` for normal development, `scopeapi-debug.sh` only when debugging
3. **Infrastructure dependency**: Development and debug scripts require infrastructure to be running
4. **Port conflicts**: Development and debug scripts use different ports, so they won't conflict

---

## 🔧 **Troubleshooting**

### **Script not found:**
```bash
# Make sure you're in the scripts directory
cd scripts

# Make scripts executable
chmod +x scopeapi-*.sh
```

### **Permission denied:**
```bash
# Fix permissions
chmod +x scopeapi-*.sh
```

### **Services won't start:**
```bash
# Check if infrastructure is running
./scripts/scopeapi-setup.sh --validate

# Start infrastructure if needed
./scripts/scopeapi-setup.sh --infrastructure
```

---

**These scripts provide a complete toolkit for ScopeAPI development!** 🛠️✨

## 🐳 **Docker Compose Files**

This directory also contains all the Docker Compose configuration files:

### **`docker-compose.yml`**
- **Purpose**: Main production-ready orchestration
- **Content**: Infrastructure + all 7 microservices
- **Usage**: `docker-compose -f scripts/docker-compose.yml up`

### **`docker-compose.override.yml`**
- **Purpose**: Development environment overrides
- **Content**: Development ports, source mounting, dev variables
- **Usage**: Automatically loaded with main compose file

### **`docker-compose.debug.yml`**
- **Purpose**: Debug mode configuration
- **Content**: Debug ports, Delve debugger, debug Dockerfiles
- **Usage**: `docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.debug.yml up`

## 🔧 **Other Configuration Files**

### **`docker-infrastructure.sh`**
- **Purpose**: Infrastructure management script
- **Content**: Docker container management utilities
- **Usage**: `./scripts/docker-infrastructure.sh [command]`

### **`scopeapi-local.sh`**
- **Purpose**: Legacy service management script
- **Content**: Alternative service management approach
- **Usage**: `./scripts/scopeapi-local.sh [command]`


## 🏗️ **New Architecture: Clear Separation of Concerns**

This directory now provides a clean separation of responsibilities:

### **🔄 Local Development Management (`scopeapi-local.sh`)**
- **Purpose**: Manages Go microservices as direct processes
- **Use Case**: Development with direct binary execution
- **Dependencies**: Requires infrastructure to be running
- **Command**: `./scripts/scopeapi-local.sh [start|stop|status|build]`

### **🐳 Infrastructure Management ()**
- **Purpose**: Manages infrastructure services only
- **Use Case**: Infrastructure setup, troubleshooting, environment management
- **Services**: PostgreSQL, Kafka, Redis, Elasticsearch, Kibana
- **Command**: `./scripts/docker-infrastructure.sh [start|stop|status|logs|setup-env]`

### **🚀 Complete Microservices Orchestration ()**
- **Purpose**: Complete container-based microservices management
- **Use Case**: Production deployment, container-based development
- **Features**: Infrastructure + microservices, container management, debugging
- **Command**: `./scripts/scopeapi-services.sh [start|stop|infrastructure|comprehensive-status]`

### **🔧 Setup & Validation ()**
- **Purpose**: Complete project setup and validation
- **Use Case**: First-time setup, environment validation
- **Features**: Infrastructure startup, database setup, test data, validation
- **Command**: `./scripts/scopeapi-setup.sh [--full|--infrastructure|--database|--validate]`

### **🐛 Debug Management ()**
- **Purpose**: Debug mode for microservices
- **Use Case**: Development debugging with Delve
- **Features**: Debug ports, container debugging, IDE integration
- **Command**: `./scripts/scopeapi-debug.sh [start|stop|logs|status]`


## 🔄 **Workflow Examples**

### **Option 1: Local Development (Process-Based)**
```bash
# 1. Start infrastructure
./scripts/docker-infrastructure.sh start

# 2. Start microservices as processes
./scripts/scopeapi-local.sh start

# 3. Check status
./scripts/scopeapi-local.sh status
```

### **Option 2: Container-Based Development**
```bash
# 1. Start everything with containers
./scripts/scopeapi-services.sh start all

# 2. Check comprehensive status
./scripts/scopeapi-services.sh comprehensive-status

# 3. Access container shell
./scripts/scopeapi-services.sh shell api-discovery
```

### **Option 3: Infrastructure Only**
```bash
# 1. Start only infrastructure
./scripts/docker-infrastructure.sh start

# 2. Or use the integrated approach
./scripts/scopeapi-services.sh infrastructure

# 3. Check infrastructure status
./scripts/docker-infrastructure.sh status
```

### **Option 4: Complete Setup**
```bash
# 1. Complete setup with validation
./scripts/scopeapi-setup.sh --full

# 2. Start services as needed
./scripts/scopeapi-services.sh start api-discovery

# 3. Debug if needed
./scripts/scopeapi-debug.sh start api-discovery
```

## 🏗️ **New Architecture: Clear Separation of Concerns**

This directory now provides a clean separation of responsibilities:

### **🔄 Local Development Management (`scopeapi-local.sh`)**
- **Purpose**: Manages Go microservices as direct processes
- **Use Case**: Development with direct binary execution
- **Dependencies**: Requires infrastructure to be running
- **Command**: `./scripts/scopeapi-local.sh [start|stop|status|build]`

### **🐳 Infrastructure Management ()**
- **Purpose**: Manages infrastructure services only
- **Use Case**: Infrastructure setup, troubleshooting, environment management
- **Services**: PostgreSQL, Kafka, Redis, Elasticsearch, Kibana
- **Command**: `./scripts/docker-infrastructure.sh [start|stop|status|logs|setup-env]`

### **🚀 Complete Microservices Orchestration ()**
- **Purpose**: Complete container-based microservices management
- **Use Case**: Production deployment, container-based development
- **Features**: Infrastructure + microservices, container management, debugging
- **Command**: `./scripts/scopeapi-services.sh [start|stop|infrastructure|comprehensive-status]`

### **🔧 Setup & Validation ()**
- **Purpose**: Complete project setup and validation
- **Use Case**: First-time setup, environment validation
- **Features**: Infrastructure startup, database setup, test data, validation
- **Command**: `./scripts/scopeapi-setup.sh [--full|--infrastructure|--database|--validate]`

### **🐛 Debug Management ()**
- **Purpose**: Debug mode for microservices
- **Use Case**: Development debugging with Delve
- **Features**: Debug ports, container debugging, IDE integration
- **Command**: `./scripts/scopeapi-debug.sh [start|stop|logs|status]`

## 🔄 **Workflow Examples**

### **Option 1: Local Development (Process-Based)**
```bash
# 1. Start infrastructure
./scripts/docker-infrastructure.sh start

# 2. Start microservices as processes
./scripts/scopeapi-local.sh start

# 3. Check status
./scripts/scopeapi-local.sh status
```

### **Option 2: Container-Based Development**
```bash
# 1. Start everything with containers
./scripts/scopeapi-services.sh start all

# 2. Check comprehensive status
./scripts/scopeapi-services.sh comprehensive-status

# 3. Access container shell
./scripts/scopeapi-services.sh shell api-discovery
```

### **Option 3: Infrastructure Only**
```bash
# 1. Start only infrastructure
./scripts/docker-infrastructure.sh start

# 2. Or use the integrated approach
./scripts/scopeapi-services.sh infrastructure

# 3. Check infrastructure status
./scripts/docker-infrastructure.sh status
```

### **Option 4: Complete Setup**
```bash
# 1. Complete setup with validation
./scripts/scopeapi-setup.sh --full

# 2. Start services as needed
./scripts/scopeapi-services.sh start api-discovery

# 3. Debug if needed
./scripts/scopeapi-debug.sh start api-discovery
```
