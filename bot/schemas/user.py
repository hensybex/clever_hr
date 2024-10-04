# schemas/user.py

from dataclasses import dataclass

@dataclass
class UserCreate:
    tg_id: int
    user_type: str

@dataclass
class UserResponse:
    id: int
    tg_id: int
    user_type: str
    created_at: str
    updated_at: str
