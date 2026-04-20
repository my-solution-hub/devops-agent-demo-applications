#!/usr/bin/env bash
# Smart Home Cat Demo — DynamoDB Local table creation script
# Creates all DynamoDB tables used by the application services.

set -euo pipefail

ENDPOINT="http://dynamodb-local:8000"
REGION="us-east-1"

aws_cmd() {
  aws dynamodb "$@" --endpoint-url "$ENDPOINT" --region "$REGION"
}

echo "Waiting for DynamoDB Local to be ready..."
until aws_cmd list-tables > /dev/null 2>&1; do
  sleep 1
done
echo "DynamoDB Local is ready."

# ── device-shadows (PK: device_id) ──────────────────────────────────────────
echo "Creating table: device-shadows"
aws_cmd create-table \
  --table-name device-shadows \
  --attribute-definitions \
    AttributeName=device_id,AttributeType=S \
  --key-schema \
    AttributeName=device_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "  table device-shadows already exists"

# ── device-commands (PK: command_id, GSI: device_id) ────────────────────────
echo "Creating table: device-commands"
aws_cmd create-table \
  --table-name device-commands \
  --attribute-definitions \
    AttributeName=command_id,AttributeType=S \
    AttributeName=device_id,AttributeType=S \
  --key-schema \
    AttributeName=command_id,KeyType=HASH \
  --global-secondary-indexes \
    '[{
      "IndexName": "device_id-index",
      "KeySchema": [{"AttributeName":"device_id","KeyType":"HASH"}],
      "Projection": {"ProjectionType":"ALL"}
    }]' \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "  table device-commands already exists"

# ── device-telemetry (PK: device_id, SK: timestamp) ────────────────────────
echo "Creating table: device-telemetry"
aws_cmd create-table \
  --table-name device-telemetry \
  --attribute-definitions \
    AttributeName=device_id,AttributeType=S \
    AttributeName=timestamp,AttributeType=S \
  --key-schema \
    AttributeName=device_id,KeyType=HASH \
    AttributeName=timestamp,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "  table device-telemetry already exists"

# ── feeding-events (PK: event_id, GSI: cat_id + timestamp SK) ──────────────
echo "Creating table: feeding-events"
aws_cmd create-table \
  --table-name feeding-events \
  --attribute-definitions \
    AttributeName=event_id,AttributeType=S \
    AttributeName=cat_id,AttributeType=S \
    AttributeName=timestamp,AttributeType=S \
  --key-schema \
    AttributeName=event_id,KeyType=HASH \
  --global-secondary-indexes \
    '[{
      "IndexName": "cat_id-timestamp-index",
      "KeySchema": [
        {"AttributeName":"cat_id","KeyType":"HASH"},
        {"AttributeName":"timestamp","KeyType":"RANGE"}
      ],
      "Projection": {"ProjectionType":"ALL"}
    }]' \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "  table feeding-events already exists"

# ── health-metrics (PK: cat_id, SK: timestamp) ─────────────────────────────
echo "Creating table: health-metrics"
aws_cmd create-table \
  --table-name health-metrics \
  --attribute-definitions \
    AttributeName=cat_id,AttributeType=S \
    AttributeName=timestamp,AttributeType=S \
  --key-schema \
    AttributeName=cat_id,KeyType=HASH \
    AttributeName=timestamp,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "  table health-metrics already exists"

# ── health-alerts (PK: alert_id, GSI: cat_id) ──────────────────────────────
echo "Creating table: health-alerts"
aws_cmd create-table \
  --table-name health-alerts \
  --attribute-definitions \
    AttributeName=alert_id,AttributeType=S \
    AttributeName=cat_id,AttributeType=S \
  --key-schema \
    AttributeName=alert_id,KeyType=HASH \
  --global-secondary-indexes \
    '[{
      "IndexName": "cat_id-index",
      "KeySchema": [{"AttributeName":"cat_id","KeyType":"HASH"}],
      "Projection": {"ProjectionType":"ALL"}
    }]' \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "  table health-alerts already exists"

echo "All DynamoDB tables created successfully."
