# Task Management API

A secure REST API for task management built with Go, implementing OAuth 2.0 Client Credentials Flow with JSON Web Signatures (JWS) for machine-to-machine authentication.

## Features

- **CRUD Operations**: Create, read, update, delete tasks
- **OAuth 2.0 Authentication**: Client Credentials Flow with JWT/JWS
- **Certificate-based Security**: No pre-shared secrets, using PKI
- **Clean Architecture**: Domain/Service/Transport layer separation
- **Type-safe Database**: SQLite with SQLC for compile-time SQL validation
- **API Documentation**: Interactive Swagger UI
- **Graceful Shutdown**: Handles SIGTERM/SIGINT properly
- **Docker Ready**: Self-contained containerized deployment

## Quick Start

### Local Development
```bash
make setup    # Generate code, docs, run migrations
make run      # Start server on :8080
```

### Docker
```bash
make docker-build && make docker-run
```

### Authentication
```bash
# 1. Generate client assertion JWT
cd util && go run generate_jwt.go

# 2. Get access token
curl -X POST http://localhost:8080/token \
  -d "grant_type=client_credentials" \
  -d "client_assertion=PASTE_JWT_HERE" \
  -d "client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

# 3. Use access token
curl -H "Authorization: Bearer ACCESS_TOKEN" http://localhost:8080/tasks
```

## API Endpoints

- `POST /token` - Get JWT access token
- `POST /tasks` - Create task
- `GET /tasks` - List all tasks
- `GET /tasks/{id}` - Get task by ID
- `PATCH /tasks/{id}` - Update task (partial)
- `DELETE /tasks/{id}` - Delete task

📖 **Full API docs**: http://localhost:8080/swagger/index.html

## Task Model

```json
{
  "id": "uuid",
  "title": "string (required)",
  "description": "string",
  "status": "to_do | in_progress | done",
  "priority": "low | medium | high",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Tech Stack

- **Go 1.24** with Chi router
- **SQLite** + SQLC for type-safe queries
- **OAuth 2.0** Client Credentials Flow
- **JWT/JWS** with RSA signatures
- **Docker** with multi-stage builds
- **Swagger/OpenAPI** documentation

## Configuration

Environment variables (see `.env.example`):
- `PORT=8080`
- `JWT_PRIVATE_KEY_FILE=secret/server.key`
- `JWT_ISSUER=ishare-task-api`
- `DB_PATH=tasks.db`

## Security Setup

⚠️ **IMPORTANT**: The certificates in `secret/` are for **development/testing only**.

### For Production:

**1. Generate your own certificates:**
```bash
# Generate private key
openssl genpkey -algorithm RSA -out secret/server.key -pkcs8 -pkeyopt rsa_keygen_bits:2048

# Generate certificate signing request
openssl req -new -key secret/server.key -out secret/server.csr \
  -subj "/CN=your-api-domain.com/O=YourOrg/C=US"

# Generate self-signed certificate (or use CA-signed)
openssl x509 -req -days 365 -in secret/server.csr -signkey secret/server.key -out secret/server.crt

# Generate client certificates (for testing)
openssl genpkey -algorithm RSA -out secret/client.key -pkcs8 -pkeyopt rsa_keygen_bits:2048
openssl req -new -key secret/client.key -out secret/client.csr \
  -subj "/CN=test-client/O=YourOrg/C=US"
openssl x509 -req -days 365 -in secret/client.csr -signkey secret/client.key -out secret/client.crt
```

**2. Update environment:**
```bash
export JWT_PRIVATE_KEY_FILE=secret/server.key
export JWT_ISSUER=https://your-api-domain.com
```

**3. Never commit real certificates to version control!**

## Development

```bash
make test              # Run all tests
make test-unit         # Unit tests only
make test-integration  # Integration tests only
make test-coverage     # Generate coverage report
```