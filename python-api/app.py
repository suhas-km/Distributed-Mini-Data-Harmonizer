"""
FastAPI application for the Data Harmonizer API.
"""

import logging
from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware

from config import settings
from database import engine, Base

# Create tables
Base.metadata.create_all(bind=engine)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler("api.log")
    ]
)
logger = logging.getLogger("harmonizer-api")

# Create FastAPI app
app = FastAPI(
    title=settings.PROJECT_NAME,
    description="API for distributed data harmonization",
    version="0.1.0",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url=f"{settings.API_V1_PREFIX}/openapi.json"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # For local development only
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Import and include API routers
from endpoints.jobs import router as jobs_router

app.include_router(
    jobs_router,
    prefix=f"{settings.API_V1_PREFIX}/jobs",
    tags=["jobs"]
)

@app.get("/", tags=["health"])
async def root():
    """Root endpoint for health check."""
    return {"status": "ok", "message": "Data Harmonizer API is running"}

@app.get("/health", tags=["health"])
async def health_check():
    """Health check endpoint."""
    return {
        "status": "ok",
        "api_version": "0.1.0",
        "environment": "development" if settings.DEBUG else "production"
    }
