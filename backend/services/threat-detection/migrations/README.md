# Threat Detection Service Migrations

This directory contains database migration scripts for the Threat Detection service.

## Migration Files

- `001_create_threats_table.sql` - Creates the threats table for storing detected threats
- `002_create_anomalies_table.sql` - Creates the anomalies table for storing detected anomalies
- `003_create_behavior_patterns_table.sql` - Creates the behavior_patterns table for behavioral analysis
- `004_create_threat_signatures_table.sql` - Creates the threat_signatures table for signature-based detection
- `005_create_baseline_profiles_table.sql` - Creates the baseline_profiles table for behavioral baselines
- `006_create_anomaly_feedback_table.sql` - Creates the anomaly_feedback table for user feedback
- `007_create_threat_statistics_table.sql` - Creates the threat_statistics table for aggregated statistics

## Running Migrations

### Prerequisites

1. PostgreSQL database running and accessible
2. Database user with CREATE TABLE permissions
3. Go 1.22+ installed

### Environment Variables

Set the following environment variables:

```bash
export DATABASE_URL="postgres://username:password@localhost/database_name?sslmode=disable"
export MIGRATIONS_DIR="./migrations"  # Optional, defaults to ./migrations
```

### Running Migrations

1. Navigate to the migrations directory:
   ```bash
   cd backend/services/threat-detection/migrations
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run all pending migrations:
   ```bash
   go run migrate.go
   ```

4. Check migration status:
   ```bash
   go run migrate.go status
   ```

### Migration Naming Convention

Migration files follow the pattern: `{version}_{description}.sql`

- Version: 3-digit zero-padded number (001, 002, 003, etc.)
- Description: Snake_case description of what the migration does

Examples:
- `001_create_threats_table.sql`
- `002_create_anomalies_table.sql`
- `003_add_indexes_to_threats.sql`

### Migration Features

- **Automatic tracking**: Migrations are tracked in the `schema_migrations` table
- **Idempotent**: Safe to run multiple times
- **Transactional**: Each migration runs in a transaction
- **Version control**: Migrations are applied in version order
- **Status checking**: Can check which migrations have been applied

### Database Schema Overview

The migrations create the following tables:

1. **threats** - Core threat detection results
2. **anomalies** - Anomaly detection results
3. **behavior_patterns** - Behavioral analysis patterns
4. **threat_signatures** - Signature-based detection rules
5. **baseline_profiles** - Behavioral baseline profiles
6. **anomaly_feedback** - User feedback on anomalies
7. **threat_statistics** - Aggregated threat statistics

### Indexes and Performance

Each table includes appropriate indexes for:
- Primary key lookups
- Common query patterns
- Time-based queries
- Filtering operations
- Composite queries

### Constraints and Validation

Tables include:
- Check constraints for enum-like fields
- Foreign key constraints where appropriate
- Not null constraints for required fields
- Unique constraints for business rules

### JSONB Usage

Several tables use JSONB for flexible data storage:
- `threats.headers` - HTTP headers
- `threats.parameters` - Request parameters
- `threats.analysis_result` - Detailed analysis results
- `threats.recommendations` - Recommended actions
- `behavior_patterns.pattern_data` - Pattern-specific data
- `baseline_profiles.baseline_data` - Baseline metrics

This allows for flexible schema evolution without requiring new migrations for additional fields.
