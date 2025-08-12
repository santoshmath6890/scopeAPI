BIN_DIR=../bin
CONFIG_DIR=config

.PHONY: help all api-discovery data-ingestion threat-detection gateway-integration admin-console data-protection attack-blocking clean

# Default target
help:
	@echo "ScopeAPI Backend Services - Available Commands:"
	@echo ""
	@echo "Build Commands:"
	@echo "  all              - Build all services"
	@echo "  api-discovery    - Build API Discovery service"
	@echo "  data-ingestion  - Build Data Ingestion service"
	@echo "  threat-detection - Build Threat Detection service"
	@echo "  gateway-integration - Build Gateway Integration service"
	@echo "  admin-console   - Build Admin Console service"
	@echo "  data-protection - Build Data Protection service"
	@echo "  attack-blocking - Build Attack Blocking service"
	@echo ""
	@echo "Utilities:"
	@echo "  clean            - Clean all build artifacts"
	@echo "  help             - Show this help message"

all: api-discovery data-ingestion threat-detection gateway-integration admin-console data-protection attack-blocking

api-discovery:
	@cd services/api-discovery && go build -o ../../bin/api-discovery ./cmd

data-ingestion:
	@cd services/data-ingestion && go build -o ../../bin/data-ingestion ./cmd

threat-detection:
	@cd services/threat-detection && go build -o ../../bin/threat-detection ./cmd

gateway-integration:
	@cd services/gateway-integration && go build -o ../../bin/gateway-integration ./cmd

admin-console:
	@cd services/admin-console && go build -o ../../bin/admin-console ./cmd

data-protection:
	@cd services/data-protection && go build -o ../../bin/data-protection ./cmd

attack-blocking:
	@cd services/attack-blocking && go build -o ../../bin/attack-blocking ./cmd

clean:
	rm -f $(BIN_DIR)/* 