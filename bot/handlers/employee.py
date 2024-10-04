# handlers/employee.py
import io
from aiogram import types, Router, F
from aiogram.filters import Command
from utils.api_client import api_client
from keyboards.employee import (
    employee_main_menu_keyboard,
    upload_employee_keyboard,
    employee_candidates_list_keyboard,
    employee_candidate_check_keyboard,
    employee_resume_analysis_keyboard,
)
import logging
from aiogram.types import BufferedInputFile
logging.basicConfig(level=logging.INFO)

router = Router()

# Correct: Using callback_query for inline buttons
@router.message(Command(commands=['employee']))
async def cmd_register_employee(callback_query: types.CallbackQuery):
    user_id = callback_query.from_user.id  # Ensure we are using the user's TG ID from the callback query
    
    logging.info(f"Registering user with TG ID: {user_id} as employee.")
    
    response = await api_client.create_user(tg_id=user_id, user_type='employee')
    logging.info(f"API Response for user registration: {response}")
    
    await callback_query.message.answer(
        "Вы успешно зарегистрировались как сотрудник." if response.get('success') else "Вы уже зарегистрированы.",
        reply_markup=employee_main_menu_keyboard()
    )

# Inline button callbacks
@router.callback_query(F.data == 'employee_main_menu')
async def employee_main_menu(callback_query: types.CallbackQuery):
    await callback_query.message.edit_text("Главное меню сотрудника:", reply_markup=employee_main_menu_keyboard())

@router.callback_query(F.data == 'employee_upload_resume')
async def employee_upload_resume(callback_query: types.CallbackQuery):
    await callback_query.message.edit_text("Пожалуйста, отправьте PDF файл резюме кандидата.", reply_markup=upload_employee_keyboard())

# Handle uploaded resume (for employees)
@router.message(F.document.mime_type == 'application/pdf')
async def handle_resume_document(message: types.Message):
    file_info = await message.bot.get_file(message.document.file_id)
    file = await message.bot.download_file(file_info.file_path)
    file_bytes = file.read()

    response = await api_client.upload_resume_employee(
        resume_file=file_bytes,
        user_id=message.from_user.id
    )

    await message.reply(
        "Резюме успешно загружено и отправлено на анализ." if response.get('success') else "Произошла ошибка при загрузке резюме.",
        reply_markup=employee_main_menu_keyboard()
    )

# Correct: Add await for api_client.get_uploaded_candidates
@router.callback_query(F.data == 'employee_list_candidates')
async def employee_list_candidates(callback_query: types.CallbackQuery):
    response = await api_client.get_uploaded_candidates(callback_query.from_user.id)
    
    # Extract the list of candidates from the response
    candidates = response.get('candidates', [])
    
    # Log the extracted candidates and their type to check if it's correct
    logging.info(f"Extracted Candidates: {candidates}, Type: {type(candidates)}")
    
    if isinstance(candidates, list) and candidates:  # Ensure candidates is a list and not empty
        keyboard = employee_candidates_list_keyboard(candidates)
        await callback_query.message.edit_text("Выберите кандидата:", reply_markup=keyboard)
    else:
        await callback_query.message.edit_text("Список кандидатов пуст.", reply_markup=employee_main_menu_keyboard())

# Correct: Add await for api_client.get_candidate_info
@router.callback_query(F.data.startswith('employee_candidate_'))
async def employee_candidate_check(callback_query: types.CallbackQuery):
    candidate_id = int(callback_query.data.split('_')[-1])
    
    # Fetch candidate information
    candidate_info = await api_client.get_candidate_info(candidate_id)
    
    # Log the fetched candidate info to inspect the structure
    logging.info(f"Candidate Info received: {candidate_info}")

    if candidate_info:
        try:
            # Access the 'candidate' nested object
            candidate = candidate_info.get('candidate', {})
            
            # Safely access candidate details with defaults to avoid KeyError
            candidate_name = candidate.get('Name', 'Имя не указано')
            candidate_email = candidate.get('Email', 'Email не указан')
            candidate_phone = candidate.get('Phone', 'Телефон не указан')
            preferable_job = candidate.get('PreferableJob', 'Предпочитаемая работа не указана')
            
            # Build the message text with candidate info
            candidate_message = (
                f"Информация о кандидате:\n"
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
        keyboard = employee_candidate_check_keyboard(candidate_id, resume_id, resume_url, resume_was_analyzed)
        
        # Send candidate info with the button to download resume
        await callback_query.message.edit_text(candidate_message, reply_markup=keyboard)
    else:
        await callback_query.message.edit_text("Кандидат не найден.", reply_markup=employee_main_menu_keyboard())




@router.callback_query(F.data == 'switch_user_type')
async def switch_user_type(callback_query: types.CallbackQuery):
    response = await api_client.switch_user_type(callback_query.from_user.id)
    await callback_query.message.edit_text(
        response.get('message', "Ошибка при смене типа пользователя."),
        reply_markup=None if response.get('message') else employee_main_menu_keyboard()
    )

@router.callback_query(F.data.startswith('download_resume_'))
async def download_resume(callback_query: types.CallbackQuery):
    candidate_id = int(callback_query.data.split('_')[-1])
    
    # Fetch the resume PDF content
    resume_file = await api_client.get_candidate_resume(candidate_id)
    
    if resume_file:
        # Log the file size for debugging
        logging.info(f"Downloaded resume file size: {len(resume_file)} bytes")
        
        # Ensure file size is what you expect
        if len(resume_file) < 100:  # This value can be adjusted according to expectations
            await callback_query.message.answer("The resume file seems to be corrupted or too small.")
            return

        # Wrap the binary data into BufferedInputFile
        input_file = BufferedInputFile(file=resume_file, filename='resume.pdf')
        
        # Send the resume PDF back to the user
        await callback_query.message.answer_document(input_file)
    else:
        await callback_query.message.edit_text("Не удалось загрузить резюме.", reply_markup=employee_main_menu_keyboard())

import logging

@router.callback_query(F.data.startswith('check_resume_analysis_'))
async def check_resume_analysis(callback_query: types.CallbackQuery):
    # Extract resume_id from the callback data
    resume_id = int(callback_query.data.split('_')[-1])
    
    # Fetch the resume analysis result from the API
    try:
        logging.info("------------------------------------HERE22------------------------------------")
        logging.info("Fetching resume analysis for ID: %s", resume_id)
        resume_analysis = await api_client.get_resume_analysis_result(resume_id)
        logging.info("Resume Analysis Data: %s", resume_analysis)
        
        # Extract data from the 'result' key
        result = resume_analysis['result']['resume_analysis_result']
        candidate_id = resume_analysis['result']['candidate_id']

        # Safely extract each field with fallback values if missing or empty
        analysis_text = (
            f"Результаты анализа резюме для ID {resume_id}\n"
            f"Статус: {result.get('AnalysisStatus', 'N/A')}\n\n"
            
            f"Профессиональное резюме и карьера:\n"
            f"  Описание: {result.get('ProfessionalSummaryAndCareerNarrative', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('ProfessionalSummaryAndCareerNarrative', {}).get('score', 0)}/10\n\n"
            
            f"Опыт работы и влияние:\n"
            f"  Описание: {result.get('WorkExperienceAndImpact', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('WorkExperienceAndImpact', {}).get('score', 0)}/10\n\n"
            
            f"Образование и непрерывное обучение:\n"
            f"  Описание: {result.get('EducationAndContinuousLearning', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('EducationAndContinuousLearning', {}).get('score', 0)}/10\n\n"
            
            f"Навыки и техническая компетентность:\n"
            f"  Описание: {result.get('SkillsAndTechnologicalProficiency', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('SkillsAndTechnologicalProficiency', {}).get('score', 0)}/10\n\n"
            
            f"Мягкие навыки и эмоциональный интеллект:\n"
            f"  Описание: {result.get('SoftSkillsAndEmotionalIntelligence', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('SoftSkillsAndEmotionalIntelligence', {}).get('score', 0)}/10\n\n"
            
            f"Лидерство, инновации и решение проблем:\n"
            f"  Описание: {result.get('LeadershipInnovationAndProblemSolvingPotential', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('LeadershipInnovationAndProblemSolvingPotential', {}).get('score', 0)}/10\n\n"
            
            f"Культурное соответствие и ценности:\n"
            f"  Описание: {result.get('CulturalFitAndValueAlignment', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('CulturalFitAndValueAlignment', {}).get('score', 0)}/10\n\n"
            
            f"Адаптивность, устойчивость и трудовая этика:\n"
            f"  Описание: {result.get('AdaptabilityResilienceAndWorkEthic', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('AdaptabilityResilienceAndWorkEthic', {}).get('score', 0)}/10\n\n"
            
            f"Языковая компетентность и коммуникация:\n"
            f"  Описание: {result.get('LanguageProficiencyAndCommunication', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('LanguageProficiencyAndCommunication', {}).get('score', 0)}/10\n\n"
            
            f"Профессиональные ассоциации и участие в сообществе:\n"
            f"  Описание: {result.get('ProfessionalAffiliationsAndCommunity', {}).get('overview', 'Нет данных')}\n"
            f"  Оценка: {result.get('ProfessionalAffiliationsAndCommunity', {}).get('score', 0)}/10\n\n"
        )
        
        # Send the analysis result to the user with a back button
        await callback_query.message.edit_text(analysis_text, reply_markup=employee_resume_analysis_keyboard(candidate_id))
    
    except Exception as e:
        # Log the error for debugging
        logging.error("Error fetching resume analysis: %s", str(e))
        
        # Handle potential errors (e.g., invalid resume_id or server issues)
        await callback_query.message.edit_text(f"Не удалось загрузить результаты анализа: {str(e)}")


@router.callback_query(F.data.startswith('run_resume_analysis_'))
async def run_resume_analysis(callback_query: types.CallbackQuery):
    # Логируем исходные данные callback_query.data
    logging.info(f"Callback query data: {callback_query.data}")

    # Пытаемся извлечь resume_id
    try:
        resume_id = int(callback_query.data.split('_')[-1])
        logging.info(f"Extracted resume_id: {resume_id}")
    except ValueError as e:
        logging.error(f"Error extracting resume_id: {e}")
        await callback_query.message.edit_text("Ошибка при извлечении идентификатора резюме.")
        return

    # Логируем попытку вызова API для анализа резюме
    logging.info(f"Calling API to analyze resume with ID: {resume_id}")

    # Вызываем API для анализа резюме
    analysis_response = await api_client.analyze_resume(resume_id)
    
    # Логируем ответ от API для отладки
    logging.info(f"Resume analysis response: {analysis_response}")

    # Проверяем ответ API и информируем пользователя
    if analysis_response.get('success'):
        await callback_query.message.edit_text("Анализ резюме успешно запущен.")
    else:
        await callback_query.message.edit_text("Не удалось запустить анализ резюме. Попробуйте позже.")
