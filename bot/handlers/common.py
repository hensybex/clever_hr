# handlers/common.py

import logging
from aiogram import types, Router, F
from aiogram.filters import CommandStart
from keyboards.common import start_keyboard
from handlers.employee import cmd_register_employee
from handlers.candidate import cmd_register_candidate
from utils.api_client import api_client

from keyboards.candidate import (
    candidate_main_menu_keyboard,
)

from keyboards.employee import (
    employee_main_menu_keyboard,
)

router = Router()

@router.message(CommandStart())
async def cmd_start(message: types.Message):
    response = await api_client.get_user_role(tg_id=message.from_user.id)
    message = "Добро пожаловать в бот для анализа резюме!",
    logging.info(response)
    user_role = response.get("role", "none")
    logging.info(user_role)
    if user_role == 'employee':
        await message.answer(
            message,
            reply_markup=employee_main_menu_keyboard()
        )
    elif user_role == 'candidate':
        await message.answer(
            message,
            reply_markup=candidate_main_menu_keyboard()
        )
    else:
        await message.answer(
            message,
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


@router.message(F.document.mime_type == 'application/pdf')
async def handle_resume_document(message: types.Message):
    # Retrieve the user role from the API, database, or context (this is just an example, you can adjust it)
    user_role = await api_client.get_user_role(tg_id=message.from_user.id)

    # Download the file
    file_info = await message.bot.get_file(message.document.file_id)
    file = await message.bot.download_file(file_info.file_path)
    file_bytes = file.read()

    # Handle based on user role
    if user_role == 'employee':
        response = await api_client.upload_resume_employee(
            resume_file=file_bytes,
            user_id=message.from_user.id
        )
        await message.reply(
            "Резюме успешно загружено и отправлено на анализ." if response.get('success') else "Произошла ошибка при загрузке резюме.",
            reply_markup=employee_main_menu_keyboard()
        )
    elif user_role == 'candidate':
        response = await api_client.upload_resume_candidate(
            resume_file=file_bytes,
            tg_id=message.from_user.id
        )
        await message.reply(
            "Ваше резюме успешно загружено и отправлено на анализ." if response.get('success') else "Произошла ошибка при загрузке резюме.",
            reply_markup=candidate_main_menu_keyboard()
        )
    else:
        await message.reply(
            "Неизвестная роль пользователя. Пожалуйста, обратитесь в службу поддержки."
        )
