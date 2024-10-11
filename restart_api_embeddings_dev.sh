# restart_api_embeddings_dev.sh

#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting API Embeddings service..."

# Bring up the db service to ensure the network is created
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d db

# Stop the api service without removing the network
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml stop api_embeddings

# Remove the api service container
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml rm -f api_embeddings

# Rebuild and bring up the api service
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d api_embeddings

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "API Embeddings service restarted successfully. Tracking logs..."
    
    # Start streaming the logs for the api service
    sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f api_embeddings
else
    echo "Failed to restart the API Embeddings service."
    exit 1
fi
