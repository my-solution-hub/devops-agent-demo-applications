# Running Locally

This guide walks you through running the Smart Home Cat Demo on your local machine using Docker Compose.

## Prerequisites

- **Docker** (v20.10+) and **Docker Compose** (v2.0+)
- No other dependencies needed — everything runs in containers

## Quick Start

```bash
# Start the core services (databases + backend + UI)
docker-compose up --build
```

This builds and starts all services. First run will take a few minutes to download images and build containers.

## Services

Once running, you'll have:

| Service | URL | Description |
|---------|-----|-------------|
| Chatbot UI | http://localhost:3000 | Chat interface for interacting with the system |
| API Gateway | http://localhost:8080 | Routes requests to backend microservices |
| Cat Profile Service | http://localhost:8081 | Cat profile CRUD (Spring Boot + PostgreSQL) |
| Device Service | http://localhost:8082 | Device management + commands (Go + PostgreSQL + DynamoDB) |
| PostgreSQL | localhost:5432 | Relational data (user: `catdemo`, password: `catdemo`, db: `catdemo`) |
| DynamoDB Local | localhost:8000 | Time-series and event data |

## Starting Only What's Available

If you only want to run the services that are fully implemented:

```bash
docker-compose up --build postgres dynamodb-local dynamodb-init api-gateway cat-profile device-service chatbot-ui
```

## Demo Walkthrough

1. Open http://localhost:3000 in your browser

2. Try these commands in the chat:

   | Command | What it does |
   |---------|--------------|
   | `add cat` | Creates a sample cat profile (Whiskers, Maine Coon, 5.5kg) |
   | `list cats` | Shows all cat profiles from PostgreSQL |
   | `add device` | Registers a sample feeder device |
   | `list devices` | Shows all registered devices |

3. You can also hit the APIs directly:

   ```bash
   # Create a cat
   curl -X POST http://localhost:8080/api/cats \
     -H "Content-Type: application/json" \
     -d '{"name": "Luna", "breed": "Siamese", "weightKg": 4.2, "ownerId": "user-1"}'

   # List cats
   curl http://localhost:8080/api/cats

   # Register a device
   curl -X POST http://localhost:8080/api/devices \
     -H "Content-Type: application/json" \
     -d '{"device_type": "feeder", "name": "Living Room Feeder", "config": {"portion_grams": 40}}'

   # List devices
   curl http://localhost:8080/api/devices

   # Send a command to a device
   curl -X POST http://localhost:8080/api/devices/{device_id}/commands \
     -H "Content-Type: application/json" \
     -d '{"action": "dispense", "params": {"amount_grams": "50"}}'

   # Check command status
   curl http://localhost:8080/api/devices/{device_id}/commands/{command_id}

   # Get device shadow (desired + reported state)
   curl http://localhost:8080/api/devices/{device_id}/shadow
   ```

## Architecture (Local)

```
┌─────────────────┐     ┌──────────────────┐
│   Chatbot UI    │────▶│   API Gateway    │
│  localhost:3000 │     │  localhost:8080   │
└─────────────────┘     └────────┬─────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼                         ▼
         ┌──────────────────┐     ┌──────────────────┐
         │ Cat Profile Svc  │     │  Device Service  │
         │  localhost:8081  │     │  localhost:8082  │
         │  (Spring Boot)   │     │      (Go)       │
         └────────┬─────────┘     └───┬─────────┬───┘
                  │                   │         │
                  ▼                   ▼         ▼
         ┌──────────────┐   ┌──────────────┐ ┌──────────────┐
         │  PostgreSQL  │   │  PostgreSQL  │ │DynamoDB Local│
         │  :5432       │   │  :5432       │ │  :8000       │
         └──────────────┘   └──────────────┘ └──────────────┘
```

## Stopping

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## Troubleshooting

### Port conflicts

If you get port-in-use errors, check what's running:

```bash
lsof -i :3000  # Chatbot UI
lsof -i :5432  # PostgreSQL
lsof -i :8000  # DynamoDB Local
lsof -i :8080  # API Gateway
lsof -i :8081  # Cat Profile
lsof -i :8082  # Device Service
```

### Database not ready

If services fail to connect to PostgreSQL or DynamoDB, they may have started before the databases were ready. Restart them:

```bash
docker-compose restart api-gateway cat-profile device-service
```

### Rebuilding after code changes

```bash
# Rebuild a specific service
docker-compose up --build cat-profile

# Rebuild everything
docker-compose up --build
```

### Viewing logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f device-service
```

## What's Not Running Yet

These services are defined in docker-compose.yml but don't have implementations yet:

- **Feeding Service** (Django) — port 8083
- **Health Monitor Service** (Django) — port 8084
- **LangGraph Agent** (Python) — port 8090
- **Strands Agents** (Python) — port 8091
- **Device Simulator** (React) — port 3001
- **Admin Console** (React) — port 3002

They'll fail to build until their code is implemented. Use the selective start command above to avoid errors.
