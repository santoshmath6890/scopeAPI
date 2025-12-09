# Gateway Integration API Documentation

## Overview

The Gateway Integration API provides comprehensive management capabilities for various API gateways and load balancers including Kong, NGINX, Traefik, Envoy, and HAProxy. This API enables centralized monitoring, configuration management, and health checking across multiple gateway instances.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API endpoints require authentication. Include the following header in your requests:

```
Authorization: Bearer <your-jwt-token>
```

## Data Models

### Integration

```json
{
  "id": "string",
  "name": "string",
  "type": "kong|nginx|traefik|envoy|haproxy",
  "description": "string",
  "status": "active|inactive|error",
  "endpoints": [
    {
      "url": "string",
      "type": "string",
      "timeout": "number"
    }
  ],
  "credentials": {
    "type": "api_key|basic_auth|oauth2",
    "username": "string",
    "password": "string",
    "api_key": "string"
  },
  "configuration": {
    "version": "string",
    "admin_url": "string",
    "status_url": "string"
  },
  "health_status": {
    "status": "healthy|unhealthy|unknown",
    "message": "string",
    "timestamp": "string"
  },
  "last_sync_at": "string",
  "created_at": "string",
  "updated_at": "string"
}
```

### HealthStatus

```json
{
  "status": "healthy|unhealthy|unknown",
  "message": "string",
  "timestamp": "string"
}
```

### SyncResult

```json
{
  "status": "completed|failed|in_progress",
  "message": "string",
  "changes": [
    {
      "type": "added|modified|deleted",
      "resource": "string",
      "details": "string"
    }
  ],
  "last_sync_at": "string"
}
```

### IntegrationStats

```json
{
  "total_integrations": "number",
  "healthy_count": "number",
  "unhealthy_count": "number",
  "unknown_count": "number",
  "by_type": {
    "kong": "number",
    "nginx": "number",
    "traefik": "number",
    "envoy": "number",
    "haproxy": "number"
  }
}
```

## API Endpoints

### 1. List Integrations

**GET** `/integrations`

Retrieve a list of all gateway integrations with optional filtering.

#### Query Parameters

- `type` (optional): Filter by gateway type (`kong`, `nginx`, `traefik`, `envoy`, `haproxy`)
- `status` (optional): Filter by status (`active`, `inactive`, `error`)
- `limit` (optional): Number of results to return (default: 50, max: 100)
- `offset` (optional): Number of results to skip (default: 0)

#### Example Request

```bash
curl -X GET "http://localhost:8080/api/v1/integrations?type=kong&status=active" \
  -H "Authorization: Bearer <your-token>"
```

#### Example Response

```json
{
  "data": [
    {
      "id": "kong-1",
      "name": "Production Kong",
      "type": "kong",
      "description": "Main Kong gateway for production APIs",
      "status": "active",
      "endpoints": [
        {
          "url": "http://kong-admin:8001",
          "type": "admin",
          "timeout": 30
        }
      ],
      "health_status": {
        "status": "healthy",
        "message": "Connection successful",
        "timestamp": "2024-01-01T12:00:00Z"
      },
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 1,
  "limit": 50,
  "offset": 0
}
```

### 2. Get Integration

**GET** `/integrations/{id}`

Retrieve detailed information about a specific integration.

#### Example Request

```bash
curl -X GET "http://localhost:8080/api/v1/integrations/kong-1" \
  -H "Authorization: Bearer <your-token>"
```

#### Example Response

```json
{
  "id": "kong-1",
  "name": "Production Kong",
  "type": "kong",
  "description": "Main Kong gateway for production APIs",
  "status": "active",
  "endpoints": [
    {
      "url": "http://kong-admin:8001",
      "type": "admin",
      "timeout": 30
    }
  ],
  "credentials": {
    "type": "api_key",
    "username": "admin",
    "password": "password"
  },
  "configuration": {
    "version": "2.8.0",
    "admin_url": "http://kong-admin:8001"
  },
  "health_status": {
    "status": "healthy",
    "message": "Connection successful",
    "timestamp": "2024-01-01T12:00:00Z"
  },
  "last_sync_at": "2024-01-01T11:30:00Z",
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 3. Create Integration

**POST** `/integrations`

Create a new gateway integration.

#### Request Body

```json
{
  "name": "Test Kong",
  "type": "kong",
  "description": "Test Kong integration",
  "status": "active",
  "endpoints": [
    {
      "url": "http://localhost:8001",
      "type": "admin",
      "timeout": 30
    }
  ],
  "credentials": {
    "type": "api_key",
    "username": "admin",
    "password": "password"
  },
  "configuration": {
    "version": "2.8.0",
    "admin_url": "http://localhost:8001"
  }
}
```

#### Example Request

```bash
curl -X POST "http://localhost:8080/api/v1/integrations" \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Kong",
    "type": "kong",
    "description": "Test Kong integration",
    "endpoints": [
      {
        "url": "http://localhost:8001",
        "type": "admin",
        "timeout": 30
      }
    ],
    "credentials": {
      "type": "api_key",
      "username": "admin",
      "password": "password"
    }
  }'
```

#### Example Response

```json
{
  "id": "new-kong-id",
  "name": "Test Kong",
  "type": "kong",
  "description": "Test Kong integration",
  "status": "active",
  "endpoints": [
    {
      "url": "http://localhost:8001",
      "type": "admin",
      "timeout": 30
    }
  ],
  "credentials": {
    "type": "api_key",
    "username": "admin",
    "password": "password"
  },
  "configuration": {
    "version": "2.8.0",
    "admin_url": "http://localhost:8001"
  },
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 4. Update Integration

**PUT** `/integrations/{id}`

Update an existing integration.

#### Request Body

```json
{
  "name": "Updated Kong",
  "description": "Updated description",
  "status": "active",
  "endpoints": [
    {
      "url": "http://localhost:8001",
      "type": "admin",
      "timeout": 60
    }
  ]
}
```

#### Example Request

```bash
curl -X PUT "http://localhost:8080/api/v1/integrations/kong-1" \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Kong",
    "description": "Updated description",
    "endpoints": [
      {
        "url": "http://localhost:8001",
        "type": "admin",
        "timeout": 60
      }
    ]
  }'
```

#### Example Response

```json
{
  "id": "kong-1",
  "name": "Updated Kong",
  "type": "kong",
  "description": "Updated description",
  "status": "active",
  "endpoints": [
    {
      "url": "http://localhost:8001",
      "type": "admin",
      "timeout": 60
    }
  ],
  "updated_at": "2024-01-01T12:30:00Z"
}
```

### 5. Delete Integration

**DELETE** `/integrations/{id}`

Delete an integration.

#### Example Request

```bash
curl -X DELETE "http://localhost:8080/api/v1/integrations/kong-1" \
  -H "Authorization: Bearer <your-token>"
```

#### Response

- **Status**: 204 No Content

### 6. Test Integration

**POST** `/integrations/{id}/test`

Test the connection to a gateway integration.

#### Example Request

```bash
curl -X POST "http://localhost:8080/api/v1/integrations/kong-1/test" \
  -H "Authorization: Bearer <your-token>"
```

#### Example Response

```json
{
  "status": "healthy",
  "message": "Connection successful",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 7. Sync Integration

**POST** `/integrations/{id}/sync`

Synchronize configuration with the gateway.

#### Example Request

```bash
curl -X POST "http://localhost:8080/api/v1/integrations/kong-1/sync" \
  -H "Authorization: Bearer <your-token>"
```

#### Example Response

```json
{
  "status": "completed",
  "message": "Sync completed successfully",
  "changes": [
    {
      "type": "added",
      "resource": "service",
      "details": "Added new service 'api-service'"
    },
    {
      "type": "modified",
      "resource": "route",
      "details": "Updated route 'api-route' configuration"
    }
  ],
  "last_sync_at": "2024-01-01T12:00:00Z"
}
```

### 8. Get Integration Statistics

**GET** `/integrations/stats`

Retrieve statistics about all integrations.

#### Example Request

```bash
curl -X GET "http://localhost:8080/api/v1/integrations/stats" \
  -H "Authorization: Bearer <your-token>"
```

#### Example Response

```json
{
  "total_integrations": 5,
  "healthy_count": 3,
  "unhealthy_count": 1,
  "unknown_count": 1,
  "by_type": {
    "kong": 2,
    "nginx": 1,
    "traefik": 1,
    "envoy": 1,
    "haproxy": 0
  }
}
```

## Gateway-Specific Operations

### Kong Gateway

For Kong-specific operations, the API supports additional endpoints:

#### Get Kong Services

**GET** `/integrations/{id}/kong/services`

#### Get Kong Routes

**GET** `/integrations/{id}/kong/routes`

#### Get Kong Plugins

**GET** `/integrations/{id}/kong/plugins`

### NGINX Gateway

For NGINX-specific operations:

#### Get NGINX Status

**GET** `/integrations/{id}/nginx/status`

#### Get NGINX Configuration

**GET** `/integrations/{id}/nginx/config`

## Error Responses

### 400 Bad Request

```json
{
  "error": "validation_error",
  "message": "Invalid request data",
  "details": {
    "field": "name",
    "issue": "Name is required"
  }
}
```

### 401 Unauthorized

```json
{
  "error": "unauthorized",
  "message": "Invalid or missing authentication token"
}
```

### 404 Not Found

```json
{
  "error": "not_found",
  "message": "Integration not found",
  "resource_id": "kong-1"
}
```

### 500 Internal Server Error

```json
{
  "error": "internal_error",
  "message": "An internal server error occurred",
  "request_id": "req-123456"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Rate Limit**: 1000 requests per hour per API key
- **Headers**: 
  - `X-RateLimit-Limit`: Maximum requests per hour
  - `X-RateLimit-Remaining`: Remaining requests in current hour
  - `X-RateLimit-Reset`: Time when the rate limit resets

## Pagination

List endpoints support pagination using `limit` and `offset` parameters:

```json
{
  "data": [...],
  "total": 100,
  "limit": 50,
  "offset": 0,
  "has_more": true
}
```

## Webhooks

The API supports webhooks for real-time notifications:

### Webhook Events

- `integration.created`
- `integration.updated`
- `integration.deleted`
- `integration.health_changed`
- `integration.sync_completed`

### Webhook Payload

```json
{
  "event": "integration.health_changed",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "integration_id": "kong-1",
    "old_status": "healthy",
    "new_status": "unhealthy",
    "message": "Connection timeout"
  }
}
```

## SDK Examples

### Go

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
)

type Integration struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"`
}

func main() {
    client := &http.Client{}
    
    req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/integrations", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    req.Header.Set("Authorization", "Bearer your-token")
    
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    var integrations []Integration
    if err := json.NewDecoder(resp.Body).Decode(&integrations); err != nil {
        log.Fatal(err)
    }
    
    for _, integration := range integrations {
        fmt.Printf("Integration: %s (%s)\n", integration.Name, integration.Type)
    }
}
```

### JavaScript/Node.js

```javascript
const axios = require('axios');

const apiClient = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
    headers: {
        'Authorization': 'Bearer your-token'
    }
});

async function getIntegrations() {
    try {
        const response = await apiClient.get('/integrations');
        return response.data;
    } catch (error) {
        console.error('Error fetching integrations:', error.response.data);
        throw error;
    }
}

async function createIntegration(integrationData) {
    try {
        const response = await apiClient.post('/integrations', integrationData);
        return response.data;
    } catch (error) {
        console.error('Error creating integration:', error.response.data);
        throw error;
    }
}

// Usage
getIntegrations()
    .then(integrations => console.log('Integrations:', integrations))
    .catch(error => console.error('Failed:', error));
```

## Testing

### Using curl

```bash
# Test the API health
curl -X GET "http://localhost:8080/api/v1/integrations/stats" \
  -H "Authorization: Bearer your-token"

# Create a test integration
curl -X POST "http://localhost:8080/api/v1/integrations" \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Kong",
    "type": "kong",
    "endpoints": [{"url": "http://localhost:8001", "type": "admin"}]
  }'
```

### Using Postman

1. Import the API collection
2. Set the base URL to `http://localhost:8080/api/v1`
3. Add the Authorization header with your Bearer token
4. Test the endpoints

## Support

For API support and questions:

- **Documentation**: [API Documentation](https://docs.scopeapi.com/gateway-integration)
- **Support Email**: api-support@scopeapi.com
- **GitHub Issues**: [GitHub Repository](https://github.com/scopeapi/scopeapi/issues) 