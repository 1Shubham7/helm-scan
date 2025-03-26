# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application
RUN go build -o helm-scan ./cmd/helm-scan/main.go

# Final stage
FROM alpine:latest

# Install necessary dependencies
RUN apk add --no-cache ca-certificates curl bash

# Install Helm
RUN wget https://get.helm.sh/helm-v3.14.2-linux-amd64.tar.gz && \
    tar -zxvf helm-v3.14.2-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf linux-amd64 helm-v3.14.2-linux-amd64.tar.gz && \
    helm version

# Set working directory
WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/helm-scan .

# Optional: Copy any additional files you might need (config, templates, etc.)
# COPY --from=builder /app/some-additional-file /root/

# Expose port (adjust to your Gin application's port)
EXPOSE 8080

# Command to run the executable
CMD ["/app/helm-scan"]