#!/bin/sh

# Start Docker daemon in the background
dockerd &

# Wait for Docker to be ready
until docker info >/dev/null 2>&1; do
  echo "Waiting for Docker to start..."
  sleep 2
done

# Start the application
exec /app/helm-scan
