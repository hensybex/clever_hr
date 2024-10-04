# Dockerfile for Telegram Bot

# Base image for Python
FROM python:3.10-slim

# Set working directory inside the container
WORKDIR /app

# Copy requirements file first (to leverage Docker layer caching)
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the bot source code into the container
COPY . .

# Expose the port (if needed) - although not mandatory for Telegram bot
EXPOSE 8000

# Define the environment variables (these will be passed through docker-compose)
ENV TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
ENV API_BASE_URL=${API_BASE_URL}

# Command to run the bot
CMD ["python", "main.py"]
