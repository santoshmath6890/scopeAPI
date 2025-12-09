-- Migration: Create threats table
-- Description: Creates the threats table for storing detected threats
-- Version: 001
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS threats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id VARCHAR(255) NOT NULL,
    threat_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'investigating', 'resolved', 'false_positive')),
    confidence_score DECIMAL(5,2) NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    risk_score DECIMAL(5,2) NOT NULL CHECK (risk_score >= 0 AND risk_score <= 100),
    
    -- Source information
    source_ip INET,
    user_agent TEXT,
    api_id VARCHAR(255),
    endpoint_id VARCHAR(255),
    user_id VARCHAR(255),
    
    -- Attack details
    attack_vector TEXT,
    attack_pattern TEXT,
    payload TEXT,
    headers JSONB,
    parameters JSONB,
    
    -- Detection information
    detection_method VARCHAR(50) NOT NULL,
    detection_engine VARCHAR(50) NOT NULL,
    detection_timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Analysis results
    analysis_result JSONB,
    recommendations JSONB,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Indexes for performance
    CONSTRAINT threats_severity_check CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT threats_status_check CHECK (status IN ('active', 'investigating', 'resolved', 'false_positive')),
    CONSTRAINT threats_confidence_check CHECK (confidence_score >= 0 AND confidence_score <= 100),
    CONSTRAINT threats_risk_check CHECK (risk_score >= 0 AND risk_score <= 100)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_threats_threat_type ON threats(threat_type);
CREATE INDEX IF NOT EXISTS idx_threats_severity ON threats(severity);
CREATE INDEX IF NOT EXISTS idx_threats_status ON threats(status);
CREATE INDEX IF NOT EXISTS idx_threats_source_ip ON threats(source_ip);
CREATE INDEX IF NOT EXISTS idx_threats_api_id ON threats(api_id);
CREATE INDEX IF NOT EXISTS idx_threats_endpoint_id ON threats(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_threats_user_id ON threats(user_id);
CREATE INDEX IF NOT EXISTS idx_threats_detection_timestamp ON threats(detection_timestamp);
CREATE INDEX IF NOT EXISTS idx_threats_created_at ON threats(created_at);
CREATE INDEX IF NOT EXISTS idx_threats_detection_method ON threats(detection_method);
CREATE INDEX IF NOT EXISTS idx_threats_detection_engine ON threats(detection_engine);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_threats_type_severity ON threats(threat_type, severity);
CREATE INDEX IF NOT EXISTS idx_threats_status_timestamp ON threats(status, detection_timestamp);
CREATE INDEX IF NOT EXISTS idx_threats_api_endpoint ON threats(api_id, endpoint_id);

-- Add comments for documentation
COMMENT ON TABLE threats IS 'Stores detected security threats and attacks';
COMMENT ON COLUMN threats.id IS 'Unique identifier for the threat';
COMMENT ON COLUMN threats.request_id IS 'Identifier for the original request that triggered the threat';
COMMENT ON COLUMN threats.threat_type IS 'Type of threat detected (e.g., sql_injection, xss, ddos)';
COMMENT ON COLUMN threats.severity IS 'Severity level of the threat';
COMMENT ON COLUMN threats.status IS 'Current status of the threat investigation';
COMMENT ON COLUMN threats.confidence_score IS 'Confidence level of the detection (0-100)';
COMMENT ON COLUMN threats.risk_score IS 'Risk score of the threat (0-100)';
COMMENT ON COLUMN threats.source_ip IS 'IP address of the threat source';
COMMENT ON COLUMN threats.user_agent IS 'User agent string from the request';
COMMENT ON COLUMN threats.api_id IS 'ID of the API that was targeted';
COMMENT ON COLUMN threats.endpoint_id IS 'ID of the specific endpoint that was targeted';
COMMENT ON COLUMN threats.user_id IS 'ID of the user associated with the request';
COMMENT ON COLUMN threats.attack_vector IS 'Description of the attack vector used';
COMMENT ON COLUMN threats.attack_pattern IS 'Pattern or signature that matched';
COMMENT ON COLUMN threats.payload IS 'The actual payload or data that triggered the threat';
COMMENT ON COLUMN threats.headers IS 'HTTP headers from the request (JSON)';
COMMENT ON COLUMN threats.parameters IS 'Request parameters (JSON)';
COMMENT ON COLUMN threats.detection_method IS 'Method used to detect the threat (signature, anomaly, behavioral)';
COMMENT ON COLUMN threats.detection_engine IS 'Engine that performed the detection';
COMMENT ON COLUMN threats.detection_timestamp IS 'When the threat was detected';
COMMENT ON COLUMN threats.analysis_result IS 'Detailed analysis results (JSON)';
COMMENT ON COLUMN threats.recommendations IS 'Recommended actions (JSON)';
COMMENT ON COLUMN threats.created_at IS 'When the record was created';
COMMENT ON COLUMN threats.updated_at IS 'When the record was last updated';
