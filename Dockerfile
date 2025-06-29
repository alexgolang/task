# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (CGO enabled for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o task-api ./cmd/main.go

# Final stage
FROM alpine:latest

# Add necessary packages
RUN apk update && \
    apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    sqlite

RUN adduser -D appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/task-api .

# Copy certificates if needed (for production, mount as volume instead)
COPY --from=builder /app/secret ./secret

# Copy migrations and other runtime files
COPY --from=builder /app/internal ./internal

# Create directory for database and set ownership
RUN mkdir -p /app/data && chown -R appuser:appuser /app

USER appuser

# Expose port
EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV PORT=8080
ENV DB_PATH=/app/tasks.db
ENV JWT_PRIVATE_KEY_FILE=/app/secret/server.key
ENV JWT_ISSUER=ishare-task-api
ENV JWT_TOKEN_EXPIRY=3600s

# Run the application
CMD ["./task-api"]