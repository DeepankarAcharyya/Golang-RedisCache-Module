#!/bin/bash

# Exit on any error
set -e

# Script constants
CONTAINER_NAME="redis-test-instance"
IMAGE_NAME="my-redis-instance"
PORT=6379

echo "Starting Redis test environment..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "Error: Docker is not running or not accessible"
    exit 1
fi

# Stop and remove existing container if it exists
if docker ps -a | grep -q $CONTAINER_NAME; then
    echo "Stopping and removing existing container..."
    docker stop $CONTAINER_NAME >/dev/null 2>&1 || true
    docker rm $CONTAINER_NAME >/dev/null 2>&1 || true
fi

# Build the Docker image
echo "Building Redis Docker image..."
docker build -t $IMAGE_NAME . || {
    echo "Error: Failed to build Docker image"
    exit 1
}

# Run the container
echo "Starting Redis container..."
docker run -d \
    --name $CONTAINER_NAME \
    -p $PORT:6379 \
    $IMAGE_NAME || {
    echo "Error: Failed to start container"
    exit 1
}

# Wait for Redis to be ready
echo "Waiting for Redis to be ready..."
for i in {1..30}; do
    if docker exec $CONTAINER_NAME redis-cli ping >/dev/null 2>&1; then
        echo "Redis is ready!"
        echo "Container '$CONTAINER_NAME' is running on port $PORT"
        exit 0
    fi
    sleep 1
done

echo "Error: Redis failed to start within 30 seconds"
exit 1
