# utils/api_client.py

import aiohttp
import json
from config import API_BASE_URL
import logging
class APIClient:
    def __init__(self):
        self.base_url = API_BASE_URL

    async def create_user(self, tg_id, user_type):
        """
        Creates a new user (employee or candidate).
        """
        url = f"{self.base_url}/users"
        data = {'tg_id': tg_id, 'user_type': user_type}
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data) as response:
                return await response.json()
    
    async def get_user_role(self, tg_id):
        """
        Creates a new user (employee or candidate).
        """
        url = f"{self.base_url}/users/{tg_id}/get_role"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def get_user(self, user_id):
        """
        Retrieves user information by ID.
        """
        url = f"{self.base_url}/users/{user_id}"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def upload_resume(self, resume_file, user_id):
        """
        Uploads a resume for a candidate by an employee (PDF).
        """
        url = f"{self.base_url}/resumes/upload"
        data = aiohttp.FormData()
        data.add_field('resume', resume_file, filename='resume.pdf', content_type='application/pdf')
        data.add_field('tg_id', str(user_id))

        async with aiohttp.ClientSession() as session:
            async with session.post(url, data=data) as response:
                return await response.json()


    async def analyze_resume(self, resume_id):
        """
        Triggers the resume analysis process using the LLM.
        """
        url = f"{self.base_url}/resumes/{resume_id}/analyze"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def get_resume_analysis_result(self, resume_id):
        """
        Retrieves the result of the resume analysis after it's processed.
        """
        url = f"{self.base_url}/resumes/{resume_id}/analysis-result"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def create_interview(self, resume_id, interview_type_id):
        """
        Creates a new interview based on a resume and interview type.
        """
        url = f"{self.base_url}/interviews"
        data = {'resume_id': resume_id, 'interview_type_id': interview_type_id}
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data) as response:
                return await response.json()

    async def analyze_interview(self, interview_id):
        """
        Runs a full analysis of an interview based on past messages.
        """
        url = f"{self.base_url}/interviews/{interview_id}/analyze"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def get_interview_analysis_result(self, interview_id):
        """
        Retrieves the result of the interview analysis.
        """
        url = f"{self.base_url}/interviews/{interview_id}/analysis-result"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def analyze_interview_message_websocket(self, interview_id, messages, on_message_callback):
        """
        Establishes a WebSocket connection for real-time interview message analysis.
        """
        uri = f"{self.base_url.replace('http', 'ws')}/ws/interview/analyse"
        async with aiohttp.ClientSession() as session:
            async with session.ws_connect(uri) as ws:
                await ws.send_str(json.dumps({'interview_id': interview_id, 'messages': messages}))

                async for msg in ws:
                    if msg.type == aiohttp.WSMsgType.TEXT:
                        await on_message_callback(msg.data)
                    elif msg.type == aiohttp.WSMsgType.ERROR:
                        break

    async def switch_user_type(self, user_id):
        """
        Switches the user type.
        """
        url = f"{self.base_url}/users/{user_id}/switch"
        async with aiohttp.ClientSession() as session:
            async with session.put(url) as response:
                return await response.json()

    async def get_uploaded_candidates(self, user_id):
        """
        Fetches candidates uploaded by the employee.
        """
        url = f"{self.base_url}/users/{user_id}/candidates"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def get_candidate_info(self, candidate_id):
        """
        Fetches candidate information by candidate ID.
        """
        url = f"{self.base_url}/candidates/{candidate_id}/get_by_id"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def get_candidate_info_by_tg_id(self, tg_id):
        """
        Fetches candidate information by Telegram ID.
        """
        url = f"{self.base_url}/candidates/{tg_id}/get_by_tg_id"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def get_candidate_resume(self, candidate_id):
        """
        Retrieves the resume PDF of a candidate by their ID and resume path.
        """
        # Construct the API endpoint with the candidate ID and resume path as query parameters
        candidate_info = await self.get_candidate_info(candidate_id)
        if 'resume_pdf' in candidate_info:
            # Remove any leading /app from the path
            resume_path = candidate_info['resume_pdf'].replace("/app", "").lstrip('/')
            async with aiohttp.ClientSession() as session:
                async with session.get(f"{self.base_url}/candidates/{candidate_id}/resume", params={"resume_path": resume_path}) as response:
                    if response.status == 200:
                        return await response.read()  # Return binary PDF content
                    else:
                        return None



    async def list_interview_types(self):
        """
        Retrieves a list of available interview types.
        """
        url = f"{self.base_url}/interview-types"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as response:
                return await response.json()

    async def create_interview_for_candidate(self, tg_id, interview_type_id):
        """
        Creates an interview for the candidate.
        """
        url = f"{self.base_url}/interviews"
        data = {'tg_id': tg_id, 'interview_type_id': interview_type_id}
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=data) as response:
                return await response.json()

# Instantiate the API client
api_client = APIClient()
