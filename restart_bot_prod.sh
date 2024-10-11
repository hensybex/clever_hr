# restart_bot_dev.sh

#!/bin/bash

# Function to handle cleanup and stopping of logs on script termination
cleanup() {
    echo "Stopping log tracking..."
    exit 0
}

# Trap SIGINT (Ctrl+C) and call cleanup function
trap cleanup SIGINT

echo "Rebuilding and restarting Telegram bot service in production..."

# Stop the telegram-bot service without removing the network
sudo docker-compose -f docker-compose.yml stop telegram-bot

# Remove the telegram-bot service container
sudo docker-compose -f docker-compose.yml rm -f telegram-bot

# Rebuild and bring up only the telegram-bot service without affecting dependencies
sudo docker-compose -f docker-compose.yml up --build -d --no-deps telegram-bot

# Check the status of the service
if [ $? -eq 0 ]; then
    echo "Telegram bot service restarted successfully in production. Tracking logs..."
    
    # Start streaming the logs for the telegram-bot service
    sudo docker-compose -f docker-compose.yml logs -f telegram-bot
else
    echo "Failed to restart the Telegram bot service in production."
    exit 1
fi
