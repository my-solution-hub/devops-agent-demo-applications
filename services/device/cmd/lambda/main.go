package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
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

	// DynamoDB client
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

	// Create Lambda adapter
	adapter := chiadapter.New(r)

	// Start Lambda handler
	lambda.Start(adapter.ProxyWithContext)
}
