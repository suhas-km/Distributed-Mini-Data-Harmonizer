"""
SQLAlchemy models for job tracking.
"""

import enum
import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Column, String, DateTime, Enum, Text
from sqlalchemy.orm import relationship

from database import Base


class JobStatus(str, enum.Enum):
    """Job status enumeration."""
    QUEUED = "queued"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"


class Job(Base):
    """Job model for tracking data harmonization tasks."""
    
    __tablename__ = "jobs"
    
    id = Column(String(36), primary_key=True, index=True, default=lambda: str(uuid.uuid4()))
    status = Column(Enum(JobStatus), default=JobStatus.QUEUED, nullable=False)
    
    # File paths
    input_file = Column(String(255), nullable=False)
    output_file = Column(String(255), nullable=True)
    
    # Job details
    file_type = Column(String(10), nullable=False)  # csv, json, etc.
    file_size = Column(String(20), nullable=False)  # Human-readable size
    harmonization_type = Column(String(50), nullable=False)  # patients, vitals, medications, lab_results
    
    # Timestamps
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow, nullable=False)
    completed_at = Column(DateTime, nullable=True)
    
    # Error information
    error_message = Column(Text, nullable=True)
    
    def __repr__(self):
        """String representation of the job."""
        return f"<Job {self.id}: {self.status.value}>"
