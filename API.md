# API Documentation
## Distributed Mini Data Harmonizer

### Base URL
```
Local Development: http://localhost:8080/api/v1
```

### Authentication

**Simple Development Setup**: No authentication required for local development.

**Optional Basic Auth** (if needed):
```http
Authorization: Basic dXNlcjpwYXNzd29yZA==
```

### Rate Limiting
**Local Development**: No rate limiting implemented (keep it simple).

## Endpoints

### Job Management

#### Create Job
Submit a new data harmonization job.

```http
POST /api/v1/jobs
```

**Request:**
```http
Content-Type: multipart/form-data

file: [binary file data]
operation: harmonize
options: {
  "remove_duplicates": true,
  "standardize_fields": true,
  "validate_schema": true,
  "output_format": "csv"
}
```

**Parameters:**
- `file` (required): Healthcare data file (CSV, JSON)
- `operation` (required): Processing operation type
  - `harmonize`: Full data harmonization
  - `validate`: Data validation only
  - `clean`: Data cleaning only
  - `standardize`: Field standardization only
- `options` (optional): Processing configuration object

**Response:**
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "queued",
  "message": "Job created successfully",
  "estimated_duration": "2-5 minutes",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `201 Created`: Job created successfully
- `400 Bad Request`: Invalid file or parameters
- `413 Payload Too Large`: File exceeds size limit (100MB)
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

#### List Jobs
Retrieve a list of jobs with pagination and filtering.

```http
GET /api/v1/jobs
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)
- `status` (optional): Filter by job status
- `operation` (optional): Filter by operation type
- `created_after` (optional): Filter by creation date (ISO 8601)
- `created_before` (optional): Filter by creation date (ISO 8601)

**Example:**
```http
GET /api/v1/jobs?page=1&limit=10&status=completed&operation=harmonize
```

**Response:**
```json
{
  "jobs": [
    {
      "job_id": "123e4567-e89b-12d3-a456-426614174000",
      "status": "completed",
      "operation": "harmonize",
      "file_name": "patients.csv",
      "created_at": "2024-01-15T10:30:00Z",
      "completed_at": "2024-01-15T10:33:45Z",
      "progress": 100
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

---

#### Get Job Details
Retrieve detailed information about a specific job.

```http
GET /api/v1/jobs/{job_id}
```

**Path Parameters:**
- `job_id` (required): Unique job identifier

**Response:**
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "completed",
  "operation": "harmonize",
  "file_name": "patients.csv",
  "file_size": 2048576,
  "options": {
    "remove_duplicates": true,
    "standardize_fields": true,
    "validate_schema": true
  },
  "progress": 100,
  "created_at": "2024-01-15T10:30:00Z",
  "started_at": "2024-01-15T10:30:15Z",
  "completed_at": "2024-01-15T10:33:45Z",
  "metrics": {
    "records_processed": 10000,
    "records_cleaned": 9850,
    "duplicates_removed": 150,
    "validation_errors": 25,
    "processing_time_seconds": 210
  },
  "result_file_size": 1987654,
  "download_url": "/api/v1/jobs/123e4567-e89b-12d3-a456-426614174000/result"
}
```

**Status Codes:**
- `200 OK`: Job found
- `404 Not Found`: Job not found
- `403 Forbidden`: Access denied

---

#### Get Job Status
Get the current status of a job (lightweight endpoint).

```http
GET /api/v1/jobs/{job_id}/status
```

**Response:**
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "processing",
  "progress": 75,
  "estimated_completion": "2024-01-15T10:35:00Z",
  "current_stage": "data_validation"
}
```

**Job Status Values:**
- `queued`: Job is waiting to be processed
- `processing`: Job is currently being processed
- `completed`: Job completed successfully
- `failed`: Job failed with errors
- `cancelled`: Job was cancelled by user

---

#### Cancel Job
Cancel a queued or processing job.

```http
DELETE /api/v1/jobs/{job_id}
```

**Response:**
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "cancelled",
  "message": "Job cancelled successfully"
}
```

**Status Codes:**
- `200 OK`: Job cancelled
- `400 Bad Request`: Job cannot be cancelled (already completed)
- `404 Not Found`: Job not found

---

#### Download Result
Download the processed file result.

```http
GET /api/v1/jobs/{job_id}/result
```

**Query Parameters:**
- `format` (optional): Output format (`csv`, `json`) - overrides job settings

**Response:**
- **Content-Type**: `application/octet-stream`
- **Content-Disposition**: `attachment; filename="processed_patients.csv"`
- **Body**: Binary file data

**Status Codes:**
- `200 OK`: File download
- `202 Accepted`: Job not yet completed
- `404 Not Found`: Job or result not found
- `410 Gone`: Result file expired or deleted

---

#### Get Job Logs
Retrieve processing logs for a job.

```http
GET /api/v1/jobs/{job_id}/logs
```

**Query Parameters:**
- `level` (optional): Filter by log level (`debug`, `info`, `warn`, `error`)
- `limit` (optional): Number of log entries (default: 100, max: 1000)

**Response:**
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "logs": [
    {
      "timestamp": "2024-01-15T10:30:15Z",
      "level": "info",
      "message": "Job processing started",
      "stage": "initialization"
    },
    {
      "timestamp": "2024-01-15T10:30:20Z",
      "level": "info",
      "message": "File validation completed",
      "stage": "validation",
      "details": {
        "records_found": 10000,
        "schema_valid": true
      }
    },
    {
      "timestamp": "2024-01-15T10:32:45Z",
      "level": "warn",
      "message": "25 validation errors found",
      "stage": "data_cleaning",
      "details": {
        "error_types": ["missing_field", "invalid_date"]
      }
    }
  ]
}
```

### Health and Monitoring

#### Health Check
Check the health status of the API service.

```http
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "dependencies": {
    "database": "healthy",
    "redis": "healthy",
    "go_worker": "healthy"
  },
  "uptime_seconds": 86400
}
```

**Status Values:**
- `healthy`: All systems operational
- `degraded`: Some non-critical issues
- `unhealthy`: Critical issues detected

---

#### Metrics
Get system metrics (Prometheus format).

```http
GET /api/v1/metrics
```

**Response:**
```
# HELP jobs_total Total number of jobs processed
# TYPE jobs_total counter
jobs_total{status="completed"} 1250
jobs_total{status="failed"} 45

# HELP job_processing_duration_seconds Time spent processing jobs
# TYPE job_processing_duration_seconds histogram
job_processing_duration_seconds_bucket{operation="harmonize",le="10"} 100
job_processing_duration_seconds_bucket{operation="harmonize",le="30"} 450
job_processing_duration_seconds_bucket{operation="harmonize",le="60"} 800
```

---

#### Version Information
Get API version and build information.

```http
GET /api/v1/version
```

**Response:**
```json
{
  "api_version": "1.0.0",
  "build_date": "2024-01-15T08:00:00Z",
  "git_commit": "abc123def456",
  "go_version": "1.19.5",
  "python_version": "3.9.16"
}
```

## Error Handling

### Error Response Format
All error responses follow a consistent format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid file format",
    "details": {
      "field": "file",
      "expected": "CSV or JSON",
      "received": "PDF"
    },
    "request_id": "req_123456789"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `FILE_TOO_LARGE` | 413 | File exceeds size limit |
| `UNSUPPORTED_FORMAT` | 400 | File format not supported |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `JOB_NOT_FOUND` | 404 | Job ID not found |
| `UNAUTHORIZED` | 401 | Invalid or missing authentication |
| `FORBIDDEN` | 403 | Access denied |
| `INTERNAL_ERROR` | 500 | Server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

## SDKs and Examples

### Python SDK Example
```python
import requests
import json

class HarmonizerClient:
    def __init__(self, base_url, api_key):
        self.base_url = base_url
        self.headers = {
            'Authorization': f'Bearer {api_key}',
            'User-Agent': 'HarmonizerClient/1.0'
        }
    
    def create_job(self, file_path, operation='harmonize', options=None):
        url = f'{self.base_url}/jobs'
        
        with open(file_path, 'rb') as f:
            files = {'file': f}
            data = {
                'operation': operation,
                'options': json.dumps(options or {})
            }
            
            response = requests.post(url, files=files, data=data, headers=self.headers)
            return response.json()
    
    def get_job_status(self, job_id):
        url = f'{self.base_url}/jobs/{job_id}/status'
        response = requests.get(url, headers=self.headers)
        return response.json()
    
    def download_result(self, job_id, output_path):
        url = f'{self.base_url}/jobs/{job_id}/result'
        response = requests.get(url, headers=self.headers)
        
        with open(output_path, 'wb') as f:
            f.write(response.content)

# Usage
client = HarmonizerClient('http://localhost:8080/api/v1', 'your-api-key')

# Create job
job = client.create_job('patients.csv', 'harmonize', {
    'remove_duplicates': True,
    'standardize_fields': True
})

print(f"Job created: {job['job_id']}")

# Check status
status = client.get_job_status(job['job_id'])
print(f"Status: {status['status']}")

# Download result (when completed)
if status['status'] == 'completed':
    client.download_result(job['job_id'], 'result.csv')
```

### cURL Examples

**Create a job:**
```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Authorization: Bearer your-api-key" \
  -F "file=@patients.csv" \
  -F "operation=harmonize" \
  -F 'options={"remove_duplicates":true,"standardize_fields":true}'
```

**Check job status:**
```bash
curl -H "Authorization: Bearer your-api-key" \
  http://localhost:8080/api/v1/jobs/123e4567-e89b-12d3-a456-426614174000/status
```

**Download result:**
```bash
curl -H "Authorization: Bearer your-api-key" \
  http://localhost:8080/api/v1/jobs/123e4567-e89b-12d3-a456-426614174000/result \
  -o processed_data.csv
```

## Simple Local Development

### File Limits
- **File Upload**: 100MB max file size (configurable)
- **Concurrent Jobs**: 3 simultaneous jobs (to prevent resource exhaustion)
- **Result Storage**: Files stored locally in `./results/` directory

### Basic Configuration
Edit `config.py` to adjust:
```python
# config.py
MAX_FILE_SIZE = 100 * 1024 * 1024  # 100MB
MAX_CONCURRENT_JOBS = 3
RESULTS_DIR = "./results"
UPLOADS_DIR = "./uploads"
DATABASE_PATH = "./data/harmonizer.db"
```
