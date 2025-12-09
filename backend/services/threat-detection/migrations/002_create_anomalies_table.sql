-- Migration: Create anomalies table
-- Description: Creates the anomalies table for storing detected anomalies
-- Version: 002
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id VARCHAR(255) NOT NULL,
    anomaly_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'investigating', 'resolved', 'false_positive')),
    confidence_score DECIMAL(5,2) NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    anomaly_score DECIMAL(5,2) NOT NULL CHECK (anomaly_score >= 0 AND anomaly_score <= 100),
    
    -- Source information
    source_ip INET,
    user_agent TEXT,
    api_id VARCHAR(255),
    endpoint_id VARCHAR(255),
    user_id VARCHAR(255),
    
    -- Anomaly details
    anomaly_description TEXT,
    detected_pattern TEXT,
    baseline_value DECIMAL(15,4),
    actual_value DECIMAL(15,4),
    deviation_percentage DECIMAL(5,2),
    
    -- Detection information
    detection_method VARCHAR(50) NOT NULL,
    detection_engine VARCHAR(50) NOT NULL,
    detection_timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Analysis results
    analysis_result JSONB,
    recommendations JSONB,
    false_positive BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT anomalies_severity_check CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT anomalies_status_check CHECK (status IN ('active', 'investigating', 'resolved', 'false_positive')),
    CONSTRAINT anomalies_confidence_check CHECK (confidence_score >= 0 AND confidence_score <= 100),
    CONSTRAINT anomalies_score_check CHECK (anomaly_score >= 0 AND anomaly_score <= 100)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_anomalies_anomaly_type ON anomalies(anomaly_type);
CREATE INDEX IF NOT EXISTS idx_anomalies_severity ON anomalies(severity);
CREATE INDEX IF NOT EXISTS idx_anomalies_status ON anomalies(status);
CREATE INDEX IF NOT EXISTS idx_anomalies_source_ip ON anomalies(source_ip);
CREATE INDEX IF NOT EXISTS idx_anomalies_api_id ON anomalies(api_id);
CREATE INDEX IF NOT EXISTS idx_anomalies_endpoint_id ON anomalies(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_anomalies_user_id ON anomalies(user_id);
CREATE INDEX IF NOT EXISTS idx_anomalies_detection_timestamp ON anomalies(detection_timestamp);
CREATE INDEX IF NOT EXISTS idx_anomalies_created_at ON anomalies(created_at);
CREATE INDEX IF NOT EXISTS idx_anomalies_detection_method ON anomalies(detection_method);
CREATE INDEX IF NOT EXISTS idx_anomalies_detection_engine ON anomalies(detection_engine);
CREATE INDEX IF NOT EXISTS idx_anomalies_false_positive ON anomalies(false_positive);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_anomalies_type_severity ON anomalies(anomaly_type, severity);
CREATE INDEX IF NOT EXISTS idx_anomalies_status_timestamp ON anomalies(status, detection_timestamp);
CREATE INDEX IF NOT EXISTS idx_anomalies_api_endpoint ON anomalies(api_id, endpoint_id);

-- Add comments for documentation
COMMENT ON TABLE anomalies IS 'Stores detected anomalies in API traffic patterns';
COMMENT ON COLUMN anomalies.id IS 'Unique identifier for the anomaly';
COMMENT ON COLUMN anomalies.request_id IS 'Identifier for the original request that triggered the anomaly';
COMMENT ON COLUMN anomalies.anomaly_type IS 'Type of anomaly detected (e.g., traffic_volume, response_time, access_pattern)';
COMMENT ON COLUMN anomalies.severity IS 'Severity level of the anomaly';
COMMENT ON COLUMN anomalies.status IS 'Current status of the anomaly investigation';
COMMENT ON COLUMN anomalies.confidence_score IS 'Confidence level of the detection (0-100)';
COMMENT ON COLUMN anomalies.anomaly_score IS 'Anomaly score indicating deviation from normal (0-100)';
COMMENT ON COLUMN anomalies.source_ip IS 'IP address of the anomaly source';
COMMENT ON COLUMN anomalies.user_agent IS 'User agent string from the request';
COMMENT ON COLUMN anomalies.api_id IS 'ID of the API that was affected';
COMMENT ON COLUMN anomalies.endpoint_id IS 'ID of the specific endpoint that was affected';
COMMENT ON COLUMN anomalies.user_id IS 'ID of the user associated with the request';
COMMENT ON COLUMN anomalies.anomaly_description IS 'Human-readable description of the anomaly';
COMMENT ON COLUMN anomalies.detected_pattern IS 'Pattern or metric that was anomalous';
COMMENT ON COLUMN anomalies.baseline_value IS 'Expected normal value';
COMMENT ON COLUMN anomalies.actual_value IS 'Actual observed value';
COMMENT ON COLUMN anomalies.deviation_percentage IS 'Percentage deviation from baseline';
COMMENT ON COLUMN anomalies.detection_method IS 'Method used to detect the anomaly (statistical, ml, etc.)';
COMMENT ON COLUMN anomalies.detection_engine IS 'Engine that performed the detection';
COMMENT ON COLUMN anomalies.detection_timestamp IS 'When the anomaly was detected';
COMMENT ON COLUMN anomalies.analysis_result IS 'Detailed analysis results (JSON)';
COMMENT ON COLUMN anomalies.recommendations IS 'Recommended actions (JSON)';
COMMENT ON COLUMN anomalies.false_positive IS 'Whether this anomaly was marked as a false positive';
COMMENT ON COLUMN anomalies.created_at IS 'When the record was created';
COMMENT ON COLUMN anomalies.updated_at IS 'When the record was last updated';
