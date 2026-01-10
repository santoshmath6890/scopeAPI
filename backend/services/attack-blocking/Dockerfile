# Multi-stage build for Attack Blocking Service
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
# Copy shared modules
COPY shared/ ./shared/
COPY services/attack-blocking/go.mod services/attack-blocking/go.sum* ./services/attack-blocking/

# Set working directory to service directory for dependency download
WORKDIR /app/services/attack-blocking

# Download dependencies
RUN go mod download

# Copy source code
COPY services/attack-blocking/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o attack-blocking cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S scopeapi && \
    adduser -u 1001 -S scopeapi -G scopeapi

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/services/attack-blocking/attack-blocking .

# Copy configuration files if they exist (ignoring errors if config dir is missing)
RUN mkdir -p /app/config
COPY --from=builder /app/services/attack-blocking/config* /app/config/

# Change ownership to non-root user
RUN chown -R scopeapi:scopeapi /app

# Switch to non-root user
USER scopeapi

# Expose ports
EXPOSE 8085

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8085/health || exit 1

# Run the application
CMD ["./attack-blocking"]
