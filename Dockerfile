# Stage 1: Build stage
FROM arm32v7/golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /run_app

# Copy dependency files first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code and static files
COPY . .

# Build the Go binary
RUN go build -o /test-website

# Stage 2: Runtime stage
FROM arm32v7/alpine:latest

# Set the working directory
WORKDIR /run_app

# Copy the compiled binary
COPY --from=builder /test-website .

# Copy static files
COPY template/ template/
COPY db/ db/
COPY image/ image/

COPY go.mod go.sum ./

# Expose the desired port
EXPOSE 5000

# Set the entry point
CMD ["/run_app/test-website"]
