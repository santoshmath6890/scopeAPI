# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
# Copy shared modules
COPY shared/ ./shared/
COPY services/gateway-integration/go.mod services/gateway-integration/go.sum* ./services/gateway-integration/

# Set working directory to service directory for dependency download
WORKDIR /app/services/gateway-integration

# Download dependencies
RUN go mod download

# Copy source code
COPY services/gateway-integration/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway-integration cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/services/gateway-integration/gateway-integration .

# Copy configuration files
COPY --from=builder /app/services/gateway-integration/config* ./config/

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8080 8081 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the application
CMD ["./gateway-integration"] 