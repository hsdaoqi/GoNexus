import os
from dotenv import load_dotenv
load_dotenv()

class Settings:
    API_KEY = os.getenv("AI_API_KEY")
    BASE_URL = os.getenv("AI_BASE_URL")
    MODEL_NAME = os.getenv("AI_MODEL_NAME")
    DB_PATH = os.getenv("CHROMA_DB_PATH")

settings = Settings()