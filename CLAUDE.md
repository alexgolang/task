# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

```bash
# Build the application
make build

# Run the server locally
make run
# or
go run cmd/main.go

# Run tests
make test                # All tests
make test-unit           # Unit tests only (excludes integration tests)
make test-integration    # Integration tests only
make test-verbose        # All tests with verbose output
make test-coverage       # Generate test coverage report

# Database operations
make migrate              # Run database migrations
make migrate-status       # Check migration status
make sqlc-generate       # Generate type-safe SQL code from queries

# Generate Swagger documentation
make swagger-generate

# Complete setup (SQLC + Swagger + migrations)
make setup

# Docker commands
make docker-build        # Build Docker image
make docker-run          # Run container with volume mount
make docker-stop         # Stop and remove container
make docker-dev          # Build and run (shortcut)
```

## Architecture Overview

This is a Go REST API for task management built with **Clean Architecture** principles:

### Layer Structure
- **Domain Layer** (`internal/app/domain/`): Core business entities (Task, auth types)
- **Service Layer** (`internal/app/service/`): Business logic implementation
- **Transport Layer** (`internal/app/transport/httpserver/`): HTTP handlers and middleware
- **Infrastructure Layer** (`internal/app/db/`, `internal/app/auth/`): Database and external services

### Key Components

**Database**: SQLite with type-safe queries via SQLC
- Schema: `internal/app/db/sqlite/migrations/`
- Queries: `internal/app/db/sqlite/queries/`
- Generated code: `internal/app/db/sqlite/sqlc/`

**Authentication**: OAuth 2.0 Client Credentials Flow with JWT
- JWT service: `internal/app/auth/jwt_service.go`
- Auth middleware: `internal/app/transport/httpserver/middleware/auth.go`
- Token endpoint: `POST /token`

**HTTP Server**: Chi router with middleware
- All task endpoints require authentication via Bearer token
- Swagger documentation at `/swagger/index.html`
- Graceful shutdown with 30s timeout for active connections

### Task Model
- Fields: `id` (UUID), `title`, `description`, `status`, `priority`, `created_at`, `updated_at`
- Status: `to_do`, `in_progress`, `done`
- Priority: `low`, `medium`, `high`

## Configuration

Environment variables (with defaults):
- `PORT` (8080)
- `DB_PATH` (tasks.db)
- `JWT_PRIVATE_KEY_FILE` (path to private key file, recommended)
- `JWT_PRIVATE_KEY` (inline private key, fallback)
- `JWT_ISSUER` (ishare-task-api)  
- `JWT_TOKEN_EXPIRY` (3600s)

**Local Development Setup:**
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your configuration
# The app will automatically load .env if present
```

**⚠️ Security Note:**
- Certificates in `secret/` are for **development/testing only**
- **Never use these certificates in production**
- Generate new certificates for any real deployment
- See README.md for certificate generation instructions

## Docker Deployment

**Build and Run:**
```bash
make docker-build    # Build image: ishare-task-api
make docker-run      # Run container (database created internally)
```

**Environment Override:**
```bash
docker run -d --name ishare-task-api \
  -p 8080:8080 \
  -e JWT_PRIVATE_KEY="your-private-key" \
  -e JWT_ISSUER="your-issuer" \
  ishare-task-api
```

**Production Docker:**
```bash
# Mount your own certificates
docker run -d --name ishare-task-api \
  -p 8080:8080 \
  -v /path/to/your/certs:/app/secret \
  -e JWT_ISSUER="https://your-domain.com" \
  ishare-task-api
```

**Key Docker Features:**
- Multi-stage build for smaller final image
- CGO enabled for SQLite support
- Self-contained: database and migrations handled internally
- Certificates included for JWT authentication
- Runs as non-root user for security
- No external volumes required
- Graceful shutdown handling (SIGTERM/SIGINT)

## Development Workflow

1. Make changes to SQL queries in `internal/app/db/sqlite/queries/`
2. Run `make sqlc-generate` to regenerate type-safe Go code
3. Update API documentation comments and run `make swagger-generate`
4. Run `make migrate` if database schema changes
5. Test with `make test` (or `make test-unit` for faster feedback)
6. Build with `make build`

## Testing Strategy

**Unit Tests**: Mock dependencies using GoMock
- Located in `*_test.go` files
- Test business logic in isolation
- Fast execution, no external dependencies

**Integration Tests**: Real database interactions
- Located in `*_integration_test.go` files  
- Test service + database layer together
- Use in-memory SQLite for speed
- Include full CRUD workflows

**Test Coverage**:
- Run `make test-coverage` to generate HTML coverage report
- Focus on business logic and error handling
- Mock external dependencies (database, HTTP calls)

## API Endpoints

- `POST /token` - Get JWT token (no auth required)
- `POST /tasks` - Create task (auth required)
- `GET /tasks` - List tasks (auth required)
- `GET /tasks/{id}` - Get task by ID (auth required)
- `PATCH /tasks/{id}` - Update task (auth required)
- `DELETE /tasks/{id}` - Delete task (auth required)

## External Dependencies & Resources

### Core Libraries
- **Chi Router** (`github.com/go-chi/chi/v5`): HTTP router and middleware
- **SQLC** (`github.com/sqlc-dev/sqlc`): Type-safe SQL code generation
- **Goose** (`github.com/pressly/goose/v3`): Database migration tool
- **UUID** (`github.com/google/uuid`): UUID generation
- **SQLite** (`github.com/mattn/go-sqlite3`): SQLite database driver

### Authentication & Security
- **JWT** (`github.com/golang-jwt/jwt/v5`): JWT token handling
- **JWX** (`github.com/lestrrat-go/jwx/v2`): Advanced JWT/JWS operations

### Documentation
- **Swaggo** (`github.com/swaggo/swag`): Swagger documentation generation
- **HTTP-Swagger** (`github.com/swaggo/http-swagger`): Swagger UI middleware

### Useful Resources
- [Chi Documentation](https://go-chi.io/#/)
- [SQLC Documentation](https://docs.sqlc.dev/en/stable/)
- [Goose Migration Guide](https://github.com/pressly/goose)
- [OAuth 2.0 Client Credentials Flow](https://datatracker.ietf.org/doc/html/rfc6749#section-4.4)
- [iSHARE Protocol Documentation](https://dev.ishare.eu/)

### iSHARE Protocol Context
This API implements authentication patterns aligned with the **iSHARE protocol** for secure data sharing:

**Key iSHARE Principles:**
- **Machine-to-Machine (M2M)** communication via REST APIs
- **Public Key Infrastructure (PKI)** for authentication
- **Signed JSON Web Tokens (JWTs)** for message integrity
- Role-based access control with specific endpoints

**Technical Standards Used:**
- OAuth 2.0 Client Credentials Flow
- JSON Web Token (JWT) with signatures
- TLS for transport security
- RESTful API design

**Implementation Notes:**
- Uses certificate-based authentication (certificates in `secret/` directory)
- JWT tokens provide both authentication and message integrity
- Designed for secure interoperability between different systems
- Follows "soft infrastructure" model for identification and authorization

## Key Files

- `cmd/main.go` - Application entry point
- `internal/app/app.go` - Application bootstrap and dependency injection
- `internal/app/domain/task.go` - Task domain model and DTOs
- `internal/app/service/task_service.go` - Core business logic
- `internal/app/transport/httpserver/handlers/task.go` - HTTP handlers
- `internal/app/db/sqlite/database.go` - Database connection and setup