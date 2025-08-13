# 🛠️ ScopeAPI Script Reference Card

## 🏗️ **Script Architecture Overview**

This project provides **three specialized scripts** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────┐
│                    ScopeAPI Scripts                            │
├─────────────────────────────────────────────────────────────────┤
│  🔄 Local Development │  🐳 Infrastructure  │  🚀 Container    │
│  Management          │  Management         │  Orchestration   │
│  scopeapi-local.sh │  docker-infrastructure.sh │ scopeapi-services.sh │
└─────────────────────────────────────────────────────────────────┘
```

## 📋 **Script Comparison Matrix**

| Feature | `scopeapi-local.sh` | `docker-infrastructure.sh` | `scopeapi-services.sh` |
|---------|----------------------|---------------------------|------------------------|
| **Purpose** | Local development (processes) | Infrastructure only | Complete orchestration |
| **Use Case** | Go development, debugging | Infrastructure setup | Production, containers |
| **Dependencies** | Requires infrastructure | None | None (starts everything) |
| **Startup Speed** | Fast (direct binaries) | Medium (containers) | Medium (containers) |
| **Resource Usage** | Low | Medium | Medium |
| **Debugging** | Easy (direct access) | N/A | Container-based |
| **Scaling** | Manual | N/A | Docker Compose |

## 🎯 **When to Use Each Script**

### **🔄 Use `scopeapi-local.sh` when:**
- ✅ **Developing Go services** - Fast iteration with direct binary execution
- ✅ **Debugging issues** - Direct process access and control
- ✅ **Performance testing** - Lower overhead for benchmarking
- ✅ **Infrastructure is already running** - Quick service management

### **🐳 Use `docker-infrastructure.sh` when:**
- ✅ **Setting up environment** - First-time infrastructure setup
- ✅ **Troubleshooting infrastructure** - Individual service management
- ✅ **Fixing permissions** - Docker permission issues
- ✅ **Environment management** - `.env` file setup and configuration

### **🚀 Use `scopeapi-services.sh` when:**
- ✅ **Production deployment** - Complete container orchestration
- ✅ **Container-based development** - Full-stack testing
- ✅ **Service orchestration** - Managing multiple services together
- ✅ **Infrastructure + microservices** - Complete system management

## 🔄 **Workflow Examples**

### **Workflow 1: Local Development (Go Services)**
```bash
# 1. Start infrastructure (if not running)
./scripts/docker-infrastructure.sh start

# 2. Start Go services as processes
./scripts/scopeapi-local.sh start

# 3. Develop and test
# 4. Check status
./scripts/scopeapi-local.sh status

# 5. Stop when done
./scripts/scopeapi-local.sh stop
```

### **Workflow 2: Container-Based Development**
```bash
# 1. Start everything with containers
./scripts/scopeapi-services.sh start all

# 2. Develop and test
# 3. Check comprehensive status
./scripts/scopeapi-services.sh comprehensive-status

# 4. Access container if needed
./scripts/scopeapi-services.sh shell api-discovery

# 5. Stop when done
./scripts/scopeapi-services.sh stop
```

### **Workflow 3: Infrastructure Management**
```bash
# 1. Check infrastructure status
./scripts/docker-infrastructure.sh status

# 2. Start infrastructure only
./scripts/docker-infrastructure.sh start

# 3. View specific service logs
./scripts/docker-infrastructure.sh logs kafka

# 4. Fix permissions if needed
./scripts/docker-infrastructure.sh fix-permissions
```

### **Workflow 4: Complete Setup**
```bash
# 1. Complete setup with validation
./scripts/scopeapi-setup.sh --full

# 2. Start specific services
./scripts/scopeapi-services.sh start api-discovery

# 3. Debug if needed
./scripts/scopeapi-debug.sh start api-discovery
```

## 🚨 **Important Notes**

### **⚠️ Infrastructure Dependencies**
- **`scopeapi-local.sh`** requires infrastructure to be running
- **`scopeapi-services.sh`** starts infrastructure automatically
- **`docker-infrastructure.sh`** manages infrastructure only

### **🔄 No Conflicts**
- Each script has a specific purpose
- No overlap in functionality
- Can be used together or separately

### **🎯 Best Practices**
- **Development**: Use `scopeapi-local.sh` for Go development
- **Testing**: Use `scopeapi-services.sh` for integration testing
- **Troubleshooting**: Use `docker-infrastructure.sh` for infrastructure issues
- **Production**: Use `scopeapi-services.sh` for deployment

## 📚 **Related Documentation**

- **[📖 Complete Documentation](README.md)** - Full documentation index
- **[💻 Development Guide](DEVELOPMENT.md)** - Development workflows
- **[🐳 Docker Setup](DOCKER_SETUP.md)** - Container setup guide
- **[🛠️ Scripts Usage](../scripts/README.md)** - Detailed script documentation

---

**🎯 This architecture provides flexibility while maintaining clear separation of concerns!**
