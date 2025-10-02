-- Migration: Create threat_statistics table
-- Description: Creates the threat_statistics table for storing aggregated threat statistics
-- Version: 007
-- Date: 2024-01-15

CREATE TABLE IF NOT EXISTS threat_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    time_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    time_granularity VARCHAR(20) NOT NULL CHECK (time_granularity IN ('minute', 'hour', 'day', 'week', 'month')),
    
    -- Threat counts by type
    total_threats INTEGER NOT NULL DEFAULT 0,
    sql_injection_count INTEGER NOT NULL DEFAULT 0,
    xss_count INTEGER NOT NULL DEFAULT 0,
    ddos_count INTEGER NOT NULL DEFAULT 0,
    brute_force_count INTEGER NOT NULL DEFAULT 0,
    data_exfiltration_count INTEGER NOT NULL DEFAULT 0,
    path_traversal_count INTEGER NOT NULL DEFAULT 0,
    command_injection_count INTEGER NOT NULL DEFAULT 0,
    other_threat_count INTEGER NOT NULL DEFAULT 0,
    
    -- Threat counts by severity
    low_severity_count INTEGER NOT NULL DEFAULT 0,
    medium_severity_count INTEGER NOT NULL DEFAULT 0,
    high_severity_count INTEGER NOT NULL DEFAULT 0,
    critical_severity_count INTEGER NOT NULL DEFAULT 0,
    
    -- Threat counts by status
    active_threat_count INTEGER NOT NULL DEFAULT 0,
    investigating_threat_count INTEGER NOT NULL DEFAULT 0,
    resolved_threat_count INTEGER NOT NULL DEFAULT 0,
    false_positive_threat_count INTEGER NOT NULL DEFAULT 0,
    
    -- Detection method counts
    signature_detection_count INTEGER NOT NULL DEFAULT 0,
    anomaly_detection_count INTEGER NOT NULL DEFAULT 0,
    behavioral_detection_count INTEGER NOT NULL DEFAULT 0,
    ml_detection_count INTEGER NOT NULL DEFAULT 0,
    
    -- Additional metrics
    unique_source_ips INTEGER NOT NULL DEFAULT 0,
    unique_user_agents INTEGER NOT NULL DEFAULT 0,
    unique_apis_affected INTEGER NOT NULL DEFAULT 0,
    unique_endpoints_affected INTEGER NOT NULL DEFAULT 0,
    unique_users_affected INTEGER NOT NULL DEFAULT 0,
    
    -- Average scores
    avg_confidence_score DECIMAL(5,2),
    avg_risk_score DECIMAL(5,2),
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT threat_statistics_granularity_check CHECK (time_granularity IN ('minute', 'hour', 'day', 'week', 'month')),
    CONSTRAINT threat_statistics_counts_check CHECK (
        total_threats >= 0 AND
        sql_injection_count >= 0 AND
        xss_count >= 0 AND
        ddos_count >= 0 AND
        brute_force_count >= 0 AND
        data_exfiltration_count >= 0 AND
        path_traversal_count >= 0 AND
        command_injection_count >= 0 AND
        other_threat_count >= 0 AND
        low_severity_count >= 0 AND
        medium_severity_count >= 0 AND
        high_severity_count >= 0 AND
        critical_severity_count >= 0 AND
        active_threat_count >= 0 AND
        investigating_threat_count >= 0 AND
        resolved_threat_count >= 0 AND
        false_positive_threat_count >= 0 AND
        signature_detection_count >= 0 AND
        anomaly_detection_count >= 0 AND
        behavioral_detection_count >= 0 AND
        ml_detection_count >= 0 AND
        unique_source_ips >= 0 AND
        unique_user_agents >= 0 AND
        unique_apis_affected >= 0 AND
        unique_endpoints_affected >= 0 AND
        unique_users_affected >= 0
    ),
    
    -- Unique constraint for time period and granularity
    UNIQUE(time_period_start, time_period_end, time_granularity)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_threat_statistics_time_period ON threat_statistics(time_period_start, time_period_end);
CREATE INDEX IF NOT EXISTS idx_threat_statistics_granularity ON threat_statistics(time_granularity);
CREATE INDEX IF NOT EXISTS idx_threat_statistics_total_threats ON threat_statistics(total_threats);
CREATE INDEX IF NOT EXISTS idx_threat_statistics_created_at ON threat_statistics(created_at);

-- Create composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_threat_statistics_time_granularity ON threat_statistics(time_period_start, time_granularity);
CREATE INDEX IF NOT EXISTS idx_threat_statistics_period_granularity ON threat_statistics(time_period_start, time_period_end, time_granularity);

-- Add comments for documentation
COMMENT ON TABLE threat_statistics IS 'Stores aggregated threat detection statistics for reporting and analysis';
COMMENT ON COLUMN threat_statistics.id IS 'Unique identifier for the statistics record';
COMMENT ON COLUMN threat_statistics.time_period_start IS 'Start of the time period';
COMMENT ON COLUMN threat_statistics.time_period_end IS 'End of the time period';
COMMENT ON COLUMN threat_statistics.time_granularity IS 'Granularity of the time period';
COMMENT ON COLUMN threat_statistics.total_threats IS 'Total number of threats detected';
COMMENT ON COLUMN threat_statistics.sql_injection_count IS 'Number of SQL injection threats';
COMMENT ON COLUMN threat_statistics.xss_count IS 'Number of XSS threats';
COMMENT ON COLUMN threat_statistics.ddos_count IS 'Number of DDoS threats';
COMMENT ON COLUMN threat_statistics.brute_force_count IS 'Number of brute force threats';
COMMENT ON COLUMN threat_statistics.data_exfiltration_count IS 'Number of data exfiltration threats';
COMMENT ON COLUMN threat_statistics.path_traversal_count IS 'Number of path traversal threats';
COMMENT ON COLUMN threat_statistics.command_injection_count IS 'Number of command injection threats';
COMMENT ON COLUMN threat_statistics.other_threat_count IS 'Number of other types of threats';
COMMENT ON COLUMN threat_statistics.low_severity_count IS 'Number of low severity threats';
COMMENT ON COLUMN threat_statistics.medium_severity_count IS 'Number of medium severity threats';
COMMENT ON COLUMN threat_statistics.high_severity_count IS 'Number of high severity threats';
COMMENT ON COLUMN threat_statistics.critical_severity_count IS 'Number of critical severity threats';
COMMENT ON COLUMN threat_statistics.active_threat_count IS 'Number of active threats';
COMMENT ON COLUMN threat_statistics.investigating_threat_count IS 'Number of threats under investigation';
COMMENT ON COLUMN threat_statistics.resolved_threat_count IS 'Number of resolved threats';
COMMENT ON COLUMN threat_statistics.false_positive_threat_count IS 'Number of false positive threats';
COMMENT ON COLUMN threat_statistics.signature_detection_count IS 'Number of threats detected by signature matching';
COMMENT ON COLUMN threat_statistics.anomaly_detection_count IS 'Number of threats detected by anomaly detection';
COMMENT ON COLUMN threat_statistics.behavioral_detection_count IS 'Number of threats detected by behavioral analysis';
COMMENT ON COLUMN threat_statistics.ml_detection_count IS 'Number of threats detected by ML models';
COMMENT ON COLUMN threat_statistics.unique_source_ips IS 'Number of unique source IPs';
COMMENT ON COLUMN threat_statistics.unique_user_agents IS 'Number of unique user agents';
COMMENT ON COLUMN threat_statistics.unique_apis_affected IS 'Number of unique APIs affected';
COMMENT ON COLUMN threat_statistics.unique_endpoints_affected IS 'Number of unique endpoints affected';
COMMENT ON COLUMN threat_statistics.unique_users_affected IS 'Number of unique users affected';
COMMENT ON COLUMN threat_statistics.avg_confidence_score IS 'Average confidence score of threats';
COMMENT ON COLUMN threat_statistics.avg_risk_score IS 'Average risk score of threats';
COMMENT ON COLUMN threat_statistics.created_at IS 'When the record was created';
COMMENT ON COLUMN threat_statistics.updated_at IS 'When the record was last updated';
