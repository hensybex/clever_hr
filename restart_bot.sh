#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting Telegram bot service..."

# Bring down the current running telegram-bot container
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down telegram-bot

# Rebuild and bring up only the telegram-bot service
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d telegram-bot

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "Telegram bot service restarted successfully. Tracking logs..."
    
    # Start streaming the logs for the telegram-bot service
    docker-compose logs -f telegram-bot
else
    echo "Failed to restart the Telegram bot service."
    exit 1
fi
