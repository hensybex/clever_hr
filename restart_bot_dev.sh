# restart_bot_dev.sh

#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting Telegram bot service..."

# Bring up the api and db services to ensure the network is created
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d db api

# Stop the telegram-bot service without removing the network
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml stop telegram-bot

# Remove the telegram-bot service container
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml rm -f telegram-bot

# Rebuild and bring up the telegram-bot service
sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d telegram-bot

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "Telegram bot service restarted successfully. Tracking logs..."
    
    # Start streaming the logs for the telegram-bot service
    sudo docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f telegram-bot
else
    echo "Failed to restart the Telegram bot service."
    exit 1
fi