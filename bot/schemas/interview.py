# schemas/interview.py

from dataclasses import dataclass

@dataclass
class InterviewCreateResponse:
    success: bool
    interview_id: int = None
    message: str = None

@dataclass
class InterviewAnalysisResult:
    interview_id: int
    result_status: str
    result: str
