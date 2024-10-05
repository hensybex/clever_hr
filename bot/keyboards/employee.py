# keyboards/employee.py

from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup, ReplyKeyboardMarkup, KeyboardButton
from aiogram.utils.keyboard import InlineKeyboardBuilder
import logging

def employee_main_menu_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Загрузить резюме кандидата", callback_data="employee_upload_resume")],
        [InlineKeyboardButton(text="Список кандидатов", callback_data="employee_list_candidates")],
        [InlineKeyboardButton(text="Сменить тип пользователя", callback_data="switch_user_type_employee")]
    ])
    return keyboard

def upload_employee_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Назад", callback_data="employee_main_menu")]
    ])
    return keyboard

def employee_candidates_list_keyboard(candidates: list, page: int = 1) -> InlineKeyboardMarkup:
    keyboard_builder = InlineKeyboardBuilder()

    # Log the candidates to see their structure
    logging.info(f"Candidates: {candidates}, Type: {type(candidates)}")
    
    # Add candidate buttons (limit to 10 per page)
    for candidate in candidates[:10]:
        
        # Use capitalized keys for accessing candidate information
        candidate_name = candidate.get('Name', 'Имя не указано')
        candidate_id = candidate.get('ID', 'unknown_id')  # Fallback in case 'ID' is missing
        
        keyboard_builder.add(
            InlineKeyboardButton(text=candidate_name, callback_data=f"employee_candidate_{candidate_id}")
        )

    # Pagination buttons
    pagination_buttons = []
    if len(candidates) > 10:
        if page > 1:
            pagination_buttons.append(
                InlineKeyboardButton(text="⬅️ Назад", callback_data=f"employee_candidates_page_{page-1}")
            )
        if len(candidates) > page * 10:
            pagination_buttons.append(
                InlineKeyboardButton(text="Вперёд ➡️", callback_data=f"employee_candidates_page_{page+1}")
            )
    
    # Add pagination buttons if they exist
    if pagination_buttons:
        keyboard_builder.row(*pagination_buttons)

    # Back button
    keyboard_builder.add(
        InlineKeyboardButton(text="Назад", callback_data="employee_main_menu")
    )

    # Adjust the layout of the buttons (1 button per row)
    keyboard_builder.adjust(1)

    # Return the keyboard with the constructed layout
    return keyboard_builder.as_markup()


def employee_candidate_check_keyboard(candidate_id: int, resume_id: int, resume_url: str = None, resume_was_analyzed: bool = False) -> InlineKeyboardMarkup:
    if resume_url:
        if resume_was_analyzed:
            keyboard = InlineKeyboardMarkup(inline_keyboard=[
                [InlineKeyboardButton(text="Назад", callback_data="employee_list_candidates")],
                [InlineKeyboardButton(text="Скачать резюме", callback_data=f"download_resume_{candidate_id}")],
                [InlineKeyboardButton(text="Посмотреть результаты анализа резюме", callback_data=f"check_resume_analysis_{resume_id}")],
                [InlineKeyboardButton(text="Главное меню", callback_data="employee_main_menu")]
            ])
        else:
            keyboard = InlineKeyboardMarkup(inline_keyboard=[
                [InlineKeyboardButton(text="Назад", callback_data="employee_list_candidates")],
                [InlineKeyboardButton(text="Скачать резюме", callback_data=f"download_resume_{candidate_id}")],
                [InlineKeyboardButton(text="Запустить анализ резюме", callback_data=f"run_resume_analysis_{resume_id}")],
                [InlineKeyboardButton(text="Главное меню", callback_data="employee_main_menu")]
            ])
    else:
        keyboard = InlineKeyboardMarkup(inline_keyboard=[
            [InlineKeyboardButton(text="Назад", callback_data="employee_list_candidates")],
            [InlineKeyboardButton(text="Главное меню", callback_data="employee_main_menu")]
        ])
    return keyboard

def employee_resume_analysis_keyboard(candidate_id: int) -> InlineKeyboardMarkup:
    keyboard_builder = InlineKeyboardBuilder()
    keyboard_builder.add(
        InlineKeyboardButton(text="Назад", callback_data=f"employee_candidate_{candidate_id}")
    )
    return keyboard_builder.as_markup()

def main_menu_reply_keyboard() -> ReplyKeyboardMarkup:
    """
    Creates a persistent reply keyboard with a 'Main Menu' button.
    """
    keyboard = ReplyKeyboardMarkup(resize_keyboard=True)
    keyboard.add(KeyboardButton(text="Главное меню"))  # "Main Menu" button
    return keyboard