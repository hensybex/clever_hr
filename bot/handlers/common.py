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
    candidate_resume_analysis_keyboard,
)

from keyboards.employee import (
    employee_main_menu_keyboard,
    employee_resume_analysis_keyboard,
)

router = Router()

@router.message(CommandStart())
async def cmd_start(message: types.Message):
    response = await api_client.get_user_role(tg_id=message.from_user.id)
    welcome_message = "Добро пожаловать в бот для анализа резюме!"  # Changed variable name
    logging.info(response)
    user_role = response.get("role", "none")
    logging.info(user_role)

    if user_role == 'employee':
        await message.answer(
            welcome_message,  # Use the new variable here
            reply_markup=employee_main_menu_keyboard()
        )
    elif user_role == 'candidate':
        await message.answer(
            welcome_message,  # Use the new variable here
            reply_markup=candidate_main_menu_keyboard()
        )
    else:
        await message.answer(
            welcome_message,  # Use the new variable here
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
    response = await api_client.get_user_role(tg_id=message.from_user.id)
    user_role = response.get("role", "none")

    # Download the file
    file_info = await message.bot.get_file(message.document.file_id)
    file = await message.bot.download_file(file_info.file_path)
    file_bytes = file.read()

    response = await api_client.upload_resume(
        resume_file=file_bytes,
        user_id=message.from_user.id
    )
    # Handle based on user role
    if user_role == 'employee':
        await message.reply(
            "Резюме успешно загружено. Вы можете отправить его на анализ во вкладке Список кандидатов -> Нужный кандидат" if response.get('success') else "Произошла ошибка при загрузке резюме.",
            reply_markup=employee_main_menu_keyboard()
        )
    elif user_role == 'candidate':
        await message.reply(
            "Ваше резюме успешно загружено. Вы можете отправить его на анализ во вкладке Профиль." if response.get('success') else "Произошла ошибка при загрузке резюме.",
            reply_markup=candidate_main_menu_keyboard()
        )
    else:
        await message.reply(
            "Неизвестная роль пользователя. Пожалуйста, обратитесь в службу поддержки."
        )

@router.callback_query(F.data.startswith('check_resume_analysis_'))
async def check_resume_analysis(callback_query: types.CallbackQuery):
    # Extract resume_id from the callback data
    response = await api_client.get_user_role(tg_id = callback_query.from_user.id)
    user_role = response.get("role", "none")
    is_employee = user_role == 'employee'

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
        if is_employee:
            await callback_query.message.edit_text(analysis_text, reply_markup=employee_resume_analysis_keyboard(candidate_id))
        else:
            await callback_query.message.edit_text(analysis_text, reply_markup=candidate_resume_analysis_keyboard())
    
    except Exception as e:
        # Log the error for debugging
        logging.error("Error fetching resume analysis: %s", str(e))
        
        # Handle potential errors (e.g., invalid resume_id or server issues)
        await callback_query.message.edit_text(f"Не удалось загрузить результаты анализа: {str(e)}")

    
""" @router.callback_query(F.data.startswith('check_interview_analysis_'))
async def check_resume_analysis(callback_query: types.CallbackQuery):
    # Extract resume_id from the callback data
    response = await api_client.get_user_role(tg_id = callback_query.from_user.id)
    user_role = response.get("role", "none")
    is_employee = user_role == 'employee'

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
        if is_employee:
            await callback_query.message.edit_text(analysis_text, reply_markup=employee_resume_analysis_keyboard(candidate_id))
        else:
            await callback_query.message.edit_text(analysis_text, reply_markup=candidate_resume_analysis_keyboard())
    
    except Exception as e:
        # Log the error for debugging
        logging.error("Error fetching resume analysis: %s", str(e))
        
        # Handle potential errors (e.g., invalid resume_id or server issues)
        await callback_query.message.edit_text(f"Не удалось загрузить результаты анализа: {str(e)}") """


@router.callback_query(F.data.startswith('check_interview_analysis_'))
async def check_interview_analysis(callback_query: types.CallbackQuery):
    # Extract interview_id from the callback data
    response = await api_client.get_user_role(tg_id=callback_query.from_user.id)
    user_role = response.get("role", "none")
    is_employee = user_role == 'employee'

    interview_id = int(callback_query.data.split('_')[-1])

    # Fetch the interview analysis result from the API
    try:
        logging.info("Fetching interview analysis for ID: %s", interview_id)
        interview_analysis = await api_client.get_interview_analysis_result(interview_id)
        logging.info("Interview Analysis Data: %s", interview_analysis)

        # Extract data from the 'result' key
        result = interview_analysis['result']
        interview_analysis_text = (
            f"Результаты анализа интервью для ID {interview_id}\n"
            f"Статус: {result.get('ResultStatus', 'N/A')}\n\n"
            f"Оценка:\n  {result.get('Assessment', 'Нет данных')}\n\n"
            f"Сильные стороны:\n  {result.get('Strengths', 'Нет данных')}\n\n"
            f"Слабые стороны:\n  {result.get('Weaknesses', 'Нет данных')}\n\n"
            f"Рекомендации:\n  {result.get('Recommendation', 'Нет данных')}\n\n"
            f"Причина:\n  {result.get('Reason', 'Нет данных')}\n"
        )

        # Send the analysis result to the user with a back button
        if is_employee:
            await callback_query.message.edit_text(interview_analysis_text, reply_markup=candidate_resume_analysis_keyboard())
        else:
            await callback_query.message.edit_text(interview_analysis_text, reply_markup=candidate_resume_analysis_keyboard())

    except Exception as e:
        # Log the error for debugging
        logging.error("Error fetching interview analysis: %s", str(e))
        
        # Handle potential errors (e.g., invalid interview_id or server issues)
        await callback_query.message.edit_text(f"Не удалось загрузить результаты анализа интервью: {str(e)}")
