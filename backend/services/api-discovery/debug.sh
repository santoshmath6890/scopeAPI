#!/bin/bash

# Debug script for API Discovery Service
# Starts the service with Delve debugger for remote debugging

set -e

echo "Starting API Discovery Service in debug mode..."

# Check if binary exists
if [ ! -f "./api-discovery" ]; then
    echo "Building service first..."
    go build -gcflags="all=-N -l" -o api-discovery ./cmd
fi

# Start with Delve debugger
# -l 0.0.0.0:2345: Listen on all interfaces for remote debugging
# --accept-multiclient: Allow multiple debugger connections
# --continue: Continue execution after attaching
# --headless: Run in headless mode (no TUI)
exec dlv --listen=0.0.0.0:2345 --headless=true --continue=true --accept-multiclient exec ./api-discovery
