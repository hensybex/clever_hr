# schemas/resume.py

from dataclasses import dataclass

@dataclass
class ResumeUploadResponse:
    success: bool
    resume_id: int = None
    message: str = None

@dataclass
class ResumeAnalysisResult:
    resume_id: int
    analysis_status: str
    results: str
