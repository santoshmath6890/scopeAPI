# Security Services Comparison: Threat Detection, Data Protection, and Attack Blocking

## 📋 Table of Contents

1. [Overview](#overview)
2. [Service Comparison](#service-comparison)
3. [Key Differences Summary](#key-differences-summary)
4. [Real-World Analogy](#real-world-analogy)
5. [Detailed Service Descriptions](#detailed-service-descriptions)
6. [Use Cases](#use-cases)
7. [How They Work Together](#how-they-work-together)
8. [Security Workflow](#security-workflow)
9. [Decision Matrix](#decision-matrix)

---

## Overview

This document explains the key differences between three critical security microservices in the ScopeAPI platform:

- **Threat Detection Service** - Identifies and analyzes security threats
- **Data Protection Service** - Protects sensitive data and ensures compliance
- **Attack Blocking Service** - Prevents and blocks malicious requests in real-time

While all three services focus on security, they serve distinct roles in the security ecosystem: **Detection**, **Protection**, and **Prevention**.

---

## Service Comparison

### Quick Reference Table

| Aspect | Threat Detection Service | Data Protection Service | Attack Blocking Service |
|--------|-------------------------|------------------------|------------------------|
| **Primary Purpose** | Identify security threats | Protect sensitive data | Block malicious requests |
| **Focus** | Threat analysis and detection | Data privacy and compliance | Real-time prevention |
| **When It Runs** | Continuous analysis | On data processing | Real-time per request |
| **Input** | Traffic data, security events | Data content, API responses | Incoming requests |
| **Output** | Threat alerts, analysis reports | PII findings, compliance reports | Block/Allow decisions |
| **Port** | 8081 | 8082 | 8083 |
| **Key Question** | "Is this a threat?" | "Is this data sensitive?" | "Should I block this request?" |
| **Action Type** | Analysis & Alerting | Classification & Reporting | Enforcement & Blocking |
| **Response Time** | Near real-time (seconds) | Near real-time (seconds) | Real-time (milliseconds) |
| **Dependencies** | PostgreSQL, Kafka, Elasticsearch | PostgreSQL, Kafka | PostgreSQL, Redis, Kafka |

---

## Key Differences Summary

### The Three Core Distinctions

#### 1. **Security Role**
- **Threat Detection**: **Detection & Analysis** - Identifies what threats exist
- **Data Protection**: **Classification & Compliance** - Identifies what data is sensitive
- **Attack Blocking**: **Prevention & Enforcement** - Stops threats from reaching systems

#### 2. **Focus Area**
- **Threat Detection**: **Security Threats** - SQL injection, XSS, DDoS, brute force, etc.
- **Data Protection**: **Data Sensitivity** - PII, credit cards, SSNs, compliance violations
- **Attack Blocking**: **Request Blocking** - IP blocking, rate limiting, geo-blocking

#### 3. **Operational Mode**
- **Threat Detection**: **Analytical** - Analyzes and reports threats
- **Data Protection**: **Protective** - Classifies and protects data
- **Attack Blocking**: **Enforcement** - Actively blocks requests

### Visual Comparison

```
┌─────────────────────────────────────────────────────────────┐
│                  Threat Detection Service                   │
│  "Is this traffic a security threat?"                      │
│                                                             │
│  Analyzes → Detects → Reports → Alerts                      │
│  Result: Threat Analysis & Alerts                          │
│  Mode: Analytical (Detection)                              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  Data Protection Service                   │
│  "Is this data sensitive or non-compliant?"                │
│                                                             │
│  Scans → Classifies → Detects PII → Reports                 │
│  Result: Data Classification & Compliance Reports         │
│  Mode: Protective (Classification)                         │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  Attack Blocking Service                   │
│  "Should I block this request?"                             │
│                                                             │
│  Evaluates → Decides → Blocks → Enforces                    │
│  Result: Block/Allow Decision                              │
│  Mode: Enforcement (Prevention)                            │
└─────────────────────────────────────────────────────────────┘
```

### Quick Decision Guide

**Ask yourself:**
- **"Is this traffic malicious or suspicious?"** → Use **Threat Detection Service**
- **"Does this data contain sensitive information?"** → Use **Data Protection Service**
- **"Should I block this request right now?"** → Use **Attack Blocking Service**

---

### 📋 Copyable Format - Key Differences Summary

```
================================================================================
KEY DIFFERENCES SUMMARY - COPYABLE FORMAT
================================================================================

1. SECURITY ROLE
   - Threat Detection: Detection & Analysis - Identifies what threats exist
   - Data Protection: Classification & Compliance - Identifies sensitive data
   - Attack Blocking: Prevention & Enforcement - Stops threats from reaching systems

2. FOCUS AREA
   - Threat Detection: Security Threats - SQL injection, XSS, DDoS, brute force
   - Data Protection: Data Sensitivity - PII, credit cards, SSNs, compliance
   - Attack Blocking: Request Blocking - IP blocking, rate limiting, geo-blocking

3. OPERATIONAL MODE
   - Threat Detection: Analytical - Analyzes and reports threats
   - Data Protection: Protective - Classifies and protects data
   - Attack Blocking: Enforcement - Actively blocks requests

QUICK DECISION GUIDE:
   Question: "Is this traffic malicious or suspicious?"
   Answer:   Use Threat Detection Service

   Question: "Does this data contain sensitive information?"
   Answer:   Use Data Protection Service

   Question: "Should I block this request right now?"
   Answer:   Use Attack Blocking Service

SERVICE COMPARISON:
   Threat Detection Service:
   - Purpose: Identify security threats
   - Focus: Threat analysis and detection
   - When: Continuous analysis
   - Output: Threat alerts, analysis reports
   - Key Question: "Is this a threat?"

   Data Protection Service:
   - Purpose: Protect sensitive data
   - Focus: Data privacy and compliance
   - When: On data processing
   - Output: PII findings, compliance reports
   - Key Question: "Is this data sensitive?"

   Attack Blocking Service:
   - Purpose: Block malicious requests
   - Focus: Real-time prevention
   - When: Real-time per request
   - Output: Block/Allow decisions
   - Key Question: "Should I block this request?"

REAL-WORLD ANALOGY:
   Threat Detection      = Security Guard (identifies suspicious activity)
   Data Protection       = Vault Inspector (identifies valuable/sensitive items)
   Attack Blocking       = Bouncer (prevents entry at the door)

SUMMARY:
   Threat Detection  = What threats exist (detection & analysis)
   Data Protection   = What data is sensitive (classification & compliance)
   Attack Blocking   = Block threats now (prevention & enforcement)
================================================================================
```

---

## Real-World Analogy

### 🏛️ Bank Security System Analogy

#### **Threat Detection Service = Security Monitoring Center**
- **Purpose**: Monitors and analyzes all activity to identify suspicious behavior
- **What it does**:
  - Watches security cameras and monitors
  - Analyzes patterns of behavior
  - Identifies potential threats (suspicious individuals, unusual activity)
  - Generates alerts and reports
  - Uses AI/ML to detect anomalies
- **Example**:
  - "Detected suspicious behavior: Person loitering near vault"
  - "Anomaly detected: Unusual access pattern from user account"
  - "Threat identified: Multiple failed login attempts"

#### **Data Protection Service = Vault Inspector**
- **Purpose**: Identifies and protects valuable/sensitive items
- **What it does**:
  - Inspects what's being stored or transferred
  - Classifies items by value/sensitivity (cash, documents, jewelry)
  - Ensures compliance with regulations (insurance requirements, storage rules)
  - Reports on sensitive data handling
- **Example**:
  - "Detected PII: Credit card number in transaction log"
  - "Data classified as: Confidential - Customer financial records"
  - "Compliance violation: GDPR - Personal data not properly encrypted"

#### **Attack Blocking Service = Security Guard at Entrance**
- **Purpose**: Prevents unauthorized access at the point of entry
- **What it does**:
  - Checks IDs and access credentials
  - Blocks known threats immediately
  - Enforces access rules (rate limiting, time restrictions)
  - Prevents entry before any damage can occur
- **Example**:
  - "Blocked: IP address on blacklist"
  - "Blocked: Rate limit exceeded (100 requests/minute)"
  - "Blocked: Request from blocked geographic region"

### 🏥 Hospital Security Analogy

#### **Threat Detection Service = Medical Monitoring System**
- Monitors patient vitals and identifies anomalies
- Detects unusual patterns that could indicate problems
- Alerts medical staff to potential issues
- **Focus**: "Is something wrong?" (Detection)

#### **Data Protection Service = HIPAA Compliance Officer**
- Ensures patient data (PHI) is properly handled
- Classifies data by sensitivity level
- Ensures compliance with HIPAA regulations
- **Focus**: "Is this data properly protected?" (Compliance)

#### **Attack Blocking Service = Hospital Security**
- Prevents unauthorized access to restricted areas
- Blocks suspicious individuals at entrances
- Enforces access policies in real-time
- **Focus**: "Should this person be allowed in?" (Prevention)

### 🚗 Traffic System Analogy

#### **Threat Detection Service = Traffic Monitoring System**
- Monitors traffic patterns
- Identifies accidents, congestion, suspicious vehicles
- Analyzes behavior to detect threats
- **Focus**: "Is there a problem?" (Detection)

#### **Data Protection Service = Cargo Inspector**
- Inspects what's being transported
- Classifies cargo by type and value
- Ensures compliance with transport regulations
- **Focus**: "Is the cargo properly handled?" (Protection)

#### **Attack Blocking Service = Traffic Police**
- Enforces traffic rules at checkpoints
- Blocks vehicles that violate rules
- Prevents dangerous vehicles from proceeding
- **Focus**: "Should this vehicle be stopped?" (Enforcement)

---

## Detailed Service Descriptions

### 1. Threat Detection Service

**Purpose**: Identifies and analyzes security threats in real-time using multiple detection methodologies.

#### What It Does

- **Threat Detection**: Identifies known attack patterns (SQL injection, XSS, DDoS, brute force, path traversal, command injection)
- **Anomaly Detection**: Detects deviations from normal traffic patterns using statistical and ML methods
- **Behavioral Analysis**: Analyzes user and entity behavior patterns to identify suspicious activities
- **Signature Matching**: Matches traffic against known threat signatures
- **Real-time Processing**: Processes incoming traffic and security events via Kafka
- **Threat Intelligence**: Integrates with threat feeds and reputation services

#### Key Features

- Multiple detection methodologies (signature, anomaly, behavioral)
- ML-based anomaly detection
- Real-time Kafka-based processing
- Threat scoring and severity classification
- Baseline profile creation for behavioral analysis
- Comprehensive threat reporting

#### Detection Capabilities

1. **Signature-based Detection**
   - SQL Injection patterns
   - XSS (Cross-Site Scripting) patterns
   - Path traversal attempts
   - Command injection patterns

2. **Anomaly Detection**
   - Traffic volume anomalies
   - Response time anomalies
   - Request pattern anomalies
   - Geolocation anomalies
   - Statistical anomaly detection
   - ML-based anomaly detection

3. **Behavioral Analysis**
   - Access pattern analysis
   - Usage pattern analysis
   - Timing pattern analysis
   - Sequence pattern analysis
   - Location pattern analysis
   - Risk scoring

#### Example Scenario

```
Threat Detection Process:
1. Receives traffic data:
   {
     "method": "POST",
     "url": "/api/users",
     "body": "username=admin'; DROP TABLE users;--",
     "ip_address": "192.168.1.100"
   }

2. Analyzes using multiple methods:
   - Signature detection: Matches SQL injection pattern
   - Anomaly detection: Unusual request pattern
   - Behavioral analysis: User behavior deviates from baseline

3. Generates threat report:
   {
     "threat_type": "SQL Injection",
     "severity": "high",
     "confidence": 0.95,
     "risk_score": 8.5,
     "indicators": ["DROP TABLE", "SQL comment", "Unusual pattern"]
   }

4. Publishes alert to Kafka
5. Stores threat in database
```

#### API Endpoints

**Threat Detection:**
- `POST /api/v1/threats/analyze` - Analyze traffic for threats
- `GET /api/v1/threats` - List threats with filtering
- `GET /api/v1/threats/:id` - Get specific threat details
- `PUT /api/v1/threats/:id/status` - Update threat status

**Anomaly Detection:**
- `POST /api/v1/anomalies/detect` - Detect anomalies in traffic
- `GET /api/v1/anomalies` - List anomalies with filtering

**Behavioral Analysis:**
- `POST /api/v1/behavioral/analyze` - Analyze behavior patterns
- `GET /api/v1/behavioral/patterns` - List behavior patterns

**Signature Management:**
- `POST /api/v1/signatures/detect` - Detect signatures in traffic
- `GET /api/v1/signatures` - List threat signatures

#### Use Cases

- **Security Monitoring**: Continuous monitoring for security threats
- **Incident Detection**: Early detection of security incidents
- **Threat Intelligence**: Integration with threat feeds
- **Behavioral Analysis**: Identifying suspicious user behavior
- **Anomaly Detection**: Detecting unusual traffic patterns

---

### 2. Data Protection Service

**Purpose**: Protects sensitive data through classification, PII detection, compliance management, and risk assessment.

#### What It Does

- **Data Classification**: Categorizes data by sensitivity (Public, Internal, Confidential, Restricted, Top Secret)
- **PII Detection**: Identifies personally identifiable information (emails, SSNs, credit cards, phone numbers, addresses)
- **Compliance Management**: Ensures compliance with GDPR, HIPAA, PCI-DSS, SOX
- **Risk Assessment**: Calculates risk scores and generates mitigation plans
- **Pattern Matching**: Uses regex and ML-based detection
- **Audit Logging**: Comprehensive audit trail for compliance

#### Key Features

- Rule-based and ML-powered classification
- Multiple PII type detection with confidence scoring
- Multi-framework compliance support (GDPR, HIPAA, PCI-DSS, SOX)
- Automated compliance reporting
- Risk scoring and mitigation planning
- Real-time scanning capabilities

#### Protection Capabilities

1. **Data Classification**
   - Rule-based classification
   - ML-powered classification
   - Pattern matching
   - Context analysis
   - Classification levels: Public, Internal, Confidential, Restricted, Top Secret

2. **PII Detection**
   - Email addresses
   - Social Security Numbers (SSN)
   - Credit card numbers
   - Phone numbers
   - Physical addresses
   - Passport numbers
   - Driver's license numbers
   - Bank account numbers
   - Health records
   - Biometric data

3. **Compliance Management**
   - GDPR compliance
   - HIPAA compliance
   - PCI-DSS compliance
   - SOX compliance
   - Custom compliance policies
   - Regional compliance support

4. **Risk Assessment**
   - Quantitative risk scoring
   - Mitigation plan generation
   - Continuous risk monitoring
   - Historical trend analysis
   - ML-based risk prediction

#### Example Scenario

```
Data Protection Process:
1. Receives API response data:
   {
     "user": {
       "name": "John Doe",
       "email": "john.doe@example.com",
       "ssn": "123-45-6789",
       "credit_card": "4532-1234-5678-9010"
     }
   }

2. Scans for PII:
   - Detects email: john.doe@example.com
   - Detects SSN: 123-45-6789
   - Detects credit card: 4532-1234-5678-9010

3. Classifies data:
   - Classification: Confidential
   - Risk Score: 8.5
   - Compliance Issues: PCI-DSS violation (credit card not encrypted)

4. Generates report:
   {
     "pii_findings": [
       {"type": "email", "location": "user.email", "risk": "medium"},
       {"type": "ssn", "location": "user.ssn", "risk": "high"},
       {"type": "credit_card", "location": "user.credit_card", "risk": "critical"}
     ],
     "compliance_violations": ["PCI-DSS"],
     "recommendations": ["Encrypt credit card data", "Mask SSN in logs"]
   }
```

#### API Endpoints

**Data Classification:**
- `POST /api/v1/classification/classify` - Classify data
- `GET /api/v1/classification/rules` - Get classification rules
- `POST /api/v1/classification/rules` - Create classification rule

**PII Detection:**
- `POST /api/v1/pii/detect` - Detect PII in content
- `GET /api/v1/pii/patterns` - Get PII patterns
- `GET /api/v1/pii/report` - Get PII detection report

**Compliance Management:**
- `GET /api/v1/compliance/frameworks` - Get compliance frameworks
- `GET /api/v1/compliance/reports` - Get compliance reports
- `POST /api/v1/compliance/reports` - Create compliance report

**Risk Assessment:**
- `POST /api/v1/risk/assess` - Assess risk
- `GET /api/v1/risk/scores` - Get risk scores
- `POST /api/v1/risk/mitigation` - Create mitigation plan

#### Use Cases

- **PII Protection**: Detect and protect personally identifiable information
- **Compliance Auditing**: Ensure compliance with regulations (GDPR, HIPAA, PCI-DSS)
- **Data Classification**: Automatically classify data by sensitivity
- **Risk Assessment**: Assess and mitigate data-related risks
- **Audit Logging**: Maintain compliance audit trails

---

### 3. Attack Blocking Service

**Purpose**: Prevents and blocks malicious requests in real-time before they reach backend systems.

#### What It Does

- **Real-time Blocking**: Evaluates and blocks requests in milliseconds
- **IP Management**: Maintains whitelists and blacklists
- **Rate Limiting**: Enforces request rate limits per IP/endpoint
- **Geo-blocking**: Blocks requests from specific geographic regions
- **Signature Detection**: Matches attack signatures in request content
- **Anomaly-based Blocking**: Blocks requests matching anomaly patterns
- **Cloud Intelligence**: Integrates threat feeds and reputation scores
- **Custom Rules**: Applies user-defined blocking rules

#### Key Features

- Multi-layered blocking (IP lists, rate limits, signatures, anomalies)
- Active block management with expiration
- Real-time processing with low latency (milliseconds)
- Event publishing to Kafka for downstream services
- Automatic cleanup of expired blocks
- Cloud threat intelligence integration

#### Blocking Mechanisms

1. **IP-based Blocking**
   - IP whitelist (always allow)
   - IP blacklist (always block)
   - Active blocks (temporary blocks)

2. **Rate Limiting**
   - Requests per minute/hour
   - Per IP address
   - Per endpoint
   - Per API key

3. **Geo-blocking**
   - Block by country
   - Block by region
   - Allow specific countries only

4. **Signature-based Blocking**
   - SQL injection patterns
   - XSS patterns
   - Path traversal patterns
   - Command injection patterns

5. **Anomaly-based Blocking**
   - Unusual request size
   - Unusual request frequency
   - Unusual headers
   - Unusual patterns

6. **Cloud Intelligence**
   - Threat feed integration
   - Reputation score checking
   - Known malicious IPs

#### Blocking Decision Flow

```
Request arrives
    ↓
1. Check IP whitelist → If whitelisted: ALLOW
    ↓
2. Check IP blacklist → If blacklisted: BLOCK
    ↓
3. Check active blocks → If blocked: BLOCK
    ↓
4. Check rate limits → If exceeded: BLOCK
    ↓
5. Check geo-blocking → If blocked region: BLOCK
    ↓
6. Check signatures → If malicious pattern: BLOCK
    ↓
7. Check anomalies → If anomalous: BLOCK
    ↓
8. Check cloud intelligence → If threat: BLOCK
    ↓
9. Check custom rules → If rule matches: BLOCK
    ↓
10. ALLOW (passed all checks)
```

#### Example Scenario

```
Attack Blocking Process:
1. Request arrives:
   {
     "ip_address": "192.168.1.100",
     "method": "POST",
     "url": "/api/users",
     "body": "username=admin'; DROP TABLE users;--"
   }

2. Evaluation process:
   - IP not whitelisted ✓
   - IP not blacklisted ✓
   - No active blocks ✓
   - Rate limit: 50/100 requests/minute ✓
   - Geo: Allowed region ✓
   - Signature: SQL injection detected ✗

3. Decision: BLOCK
   {
     "action": "block",
     "reason": "Malicious signature detected: SQL Injection Pattern 1",
     "block_id": "block-12345",
     "blocked_until": "2024-01-15T11:00:00Z"
   }

4. Creates active block
5. Publishes blocking event to Kafka
6. Returns 403 Forbidden to client
```

#### API Endpoints

**Blocking Operations:**
- `POST /api/v1/blocking/process` - Process request for blocking decision
- `GET /api/v1/blocking/active` - Get active blocks
- `POST /api/v1/blocking/unblock` - Unblock IP address

**Blocking Rules:**
- `POST /api/v1/blocking/rules` - Create blocking rule
- `GET /api/v1/blocking/rules` - List blocking rules
- `PUT /api/v1/blocking/rules/:id` - Update blocking rule

**IP Management:**
- `POST /api/v1/blocking/whitelist` - Add IP to whitelist
- `POST /api/v1/blocking/blacklist` - Add IP to blacklist
- `DELETE /api/v1/blocking/whitelist/:ip` - Remove from whitelist

**Statistics:**
- `GET /api/v1/blocking/stats` - Get blocking statistics
- `GET /api/v1/blocking/health` - Get blocking health status

#### Use Cases

- **Real-time Protection**: Block attacks before they reach backend
- **DDoS Mitigation**: Rate limiting and IP blocking for DDoS attacks
- **Brute Force Prevention**: Block repeated failed login attempts
- **Geo-blocking**: Restrict access by geographic region
- **Custom Security Rules**: Enforce organization-specific security policies

---

## Use Cases

### When to Use Threat Detection Service

✅ **Use when you need to:**
- Identify security threats in traffic
- Detect SQL injection, XSS, or other attacks
- Analyze traffic patterns for anomalies
- Monitor user behavior for suspicious activity
- Generate security alerts and reports
- Integrate with threat intelligence feeds

❌ **Don't use for:**
- Blocking requests in real-time
- Detecting PII or sensitive data
- Enforcing access control
- Rate limiting

### When to Use Data Protection Service

✅ **Use when you need to:**
- Detect PII in data (emails, SSNs, credit cards)
- Classify data by sensitivity level
- Ensure GDPR, HIPAA, or PCI-DSS compliance
- Assess data-related risks
- Generate compliance reports
- Audit data handling practices

❌ **Don't use for:**
- Detecting security attacks
- Blocking malicious requests
- Rate limiting
- IP management

### When to Use Attack Blocking Service

✅ **Use when you need to:**
- Block malicious requests in real-time
- Enforce rate limiting
- Manage IP whitelists/blacklists
- Block requests from specific regions
- Prevent attacks before they reach backend
- Enforce custom security rules

❌ **Don't use for:**
- Analyzing threats (use Threat Detection)
- Detecting PII (use Data Protection)
- Generating compliance reports
- Behavioral analysis

---

## How They Work Together

These three services work together to provide comprehensive security:

```
┌─────────────────────────────────────────────────────────────┐
│                    Complete Security Flow                   │
└─────────────────────────────────────────────────────────────┘

1. REQUEST ARRIVES
   ┌─────────────────────┐
   │ Attack Blocking      │ → Evaluates request
   │ Service              │ → Blocks if malicious
   └─────────────────────┘
           ↓
   Request passes blocking checks

2. TRAFFIC PROCESSING
   ┌─────────────────────┐
   │ Threat Detection     │ → Analyzes traffic
   │ Service              │ → Detects threats
   └─────────────────────┘
           ↓
   Threat detected → Alert generated

3. DATA PROTECTION
   ┌─────────────────────┐
   │ Data Protection      │ → Scans data
   │ Service              │ → Detects PII
   └─────────────────────┘
           ↓
   PII detected → Compliance report
```

### Integration Flow

```
Incoming Request
    │
    ↓
┌─────────────────────┐
│ Attack Blocking      │ → Real-time evaluation
│ Service              │ → Block/Allow decision
└─────────────────────┘
    │
    ↓ (if allowed)
Request reaches backend
    │
    ↓
Traffic data flows to:
    │
    ├─→ Threat Detection Service
    │   → Analyzes for security threats
    │   → Generates alerts
    │
    └─→ Data Protection Service
        → Scans for PII
        → Classifies data
        → Checks compliance
    │
    ↓
If threat detected:
    │
    └─→ Attack Blocking Service
        → Updates blacklist
        → Creates active block
        → Prevents future requests
```

### Example: Complete Security Workflow

**Scenario**: A malicious request attempting SQL injection while also containing PII

1. **Attack Blocking Service** evaluates the request:
   ```
   - Checks IP: Not blacklisted
   - Checks rate limit: Within limits
   - Checks signature: SQL injection pattern detected
   - Decision: BLOCK
   - Creates active block for IP
   ```

2. **Threat Detection Service** (if request had passed):
   ```
   - Analyzes traffic pattern
   - Detects SQL injection attempt
   - Generates threat alert
   - Stores threat in database
   - Publishes to Kafka
   ```

3. **Data Protection Service** (if analyzing response data):
   ```
   - Scans response for PII
   - Detects credit card number
   - Classifies as Confidential
   - Flags PCI-DSS compliance issue
   - Generates compliance report
   ```

---

## Security Workflow

### Complete Security Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│                    Security Pipeline                            │
└─────────────────────────────────────────────────────────────────┘

Request → Attack Blocking (Prevention)
            ↓ (if allowed)
         Backend Processing
            ↓
         Response Data
            ↓
    ┌───────┴───────┐
    │               │
Threat Detection  Data Protection
(Analysis)        (Classification)
    │               │
    └───────┬───────┘
            ↓
    Security Dashboard
    (Alerts & Reports)
```

### Real-time Security Flow

```
1. REQUEST ARRIVAL
   ┌─────────────────────┐
   │ Attack Blocking      │
   │ Service              │
   │                      │
   │ • IP check           │
   │ • Rate limit         │
   │ • Signature match    │
   │ • Anomaly check      │
   │                      │
   │ Decision: Block/Allow│
   └─────────────────────┘
           │
           ↓ (if allowed)
   
2. TRAFFIC ANALYSIS
   ┌─────────────────────┐      ┌─────────────────────┐
   │ Threat Detection     │      │ Data Protection     │
   │ Service              │      │ Service             │
   │                      │      │                     │
   │ • Threat analysis    │      │ • PII detection     │
   │ • Anomaly detection   │      │ • Data classification│
   │ • Behavioral analysis│      │ • Compliance check  │
   │                      │      │                     │
   │ Output: Threat alerts│      │ Output: PII reports│
   └─────────────────────┘      └─────────────────────┘
           │                              │
           └──────────────┬───────────────┘
                          ↓
           3. SECURITY RESPONSE
           ┌─────────────────────┐
           │ Attack Blocking      │
           │ Service              │
           │                      │
           │ • Update blacklist   │
           │ • Create active block│
           │ • Block future reqs  │
           └─────────────────────┘
```

---

## Decision Matrix

Use this matrix to determine which service to use for your specific security need:

| Need | Threat Detection | Data Protection | Attack Blocking |
|------|-----------------|-----------------|----------------|
| Detect SQL injection attacks | ✅ | ❌ | ❌ |
| Block malicious requests | ❌ | ❌ | ✅ |
| Detect PII in data | ❌ | ✅ | ❌ |
| Rate limiting | ❌ | ❌ | ✅ |
| Detect XSS attacks | ✅ | ❌ | ❌ |
| Ensure GDPR compliance | ❌ | ✅ | ❌ |
| IP whitelisting/blacklisting | ❌ | ❌ | ✅ |
| Anomaly detection | ✅ | ❌ | ❌ |
| Classify data sensitivity | ❌ | ✅ | ❌ |
| Real-time request blocking | ❌ | ❌ | ✅ |
| Behavioral analysis | ✅ | ❌ | ❌ |
| Detect credit card numbers | ❌ | ✅ | ❌ |
| Geo-blocking | ❌ | ❌ | ✅ |
| Threat intelligence | ✅ | ❌ | ❌ |
| Compliance reporting | ❌ | ✅ | ❌ |
| Generate security alerts | ✅ | ❌ | ❌ |

---

## Summary

### Threat Detection Service
- **Focus**: Security threat analysis - "Is this a threat?"
- **When**: Continuous analysis of traffic
- **Output**: Threat alerts and analysis reports
- **Use Case**: Identifying and analyzing security threats

### Data Protection Service
- **Focus**: Data privacy and compliance - "Is this data sensitive?"
- **When**: On data processing and analysis
- **Output**: PII findings and compliance reports
- **Use Case**: Protecting sensitive data and ensuring compliance

### Attack Blocking Service
- **Focus**: Real-time prevention - "Should I block this request?"
- **When**: Real-time per request evaluation
- **Output**: Block/Allow decisions
- **Use Case**: Preventing attacks before they reach systems

---

## Related Documentation

- [Backend README](../backend/README.md) - Overall backend architecture
- [Threat Detection README](../backend/services/threat-detection/README.md) - Detailed Threat Detection docs
- [Data Protection README](../backend/services/data-protection/README.md) - Detailed Data Protection docs
- [Attack Blocking Service](../backend/services/attack-blocking/) - Attack Blocking implementation
- [Microservices Comparison](./MICROSERVICES_COMPARISON.md) - Other microservices comparison

---

**Last Updated**: 2024-01-15  
**Version**: 1.0  
**Maintained by**: ScopeAPI Development Team

