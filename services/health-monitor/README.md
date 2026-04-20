# Health Monitor Service

Django (Python) service for health telemetry aggregation and threshold-based alert generation. Reads telemetry from DynamoDB, computes per-cat health summaries, and generates alerts when metrics deviate beyond configurable thresholds.

Deployed to ECS Fargate.
