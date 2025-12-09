# Microservices Comparison: Data Ingestion, Gateway Integration, and API Discovery

## 📋 Table of Contents

1. [Overview](#overview)
2. [Service Comparison](#service-comparison)
3. [Key Differences Summary](#key-differences-summary)
4. [Real-World Analogy](#real-world-analogy)
5. [Detailed Service Descriptions](#detailed-service-descriptions)
6. [Use Cases](#use-cases)
7. [How They Work Together](#how-they-work-together)
8. [Architecture Diagram](#architecture-diagram)
9. [Decision Matrix](#decision-matrix)

---

## Overview

This document explains the key differences between three critical microservices in the ScopeAPI platform:

- **API Discovery Service** - Discovers and catalogs API endpoints
- **Gateway Integration Service** - Manages API gateway infrastructure
- **Data Ingestion Service** - Processes runtime traffic data

While these services may seem related, they serve distinct purposes in the API security and management ecosystem.

---

## Service Comparison

### Quick Reference Table

| Aspect | API Discovery Service | Gateway Integration Service | Data Ingestion Service |
|--------|----------------------|----------------------------|----------------------|
| **Primary Purpose** | Find and catalog APIs | Manage gateway infrastructure | Process runtime traffic |
| **Focus** | Static API inventory | Gateway configuration | Real-time data processing |
| **When It Runs** | Periodic scans / On-demand | On configuration changes | Continuous (real-time) |
| **Input** | Infrastructure scans | Gateway management APIs | Traffic logs/events |
| **Output** | API catalog database | Gateway state/config | Normalized traffic stream |
| **Port** | 8080 | 8084 | 8085 |
| **Key Question** | "What APIs exist?" | "How are gateways configured?" | "What traffic is happening?" |
| **Data Type** | Endpoint metadata | Configuration data | Request/response data |
| **Frequency** | Scheduled/Manual | Event-driven | Continuous stream |
| **Dependencies** | PostgreSQL, Kafka | PostgreSQL, Kafka, Gateways | PostgreSQL, Kafka |

---

## Key Differences Summary

### The Three Core Distinctions

#### 1. **Temporal Focus**
- **API Discovery**: **Static/Periodic** - Works with snapshots of your infrastructure at specific points in time
- **Gateway Integration**: **Event-driven** - Responds to configuration changes and management operations
- **Data Ingestion**: **Continuous/Real-time** - Processes data as it happens, 24/7 streaming

#### 2. **Data Perspective**
- **API Discovery**: **What exists** - Focuses on the structure and metadata of APIs
- **Gateway Integration**: **How it's configured** - Focuses on gateway settings and infrastructure state
- **Data Ingestion**: **What's happening** - Focuses on actual runtime traffic and events

#### 3. **Operational Role**
- **API Discovery**: **Inventory Management** - Creates and maintains a catalog
- **Gateway Integration**: **Infrastructure Control** - Manages and configures gateways
- **Data Ingestion**: **Data Processing** - Transforms and streams traffic data

### Visual Comparison

```
┌─────────────────────────────────────────────────────────────┐
│                    API Discovery Service                    │
│  "What APIs exist in our infrastructure?"                   │
│                                                             │
│  Scans → Finds → Catalogs → Documents                       │
│  Result: API Inventory Database                             │
│  Frequency: Periodic / On-demand                            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                 Gateway Integration Service                 │
│  "How do we configure our gateways?"                        │
│                                                             │
│  Connects → Configures → Manages → Syncs                   │
│  Result: Gateway Configuration Management                  │
│  Frequency: Event-driven / On configuration change          │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                   Data Ingestion Service                    │
│  "What traffic is happening right now?"                     │
│                                                             │
│  Receives → Parses → Normalizes → Publishes                │
│  Result: Processed Traffic Data Stream                     │
│  Frequency: Continuous / Real-time                          │
└─────────────────────────────────────────────────────────────┘
```

### Quick Decision Guide

**Ask yourself:**
- **"What APIs do we have?"** → Use **API Discovery Service**
- **"How should I configure the gateway?"** → Use **Gateway Integration Service**
- **"What traffic is flowing through?"** → Use **Data Ingestion Service**

---

### 📋 Copyable Format - Key Differences Summary

```
================================================================================
KEY DIFFERENCES SUMMARY - COPYABLE FORMAT
================================================================================

1. TEMPORAL FOCUS
   - API Discovery:      Static/Periodic - Works with snapshots at specific points in time
   - Gateway Integration: Event-driven - Responds to configuration changes
   - Data Ingestion:      Continuous/Real-time - Processes data 24/7 as it happens

2. DATA PERSPECTIVE
   - API Discovery:      What exists - Structure and metadata of APIs
   - Gateway Integration: How it's configured - Gateway settings and infrastructure state
   - Data Ingestion:      What's happening - Actual runtime traffic and events

3. OPERATIONAL ROLE
   - API Discovery:      Inventory Management - Creates and maintains catalog
   - Gateway Integration: Infrastructure Control - Manages and configures gateways
   - Data Ingestion:      Data Processing - Transforms and streams traffic data

QUICK DECISION GUIDE:
   Question: "What APIs do we have?"
   Answer:   Use API Discovery Service

   Question: "How should I configure the gateway?"
   Answer:   Use Gateway Integration Service

   Question: "What traffic is flowing through?"
   Answer:   Use Data Ingestion Service

SERVICE COMPARISON:
   API Discovery Service:
   - Purpose: Find and catalog APIs
   - Focus: Static API inventory
   - When: Periodic scans / On-demand
   - Output: API catalog database
   - Key Question: "What APIs exist?"

   Gateway Integration Service:
   - Purpose: Manage gateway infrastructure
   - Focus: Gateway configuration
   - When: On configuration changes
   - Output: Gateway state/config
   - Key Question: "How are gateways configured?"

   Data Ingestion Service:
   - Purpose: Process runtime traffic
   - Focus: Real-time data processing
   - When: Continuous (real-time)
   - Output: Normalized traffic stream
   - Key Question: "What traffic is happening?"

REAL-WORLD ANALOGY:
   API Discovery      = Library catalog (lists what books/APIs exist)
   Gateway Integration = Library management system (controls organization and access)
   Data Ingestion     = Checkout system (processes who is checking out books/traffic)

SUMMARY:
   Discovery          = What exists (static inventory)
   Gateway Integration = How it's configured (infrastructure management)
   Data Ingestion     = What's happening (runtime traffic)
================================================================================
```

---

## Real-World Analogy

To better understand these services, think of them in terms of a **library system**:

### 📚 Library Analogy

#### **API Discovery Service = Library Catalog System**
- **Purpose**: Creates and maintains a catalog of all books (APIs) in the library
- **What it does**: 
  - Scans the library shelves to find all books
  - Creates a catalog card for each book with details (title, author, ISBN, location)
  - Tracks when new books arrive or old books are removed
  - Documents what books exist, but doesn't control access to them
- **Example**: 
  - "We have 1,000 books across 5 sections: Fiction, Non-fiction, Science, History, and Technology"
  - "A new book 'API Security Guide' was added to the Technology section"

#### **Gateway Integration Service = Library Management System**
- **Purpose**: Controls how the library operates and manages access
- **What it does**:
  - Sets up borrowing rules (rate limits)
  - Configures security systems (authentication)
  - Manages library hours and access policies
  - Controls how books are organized and accessed
  - Monitors library infrastructure (security cameras, checkout systems)
- **Example**:
  - "Configured checkout system: 5 books max per person, 2-week loan period"
  - "Set up security gates to prevent unauthorized access"
  - "Configured different access levels: public, member, staff"

#### **Data Ingestion Service = Checkout/Return System**
- **Purpose**: Processes actual transactions as they happen in real-time
- **What it does**:
  - Records every book checkout and return
  - Processes who is checking out what books, when, and how often
  - Normalizes transaction data (format, timestamps, user info)
  - Streams transaction data to analytics systems
- **Example**:
  - "John checked out 'API Design Patterns' at 2:30 PM"
  - "Processed 500 checkouts today"
  - "Streaming checkout data to analytics for usage patterns"

### 🏢 Office Building Analogy

Another way to think about it:

#### **API Discovery Service = Building Directory**
- Creates a directory of all offices (APIs) in the building
- Lists what each office does and who works there
- Updates when offices move or new ones are added
- **Doesn't control access** - just documents what exists

#### **Gateway Integration Service = Building Security & Access Control**
- Manages security systems (keycards, cameras)
- Configures access rules (who can enter which floors)
- Sets up elevators and door controls
- Monitors building infrastructure health
- **Controls how people access** the building

#### **Data Ingestion Service = Visitor Log System**
- Records every person entering and leaving
- Processes visitor data in real-time
- Tracks who visited which office, when, and for how long
- Streams visitor data for security analysis
- **Processes actual activity** as it happens

### 🚗 Traffic System Analogy

#### **API Discovery Service = Road Map**
- Maps all roads (APIs) in the city
- Documents road types, speed limits, and destinations
- Updates when new roads are built
- **Shows what roads exist**, not who's using them

#### **Gateway Integration Service = Traffic Control System**
- Manages traffic lights and signs
- Configures speed limits and lane rules
- Sets up toll booths and access controls
- Monitors traffic infrastructure
- **Controls how traffic flows**

#### **Data Ingestion Service = Traffic Monitoring System**
- Records every vehicle passing through
- Processes real-time traffic data
- Tracks speed, volume, and patterns
- Streams data to traffic analysis systems
- **Processes actual traffic** as it happens

---

## Detailed Service Descriptions

### 1. API Discovery Service

**Purpose**: Automatically discovers and catalogs API endpoints across your infrastructure.

#### What It Does

- **Scans Infrastructure**: Automatically crawls and scans your infrastructure to find API endpoints
- **Creates Inventory**: Builds a comprehensive catalog of all discovered APIs
- **Documents Metadata**: Extracts and stores API metadata (methods, parameters, schemas, versions)
- **Tracks Changes**: Monitors and detects changes in API endpoints over time
- **Analyzes Specifications**: Parses OpenAPI/Swagger specifications to understand API structure

#### Key Features

- Automated endpoint scanning
- API catalog management
- Change detection and monitoring
- Metadata extraction and analysis
- Real-time discovery status tracking
- Version tracking for APIs

#### Example Scenario

```
Discovery Process:
1. Scans gateway at https://api.company.com
2. Discovers endpoints:
   - GET    /api/v1/users
   - POST   /api/v1/users
   - GET    /api/v1/users/:id
   - PUT    /api/v1/users/:id
   - DELETE /api/v1/users/:id
3. Creates catalog entry:
   {
     "api_name": "User Management API",
     "version": "v1",
     "endpoints": 5,
     "base_url": "/api/v1/users"
   }
4. Stores in inventory database
```

#### API Endpoints

- `POST /api/v1/discovery/scan` - Start API discovery scan
- `GET /api/v1/discovery/status/:id` - Get discovery scan status
- `GET /api/v1/inventory/apis` - List all discovered APIs
- `GET /api/v1/inventory/apis/:id` - Get specific API details
- `POST /api/v1/endpoints/analyze` - Analyze endpoint metadata
- `GET /api/v1/endpoints/:id/metadata` - Get endpoint metadata

#### Use Cases

- **API Inventory Management**: Maintain a complete list of all APIs in your organization
- **Shadow API Detection**: Find APIs that aren't documented or managed
- **Compliance Auditing**: Document all APIs for compliance requirements
- **API Documentation**: Auto-generate API documentation from discovered endpoints
- **Change Management**: Track when APIs are added, modified, or removed

---

### 2. Gateway Integration Service

**Purpose**: Centralized management and monitoring of multiple API gateways.

#### What It Does

- **Connects to Gateways**: Establishes connections to Kong, NGINX, Traefik, Envoy, and HAProxy
- **Manages Configuration**: Creates, updates, and deploys gateway configurations
- **Synchronizes Settings**: Syncs configurations across multiple gateways and environments
- **Monitors Health**: Performs real-time health checks and status monitoring
- **Manages Credentials**: Securely stores and manages gateway authentication credentials
- **Version Control**: Maintains versioned configurations with rollback capabilities

#### Key Features

- Multi-gateway support (Kong, NGINX, Traefik, Envoy, HAProxy)
- Configuration versioning and rollback
- Real-time synchronization
- Health monitoring and status checks
- Secure credential management
- Gateway-specific operations

#### Example Scenario

```
Gateway Integration Process:
1. Connects to Kong gateway at 192.168.1.100:8001
2. Creates service:
   {
     "name": "user-service",
     "url": "http://backend:8080"
   }
3. Creates route:
   {
     "paths": ["/api/users"],
     "service": "user-service"
   }
4. Adds rate limiting plugin:
   {
     "name": "rate-limiting",
     "config": {
       "minute": 100,
       "hour": 1000
     }
   }
5. Deploys configuration to production gateway
6. Monitors gateway health status
```

#### Gateway-Specific Capabilities

| Gateway | Capabilities |
|---------|-------------|
| **Kong** | Services, Routes, Plugins, Consumers management |
| **NGINX** | Config management, Upstreams, Reload operations |
| **Traefik** | Providers, Middlewares, Routers management |
| **Envoy** | Clusters, Listeners, Filters management |
| **HAProxy** | Config management, Backends, Reload operations |

#### API Endpoints

**Integration Management:**
- `GET /api/v1/integrations` - List all integrations
- `POST /api/v1/integrations` - Create new integration
- `PUT /api/v1/integrations/:id` - Update integration
- `DELETE /api/v1/integrations/:id` - Delete integration
- `POST /api/v1/integrations/:id/test` - Test integration
- `POST /api/v1/integrations/:id/sync` - Sync integration

**Configuration Management:**
- `GET /api/v1/configs` - List configurations
- `POST /api/v1/configs` - Create configuration
- `POST /api/v1/configs/:id/deploy` - Deploy configuration
- `POST /api/v1/configs/:id/validate` - Validate configuration

**Gateway-Specific (Kong example):**
- `GET /api/v1/kong/services` - List Kong services
- `GET /api/v1/kong/routes` - List Kong routes
- `POST /api/v1/kong/plugins` - Create Kong plugin

#### Use Cases

- **Multi-Gateway Management**: Manage multiple gateways from a single interface
- **Configuration Deployment**: Deploy configurations across dev, staging, and production
- **Centralized Control**: Control all gateway settings from one place
- **Configuration Versioning**: Track and rollback gateway configurations
- **Gateway Health Monitoring**: Monitor the health of all gateway instances

---

### 3. Data Ingestion Service

**Purpose**: Processes and normalizes runtime traffic data from various sources.

#### What It Does

- **Receives Traffic Data**: Accepts API traffic data from multiple sources (gateways, proxies, logs)
- **Parses Data**: Parses traffic in various formats (JSON, XML, log formats, binary)
- **Normalizes Data**: Converts data into a standardized schema for downstream processing
- **Publishes to Kafka**: Sends processed traffic data to Kafka topics for real-time analysis
- **Tracks Statistics**: Maintains ingestion statistics and metrics
- **Handles Batches**: Processes both single requests and batch traffic data

#### Key Features

- Multiple format support (JSON, XML, logs, binary)
- Schema-based normalization
- Batch and streaming processing
- Kafka integration for downstream services
- Real-time status tracking
- Configurable parsing rules

#### Example Scenario

```
Data Ingestion Process:
1. Receives traffic data:
   {
     "method": "POST",
     "url": "/api/users",
     "ip_address": "192.168.1.50",
     "timestamp": "2024-01-15T10:30:00Z",
     "status_code": 201,
     "response_time_ms": 150,
     "request_body": "{\"name\":\"John\"}",
     "response_body": "{\"id\":123,\"name\":\"John\"}"
   }

2. Parses and validates data

3. Normalizes to standard schema:
   {
     "id": "req-12345",
     "timestamp": "2024-01-15T10:30:00Z",
     "source_ip": "192.168.1.50",
     "method": "POST",
     "endpoint": "/api/users",
     "status_code": 201,
     "response_time": 150,
     "request_size": 15,
     "response_size": 25
   }

4. Publishes to Kafka topic: "api-traffic"
5. Updates ingestion statistics
```

#### Processing Flow

```
Traffic Source → Data Ingestion Service
                    ↓
              Parse & Validate
                    ↓
              Normalize Schema
                    ↓
              Publish to Kafka
                    ↓
         Threat Detection Service
         Data Protection Service
         Analytics Services
```

#### API Endpoints

**Data Ingestion:**
- `POST /api/v1/ingestion/traffic` - Ingest single traffic data
- `POST /api/v1/ingestion/batch` - Ingest batch traffic data
- `GET /api/v1/ingestion/status/:id` - Get ingestion status
- `GET /api/v1/ingestion/stats` - Get ingestion statistics

**Parser:**
- `POST /api/v1/parser/parse` - Parse data
- `GET /api/v1/parser/formats` - Get supported formats
- `POST /api/v1/parser/validate` - Validate format

**Normalizer:**
- `POST /api/v1/normalizer/normalize` - Normalize data
- `GET /api/v1/normalizer/schemas` - Get normalization schemas
- `POST /api/v1/normalizer/schema` - Create schema

#### Use Cases

- **Traffic Analysis**: Process traffic for security analysis and threat detection
- **Real-time Monitoring**: Stream traffic data for real-time monitoring dashboards
- **Log Aggregation**: Collect and normalize logs from multiple sources
- **Data Pipeline**: Feed processed traffic data to analytics and ML services
- **Compliance Logging**: Process and store traffic for compliance requirements

---

## Use Cases

### When to Use API Discovery Service

✅ **Use when you need to:**
- Build an inventory of all APIs in your organization
- Find undocumented or "shadow" APIs
- Track API changes over time
- Generate API documentation automatically
- Audit APIs for compliance
- Understand your API landscape

❌ **Don't use for:**
- Processing real-time traffic
- Configuring gateways
- Managing runtime requests

### When to Use Gateway Integration Service

✅ **Use when you need to:**
- Manage multiple API gateways from one place
- Deploy configurations across environments
- Monitor gateway health
- Version control gateway configurations
- Synchronize settings between gateways
- Configure routes, plugins, and middleware

❌ **Don't use for:**
- Discovering what APIs exist
- Processing traffic data
- Analyzing request/response content

### When to Use Data Ingestion Service

✅ **Use when you need to:**
- Process real-time API traffic
- Normalize traffic data from multiple sources
- Feed data to analytics or ML services
- Stream traffic for real-time monitoring
- Aggregate logs from different sources
- Process traffic for security analysis

❌ **Don't use for:**
- Discovering API endpoints
- Configuring gateways
- Managing gateway settings

---

## How They Work Together

These three services work together to provide a complete API security and management solution:

```
┌─────────────────────────────────────────────────────────────┐
│                    Complete API Lifecycle                   │
└─────────────────────────────────────────────────────────────┘

1. DISCOVERY PHASE
   ┌─────────────────────┐
   │ API Discovery        │ → Finds all APIs in infrastructure
   │ Service              │ → Creates inventory catalog
   └─────────────────────┘
           ↓
   "We have 50 endpoints across 5 APIs"

2. CONFIGURATION PHASE
   ┌─────────────────────┐
   │ Gateway Integration  │ → Configures gateways
   │ Service              │ → Sets up routes, plugins, security
   └─────────────────────┘
           ↓
   "Kong configured with rate limiting on /api/users"

3. RUNTIME PHASE
   ┌─────────────────────┐
   │ Data Ingestion      │ → Processes traffic
   │ Service              │ → Normalizes and streams data
   └─────────────────────┘
           ↓
   "Processing 10,000 requests/hour through gateways"
```

### Integration Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    Service Integration Flow                  │
└──────────────────────────────────────────────────────────────┘

API Discovery Service
    │
    │ Discovers endpoints from gateways
    ↓
Gateway Integration Service
    │
    │ Manages gateway configurations
    │ Sets up routes, security, plugins
    ↓
Traffic flows through configured gateways
    │
    │ Real-time API requests/responses
    ↓
Data Ingestion Service
    │
    │ Processes and normalizes traffic
    │ Publishes to Kafka
    ↓
Downstream Services
    ├─ Threat Detection Service (analyzes for threats)
    ├─ Data Protection Service (scans for PII)
    └─ Analytics Services (monitoring, dashboards)
```

### Example: Complete Workflow

**Scenario**: Setting up security for a new API endpoint

1. **API Discovery Service** discovers the new endpoint:
   ```
   POST /api/v1/orders
   ```

2. **Gateway Integration Service** configures the gateway:
   ```
   - Creates route: /api/v1/orders → order-service:8080
   - Adds rate limiting: 100 requests/minute
   - Adds authentication plugin
   - Configures CORS
   ```

3. **Data Ingestion Service** processes traffic:
   ```
   - Receives requests to /api/v1/orders
   - Normalizes traffic data
   - Publishes to Kafka
   - Threat Detection analyzes for attacks
   - Data Protection scans for sensitive data
   ```

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    ScopeAPI Architecture                        │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────┐      ┌──────────────────┐      ┌──────────────────┐
│   API Discovery  │      │ Gateway          │      │ Data Ingestion  │
│   Service        │      │ Integration      │      │ Service         │
│                  │      │ Service          │      │                 │
│ • Scans APIs     │      │ • Manages        │      │ • Processes     │
│ • Creates        │      │   Gateways       │      │   Traffic       │
│   Inventory      │      │ • Configures     │      │ • Normalizes    │
│ • Tracks Changes │      │   Routes         │      │ • Streams Data  │
└──────────────────┘      └──────────────────┘      └──────────────────┘
         │                        │                          │
         │                        │                          │
         └────────────────────────┼──────────────────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │      PostgreSQL            │
                    │   (Shared Database)       │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │      Kafka                │
                    │   (Message Queue)         │
                    └─────────────┬─────────────┘
                                  │
         ┌────────────────────────┼────────────────────────┐
         │                        │                          │
┌────────▼────────┐    ┌─────────▼─────────┐    ┌─────────▼─────────┐
│ Threat          │    │ Data Protection   │    │ Analytics         │
│ Detection       │    │ Service           │    │ Services           │
│ Service         │    │                   │    │                    │
└─────────────────┘    └───────────────────┘    └────────────────────┘
```

---

## Decision Matrix

Use this matrix to determine which service to use for your specific need:

| Need | API Discovery | Gateway Integration | Data Ingestion |
|------|---------------|---------------------|----------------|
| Find all APIs in infrastructure | ✅ | ❌ | ❌ |
| Configure rate limiting on gateway | ❌ | ✅ | ❌ |
| Process real-time traffic logs | ❌ | ❌ | ✅ |
| Create API inventory | ✅ | ❌ | ❌ |
| Deploy gateway configuration | ❌ | ✅ | ❌ |
| Normalize traffic data | ❌ | ❌ | ✅ |
| Track API changes over time | ✅ | ❌ | ❌ |
| Monitor gateway health | ❌ | ✅ | ❌ |
| Stream traffic to analytics | ❌ | ❌ | ✅ |
| Manage multiple gateways | ❌ | ✅ | ❌ |
| Find undocumented APIs | ✅ | ❌ | ❌ |
| Process request/response data | ❌ | ❌ | ✅ |

---

## Summary

### API Discovery Service
- **Focus**: Static inventory - "What APIs exist?"
- **When**: Periodic scans or on-demand discovery
- **Output**: API catalog database
- **Use Case**: Building and maintaining API inventory

### Gateway Integration Service
- **Focus**: Infrastructure management - "How are gateways configured?"
- **When**: Configuration changes or management operations
- **Output**: Gateway state and configuration
- **Use Case**: Managing and configuring API gateways

### Data Ingestion Service
- **Focus**: Runtime processing - "What traffic is happening?"
- **When**: Continuous real-time processing
- **Output**: Normalized traffic data stream
- **Use Case**: Processing and analyzing runtime traffic

---

## Related Documentation

- [Backend README](../backend/README.md) - Overall backend architecture
- [API Discovery README](../backend/services/api-discovery/README.md) - Detailed API Discovery docs
- [Gateway Integration README](../backend/services/gateway-integration/README.md) - Detailed Gateway Integration docs
- [Data Ingestion Service](../backend/services/data-ingestion/) - Data Ingestion implementation
- [Architecture Guide](./ARCHITECTURE.md) - System architecture overview


**Last Updated**: 2024-01-15  
**Version**: 1.0  
**Maintained by**: ScopeAPI Development Team

