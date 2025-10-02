# Threat Detection Service

The Threat Detection Service is a comprehensive security component of ScopeAPI that provides real-time threat detection, anomaly analysis, behavioral monitoring, and signature-based detection capabilities.

## Overview

This service implements multiple detection methodologies to identify and mitigate various types of security threats targeting APIs:

- **Signature-based Detection**: Identifies known attack patterns using predefined signatures
- **Anomaly Detection**: Detects deviations from normal traffic patterns using statistical and ML methods
- **Behavioral Analysis**: Analyzes user and entity behavior patterns to identify suspicious activities
- **Real-time Processing**: Processes incoming traffic and security events in real-time via Kafka

## Features

### Core Detection Capabilities

1. **Threat Detection**
   - SQL Injection detection
   - Cross-Site Scripting (XSS) detection
   - DDoS attack detection
   - Brute force attack detection
   - Data exfiltration detection
   - Path traversal detection
   - Command injection detection

2. **Anomaly Detection**
   - Traffic volume anomalies
   - Response time anomalies
   - Request pattern anomalies
   - Geolocation anomalies
   - Statistical anomaly detection
   - Machine learning-based anomaly detection

3. **Behavioral Analysis**
   - Access pattern analysis
   - Usage pattern analysis
   - Timing pattern analysis
   - Sequence pattern analysis
   - Location pattern analysis
   - Risk scoring and assessment

4. **Signature Management**
   - Signature-based threat detection
   - Custom signature creation
   - Signature import/export
   - Signature testing and validation
   - Signature performance metrics

### API Endpoints

#### Threat Detection
- `GET /api/v1/threats` - List threats with filtering
- `GET /api/v1/threats/:id` - Get specific threat details
- `POST /api/v1/threats/analyze` - Analyze traffic for threats
- `PUT /api/v1/threats/:id/status` - Update threat status
- `DELETE /api/v1/threats/:id` - Delete threat record

#### Anomaly Detection
- `GET /api/v1/anomalies` - List anomalies with filtering
- `GET /api/v1/anomalies/:id` - Get specific anomaly details
- `POST /api/v1/anomalies/detect` - Detect anomalies in traffic
- `PUT /api/v1/anomalies/:id/feedback` - Provide feedback on anomalies

#### Behavioral Analysis
- `GET /api/v1/behavioral/patterns` - List behavior patterns
- `POST /api/v1/behavioral/analyze` - Analyze behavior patterns
- `GET /api/v1/behavioral/baselines` - Get baseline profiles
- `POST /api/v1/behavioral/baselines` - Create baseline profiles

#### Signature Management
- `GET /api/v1/signatures` - List threat signatures
- `GET /api/v1/signatures/:id` - Get specific signature details
- `POST /api/v1/signatures/detect` - Detect signatures in traffic
- `POST /api/v1/signatures/test` - Test signature against data
- `POST /api/v1/signatures/import` - Import signature set
- `GET /api/v1/signatures/export/:set` - Export signature set

## Architecture

### Components

1. **Handlers** (`internal/handlers/`)
   - HTTP request handlers for all API endpoints
   - Request validation and response formatting
   - Error handling and logging

2. **Services** (`internal/services/`)
   - Core business logic for threat detection
   - Anomaly detection algorithms
   - Behavioral analysis engines
   - Signature matching logic

3. **Models** (`internal/models/`)
   - Data structures for threats, anomalies, behaviors
   - Request/response models
   - Configuration models

4. **Repository** (`internal/repository/`)
   - Database access layer
   - Data persistence operations
   - Query optimization

5. **Configuration** (`internal/config/`)
   - Configuration management
   - Environment variable handling
   - Service configuration

### Database Schema

The service uses PostgreSQL with the following main tables:

- `threats` - Stores detected threats and attacks
- `anomalies` - Stores detected anomalies
- `behavior_patterns` - Stores behavioral analysis patterns
- `threat_signatures` - Stores detection signatures
- `baseline_profiles` - Stores behavioral baselines
- `anomaly_feedback` - Stores user feedback on anomalies
- `threat_statistics` - Stores aggregated statistics

## Installation and Setup

### Prerequisites

- Go 1.22+
- PostgreSQL 12+
- Kafka 2.8+
- Redis (for caching)

### Database Setup

1. Create the database:
   ```sql
   CREATE DATABASE threat_detection;
   ```

2. Run migrations:
   ```bash
   cd migrations
   export DATABASE_URL="postgres://user:password@localhost/threat_detection?sslmode=disable"
   go run migrate.go
   ```

### Configuration

Create a configuration file `config/threat-detection.yaml`:

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  postgresql:
    host: "localhost"
    port: 5432
    user: "threat_detection"
    password: "password"
    database: "threat_detection"
    ssl_mode: "disable"

messaging:
  kafka:
    brokers: ["localhost:9092"]

auth:
  jwt:
    secret: "your-jwt-secret"
    expiration: "24h"

logging:
  level: "info"
  format: "json"

metrics:
  enabled: true
  port: "9090"
```

### Running the Service

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Set environment variables:
   ```bash
   export DATABASE_URL="postgres://user:password@localhost/threat_detection?sslmode=disable"
   export KAFKA_BROKERS="localhost:9092"
   ```

3. Run the service:
   ```bash
   go run cmd/main.go
   ```

## Usage Examples

### Detecting Threats

```bash
curl -X POST http://localhost:8080/api/v1/threats/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "traffic_data": {
      "method": "POST",
      "url": "/api/users",
      "headers": {
        "User-Agent": "Mozilla/5.0...",
        "Content-Type": "application/json"
      },
      "body": "{\"username\":\"admin\"; DROP TABLE users;--\"}",
      "ip_address": "192.168.1.100",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  }'
```

### Detecting Anomalies

```bash
curl -X POST http://localhost:8080/api/v1/anomalies/detect \
  -H "Content-Type: application/json" \
  -d '{
    "traffic_data": {
      "api_id": "api-123",
      "endpoint_id": "endpoint-456",
      "request_count": 1000,
      "response_time_avg": 5000,
      "error_rate": 0.8,
      "timestamp": "2024-01-15T10:30:00Z"
    },
    "sensitivity": 0.7
  }'
```

### Analyzing Behavior

```bash
curl -X POST http://localhost:8080/api/v1/behavioral/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "entity_id": "user-123",
    "entity_type": "user",
    "analysis_period": {
      "start": "2024-01-01T00:00:00Z",
      "end": "2024-01-15T23:59:59Z"
    },
    "traffic_data": [...]
  }'
```

## Monitoring and Metrics

The service exposes Prometheus metrics at `/metrics`:

- `threat_detection_threats_total` - Total threats detected
- `threat_detection_anomalies_total` - Total anomalies detected
- `threat_detection_signatures_matched_total` - Total signature matches
- `threat_detection_processing_duration_seconds` - Processing duration
- `threat_detection_requests_total` - Total API requests

## Security Considerations

1. **Authentication**: JWT-based authentication (configurable)
2. **Authorization**: Role-based access control
3. **Data Encryption**: Sensitive data encrypted at rest
4. **Audit Logging**: Comprehensive audit trails
5. **Rate Limiting**: Built-in rate limiting for API endpoints

## Performance Optimization

1. **Database Indexing**: Optimized indexes for common queries
2. **Caching**: Redis-based caching for frequently accessed data
3. **Async Processing**: Kafka-based asynchronous processing
4. **Connection Pooling**: Database connection pooling
5. **Batch Processing**: Batch operations for bulk data

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check database credentials
   - Verify network connectivity
   - Check PostgreSQL logs

2. **Kafka Connection Issues**
   - Verify Kafka brokers are running
   - Check network connectivity
   - Verify topic configuration

3. **High Memory Usage**
   - Check for memory leaks in detection algorithms
   - Monitor cache usage
   - Review batch processing sizes

### Logs

The service provides structured logging with different levels:
- `DEBUG`: Detailed debugging information
- `INFO`: General information
- `WARN`: Warning messages
- `ERROR`: Error conditions
- `FATAL`: Fatal errors

## Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Follow semantic versioning

## License

This project is part of the ScopeAPI platform and follows the same licensing terms.
