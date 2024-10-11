# restart_api_embeddings_prod.sh

#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting API Embeddings service in production..."

# Stop the api service without removing the network
sudo docker-compose -f docker-compose.yml stop api_embeddings

# Remove the api service container
sudo docker-compose -f docker-compose.yml rm -f api_embeddings

# Rebuild and bring up only the api service without affecting dependencies
sudo docker-compose -f docker-compose.yml up --build -d --no-deps api_embeddings

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "API Embeddings service restarted successfully in production. Tracking logs..."
    
    # Start streaming the logs for the api service
    sudo docker-compose -f docker-compose.yml logs -f api_embeddings
else
    echo "Failed to restart the API Embeddings service in production."
    exit 1
fi
