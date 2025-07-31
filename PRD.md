# Project Requirements Document (PRD)
## Distributed Mini Data Harmonizer

### Introduction
The Distributed Mini Data Harmonizer is a distributed data processing pipeline designed to harmonize healthcare data using a Python orchestrator and Go workers. This project demonstrates modern distributed system patterns, concurrency handling, and healthcare data processing capabilities.

### Project Goals

#### Business Goals
- Create a scalable, distributed data harmonization system for healthcare data
- Demonstrate proficiency in Go concurrency and Python orchestration
- Build a portfolio project showcasing distributed system design
- Prepare for technical interviews in healthcare data engineering roles

#### Technical Goals
- Implement efficient concurrent data processing using Go goroutines and channels
- Design a robust REST API for job orchestration using Python
- Demonstrate inter-service communication patterns
- Implement basic observability and monitoring
- Handle healthcare data formats (CSV, JSON) with proper validation

### Target Audience

#### Primary Users
- Healthcare data engineers processing EHR data
- Data scientists requiring cleaned and standardized healthcare datasets
- System administrators managing data pipeline workflows

#### Secondary Users
- Developers learning distributed system patterns
- Technical interviewers evaluating system design skills

### Core Features

#### Must-Have Features
1. **File Upload & Processing**
   - Accept CSV/JSON healthcare data files
   - Queue processing jobs asynchronously
   - Return processed results with metadata

2. **Go Worker Service**
   - HTTP/gRPC endpoint for receiving jobs
   - Concurrent file processing using goroutines
   - Data harmonization operations (deduplication, standardization, validation)

3. **Python Orchestrator**
   - REST API for job submission and status tracking
   - Job queue management
   - Database integration for job metadata storage

4. **Data Harmonization**
   - Remove duplicate records
   - Standardize field names and formats
   - Basic data validation and checksums

#### Nice-to-Have Features
1. **Web UI**
   - Simple interface for file uploads
   - Job status monitoring dashboard
   - Results visualization

2. **Advanced Observability**
   - Prometheus metrics export
   - Structured logging with correlation IDs
   - Performance monitoring dashboards

3. **Extended Data Support**
   - FASTA/VCF file processing for bioinformatics
   - R worker integration for statistical processing
   - Batch processing capabilities

### Out of Scope
- Production-grade security implementation
- Complex data transformation algorithms
- Real-time streaming data processing
- Multi-tenant architecture
- Advanced machine learning features

### Technical Requirements

#### System Architecture
- **Python Backend**: Flask/FastAPI for REST API
- **Go Worker**: HTTP server with concurrent processing
- **Database**: SQLite (development) / PostgreSQL (production)
- **Queue**: In-memory (demo) / Redis/RabbitMQ (production)

#### Performance Requirements
- Process multiple files concurrently (minimum 5 simultaneous jobs)
- Handle files up to 100MB in size
- Response time under 2 seconds for job submission
- Job processing time proportional to file size

#### Data Requirements
- Support CSV and JSON input formats
- Handle mock EHR data: patients, lab results, imaging metadata
- Maintain data integrity throughout processing
- Generate processing audit logs

### Success Metrics

#### Technical Metrics
- Job processing throughput (files/minute)
- System uptime and reliability
- Error rate and recovery time
- Resource utilization efficiency

#### Learning Metrics
- Demonstration of Go concurrency patterns
- Implementation of distributed system principles
- Code quality and documentation standards
- Interview readiness and technical communication

### Timeline
- **Week 1**: Core Go worker and Python orchestrator
- **Week 2**: Integration, testing, and basic observability
- **Week 3**: Documentation, polish, and optional features

### Risk Assessment
- **Technical Risk**: Go-Python integration complexity
- **Timeline Risk**: Learning curve for new technologies
- **Scope Risk**: Feature creep beyond core requirements

### Dependencies
- Go 1.19+ development environment
- Python 3.8+ with Flask/FastAPI
- Database system (SQLite/PostgreSQL)
- Mock healthcare datasets
