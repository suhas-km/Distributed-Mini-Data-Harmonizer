# Distributed Mini Data Harmonizer

A distributed data processing pipeline that combines Python orchestration with Go workers to harmonize healthcare data. This project demonstrates modern distributed system patterns, concurrent processing, and inter-service communication in a local environment.

## 🏗️ Architecture Overview

```
┌─────────────────┐    HTTP/gRPC    ┌─────────────────┐
│  Python API     │ ──────────────► │   Go Worker     │
│  (Orchestrator) │                 │   (Processor)   │
│                 │ ◄────────────── │                 │
└─────────┬───────┘                 └─────────────────┘
          │
          ▼
┌─────────────────┐
│   Database      │
│ (Job Metadata)  │
└─────────────────┘
```

- **Python Backend**: REST API, job orchestration, and database management
- **Go Worker**: Concurrent data processing with goroutines and channels
- **Database**: Job metadata, results, and audit logs

## 🚀 Features

- **Concurrent Processing**: Go workers handle multiple files simultaneously
- **Healthcare Data Support**: Process EHR data (CSV/JSON formats)
- **Data Harmonization**: Deduplication, standardization, and validation
- **REST API**: Submit jobs, track status, retrieve results
- **Observability**: Logging, metrics, and job monitoring
- **Scalable Design**: Easy to extend with additional workers

## 📋 Prerequisites

- **Docker**: 20.10 or higher
- **Docker Compose**: 2.0 or higher
- **Git**: For version control

*Note: Direct installation of Go and Python is not required as the application runs in Docker containers.*

## 🛠️ Installation

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/Distributed-Mini-Data-Harmonizer.git
cd Distributed-Mini-Data-Harmonizer
```

### 2. One-Command Setup
The project includes a startup script that creates necessary directories, builds Docker images, and starts all services:

```bash
./start.sh
```

This script will:
- Create required directories (uploads, results, data)
- Build and start all Docker containers
- Configure networking between services
- Initialize the database

### 3. Sample Data
Sample healthcare data is available in the `sample_data/` directory, including:
- patients.csv - Patient demographic information
- vitals.csv - Patient vital signs
- medications.csv - Medication records
- lab_results.csv - Laboratory test results

## 🎯 Usage

### Access the Services

After running `./start.sh`, the following services are available:

- **Web UI**: http://localhost:8082
- **API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/docs
- **Go Worker**: http://localhost:8081
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001

### Using the Web UI

1. Navigate to http://localhost:8082
2. Upload a file for processing
3. Select the harmonization type
4. Submit the job and monitor progress

### Using the API Directly

#### Submit a Job

```bash
# Upload a file for processing
curl -X POST http://localhost:8080/api/v1/jobs/ \
  -F "file=@sample_data/patients.csv" \
  -F "harmonization_type=patients"
```

#### Check Job Status

```bash
# Get job status
curl http://localhost:8080/api/v1/jobs/{job_id}
```

#### List All Jobs

```bash
# Get all jobs
curl http://localhost:8080/api/v1/jobs/
```

#### Retrieve Results

```bash
# Download processed file
curl http://localhost:8080/api/v1/jobs/{job_id}/result -o result.csv
```

## 🧪 Testing

### Run Unit Tests
```bash
# Python tests (inside container)
docker compose exec python-api pytest tests/

# Go tests (inside container)
docker compose exec go-worker go test ./...
```

### Run Integration Tests
```bash
# End-to-end pipeline test
docker compose exec python-api python tests/integration/test_pipeline.py
```

## 📊 Monitoring

- **Logs**: View logs with `docker compose logs -f [service-name]`
- **Health Checks**: 
  - Python API: `http://localhost:8080/health`
  - Go Worker: `http://localhost:8081/health`
- **Metrics**: 
  - Prometheus: `http://localhost:9090`
  - Grafana: `http://localhost:3001` (admin/admin)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Start for Contributors
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run the test suite: `make test`
5. Submit a pull request

## 📚 Documentation

- [Project Requirements Document (PRD)](PRD.md)
- [Architecture Guide](ARCHITECTURE.md)
- [API Documentation](API.md)
- [Contributing Guidelines](CONTRIBUTING.md)

## 🏷️ Tech Stack

- **Backend**: Python with FastAPI and SQLAlchemy ORM
- **Worker**: Go with goroutines and worker pools
- **Database**: SQLite (file-based)
- **API Documentation**: OpenAPI/Swagger
- **Monitoring**: Prometheus, Grafana, structured logging
- **Containerization**: Docker, Docker Compose
- **UI**: Simple HTML/CSS/JS with Nginx

## 📈 Performance

- Processes 5+ files concurrently
- Handles files up to 100MB
- Sub-2-second job submission response time

## 🔄 Implementation Plan

1. **Phase 1**: Project structure with FastAPI + SQLAlchemy + Docker
   - Setup Python API with FastAPI framework
   - Implement SQLAlchemy ORM models
   - Create Docker configuration

2. **Phase 2**: Go worker with concurrency patterns
   - Implement HTTP server in Go
   - Create worker pool with bounded concurrency
   - Add context-based cancellation

3. **Phase 3**: Monitoring and observability
   - Add structured logging
   - Implement health check endpoints
   - Create basic metrics collection

4. **Phase 4**: Data processing
   - Implement various data harmonization processors
   - Add support for different healthcare data types
   - Create file handling utilities

5. **Phase 5**: Integration and deployment
   - Implement Docker Compose orchestration
   - Create startup script for one-command deployment
   - Add minimal UI for job submission and monitoring

## 🔧 Troubleshooting

### Common Issues

#### File Upload Issues
- **Problem**: Go worker cannot find uploaded files
- **Solution**: Ensure the `uploads` directory exists and is properly mounted in both containers. The startup script creates this directory automatically.

#### Container Communication
- **Problem**: Python API cannot connect to Go worker
- **Solution**: Use the Docker service name (`go-worker`) instead of `localhost` in the Python API configuration.

#### Job Status Updates
- **Problem**: Go worker receives 405 Method Not Allowed error when updating job status
- **Solution**: Ensure the Python API has a POST endpoint for job status updates at `/api/v1/jobs/{job_id}/status`.

#### Docker Volume Permissions
- **Problem**: Permission denied when accessing mounted volumes
- **Solution**: Ensure the host directories have appropriate permissions (e.g., `chmod -R 777 uploads results data`).

#### Port Conflicts
- **Problem**: Services fail to start due to port conflicts
- **Solution**: Check if other applications are using the same ports and modify the Docker Compose file if needed.

### Debugging Commands

```bash
# View logs for a specific service
docker compose logs -f python-api

# Check container status
docker compose ps

# Restart a specific service
docker compose restart go-worker

# Rebuild and restart all services
docker compose up -d --build
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Authors

- **Suhas KM** - *Initial work* - [Suhas KM](https://github.com/suhas-km)

## Support

If you have questions or need help:
- Open an issue on GitHub
- Check the [Architecture Guide](ARCHITECTURE.md)
- Check the [API Documentation](API.md)
- Check the [Contributing Guidelines](CONTRIBUTING.md)
- Check the [Project Requirements Document (PRD)](PRD.md)
- Check the [README](README.md)
