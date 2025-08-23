# 🛡️ Data Protection Service (DPS)

The Data Protection Service is a comprehensive microservice designed to protect sensitive data through classification, PII detection, compliance management, and risk assessment.

## 🚀 Features

### **Data Classification**
- **Rule-based Classification**: Define custom rules for data categorization
- **ML-powered Classification**: Machine learning algorithms for automatic data classification
- **Pattern Matching**: Regex and keyword-based classification
- **Context Analysis**: Intelligent context-aware classification
- **Classification Levels**: Public, Internal, Confidential, Restricted, Top Secret

### **PII Detection**
- **Pattern Recognition**: Regex-based PII pattern detection
- **Machine Learning**: ML models for enhanced PII detection
- **Multiple PII Types**: Email, SSN, Credit Card, Phone, Address, etc.
- **Confidence Scoring**: Confidence levels for detection accuracy
- **Real-time Scanning**: Live content scanning for PII

### **Compliance Management**
- **Framework Support**: GDPR, HIPAA, PCI-DSS, SOX compliance
- **Audit Logging**: Comprehensive audit trail for compliance
- **Report Generation**: Automated compliance reports
- **Policy Management**: Custom compliance policies
- **Regional Compliance**: Multi-region compliance support

### **Risk Assessment**
- **Risk Scoring**: Quantitative risk assessment algorithms
- **Mitigation Planning**: Automated mitigation plan generation
- **Risk Monitoring**: Continuous risk monitoring and alerting
- **Trend Analysis**: Historical risk trend analysis
- **ML Risk Models**: Machine learning-based risk prediction

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Data Protection Service                  │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   HTTP      │  │   Kafka     │  │  Database   │        │
│  │  Handlers   │  │  Producer   │  │ Connection  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Data      │  │    PII      │  │ Compliance  │        │
│  │Classification│  │ Detection   │  │  Service    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Risk      │  │ Repository  │  │   Config   │        │
│  │ Assessment  │  │    Layer    │  │ Management │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## 📋 API Endpoints

### **Data Classification**
- `POST /api/v1/classification/classify` - Classify data
- `GET /api/v1/classification/rules` - Get classification rules
- `POST /api/v1/classification/rules` - Create classification rule
- `PUT /api/v1/classification/rules/:id` - Update classification rule
- `DELETE /api/v1/classification/rules/:id` - Delete classification rule
- `GET /api/v1/classification/rules/:id` - Get specific rule
- `POST /api/v1/classification/rules/:id/enable` - Enable rule
- `POST /api/v1/classification/rules/:id/disable` - Disable rule
- `GET /api/v1/classification/report` - Get classification report

### **PII Detection**
- `POST /api/v1/pii/detect` - Detect PII in content
- `GET /api/v1/pii/patterns` - Get PII patterns
- `POST /api/v1/pii/patterns` - Create PII pattern
- `PUT /api/v1/pii/patterns/:id` - Update PII pattern
- `DELETE /api/v1/pii/patterns/:id` - Delete PII pattern
- `GET /api/v1/pii/scan` - Scan content for PII
- `GET /api/v1/pii/report` - Get PII detection report

### **Compliance Management**
- `GET /api/v1/compliance/frameworks` - Get compliance frameworks
- `GET /api/v1/compliance/frameworks/:id` - Get specific framework
- `POST /api/v1/compliance/frameworks` - Create framework
- `PUT /api/v1/compliance/frameworks/:id` - Update framework
- `DELETE /api/v1/compliance/frameworks/:id` - Delete framework
- `GET /api/v1/compliance/reports` - Get compliance reports
- `POST /api/v1/compliance/reports` - Create compliance report
- `GET /api/v1/compliance/reports/:id` - Get specific report
- `PUT /api/v1/compliance/reports/:id` - Update report
- `DELETE /api/v1/compliance/reports/:id` - Delete report
- `GET /api/v1/compliance/audit` - Get audit log

### **Risk Assessment**
- `POST /api/v1/risk/assess` - Assess risk
- `GET /api/v1/risk/scores` - Get risk scores
- `GET /api/v1/risk/scores/:id` - Get specific risk score
- `POST /api/v1/risk/mitigation` - Create mitigation plan
- `PUT /api/v1/risk/mitigation/:id` - Update mitigation plan
- `GET /api/v1/risk/mitigation/:id` - Get mitigation plan
- `DELETE /api/v1/risk/mitigation/:id` - Delete mitigation plan

### **Health & Monitoring**
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics endpoint

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Configuration**: Viper
- **Database**: PostgreSQL (via shared package)
- **Messaging**: Kafka (via shared package)
- **Logging**: Structured logging (via shared package)
- **Monitoring**: Prometheus metrics (via shared package)
- **Containerization**: Docker
- **Orchestration**: Kubernetes ready

## 🚀 Quick Start

### **Prerequisites**
- Go 1.21+
- PostgreSQL
- Kafka
- Docker (optional)

### **Local Development**

1. **Clone the repository**
   ```bash
   cd backend/services/data-protection
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=scopeapi
   export DB_PASSWORD=your_password
   export DB_NAME=scopeapi
   export KAFKA_BROKERS=localhost:9092
   export JWT_SECRET=your_jwt_secret
   ```

4. **Run the service**
   ```bash
   go run cmd/main.go
   ```

### **Docker Deployment**

1. **Build the image**
   ```bash
   docker build -t data-protection-service .
   ```

2. **Run the container**
   ```bash
   docker run -p 8084:8084 \
     -e DB_HOST=your_db_host \
     -e DB_PORT=5432 \
     -e DB_USER=your_user \
     -e DB_PASSWORD=your_password \
     -e DB_NAME=your_db \
     -e KAFKA_BROKERS=your_kafka_host:9092 \
     -e JWT_SECRET=your_secret \
     data-protection-service
   ```

## ⚙️ Configuration

The service supports configuration through:
- Environment variables
- Configuration files (YAML)
- Default values

### **Key Configuration Options**

```yaml
server:
  port: 8084
  host: "0.0.0.0"
  environment: "development"
  log_level: "info"

database:
  postgresql:
    host: "localhost"
    port: 5432
    user: "scopeapi"
    password: "your_password"
    database: "scopeapi"
    ssl_mode: "disable"

messaging:
  kafka:
    brokers: ["localhost:9092"]
    topic_prefix: "scopeapi.data-protection"

features:
  data_classification:
    enabled: true
    ml_enabled: true
    confidence_threshold: 0.8
  
  pii_detection:
    enabled: true
    pattern_matching: true
    ml_detection: true
    confidence_threshold: 0.85
  
  compliance:
    enabled: true
    gdpr: true
    hipaa: true
    pci_dss: true
    sox: true
  
  risk_assessment:
    enabled: true
    ml_enabled: true
    risk_threshold: 0.7
```

## 📊 Monitoring & Observability

### **Health Checks**
- Service health endpoint at `/health`
- Database connectivity checks
- Kafka producer health

### **Metrics**
- Prometheus metrics at `/metrics`
- Request counts and latencies
- Error rates and success rates
- Business metrics (classifications, PII detections, etc.)

### **Logging**
- Structured JSON logging
- Request/response logging
- Error logging with context
- Performance logging

## 🔒 Security Features

- **JWT Authentication**: Secure API access
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries
- **Rate Limiting**: API rate limiting (configurable)
- **Audit Logging**: Complete audit trail
- **Data Encryption**: Sensitive data encryption

## 🧪 Testing

### **Run Tests**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/services
```

### **Test Coverage**
- Unit tests for all services
- Integration tests for repositories
- HTTP handler tests
- Configuration tests

## 📈 Performance

### **Optimizations**
- Connection pooling for database
- Kafka producer batching
- In-memory caching for rules
- Efficient regex compilation
- ML model optimization

### **Benchmarks**
- Classification: ~1000 items/second
- PII Detection: ~500 items/second
- Risk Assessment: ~200 assessments/second

## 🔄 Deployment

### **Kubernetes**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: data-protection-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: data-protection-service
  template:
    metadata:
      labels:
        app: data-protection-service
    spec:
      containers:
      - name: data-protection-service
        image: data-protection-service:latest
        ports:
        - containerPort: 8084
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
```

### **Docker Compose**
```yaml
version: '3.8'
services:
  data-protection:
    build: .
    ports:
      - "8084:8084"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=scopeapi
      - DB_PASSWORD=scopeapi
      - DB_NAME=scopeapi
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - postgres
      - kafka
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the API specifications

## 🔮 Roadmap

### **Phase 1 (Current)**
- ✅ Core data classification
- ✅ PII detection
- ✅ Compliance management
- ✅ Risk assessment

### **Phase 2 (Next)**
- 🔄 Advanced ML models
- 🔄 Real-time streaming
- 🔄 Advanced analytics
- 🔄 Integration APIs

### **Phase 3 (Future)**
- 🔮 AI-powered insights
- 🔮 Predictive analytics
- 🔮 Advanced compliance
- 🔮 Global scale support
