# 🚀 ScopeAPI Deployment Guide

This guide covers production deployment, monitoring, and operations for the ScopeAPI platform.

## 📋 **Table of Contents**

- [Deployment Overview](#deployment-overview)
- [Environment Strategy & Security](#environment-strategy--security)
- [Local Development Setup](#local-development-setup)
- [Secrets Management](#secrets-management)
- [Deployment Commands](#deployment-commands)
- [Validation and Checks](#validation-and-checks)
- [Troubleshooting](#troubleshooting)
- [Kubernetes Migration](#kubernetes-migration)
- [Prerequisites](#prerequisites)
- [Production Deployment](#production-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Monitoring and Observability](#monitoring-and-observability)
- [Security Configuration](#security-configuration)
- [Backup and Recovery](#backup-and-recovery)
- [Scaling and Performance](#scaling-and-performance)
- [Maintenance](#maintenance)

## 🎯 **Deployment Overview**

ScopeAPI supports multiple deployment strategies:

- **Docker Compose** - Simple, single-server deployment
- **Kubernetes** - Production-grade, scalable deployment
- **Hybrid Cloud** - Multi-environment deployment
- **On-Premise** - Self-hosted deployment

## 🔒 **Environment Strategy & Security**

### **Local Development**
- **File**: `.env.local` (your machine only)
- **Purpose**: Local development and testing
- **Security**: Never commit to version control
- **Usage**: `./scripts/deploy.sh -e dev -p docker`

### **Staging Environment**
- **Method**: Kubernetes Secrets
- **Purpose**: Team testing and integration
- **Security**: Encrypted, RBAC controlled
- **Usage**: `./scripts/deploy.sh -e staging -p k8s`

### **Production Environment**
- **Method**: Kubernetes Secrets + External Secrets Manager
- **Purpose**: Live production deployment
- **Security**: Enterprise-grade encryption and access control
- **Usage**: `./scripts/deploy.sh -e prod -p k8s`

### **🚨 Security Rules**

#### **❌ NEVER DO:**
- Commit `.env.local` files to version control
- Use `.env` files for staging or production
- Store real passwords in plain text files
- Share environment files between team members

#### **✅ ALWAYS DO:**
- Use `.env.local` only for local development
- Use Kubernetes Secrets for staging/production
- Keep `.env.local` files local to your machine
- Use `env.example` as a template only

### **📁 File Usage Matrix**

| **File** | **Local Dev** | **Staging** | **Production** | **Commit?** |
|----------|---------------|-------------|----------------|-------------|
| `env.example` | ✅ Template | ✅ Template | ✅ Template | ✅ Yes |
| `.env.local` | ✅ Use | ❌ Never | ❌ Never | ❌ Never |
| `.env` | ❌ Deprecated | ❌ Never | ❌ Never | ❌ Never |
| `k8s/secrets.yaml` | ❌ Not used | ✅ Use | ✅ Use | ❌ Never |

## 🏠 **Local Development Setup**

### **1. Environment Configuration**
```bash
# Copy the template
cp env.example .env.local

# Edit with your local values
nano .env.local

# Key environment variables:
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_secure_password
KAFKA_BROKER_ID=1
```

### **2. Deploy Locally**
```bash
# Deploy to Docker (local development)
./scripts/deploy.sh -e dev -p docker

# This will:
# 1. Use .env.local file
# 2. Deploy to Docker Compose
# 3. Only work on your local machine
```

### **3. Local Development Commands**
```bash
# Start local development environment
./scripts/dev.sh start all

# Check status
./scripts/scopeapi.sh status

# View logs
./scripts/dev.sh logs api-discovery

# Stop when done
./scripts/dev.sh stop
```

## 🔐 **Secrets Management**

### **Local Development (.env.local)**
```bash
# Example .env.local content
POSTGRES_PASSWORD=my_local_password_123
REDIS_PASSWORD=my_local_redis_456
JWT_SECRET=my_local_jwt_secret_789
```

### **Staging/Production (Kubernetes Secrets)**
```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: scopeapi-secrets
  namespace: scopeapi
type: Opaque
data:
  POSTGRES_PASSWORD: <base64-encoded-staging-password>
  REDIS_PASSWORD: <base64-encoded-staging-redis-password>
  JWT_SECRET: <base64-encoded-staging-jwt-secret>
```

### **Generate Base64 Encoded Secrets**
```bash
# Run the secret generation script
./scripts/generate-secrets.sh

# This will output base64 encoded values like:
# POSTGRES_PASSWORD: eW91cl9zZWN1cmVfcG9zdGdyZXNfcGFzc3dvcmQ=
# REDIS_PASSWORD: eW91cl9zZWN1cmVfcmVkaXNfcGFzc3dvcmQ=
# JWT_SECRET: eW91cl9zZWN1cmVfand0X3NlY3JldA==

# Copy these values to your k8s/secrets.yaml file
```

## 🚀 **Deployment Commands**

### **Local Development**
```bash
# Start local development environment
./scripts/deploy.sh -e dev -p docker

# This will:
# 1. Use .env.local file
# 2. Deploy to Docker Compose
# 3. Only work on your local machine
```

### **Staging Environment**
```bash
# Deploy to staging
./scripts/deploy.sh -e staging -p k8s

# This will:
# 1. Use Kubernetes Secrets
# 2. Deploy to staging cluster
# 3. Use encrypted secrets
```

### **Production Environment**
```bash
# Deploy to production
./scripts/deploy.sh -e prod -p k8s

# This will:
# 1. Use Kubernetes Secrets
# 2. Deploy to production cluster
# 3. Use enterprise-grade security
```

### **Quick Deployment Reference**
```bash
# Local development
./scripts/deploy.sh                    # Default: dev + docker

# Staging deployment
./scripts/deploy.sh -e staging -p k8s  # Deploy to staging

# Production deployment
./scripts/deploy.sh -e prod -p k8s     # Deploy to production
```

## 🔍 **Validation and Checks**

### **Check Current Environment**
```bash
# See which environment files exist
ls -la | grep "\.env"

# Should only show:
# .env.local (for local development)
# env.example (template)
```

### **Validate Security**
```bash
# Check if .env.local is tracked by git
git status .env.local

# Should show: "Untracked files" or nothing
# If it shows as tracked, remove it:
git rm --cached .env.local
```

### **Check Service Health**
```bash
# Local development
./scripts/scopeapi.sh status

# Kubernetes deployment
kubectl get pods -n scopeapi
kubectl get services -n scopeapi
kubectl get secrets -n scopeapi
```

## 🆘 **Troubleshooting**

### **Common Issues**

#### **1. Script Refuses to Use .env File**
```bash
# Error: "Docker deployment is only allowed for LOCAL DEVELOPMENT"
# Solution: Use .env.local instead of .env
mv .env .env.local
```

#### **2. Script Refuses Docker for Staging/Production**
```bash
# Error: "Docker deployment is only allowed for LOCAL DEVELOPMENT"
# Solution: Use Kubernetes for staging/production
./scripts/deploy.sh -e staging -p k8s
./scripts/deploy.sh -e prod -p k8s
```

#### **3. Missing Environment File**
```bash
# Error: "No .env.local file found"
# Solution: Create from template
cp env.example .env.local
nano .env.local
```

#### **4. Kubernetes Secrets Not Loading**
```bash
# Check if secrets exist
kubectl get secrets -n scopeapi

# Check secret contents (base64 encoded)
kubectl get secret scopeapi-secrets -n scopeapi -o yaml

# Verify secret is mounted in pod
kubectl describe pod <pod-name> -n scopeapi
```

#### **5. Services Not Starting**
```bash
# Check pod status
kubectl get pods -n scopeapi

# Check pod events
kubectl describe pod <pod-name> -n scopeapi

# Check pod logs
kubectl logs <pod-name> -n scopeapi
```

### **Deployment Architecture**

```
┌─────────────────────────────────────────────────────────────────┐
│                    Production Environment                       │
├─────────────────────────────────────────────────────────────────┤
│  Load Balancer (Nginx/HAProxy) │  Monitoring (Prometheus/Grafana) │
├─────────────────────────────────────────────────────────────────┤
│                    Application Layer                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │API Discovery│ │Threat Detect│ │Data Protect │ │Attack Block │ │
│  │   (3 pods)  │ │   (3 pods)  │ │   (2 pods)  │ │   (3 pods)  │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │Gateway Integ│ │Data Ingest  │ │Admin Console│               │
│  │   (2 pods)  │ │   (3 pods)  │ │   (2 pods)  │               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
├─────────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │
│  │  PostgreSQL │ │    Kafka    │ │    Redis    │ │Elasticsearch│ │
│  │   (Primary) │ │   (Cluster) │ │   (Cluster) │ │   (Cluster) │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │PostgreSQL   │ │   Logging   │ │   Tracing   │               │
│  │ (Replicas)  │ │   (Fluentd) │ │  (Jaeger)   │               │
│  └─────────────┘ └─────────────┘ └─────────────┘               │
└─────────────────────────────────────────────────────────────────┘
```

## ✅ **Prerequisites**

### **System Requirements**

#### **Minimum Requirements**
- **CPU**: 8 cores (16 cores recommended)
- **Memory**: 16GB RAM (32GB recommended)
- **Storage**: 100GB SSD (500GB recommended)
- **Network**: 1Gbps (10Gbps recommended)

#### **Recommended Requirements**
- **CPU**: 16+ cores
- **Memory**: 64GB+ RAM
- **Storage**: 1TB+ NVMe SSD
- **Network**: 10Gbps with redundancy

### **Software Requirements**

#### **Operating System**
- **Linux**: Ubuntu 20.04+, CentOS 8+, RHEL 8+
- **Container Runtime**: Docker 24.0+ or containerd
- **Kernel**: 5.4+ with proper modules enabled

#### **Dependencies**
- **Docker**: 24.0+ with Docker Compose
- **Kubernetes**: 1.25+ (if using K8s)
- **Database**: PostgreSQL 15+
- **Message Queue**: Apache Kafka 3.4+
- **Cache**: Redis 7+

### **Network Requirements**

#### **Ports and Protocols**
- **HTTP/HTTPS**: 80, 443 (external)
- **Service Ports**: 8080-8086 (internal)
- **Database**: 5432 (PostgreSQL)
- **Message Queue**: 9092 (Kafka)
- **Cache**: 6379 (Redis)
- **Monitoring**: 9090 (Prometheus), 3000 (Grafana)

#### **Firewall Configuration**
```bash
# Allow external access
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp

# Allow internal service communication
sudo ufw allow 8080:8086/tcp
sudo ufw allow 5432/tcp
sudo ufw allow 9092/tcp
sudo ufw allow 6379/tcp
```

## ⚙️ **Environment Configuration**

> **📝 Note**: This section has been consolidated into the [Environment Strategy & Security](#environment-strategy--security) and [Local Development Setup](#local-development-setup) sections above. For local development, use `.env.local`. For staging/production, use Kubernetes Secrets.

  redis-primary:
    image: redis:7-alpine
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - scopeapi-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis-replica:
    image: redis:7-alpine
    command: redis-server --slaveof redis-primary 6379 --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_replica_data:/data
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      redis-primary:
        condition: service_healthy

  kafka-1:
    image: confluentinc/cp-kafka:7.4.0
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 2
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 3
    ports:
      - "9092:9092"
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      - zookeeper

  kafka-2:
    image: confluentinc/cp-kafka:7.4.0
    environment:
      KAFKA_BROKER_ID: 2
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:29092,PLAINTEXT_HOST://localhost:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 2
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 3
    ports:
      - "9093:9092"
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      - zookeeper

  kafka-3:
    image: confluentinc/cp-kafka:7.4.0
    environment:
      KAFKA_BROKER_ID: 3
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-3:29092,PLAINTEXT_HOST://localhost:9094
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
      KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 2
      KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 3
    ports:
      - "9094:9092"
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      - zookeeper

  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - scopeapi-network
    restart: unless-stopped

  # Monitoring Services
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
    networks:
      - scopeapi-network
    restart: unless-stopped
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'

  grafana:
    image: grafana/grafana:latest
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD:-admin}
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards:ro
      - ./monitoring/grafana/datasources:/etc/grafana/provisioning/datasources:ro
    ports:
      - "3000:3000"
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      - prometheus

  jaeger:
    image: jaegertracing/all-in-one:latest
    environment:
      COLLECTOR_OTLP_ENABLED: true
    ports:
      - "16686:16686"
      - "14268:14268"
    networks:
      - scopeapi-network
    restart: unless-stopped

  # Load Balancer
  nginx:
    image: nginx:alpine
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    ports:
      - "80:80"
      - "443:443"
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      - api-discovery
      - gateway-integration
      - data-ingestion
      - threat-detection
      - data-protection
      - attack-blocking
      - admin-console

  # Application Services
  api-discovery:
    build:
      context: ./backend/services/api-discovery
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

  gateway-integration:
    build:
      context: ./backend/services/gateway-integration
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

  data-ingestion:
    build:
      context: ./backend/services/data-ingestion
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

  threat-detection:
    build:
      context: ./backend/services/threat-detection
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

  data-protection:
    build:
      context: ./backend/services/data-protection
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

  attack-blocking:
    build:
      context: ./backend/services/attack-blocking
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

  admin-console:
    build:
      context: ./adminConsole
      dockerfile: Dockerfile
    environment:
      - DB_HOST=postgres-primary
      - DB_PORT=${POSTGRES_PORT}
      - DB_USER=${POSTGRES_USER}
      - DB_PASSWORD=${POSTGRES_PASSWORD}
      - DB_NAME=${POSTGRES_DB}
      - REDIS_HOST=redis-primary
      - REDIS_PORT=${REDIS_PORT}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - LOG_LEVEL=${LOG_LEVEL}
      - ENVIRONMENT=${ENVIRONMENT}
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy
      redis-primary:
        condition: service_healthy
      kafka-1:
        condition: service_started

volumes:
  postgres_data:
  postgres_replica_data:
  redis_data:
  redis_replica_data:
  prometheus_data:
  grafana_data:

networks:
  scopeapi-network:
    driver: bridge
```

## 🔄 **Kubernetes Migration**

### **Migration Overview**

#### **What We're Migrating From:**
- **Docker Compose** with `.env` files for configuration
- **Plain text passwords** stored in environment files
- **Single-server deployment** with limited scalability

#### **What We're Migrating To:**
- **Kubernetes** with proper secrets management
- **Encrypted secrets** stored securely in Kubernetes
- **Production-ready deployment** with auto-scaling and high availability

### **Phase 1: Prepare Your Environment**

#### **1.1 Install Kubernetes Tools**
```bash
# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# Install Docker Desktop with Kubernetes (Windows/macOS)
# Or use Minikube for local development
minikube start --driver=docker
```

#### **1.2 Backup Current Configuration**
```bash
# Backup your current .env file
cp .env .env.backup

# Backup Docker Compose configuration
cp scripts/docker-compose.yml scripts/docker-compose.yml.backup
```

### **Phase 2: Generate Kubernetes Secrets**

#### **2.1 Generate Base64 Encoded Secrets**
```bash
# Run the secret generation script
./scripts/generate-secrets.sh

# This will output base64 encoded values like:
# POSTGRES_PASSWORD: eW91cl9zZWN1cmVfcG9zdGdyZXNfcGFzc3dvcmQ=
# REDIS_PASSWORD: eW91cl9zZWN1cmVfcmVkaXNfcGFzc3dvcmQ=
```

#### **2.2 Update Kubernetes Secrets File**
```bash
# Edit the secrets file with your actual values
nano k8s/secrets.yaml

# Replace the placeholder values:
# POSTGRES_PASSWORD: <base64-encoded-postgres-password>
# With your actual base64 encoded values:
# POSTGRES_PASSWORD: eW91cl9zZWN1cmVfcG9zdGdyZXNfcGFzc3dvcmQ=
```

### **Phase 3: Deploy to Kubernetes**

#### **3.1 Deploy Infrastructure**
```bash
# Deploy to Kubernetes
./scripts/deploy.sh -e staging -p k8s

# Or use the comprehensive deployment script
./scripts/deploy.sh -e prod -p k8s
```

#### **3.2 Verify Deployment**
```bash
# Check namespace
kubectl get namespace scopeapi

# Check pods
kubectl get pods -n scopeapi

# Check services
kubectl get services -n scopeapi

# Check secrets
kubectl get secrets -n scopeapi
```

### **Phase 4: Test and Validate**

#### **4.1 Test Services**
```bash
# Test API Discovery service
kubectl port-forward service/api-discovery-service 8080:80 -n scopeapi
curl http://localhost:8080/health

# Test Gateway Integration service
kubectl port-forward service/gateway-integration-service 8081:80 -n scopeapi
curl http://localhost:8081/health
```

#### **4.2 Check Logs**
```bash
# View service logs
kubectl logs -f deployment/api-discovery -n scopeapi
kubectl logs -f deployment/gateway-integration -n scopeapi
```

### **Secrets Management Migration**

#### **Before (Docker Compose):**
```yaml
# docker-compose.yml
services:
  postgres:
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}  # From .env file
```

```bash
# .env file
POSTGRES_PASSWORD=your_plain_text_password  # ⚠️ UNSAFE!
```

#### **After (Kubernetes):**
```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: scopeapi-secrets
  namespace: scopeapi
type: Opaque
data:
  POSTGRES_PASSWORD: eW91cl9zZWN1cmVfcG9zdGdyZXNfcGFzc3dvcmQ=  # Base64 encoded
```

```yaml
# k8s/deployments/api-discovery-deployment.yaml
spec:
  template:
    spec:
      containers:
      - name: api-discovery
        envFrom:
        - secretRef:
            name: scopeapi-secrets  # Secrets injected automatically
```

### **Configuration Updates**

#### **Update Your Application Code**

##### **Go Services:**
```go
// Before: Direct environment variable access
dbPassword := os.Getenv("POSTGRES_PASSWORD")

// After: Same code works, but secrets are injected by Kubernetes
dbPassword := os.Getenv("POSTGRES_PASSWORD")
```

##### **Angular Frontend:**
```typescript
// Before: Environment variables in Angular
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080'
};

// After: Same configuration, but deployed via Kubernetes
export const environment = {
  production: true,
  apiUrl: 'https://api.yourdomain.com'
};
```

### **Security Improvements**

#### **What's More Secure Now:**

1. **Secrets Encryption**: Kubernetes encrypts secrets at rest
2. **RBAC**: Role-based access control for secrets
3. **Network Policies**: Pod-to-pod communication rules
4. **Security Contexts**: Non-root containers
5. **TLS Everywhere**: HTTPS for all external communication

#### **Access Control:**
```yaml
# k8s/rbac/service-account.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: scopeapi-role
  namespace: scopeapi
rules:
- apiGroups: [""]
  resources: ["secrets", "configmaps"]
  verbs: ["get", "list", "watch"]
```

### **Migration Checklist**

- [ ] **Environment Setup**
  - [ ] Install kubectl
  - [ ] Set up Kubernetes cluster
  - [ ] Backup current configuration

- [ ] **Secrets Migration**
  - [ ] Generate base64 encoded secrets
  - [ ] Update k8s/secrets.yaml
  - [ ] Verify secrets are properly formatted

- [ ] **Deployment**
  - [ ] Deploy to Kubernetes
  - [ ] Verify all services are running
  - [ ] Test service communication

- [ ] **Validation**
  - [ ] Run all tests
  - [ ] Verify functionality
  - [ ] Check performance

- [ ] **Cleanup**
  - [ ] Remove old Docker Compose deployment
  - [ ] Update documentation
  - [ ] Train team on new deployment process

## 🚀 **Production Deployment**

### **1. Prepare Production Environment**

```bash
# Create production directory
mkdir -p /opt/scopeapi
cd /opt/scopeapi

# Clone repository
git clone https://github.com/your-org/scopeapi.git .
git checkout production

# For production, use Kubernetes Secrets (see k8s/secrets.yaml)
# For development/staging, create environment file:
cp env.example .env.local
nano .env.local  # Edit with development values

# Create production compose file
cp docker-compose.prod.yml docker-compose.yml
```

### **2. SSL Certificate Setup**

```bash
# Create SSL directory
mkdir -p nginx/ssl

# Generate self-signed certificate (for testing)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/nginx.key \
  -out nginx/ssl/nginx.crt \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=yourdomain.com"

# For production, use Let's Encrypt or purchase certificates
# Let's Encrypt example:
certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com
cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem nginx/ssl/nginx.crt
cp /etc/letsencrypt/live/yourdomain.com/privkey.pem nginx/ssl/nginx.key
```

### **3. Nginx Configuration**

Create `nginx/nginx.conf`:

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api_discovery {
        server api-discovery:8080;
    }

    upstream gateway_integration {
        server gateway-integration:8081;
    }

    upstream data_ingestion {
        server data-ingestion:8082;
    }

    upstream threat_detection {
        server threat-detection:8083;
    }

    upstream data_protection {
        server data-protection:8084;
    }

    upstream attack_blocking {
        server attack-blocking:8085;
    }

    upstream admin_console {
        server admin-console:8086;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=admin:10m rate=5r/s;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

    # API Discovery Service
    server {
        listen 80;
        server_name api.yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name api.yourdomain.com;

        ssl_certificate /etc/nginx/ssl/nginx.crt;
        ssl_certificate_key /etc/nginx/ssl/nginx.key;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        location / {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://api_discovery;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }

    # Admin Console
    server {
        listen 80;
        server_name admin.yourdomain.com;
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name admin.yourdomain.com;

        ssl_certificate /etc/nginx/ssl/nginx.crt;
        ssl_certificate_key /etc/nginx/ssl/nginx.key;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        location / {
            limit_req zone=admin burst=10 nodelay;
            proxy_pass http://admin_console;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

### **4. Deploy Services**

```bash
# Start infrastructure services first
docker-compose up -d postgres-primary redis-primary zookeeper kafka-1 kafka-2 kafka-3

# Wait for infrastructure to be ready
sleep 30

# Start application services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### **5. Initialize Database**

```bash
# Run database setup
./scripts/setup-database.sh --test-data

# Verify database
./scripts/setup-database.sh --validate
```

## ☸️ **Kubernetes Deployment**

### **1. Create Namespace**

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: scopeapi
  labels:
    name: scopeapi
```

### **2. Create ConfigMap and Secrets**

**ConfigMap** (non-sensitive configuration):

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: scopeapi-config
  namespace: scopeapi
data:
  POSTGRES_HOST: "postgres-primary"
  POSTGRES_PORT: "5432"
  POSTGRES_DB: "scopeapi"
  REDIS_HOST: "redis-primary"
  REDIS_PORT: "6379"
  KAFKA_BROKERS: "kafka-1:9092,kafka-2:9092,kafka-3:9092"
  LOG_LEVEL: "info"
  ENVIRONMENT: "production"
```

### **3. Create Secret**

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: scopeapi-secrets
  namespace: scopeapi
type: Opaque
data:
  POSTGRES_PASSWORD: <base64-encoded-password>
  REDIS_PASSWORD: <base64-encoded-redis-password>
  JWT_SECRET: <base64-encoded-jwt-secret>
  ENCRYPTION_KEY: <base64-encoded-encryption-key>
```

### **4. Create Deployment**

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-discovery
  namespace: scopeapi
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-discovery
  template:
    metadata:
      labels:
        app: api-discovery
    spec:
      containers:
      - name: api-discovery
        image: scopeapi/api-discovery:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: scopeapi-config
        - secretRef:
            name: scopeapi-secrets
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### **5. Create Service**

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: api-discovery-service
  namespace: scopeapi
spec:
  selector:
    app: api-discovery
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

### **6. Create Ingress**

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: scopeapi-ingress
  namespace: scopeapi
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    - admin.yourdomain.com
    secretName: scopeapi-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-discovery-service
            port:
              number: 80
  - host: admin.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: admin-console-service
            port:
              number: 80
```

### **7. Deploy to Kubernetes**

```bash
# Apply all configurations
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n scopeapi
kubectl get services -n scopeapi
kubectl get ingress -n scopeapi

# View logs
kubectl logs -f deployment/api-discovery -n scopeapi
```

## 📊 **Monitoring and Observability**

### **Prometheus Configuration**

Create `monitoring/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'scopeapi-services'
    static_configs:
      - targets:
        - 'api-discovery:8080'
        - 'gateway-integration:8081'
        - 'data-ingestion:8082'
        - 'threat-detection:8083'
        - 'data-protection:8084'
        - 'attack-blocking:8085'
        - 'admin-console:8086'
    metrics_path: /metrics
    scrape_interval: 30s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-primary:5432']
    metrics_path: /metrics
    scrape_interval: 30s

  - job_name: 'redis'
    static_configs:
      - targets: ['redis-primary:6379']
    metrics_path: /metrics
    scrape_interval: 30s

  - job_name: 'kafka'
    static_configs:
      - targets: ['kafka-1:9092', 'kafka-2:9092', 'kafka-3:9092']
    metrics_path: /metrics
    scrape_interval: 30s
```

### **Grafana Dashboards**

Create `monitoring/grafana/dashboards/scopeapi-overview.json`:

```json
{
  "dashboard": {
    "id": null,
    "title": "ScopeAPI Overview",
    "tags": ["scopeapi", "overview"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Service Health",
        "type": "stat",
        "targets": [
          {
            "expr": "up{job=\"scopeapi-services\"}",
            "legendFormat": "{{instance}}"
          }
        ]
      },
      {
        "id": 2,
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{instance}}"
          }
        ]
      },
      {
        "id": 3,
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "{{instance}}"
          }
        ]
      }
    ]
  }
}
```

### **Alerting Rules**

Create `monitoring/rules/alerts.yml`:

```yaml
groups:
  - name: scopeapi
    rules:
      - alert: ServiceDown
        expr: up{job="scopeapi-services"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.instance }} is down"
          description: "Service {{ $labels.instance }} has been down for more than 1 minute"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time for {{ $labels.instance }}"
          description: "95th percentile response time is above 1 second"

      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate for {{ $labels.instance }}"
          description: "Error rate is above 5%"
```

## 🔒 **Security Configuration**

### **Network Security**

```bash
# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Configure Docker network security
docker network create --driver bridge --subnet=172.20.0.0/16 scopeapi-network
```

### **Container Security**

```dockerfile
# Example secure Dockerfile
FROM golang:1.21-alpine AS builder

# Install security updates
RUN apk update && apk upgrade

# Create non-root user
RUN addgroup -g 1001 -S scopeapi && \
    adduser -u 1001 -S scopeapi -G scopeapi

# Build application
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

# Production stage
FROM alpine:latest

# Install security updates
RUN apk update && apk upgrade

# Create non-root user
RUN addgroup -g 1001 -S scopeapi && \
    adduser -u 1001 -S scopeapi -G scopeapi

# Copy binary
COPY --from=builder /app/main /app/main

# Set ownership
RUN chown -R scopeapi:scopeapi /app

# Switch to non-root user
USER scopeapi

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["/app/main"]
```

### **Secrets Management**

```bash
# Use Docker secrets in production
echo "your_secure_password" | docker secret create postgres_password -

# Or use external secrets management
# HashiCorp Vault, AWS Secrets Manager, Azure Key Vault, etc.
```

## 💾 **Backup and Recovery**

### **Database Backup**

```bash
#!/bin/bash
# backup-database.sh

BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="scopeapi_$DATE.sql"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
docker exec scopeapi-postgres-primary pg_dump -U scopeapi scopeapi > $BACKUP_DIR/$BACKUP_FILE

# Compress backup
gzip $BACKUP_DIR/$BACKUP_FILE

# Keep only last 7 days of backups
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/$BACKUP_FILE.gz"
```

### **Configuration Backup**

```bash
#!/bin/bash
# backup-config.sh

BACKUP_DIR="/backups/config"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="config_$DATE.tar.gz"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup configuration files
tar -czf $BACKUP_DIR/$BACKUP_FILE \
    docker-compose.yml \
    .env.production \
    nginx/ \
    monitoring/ \
    k8s/

echo "Configuration backup completed: $BACKUP_DIR/$BACKUP_FILE"
```

### **Recovery Procedures**

```bash
# Database recovery
docker exec -i scopeapi-postgres-primary psql -U scopeapi scopeapi < backup_file.sql

# Configuration recovery
tar -xzf config_backup.tar.gz
docker-compose down
docker-compose up -d
```

## 📈 **Scaling and Performance**

### **Horizontal Scaling**

```bash
# Scale services
docker-compose up -d --scale api-discovery=3
docker-compose up -d --scale threat-detection=3
docker-compose up -d --scale data-ingestion=3

# Or in Kubernetes
kubectl scale deployment api-discovery --replicas=5 -n scopeapi
```

### **Load Balancing**

```nginx
# Nginx load balancing configuration
upstream api_discovery {
    server api-discovery-1:8080 weight=1;
    server api-discovery-2:8080 weight=1;
    server api-discovery-3:8080 weight=1;
}
```

### **Performance Tuning**

```bash
# Database tuning
docker exec -it scopeapi-postgres-primary psql -U scopeapi -c "
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
"

# Redis tuning
docker exec -it scopeapi-redis-primary redis-cli CONFIG SET maxmemory 512mb
docker exec -it scopeapi-redis-primary redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

## 🔍 **Troubleshooting**

### **Common Issues**

#### **Service Won't Start**
```bash
# Check logs
docker-compose logs service-name

# Check resource usage
docker stats

# Check network connectivity
docker exec service-name ping postgres-primary
```

#### **Database Connection Issues**
```bash
# Check PostgreSQL status
docker exec scopeapi-postgres-primary pg_isready -U scopeapi

# Check connection from service
docker exec service-name wget -qO- http://postgres-primary:5432
```

#### **Performance Issues**
```bash
# Check resource usage
docker stats
docker exec service-name top

# Check database performance
docker exec scopeapi-postgres-primary psql -U scopeapi -c "
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY total_time DESC
LIMIT 10;
"
```

### **Debugging Commands**

```bash
# Get service status
docker-compose ps

# View real-time logs
docker-compose logs -f

# Execute command in container
docker exec -it service-name sh

# Check network
docker network inspect scopeapi-network

# Check volumes
docker volume ls
docker volume inspect volume-name
```

## 🛠️ **Maintenance**

### **Regular Maintenance Tasks**

```bash
#!/bin/bash
# maintenance.sh

echo "Starting ScopeAPI maintenance..."

# Update images
docker-compose pull

# Restart services
docker-compose restart

# Clean up unused resources
docker system prune -f

# Check disk usage
df -h

# Check log sizes
du -sh /var/lib/docker/containers/*/*-json.log

echo "Maintenance completed"
```

### **Update Procedures**

```bash
# Update application
git pull origin production
docker-compose build
docker-compose up -d

# Update infrastructure
docker-compose pull postgres redis kafka
docker-compose up -d postgres redis kafka

# Verify update
docker-compose ps
./scripts/setup-database.sh --validate
```

### **Monitoring and Alerts**

```bash
# Check service health
curl -f http://localhost:8080/health
curl -f http://localhost:8081/health
curl -f http://localhost:8082/health

# Check monitoring
curl -f http://localhost:9090/-/healthy
curl -f http://localhost:3000/api/health
```

---

**🎯 This deployment guide helps you:**
- **Deploy** ScopeAPI to production environments
- **Configure** monitoring and observability
- **Secure** your deployment with best practices
- **Scale** your infrastructure as needed
- **Maintain** and troubleshoot your deployment

**For development setup, see our [Development Guide](DEVELOPMENT.md).**
