"""
Utility functions for file handling in the Data Harmonizer API.
"""

import os
import uuid
import shutil
from pathlib import Path
from typing import Tuple, Optional

from fastapi import UploadFile, HTTPException
from config import settings


def validate_file_extension(filename: str) -> bool:
    """
    Validate that the file has an allowed extension.
    
    Args:
        filename: Name of the file to validate
        
    Returns:
        bool: True if the file extension is allowed, False otherwise
    """
    ext = filename.split('.')[-1].lower()
    return ext in settings.ALLOWED_EXTENSIONS


def validate_file_size(file_size: int) -> bool:
    """
    Validate that the file size is within allowed limits.
    
    Args:
        file_size: Size of the file in bytes
        
    Returns:
        bool: True if the file size is allowed, False otherwise
    """
    return file_size <= settings.MAX_UPLOAD_SIZE


def get_file_type(filename: str) -> str:
    """
    Get the file type from the filename.
    
    Args:
        filename: Name of the file
        
    Returns:
        str: File extension
    """
    return filename.split('.')[-1].lower()


def format_file_size(size_bytes: int) -> str:
    """
    Format file size in a human-readable format.
    
    Args:
        size_bytes: Size in bytes
        
    Returns:
        str: Human-readable file size
    """
    for unit in ['B', 'KB', 'MB', 'GB']:
        if size_bytes < 1024 or unit == 'GB':
            return f"{size_bytes:.2f} {unit}"
        size_bytes /= 1024


async def save_upload_file(upload_file: UploadFile) -> Tuple[str, str, int, str]:
    """
    Save an uploaded file to the uploads directory.
    
    Args:
        upload_file: The uploaded file
        
    Returns:
        Tuple containing:
            - Saved file path
            - File type
            - File size in bytes
            - Human-readable file size
            
    Raises:
        HTTPException: If the file is invalid
    """
    # Validate file extension
    if not validate_file_extension(upload_file.filename):
        raise HTTPException(
            status_code=400,
            detail=f"File type not allowed. Allowed types: {', '.join(settings.ALLOWED_EXTENSIONS)}"
        )
    
    # Create a unique filename
    file_type = get_file_type(upload_file.filename)
    unique_filename = f"{uuid.uuid4()}.{file_type}"
    file_path = os.path.join(settings.UPLOAD_DIR, unique_filename)
    
    # Save the file
    with open(file_path, "wb") as buffer:
        # Get file size while saving
        size_bytes = 0
        while content := await upload_file.read(1024 * 1024):  # Read 1MB at a time
            size_bytes += len(content)
            
            # Check file size
            if not validate_file_size(size_bytes):
                # Remove partially saved file
                os.remove(file_path)
                raise HTTPException(
                    status_code=413,
                    detail=f"File too large. Maximum size: {format_file_size(settings.MAX_UPLOAD_SIZE)}"
                )
            
            buffer.write(content)
    
    # Format file size for display
    human_readable_size = format_file_size(size_bytes)
    
    return file_path, file_type, size_bytes, human_readable_size


def get_harmonization_type_from_filename(filename: str) -> str:
    """
    Determine the harmonization type based on the filename.
    
    Args:
        filename: Name of the file
        
    Returns:
        str: Harmonization type (patients, vitals, medications, lab_results)
    """
    basename = os.path.basename(filename).lower()
    
    if "patient" in basename:
        return "patients"
    elif "vital" in basename:
        return "vitals"
    elif "medication" in basename or "med" in basename:
        return "medications"
    elif "lab" in basename or "test" in basename:
        return "lab_results"
    else:
        # Default to generic harmonization
        return "generic"
