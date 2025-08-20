# 🛠️ ScopeAPI Scripts & Configuration

This directory contains all the scripts, Docker Compose files, and configuration files for managing ScopeAPI development, debugging, and setup.

## 📋 **Script Overview**

### **1. `scopeapi.sh` - Main Orchestrator Script**
- **Purpose**: Unified script for all ScopeAPI operations
- **When to use**: Setup, main operations, service management
- **Features**: 
  - Complete setup orchestration (infrastructure + database + validation)
  - Service lifecycle management (start/stop/restart/logs/status)
  - Comprehensive status reporting
  - Cleanup operations

**Usage Examples:**
```bash
# Complete setup with validation (recommended for first time)
cd scripts
./scopeapi.sh setup --full

# Start infrastructure only
./scopeapi.sh setup --infrastructure

# Setup database only
./scopeapi.sh setup --database

# Create test data
./scopeapi.sh setup --test-data

# Validate existing setup
./scopeapi.sh setup --validate

# Start all services
./scopeapi.sh start all

# Check comprehensive status
./scopeapi.sh comprehensive-status

# Cleanup services
./scopeapi.sh setup --cleanup
```

---

### **2. `dev.sh` - Development Workflow Script**
- **Purpose**: Daily development tasks, debugging, testing, and development utilities
- **When to use**: Daily development work, debugging, testing, code quality
- **Features**:
  - Development environment management
  - Debug mode with Delve debugger
  - Testing (backend/frontend)
  - Code quality (linting, formatting)
  - Development utilities

**Usage Examples:**
```bash
# Start development environment
cd scripts
./dev.sh start all

# Start specific service
./dev.sh start api-discovery

# Debug service
./dev.sh debug api-discovery

# Run tests
./dev.sh test

# Lint code
./dev.sh lint

# Format code
./dev.sh format
```

---

### **3. `infrastructure.sh` - Infrastructure Management Script**
- **Purpose**: Manage infrastructure services (PostgreSQL, Kafka, Redis, etc.)
- **When to use**: Infrastructure setup, monitoring, troubleshooting
- **Features**:
  - Infrastructure service management
  - Health monitoring
  - Troubleshooting utilities
  - Resource management

**Usage Examples:**
```bash
# Start infrastructure services
cd scripts
./infrastructure.sh start

# Check infrastructure health
./infrastructure.sh health

# View infrastructure logs
./infrastructure.sh logs

# Stop infrastructure
./infrastructure.sh stop
```

---

### **4. `deploy.sh` - Deployment Script**
- **Purpose**: Unified deployment for Docker (local) and Kubernetes (staging/production)
- **When to use**: Environment deployment, secrets management
- **Features**:
  - Environment-specific deployment
  - Platform selection (Docker/Kubernetes)
  - Secrets management
  - Validation and checks

**Usage Examples:**
```bash
# Deploy to Docker (local development)
cd scripts
./deploy.sh -e dev -p docker

# Deploy to Kubernetes (staging)
./deploy.sh -e staging -p k8s

# Deploy to Kubernetes (production)
./deploy.sh -e prod -p k8s
```

---

## 🔄 **Migration from Old Scripts**

The following table shows how to migrate from the old scripts to the new consolidated ones:

| **Old Script** | **New Command** | **Notes** |
|----------------|-----------------|-----------|
| `./scopeapi-setup.sh --full` | `./scopeapi.sh setup --full` | Complete setup with validation |
| `./scopeapi-setup.sh --infrastructure` | `./scopeapi.sh setup --infrastructure` | Infrastructure only |
| `./scopeapi-setup.sh --database` | `./scopeapi.sh setup --database` | Database setup only |
| `./scopeapi-setup.sh --test-data` | `./scopeapi.sh setup --test-data` | Setup with test data |
| `./scopeapi-setup.sh --validate` | `./scopeapi.sh setup --validate` | Validate setup |
| `./scopeapi-setup.sh --cleanup` | `./scopeapi.sh setup --cleanup` | Cleanup services |
| `./scopeapi-setup.sh --cleanup-full` | `./scopeapi.sh setup --cleanup-full` | Full cleanup |
| `./scopeapi-services.sh start all` | `./scopeapi.sh start all` | Start all services |
| `./scopeapi-services.sh start [service]` | `./scopeapi.sh start [service]` | Start specific service |
| `./scopeapi-services.sh stop` | `./scopeapi.sh stop` | Stop all services |
| `./scopeapi-services.sh status` | `./scopeapi.sh status` | Show status |
| `./scopeapi-services.sh comprehensive-status` | `./scopeapi.sh comprehensive-status` | Detailed status |
| `./scopeapi-services.sh logs [service]` | `./scopeapi.sh logs [service]` | Show logs |
| `./scopeapi-debug.sh start [service]` | `./dev.sh debug [service]` | Debug mode |
| `./scopeapi-local.sh start [service]` | `./dev.sh start [service]` | Local development |
| `./docker-infrastructure.sh start` | `./infrastructure.sh start` | Infrastructure |
| `./deploy-all.sh -e dev -p docker` | `./deploy.sh -e dev -p docker` | Local deployment |
| `./deploy-k8s.sh` | `./deploy.sh -e staging -p k8s` | Kubernetes deployment |

---

## 🔄 **Workflow Progression**

### **First Time Setup:**
```bash
# 1. Complete setup
cd scripts
./scopeapi.sh setup --full

# 2. Verify everything is working
./scopeapi.sh setup --validate
```

### **Daily Development:**
```bash
# 1. Start services for development
cd scripts
./scopeapi.sh start all

# 2. Make code changes
# 3. View logs if needed
./scopeapi.sh logs api-discovery

# 4. Stop when done
./scopeapi.sh stop
```

### **When Debugging is Needed:**
```bash
# 1. Start service in debug mode
cd scripts
./dev.sh debug api-discovery

# 2. Connect IDE to localhost:2345
# 3. Set breakpoints and debug
# 4. Stop debug service when done
./dev.sh stop
```

---

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

---

## ☸️ **Kubernetes Configuration**

For Kubernetes deployment configurations, see the **[k8s/](../k8s/README.md)** directory at the project root, which contains:
- **Deployments** for all microservices and admin console
- **Services** for network communication
- **Ingress** for traffic routing and load balancing
- **Secrets** for secure environment variable management
- **ConfigMaps** for application configuration

---

**These scripts provide a complete toolkit for ScopeAPI development!** 🛠️✨
