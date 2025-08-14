# 🐳 ScopeAPI Docker Setup

## 📁 **Three-File Architecture**

This project uses a **clean three-file approach** for Docker Compose:

### **1. `scripts/docker-compose.yml` (Main)**
- **Purpose**: Production-ready orchestration
- **Content**: Infrastructure + microservices with basic settings
- **Usage**: `docker-compose -f scripts/docker-compose.yml up` (production-like)

### **2. `scripts/docker-compose.override.yml` (Development)**
- **Purpose**: Development environment enhancements
- **Content**: Service ports, development environment variables, source mounting
- **Usage**: `docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.override.yml up` (automatically loaded for development)

### **3. `scripts/docker-compose.debug.yml` (Debug)**
- **Purpose**: Debugging capabilities only
- **Content**: Debug ports, Delve debugger, debug Dockerfiles
- **Usage**: `docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.debug.yml up`

## 🚀 **Quick Start**

### **Complete Setup (Recommended for First Time)**
```bash
# Complete setup with validation
cd scripts
./scopeapi-setup.sh --full

# Or step by step:
./scopeapi-setup.sh --infrastructure  # Start infrastructure
./scopeapi-setup.sh --database        # Setup database
./scopeapi-setup.sh --validate        # Validate setup
```

### **Development Mode (After Setup)**
```bash
# Start infrastructure + services with development overrides
cd scripts
./scopeapi-services.sh start api-discovery

# Or use docker-compose directly (automatically loads override)
docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.override.yml up api-discovery
```

### **Debug Mode**
```bash
# Start services in debug mode
./scripts/./scripts/scopeapi-debug.sh start api-discovery

# Or use docker-compose with debug config
docker-compose -f scripts/docker-compose.yml -f scripts/docker-compose.debug.yml up api-discovery
```

### **Production Mode**
```bash
# Start services without development overrides
docker-compose -f scripts/docker-compose.yml up api-discovery
```

## 🔍 **What Each File Provides**

### **`scripts/docker-compose.yml`**
- ✅ PostgreSQL, Kafka, Redis, Elasticsearch, Kibana
- ✅ All 7 microservices with basic configuration
- ✅ Health checks and dependencies
- ✅ Production-ready settings

### **`scripts/docker-compose.override.yml`**
- 🔧 **Service ports** (8080-8086) for local access
- 🌱 **Development environment** variables
- 📁 **Source code mounting** for live development
- 🚀 **Development-friendly** configurations

### **`scripts/docker-compose.debug.yml`**
- 🐛 **Debug ports** (2345-2351) for Delve debugger
- 🔍 **Debug environment** variables
- 🚀 **Debug Dockerfiles** with Delve included
- 🔍 **Interactive debugging** support

## 🛠️ **Development Workflows**

### **Workflow 1: Complete Setup (First Time)**
```bash
# Complete setup with validation
cd scripts
./scopeapi-setup.sh --full

# This will:
# 1. Start infrastructure services
# 2. Setup PostgreSQL database
# 3. Run migrations
# 4. Create test data
# 5. Validate everything
```

### **Workflow 2: Services Management (Daily Work)**
```bash
# Start everything with development overrides
cd scripts
./scopeapi-services.sh start all

# Make code changes
# Services run with development settings
```

### **Workflow 3: Debug Mode**
```bash
# Start in debug mode
cd scripts
./scopeapi-debug.sh start api-discovery

# Connect IDE to localhost:2345
# Set breakpoints and debug
```

### **Workflow 4: Production Mode**
```bash
# Start without development overrides
docker-compose -f docker-compose.yml up api-discovery

# Services run with production settings
```

### **Workflow 5: Infrastructure Management**
```bash
# Start only infrastructure
cd scripts
./scopeapi-setup.sh --infrastructure

# Setup database separately
./scopeapi-setup.sh --database

# Validate setup
./scopeapi-setup.sh --validate
```

## 📋 **Available Commands**

### **Setup & Infrastructure (`cd scripts
./scopeapi-setup.sh`)**
```bash
cd scripts
./scopeapi-setup.sh --full               # Complete setup with validation
./scopeapi-setup.sh --infrastructure     # Start infrastructure only
./scopeapi-setup.sh --database           # Setup database only
./scopeapi-setup.sh --test-data          # Setup + create test data
./scopeapi-setup.sh --validate           # Setup + run validation tests
```

### **Services Management (`cd scripts
./scopeapi-services.sh`)**
```bash
cd scripts
./scopeapi-services.sh start [services...]     # Start services with dev overrides
./scopeapi-services.sh stop                    # Stop all services
./scopeapi-services.sh logs [service]          # View logs
./scopeapi-services.sh status                  # Check status
./scopeapi-services.sh build [services...]     # Build services
```

### **Debug Mode (`cd scripts
./scopeapi-debug.sh`)**
```bash
cd scripts
./scopeapi-debug.sh start [services...]   # Start in debug mode
./scopeapi-debug.sh stop                  # Stop debug services
./scopeapi-debug.sh logs [service]        # View debug logs
./scopeapi-debug.sh status                # Check debug status
./scopeapi-debug.sh build [services...]   # Build debug images
```

## 🔧 **Service Ports**

| Service | Normal Port | Debug Port | Purpose |
|---------|-------------|------------|---------|
| API Discovery | 8080 | 2345 | Service + Delve debugger |
| Gateway Integration | 8081 | 2346 | Service + Delve debugger |
| Data Ingestion | 8082 | 2347 | Service + Delve debugger |
| Threat Detection | 8083 | 2348 | Service + Delve debugger |
| Data Protection | 8084 | 2349 | Service + Delve debugger |
| Attack Blocking | 8085 | 2350 | Service + Delve debugger |
| Admin Console | 8086 | 2351 | Service + Delve debugger |

## 🎯 **Benefits of This Approach**

1. **✅ Clean separation** - Production vs development vs debug configuration
2. **✅ Automatic loading** - Development overrides load automatically
3. **✅ Easy switching** - Development mode vs debug mode vs production mode
4. **✅ Team flexibility** - Developers choose their workflow
5. **✅ Production ready** - Main file stays clean
6. **✅ Debug support** - Full debugging capabilities when needed

## 🚀 **IDE Integration**

### **VS Code**
1. Install Go extension
2. Use launch configurations in `.vscode/launch.json`
3. Connect to appropriate debug port (2345-2351)

### **GoLand**
1. Create Go Remote configuration
2. Set host: localhost, port: 2345-2351
3. Start debugging

## 🐛 **Debugging Guide**

### **Quick Debug Start:**
```bash
# Start service in debug mode
./scripts/./scripts/scopeapi-debug.sh start api-discovery

# Connect IDE to localhost:2345
# Set breakpoints and debug
```

### **Debug Ports by Service:**
| Service | Debug Port | IDE Configuration |
|---------|------------|-------------------|
| API Discovery | 2345 | localhost:2345 |
| Gateway Integration | 2346 | localhost:2346 |
| Data Ingestion | 2347 | localhost:2347 |
| Threat Detection | 2348 | localhost:2348 |
| Data Protection | 2349 | localhost:2349 |
| Attack Blocking | 2350 | localhost:2350 |
| Admin Console | 2351 | localhost:2351 |

### **Debug Workflow:**
1. **Start debug mode**: `cd scripts
./scopeapi-debug.sh start [service]`
2. **Connect IDE**: Use appropriate debug port
3. **Set breakpoints**: In your Go code
4. **Debug**: Step through, inspect variables
5. **Stop debug**: `./scopeapi-debug.sh stop`

### **Common Debug Issues:**
- **Port already in use**: Check if another debug session is running
- **Can't connect**: Ensure service is started in debug mode
- **Breakpoints not hit**: Verify source code mounting in debug Dockerfile

## 🔍 **Troubleshooting**

### **Service Won't Start**
```bash
# Check if infrastructure is ready
./scripts/dev.sh status

# Check service logs
./scripts/dev.sh logs [service]
```

### **Debug Port Not Working**
```bash
# Check if debug service is running
./scripts/debug.sh status

# Check debug logs
./scripts/debug.sh logs [service]
```

### **Port Already in Use**
```bash
# Check what's using the port
sudo lsof -i :8080

# Kill the process or use different port
```

---

**🎯 This setup gives you the best of all worlds:**
- **Clean production configuration** in main file
- **Automatic development overrides** for daily development
- **Full debugging capabilities** when you need them
- **Easy switching** between development, debug, and production modes
- **Team flexibility** for different development preferences
- **Clear script naming** that indicates purpose and scope

## 🏗️ **Script Architecture: Clear Separation of Concerns**

This project provides three specialized scripts for different use cases:

### **🔄 Local Development Management (`scopeapi-local.sh`)**
- **Purpose**: Manages Go microservices as direct processes
- **Use Case**: Development with direct binary execution
- **Dependencies**: Requires infrastructure to be running
- **Infrastructure Check**: Automatically verifies PostgreSQL, Kafka, Redis are running

### **🐳 Infrastructure Management ()**
- **Purpose**: Manages infrastructure services only
- **Use Case**: Infrastructure setup, troubleshooting, environment management
- **Services**: ZooKeeper, Kafka, PostgreSQL, Redis, Elasticsearch, Kibana
- **Features**: Environment setup, permission fixing, individual service management

### **🚀 Complete Microservices Orchestration ()**
- **Purpose**: Complete container-based microservices management
- **Use Case**: Production deployment, container-based development
- **Features**: Infrastructure + microservices, container management, debugging
- **New Commands**: `infrastructure`, `comprehensive-status`

### **🔧 Setup & Validation ()**
- **Purpose**: Complete project setup and validation
- **Use Case**: First-time setup, environment validation
- **Features**: Infrastructure startup, database setup, test data, validation

### **🐛 Debug Management ()**
- **Purpose**: Debug mode for microservices
- **Use Case**: Development debugging with Delve
- **Features**: Debug ports, container debugging, IDE integration

## 🎯 **When to Use Each Script:**

| Use Case | Primary Script | Secondary Script |
|----------|----------------|------------------|
| **First-time setup** | `cd scripts
./scopeapi-setup.sh` | - |
| **Daily development (containers)** | `cd scripts
./scopeapi-services.sh` | - |
| **Local development (processes)** | `cd scripts
./scopeapi-local.sh` | `cd scripts
./docker-infrastructure.sh` |
| **Infrastructure troubleshooting** | `cd scripts
./docker-infrastructure.sh` | - |
| **Debugging** | `cd scripts
./scopeapi-debug.sh` | - |
| **Production deployment** | `cd scripts
./scopeapi-services.sh` | - |
