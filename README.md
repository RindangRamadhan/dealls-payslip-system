# Payslip System Documentation

## Quick Start Guide

### Prerequisites

- Docker and Docker Compose
- Go 1.19+
- Make

### Running the Application

```bash
# Start services with Docker
make compose-up

# Run application with database migrations
make run
```

### Service Endpoints

- **Health Check**: http://127.0.0.1:8080/health
- **Metrics**: http://127.0.0.1:8080/metrics
- **API Documentation**: http://127.0.0.1:8080/swagger

## API Usage

### REST API

The system provides RESTful endpoints accessible at `http://127.0.0.1:8080/v1/`

**Authentication**: Include appropriate headers as defined in your configuration

**Content-Type**: `application/json`

**Example Request**:

```bash
curl -X GET http://127.0.0.1:8080/v1/payslips \
  -H "Content-Type: application/json"
```

### gRPC API

gRPC services are available with protocol buffer definitions in `docs/proto/`

### Message Queue (AMQP)

RabbitMQ-based RPC pattern for asynchronous operations

## Software Architecture

### Architecture Pattern: Clean Architecture

The system follows Clean Architecture principles with clear separation of concerns:
![Example](docs/img/layers-1.png)
![Example](docs/img/layers-2.png)

### Key Principles

- **Dependency Inversion**: Inner layers don't depend on outer layers
- **Interface-based Communication**: Layers communicate through interfaces
- **Business Logic Independence**: Core logic remains framework-agnostic

### Project Structure

```
├── cmd/app/               # Application entry point
├── config/                # Configuration management
├── docs/                  # API documentation
├── internal/
│   ├── app/               # Dependency injection & app setup
│   ├── controller/        # HTTP, gRPC
│   │   ├── http/v1/       # REST API endpoints
│   │   ├── grpc/v1/       # gRPC service handlers
│   ├── usecase/           # Business logic
│   ├── entity/            # Domain models
│   └── repo/              # Data access interfaces
│       ├── persistent/    # Database repositories
│       └── webapi/        # External API clients
├── pkg/                   # Shared packages
└── integration-test/      # Integration tests
```

### Data Flow Example

```
HTTP > usecase
        usecase > repository (Postgres)
        usecase < repository (Postgres)
HTTP < usecase
```

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```env
# App
APP_NAME=dealls-payslip-system
APP_VERSION=1.0.0

# HTTP settings
HTTP_PORT=8080
HTTP_USE_PREFORK_MODE=false

# Logger
LOG_LEVEL=debug

# PG
PG_POOL_MAX=2
PG_URL=postgres://user:pass@localhost:5433/db_name

# Metrics
METRICS_ENABLED=true

# Swagger
SWAGGER_ENABLED=true

# JWT
JWT_SECRET=somejwtsecret
```

### Database Migrations

Migrations run automatically with the `migrate` tag:

```bash
# Create new migration
make migrate-create

# Run database migration
make migrate-up
```

## Development Guidelines

### API Versioning

- REST: `/v1/`, `/v2/` URL prefixes
- gRPC: Separate `v1/`, `v2/` service packages

### Testing

```bash
# Unit tests
go test ./...

or

make test
```

## Monitoring & Observability

- **Health Checks**: `/health` endpoint
- **Metrics**: Prometheus-compatible `/metrics`
- **Logging**: Structured logging with configurable levels
- **Tracing**: Request correlation IDs

---

For detailed API specifications, visit the Swagger documentation at `/swagger` when the application is running.
