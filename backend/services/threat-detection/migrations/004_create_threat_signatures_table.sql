-- Migration: Create threat_signatures table
-- Description: Creates the threat_signatures table for storing threat detection signatures
-- Version: 004
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS threat_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    signature_type VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    
    -- Signature data
    pattern TEXT NOT NULL,
    regex_pattern TEXT,
    signature_data JSONB,
    
    -- Signature metadata
    version VARCHAR(20) NOT NULL DEFAULT '1.0',
    signature_set VARCHAR(100),
    author VARCHAR(255),
    created_by VARCHAR(255),
    
    -- Status and configuration
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    priority INTEGER NOT NULL DEFAULT 0,
    
    -- Performance metrics
    match_count INTEGER NOT NULL DEFAULT 0,
    false_positive_count INTEGER NOT NULL DEFAULT 0,
    last_matched_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT threat_signatures_severity_check CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT threat_signatures_match_count_check CHECK (match_count >= 0),
    CONSTRAINT threat_signatures_false_positive_check CHECK (false_positive_count >= 0),
    CONSTRAINT threat_signatures_priority_check CHECK (priority >= 0)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_threat_signatures_name ON threat_signatures(name);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_signature_type ON threat_signatures(signature_type);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_category ON threat_signatures(category);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_severity ON threat_signatures(severity);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_signature_set ON threat_signatures(signature_set);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_is_active ON threat_signatures(is_active);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_is_custom ON threat_signatures(is_custom);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_priority ON threat_signatures(priority);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_created_at ON threat_signatures(created_at);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_last_matched ON threat_signatures(last_matched_at);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_threat_signatures_type_category ON threat_signatures(signature_type, category);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_active_priority ON threat_signatures(is_active, priority);
CREATE INDEX IF NOT EXISTS idx_threat_signatures_set_active ON threat_signatures(signature_set, is_active);

-- Add comments for documentation
COMMENT ON TABLE threat_signatures IS 'Stores threat detection signatures and patterns';
COMMENT ON COLUMN threat_signatures.id IS 'Unique identifier for the signature';
COMMENT ON COLUMN threat_signatures.name IS 'Name of the signature';
COMMENT ON COLUMN threat_signatures.description IS 'Description of what the signature detects';
COMMENT ON COLUMN threat_signatures.signature_type IS 'Type of signature (sql_injection, xss, ddos, etc.)';
COMMENT ON COLUMN threat_signatures.category IS 'Category of the signature';
COMMENT ON COLUMN threat_signatures.severity IS 'Severity level of threats detected by this signature';
COMMENT ON COLUMN threat_signatures.pattern IS 'The signature pattern or rule';
COMMENT ON COLUMN threat_signatures.regex_pattern IS 'Compiled regex pattern for matching';
COMMENT ON COLUMN threat_signatures.signature_data IS 'Additional signature data and configuration (JSON)';
COMMENT ON COLUMN threat_signatures.version IS 'Version of the signature';
COMMENT ON COLUMN threat_signatures.signature_set IS 'Signature set this belongs to';
COMMENT ON COLUMN threat_signatures.author IS 'Author of the signature';
COMMENT ON COLUMN threat_signatures.created_by IS 'User who created/imported the signature';
COMMENT ON COLUMN threat_signatures.is_active IS 'Whether the signature is currently active';
COMMENT ON COLUMN threat_signatures.is_custom IS 'Whether this is a custom signature';
COMMENT ON COLUMN threat_signatures.priority IS 'Priority of the signature (higher = more important)';
COMMENT ON COLUMN threat_signatures.match_count IS 'Number of times this signature has matched';
COMMENT ON COLUMN threat_signatures.false_positive_count IS 'Number of false positives for this signature';
COMMENT ON COLUMN threat_signatures.last_matched_at IS 'When this signature was last matched';
COMMENT ON COLUMN threat_signatures.created_at IS 'When the record was created';
COMMENT ON COLUMN threat_signatures.updated_at IS 'When the record was last updated';
