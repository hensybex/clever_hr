#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting API service..."

# Bring up the db service to ensure the network is created
sudo docker-compose -p app -f docker-compose.yml -f docker-compose.dev.yml up -d db

echo "---1---"

# Stop the api service without removing the network
sudo docker-compose -p app -f docker-compose.yml -f docker-compose.dev.yml stop api

echo "---2---"

# Remove the api service container and any orphan containers
sudo docker-compose -p app -f docker-compose.yml -f docker-compose.dev.yml rm -f api

echo "---3---"

# Rebuild and bring up the api service
sudo docker-compose -p app -f docker-compose.yml -f docker-compose.dev.yml up --build -d api

echo "---4---"

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "API service restarted successfully. Tracking logs..."
    
    # Start streaming the logs for the api service
    sudo docker-compose -p app -f docker-compose.yml -f docker-compose.dev.yml logs -f api
else
    echo "Failed to restart the API service."
    exit 1
fi
