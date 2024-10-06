# handlers/candidate.py

import asyncio
import json
import aiogram
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
    
    if response.get('success'):
        await callback_query.message.answer(
            "Вы успешно зарегистрировались как кандидат.",
            reply_markup=candidate_main_menu_keyboard()
        )
    else:
        response = await api_client.switch_user_type(callback_query.from_user.id)
        await callback_query.message.edit_text(
            "Тип пользователя успешно изменен - кандидат",
            reply_markup=candidate_main_menu_keyboard(),
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
""" @router.message(F.document.mime_type == 'application/pdf')
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
    ) """

# Profile check for candidate
@router.callback_query(F.data == 'candidate_profile')
async def candidate_profile(callback_query: types.CallbackQuery):
    candidate_info = await api_client.get_candidate_info_by_tg_id(callback_query.from_user.id)
    # Log the fetched candidate info to inspect the structure
    logging.info(f"Candidate Info received: {candidate_info}")

    if candidate_info:
        try:
            # Access the 'candidate' nested object
            candidate = candidate_info.get('candidate', {})
            
            # Safely access candidate details with defaults to avoid KeyError
            candidate_id = candidate.get('id', '0')
            candidate_name = candidate.get('Name', 'Имя не указано')
            candidate_email = candidate.get('Email', 'Email не указан')
            candidate_phone = candidate.get('Phone', 'Телефон не указан')
            preferable_job = candidate.get('PreferableJob', 'Предпочитаемая работа не указана')
            
            # Build the message text with candidate info
            candidate_message = (
                f"Ваш профиль:\n"
                f"Имя: {candidate_name}\n"
                f"Email: {candidate_email}\n"
                f"Телефон: {candidate_phone}\n"
                f"Предпочитаемая работа: {preferable_job}"
            )
        except KeyError as e:
            # Log if there are missing keys
            logging.error(f"KeyError: Missing key {e} in candidate_info: {candidate_info}")
            await callback_query.message.edit_text("Ошибка: данные кандидата неполные.", reply_markup=employee_main_menu_keyboard())
            return

        # Prepare the keyboard using a dedicated function
        resume_url = candidate_info.get('resume_pdf')
        resume_id = candidate_info.get('resume_id')
        resume_was_analyzed = candidate_info.get('was_resume_analysed')
        interview_analysis_result_id = candidate_info.get('interview_analysis_result_id')
        logging.info("BEFORE MOVE TO CANDIDATE PROFILE KEYBOARD")
        logging.info(interview_analysis_result_id)
        keyboard = candidate_profile_keyboard(candidate_id, resume_id, resume_url, resume_was_analyzed, interview_analysis_result_id)
        
        # Send candidate info with the button to download resume
        await callback_query.message.edit_text(candidate_message, reply_markup=keyboard)
    else:
        await callback_query.message.edit_text("Кандидат не найден.", reply_markup=employee_main_menu_keyboard())

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
    logging.info("!!!!!!HERE!!!!!!")
    logging.info(response)
    if response.get('success'):
        await state.set_state(InterviewState.in_progress)
        await state.update_data(interview_id=response['interview_id'])
        await callback_query.message.edit_text("Интервью начато. Вы можете отправлять сообщения.", reply_markup=interview_keyboard())
    else:
        await callback_query.message.edit_text("Перед началом интервью, загрузите резюме и отправьте его на анализ.", reply_markup=candidate_main_menu_keyboard())



@router.message(InterviewState.in_progress)
async def handle_interview_message(message: types.Message, state: FSMContext):
    data = await state.get_data()
    interview_id = data.get('interview_id')
    user_message = message.text

    full_response = ""
    reply_message = None  # Will hold the reply message object once it is sent
    last_update_time = asyncio.get_event_loop().time()

    async def on_message_callback(response_chunk):
        nonlocal full_response, reply_message, last_update_time
        try:
            # Parse the JSON response chunk
            response = json.loads(response_chunk)
            logging.info(response)
            
            # Handle 'result' chunks and accumulate them
            if 'result' in response:
                full_response += response['result']

            # Handle the end of the interview
            if 'status' in response and response['status'] == 'End of interview':
                logging.info("WTF WHY AM I HERE???")
                await message.reply("Спасибо за ваше участие! Возвращаемся в главное меню.", reply_markup=candidate_main_menu_keyboard())
                return

            # Send the initial message only after the first set of chunks is obtained
            if reply_message is None and full_response:
                reply_message = await message.reply(full_response)  # Send only the full response
            elif reply_message:
                # Update the message every 2 seconds
                current_time = asyncio.get_event_loop().time()
                if current_time - last_update_time >= 2:
                    try:
                        await reply_message.edit_text(full_response)  # Update only with the full response
                        last_update_time = current_time  # Update the last update time
                    except aiogram.utils.exceptions.MessageNotModified:
                        pass  # Ignore if message content hasn't changed

        except json.JSONDecodeError:
            pass  # If the chunk isn't a valid JSON, we skip it

    # Start the WebSocket interaction and message handling
    await api_client.analyze_interview_message_websocket(interview_id, user_message, on_message_callback)

    # Final update after WebSocket communication is complete
    if full_response and reply_message:
        try:
            await reply_message.edit_text(full_response)
        except aiogram.utils.exceptions.MessageNotModified:
            pass
"""     else:
        # If no response has been sent, send a thank you message
        await message.reply("Спасибо за ваше участие! Возвращаемся в главное меню.", reply_markup=candidate_main_menu_keyboard())


 """



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