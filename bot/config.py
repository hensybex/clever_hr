# config.py

import os
from dotenv import load_dotenv

load_dotenv()

# Telegram Bot Token
TELEGRAM_BOT_TOKEN = os.getenv('TELEGRAM_BOT_TOKEN')

# API Base URL
API_BASE_URL = os.getenv('API_BASE_URL', 'http://localhost:8080')
