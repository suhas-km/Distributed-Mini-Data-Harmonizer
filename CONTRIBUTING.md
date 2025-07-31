# Contributing to Distributed Mini Data Harmonizer

Thank you for your interest in contributing to the Distributed Mini Data Harmonizer! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- Python 3.8+
- Git
- Basic understanding of distributed systems
- Familiarity with healthcare data formats (helpful but not required)

### Development Setup
1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/Distributed-Mini-Data-Harmonizer.git`
3. Follow the installation instructions in [README.md](README.md)
4. Create a new branch: `git checkout -b feature/your-feature-name`

## 📋 How to Contribute

### Reporting Bugs
Before creating a bug report, please check existing issues to avoid duplicates.

**Bug Report Template:**
```markdown
**Bug Description**
A clear description of what the bug is.

**Steps to Reproduce**
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- OS: [e.g., macOS 12.0]
- Go version: [e.g., 1.19]
- Python version: [e.g., 3.9]
- Database: [e.g., SQLite/PostgreSQL]

**Additional Context**
Add any other context about the problem here.
```

### Suggesting Features
We welcome feature suggestions! Please use the following template:

**Feature Request Template:**
```markdown
**Feature Description**
A clear description of the feature you'd like to see.

**Problem Statement**
What problem does this feature solve?

**Proposed Solution**
How would you like this feature to work?

**Alternatives Considered**
Any alternative solutions you've considered.

**Additional Context**
Any other context or screenshots about the feature request.
```

## 💻 Development Guidelines

### Code Style

#### Python Code Style
- Follow [PEP 8](https://www.python.org/dev/peps/pep-0008/)
- Use [Black](https://black.readthedocs.io/) for code formatting
- Use [isort](https://pycqa.github.io/isort/) for import sorting
- Maximum line length: 88 characters
- Use type hints where appropriate

**Example:**
```python
from typing import List, Optional

def process_healthcare_data(
    file_path: str, 
    operations: List[str]
) -> Optional[dict]:
    """Process healthcare data with specified operations.
    
    Args:
        file_path: Path to the input file
        operations: List of operations to perform
        
    Returns:
        Processing results or None if failed
    """
    # Implementation here
    pass
```

#### Go Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golint` for linting
- Use `go vet` for static analysis
- Follow Go naming conventions

**Example:**
```go
// ProcessRequest handles data processing requests
func ProcessRequest(ctx context.Context, req *ProcessingRequest) (*ProcessingResponse, error) {
    if req == nil {
        return nil, errors.New("request cannot be nil")
    }
    
    // Implementation here
    return &ProcessingResponse{}, nil
}
```

### Testing Requirements

#### Python Tests
- Use `pytest` for testing
- Aim for >80% code coverage
- Write both unit and integration tests
- Use descriptive test names

```python
def test_data_harmonizer_removes_duplicates():
    """Test that the harmonizer correctly removes duplicate records."""
    # Test implementation
    pass

def test_api_returns_job_status():
    """Test that the API correctly returns job status."""
    # Test implementation
    pass
```

#### Go Tests
- Use Go's built-in testing package
- Follow table-driven test patterns
- Test concurrent operations thoroughly

```go
func TestDataProcessor_ProcessFile(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid CSV file",
            input:    "test_data.csv",
            expected: "processed_data.csv",
            wantErr:  false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Documentation Requirements
- Update relevant documentation for any changes
- Include docstrings/comments for public functions
- Update API documentation if endpoints change
- Add examples for new features

## 🔄 Pull Request Process

### Before Submitting
1. **Run Tests**: Ensure all tests pass
   ```bash
   # Python tests
   pytest tests/
   
   # Go tests
   cd go-worker && go test ./...
   ```

2. **Code Quality Checks**:
   ```bash
   # Python
   black --check .
   isort --check-only .
   flake8 .
   
   # Go
   gofmt -d .
   golint ./...
   go vet ./...
   ```

3. **Update Documentation**: Update relevant docs and README if needed

### Pull Request Template
```markdown
## Description
Brief description of changes made.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added and passing
- [ ] No breaking changes (or breaking changes documented)
```

### Review Process
1. **Automated Checks**: All CI checks must pass
2. **Code Review**: At least one maintainer review required
3. **Testing**: Reviewer will test functionality
4. **Merge**: Squash and merge after approval

## 🏗️ Architecture Considerations

### Adding New Features
- Consider impact on both Python and Go components
- Maintain backward compatibility when possible
- Follow existing patterns and conventions
- Consider performance implications
- Add appropriate logging and error handling

### Database Changes
- Create migration scripts for schema changes
- Test migrations on sample data
- Consider impact on existing data
- Document any breaking changes

### API Changes
- Follow RESTful principles
- Maintain API versioning
- Update API documentation
- Consider backward compatibility

## 🐛 Debugging Guidelines

### Common Issues
1. **Go-Python Communication**: Check HTTP/gRPC endpoints
2. **Concurrency Issues**: Use Go race detector: `go test -race`
3. **Database Locks**: Check for proper connection handling
4. **File Processing**: Verify file permissions and formats

### Debugging Tools
- **Python**: Use `pdb` or IDE debugger
- **Go**: Use `delve` debugger or print statements
- **Logs**: Check application logs in `logs/` directory
- **Monitoring**: Use health check endpoints

## 📞 Getting Help

### Communication Channels
- **Issues**: GitHub Issues for bugs and feature requests
- **Discussions**: GitHub Discussions for questions
- **Documentation**: Check existing docs first

### Maintainer Response Time
- Bug reports: 2-3 business days
- Feature requests: 1 week
- Pull requests: 3-5 business days

## 🎯 Project Priorities

### Current Focus Areas
1. Core functionality stability
2. Performance optimization
3. Documentation improvements
4. Test coverage increase

### Future Roadmap
- Advanced data validation
- Real-time processing capabilities
- Enhanced monitoring and observability
- Multi-tenant support

## 📜 Code of Conduct

### Our Standards
- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain professional communication

### Enforcement
Violations of the code of conduct should be reported to project maintainers.

Thank you for contributing to the Distributed Mini Data Harmonizer! 🚀
