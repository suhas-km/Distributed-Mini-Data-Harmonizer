FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create necessary directories
RUN mkdir -p uploads results data

# Expose port
EXPOSE 8080

# Create startup script
RUN echo '#!/bin/bash\npython scripts/init_db.py\nuvicorn app:app --host 0.0.0.0 --port 8080' > /app/start.sh && chmod +x /app/start.sh

# Run the application
CMD ["/app/start.sh"]
