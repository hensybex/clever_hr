# keyboards/candidate.py

from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup

def candidate_main_menu_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Загрузить резюме", callback_data="candidate_upload_resume")],
        [InlineKeyboardButton(text="Профиль", callback_data="candidate_profile")],
        [InlineKeyboardButton(text="Начать интервью", callback_data="candidate_start_interview")],
        [InlineKeyboardButton(text="Сменить тип пользователя", callback_data="switch_user_type")]
    ])
    return keyboard

def upload_candidate_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")]
    ])
    return keyboard

def candidate_profile_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")]
    ])
    return keyboard

def before_interview_keyboard(interview_types: list, page: int =1) -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup()
    for interview_type in interview_types[:10]:  # Limiting to 10 per page
        keyboard.add(
            InlineKeyboardButton(text=interview_type['name'], callback_data=f"candidate_interview_type_{interview_type['id']}")
        )
    
    # Pagination Buttons
    pagination_buttons = []
    if len(interview_types) > 10:
        if page > 1:
            pagination_buttons.append(InlineKeyboardButton(text="⬅️ Назад", callback_data=f"before_interview_page_{page-1}"))
        if len(interview_types) > page * 10:
            pagination_buttons.append(InlineKeyboardButton(text="Вперёд ➡️", callback_data=f"before_interview_page_{page+1}"))
        if pagination_buttons:
            keyboard.add(*pagination_buttons)
    
    # Back Button
    keyboard.add(
        InlineKeyboardButton(text="Назад", callback_data="candidate_main_menu")
    )
    return keyboard

def interview_keyboard():
    keyboard = InlineKeyboardMarkup(inline_keyboard=[
        [InlineKeyboardButton(text="Закончить интервью", callback_data="stop_interview")]
    ])
    return keyboard
