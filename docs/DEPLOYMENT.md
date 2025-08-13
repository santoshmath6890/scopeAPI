# 🚀 ScopeAPI Deployment Guide

This guide covers production deployment, monitoring, and operations for the ScopeAPI platform.

## 📋 **Table of Contents**

- [Deployment Overview](#deployment-overview)
- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Production Deployment](#production-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Monitoring and Observability](#monitoring-and-observability)
- [Security Configuration](#security-configuration)
- [Backup and Recovery](#backup-and-recovery)
- [Scaling and Performance](#scaling-and-performance)
- [Troubleshooting](#troubleshooting)
- [Maintenance](#maintenance)

## 🎯 **Deployment Overview**

ScopeAPI supports multiple deployment strategies:

- **Docker Compose** - Simple, single-server deployment
- **Kubernetes** - Production-grade, scalable deployment
- **Hybrid Cloud** - Multi-environment deployment
- **On-Premise** - Self-hosted deployment

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

### **Environment Variables**

Create a comprehensive `.env` file for production:

```bash
# Database Configuration
POSTGRES_HOST=postgres-primary
POSTGRES_PORT=5432
POSTGRES_USER=scopeapi
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_DB=scopeapi
POSTGRES_SSL_MODE=require

# Redis Configuration
REDIS_HOST=redis-primary
REDIS_PORT=6379
REDIS_PASSWORD=your_secure_redis_password
REDIS_DB=0

# Kafka Configuration
KAFKA_BROKERS=kafka-1:9092,kafka-2:9092,kafka-3:9092
KAFKA_TOPIC_PREFIX=scopeapi
KAFKA_CONSUMER_GROUP=scopeapi-consumer

# Security Configuration
JWT_SECRET=your_very_long_jwt_secret_key_here
JWT_EXPIRY=24h
ENCRYPTION_KEY=your_32_character_encryption_key

# Service Configuration
LOG_LEVEL=info
ENVIRONMENT=production
API_VERSION=v1
CORS_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com

# Monitoring Configuration
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
JAEGER_ENABLED=true
JAEGER_ENDPOINT=http://jaeger:14268/api/traces

# External Services
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=alerts@yourdomain.com
SMTP_PASSWORD=your_smtp_password
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/your/webhook/url
```

### **Configuration Files**

#### **Production Docker Compose**

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  # Infrastructure Services
  postgres-primary:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql:ro
    ports:
      - "5432:5432"
    networks:
      - scopeapi-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres-replica:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_MASTER_HOST: postgres-primary
      POSTGRES_MASTER_PORT: 5432
    volumes:
      - postgres_replica_data:/var/lib/postgresql/data
    networks:
      - scopeapi-network
    restart: unless-stopped
    depends_on:
      postgres-primary:
        condition: service_healthy

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

## 🚀 **Production Deployment**

### **1. Prepare Production Environment**

```bash
# Create production directory
mkdir -p /opt/scopeapi
cd /opt/scopeapi

# Clone repository
git clone https://github.com/your-org/scopeapi.git .
git checkout production

# Create production environment file
cp env.example .env.production
nano .env.production  # Edit with production values

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

### **2. Create ConfigMap**

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
