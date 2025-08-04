"""
Tests for job endpoints.
"""

import os
import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from app import app
from database import Base, get_db
from models.job import Job, JobStatus

# Create test database
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Test client
client = TestClient(app)


@pytest.fixture(scope="function")
def test_db():
    """Create test database and tables."""
    Base.metadata.create_all(bind=engine)
    yield
    Base.metadata.drop_all(bind=engine)


@pytest.fixture(scope="function")
def db_session(test_db):
    """Get test database session."""
    connection = engine.connect()
    transaction = connection.begin()
    session = TestingSessionLocal(bind=connection)
    
    # Override the get_db dependency
    def override_get_db():
        try:
            yield session
        finally:
            session.close()
    
    app.dependency_overrides[get_db] = override_get_db
    
    yield session
    
    transaction.rollback()
    connection.close()
    app.dependency_overrides.clear()


def test_health_check():
    """Test health check endpoint."""
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"


def test_list_jobs_empty(db_session):
    """Test listing jobs when no jobs exist."""
    response = client.get("/api/v1/jobs/")
    assert response.status_code == 200
    assert response.json() == []


def test_get_job_not_found(db_session):
    """Test getting a job that doesn't exist."""
    response = client.get("/api/v1/jobs/nonexistent-id")
    assert response.status_code == 404
    assert "not found" in response.json()["detail"]


def test_list_jobs(db_session):
    """Test listing jobs."""
    # Create test jobs
    job1 = Job(
        id="test-job-1",
        status=JobStatus.QUEUED,
        input_file="/path/to/input1.csv",
        file_type="csv",
        file_size="10.5 KB",
        harmonization_type="patients"
    )
    job2 = Job(
        id="test-job-2",
        status=JobStatus.COMPLETED,
        input_file="/path/to/input2.csv",
        output_file="/path/to/output2.csv",
        file_type="csv",
        file_size="15.2 KB",
        harmonization_type="vitals"
    )
    
    db_session.add(job1)
    db_session.add(job2)
    db_session.commit()
    
    # Test listing all jobs
    response = client.get("/api/v1/jobs/")
    assert response.status_code == 200
    jobs = response.json()
    assert len(jobs) == 2
    
    # Test filtering by status
    response = client.get("/api/v1/jobs/?status=completed")
    assert response.status_code == 200
    jobs = response.json()
    assert len(jobs) == 1
    assert jobs[0]["id"] == "test-job-2"


def test_get_job(db_session):
    """Test getting a job."""
    # Create test job
    job = Job(
        id="test-job",
        status=JobStatus.QUEUED,
        input_file="/path/to/input.csv",
        file_type="csv",
        file_size="10.5 KB",
        harmonization_type="patients"
    )
    
    db_session.add(job)
    db_session.commit()
    
    # Test getting the job
    response = client.get("/api/v1/jobs/test-job")
    assert response.status_code == 200
    job_data = response.json()
    assert job_data["id"] == "test-job"
    assert job_data["status"] == "queued"
