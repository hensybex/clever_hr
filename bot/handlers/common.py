# handlers/common.py

import logging
from aiogram import types, Router, F
from aiogram.filters import CommandStart
from keyboards.common import start_keyboard
from handlers.employee import cmd_register_employee
from handlers.candidate import cmd_register_candidate

router = Router()

@router.message(CommandStart())
async def cmd_start(message: types.Message):
    await message.answer(
        "Добро пожаловать в бот для анализа резюме!",
        reply_markup=start_keyboard()
    )

@router.callback_query(F.data.startswith("register_"))
async def process_register_callback(callback_query: types.CallbackQuery):
    user_id = callback_query.from_user.id  # Correctly grab the ID of the user who clicked the inline button
    
    logging.info(f"Received callback from user with TG ID: {user_id}")
    
    data = callback_query.data
    if data == "register_employee":
        await cmd_register_employee(callback_query)
    elif data == "register_candidate":
        await cmd_register_candidate(callback_query)
    else:
        await callback_query.answer("Неизвестная команда.")
