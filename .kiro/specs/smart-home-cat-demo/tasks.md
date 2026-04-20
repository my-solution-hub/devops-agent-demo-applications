# Implementation Plan: Smart Home Cat Demo

## Overview

This plan breaks the Smart Home Cat Demo into incremental coding tasks across the polyglot monorepo. Tasks are ordered so each step builds on the previous: project scaffolding first, then data models and core services, then AI agents, then frontends, then deployment artifacts and CI/CD. Property-based tests and unit tests are included as sub-tasks close to the code they validate.

## Tasks

- [x] 1. Scaffold project structure and shared configuration
  - [x] 1.1 Create monorepo directory structure
    - Create top-level directories: `langgraph-agent/`, `langgraph-agent/rag/`, `strands-agents/`, `chatbot/`, `device-simulator/`, `admin-console/`, `services/api-gateway/`, `services/cat-profile/`, `services/device/`, `services/feeding/`, `services/health-monitor/`, `deploy/k8s/`, `deploy/ecs/`, `deploy/lambda/`, `.github/workflows/`
    - Add placeholder README files in each directory
    - _Requirements: 15.1_

  - [x] 1.2 Create root README with architecture documentation
    - Document the three AI agent patterns (LangGraph, Strands, RAG)
    - Document project structure, prerequisites, and getting-started instructions
    - Include the high-level architecture diagram from the design
    - _Requirements: 15.2_

  - [x] 1.3 Create Docker Compose local development setup
    - Define services: `postgres`, `dynamodb-local`, `api-gateway`, `cat-profile`, `device-service`, `feeding-service`, `health-monitor`, `langgraph-agent`, `strands-agents`, `chatbot-ui`, `device-simulator`, `admin-console`
    - Configure PostgreSQL with initialization scripts for all tables (`cat_profiles`, `devices`, `device_assignments`, `feeding_schedules`)
    - Configure DynamoDB Local with table creation scripts (`device-shadows`, `device-commands`, `device-telemetry`, `feeding-events`, `health-metrics`, `health-alerts`)
    - Set up environment variables and service networking
    - _Requirements: 15.4_

- [x] 2. Implement Cat Profile Service (Spring Boot / Java / EKS)
  - [x] 2.1 Initialize Spring Boot project for Cat Profile Service
    - Create `services/cat-profile/` with Spring Boot starter (Web, JPA, PostgreSQL driver, Validation)
    - Define JPA entities: `CatProfile` (cat_id, owner_id, name, breed, age_months, weight_kg, dietary_restrictions, created_at, updated_at)
    - Define JPA entity: `DeviceAssignment` (assignment_id, cat_id, device_id, device_type, assigned_at) with unique constraint on (cat_id, device_type)
    - Create Spring Data JPA repositories for both entities
    - _Requirements: 8.1_

  - [x] 2.2 Implement Cat Profile REST endpoints
    - `POST /cats` — create cat profile with validation (name and weight required)
    - `GET /cats/{id}` — get cat profile by ID
    - `PUT /cats/{id}` — update cat profile with validation
    - `GET /cats` — list cats with optional `owner_id` query filter
    - `GET /cats/{id}/health` — get cat health summary (placeholder, delegates to Health Monitor later)
    - Return 400 with field-level errors for missing required fields
    - Return 404 for non-existent cat profiles
    - _Requirements: 8.1, 8.2_

  - [ ]* 2.3 Write property tests for Cat Profile Service (jqwik)
    - **Property 10: Cat profile round-trip persistence** — For any valid cat profile, store and retrieve returns matching fields
    - **Property 11: Cat profile validation rejects missing required fields** — For any profile missing name or weight, create/update returns descriptive error
    - **Validates: Requirements 8.1, 8.2**

  - [ ]* 2.4 Write unit tests for Cat Profile Service
    - Test CRUD operations with valid data
    - Test validation error responses for missing name, missing weight, missing both
    - Test owner_id filtering on list endpoint
    - Test 404 response for non-existent profile
    - _Requirements: 8.1, 8.2_

- [x] 3. Implement Device Service (Go / Lambda)
  - [x] 3.1 Initialize Go project for Device Service
    - Create `services/device/` with Go module, Lambda handler entry point
    - Set up PostgreSQL connection (pgx or sqlx) for device registry
    - Set up DynamoDB client for device commands, telemetry, and shadow tables
    - Define Go structs: `Device`, `CommandRecord`, `DeviceShadow`, `TelemetryRecord`
    - _Requirements: 10.1, 10.5_

  - [x] 3.2 Implement device CRUD endpoints
    - `POST /devices` — register device (persist to PostgreSQL)
    - `GET /devices/{id}` — get device including shadow state from DynamoDB
    - `PUT /devices/{id}` — update device configuration
    - `GET /devices` — list devices
    - `GET /devices/{id}/shadow` — get device shadow (desired + reported) from DynamoDB
    - `PUT /devices/{id}/shadow` — update desired state in DynamoDB
    - _Requirements: 10.5_

  - [x] 3.3 Implement device command endpoints
    - `POST /devices/{id}/commands` — submit command, create command record (status: pending), update desired state in DynamoDB
    - `GET /devices/{id}/commands/{cmd_id}` — get command status
    - `GET /devices/{id}/commands` — list recent commands (with optional `status` filter for polling)
    - `POST /devices/{id}/commands/{cmd_id}/ack` — acknowledge command, update status to acknowledged, update reported state
    - Implement background timeout logic: mark commands as `timed_out` if not acknowledged within 10 seconds
    - _Requirements: 10.1, 10.2, 10.3, 10.4_

  - [ ]* 3.4 Write property tests for Device Service (rapid)
    - **Property 18: Device command round-trip** — For any valid command, submit and retrieve returns matching fields with pending or acknowledged status
    - **Property 19: Device shadow round-trip** — For any device shadow update, retrieve returns matching desired and reported state
    - **Validates: Requirements 10.1, 10.2, 10.5**

  - [ ]* 3.5 Write unit tests for Device Service
    - Test device CRUD operations
    - Test command lifecycle (pending → acknowledged, pending → timed_out)
    - Test shadow state updates (desired and reported)
    - Test command status filtering on list endpoint
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [ ] 4. Checkpoint — Core data services
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. Implement API Gateway Service (Spring Boot / Java / EKS)
  - [ ] 5.1 Initialize Spring Boot project for API Gateway Service
    - Create `services/api-gateway/` with Spring Boot starter (Web, Security, WebClient)
    - Configure JWT validation using Cognito JWKS endpoint
    - Implement role-based access control: reject 401 for missing/invalid JWT, reject 403 for non-admin on admin endpoints
    - _Requirements: 14.4, 14.5_

  - [ ] 5.2 Implement request routing to downstream services
    - Route `/api/cats/**` → Cat Profile Service
    - Route `/api/devices/**` → Device Service
    - Route `/api/feeding/**` → Feeding Service
    - Route `/api/health/**` → Health Monitor Service
    - Propagate W3C Trace Context headers (`traceparent`, `tracestate`) to downstream requests
    - Configure downstream service URLs via environment variables (for multi-compute discovery)
    - _Requirements: 11.4, 11.5_

  - [ ]* 5.3 Write property tests for API Gateway Service (jqwik)
    - **Property 20: API Gateway request path routing** — For any request with a known path pattern, route to the correct downstream service
    - **Property 21: Trace context header propagation** — For any request with W3C Trace Context headers, propagate them to downstream
    - **Property 22: Invalid JWT rejection** — For any request without a valid JWT, respond with 401
    - **Property 23: Non-admin role rejection on admin endpoints** — For any non-admin JWT on admin endpoint, respond with 403
    - **Validates: Requirements 11.4, 11.5, 14.4, 14.5**

  - [ ]* 5.4 Write unit tests for API Gateway Service
    - Test routing for each path pattern
    - Test JWT validation (valid, expired, malformed, missing)
    - Test role-based access (admin vs owner)
    - Test trace header propagation
    - _Requirements: 11.4, 11.5, 14.4, 14.5_

- [ ] 6. Implement Feeding Service (Django / Python / ECS Fargate)
  - [ ] 6.1 Initialize Django project for Feeding Service
    - Create `services/feeding/` with Django project, REST framework
    - Define Django model: `FeedingSchedule` (schedule_id, cat_id, device_id, meal_times, portion_grams, max_daily_grams, min_interval_minutes, enabled, created_at) in PostgreSQL
    - Set up DynamoDB client (boto3) for feeding events table
    - _Requirements: 8.5_

  - [ ] 6.2 Implement Feeding Service REST endpoints
    - `POST /feeding/schedules` — create feeding schedule
    - `GET /feeding/schedules/{cat_id}` — get schedules for a cat
    - `PUT /feeding/schedules/{id}` — update schedule
    - `GET /feeding/history/{cat_id}` — get feeding history from DynamoDB with time range filter
    - `POST /feeding/dispense` — trigger immediate feeding (invoke Device Service to update feeder desired state)
    - Implement duplicate feeding prevention: reject feed if cat was fed within configurable minimum interval
    - _Requirements: 8.5, 9.1, 9.2, 9.4, 9.5_

  - [ ] 6.3 Implement scheduled feeding automation
    - Set up Celery beat (or Django-Q) to check schedules every minute
    - On schedule match: invoke Device Service to update feeder desired state to dispense
    - On device offline: retry 3× at 60-second intervals, then generate alert via Health Monitor Service
    - Record feeding events as append-only entries in DynamoDB
    - _Requirements: 9.1, 9.2, 9.3_

  - [ ]* 6.4 Write property tests for Feeding Service (Hypothesis)
    - **Property 14: Feeding schedule round-trip persistence** — For any valid schedule, store and retrieve returns matching fields
    - **Property 15: Feeding event recording round-trip** — For any completed feeding event, retrieve history returns matching event
    - **Property 16: Feeding history query filtering** — For any cat_id and time range, return only matching events
    - **Property 17: Duplicate feeding prevention** — For two feed commands within minimum interval, reject the second
    - **Validates: Requirements 8.5, 9.2, 9.4, 9.5**

  - [ ]* 6.5 Write unit tests for Feeding Service
    - Test schedule CRUD operations
    - Test feeding history time range filtering
    - Test duplicate feeding rejection (409 Conflict)
    - Test scheduled feeding trigger logic
    - Test retry logic for offline devices
    - _Requirements: 8.5, 9.1, 9.2, 9.3, 9.4, 9.5_

- [ ] 7. Implement Health Monitor Service (Django / Python / ECS Fargate)
  - [ ] 7.1 Initialize Django project for Health Monitor Service
    - Create `services/health-monitor/` with Django project, REST framework
    - Set up DynamoDB client (boto3) for health-metrics and health-alerts tables
    - _Requirements: 8.3, 8.4_

  - [ ] 7.2 Implement Health Monitor REST endpoints
    - `GET /health/{cat_id}/summary` — aggregated health summary from DynamoDB telemetry
    - `GET /health/{cat_id}/alerts` — active health alerts for a cat
    - `GET /health/{cat_id}/metrics` — raw health metrics with time range filter
    - `GET /health/alerts` — all active alerts (admin endpoint)
    - Implement telemetry aggregation: read telemetry from DynamoDB, compute per-cat averages, counts, min/max
    - Implement threshold-based alert generation: generate alerts when metrics deviate beyond configurable thresholds
    - _Requirements: 8.3, 8.4_

  - [ ]* 7.3 Write property tests for Health Monitor Service (Hypothesis)
    - **Property 12: Health telemetry aggregation correctness** — For any set of telemetry data points, aggregated values are mathematically consistent
    - **Property 13: Health alert threshold detection** — For any metric value and threshold, alert generated if and only if value deviates beyond threshold
    - **Validates: Requirements 8.3, 8.4**

  - [ ]* 7.4 Write unit tests for Health Monitor Service
    - Test health summary aggregation with known data
    - Test alert generation for threshold violations
    - Test no alert for values within threshold
    - Test time range filtering on metrics endpoint
    - _Requirements: 8.3, 8.4_

- [ ] 8. Checkpoint — All backend microservices
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. Implement LangGraph Workflow Agent (Python / AgentCore)
  - [ ] 9.1 Initialize LangGraph agent project
    - Create `langgraph-agent/` with Python project (pyproject.toml or requirements.txt)
    - Install dependencies: langgraph, langchain, boto3, httpx
    - Define `CatCareState` TypedDict with fields: messages, intent, entities, cat_profile, rag_context, action_result, needs_clarification
    - _Requirements: 1.9_

  - [ ] 9.2 Implement LangGraph workflow graph nodes
    - Implement `intent_node` — classify user intent (device_command, query, knowledge, clarification) and extract entities
    - Implement `cat_profile_node` — look up cat profile via API Gateway (`GET /api/cats?name=...`)
    - Implement `action_node` — determine which action to take based on intent and context
    - Implement `device_command_node` — execute device commands via API Gateway (`POST /api/devices/{id}/commands`)
    - Implement `rag_node` — invoke RAG pipeline to retrieve domain knowledge
    - Implement `response_node` — generate natural language response using accumulated context
    - Implement `clarification_node` — ask user for clarification when intent is ambiguous
    - _Requirements: 1.1, 1.3, 1.4, 1.5, 1.6_

  - [ ] 9.3 Wire LangGraph graph with conditional edges
    - Define conditional routing from `intent_node` → `cat_profile_node` (device/query), `rag_node` (knowledge), `clarification_node` (ambiguous)
    - Define edges: `cat_profile_node` → `action_node`, `action_node` → `device_command_node` or `rag_node`, `device_command_node` → `response_node`, `rag_node` → `response_node`
    - Ensure state persists across all node transitions within a session
    - Implement the `invoke` async function as the AgentCore entrypoint
    - _Requirements: 1.1, 1.2, 1.7_

  - [ ]* 9.4 Write property tests for LangGraph Workflow (Hypothesis)
    - **Property 1: LangGraph intent-based routing** — For any intent classification, route through the correct node sequence
    - **Property 2: LangGraph state persistence across nodes** — For any state key-value pair set by a node, subsequent nodes have access to it
    - **Validates: Requirements 1.3, 1.4, 1.5, 1.7**

  - [ ]* 9.5 Write unit tests for LangGraph Workflow
    - Test intent classification for each intent type
    - Test node transitions for device command flow
    - Test node transitions for knowledge query flow
    - Test clarification flow for ambiguous input
    - Test state accumulation across nodes
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_

- [ ] 10. Implement RAG Pipeline (Python / within LangGraph)
  - [ ] 10.1 Implement RAG pipeline components
    - Create `langgraph-agent/rag/` directory
    - Implement `document_loader` — load and chunk cat-care documents from Knowledge Base
    - Implement `embedder` — convert text chunks into vector embeddings (sentence-transformers locally, Titan Embeddings for cloud)
    - Implement `vector_store` — FAISS for local dev, Amazon Bedrock Knowledge Base for cloud
    - Implement `retriever` — similarity search returning top-k relevant chunks ordered by descending relevance score
    - Implement `indexer` — incremental indexing of new documents without full re-index
    - _Requirements: 2.1, 2.2, 2.3, 2.6, 2.7_

  - [ ] 10.2 Create Knowledge Base documents
    - Create cat-care domain documents: feeding guidelines (portion sizes by breed/age/weight), breed-specific dietary needs, health monitoring tips, device troubleshooting guides
    - Index documents into the vector store
    - _Requirements: 2.5_

  - [ ] 10.3 Integrate RAG pipeline with LangGraph rag_node
    - Wire `rag_node` to invoke `RAGPipeline.retrieve()` and pass retrieved context to `response_node`
    - Ensure retrieved context is included in the LLM prompt for response generation
    - _Requirements: 2.2, 2.4_

  - [ ]* 10.4 Write property tests for RAG Pipeline (Hypothesis)
    - **Property 3: RAG retrieval returns bounded, ordered results** — For any query and seeded vector store, return at most top-k chunks ordered by descending relevance
    - **Property 4: Incremental RAG indexing preserves existing documents** — For any new document added, previously indexed documents remain retrievable
    - **Validates: Requirements 2.2, 2.6**

  - [ ]* 10.5 Write unit tests for RAG Pipeline
    - Test document loading and chunking
    - Test embedding generation
    - Test retrieval with known documents and queries
    - Test incremental indexing preserves existing documents
    - Test empty query handling
    - _Requirements: 2.1, 2.2, 2.3, 2.6_

- [ ] 11. Implement Strands Multi-Agent System (Python / AgentCore)
  - [ ] 11.1 Initialize Strands agents project
    - Create `strands-agents/` with Python project (pyproject.toml or requirements.txt)
    - Install dependencies: strands-agents SDK, boto3, httpx
    - _Requirements: 3.10_

  - [ ] 11.2 Implement specialist agents
    - Implement `FeedingAgent` — handle feeding commands, schedules, history via Feeding Service and Device Service
    - Implement `HealthAgent` — handle health queries, alerts, metrics via Health Monitor Service and Cat Profile Service
    - Implement `DeviceAgent` — handle device control, status, troubleshooting via Device Service
    - Each agent uses Strands Agent class with domain-specific system prompt and tools
    - _Requirements: 3.3, 3.4, 3.5_

  - [ ] 11.3 Implement Orchestrator Agent
    - Implement `OrchestratorAgent` — classify intent and route to appropriate specialist agent
    - Handle multi-domain requests by coordinating across specialist agents and combining results
    - Implement error handling: catch specialist errors and return descriptive error messages
    - Wire AgentCore Memory for conversation history persistence
    - _Requirements: 3.1, 3.2, 3.6, 3.8, 3.9_

  - [ ]* 11.4 Write property tests for Strands Agents (Hypothesis)
    - **Property 5: Strands Orchestrator routes to correct specialist** — For any classifiable intent, route to the corresponding specialist agent
    - **Property 6: Strands Orchestrator returns descriptive errors on specialist failure** — For any specialist error, return non-empty error message with context
    - **Validates: Requirements 3.2, 3.9**

  - [ ]* 11.5 Write unit tests for Strands Agents
    - Test Orchestrator routing for feeding, health, and device intents
    - Test multi-domain request coordination
    - Test error handling when specialist agent fails
    - Test individual specialist agent logic
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.9_

- [ ] 12. Implement AgentCore Entrypoint
  - [ ] 12.1 Implement AgentCore Entrypoint routing
    - Implement `POST /invoke` endpoint — route to LangGraph or Strands based on `agent_type` parameter
    - Implement `GET /agents` endpoint — list available agent implementations
    - Return error with available options for missing/invalid `agent_type`
    - Pass JWT token through to downstream agent without modification
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

  - [ ]* 12.2 Write property tests for AgentCore Entrypoint (Hypothesis)
    - **Property 7: AgentCore Entrypoint routes to correct agent implementation** — For any valid agent selection, forward to the corresponding agent
    - **Property 8: AgentCore Entrypoint rejects invalid agent selection** — For any invalid agent_type, return error listing available agents
    - **Property 9: AgentCore JWT passthrough** — For any JWT token, forward identical token to downstream agent
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5**

  - [ ]* 12.3 Write unit tests for AgentCore Entrypoint
    - Test routing to LangGraph agent
    - Test routing to Strands agent
    - Test error response for invalid agent_type
    - Test JWT passthrough
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 13. Checkpoint — All backend services and AI agents
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 14. Implement Chatbot UI (React / TypeScript / ECS Fargate)
  - [ ] 14.1 Initialize React project for Chatbot UI
    - Create `chatbot/` with React + TypeScript project (Vite or Create React App)
    - Set up Cognito authentication (redirect to login if unauthenticated)
    - _Requirements: 6.2_

  - [ ] 14.2 Implement chat interface components
    - Build message input and send functionality
    - Implement agent selection toggle (LangGraph vs Strands)
    - Implement real-time streamed response display from AgentCore Gateway
    - Build scrollable conversation history
    - Implement connection status indicator with auto-retry on disconnect
    - _Requirements: 6.1, 6.3, 6.4, 6.5, 6.6_

  - [ ]* 14.3 Write unit tests for Chatbot UI
    - Test message rendering
    - Test agent selection toggle state
    - Test connection status indicator states
    - _Requirements: 6.1, 6.3, 6.5, 6.6_

- [ ] 15. Implement Device Simulator (React / TypeScript / ECS Fargate)
  - [ ] 15.1 Initialize React project for Device Simulator
    - Create `device-simulator/` with React + TypeScript project
    - Set up Cognito authentication
    - _Requirements: 5.1_

  - [ ] 15.2 Implement device simulation logic
    - Implement simulated device types: automatic feeder, water fountain, litter box monitor, activity tracker
    - Implement command polling: poll Device Service for pending commands and acknowledge via REST
    - Implement periodic telemetry publishing: write telemetry records to Device Service at configurable intervals
    - Implement fault condition simulation: create alert records via Health Monitor Service REST API
    - Build React UI displaying real-time device states and telemetry
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

  - [ ]* 15.3 Write unit tests for Device Simulator
    - Test device type rendering
    - Test command polling and acknowledgment logic
    - Test telemetry publishing intervals
    - _Requirements: 5.1, 5.2, 5.3_

- [ ] 16. Implement Admin Console (React / TypeScript / ECS Fargate)
  - [ ] 16.1 Initialize React project for Admin Console
    - Create `admin-console/` with React + TypeScript project
    - Set up Cognito authentication with admin role verification
    - _Requirements: 7.5_

  - [ ] 16.2 Implement admin dashboard and management views
    - Build dashboard: device status overview, microservice health, agent status
    - Build device management: register devices, view device state, telemetry, command history
    - Build cat profile management: view and edit cat profiles
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [ ]* 16.3 Write unit tests for Admin Console
    - Test dashboard rendering with mock data
    - Test device management CRUD interactions
    - Test admin role gate (redirect non-admin users)
    - _Requirements: 7.1, 7.2, 7.3, 7.5_

- [ ] 17. Checkpoint — All frontends
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 18. Create deployment artifacts
  - [ ] 18.1 Create Dockerfiles for all services and frontends
    - Dockerfile for API Gateway Service (Spring Boot, multi-stage build)
    - Dockerfile for Cat Profile Service (Spring Boot, multi-stage build)
    - Dockerfile for Feeding Service (Django)
    - Dockerfile for Health Monitor Service (Django)
    - Dockerfile for Device Service (Go, compile to binary for Lambda; also containerized for local dev)
    - Dockerfile for Chatbot UI (React, build static assets, serve via Nginx)
    - Dockerfile for Device Simulator (React, build static assets, serve via Nginx)
    - Dockerfile for Admin Console (React, build static assets, serve via Nginx)
    - Dockerfile for LangGraph Agent (Python)
    - Dockerfile for Strands Agents (Python)
    - _Requirements: 12.1, 12.5, 12.7_

  - [ ] 18.2 Create Kubernetes manifests for EKS services
    - Create `deploy/k8s/` with Deployment, Service, ConfigMap, and Ingress resources for API Gateway Service and Cat Profile Service
    - Configure liveness and readiness probes
    - Reference environment-specific variables from sibling infra project
    - _Requirements: 12.2, 12.6_

  - [ ] 18.3 Create ECS task definitions for Fargate services
    - Create `deploy/ecs/` with task definitions and service configurations for: Chatbot UI, Admin Console, Device Simulator, Feeding Service, Health Monitor Service
    - Configure health checks
    - Reference environment-specific variables from sibling infra project
    - _Requirements: 12.3, 12.6_

  - [ ] 18.4 Create Lambda deployment configuration
    - Create `deploy/lambda/` with Lambda deployment package configuration for Device Service (Go)
    - Include compiled Go binary packaging and configuration
    - Reference environment-specific variables from sibling infra project
    - _Requirements: 12.4, 12.6_

- [ ] 19. Create GitHub Actions CI/CD pipelines
  - [ ] 19.1 Create CI workflow for pull requests
    - Create `.github/workflows/ci.yml` triggered on pull requests
    - Run linting, unit tests, and build validation for all components
    - Configure OIDC authentication to AWS (no stored access keys)
    - _Requirements: 13.1, 13.4_

  - [ ] 19.2 Create deployment workflows for main branch
    - Create `.github/workflows/deploy-services.yml` — build Spring Boot images → push to ECR → deploy to EKS
    - Create `.github/workflows/deploy-ecs-services.yml` — build Django images → push to ECR → deploy to ECS Fargate
    - Create `.github/workflows/deploy-frontends.yml` — build React images → push to ECR → deploy to ECS Fargate
    - Create `.github/workflows/deploy-lambda.yml` — build Go binary → package → deploy to Lambda
    - Create `.github/workflows/deploy-agents.yml` — build agent packages → deploy to AgentCore
    - All workflows use OIDC authentication and environment-specific configuration
    - _Requirements: 13.1, 13.2, 13.3, 13.5_

- [ ] 20. Final checkpoint — Full system
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation at logical boundaries
- Property tests validate universal correctness properties from the design document
- Unit tests validate specific examples and edge cases
- All infrastructure (VPC, EKS, ECS, Lambda, Cognito, OIDC) is provisioned by the sibling infra project — this plan covers application code and deployment artifacts only
- The Docker Compose setup (task 1.3) enables local end-to-end testing without AWS dependencies
