-- ScopeAPI Database Initialization Script
-- This script creates the basic database structure for ScopeAPI

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS api_discovery;
CREATE SCHEMA IF NOT EXISTS threat_detection;
CREATE SCHEMA IF NOT EXISTS data_protection;
CREATE SCHEMA IF NOT EXISTS attack_blocking;
CREATE SCHEMA IF NOT EXISTS gateway_integration;
CREATE SCHEMA IF NOT EXISTS audit;

-- Create basic tables for API Discovery
CREATE TABLE IF NOT EXISTS api_discovery.endpoints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url VARCHAR(500) NOT NULL,
    method VARCHAR(10) NOT NULL,
    service_name VARCHAR(100),
    discovered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    metadata JSONB
);

-- Create basic tables for Threat Detection
CREATE TABLE IF NOT EXISTS threat_detection.threats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    threat_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'new',
    source_ip INET,
    target_endpoint VARCHAR(500),
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB
);

-- Create basic tables for Data Protection
CREATE TABLE IF NOT EXISTS data_protection.classifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    data_type VARCHAR(100) NOT NULL,
    sensitivity_level VARCHAR(20) NOT NULL,
    classification_rules JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create basic tables for Attack Blocking
CREATE TABLE IF NOT EXISTS attack_blocking.blocking_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name VARCHAR(200) NOT NULL,
    rule_type VARCHAR(50) NOT NULL,
    conditions JSONB NOT NULL,
    action VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create basic tables for Gateway Integration
CREATE TABLE IF NOT EXISTS gateway_integration.integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    gateway_type VARCHAR(50) NOT NULL,
    gateway_name VARCHAR(200) NOT NULL,
    configuration JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create audit log table
CREATE TABLE IF NOT EXISTS audit.logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    user_id VARCHAR(100),
    resource_type VARCHAR(100),
    resource_id UUID,
    details JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_endpoints_url ON api_discovery.endpoints(url);
CREATE INDEX IF NOT EXISTS idx_endpoints_service ON api_discovery.endpoints(service_name);
CREATE INDEX IF NOT EXISTS idx_threats_status ON threat_detection.threats(status);
CREATE INDEX IF NOT EXISTS idx_threats_severity ON threat_detection.threats(severity);
CREATE INDEX IF NOT EXISTS idx_threats_detected_at ON threat_detection.threats(detected_at);
CREATE INDEX IF NOT EXISTS idx_blocking_rules_type ON attack_blocking.blocking_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_integrations_type ON gateway_integration.integrations(gateway_type);
CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit.logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_service ON audit.logs(service_name);

-- Insert some sample data
INSERT INTO api_discovery.endpoints (url, method, service_name) VALUES
    ('/api/v1/health', 'GET', 'health-check'),
    ('/api/v1/discovery', 'POST', 'api-discovery'),
    ('/api/v1/threats', 'GET', 'threat-detection')
ON CONFLICT DO NOTHING;

INSERT INTO threat_detection.threats (threat_type, severity, source_ip, target_endpoint) VALUES
    ('SQL_INJECTION', 'HIGH', '192.168.1.100', '/api/v1/users'),
    ('XSS_ATTACK', 'MEDIUM', '10.0.0.50', '/api/v1/comments')
ON CONFLICT DO NOTHING;

INSERT INTO attack_blocking.blocking_rules (rule_name, rule_type, conditions, action) VALUES
    ('Block SQL Injection', 'PATTERN_MATCH', '{"pattern": ".*(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE).*", "field": "query"}', 'BLOCK'),
    ('Block XSS Attacks', 'PATTERN_MATCH', '{"pattern": ".*<script.*>.*</script>.*", "field": "content"}', 'BLOCK')
ON CONFLICT DO NOTHING;

-- Grant permissions to scopeapi user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA api_discovery TO scopeapi;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA threat_detection TO scopeapi;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA data_protection TO scopeapi;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA attack_blocking TO scopeapi;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA gateway_integration TO scopeapi;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA audit TO scopeapi;

GRANT USAGE ON SCHEMA api_discovery TO scopeapi;
GRANT USAGE ON SCHEMA threat_detection TO scopeapi;
GRANT USAGE ON SCHEMA data_protection TO scopeapi;
GRANT USAGE ON SCHEMA attack_blocking TO scopeapi;
GRANT USAGE ON SCHEMA gateway_integration TO scopeapi;
GRANT USAGE ON SCHEMA audit TO scopeapi;

-- Create a function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updated_at
CREATE TRIGGER update_gateway_integration_updated_at 
    BEFORE UPDATE ON gateway_integration.integrations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'ScopeAPI database initialization completed successfully!';
    RAISE NOTICE 'Created schemas: api_discovery, threat_detection, data_protection, attack_blocking, gateway_integration, audit';
    RAISE NOTICE 'Sample data inserted for testing';
END $$; 