# keyboards/__init__.py

from .common import start_keyboard
from .employee import (
    employee_main_menu_keyboard,
    upload_employee_keyboard,
    employee_candidates_list_keyboard,
    employee_candidate_check_keyboard,
    employee_resume_analysis_keyboard
)
from .candidate import (
    candidate_main_menu_keyboard,
    upload_candidate_keyboard,
    candidate_profile_keyboard,
    before_interview_keyboard,
    interview_keyboard
)
