# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy proto files (needed for local replace directive)
COPY proto ./proto

# Copy cart service source code first
COPY services/cart .

# Update go.mod to use local proto path for Docker context
RUN sed -i 's|=> ../../proto/catalog|=> ./proto/catalog|g' go.mod

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cart-service .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/cart-service .

# Make binary executable
RUN chmod +x cart-service

# Expose port
EXPOSE 8080

# Run the application
CMD ["./cart-service"]