#!/bin/bash

# ScopeAPI Database Setup Script
# This script sets up PostgreSQL database, runs migrations, and validates setup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    if [ "$status" -eq 0 ]; then
        echo -e "${GREEN}[SUCCESS]${NC} $message"
    else
        echo -e "${RED}[ERROR]${NC} $message"
    fi
}

# Function to print info
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Default configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-scopeapi}
DB_SSL_MODE=${DB_SSL_MODE:-disable}

print_info "ScopeAPI Database Setup"
print_info "======================="

# Check if PostgreSQL is installed
check_postgresql() {
    print_info "Checking PostgreSQL installation..."
    
    if command -v psql >/dev/null 2>&1; then
        print_status 0 "PostgreSQL client found"
        return 0
    else
        print_status 1 "PostgreSQL client not found"
        print_warning "Please install PostgreSQL client:"
        print_warning "  Ubuntu/Debian: sudo apt-get install postgresql-client"
        print_warning "  CentOS/RHEL: sudo yum install postgresql"
        print_warning "  macOS: brew install postgresql"
        return 1
    fi
}

# Check if PostgreSQL server is running
check_postgresql_server() {
    print_info "Checking PostgreSQL server connection..."
    
    if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1; then
        print_status 0 "PostgreSQL server is running and accessible"
        return 0
    else
        print_status 1 "PostgreSQL server is not accessible"
        print_warning "Please ensure PostgreSQL server is running:"
        print_warning "  Ubuntu/Debian: sudo systemctl start postgresql"
        print_warning "  CentOS/RHEL: sudo systemctl start postgresql"
        print_warning "  macOS: brew services start postgresql"
        print_warning "  Or start with Docker: docker run --name postgres -e POSTGRES_PASSWORD=$DB_PASSWORD -p 5432:5432 -d postgres:13"
        return 1
    fi
}

# Create database if it doesn't exist
create_database() {
    print_info "Creating database '$DB_NAME' if it doesn't exist..."
    
    # Set password for psql
    export PGPASSWORD="$DB_PASSWORD"
    
    # Check if database exists
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
        print_status 0 "Database '$DB_NAME' already exists"
    else
        # Create database
        if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;" >/dev/null 2>&1; then
            print_status 0 "Database '$DB_NAME' created successfully"
        else
            print_status 1 "Failed to create database '$DB_NAME'"
            return 1
        fi
    fi
}

# Build and run migrations
run_migrations() {
    print_info "Running database migrations..."
    
    # Change to project root directory
    cd "$(dirname "$0")/.."
    
    # Build the migration tool
    print_info "Building migration tool..."
    if go build -o bin/migrate backend/shared/database/postgresql/migrator.go; then
        print_status 0 "Migration tool built successfully"
    else
        print_status 1 "Failed to build migration tool"
        return 1
    fi
    
    # Run migrations
    print_info "Applying database migrations..."
    if ./bin/migrate migrate; then
        print_status 0 "Database migrations applied successfully"
    else
        print_status 1 "Failed to apply database migrations"
        return 1
    fi
}

# Create a simple migration runner
create_migration_runner() {
    print_info "Creating migration runner..."
    
    cat > bin/migrate.go << 'EOF'
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"scopeapi.local/backend/shared/database/postgresql"
)

func main() {
	var (
		host     = flag.String("host", "localhost", "Database host")
		port     = flag.String("port", "5432", "Database port")
		user     = flag.String("user", "postgres", "Database user")
		password = flag.String("password", "password", "Database password")
		dbname   = flag.String("dbname", "scopeapi", "Database name")
		sslmode  = flag.String("sslmode", "disable", "SSL mode")
		action   = flag.String("action", "migrate", "Action: migrate, rollback, status")
	)
	flag.Parse()

	config := postgresql.Config{
		Host:     *host,
		Port:     *port,
		User:     *user,
		Password: *password,
		DBName:   *dbname,
		SSLMode:  *sslmode,
	}

	migrator, err := postgresql.NewMigrator(config)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	migrationsDir := "backend/shared/database/postgresql/migrations"

	switch *action {
	case "migrate":
		if err := migrator.Migrate(migrationsDir); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "rollback":
		if err := migrator.Rollback(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")
	case "status":
		if err := migrator.Status(migrationsDir); err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: migrate, rollback, status")
		os.Exit(1)
	}
}
EOF

    print_status 0 "Migration runner created"
}

# Function to create test data
create_test_data() {
    print_info "Creating sample test data..."
    
    # Set password for psql
    export PGPASSWORD="$DB_PASSWORD"
    
    if [ "$basic_flag" = true ]; then
        print_info "Basic mode: Skipping test data creation (requires database schema)"
        print_status 0 "Test data creation skipped in basic mode"
        return 0
    fi
    
    # Insert test endpoints (only if not in basic mode)
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "
    INSERT INTO api_discovery.endpoints (url, method, service_name) VALUES
    ('/api/v1/users', 'GET', 'user-service'),
    ('/api/v1/users', 'POST', 'user-service'),
    ('/api/v1/auth/login', 'POST', 'auth-service'),
    ('/api/v1/health', 'GET', 'health-check')
    ON CONFLICT (url, method) DO NOTHING;
    " >/dev/null 2>&1; then
        print_status 0 "Test data created successfully"
    else
        print_warning "Failed to create test data (may already exist)"
    fi
}

# Function to validate database setup
validate_database() {
    print_info "Validating database setup..."
    
    # Set password for psql
    export PGPASSWORD="$DB_PASSWORD"
    
    # Test basic connection
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" >/dev/null 2>&1; then
        print_status 0 "Basic database connection successful"
    else
        print_status 1 "Basic database connection failed"
        return 1
    fi
    
    if [ "$basic_flag" = true ]; then
        print_info "Basic mode: Skipping schema validation (requires database schema)"
        print_status 0 "Basic validation completed successfully"
        return 0
    fi
    
    # Test schema and tables (only if not in basic mode)
    print_info "Testing database schema..."
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT COUNT(*) FROM api_discovery.endpoints;" >/dev/null 2>&1; then
        print_status 0 "Database schema is valid"
    else
        print_status 1 "Database schema validation failed"
        return 1
    fi
    
    print_status 0 "Database validation completed successfully"
}

# Function to show help
show_help() {
    echo "ScopeAPI Database Setup Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --help, -h     Show this help message"
    echo "  --verbose, -v  Enable verbose output"
    echo "  --basic        Basic setup only (skip migrations)"
    echo "  --test-data    Create sample test data after setup"
    echo "  --validate     Run validation tests after setup"
    echo ""
    echo "Environment Variables:"
    echo "  DB_HOST        Database host (default: localhost)"
    echo "  DB_PORT        Database port (default: 5432)"
    echo "  DB_USER        Database user (default: postgres)"
    echo "  DB_PASSWORD    Database password (default: password)"
    echo "  DB_NAME        Database name (default: scopeapi)"
    echo "  DB_SSL_MODE    SSL mode (default: disable)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Use default settings"
    echo "  $0 --test-data        # Setup + create test data"
    echo "  $0 --validate         # Setup + run validation tests"
    echo "  DB_HOST=192.168.1.100 $0  # Use custom host"
    echo "  $0 --verbose          # Enable verbose output"
}

# Main execution
main() {
    local create_test_data_flag=false
    local validate_flag=false
    local basic_flag=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --basic)
                basic_flag=true
                shift
                ;;
            --test-data)
                create_test_data_flag=true
                shift
                ;;
            --validate)
                validate_flag=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            --verbose|-v)
                set -x
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "Starting database setup..."
    
    # Check prerequisites
    if ! check_postgresql; then
        exit 1
    fi
    
    if ! check_postgresql_server; then
        exit 1
    fi
    
    # Create database
    if ! create_database; then
        exit 1
    fi
    
    # Create migration runner and run migrations (skip if basic mode)
    if [ "$basic_flag" = false ]; then
        create_migration_runner
        
        # Run migrations
        if ! run_migrations; then
            exit 1
        fi
    else
        print_info "Skipping migrations in basic mode"
    fi
    
    # Create test data if requested
    if [ "$create_test_data_flag" = true ]; then
        create_test_data
    fi
    
    # Validate setup if requested
    if [ "$validate_flag" = true ]; then
        validate_database
    fi
    
    print_info "Database setup completed successfully!"
    print_info "You can now start the ScopeAPI services."
}

# Run main function
main "$@" 