# ScopeAPI Docker Setup

## Overview

This project uses a **hybrid approach** that follows industry best practices:

- **Individual Dockerfiles** for each microservice (for building and CI/CD)
- **Extended docker-compose.yml** for local development and testing
- **Development script** for easy service management

## 🏗️ **Architecture**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Docker Setup                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  Individual     │  │  docker-compose │  │  Development    │             │
│  │  Dockerfiles    │  │  .yml           │  │  Script         │             │
│  │                 │  │                 │  │                 │             │
│  │ • api-discovery │  │ • Infrastructure│  │ • Easy service  │             │
│  │ • gateway-int   │  │ • Microservices │  │   management    │             │
│  │ • data-ingest   │  │ • Networking    │  │ • Start/stop    │             │
│  │ • etc...        │  │ • Volumes       │  │ • Logs/status   │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 🚀 **Quick Start**

### **1. Start Infrastructure Only**
```bash
# Start databases, Kafka, Redis, etc.
./scripts/dev.sh start infrastructure
```

### **2. Start Specific Services**
```bash
# Start infrastructure + API Discovery
./scripts/dev.sh start api-discovery

# Start infrastructure + multiple services
./scripts/dev.sh start api-discovery gateway-integration

# Start everything
./scripts/dev.sh start all
```

### **3. Check Status**
```bash
# View all services
./scripts/dev.sh status

# View logs
./scripts/dev.sh logs api-discovery
```

### **4. Stop Everything**
```bash
./scripts/dev.sh stop
```

## 🛠️ **Development Workflow**

### **Option 1: Full Stack (Recommended)**
```bash
# Start everything for full development
./scripts/dev.sh start all

# Make code changes
# Services auto-reload (if using go run)

# View logs
./scripts/dev.sh logs api-discovery

# Restart specific service
./scripts/dev.sh restart api-discovery
```

### **Option 2: Infrastructure + Local Services**
```bash
# Start only infrastructure
./scripts/dev.sh start infrastructure

# Run services locally with go run
cd services/api-discovery
go run cmd/main.go

# In another terminal
cd services/gateway-integration
go run cmd/main.go
```

### **Option 3: Individual Service Development**
```bash
# Start infrastructure
./scripts/dev.sh start infrastructure

# Build and run specific service
./scripts/dev.sh build api-discovery
./scripts/dev.sh start api-discovery
```

## 📋 **Available Commands**

```bash
./scripts/dev.sh help                    # Show all commands
./scripts/dev.sh start [services...]     # Start services
./scripts/dev.sh stop                    # Stop all services
./scripts/dev.sh restart [services...]   # Restart services
./scripts/dev.sh logs [service]          # Show logs
./scripts/dev.sh status                  # Show service status
./scripts/dev.sh build [services...]     # Build services
./scripts/dev.sh clean                   # Clean everything
```

## 🔧 **Service Configuration**

### **Port Mapping**
- **API Discovery**: `localhost:8080`
- **Gateway Integration**: `localhost:8081`
- **Data Ingestion**: `localhost:8082`
- **Threat Detection**: `localhost:8083`
- **Data Protection**: `localhost:8084`
- **Attack Blocking**: `localhost:8085`
- **Admin Console**: `localhost:8086`

### **Environment Variables**
All services automatically get:
- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=scopeapi`
- `DB_PASSWORD=${POSTGRES_PASSWORD}`
- `DB_NAME=scopeapi`
- `KAFKA_BROKERS=kafka:9092`

## 🐳 **Docker Commands (Alternative)**

If you prefer using docker-compose directly:

```bash
# Start everything
docker-compose up -d

# Start specific services
docker-compose up -d api-discovery gateway-integration

# View logs
docker-compose logs -f api-discovery

# Stop everything
docker-compose down

# Rebuild services
docker-compose build api-discovery
```

## 🔍 **Troubleshooting**

### **Service Won't Start**
```bash
# Check service status
./scripts/dev.sh status

# View service logs
./scripts/dev.sh logs api-discovery

# Check if infrastructure is ready
docker-compose ps
```

### **Port Already in Use**
```bash
# Check what's using the port
sudo lsof -i :8080

# Kill the process or change port in docker-compose.yml
```

### **Database Connection Issues**
```bash
# Wait for PostgreSQL to be ready
docker-compose logs postgres

# Check health status
docker-compose ps postgres
```

## 🚀 **Production Deployment**

For production, you can still use the individual Dockerfiles:

```bash
# Build production images
docker build -t scopeapi/api-discovery:latest ./services/api-discovery
docker build -t scopeapi/gateway-integration:latest ./services/gateway-integration

# Push to registry
docker push scopeapi/api-discovery:latest
docker push scopeapi/gateway-integration:latest
```

## 📚 **Best Practices**

1. **Always start infrastructure first** - Services depend on databases and message queues
2. **Use health checks** - Services wait for dependencies to be healthy
3. **Environment variables** - Use `.env` file for sensitive data
4. **Service isolation** - Each service runs in its own container
5. **Network isolation** - All services use the same network for communication

## 🔗 **Related Files**

- `docker-compose.yml` - Main orchestration file
- `docker-compose.override.yml` - Development overrides (optional)
- `scripts/dev.sh` - Development management script
- `services/*/Dockerfile` - Individual service Dockerfiles
- `.env` - Environment variables (create from `.env.example`)

---

**🎯 This setup gives you the best of both worlds:**
- **Easy development** with docker-compose
- **Production flexibility** with individual Dockerfiles
- **Industry standard** approach used by Netflix, Uber, Spotify
