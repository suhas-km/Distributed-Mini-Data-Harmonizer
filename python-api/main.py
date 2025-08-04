"""
Entry point for the Data Harmonizer API.
"""

import uvicorn
import os
from app import app

if __name__ == "__main__":
    # Create required directories
    os.makedirs("uploads", exist_ok=True)
    os.makedirs("results", exist_ok=True)
    os.makedirs("data", exist_ok=True)
    
    # Run the application
    uvicorn.run(
        "app:app",
        host="0.0.0.0",
        port=8080,
        reload=True
    )
