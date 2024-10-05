# keyboards/candidate.py

from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup
from aiogram.utils.keyboard import InlineKeyboardBuilder
import logging

def candidate_main_menu_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Загрузить резюме", callback_data="candidate_upload_resume")],
        [InlineKeyboardButton(text="Профиль", callback_data="candidate_profile")],
        [InlineKeyboardButton(text="Начать интервью", callback_data="candidate_start_interview")],
        [InlineKeyboardButton(text="Сменить тип пользователя", callback_data="switch_user_type_candidate")]
    ])
    return keyboard

def upload_candidate_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")]
    ])
    return keyboard


def candidate_profile_keyboard(candidate_id: int, resume_id: int, resume_url: str = None, resume_was_analyzed: bool = False) -> InlineKeyboardMarkup:
    if resume_url:
        if resume_was_analyzed:
            keyboard = InlineKeyboardMarkup(inline_keyboard=[
                [InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")],
                [InlineKeyboardButton(text="Скачать резюме", callback_data=f"download_resume_{candidate_id}")],
                [InlineKeyboardButton(text="Посмотреть результаты анализа резюме", callback_data=f"check_resume_analysis_{resume_id}")],
            ])
        else:
            keyboard = InlineKeyboardMarkup(inline_keyboard=[
                [InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")],
                [InlineKeyboardButton(text="Скачать резюме", callback_data=f"download_resume_{candidate_id}")],
                [InlineKeyboardButton(text="Запустить анализ резюме", callback_data=f"run_resume_analysis_{resume_id}")],
            ])
    else:
        keyboard = InlineKeyboardMarkup(inline_keyboard=[
            [InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")],
        ])
    return keyboard

def before_interview_keyboard(interview_types: list, page: int =1) -> InlineKeyboardMarkup:
    keyboard_builder = InlineKeyboardBuilder()

    logging.info(f"Interview types: {interview_types}")

    for interview_type in interview_types[:10]:
        keyboard_builder.add(
            InlineKeyboardButton(text=interview_type['Name'], callback_data=f"candidate_interview_type_{interview_type['ID']}")
        )
    
    # Pagination Buttons
    pagination_buttons = []
    if len(interview_types) > 10:
        if page > 1:
            pagination_buttons.append(InlineKeyboardButton(text="⬅️ Назад", callback_data=f"before_interview_page_{page-1}"))
        if len(interview_types) > page * 10:
            pagination_buttons.append(InlineKeyboardButton(text="Вперёд ➡️", callback_data=f"before_interview_page_{page+1}"))
    
    if pagination_buttons:
        keyboard_builder.row(*pagination_buttons)
    
    # Back button
    keyboard_builder.add(
        InlineKeyboardButton(text="Назад", callback_data="employee_main_menu")
    )

    keyboard_builder.adjust(1)
    
    return keyboard_builder.as_markup()


def interview_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Закончить интервью", callback_data="stop_interview")]
    ])
    return keyboard
