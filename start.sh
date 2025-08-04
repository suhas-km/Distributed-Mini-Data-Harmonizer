#!/bin/bash

# Distributed Mini Data Harmonizer Startup Script
# This script brings up the entire application with Docker Compose

echo "Starting Distributed Mini Data Harmonizer..."

# Create necessary directories if they don't exist
echo "Creating necessary directories..."
mkdir -p uploads results data

# Build and start all services
echo "Building and starting all services..."
docker compose up -d --build

echo ""
echo "Distributed Mini Data Harmonizer is now running!"
echo "Access the services at:"
echo "- UI: http://localhost:8082"
echo "- API: http://localhost:8080"
echo "- API Documentation: http://localhost:8080/docs"
echo ""
echo "To view logs, run: docker compose logs -f"
echo "To stop the application, run: docker compose down"
