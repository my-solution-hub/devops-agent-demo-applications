package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/router"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/store"
)

func main() {
	ctx := context.Background()

	// PostgreSQL connection
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	// AWS SDK config
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// DynamoDB client (with optional local endpoint)
	dynamoOpts := func(o *dynamodb.Options) {}
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		dynamoOpts = func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		}
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg, dynamoOpts)

	// Create stores
	pgStore := store.NewPostgresStore(pool)
	dynamoStore := store.NewDynamoStore(dynamoClient)

	// Create router
	r := router.New(pgStore, dynamoStore)

	// Start background timeout goroutine
	// In production Lambda, this would be a separate scheduled invocation.
	go runTimeoutChecker(ctx, dynamoStore, pgStore)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Device Service starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// runTimeoutChecker periodically checks for pending commands that have exceeded
// the 10-second timeout and marks them as timed_out.
// Note: In Lambda deployment, this logic would be handled by a separate scheduled
// Lambda invocation (e.g., via EventBridge). This goroutine is only for local dev.
func runTimeoutChecker(ctx context.Context, dynamo *store.DynamoStore, pg *store.PostgresStore) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			devices, err := pg.ListDevices(ctx)
			if err != nil {
				log.Printf("timeout checker: failed to list devices: %v", err)
				continue
			}
			for _, device := range devices {
				if err := dynamo.TimeoutPendingCommands(ctx, device.DeviceID, 10*time.Second); err != nil {
					log.Printf("timeout checker: failed for device %s: %v", device.DeviceID, err)
				}
			}
		}
	}
}
