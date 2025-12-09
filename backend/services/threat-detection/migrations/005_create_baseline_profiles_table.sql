-- Migration: Create baseline_profiles table
-- Description: Creates the baseline_profiles table for storing behavioral baselines
-- Version: 005
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS baseline_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id VARCHAR(255) NOT NULL,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('user', 'ip', 'api', 'endpoint')),
    baseline_type VARCHAR(100) NOT NULL,
    
    -- Baseline data
    baseline_data JSONB NOT NULL,
    metrics JSONB NOT NULL,
    thresholds JSONB,
    
    -- Training information
    training_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    training_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    training_sample_count INTEGER NOT NULL DEFAULT 0,
    training_data_hash VARCHAR(64),
    
    -- Profile metadata
    version VARCHAR(20) NOT NULL DEFAULT '1.0',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    confidence_score DECIMAL(5,2) NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    
    -- Update tracking
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    update_count INTEGER NOT NULL DEFAULT 0,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT baseline_profiles_entity_type_check CHECK (entity_type IN ('user', 'ip', 'api', 'endpoint')),
    CONSTRAINT baseline_profiles_sample_count_check CHECK (training_sample_count >= 0),
    CONSTRAINT baseline_profiles_confidence_check CHECK (confidence_score >= 0 AND confidence_score <= 100),
    CONSTRAINT baseline_profiles_update_count_check CHECK (update_count >= 0),
    
    -- Unique constraint for entity-type combination
    UNIQUE(entity_id, entity_type, baseline_type)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_entity_id ON baseline_profiles(entity_id);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_entity_type ON baseline_profiles(entity_type);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_baseline_type ON baseline_profiles(baseline_type);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_is_active ON baseline_profiles(is_active);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_confidence_score ON baseline_profiles(confidence_score);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_last_updated ON baseline_profiles(last_updated_at);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_created_at ON baseline_profiles(created_at);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_entity ON baseline_profiles(entity_id, entity_type);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_type_active ON baseline_profiles(baseline_type, is_active);
CREATE INDEX IF NOT EXISTS idx_baseline_profiles_entity_active ON baseline_profiles(entity_id, entity_type, is_active);

-- Add comments for documentation
COMMENT ON TABLE baseline_profiles IS 'Stores behavioral baseline profiles for anomaly detection';
COMMENT ON COLUMN baseline_profiles.id IS 'Unique identifier for the baseline profile';
COMMENT ON COLUMN baseline_profiles.entity_id IS 'ID of the entity this baseline is for';
COMMENT ON COLUMN baseline_profiles.entity_type IS 'Type of entity (user, ip, api, endpoint)';
COMMENT ON COLUMN baseline_profiles.baseline_type IS 'Type of baseline (access_pattern, usage_pattern, timing_pattern, etc.)';
COMMENT ON COLUMN baseline_profiles.baseline_data IS 'The actual baseline data and patterns (JSON)';
COMMENT ON COLUMN baseline_profiles.metrics IS 'Statistical metrics for the baseline (JSON)';
COMMENT ON COLUMN baseline_profiles.thresholds IS 'Anomaly detection thresholds (JSON)';
COMMENT ON COLUMN baseline_profiles.training_period_start IS 'Start of the training period';
COMMENT ON COLUMN baseline_profiles.training_period_end IS 'End of the training period';
COMMENT ON COLUMN baseline_profiles.training_sample_count IS 'Number of samples used for training';
COMMENT ON COLUMN baseline_profiles.training_data_hash IS 'Hash of the training data for integrity';
COMMENT ON COLUMN baseline_profiles.version IS 'Version of the baseline profile';
COMMENT ON COLUMN baseline_profiles.is_active IS 'Whether this baseline is currently active';
COMMENT ON COLUMN baseline_profiles.confidence_score IS 'Confidence level of the baseline (0-100)';
COMMENT ON COLUMN baseline_profiles.last_updated_at IS 'When the baseline was last updated';
COMMENT ON COLUMN baseline_profiles.update_count IS 'Number of times this baseline has been updated';
COMMENT ON COLUMN baseline_profiles.created_at IS 'When the record was created';
COMMENT ON COLUMN baseline_profiles.updated_at IS 'When the record was last updated';
