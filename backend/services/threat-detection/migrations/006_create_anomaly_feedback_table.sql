-- Migration: Create anomaly_feedback table
-- Description: Creates the anomaly_feedback table for storing user feedback on anomalies
-- Version: 006
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS anomaly_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anomaly_id UUID NOT NULL,
    feedback_type VARCHAR(50) NOT NULL CHECK (feedback_type IN ('false_positive', 'true_positive', 'investigation_result')),
    feedback_value BOOLEAN,
    feedback_text TEXT,
    
    -- Feedback source
    provided_by VARCHAR(255) NOT NULL,
    provided_by_type VARCHAR(50) NOT NULL CHECK (provided_by_type IN ('user', 'system', 'admin', 'analyst')),
    
    -- Feedback context
    investigation_notes TEXT,
    action_taken TEXT,
    resolution_status VARCHAR(50) CHECK (resolution_status IN ('resolved', 'investigating', 'escalated', 'closed')),
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT anomaly_feedback_type_check CHECK (feedback_type IN ('false_positive', 'true_positive', 'investigation_result')),
    CONSTRAINT anomaly_feedback_provider_type_check CHECK (provided_by_type IN ('user', 'system', 'admin', 'analyst')),
    CONSTRAINT anomaly_feedback_resolution_check CHECK (resolution_status IN ('resolved', 'investigating', 'escalated', 'closed')),
    
    -- Foreign key constraint (if anomalies table exists)
    CONSTRAINT fk_anomaly_feedback_anomaly_id FOREIGN KEY (anomaly_id) REFERENCES anomalies(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_anomaly_id ON anomaly_feedback(anomaly_id);
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_feedback_type ON anomaly_feedback(feedback_type);
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_provided_by ON anomaly_feedback(provided_by);
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_provided_by_type ON anomaly_feedback(provided_by_type);
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_resolution_status ON anomaly_feedback(resolution_status);
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_created_at ON anomaly_feedback(created_at);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_anomaly_type ON anomaly_feedback(anomaly_id, feedback_type);
CREATE INDEX IF NOT EXISTS idx_anomaly_feedback_provider_type ON anomaly_feedback(provided_by, provided_by_type);

-- Add comments for documentation
COMMENT ON TABLE anomaly_feedback IS 'Stores user feedback and investigation results for anomalies';
COMMENT ON COLUMN anomaly_feedback.id IS 'Unique identifier for the feedback';
COMMENT ON COLUMN anomaly_feedback.anomaly_id IS 'ID of the anomaly this feedback is for';
COMMENT ON COLUMN anomaly_feedback.feedback_type IS 'Type of feedback provided';
COMMENT ON COLUMN anomaly_feedback.feedback_value IS 'Boolean value indicating true/false positive';
COMMENT ON COLUMN anomaly_feedback.feedback_text IS 'Text description of the feedback';
COMMENT ON COLUMN anomaly_feedback.provided_by IS 'Who provided the feedback';
COMMENT ON COLUMN anomaly_feedback.provided_by_type IS 'Type of person who provided the feedback';
COMMENT ON COLUMN anomaly_feedback.investigation_notes IS 'Notes from the investigation';
COMMENT ON COLUMN anomaly_feedback.action_taken IS 'Actions taken based on the feedback';
COMMENT ON COLUMN anomaly_feedback.resolution_status IS 'Current resolution status';
COMMENT ON COLUMN anomaly_feedback.created_at IS 'When the feedback was provided';
COMMENT ON COLUMN anomaly_feedback.updated_at IS 'When the feedback was last updated';
