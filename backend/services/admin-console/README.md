# Admin Console Microservice

This microservice provides the admin console functionality for the ScopeAPI platform. It serves both the Angular frontend application and provides backend APIs for admin operations.

## Features

- **Angular Frontend**: Serves the admin console Angular application
- **User Management**: CRUD operations for user management
- **Dashboard Statistics**: Real-time dashboard metrics and statistics
- **System Monitoring**: System health, logs, and performance metrics
- **API Gateway**: Centralized API endpoints for admin operations

## Architecture

The service is built using:
- **Go** with Gin framework for the backend API
- **Angular** for the frontend application
- **Multi-stage Docker** build for optimized containerization

## API Endpoints

### Health Check
- `GET /api/v1/health` - Service health status

### Dashboard
- `GET /api/v1/dashboard/stats` - Dashboard statistics

### User Management
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### System
- `GET /api/v1/system/info` - System information
- `GET /api/v1/system/logs` - System logs

## Configuration

The service uses a YAML configuration file (`config/config.yaml`) with the following sections:

- **Server**: Host and port configuration
- **CORS**: Cross-origin resource sharing settings
- **Static**: Static file serving configuration
- **Database**: Database connection settings
- **Redis**: Redis connection settings
- **Services**: Other microservice endpoints

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker (optional)

### Local Development

1. **Build Angular application**:
   ```bash
   cd ../../adminConsole
   npm install
   npm run build
   ```

2. **Run the Go service**:
   ```bash
   cd backend/services/admin-console
   go mod tidy
   go run cmd/main.go
   ```

3. **Access the application**:
   - Frontend: http://localhost:8080
   - API: http://localhost:8080/api/v1/health

### Docker Development

1. **Build the Docker image**:
   ```bash
   docker build -t scopeapi-admin-console .
   ```

2. **Run the container**:
   ```bash
   docker run -p 8080:8080 scopeapi-admin-console
   ```

## Production Deployment

### Docker Compose
Add the following service to your `docker-compose.yml`:

```yaml
admin-console:
  build: ./backend/services/admin-console
  ports:
    - "8080:8080"
  environment:
    - ENVIRONMENT=production
  depends_on:
    - postgres
    - redis
```

### Kubernetes
The service can be deployed to Kubernetes using the provided manifests in the `k8s/` directory.

## Environment Variables

- `ENVIRONMENT`: Set to "production" for production mode
- `SERVER_PORT`: Port to run the server on (default: 8080)
- `SERVER_HOST`: Host to bind to (default: 0.0.0.0)

## Monitoring

The service provides health check endpoints and logs in JSON format for easy integration with monitoring systems.

## Security

- CORS configuration for cross-origin requests
- Input validation for all API endpoints
- Secure headers and middleware
- Authentication and authorization (to be implemented)

## Contributing

1. Follow the Go coding standards
2. Add tests for new functionality
3. Update documentation for API changes
4. Ensure the Angular build works correctly

## License

This project is part of the ScopeAPI platform and follows the same licensing terms. 