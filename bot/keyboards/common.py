# keyboards/common.py

from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup

def start_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [
            InlineKeyboardButton(text="HR", callback_data="register_employee"),
            InlineKeyboardButton(text="Кандидат", callback_data="register_candidate")
        ]
    ])
    return keyboard
