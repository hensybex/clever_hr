#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting API service..."

# Bring down the current running api container
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down api

# Rebuild and bring up only the api service
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d api

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "API service restarted successfully. Tracking logs..."
    
    # Start streaming the logs for the api service
    docker-compose logs -f api
else
    echo "Failed to restart the API service."
    exit 1
fi
