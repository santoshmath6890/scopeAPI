# 🔌 ScopeAPI API Reference

This document provides comprehensive API documentation for all ScopeAPI services, including endpoints, request/response formats, authentication, and examples.

## 📋 **Table of Contents**

- [API Overview](#api-overview)
- [Authentication](#authentication)
- [Common Patterns](#common-patterns)
- [API Discovery Service](#api-discovery-service)
- [Threat Detection Service](#threat-detection-service)
- [Data Protection Service](#data-protection-service)
- [Attack Blocking Service](#attack-blocking-service)
- [Gateway Integration Service](#gateway-integration-service)
- [Data Ingestion Service](#data-ingestion-service)
- [Admin Console Service](#admin-console-service)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [API Versioning](#api-versioning)

## 🎯 **API Overview**

ScopeAPI provides **RESTful APIs** for all services with consistent patterns:

- **Base URLs**: Each service has its own base URL and port
- **HTTP Methods**: Standard REST methods (GET, POST, PUT, DELETE)
- **Response Format**: JSON with consistent structure
- **Authentication**: JWT-based authentication
- **Versioning**: API versioning via URL path
- **Documentation**: OpenAPI/Swagger specifications

### **Service Endpoints**

| Service | Base URL | Port | Health Check |
|---------|----------|------|--------------|
| API Discovery | `http://localhost:8080` | 8080 | `/health` |
| Gateway Integration | `http://localhost:8081` | 8081 | `/health` |
| Data Ingestion | `http://localhost:8082` | 8082 | `/health` |
| Threat Detection | `http://localhost:8083` | 8083 | `/health` |
| Data Protection | `http://localhost:8084` | 8084 | `/health` |
| Attack Blocking | `http://localhost:8085` | 8085 | `/health` |
| Admin Console | `http://localhost:8086` | 8086 | `/health` |

### **API Response Format**

All APIs return responses in this consistent format:

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "message": "Operation completed successfully",
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_123456789"
}
```

### **Error Response Format**

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input parameters",
    "details": [
      {
        "field": "email",
        "message": "Email format is invalid"
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_123456789"
}
```

## 🔐 **Authentication**

### **JWT Authentication**

Most APIs require JWT authentication via the `Authorization` header:

```bash
Authorization: Bearer <jwt_token>
```

### **Getting a JWT Token**

```bash
# Login to get token
curl -X POST http://localhost:8086/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }'

# Response
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T18:30:00Z",
    "user": {
      "id": "user_123",
      "username": "admin",
      "role": "admin"
    }
  }
}
```

### **Using the Token**

```bash
# Use token in subsequent requests
curl -X GET http://localhost:8080/api/v1/endpoints \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 🔄 **Common Patterns**

### **Pagination**

Many list endpoints support pagination:

```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

### **Filtering and Sorting**

```bash
# Filter by status
GET /api/v1/endpoints?status=active

# Sort by creation date
GET /api/v1/endpoints?sort=created_at&order=desc

# Multiple filters
GET /api/v1/endpoints?status=active&method=GET&limit=50
```

### **Bulk Operations**

```bash
# Bulk update
PUT /api/v1/endpoints/bulk
{
  "ids": ["endpoint_1", "endpoint_2"],
  "updates": {
    "status": "inactive"
  }
}
```

## 🔍 **API Discovery Service**

**Base URL**: `http://localhost:8080`

### **Health Check**

```bash
GET /health
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "1.0.0",
    "uptime": "2h 15m 30s"
  }
}
```

### **Endpoints Management**

#### **List Endpoints**

```bash
GET /api/v1/endpoints
```

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 20, max: 100)
- `status` (string): Filter by status (active, inactive, archived)
- `method` (string): Filter by HTTP method
- `service` (string): Filter by service name
- `sort` (string): Sort field (created_at, updated_at, url)
- `order` (string): Sort order (asc, desc)

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "endpoint_123",
        "url": "/api/v1/users",
        "method": "GET",
        "service_name": "user-service",
        "status": "active",
        "discovered_at": "2024-01-15T10:30:00Z",
        "last_seen": "2024-01-15T10:30:00Z",
        "metadata": {
          "response_time_avg": 150,
          "success_rate": 99.5,
          "traffic_volume": 1000
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

#### **Get Endpoint by ID**

```bash
GET /api/v1/endpoints/{id}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "endpoint_123",
    "url": "/api/v1/users",
    "method": "GET",
    "service_name": "user-service",
    "status": "active",
    "discovered_at": "2024-01-15T10:30:00Z",
    "last_seen": "2024-01-15T10:30:00Z",
    "parameters": [
      {
        "name": "page",
        "type": "integer",
        "required": false,
        "default": 1
      },
      {
        "name": "limit",
        "type": "integer",
        "required": false,
        "default": 20
      }
    ],
    "responses": [
      {
        "status_code": 200,
        "content_type": "application/json",
        "schema": {
          "type": "object",
          "properties": {
            "users": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "id": {"type": "string"},
                  "name": {"type": "string"},
                  "email": {"type": "string"}
                }
              }
            }
          }
        }
      }
    ],
    "metadata": {
      "response_time_avg": 150,
      "success_rate": 99.5,
      "traffic_volume": 1000,
      "error_rate": 0.5,
      "last_error": null
    }
  }
}
```

#### **Create Endpoint**

```bash
POST /api/v1/endpoints
```

**Request Body:**
```json
{
  "url": "/api/v1/products",
  "method": "POST",
  "service_name": "product-service",
  "description": "Create a new product",
  "parameters": [
    {
      "name": "name",
      "type": "string",
      "required": true
    },
    {
      "name": "price",
      "type": "number",
      "required": true
    }
  ]
}
```

#### **Update Endpoint**

```bash
PUT /api/v1/endpoints/{id}
```

**Request Body:**
```json
{
  "status": "inactive",
  "description": "Updated description"
}
```

#### **Delete Endpoint**

```bash
DELETE /api/v1/endpoints/{id}
```

### **Discovery Operations**

#### **Start Discovery Scan**

```bash
POST /api/v1/discovery/scan
```

**Request Body:**
```json
{
  "target_urls": [
    "https://api.example.com",
    "https://api2.example.com"
  ],
  "scan_options": {
    "max_depth": 3,
    "include_swagger": true,
    "include_graphql": true,
    "rate_limit": 10
  }
}
```

#### **Get Scan Status**

```bash
GET /api/v1/discovery/scan/{scan_id}
```

#### **Get Discovery History**

```bash
GET /api/v1/discovery/history
```

## 🛡️ **Threat Detection Service**

**Base URL**: `http://localhost:8083`

### **Threat Detection**

#### **Analyze Request**

```bash
POST /api/v1/threats/analyze
```

**Request Body:**
```json
{
  "request": {
    "method": "POST",
    "url": "/api/v1/users",
    "headers": {
      "User-Agent": "Mozilla/5.0...",
      "Content-Type": "application/json"
    },
    "body": "{\"username\":\"admin\",\"password\":\"test123\"}",
    "ip_address": "192.168.1.100",
    "user_id": "user_123",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "context": {
    "session_id": "session_456",
    "user_agent": "Mozilla/5.0...",
    "referrer": "https://example.com/login"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "threat_level": "medium",
    "threats": [
      {
        "type": "brute_force",
        "confidence": 0.85,
        "description": "Multiple failed login attempts detected",
        "severity": "medium",
        "recommendations": [
          "Implement rate limiting",
          "Add CAPTCHA verification"
        ]
      }
    ],
    "risk_score": 75,
    "action_required": "monitor",
    "analysis_id": "analysis_789"
  }
}
```

#### **Get Threat History**

```bash
GET /api/v1/threats/history
```

**Query Parameters:**
- `threat_type` (string): Filter by threat type
- `severity` (string): Filter by severity (low, medium, high, critical)
- `start_date` (string): Start date (ISO 8601)
- `end_date` (string): End date (ISO 8601)
- `user_id` (string): Filter by user ID
- `ip_address` (string): Filter by IP address

#### **Get Threat Statistics**

```bash
GET /api/v1/threats/statistics
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_threats": 1250,
    "threats_by_type": {
      "brute_force": 450,
      "sql_injection": 200,
      "xss": 150,
      "rate_limit_exceeded": 300,
      "suspicious_pattern": 150
    },
    "threats_by_severity": {
      "low": 300,
      "medium": 600,
      "high": 300,
      "critical": 50
    },
    "threats_by_time": {
      "last_hour": 25,
      "last_24h": 450,
      "last_7d": 1250
    }
  }
}
```

## 🔒 **Data Protection Service**

**Base URL**: `http://localhost:8084`

### **Data Classification**

#### **Classify Data**

```bash
POST /api/v1/data/classify
```

**Request Body:**
```json
{
  "content": "User email: john.doe@example.com, SSN: 123-45-6789",
  "content_type": "text/plain",
  "context": {
    "source": "user_input",
    "field_name": "description",
    "form_id": "user_registration"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "classification": {
      "overall_risk": "high",
      "pii_detected": true,
      "categories": [
        {
          "type": "email",
          "confidence": 0.95,
          "value": "john.doe@example.com",
          "risk_level": "medium"
        },
        {
          "type": "ssn",
          "confidence": 0.98,
          "value": "123-45-6789",
          "risk_level": "high"
        }
      ]
    },
    "recommendations": [
      "Mask SSN in logs",
      "Encrypt sensitive data",
      "Implement data retention policy"
    ]
  }
}
```

#### **Get Data Classification Rules**

```bash
GET /api/v1/data/rules
```

#### **Create Classification Rule**

```bash
POST /api/v1/data/rules
```

**Request Body:**
```json
{
  "name": "Custom Credit Card Pattern",
  "pattern": "\\b\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}[\\s-]?\\d{4}\\b",
  "type": "credit_card",
  "risk_level": "high",
  "description": "Detects credit card numbers in various formats"
}
```

### **Compliance Monitoring**

#### **Get Compliance Report**

```bash
GET /api/v1/compliance/report
```

**Query Parameters:**
- `framework` (string): Compliance framework (GDPR, HIPAA, PCI-DSS)
- `start_date` (string): Start date
- `end_date` (string): End date

## ⚡ **Attack Blocking Service**

**Base URL**: `http://localhost:8085`

### **Blocking Rules**

#### **List Blocking Rules**

```bash
GET /api/v1/blocking/rules
```

#### **Create Blocking Rule**

```bash
POST /api/v1/blocking/rules
```

**Request Body:**
```json
{
  "name": "Block Suspicious IPs",
  "type": "ip_block",
  "conditions": {
    "ip_addresses": ["192.168.1.100", "10.0.0.50"],
    "reason": "Suspicious activity detected"
  },
  "action": "block",
  "duration": "24h",
  "enabled": true
}
```

#### **Update Blocking Rule**

```bash
PUT /api/v1/blocking/rules/{id}
```

#### **Delete Blocking Rule**

```bash
DELETE /api/v1/blocking/rules/{id}
```

### **Blocking Status**

#### **Check IP Status**

```bash
GET /api/v1/blocking/status/{ip_address}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "ip_address": "192.168.1.100",
    "status": "blocked",
    "blocked_at": "2024-01-15T10:30:00Z",
    "blocked_until": "2024-01-16T10:30:00Z",
    "reason": "Multiple failed login attempts",
    "rule_id": "rule_123"
  }
}
```

## 🌐 **Gateway Integration Service**

**Base URL**: `http://localhost:8081`

### **Gateway Management**

#### **List Gateways**

```bash
GET /api/v1/gateways
```

#### **Get Gateway Status**

```bash
GET /api/v1/gateways/{gateway_id}/status
```

#### **Deploy Configuration**

```bash
POST /api/v1/gateways/{gateway_id}/deploy
```

**Request Body:**
```json
{
  "configuration": {
    "routes": [
      {
        "path": "/api/v1/users",
        "upstream": "user-service:8080",
        "methods": ["GET", "POST"],
        "rate_limit": 100
      }
    ],
    "policies": [
      {
        "name": "authentication",
        "type": "jwt",
        "config": {
          "secret": "your-secret-key"
        }
      }
    ]
  },
  "validate_only": false
}
```

### **Gateway Types**

#### **Kong Gateway**

```bash
# Kong-specific configuration
POST /api/v1/gateways/kong/configure
```

#### **Envoy Gateway**

```bash
# Envoy-specific configuration
POST /api/v1/gateways/envoy/configure
```

#### **HAProxy Gateway**

```bash
# HAProxy-specific configuration
POST /api/v1/gateways/haproxy/configure
```

## 📥 **Data Ingestion Service**

**Base URL**: `http://localhost:8082`

### **Data Ingestion**

#### **Ingest Traffic Data**

```bash
POST /api/v1/ingestion/traffic
```

**Request Body:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "source": "api_gateway",
  "data": {
    "request_id": "req_123",
    "method": "POST",
    "url": "/api/v1/users",
    "status_code": 201,
    "response_time": 150,
    "user_id": "user_123",
    "ip_address": "192.168.1.100"
  }
}
```

#### **Ingest Security Events**

```bash
POST /api/v1/ingestion/security
```

**Request Body:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "event_type": "authentication_failure",
  "severity": "medium",
  "data": {
    "user_id": "user_123",
    "ip_address": "192.168.1.100",
    "reason": "Invalid credentials",
    "attempt_count": 5
  }
}
```

### **Data Processing**

#### **Get Processing Status**

```bash
GET /api/v1/ingestion/status
```

#### **Get Processing Statistics**

```bash
GET /api/v1/ingestion/statistics
```

## 📊 **Admin Console Service**

**Base URL**: `http://localhost:8086`

### **Authentication**

#### **Login**

```bash
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "username": "admin",
  "password": "password"
}
```

#### **Logout**

```bash
POST /api/v1/auth/logout
```

#### **Refresh Token**

```bash
POST /api/v1/auth/refresh
```

### **User Management**

#### **List Users**

```bash
GET /api/v1/users
```

#### **Create User**

```bash
POST /api/v1/users
```

**Request Body:**
```json
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "securepassword",
  "role": "analyst",
  "permissions": ["read", "write"]
}
```

#### **Update User**

```bash
PUT /api/v1/users/{id}
```

#### **Delete User**

```bash
DELETE /api/v1/users/{id}
```

### **System Configuration**

#### **Get System Settings**

```bash
GET /api/v1/system/settings
```

#### **Update System Settings**

```bash
PUT /api/v1/system/settings
```

**Request Body:**
```json
{
  "security": {
    "session_timeout": 3600,
    "max_login_attempts": 5,
    "password_policy": {
      "min_length": 8,
      "require_special_chars": true,
      "require_numbers": true
    }
  },
  "notifications": {
    "email_enabled": true,
    "slack_enabled": false,
    "alert_threshold": 75
  }
}
```

## ❌ **Error Handling**

### **HTTP Status Codes**

- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request parameters
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict
- **422 Unprocessable Entity**: Validation error
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Server error

### **Error Codes**

| Code | Description | HTTP Status |
|------|-------------|-------------|
| `VALIDATION_ERROR` | Input validation failed | 400 |
| `AUTHENTICATION_FAILED` | Invalid credentials | 401 |
| `INSUFFICIENT_PERMISSIONS` | User lacks required permissions | 403 |
| `RESOURCE_NOT_FOUND` | Requested resource not found | 404 |
| `RESOURCE_CONFLICT` | Resource already exists | 409 |
| `RATE_LIMIT_EXCEEDED` | Too many requests | 429 |
| `INTERNAL_ERROR` | Unexpected server error | 500 |

### **Error Response Examples**

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input parameters",
    "details": [
      {
        "field": "email",
        "message": "Email format is invalid",
        "value": "invalid-email"
      },
      {
        "field": "password",
        "message": "Password must be at least 8 characters",
        "value": "123"
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_123456789"
}
```

## 🚦 **Rate Limiting**

### **Rate Limit Headers**

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642239000
X-RateLimit-Reset-Time: 2024-01-15T11:30:00Z
```

### **Rate Limit Exceeded Response**

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Try again in 60 seconds.",
    "retry_after": 60
  }
}
```

## 🔄 **API Versioning**

### **Version in URL Path**

```bash
# Current version
GET /api/v1/endpoints

# Future version
GET /api/v2/endpoints
```

### **Version Header**

```http
Accept: application/vnd.scopeapi.v1+json
```

### **Version Deprecation**

```http
X-API-Version-Deprecated: true
X-API-Version-Sunset: 2024-12-31T23:59:59Z
```

## 📚 **Additional Resources**

### **OpenAPI Specifications**

Each service provides OpenAPI/Swagger documentation:

- **API Discovery**: `http://localhost:8080/swagger/`
- **Gateway Integration**: `http://localhost:8081/swagger/`
- **Data Ingestion**: `http://localhost:8082/swagger/`
- **Threat Detection**: `http://localhost:8083/swagger/`
- **Data Protection**: `http://localhost:8084/swagger/`
- **Attack Blocking**: `http://localhost:8085/swagger/`
- **Admin Console**: `http://localhost:8086/swagger/`

### **SDK and Client Libraries**

- **Go Client**: `github.com/scopeapi/go-client`
- **Python Client**: `pip install scopeapi-client`
- **JavaScript Client**: `npm install scopeapi-client`
- **Postman Collection**: Available in `/docs/postman/`

### **Testing and Examples**

- **Postman Collection**: Complete API testing collection
- **cURL Examples**: Command-line examples for all endpoints
- **Code Samples**: Examples in Go, Python, JavaScript
- **Integration Tests**: Automated API testing suite

---

**🎯 This API reference helps you:**
- **Understand** all available endpoints and their usage
- **Integrate** ScopeAPI into your applications
- **Test** API functionality and behavior
- **Build** custom clients and integrations
- **Troubleshoot** API issues and errors

**For more information, see our [Development Guide](DEVELOPMENT.md) and [Architecture Guide](ARCHITECTURE.md).**
