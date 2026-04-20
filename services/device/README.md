# Device Service

Go service for device registration, command handling, telemetry, and device shadow management. Uses PostgreSQL (pgx) for the device registry and DynamoDB for commands, telemetry, and shadow state.

## Architecture

The service can run in two modes:
- **Lambda** — Production deployment via AWS Lambda using `aws-lambda-go-api-proxy/chiadapter`
- **Standalone HTTP server** — Local development via Docker Compose

## Project Structure

```
services/device/
├── cmd/
│   ├── lambda/main.go       # Lambda handler entry point
│   └── server/main.go       # Standalone HTTP server for local dev
├── internal/
│   ├── handler/
│   │   ├── device.go        # Device CRUD + shadow handlers
│   │   └── command.go       # Device command handlers
│   ├── model/
│   │   └── models.go        # Go structs (Device, CommandRecord, DeviceShadow, TelemetryRecord)
│   ├── store/
│   │   ├── postgres.go      # PostgreSQL operations (device registry)
│   │   └── dynamo.go        # DynamoDB operations (commands, shadow)
│   └── router/
│       └── router.go        # Route definitions (chi)
├── go.mod
└── go.sum
```

## Endpoints

### Device CRUD (PostgreSQL)
- `POST /devices` — Register a new device
- `GET /devices` — List all devices
- `GET /devices/{id}` — Get device (merged with shadow from DynamoDB)
- `PUT /devices/{id}` — Update device configuration

### Device Shadow (DynamoDB)
- `GET /devices/{id}/shadow` — Get device shadow (desired + reported state)
- `PUT /devices/{id}/shadow` — Update desired state

### Device Commands (DynamoDB)
- `POST /devices/{id}/commands` — Submit command (creates record with status=pending, updates desired state)
- `GET /devices/{id}/commands` — List commands (with optional `?status=pending` filter)
- `GET /devices/{id}/commands/{cmd_id}` — Get command status
- `POST /devices/{id}/commands/{cmd_id}/ack` — Acknowledge command (updates status, updates reported state)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | (required) |
| `DYNAMODB_ENDPOINT` | DynamoDB endpoint (for local dev) | (uses AWS default) |
| `AWS_REGION` | AWS region | (from AWS config) |
| `PORT` | HTTP server port (standalone mode only) | `8082` |

## Command Timeout

Commands that are not acknowledged within 10 seconds are automatically marked as `timed_out`.

- **Standalone server**: A background goroutine checks every 5 seconds
- **Lambda**: Would be handled by a separate scheduled Lambda invocation (e.g., via EventBridge)

## Building

```bash
# Build all packages
go build ./...

# Build Lambda binary
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go

# Build standalone server
go build -o device-server cmd/server/main.go
```

## Running Locally

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/catdemo?sslmode=disable"
export DYNAMODB_ENDPOINT="http://localhost:8000"
export AWS_REGION="us-east-1"
go run cmd/server/main.go
```
