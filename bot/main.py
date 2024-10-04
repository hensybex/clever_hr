# main.py

import asyncio
from aiogram import Bot, Dispatcher
from config import TELEGRAM_BOT_TOKEN
from handlers.common import router as common_router
from handlers.employee import router as employee_router
from handlers.candidate import router as candidate_router
from handlers.interview import router as interview_router

async def main():
    bot = Bot(token=TELEGRAM_BOT_TOKEN)
    dp = Dispatcher()

    # Include routers
    dp.include_router(common_router)
    dp.include_router(employee_router)
    dp.include_router(candidate_router)
    dp.include_router(interview_router)

    # Start polling
    await dp.start_polling(bot)

if __name__ == '__main__':
    asyncio.run(main())
