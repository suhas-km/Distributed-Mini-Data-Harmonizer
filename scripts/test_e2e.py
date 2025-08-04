#!/usr/bin/env python3
"""
End-to-End Test Script for Distributed Mini Data Harmonizer

This script tests the complete pipeline:
1. Starts all services via Docker Compose
2. Waits for services to be healthy
3. Uploads test data files
4. Monitors job processing
5. Validates harmonized output
6. Tests the UI endpoints
7. Checks monitoring endpoints
"""

import os
import sys
import time
import json
import requests
import subprocess
import pandas as pd
from pathlib import Path

# Configuration
API_BASE_URL = "http://localhost:8080"
WORKER_BASE_URL = "http://localhost:8081"
UI_BASE_URL = "http://localhost:3000"
PROMETHEUS_URL = "http://localhost:9090"
GRAFANA_URL = "http://localhost:3001"

SAMPLE_DATA_DIR = Path("sample_data")
RESULTS_DIR = Path("results")
TIMEOUT_SECONDS = 300  # 5 minutes


class E2ETestRunner:
    def __init__(self):
        self.test_results = []
        self.failed_tests = 0
        
    def log(self, message, level="INFO"):
        timestamp = time.strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] [{level}] {message}")
        
    def run_test(self, test_name, test_func):
        """Run a test and record results"""
        self.log(f"Running test: {test_name}")
        try:
            test_func()
            self.log(f"✅ PASSED: {test_name}")
            self.test_results.append({"test": test_name, "status": "PASSED"})
        except Exception as e:
            self.log(f"❌ FAILED: {test_name} - {str(e)}", "ERROR")
            self.test_results.append({"test": test_name, "status": "FAILED", "error": str(e)})
            self.failed_tests += 1
            
    def wait_for_service(self, url, service_name, timeout=60):
        """Wait for a service to become available"""
        self.log(f"Waiting for {service_name} at {url}")
        start_time = time.time()
        
        while time.time() - start_time < timeout:
            try:
                response = requests.get(f"{url}/health", timeout=5)
                if response.status_code == 200:
                    self.log(f"✅ {service_name} is ready")
                    return True
            except requests.exceptions.RequestException:
                pass
            time.sleep(2)
            
        raise Exception(f"{service_name} failed to start within {timeout} seconds")
        
    def test_docker_compose_up(self):
        """Start all services with Docker Compose"""
        self.log("Starting Docker Compose services...")
        
        # Stop any existing services
        subprocess.run(["docker-compose", "down"], capture_output=True)
        
        # Start services
        result = subprocess.run(
            ["docker-compose", "up", "-d", "--build"],
            capture_output=True,
            text=True
        )
        
        if result.returncode != 0:
            raise Exception(f"Docker Compose failed: {result.stderr}")
            
        # Wait for services to be ready
        self.wait_for_service(API_BASE_URL, "Python API")
        self.wait_for_service(WORKER_BASE_URL, "Go Worker")
        
        # Give UI a moment to start
        time.sleep(10)
        
    def test_api_health(self):
        """Test API health endpoints"""
        response = requests.get(f"{API_BASE_URL}/health")
        assert response.status_code == 200
        
        health_data = response.json()
        assert health_data["status"] == "ok"
        
    def test_worker_health(self):
        """Test Go Worker health endpoint"""
        response = requests.get(f"{WORKER_BASE_URL}/health")
        assert response.status_code == 200
        
        health_data = response.json()
        assert health_data["status"] == "ok"
        
    def test_ui_accessibility(self):
        """Test that UI is accessible"""
        response = requests.get(UI_BASE_URL)
        assert response.status_code == 200
        assert "Data Harmonizer" in response.text
        
    def test_file_upload_and_processing(self):
        """Test file upload and job processing"""
        # Use the generated sample data
        test_file = SAMPLE_DATA_DIR / "patients.csv"
        if not test_file.exists():
            raise Exception(f"Test file not found: {test_file}")
            
        # Upload file
        with open(test_file, 'rb') as f:
            files = {'file': f}
            data = {'harmonization_type': 'patients'}
            
            response = requests.post(
                f"{API_BASE_URL}/api/v1/jobs/",
                files=files,
                data=data
            )
            
        assert response.status_code == 201
        job_data = response.json()
        job_id = job_data["id"]
        
        self.log(f"Created job: {job_id}")
        
        # Monitor job progress
        start_time = time.time()
        while time.time() - start_time < TIMEOUT_SECONDS:
            response = requests.get(f"{API_BASE_URL}/api/v1/jobs/{job_id}")
            assert response.status_code == 200
            
            job_status = response.json()
            self.log(f"Job {job_id} status: {job_status['status']}")
            
            if job_status["status"] == "completed":
                assert job_status["output_file"] is not None
                self.log(f"✅ Job completed successfully")
                return job_id
            elif job_status["status"] == "failed":
                raise Exception(f"Job failed: {job_status.get('error_message', 'Unknown error')}")
                
            time.sleep(5)
            
        raise Exception(f"Job {job_id} did not complete within timeout")
        
    def test_result_download(self):
        """Test downloading harmonized results"""
        # First create a job
        job_id = self.test_file_upload_and_processing()
        
        # Download result
        response = requests.get(f"{API_BASE_URL}/api/v1/jobs/{job_id}/result")
        assert response.status_code == 200
        assert response.headers.get('content-type') == 'text/csv'
        
        # Validate CSV content
        csv_content = response.text
        lines = csv_content.strip().split('\n')
        assert len(lines) > 1  # Should have header + data
        
        self.log(f"✅ Successfully downloaded result with {len(lines)} lines")
        
    def test_multiple_data_types(self):
        """Test processing different data types"""
        test_files = [
            ("patients.csv", "patients"),
            ("vitals.csv", "vitals"),
            ("medications.csv", "medications"),
            ("lab_results.csv", "lab_results")
        ]
        
        job_ids = []
        
        for filename, harmonization_type in test_files:
            test_file = SAMPLE_DATA_DIR / filename
            if not test_file.exists():
                self.log(f"Skipping {filename} - file not found")
                continue
                
            # Upload file
            with open(test_file, 'rb') as f:
                files = {'file': f}
                data = {'harmonization_type': harmonization_type}
                
                response = requests.post(
                    f"{API_BASE_URL}/api/v1/jobs/",
                    files=files,
                    data=data
                )
                
            assert response.status_code == 201
            job_data = response.json()
            job_ids.append((job_data["id"], harmonization_type))
            
        # Wait for all jobs to complete
        for job_id, data_type in job_ids:
            start_time = time.time()
            while time.time() - start_time < TIMEOUT_SECONDS:
                response = requests.get(f"{API_BASE_URL}/api/v1/jobs/{job_id}")
                job_status = response.json()
                
                if job_status["status"] in ["completed", "failed"]:
                    break
                    
                time.sleep(5)
                
            assert job_status["status"] == "completed", f"Job {job_id} ({data_type}) failed"
            
        self.log(f"✅ Successfully processed {len(job_ids)} different data types")
        
    def test_monitoring_endpoints(self):
        """Test monitoring endpoints"""
        # Test Prometheus
        try:
            response = requests.get(f"{PROMETHEUS_URL}/api/v1/targets")
            assert response.status_code == 200
            self.log("✅ Prometheus is accessible")
        except requests.exceptions.RequestException:
            self.log("⚠️  Prometheus not accessible (may not be fully started yet)")
            
        # Test Grafana
        try:
            response = requests.get(GRAFANA_URL)
            assert response.status_code == 200
            self.log("✅ Grafana is accessible")
        except requests.exceptions.RequestException:
            self.log("⚠️  Grafana not accessible (may not be fully started yet)")
            
    def test_api_list_jobs(self):
        """Test listing jobs via API"""
        response = requests.get(f"{API_BASE_URL}/api/v1/jobs/")
        assert response.status_code == 200
        
        jobs = response.json()
        assert isinstance(jobs, list)
        self.log(f"✅ Found {len(jobs)} jobs in the system")
        
    def test_concurrent_processing(self):
        """Test concurrent job processing"""
        test_file = SAMPLE_DATA_DIR / "patients.csv"
        if not test_file.exists():
            raise Exception(f"Test file not found: {test_file}")
            
        # Submit multiple jobs simultaneously
        job_ids = []
        for i in range(3):
            with open(test_file, 'rb') as f:
                files = {'file': f}
                data = {'harmonization_type': 'patients'}
                
                response = requests.post(
                    f"{API_BASE_URL}/api/v1/jobs/",
                    files=files,
                    data=data
                )
                
            assert response.status_code == 201
            job_data = response.json()
            job_ids.append(job_data["id"])
            
        # Wait for all jobs to complete
        completed_jobs = 0
        start_time = time.time()
        
        while completed_jobs < len(job_ids) and time.time() - start_time < TIMEOUT_SECONDS:
            completed_jobs = 0
            for job_id in job_ids:
                response = requests.get(f"{API_BASE_URL}/api/v1/jobs/{job_id}")
                job_status = response.json()
                
                if job_status["status"] == "completed":
                    completed_jobs += 1
                elif job_status["status"] == "failed":
                    raise Exception(f"Job {job_id} failed")
                    
            time.sleep(2)
            
        assert completed_jobs == len(job_ids), f"Only {completed_jobs}/{len(job_ids)} jobs completed"
        self.log(f"✅ Successfully processed {len(job_ids)} concurrent jobs")
        
    def cleanup(self):
        """Clean up test environment"""
        self.log("Cleaning up...")
        subprocess.run(["docker-compose", "down"], capture_output=True)
        
    def run_all_tests(self):
        """Run all tests"""
        self.log("🚀 Starting End-to-End Tests")
        
        try:
            # Infrastructure tests
            self.run_test("Docker Compose Up", self.test_docker_compose_up)
            self.run_test("API Health Check", self.test_api_health)
            self.run_test("Worker Health Check", self.test_worker_health)
            self.run_test("UI Accessibility", self.test_ui_accessibility)
            
            # Functional tests
            self.run_test("File Upload and Processing", self.test_file_upload_and_processing)
            self.run_test("Result Download", self.test_result_download)
            self.run_test("Multiple Data Types", self.test_multiple_data_types)
            self.run_test("API List Jobs", self.test_api_list_jobs)
            self.run_test("Concurrent Processing", self.test_concurrent_processing)
            
            # Monitoring tests
            self.run_test("Monitoring Endpoints", self.test_monitoring_endpoints)
            
        except KeyboardInterrupt:
            self.log("Tests interrupted by user", "WARNING")
        finally:
            # Always cleanup
            if "--no-cleanup" not in sys.argv:
                self.cleanup()
                
        # Print summary
        self.print_summary()
        
    def print_summary(self):
        """Print test summary"""
        total_tests = len(self.test_results)
        passed_tests = total_tests - self.failed_tests
        
        self.log("=" * 60)
        self.log("TEST SUMMARY")
        self.log("=" * 60)
        self.log(f"Total Tests: {total_tests}")
        self.log(f"Passed: {passed_tests}")
        self.log(f"Failed: {self.failed_tests}")
        
        if self.failed_tests > 0:
            self.log("\nFailed Tests:")
            for result in self.test_results:
                if result["status"] == "FAILED":
                    self.log(f"  ❌ {result['test']}: {result.get('error', 'Unknown error')}")
                    
        success_rate = (passed_tests / total_tests) * 100 if total_tests > 0 else 0
        self.log(f"\nSuccess Rate: {success_rate:.1f}%")
        
        if self.failed_tests == 0:
            self.log("🎉 All tests passed!")
            sys.exit(0)
        else:
            self.log("💥 Some tests failed!")
            sys.exit(1)


if __name__ == "__main__":
    # Change to project root directory
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    os.chdir(project_root)
    
    # Run tests
    runner = E2ETestRunner()
    runner.run_all_tests()
