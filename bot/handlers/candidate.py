# handlers/candidate.py

from aiogram import types, Router, F
from aiogram.filters import Command
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from utils.api_client import api_client
from keyboards.candidate import (
    candidate_main_menu_keyboard,
    upload_candidate_keyboard,
    candidate_profile_keyboard,
    before_interview_keyboard,
    interview_keyboard
)

from keyboards.employee import (
    employee_main_menu_keyboard,
)

router = Router()

class InterviewState(StatesGroup):
    in_progress = State()

import logging

# Register candidate
@router.message(Command(commands=['candidate']))
async def cmd_register_candidate(callback_query: types.CallbackQuery):
    user_id = callback_query.from_user.id  # Ensure we are using the user's TG ID from the callback query
    
    logging.info(f"Registering user with TG ID: {user_id} as employee.")
    
    response = await api_client.create_user(tg_id=user_id, user_type='candidate')
    logging.info(f"API Response for user registration: {response}")
    
    await callback_query.message.answer(
        "Вы успешно зарегистрировались как кандидат." if response.get('success') else "Вы уже зарегистрированы.",
        reply_markup=employee_main_menu_keyboard()
    )


# Handle inline buttons for candidate main menu
@router.callback_query(F.data == 'candidate_main_menu')
async def candidate_main_menu(callback_query: types.CallbackQuery):
    await callback_query.message.edit_text("Главное меню кандидата:", reply_markup=candidate_main_menu_keyboard())

# Handle resume upload for candidate
@router.callback_query(F.data == 'candidate_upload_resume')
async def candidate_upload_resume(callback_query: types.CallbackQuery):
    await callback_query.message.edit_text("Пожалуйста, отправьте ваш PDF файл резюме.", reply_markup=upload_candidate_keyboard())

# Handle document upload
@router.message(F.document.mime_type == 'application/pdf')
async def handle_resume_document(message: types.Message):
    file_info = await message.bot.get_file(message.document.file_id)
    file = await message.bot.download_file(file_info.file_path)
    file_bytes = file.read()

    response = await api_client.upload_resume_candidate(
        resume_file=file_bytes,
        tg_id=message.from_user.id
    )

    await message.reply(
        "Ваше резюме успешно загружено и отправлено на анализ." if response.get('success') else "Произошла ошибка при загрузке резюме.",
        reply_markup=candidate_main_menu_keyboard()
    )

# Profile check for candidate
@router.callback_query(F.data == 'candidate_profile')
async def candidate_profile(callback_query: types.CallbackQuery):
    candidate_info = await api_client.get_candidate_info_by_tg_id(callback_query.from_user.id)
    if candidate_info:
        await callback_query.message.edit_text(
            f"Ваш профиль:\nИмя: {candidate_info['name']}\nEmail: {candidate_info['email']}",
            reply_markup=candidate_profile_keyboard()
        )
    else:
        await callback_query.message.edit_text("Профиль не найден.", reply_markup=candidate_main_menu_keyboard())

# Start interview flow for candidate
@router.callback_query(F.data == 'candidate_start_interview')
async def candidate_start_interview(callback_query: types.CallbackQuery):
    interview_types = await api_client.list_interview_types()
    if interview_types:
        keyboard = before_interview_keyboard(interview_types)
        await callback_query.message.edit_text("Выберите тип интервью:", reply_markup=keyboard)
    else:
        await callback_query.message.edit_text("Типы интервью не найдены.", reply_markup=candidate_main_menu_keyboard())

# Handle interview state and messaging
@router.callback_query(F.data.startswith('candidate_interview_type_'))
async def start_interview(callback_query: types.CallbackQuery, state: FSMContext):
    interview_type_id = int(callback_query.data.split('_')[-1])
    response = await api_client.create_interview_for_candidate(callback_query.from_user.id, interview_type_id)
    if response.get('success'):
        await state.set_state(InterviewState.in_progress)
        await state.update_data(interview_id=response['interview_id'])
        await callback_query.message.edit_text("Интервью начато. Вы можете отправлять сообщения.", reply_markup=interview_keyboard())
    else:
        await callback_query.message.edit_text("Ошибка при создании интервью.", reply_markup=candidate_main_menu_keyboard())

# Handle messages during the interview
@router.message(InterviewState.in_progress)
async def handle_interview_message(message: types.Message, state: FSMContext):
    data = await state.get_data()
    interview_id = data.get('interview_id')
    user_message = message.text

    async def on_message_callback(response):
        await message.reply(f"Ответ:\n{response}")

    await api_client.analyze_interview_message_websocket(interview_id, user_message, on_message_callback)

# Stop the interview
@router.callback_query(F.data == 'stop_interview')
async def stop_interview(callback_query: types.CallbackQuery, state: FSMContext):
    await state.clear()
    await callback_query.message.edit_text("Интервью завершено.", reply_markup=candidate_main_menu_keyboard())

# Switching user type for candidate
@router.callback_query(F.data == 'switch_user_type_candidate')
async def switch_user_type(callback_query: types.CallbackQuery):
    response = await api_client.switch_user_type(callback_query.from_user.id)
    await callback_query.message.edit_text(
        "Тип пользователя успешно изменен - сотрудник",
        reply_markup=employee_main_menu_keyboard() if response.get('message') else candidate_main_menu_keyboard()
    )