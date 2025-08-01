# Distributed Mini Data Harmonizer

A distributed data processing pipeline that combines Python orchestration with Go workers to harmonize healthcare data. This project demonstrates modern distributed system patterns, concurrent processing, and inter-service communication.

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

- **Go**: 1.19 or higher
- **Python**: 3.8 or higher
- **Database**: SQLite (development) or PostgreSQL (production)
- **Git**: For version control

## 🛠️ Installation

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/Distributed-Mini-Data-Harmonizer.git
cd Distributed-Mini-Data-Harmonizer
```

### 2. Set Up Python Environment
```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

### 3. Set Up Go Environment
```bash
cd go-worker
go mod download
go build -o worker ./cmd/worker
```

### 4. Initialize Database
```bash
# Run database migrations
python scripts/init_db.py
```

### 5. Prepare Sample Data
```bash
# Generate mock healthcare data
python scripts/generate_sample_data.py
```

## 🎯 Usage

### Start the Services

1. **Start the Go Worker**:
```bash
cd go-worker
./worker --port=8081
```

2. **Start the Python API**:
```bash
python app.py
```

### Submit a Job

```bash
# Upload a file for processing
curl -X POST http://localhost:8080/api/jobs \
  -F "file=@sample_data/patients.csv" \
  -F "operation=harmonize"
```

### Check Job Status

```bash
# Get job status
curl http://localhost:8080/api/jobs/{job_id}/status
```

### Retrieve Results

```bash
# Download processed file
curl http://localhost:8080/api/jobs/{job_id}/result -o result.csv
```

## 🧪 Testing

### Run Unit Tests
```bash
# Python tests
pytest tests/

# Go tests
cd go-worker
go test ./...
```

### Run Integration Tests
```bash
# End-to-end pipeline test
python tests/integration/test_pipeline.py
```

## 📊 Monitoring

- **Logs**: Check `logs/` directory for application logs
- **Metrics**: Access metrics at `http://localhost:8080/metrics`
- **Health Check**: `http://localhost:8080/health`

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

- **Backend**: Python (Flask/FastAPI)
- **Worker**: Go with goroutines and channels
- **Database**: SQLite/PostgreSQL
- **Queue**: Redis (optional)
- **Monitoring**: Prometheus (optional)

## 📈 Performance

- Processes 5+ files concurrently
- Handles files up to 100MB
- Sub-2-second job submission response time
- Horizontal scaling ready

## 🔒 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Suhas KM** - *Initial work* - [Suhas KM](https://github.com/suhas-km)

## 🙏 Acknowledgments

- Healthcare data processing patterns
- Go concurrency best practices
- Distributed system design principles

## 📞 Support

If you have questions or need help:
- Open an issue on GitHub
- Check the [Architecture Guide](ARCHITECTURE.md)
- Check the [API Documentation](API.md)
- Check the [Contributing Guidelines](CONTRIBUTING.md)
- Check the [Project Requirements Document (PRD)](PRD.md)
- Check the [README](README.md)
