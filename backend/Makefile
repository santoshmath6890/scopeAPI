BIN_DIR=../bin
CONFIG_DIR=config

all: api-discovery data-ingestion threat-detection

api-discovery:
	@cd services/api-discovery && go build -o ../../../bin/api-discovery ./cmd

data-ingestion:
	@cd services/data-ingestion && go build -o ../../../bin/data-ingestion ./cmd

threat-detection:
	@cd services/threat-detection && go build -o ../../../bin/threat-detection ./cmd

clean:
	rm -f $(BIN_DIR)/* 