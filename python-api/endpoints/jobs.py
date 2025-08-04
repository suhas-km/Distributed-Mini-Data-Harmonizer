"""
Job management endpoints for the Data Harmonizer API.
"""

import os
import logging
import httpx
from typing import List
from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, UploadFile, File, Form, BackgroundTasks
from fastapi.responses import FileResponse
from sqlalchemy.orm import Session

from database import get_db
from models.job import Job, JobStatus
from models.schemas import JobCreate, JobResponse, JobStatusResponse, ErrorResponse
from utils.file_utils import save_upload_file, get_harmonization_type_from_filename
from config import settings

router = APIRouter()
logger = logging.getLogger("harmonizer-api")


async def process_job_async(job_id: str, db: Session):
    """
    Process a job asynchronously by sending it to the Go worker.
    
    Args:
        job_id: ID of the job to process
        db: Database session
    """
    # Get job from database
    job = db.query(Job).filter(Job.id == job_id).first()
    if not job:
        logger.error(f"Job {job_id} not found")
        return
    
    # Update job status
    job.status = JobStatus.PROCESSING
    db.commit()
    
    try:
        # Send job to Go worker
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{settings.WORKER_URL}/process",
                json={
                    "job_id": job.id,
                    "input_file": job.input_file,
                    "harmonization_type": job.harmonization_type
                },
                timeout=30.0
            )
            
            if response.status_code != 202:
                # Handle error
                logger.error(f"Worker error: {response.text}")
                job.status = JobStatus.FAILED
                job.error_message = f"Worker error: {response.text}"
                db.commit()
                return
            
            # Job successfully sent to worker
            logger.info(f"Job {job_id} sent to worker")
            
    except Exception as e:
        # Handle connection error
        logger.error(f"Error sending job to worker: {str(e)}")
        job.status = JobStatus.FAILED
        job.error_message = f"Error sending job to worker: {str(e)}"
        db.commit()


@router.post(
    "/",
    response_model=JobResponse,
    status_code=202,
    responses={
        400: {"model": ErrorResponse, "description": "Bad Request"},
        413: {"model": ErrorResponse, "description": "File too large"}
    }
)
async def create_job(
    background_tasks: BackgroundTasks,
    file: UploadFile = File(...),
    harmonization_type: str = Form(None),
    db: Session = Depends(get_db)
):
    """
    Create a new data harmonization job.
    
    Args:
        background_tasks: FastAPI background tasks
        file: File to harmonize
        harmonization_type: Type of harmonization to perform (optional)
        db: Database session
        
    Returns:
        JobResponse: Created job details
    """
    try:
        # Save uploaded file
        file_path, file_type, file_size_bytes, file_size = await save_upload_file(file)
        
        # Determine harmonization type if not provided
        if not harmonization_type:
            harmonization_type = get_harmonization_type_from_filename(file.filename)
        
        # Create job in database
        job = Job(
            input_file=file_path,
            file_type=file_type,
            file_size=file_size,
            harmonization_type=harmonization_type
        )
        
        db.add(job)
        db.commit()
        db.refresh(job)
        
        # Process job in background
        background_tasks.add_task(process_job_async, job.id, db)
        
        logger.info(f"Created job {job.id} for file {file.filename}")
        return job
        
    except HTTPException as e:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        logger.error(f"Error creating job: {str(e)}")
        raise HTTPException(
            status_code=500,
            detail=f"Error creating job: {str(e)}"
        )


@router.get(
    "/{job_id}",
    response_model=JobResponse,
    responses={404: {"model": ErrorResponse, "description": "Job not found"}}
)
async def get_job(job_id: str, db: Session = Depends(get_db)):
    """
    Get job details.
    
    Args:
        job_id: ID of the job
        db: Database session
        
    Returns:
        JobResponse: Job details
    """
    job = db.query(Job).filter(Job.id == job_id).first()
    if not job:
        raise HTTPException(
            status_code=404,
            detail=f"Job {job_id} not found"
        )
    
    return job


@router.get(
    "/{job_id}/status",
    response_model=JobStatusResponse,
    responses={404: {"model": ErrorResponse, "description": "Job not found"}}
)
async def get_job_status(job_id: str, db: Session = Depends(get_db)):
    """
    Get job status.
    
    Args:
        job_id: ID of the job
        db: Database session
        
    Returns:
        JobStatusResponse: Job status details
    """
    job = db.query(Job).filter(Job.id == job_id).first()
    if not job:
        raise HTTPException(
            status_code=404,
            detail=f"Job {job_id} not found"
        )
    
    return job


@router.get(
    "/{job_id}/result",
    responses={
        404: {"model": ErrorResponse, "description": "Job or result not found"},
        400: {"model": ErrorResponse, "description": "Job not completed"}
    }
)
async def get_job_result(job_id: str, db: Session = Depends(get_db)):
    """
    Get job result file.
    
    Args:
        job_id: ID of the job
        db: Database session
        
    Returns:
        FileResponse: Result file
    """
    job = db.query(Job).filter(Job.id == job_id).first()
    if not job:
        raise HTTPException(
            status_code=404,
            detail=f"Job {job_id} not found"
        )
    
    if job.status != JobStatus.COMPLETED:
        raise HTTPException(
            status_code=400,
            detail=f"Job {job_id} is not completed (status: {job.status.value})"
        )
    
    if not job.output_file or not os.path.exists(job.output_file):
        raise HTTPException(
            status_code=404,
            detail=f"Result file for job {job_id} not found"
        )
    
    return FileResponse(
        job.output_file,
        filename=os.path.basename(job.output_file),
        media_type="text/csv" if job.file_type == "csv" else "application/json"
    )


@router.get("/", response_model=List[JobResponse])
async def list_jobs(
    skip: int = 0,
    limit: int = 100,
    status: JobStatus = None,
    db: Session = Depends(get_db)
):
    """
    List jobs with optional filtering.
    
    Args:
        skip: Number of jobs to skip
        limit: Maximum number of jobs to return
        status: Filter by job status
        db: Database session
        
    Returns:
        List[JobResponse]: List of jobs
    """
    query = db.query(Job)
    
    if status:
        query = query.filter(Job.status == status)
    
    jobs = query.order_by(Job.created_at.desc()).offset(skip).limit(limit).all()
    return jobs
