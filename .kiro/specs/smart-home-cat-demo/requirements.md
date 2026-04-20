# Requirements Document

## Introduction

Smart Home Cat Demo is an AI-first microservice demo application for managing cat-care IoT devices through natural language interaction. The system combines multiple AI agent patterns — a LangGraph stateful workflow with embedded RAG, a Strands-based multi-agent system on AWS AgentCore, and a simple RAG pipeline — with a cat-focused IoT device ecosystem, React frontends, and polyglot microservices (Spring Boot, Python Django, Go). Everything is deployed across multiple AWS compute platforms (EKS, ECS Fargate, Lambda, AgentCore Runtime) using GitHub Actions CI/CD with OIDC authentication. The system demonstrates polyglot persistence by combining PostgreSQL (relational/structured data) with DynamoDB (time-series/event data) across services.

This repository contains **application code only**. All AWS infrastructure (VPC, EKS cluster, ECS services, Lambda functions, OIDC provider, IAM roles) is provisioned by a sibling infrastructure project in a separate repository. This project produces deployment artifacts (Dockerfiles, Kubernetes manifests, ECS task definitions, Lambda deployment packages) that target the pre-existing infrastructure.

The project showcases three distinct AI agent design patterns — graph-based workflow orchestration (LangGraph), multi-agent collaboration (Strands Agents), and retrieval-augmented generation (RAG) — applied to a conversational AI agent that controls simulated cat-care devices (feeders, water fountains, litter box monitors, activity trackers, etc.). The backend uses a polyglot microservice architecture with Spring Boot (Java), Python Django, and Go to demonstrate multi-language service development, distributed across all four compute types (EKS, ECS Fargate, Lambda, AgentCore Runtime). The system uses a dual storage pattern: PostgreSQL for relational/structured data (cat profiles, feeding schedules, device registry) and DynamoDB for high-throughput time-series/event data (telemetry, commands, feeding events, health metrics). Device simulation is handled entirely through REST API calls to the Device Service, which persists device state across PostgreSQL (registry) and DynamoDB (commands, telemetry, shadow). There is no message broker — "devices" are simply data records whose state is updated via REST endpoints.

## Glossary

- **AI_Agent**: A general term covering all three AI agent implementations in this project: the LangGraph_Workflow, the Strands_Agents multi-agent system, and the RAG_Pipeline
- **LangGraph_Workflow**: A Python-based stateful workflow agent built with LangGraph that orchestrates cat-care operations as a directed graph of nodes, including an embedded RAG_Node for knowledge retrieval
- **Strands_Agents**: A Python-based multi-agent system built with the Strands Agents SDK, consisting of multiple specialized agents (Feeding_Agent, Health_Agent, Device_Agent) coordinated by an Orchestrator_Agent, deployed to AgentCore_Runtime
- **RAG_Pipeline**: A simple Retrieval-Augmented Generation pipeline embedded as a node within the LangGraph_Workflow that retrieves cat-care domain knowledge from a Vector_Store before generating responses
- **Knowledge_Base**: A collection of cat-care domain documents (feeding guidelines, breed-specific dietary needs, health monitoring tips, device troubleshooting) indexed in the Vector_Store for retrieval by the RAG_Pipeline
- **Vector_Store**: A vector database (e.g., FAISS or ChromaDB locally, or Amazon Bedrock Knowledge Base) used by the RAG_Pipeline to store and retrieve document embeddings from the Knowledge_Base
- **Orchestrator_Agent**: A Strands agent that receives user requests and routes them to the appropriate specialist agent (Feeding_Agent, Health_Agent, or Device_Agent)
- **Feeding_Agent**: A Strands specialist agent responsible for handling feeding-related commands and queries (scheduling, dispensing, feeding history)
- **Health_Agent**: A Strands specialist agent responsible for handling health-related queries (health summaries, alerts, activity metrics)
- **Device_Agent**: A Strands specialist agent responsible for handling device control and status commands (device state, telemetry, troubleshooting)
- **RAG_Node**: A node within the LangGraph_Workflow graph that invokes the RAG_Pipeline to retrieve relevant domain knowledge before generating a response
- **Chatbot_UI**: A React frontend application providing a conversational interface for users to interact with any of the AI agent implementations via natural language
- **Device_Simulator**: A React frontend application that simulates cat-care IoT devices by making REST API calls to the Device_Service to update device state and telemetry in DynamoDB
- **Admin_Console**: A React frontend application for managing devices, viewing system health, and configuring the AI agent implementations
- **API_Gateway_Service**: A Spring Boot (Java) microservice acting as the entry point for backend API requests, routing to downstream services
- **Cat_Profile_Service**: A Spring Boot (Java) microservice managing cat profiles, owner information, and cat health records. Stores relational domain data (cat profiles, owner associations, device assignments) in PostgreSQL using JPA/Hibernate.
- **Device_Service**: A Go microservice managing IoT device registration, configuration, state persistence, and device commands via REST API. Stores the device registry (device metadata, cat assignments) in PostgreSQL using pgx or sqlx, and stores device commands and telemetry in DynamoDB for high-throughput time-series access. Maintains device shadow (desired/reported state) in DynamoDB and handles command execution by updating desired state.
- **Feeding_Service**: A Python Django microservice managing feeding schedules, portion control, and feeding history for cats. Stores feeding schedules in PostgreSQL using Django ORM for structured querying by time, and stores feeding events in DynamoDB as an append-only log.
- **Health_Monitor_Service**: A Python Django microservice tracking cat health metrics from device telemetry data (weight, activity, litter box usage). Reads telemetry from DynamoDB and stores health metrics and health alerts in DynamoDB for time-series and event-driven access.
- **AgentCore_Runtime**: AWS AgentCore Runtime service hosting all AI agent implementations (LangGraph_Workflow and Strands_Agents)
- **AgentCore_Memory**: AWS AgentCore Memory service providing conversation history and context persistence for the AI agent implementations
- **AgentCore_Entrypoint**: A routing layer within AgentCore_Gateway that receives incoming requests and dispatches them to the selected AI agent implementation (LangGraph_Workflow or Strands_Agents) based on the agent selection parameter
- **AgentCore_Gateway**: AWS AgentCore Gateway service managing API access to the AI agent implementations
- **EKS_Cluster**: Amazon Elastic Kubernetes Service cluster for deploying long-running traditional microservices (API Gateway Service, Cat Profile Service) that benefit from Kubernetes orchestration
- **ECS_Fargate**: Amazon ECS with Fargate launch type for deploying containerized applications that don't need Kubernetes complexity — Django microservices (Feeding Service, Health Monitor Service) and React frontends (Device Simulator, Admin Console)
- **Lambda**: AWS Lambda for deploying lightweight, event-driven microservices. The Device Service (Go) runs on Lambda, leveraging Go's fast cold starts and Lambda's request-response model for device command and telemetry operations
- **PostgreSQL**: PostgreSQL relational database for storing structured domain data (cat profiles, feeding schedules, device registry) with proper relational modeling, foreign keys, and queryable schemas. Each service uses a language-idiomatic ORM or driver: Spring Boot uses JPA/Hibernate, Django uses Django ORM, and Go uses pgx or sqlx.
- **DynamoDB**: Amazon DynamoDB tables for storing high-throughput, time-series, and event-driven data (device telemetry, device commands, feeding events, health metrics, health alerts). Chosen for its append-only write patterns, TTL-based cleanup, and time-series query efficiency.
- **Dual_Storage_Pattern**: An architectural pattern where each microservice uses PostgreSQL for structured domain/relational data and DynamoDB for operational event and time-series data, demonstrating how to combine an RDBMS with a NoSQL store in a polyglot persistence architecture.
- **Cognito**: Amazon Cognito user pool providing authentication and authorization for all frontend applications
- **ECR**: Amazon Elastic Container Registry for storing container images
- **Sibling_Infra_Project**: An external Terraform project (in a sibling repository) that provisions the shared VPC, EKS_Cluster, ECS_Fargate, Lambda function infrastructure, OIDC provider, and IAM roles; this application project assumes all infrastructure already exists
- **CI_CD_Pipeline**: GitHub Actions workflows automating build, test, packaging, and deployment of all application components to pre-existing infrastructure
- **Deployment_Artifacts**: Dockerfiles, Kubernetes manifests, Helm charts, ECS task definitions, and Lambda deployment packages produced by this project for deployment to infrastructure provisioned by the Sibling_Infra_Project

## Requirements

### Requirement 1: LangGraph Workflow Agent with Embedded RAG

**User Story:** As a cat owner, I want a stateful workflow agent that orchestrates cat-care operations step by step and retrieves domain knowledge when needed, so that I receive accurate, context-aware responses to my natural language commands.

#### Acceptance Criteria

1. THE LangGraph_Workflow SHALL define a directed graph with distinct nodes for: intent understanding, cat profile lookup, action determination, device command execution, RAG-based knowledge retrieval, and response generation
2. WHEN a user sends a natural language command via the Chatbot_UI, THE LangGraph_Workflow SHALL route the command through the graph nodes in the correct order and return a response within 5 seconds
3. WHEN the LangGraph_Workflow determines that domain knowledge is needed (e.g., feeding guidelines, breed-specific dietary needs, health tips), THE LangGraph_Workflow SHALL invoke the RAG_Node to retrieve relevant context from the Knowledge_Base before generating a response
4. WHEN the LangGraph_Workflow receives a command to control a device, THE LangGraph_Workflow SHALL invoke the Device_Service via the device command execution node to perform the action
5. WHEN the LangGraph_Workflow receives a query about a cat's status, THE LangGraph_Workflow SHALL retrieve data from the Cat_Profile_Service and Health_Monitor_Service via the cat profile lookup node and return a natural language summary
6. WHEN the LangGraph_Workflow cannot determine the user's intent, THE LangGraph_Workflow SHALL transition to a clarification node and ask the user a clarifying question rather than executing an incorrect action
7. THE LangGraph_Workflow SHALL maintain conversation state across graph transitions so that prior context is available to subsequent nodes within the same session
8. THE LangGraph_Workflow SHALL be deployed to the AgentCore_Runtime alongside the Strands_Agents system
9. THE LangGraph_Workflow SHALL be implemented in Python and reside in the `langgraph-agent/` directory

### Requirement 2: Simple RAG Pipeline within LangGraph Workflow

**User Story:** As a cat owner, I want the AI agent to draw on cat-care domain knowledge when answering my questions, so that I receive accurate advice about feeding, health, and device usage.

#### Acceptance Criteria

1. THE RAG_Pipeline SHALL index cat-care domain documents (feeding guidelines, breed-specific dietary needs, health monitoring tips, device troubleshooting guides) into the Vector_Store
2. WHEN the RAG_Node is invoked by the LangGraph_Workflow, THE RAG_Pipeline SHALL convert the query into an embedding, retrieve the top-k most relevant document chunks from the Vector_Store, and return the retrieved context to the LangGraph_Workflow
3. THE RAG_Pipeline SHALL use a vector database (FAISS or ChromaDB for local development, or Amazon Bedrock Knowledge Base for cloud deployment) as the Vector_Store
4. WHEN the RAG_Pipeline retrieves context, THE LangGraph_Workflow SHALL include the retrieved context in the prompt sent to the language model for response generation
5. THE Knowledge_Base SHALL contain at least the following document categories: feeding guidelines, breed-specific dietary needs, health monitoring tips, and device troubleshooting guides
6. WHEN a new document is added to the Knowledge_Base, THE RAG_Pipeline SHALL re-index the document into the Vector_Store without requiring a full re-index of existing documents
7. THE RAG_Pipeline SHALL be implemented in Python and reside in the `langgraph-agent/rag/` directory

### Requirement 3: Strands Multi-Agent System

**User Story:** As a cat owner, I want a multi-agent system where specialized agents handle different aspects of cat care, so that each domain (feeding, health, device control) is managed by a focused expert agent.

#### Acceptance Criteria

1. THE Strands_Agents system SHALL consist of at least four agents: an Orchestrator_Agent, a Feeding_Agent, a Health_Agent, and a Device_Agent
2. WHEN a user sends a natural language command via the Chatbot_UI, THE Orchestrator_Agent SHALL classify the intent and route the request to the appropriate specialist agent (Feeding_Agent, Health_Agent, or Device_Agent)
3. WHEN the Feeding_Agent receives a feeding-related request, THE Feeding_Agent SHALL interact with the Feeding_Service and Device_Service to execute feeding commands or retrieve feeding history
4. WHEN the Health_Agent receives a health-related query, THE Health_Agent SHALL interact with the Health_Monitor_Service and Cat_Profile_Service to retrieve health summaries and alerts
5. WHEN the Device_Agent receives a device control or status request, THE Device_Agent SHALL interact with the Device_Service to execute device commands or retrieve device state
6. WHEN a request spans multiple domains (e.g., "feed my cat and check her health"), THE Orchestrator_Agent SHALL coordinate across the relevant specialist agents and combine the results into a single response
7. THE Strands_Agents system SHALL be deployed to the AgentCore_Runtime
8. THE AgentCore_Memory SHALL persist conversation history so that the Orchestrator_Agent can reference prior interactions within the same session
9. IF a specialist agent encounters an error, THEN THE Orchestrator_Agent SHALL return a descriptive error message to the user rather than failing silently
10. THE Strands_Agents system SHALL be implemented in Python using the Strands Agents SDK and reside in the `strands-agents/` directory

### Requirement 4: AgentCore Entrypoint and Agent Routing

**User Story:** As a platform engineer, I want a unified entrypoint in AgentCore that routes requests to the correct AI agent implementation, so that all agent deployments are managed under a single gateway.

#### Acceptance Criteria

1. THE AgentCore_Entrypoint SHALL accept incoming requests from the Chatbot_UI via the AgentCore_Gateway and route each request to the selected AI agent implementation (LangGraph_Workflow or Strands_Agents) based on an agent selection parameter in the request
2. WHEN a request specifies the LangGraph_Workflow agent, THE AgentCore_Entrypoint SHALL forward the request to the LangGraph_Workflow instance running on AgentCore_Runtime
3. WHEN a request specifies the Strands_Agents agent, THE AgentCore_Entrypoint SHALL forward the request to the Orchestrator_Agent of the Strands_Agents system running on AgentCore_Runtime
4. IF a request does not include a valid agent selection parameter, THEN THE AgentCore_Entrypoint SHALL return a descriptive error indicating the available agent implementations
5. THE AgentCore_Entrypoint SHALL pass authentication context (JWT tokens from Cognito) through to the selected AI agent implementation without modification

### Requirement 5: Cat-Care IoT Device Simulation

**User Story:** As a developer, I want a device simulator that emulates cat-care IoT devices, so that I can test the system end-to-end without physical hardware.

#### Acceptance Criteria

1. THE Device_Simulator SHALL simulate at least four cat-care device types: automatic feeder, water fountain, litter box monitor, and activity tracker
2. WHEN a simulated device receives a command (via polling the Device_Service for desired state changes), THE Device_Simulator SHALL update the device's reported state via a REST call to the Device_Service within 2 seconds
3. THE Device_Simulator SHALL publish periodic telemetry data (feeding events, water level, litter box usage, activity metrics) by writing telemetry records to the Device_Service via REST at configurable intervals
4. WHEN a simulated device encounters a fault condition (low food, empty water, full litter box), THE Device_Simulator SHALL create an alert record via a REST call to the Health_Monitor_Service
5. THE Device_Simulator SHALL display real-time device states and telemetry in its React UI

### Requirement 6: Chatbot Frontend

**User Story:** As a cat owner, I want a web-based chat interface to communicate with the AI agent, so that I can manage my cats' care from any browser.

#### Acceptance Criteria

1. THE Chatbot_UI SHALL provide a text-based conversational interface for sending messages to any of the AI agent implementations (LangGraph_Workflow or Strands_Agents) via AgentCore_Gateway
2. WHEN a user opens the Chatbot_UI, THE Cognito SHALL authenticate the user before granting access to the AI agent implementations
3. THE Chatbot_UI SHALL display the AI agent's responses in real time as streamed text
4. THE Chatbot_UI SHALL maintain a scrollable conversation history for the current session
5. WHEN the Chatbot_UI loses connectivity to the AgentCore_Gateway, THE Chatbot_UI SHALL display a connection status indicator and retry automatically
6. THE Chatbot_UI SHALL allow the user to select which AI agent implementation (LangGraph_Workflow or Strands_Agents) to interact with via a toggle or dropdown control

### Requirement 7: Admin Console

**User Story:** As a system administrator, I want an admin console to manage devices, view system health, and configure the AI agent implementations, so that I can maintain and monitor the platform.

#### Acceptance Criteria

1. THE Admin_Console SHALL display a dashboard showing the status of all registered IoT devices, microservices, and the AI agent implementations (LangGraph_Workflow and Strands_Agents)
2. WHEN an administrator registers a new device via the Admin_Console, THE Device_Service SHALL persist the device configuration and add the device to the device registry
3. THE Admin_Console SHALL allow administrators to view and edit cat profiles managed by the Cat_Profile_Service
4. WHEN an administrator views a specific device, THE Admin_Console SHALL display the device's current state, recent telemetry, and command history
5. WHEN a user opens the Admin_Console, THE Cognito SHALL authenticate the user and verify the user holds an administrator role before granting access

### Requirement 8: Cat Profile and Health Management

**User Story:** As a cat owner, I want to manage profiles for each of my cats with health tracking, so that the system can personalize feeding schedules and monitor well-being.

#### Acceptance Criteria

1. THE Cat_Profile_Service SHALL store cat profiles containing name, breed, age, weight, dietary restrictions, and owner association in PostgreSQL using JPA/Hibernate, with proper relational modeling (foreign keys for owner and device associations)
2. WHEN a cat profile is created or updated, THE Cat_Profile_Service SHALL validate that required fields (name, weight) are present and return a descriptive error for missing fields
3. THE Health_Monitor_Service SHALL aggregate telemetry data from DynamoDB (written by IoT devices via the Device_Service) into per-cat health summaries
4. WHEN a health metric for a cat deviates beyond a configurable threshold, THE Health_Monitor_Service SHALL generate a health alert in DynamoDB and make the alert available to the AI agent implementations
5. THE Feeding_Service SHALL manage per-cat feeding schedules with configurable meal times, portion sizes, and dietary constraints, stored in PostgreSQL using Django ORM for structured time-based querying

### Requirement 9: Feeding Schedule Automation

**User Story:** As a cat owner, I want automated feeding schedules for each cat, so that my cats are fed consistently even when I am not available to issue commands.

#### Acceptance Criteria

1. WHEN a scheduled feeding time arrives, THE Feeding_Service SHALL invoke the Device_Service to update the target feeder device's desired state to dispense the configured portion for the target cat
2. WHEN a feeding event completes, THE Feeding_Service SHALL record the event (timestamp, portion dispensed, device used, cat fed) as an append-only entry in DynamoDB
3. IF a feeder device is offline at the scheduled feeding time, THEN THE Feeding_Service SHALL retry the command three times at 60-second intervals and generate an alert after all retries fail
4. WHEN a user requests feeding history via any AI agent implementation, THE Feeding_Service SHALL return the feeding log for the specified cat and time range from DynamoDB
5. THE Feeding_Service SHALL prevent duplicate feedings by rejecting a feed command for a cat that was fed within a configurable minimum interval

### Requirement 10: Device Command and State Management

**User Story:** As a developer, I want reliable REST-based device command and state management, so that commands and telemetry flow correctly through the system via the Device Service.

#### Acceptance Criteria

1. THE Device_Service SHALL accept device commands via a REST endpoint (`POST /devices/{device_id}/commands`) that updates the device's desired state in DynamoDB
2. THE Device_Service SHALL track command status (pending, acknowledged, timed_out, failed) and allow callers to query command status via `GET /devices/{device_id}/commands/{command_id}`
3. WHEN the Device_Simulator polls for pending commands and acknowledges a command, THE Device_Service SHALL update the command status to acknowledged and update the device's reported state
4. IF a device command is not acknowledged within 10 seconds, THEN THE Device_Service SHALL mark the command as timed out and return a failure status to the calling service
5. THE Device_Service SHALL maintain a device shadow (desired and reported state) for each registered device in DynamoDB

### Requirement 11: Polyglot Microservices across Multi-Compute

**User Story:** As a platform engineer, I want the backend microservices (built with Spring Boot, Python Django, and Go) distributed across EKS, ECS Fargate, and Lambda with proper service discovery and observability, so that the system demonstrates polyglot service development across multiple compute types.

#### Acceptance Criteria

1. THE API_Gateway_Service (Spring Boot) and Cat_Profile_Service (Spring Boot) SHALL be deployed as separate Kubernetes deployments on the EKS_Cluster; THE Feeding_Service (Django) and Health_Monitor_Service (Django) SHALL be deployed as ECS Fargate services; THE Device_Service (Go) SHALL be deployed as a Lambda function
2. THE microservices SHALL use appropriate service discovery mechanisms for cross-compute communication: Kubernetes Service resources within EKS, ECS service discovery for Fargate services, and API Gateway or direct invocation for Lambda functions
3. WHEN a microservice instance becomes unhealthy, THE respective compute platform SHALL restart the instance: Kubernetes liveness/readiness probes for EKS, ECS health checks for Fargate, and Lambda's built-in retry/concurrency management for Lambda functions
4. THE API_Gateway_Service SHALL route incoming requests to the appropriate downstream microservice based on the request path, regardless of which compute platform hosts the target service
5. WHEN a request passes through the API_Gateway_Service, THE API_Gateway_Service SHALL propagate trace context headers for distributed tracing across all compute types

### Requirement 12: Multi-Compute Deployment Packaging

**User Story:** As a platform engineer, I want the application code packaged with deployment artifacts for EKS and ECS Fargate, so that the application can be deployed to the multi-compute infrastructure provisioned by the sibling infrastructure project.

#### Acceptance Criteria

1. THE project SHALL include a Dockerfile for each microservice (Spring Boot, Django, and Go) that produces a container image deployable to the target compute platform (EKS for Spring Boot, ECS Fargate for Django)
2. THE project SHALL include Kubernetes manifests (or Helm charts) in a `deploy/k8s/` directory defining deployments, services, config maps, and ingress resources for the Spring Boot microservices (API Gateway, Cat Profile) on the EKS_Cluster
3. THE project SHALL include ECS task definitions and service configurations in a `deploy/ecs/` directory for deploying the Chatbot_UI, Admin_Console, Device_Simulator, Feeding_Service, and Health_Monitor_Service as containerized applications on ECS_Fargate
4. THE project SHALL include a Lambda deployment package configuration in a `deploy/lambda/` directory for deploying the Device_Service (Go) as a Lambda function, including the compiled Go binary and any required configuration
5. THE project SHALL include a Dockerfile for each React frontend application that produces a container image serving the built static assets via a lightweight web server
6. THE project SHALL include deployment configuration files in a `deploy/` directory that reference environment-specific variables (e.g., VPC IDs, subnet IDs, cluster names, Lambda function names) provided by the Sibling_Infra_Project
7. WHEN a developer builds any component, THE build system SHALL produce a deployment-ready artifact (container image for EKS/ECS, compiled binary for Lambda) without requiring access to the infrastructure provisioning tools

### Requirement 13: GitHub Actions CI/CD with OIDC Authentication

**User Story:** As a developer, I want automated CI/CD pipelines using GitHub Actions with OIDC authentication to AWS, so that deployments are secure and require no long-lived credentials.

#### Acceptance Criteria

1. THE CI_CD_Pipeline SHALL authenticate to AWS using OIDC (with the IAM OIDC provider and roles provisioned by the Sibling_Infra_Project) without storing AWS access keys in GitHub secrets
2. WHEN code is pushed to the main branch, THE CI_CD_Pipeline SHALL build container images, push the images to ECR, and deploy updated services to the target compute platform
3. THE CI_CD_Pipeline SHALL include separate workflow jobs for building and deploying each component (microservices, frontends, AI agent implementations)
4. WHEN a pull request is opened, THE CI_CD_Pipeline SHALL run linting, unit tests, and build validation without deploying
5. THE CI_CD_Pipeline SHALL use environment-specific configuration (e.g., dev, prod) to determine target ECR repositories, EKS cluster names, and ECS service identifiers provided by the Sibling_Infra_Project

### Requirement 14: Authentication and Authorization

**User Story:** As a system administrator, I want centralized authentication via Cognito for all frontends, so that user access is secure and manageable.

#### Acceptance Criteria

1. THE Cognito SHALL provide a user pool with sign-up, sign-in, and password recovery flows for the Chatbot_UI, Device_Simulator, and Admin_Console
2. WHEN a user authenticates via Cognito, THE Cognito SHALL issue JWT tokens that the frontend applications include in API requests to the AgentCore_Gateway and API_Gateway_Service
3. THE Cognito SHALL define at least two user groups: "owner" for cat owners and "admin" for administrators
4. WHEN a request arrives without a valid JWT token, THE API_Gateway_Service SHALL reject the request with a 401 Unauthorized response
5. WHEN a non-admin user attempts to access an admin-only endpoint, THE API_Gateway_Service SHALL reject the request with a 403 Forbidden response

### Requirement 15: Project Structure and Developer Experience

**User Story:** As a developer, I want a well-organized monorepo with clear project structure, so that I can navigate, build, and contribute to the project efficiently.

#### Acceptance Criteria

1. THE project SHALL organize code into top-level directories: `langgraph-agent/` (LangGraph_Workflow with embedded RAG), `strands-agents/` (Strands_Agents multi-agent system), `chatbot/` (Chatbot_UI), `device-simulator/` (Device_Simulator), `admin-console/` (Admin_Console), `services/` (Spring Boot, Django, and Go microservices), `deploy/` (Kubernetes manifests, ECS task definitions, Lambda deployment packages), and `.github/workflows/` (CI_CD_Pipeline)
2. THE project SHALL include a root README documenting the architecture, the three AI agent patterns (LangGraph workflow in `langgraph-agent/`, Strands multi-agent in `strands-agents/`, RAG pipeline in `langgraph-agent/rag/`), project structure, prerequisites, and getting-started instructions
3. WHEN a developer runs the build command for any component, THE build system SHALL produce a deployable artifact (container image) without requiring manual configuration steps beyond setting environment variables
4. THE project SHALL include a local development setup using Docker Compose that starts the microservices, the Device_Simulator, the LangGraph_Workflow agent, the Strands_Agents system, a PostgreSQL instance, and a DynamoDB Local instance for end-to-end testing without AWS dependencies
5. WHEN a developer adds a new microservice, THE project SHALL provide a service template or archetype in the `services/` directory to ensure consistent structure
