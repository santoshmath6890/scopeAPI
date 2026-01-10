# Multi-stage build for Data Ingestion Service
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
# Copy shared modules
COPY shared/ ./shared/
COPY services/data-ingestion/go.mod services/data-ingestion/go.sum* ./services/data-ingestion/

# Set working directory to service directory for dependency download
WORKDIR /app/services/data-ingestion

# Download dependencies
RUN go mod download

# Copy source code
COPY services/data-ingestion/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o data-ingestion cmd/main.go

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
COPY --from=builder /app/services/data-ingestion/data-ingestion .

# Copy configuration files
COPY --from=builder /app/services/data-ingestion/config* ./config/

# Change ownership to non-root user
RUN chown -R scopeapi:scopeapi /app

# Switch to non-root user
USER scopeapi

# Expose ports
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application
CMD ["./data-ingestion"]
