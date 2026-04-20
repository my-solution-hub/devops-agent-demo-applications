# GitHub Actions Workflows

CI/CD pipelines for the Smart Home Cat Demo. Includes:
- `ci.yml` — Pull request validation (lint, test, build)
- `deploy-services.yml` — Spring Boot services to EKS
- `deploy-ecs-services.yml` — Django services to ECS Fargate
- `deploy-frontends.yml` — React apps to ECS Fargate
- `deploy-lambda.yml` — Go Device Service to Lambda
- `deploy-agents.yml` — AI agents to AgentCore

All workflows use OIDC authentication (no stored AWS access keys).
