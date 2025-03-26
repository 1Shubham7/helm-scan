# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o helm-scan ./cmd/helm-scan/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates curl bash

# Install Helm
RUN wget https://get.helm.sh/helm-v3.14.2-linux-amd64.tar.gz && \
    tar -zxvf helm-v3.14.2-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -rf linux-amd64 helm-v3.14.2-linux-amd64.tar.gz && \
    helm version

WORKDIR /app

COPY --from=builder /app/helm-scan .

EXPOSE 8080

CMD ["/app/helm-scan"]