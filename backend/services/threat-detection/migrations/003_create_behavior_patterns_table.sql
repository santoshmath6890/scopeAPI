-- Migration: Create behavior_patterns table
-- Description: Creates the behavior_patterns table for storing behavioral analysis patterns
-- Version: 003
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS behavior_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id VARCHAR(255) NOT NULL,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('user', 'ip', 'api', 'endpoint')),
    pattern_type VARCHAR(100) NOT NULL,
    pattern_category VARCHAR(50) NOT NULL,
    
    -- Pattern data
    pattern_data JSONB NOT NULL,
    risk_score DECIMAL(5,2) NOT NULL CHECK (risk_score >= 0 AND risk_score <= 100),
    confidence_score DECIMAL(5,2) NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    
    -- Analysis context
    analysis_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    analysis_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    sample_size INTEGER NOT NULL DEFAULT 0,
    
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
    
    -- Constraints
    CONSTRAINT behavior_patterns_entity_type_check CHECK (entity_type IN ('user', 'ip', 'api', 'endpoint')),
    CONSTRAINT behavior_patterns_risk_check CHECK (risk_score >= 0 AND risk_score <= 100),
    CONSTRAINT behavior_patterns_confidence_check CHECK (confidence_score >= 0 AND confidence_score <= 100),
    CONSTRAINT behavior_patterns_sample_size_check CHECK (sample_size >= 0)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_entity_id ON behavior_patterns(entity_id);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_entity_type ON behavior_patterns(entity_type);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_pattern_type ON behavior_patterns(pattern_type);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_pattern_category ON behavior_patterns(pattern_category);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_risk_score ON behavior_patterns(risk_score);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_detection_timestamp ON behavior_patterns(detection_timestamp);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_created_at ON behavior_patterns(created_at);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_detection_method ON behavior_patterns(detection_method);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_detection_engine ON behavior_patterns(detection_engine);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_entity ON behavior_patterns(entity_id, entity_type);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_type_category ON behavior_patterns(pattern_type, pattern_category);
CREATE INDEX IF NOT EXISTS idx_behavior_patterns_risk_timestamp ON behavior_patterns(risk_score, detection_timestamp);

-- Add comments for documentation
COMMENT ON TABLE behavior_patterns IS 'Stores behavioral analysis patterns and risk assessments';
COMMENT ON COLUMN behavior_patterns.id IS 'Unique identifier for the behavior pattern';
COMMENT ON COLUMN behavior_patterns.entity_id IS 'ID of the entity being analyzed (user, IP, API, endpoint)';
COMMENT ON COLUMN behavior_patterns.entity_type IS 'Type of entity being analyzed';
COMMENT ON COLUMN behavior_patterns.pattern_type IS 'Type of behavior pattern detected';
COMMENT ON COLUMN behavior_patterns.pattern_category IS 'Category of the behavior pattern';
COMMENT ON COLUMN behavior_patterns.pattern_data IS 'Detailed pattern data and metrics (JSON)';
COMMENT ON COLUMN behavior_patterns.risk_score IS 'Risk score associated with this pattern (0-100)';
COMMENT ON COLUMN behavior_patterns.confidence_score IS 'Confidence level of the pattern detection (0-100)';
COMMENT ON COLUMN behavior_patterns.analysis_period_start IS 'Start of the analysis period';
COMMENT ON COLUMN behavior_patterns.analysis_period_end IS 'End of the analysis period';
COMMENT ON COLUMN behavior_patterns.sample_size IS 'Number of samples used in the analysis';
COMMENT ON COLUMN behavior_patterns.detection_method IS 'Method used to detect the pattern';
COMMENT ON COLUMN behavior_patterns.detection_engine IS 'Engine that performed the detection';
COMMENT ON COLUMN behavior_patterns.detection_timestamp IS 'When the pattern was detected';
COMMENT ON COLUMN behavior_patterns.analysis_result IS 'Detailed analysis results (JSON)';
COMMENT ON COLUMN behavior_patterns.recommendations IS 'Recommended actions (JSON)';
COMMENT ON COLUMN behavior_patterns.created_at IS 'When the record was created';
COMMENT ON COLUMN behavior_patterns.updated_at IS 'When the record was last updated';
