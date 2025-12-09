# 🚀 ScopeAPI Kubernetes Configuration

This directory contains all Kubernetes configuration files for deploying ScopeAPI to production.

## 📁 Directory Structure

```
k8s/
├── namespace.yaml           # Namespace and resource quotas
├── configmap.yaml          # Non-sensitive configuration
├── secrets.yaml            # Sensitive data (DO NOT COMMIT)
├── rbac/                   # Role-based access control
├── deployments/            # Application deployments
├── services/               # Service definitions
├── ingress/                # Ingress and routing
├── monitoring/             # Monitoring and observability
├── policies/               # Security policies
└── external-secrets/       # External secrets management
```

## 🔐 Security Notice

**⚠️ IMPORTANT: Never commit the `secrets.yaml` file to version control!**

- Use `scripts/generate-secrets.sh` to generate base64 encoded values
- Store actual secrets in external secrets managers (Vault, AWS Secrets Manager)
- Use Kubernetes secrets for sensitive data
- Rotate secrets regularly

## 🚀 Quick Start

1. **Generate secrets:**
   ```bash
   ./scripts/generate-secrets.sh
   ```

2. **Update secrets.yaml with real values**

3. **Deploy to Kubernetes:**
   ```bash
   ./scripts/deploy-k8s.sh
   ```

## 🔧 Configuration

### Environment Variables
- **Development**: Use `.env.local` files
- **Staging**: Use `.env.staging` files  
- **Production**: Use Kubernetes secrets

### Secrets Management
- **Local Development**: Docker Compose with .env files
- **Staging**: Docker Secrets
- **Production**: Kubernetes Secrets + External Secrets Manager

## 📊 Monitoring

- **Prometheus**: Metrics collection
- **Grafana**: Visualization and dashboards
- **Jaeger**: Distributed tracing
- **Health Checks**: Built into all services

## 🔒 Security Features

- **RBAC**: Role-based access control
- **Network Policies**: Pod-to-pod communication rules
- **Security Contexts**: Non-root containers
- **TLS**: HTTPS everywhere
- **Rate Limiting**: API protection

## 🛠️ Maintenance

- **Updates**: Rolling updates with zero downtime
- **Scaling**: Horizontal pod autoscaling
- **Backups**: Automated database backups
- **Logs**: Centralized logging with Fluentd

## 📚 Documentation

- [Deployment Guide](../docs/DEPLOYMENT.md)
- [Development Guide](../docs/DEVELOPMENT.md)
- [Architecture Guide](../docs/ARCHITECTURE.md)
