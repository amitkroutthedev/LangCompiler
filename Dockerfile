# Stage 1: Build the Go application
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create the final image
FROM debian:bullseye-slim

# Install necessary compilers and runtimes
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    g++ \
    nodejs \
    npm \
    default-jdk \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -s /bin/bash coderunner

# Create necessary directories
WORKDIR /app

# Copy the built executable from builder stage
COPY --from=builder /app/main .

# Create a directory for temporary files with appropriate permissions
RUN mkdir -p /tmp/executions && \
    chown -R coderunner:coderunner /tmp/executions && \
    chmod 755 /tmp/executions

# Switch to non-root user
USER coderunner

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]