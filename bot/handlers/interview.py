# handlers/interview.py

from aiogram import types, Router, F
from aiogram.filters import Command
from utils.api_client import api_client

router = Router()

@router.message(Command(commands=['create_interview']))
async def cmd_create_interview(message: types.Message):
    # Parse arguments for resume_id and interview_type_id
    args = message.get_args().split()
    if len(args) != 2:
        await message.reply("Использование: /create_interview <resume_id> <interview_type_id>")
        return
    resume_id, interview_type_id = map(int, args)
    
    # Call the API and await the response
    response = await api_client.create_interview(resume_id, interview_type_id)
    await message.reply(
        f"Интервью успешно создано (ID: {response['interview_id']})." if response.get('success') else "Ошибка при создании интервью."
    )

@router.message(Command(commands=['analyze_interview']))
async def cmd_analyze_interview(message: types.Message):
    # Parse interview_id from arguments
    args = message.get_args()
    if not args:
        await message.reply("Использование: /analyze_interview <interview_id>")
        return
    interview_id = int(args)
    
    # Call the API and await the response
    response = await api_client.analyze_interview(interview_id)
    await message.reply(
        "Анализ интервью начат." if response.get('success') else "Ошибка при запуске анализа интервью."
    )

""" @router.message(Command(commands=['get_interview_analysis']))
async def cmd_get_interview_analysis_result(message: types.Message):
    args = message.get_args()
    if not args:
        await message.reply("Использование: /get_interview_analysis <interview_id>")
        return
    interview_id = int(args)
    
    # Call the API and await the response
    response = await api_client.get_interview_analysis_result(interview_id)
    if 'result' in response:
        await message.reply(f"Результаты анализа интервью:\n{response['result']}")
    else:
        await message.reply("Анализ интервью еще не завершен.")
 """