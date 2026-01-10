# Multi-stage build for admin console microservice
FROM node:18-alpine AS angular-builder

# Set working directory
WORKDIR /app

# Copy Angular application (from project root)
COPY adminConsole/package*.json ./
RUN npm ci --only=production

# Copy Angular source code
COPY adminConsole/ ./

# Build Angular application  
RUN npm run build -- --configuration production

# Go build stage
FROM golang:1.21-alpine AS go-builder

# Set working directory
WORKDIR /app

# Copy source code first (this contains go.mod)
COPY backend/services/admin-console/ ./

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=go-builder /app/main .

# Copy Angular build from angular-builder stage
COPY --from=angular-builder /app/dist/admin-console ./dist

# Copy configuration
COPY backend/services/admin-console/config/ ./config/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
