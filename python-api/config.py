"""
Configuration settings for the Data Harmonizer API.
"""

import os
from pathlib import Path
from typing import Optional

from pydantic import Field
from pydantic_settings import BaseSettings

# Get the base directory of the project
BASE_DIR = Path(__file__).parent.parent.absolute()


class Settings(BaseSettings):
    """Application settings."""
    
    # API settings
    API_V1_PREFIX: str = "/api/v1"
    PROJECT_NAME: str = "Distributed Mini Data Harmonizer"
    DEBUG: bool = Field(default=True)
    
    # File paths
    UPLOAD_DIR: str = Field(default="/app/uploads")
    RESULTS_DIR: str = Field(default="/app/results")
    
    # Database settings
    DATABASE_URL: str = Field(
        default=f"sqlite:///{BASE_DIR}/data/harmonizer.db"
    )
    
    # Worker settings
    WORKER_URL: str = "http://go-worker:8081"
    MAX_CONCURRENT_JOBS: int = Field(default=3)
    
    # File settings
    MAX_UPLOAD_SIZE: int = Field(default=104857600)  # 100MB in bytes
    ALLOWED_EXTENSIONS: list = Field(default=["csv", "json"])
    
    class Config:
        env_file = ".env"
        case_sensitive = True


# Create settings instance
settings = Settings()

# Ensure directories exist
os.makedirs(settings.UPLOAD_DIR, exist_ok=True)
os.makedirs(settings.RESULTS_DIR, exist_ok=True)
os.makedirs(os.path.dirname(settings.DATABASE_URL.replace("sqlite:///", "")), exist_ok=True)
