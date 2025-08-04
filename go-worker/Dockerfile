FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# Use a minimal alpine image for the final stage
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/worker .

# Create directories for input and output files
RUN mkdir -p /app/uploads /app/results

# Set environment variables
ENV PORT=8081
ENV HOST=0.0.0.0
ENV WORKER_COUNT=3
ENV QUEUE_SIZE=100
ENV INPUT_DIR=/app/uploads
ENV OUTPUT_DIR=/app/results
ENV PYTHON_API_URL=http://python-api:8080/api/v1

# Expose port
EXPOSE 8081

# Run the application
CMD ["./worker"]
