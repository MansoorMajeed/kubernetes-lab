# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy proto files (needed for local replace directive)
COPY proto ./proto

# Copy catalog service files
COPY services/catalog/go.mod services/catalog/go.sum ./

# Download dependencies
RUN go mod download

# Copy catalog service source code
COPY services/catalog .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o catalog-service .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/catalog-service .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./catalog-service"] 