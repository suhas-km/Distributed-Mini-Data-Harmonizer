"""
Pydantic schemas for API request/response validation.
"""

import uuid
from datetime import datetime
from typing import Optional, List

from pydantic import BaseModel, Field, validator

from models.job import JobStatus


class JobBase(BaseModel):
    """Base schema for job data."""
    harmonization_type: str = Field(..., description="Type of data to harmonize (patients, vitals, medications, lab_results)")


class JobCreate(JobBase):
    """Schema for job creation."""
    # No additional fields needed - file will be uploaded separately


class JobUpdate(BaseModel):
    """Schema for job updates."""
    status: Optional[JobStatus] = None
    error_message: Optional[str] = None
    completed_at: Optional[datetime] = None


class JobResponse(JobBase):
    """Schema for job response."""
    id: str
    status: JobStatus
    file_type: str
    file_size: str
    input_file: str
    output_file: Optional[str] = None
    created_at: datetime
    updated_at: datetime
    completed_at: Optional[datetime] = None
    error_message: Optional[str] = None

    class Config:
        orm_mode = True


class JobStatusResponse(BaseModel):
    """Schema for job status response."""
    id: str
    status: JobStatus
    created_at: datetime
    updated_at: datetime
    completed_at: Optional[datetime] = None
    error_message: Optional[str] = None

    class Config:
        orm_mode = True


class ErrorResponse(BaseModel):
    """Schema for error responses."""
    detail: str
