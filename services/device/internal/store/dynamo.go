package store

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/model"
)

const (
	shadowsTable  = "device-shadows"
	commandsTable = "device-commands"
)

// DynamoStore handles DynamoDB operations for device shadows and commands.
type DynamoStore struct {
	client *dynamodb.Client
}

// NewDynamoStore creates a new DynamoStore with the given DynamoDB client.
func NewDynamoStore(client *dynamodb.Client) *DynamoStore {
	return &DynamoStore{client: client}
}

// GetShadow retrieves the device shadow for a given device ID.
func (s *DynamoStore) GetShadow(ctx context.Context, deviceID string) (*model.DeviceShadow, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(shadowsTable),
		Key: map[string]types.AttributeValue{
			"device_id": &types.AttributeValueMemberS{Value: deviceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get shadow: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}

	var shadow model.DeviceShadow
	if err := attributevalue.UnmarshalMap(out.Item, &shadow); err != nil {
		return nil, fmt.Errorf("unmarshal shadow: %w", err)
	}
	return &shadow, nil
}

// PutShadow creates or replaces the device shadow for a given device.
func (s *DynamoStore) PutShadow(ctx context.Context, shadow *model.DeviceShadow) error {
	item, err := attributevalue.MarshalMap(shadow)
	if err != nil {
		return fmt.Errorf("marshal shadow: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(shadowsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put shadow: %w", err)
	}
	return nil
}

// UpdateDesiredState updates only the desired_state and updated_at fields of a shadow.
func (s *DynamoStore) UpdateDesiredState(ctx context.Context, deviceID string, desiredState map[string]interface{}) error {
	desiredAV, err := attributevalue.MarshalMap(desiredState)
	if err != nil {
		return fmt.Errorf("marshal desired state: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(shadowsTable),
		Key: map[string]types.AttributeValue{
			"device_id": &types.AttributeValueMemberS{Value: deviceID},
		},
		UpdateExpression: aws.String("SET desired_state = :ds, updated_at = :ua"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ds": &types.AttributeValueMemberM{Value: desiredAV},
			":ua": &types.AttributeValueMemberS{Value: now},
		},
	})
	if err != nil {
		return fmt.Errorf("update desired state: %w", err)
	}
	return nil
}

// UpdateReportedState updates only the reported_state and updated_at fields of a shadow.
func (s *DynamoStore) UpdateReportedState(ctx context.Context, deviceID string, reportedState map[string]interface{}) error {
	reportedAV, err := attributevalue.MarshalMap(reportedState)
	if err != nil {
		return fmt.Errorf("marshal reported state: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(shadowsTable),
		Key: map[string]types.AttributeValue{
			"device_id": &types.AttributeValueMemberS{Value: deviceID},
		},
		UpdateExpression: aws.String("SET reported_state = :rs, updated_at = :ua"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":rs": &types.AttributeValueMemberM{Value: reportedAV},
			":ua": &types.AttributeValueMemberS{Value: now},
		},
	})
	if err != nil {
		return fmt.Errorf("update reported state: %w", err)
	}
	return nil
}

// CreateCommand inserts a new command record into DynamoDB.
func (s *DynamoStore) CreateCommand(ctx context.Context, cmd *model.CommandRecord) error {
	item, err := attributevalue.MarshalMap(cmd)
	if err != nil {
		return fmt.Errorf("marshal command: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(commandsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("put command: %w", err)
	}
	return nil
}

// GetCommand retrieves a command record by command ID.
func (s *DynamoStore) GetCommand(ctx context.Context, commandID string) (*model.CommandRecord, error) {
	out, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(commandsTable),
		Key: map[string]types.AttributeValue{
			"command_id": &types.AttributeValueMemberS{Value: commandID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get command: %w", err)
	}
	if out.Item == nil {
		return nil, nil
	}

	var cmd model.CommandRecord
	if err := attributevalue.UnmarshalMap(out.Item, &cmd); err != nil {
		return nil, fmt.Errorf("unmarshal command: %w", err)
	}
	return &cmd, nil
}

// ListCommands retrieves commands for a device, optionally filtered by status.
func (s *DynamoStore) ListCommands(ctx context.Context, deviceID string, status string) ([]model.CommandRecord, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(commandsTable),
		IndexName:             aws.String("device_id-index"),
		KeyConditionExpression: aws.String("device_id = :did"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":did": &types.AttributeValueMemberS{Value: deviceID},
		},
	}

	if status != "" {
		input.FilterExpression = aws.String("#s = :status")
		input.ExpressionAttributeNames = map[string]string{
			"#s": "status",
		}
		input.ExpressionAttributeValues[":status"] = &types.AttributeValueMemberS{Value: status}
	}

	out, err := s.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("query commands: %w", err)
	}

	var commands []model.CommandRecord
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &commands); err != nil {
		return nil, fmt.Errorf("unmarshal commands: %w", err)
	}
	if commands == nil {
		commands = []model.CommandRecord{}
	}
	return commands, nil
}

// AcknowledgeCommand updates a command's status to acknowledged and sets acked_at.
func (s *DynamoStore) AcknowledgeCommand(ctx context.Context, commandID string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(commandsTable),
		Key: map[string]types.AttributeValue{
			"command_id": &types.AttributeValueMemberS{Value: commandID},
		},
		UpdateExpression: aws.String("SET #s = :status, acked_at = :acked"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "acknowledged"},
			":acked":  &types.AttributeValueMemberS{Value: now},
		},
		ConditionExpression: aws.String("#s = :pending"),
	})
	if err != nil {
		return fmt.Errorf("acknowledge command: %w", err)
	}
	return nil
}

// TimeoutPendingCommands marks commands as timed_out if they've been pending longer than the timeout duration.
func (s *DynamoStore) TimeoutPendingCommands(ctx context.Context, deviceID string, timeout time.Duration) error {
	commands, err := s.ListCommands(ctx, deviceID, "pending")
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, cmd := range commands {
		createdAt, err := time.Parse(time.RFC3339, cmd.CreatedAt)
		if err != nil {
			continue
		}
		if now.Sub(createdAt) > timeout {
			_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				TableName: aws.String(commandsTable),
				Key: map[string]types.AttributeValue{
					"command_id": &types.AttributeValueMemberS{Value: cmd.CommandID},
				},
				UpdateExpression: aws.String("SET #s = :status"),
				ExpressionAttributeNames: map[string]string{
					"#s": "status",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":status": &types.AttributeValueMemberS{Value: "timed_out"},
				},
			})
			if err != nil {
				return fmt.Errorf("timeout command %s: %w", cmd.CommandID, err)
			}
		}
	}
	return nil
}
