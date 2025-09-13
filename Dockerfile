# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for private repos and SSL)
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for runtime
RUN adduser -D -g '' devpod

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o devpod-provider-upcloud \
    main.go

# Final stage
FROM scratch

# Copy certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy user account
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary
COPY --from=builder /app/devpod-provider-upcloud /devpod-provider-upcloud

# Copy provider.yaml
COPY --from=builder /app/provider.yaml /provider.yaml

# Switch to non-root user
USER devpod

# Set entrypoint
ENTRYPOINT ["/devpod-provider-upcloud"]

# Add labels
LABEL org.opencontainers.image.title="UpCloud DevPod Provider"
LABEL org.opencontainers.image.description="DevPod provider for UpCloud infrastructure"
LABEL org.opencontainers.image.vendor="neuralmux"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/neuralmux/devpod-provider-upcloud"
LABEL org.opencontainers.image.documentation="https://github.com/neuralmux/devpod-provider-upcloud/blob/main/README.md"