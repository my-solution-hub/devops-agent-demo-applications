package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/model"
)

// PostgresStore handles device registry operations in PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

// NewPostgresStore creates a new PostgresStore with the given connection pool.
func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// CreateDevice inserts a new device into the devices table.
func (s *PostgresStore) CreateDevice(ctx context.Context, d *model.Device) error {
	configJSON, err := json.Marshal(d.Config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO devices (device_id, device_type, name, status, config, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		d.DeviceID, d.DeviceType, d.Name, d.Status, configJSON, d.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert device: %w", err)
	}
	return nil
}

// GetDevice retrieves a device by ID.
func (s *PostgresStore) GetDevice(ctx context.Context, deviceID string) (*model.Device, error) {
	var d model.Device
	var configJSON []byte
	var lastSeen *string

	err := s.pool.QueryRow(ctx,
		`SELECT device_id, device_type, name, status, config, last_seen, created_at
		 FROM devices WHERE device_id = $1`, deviceID,
	).Scan(&d.DeviceID, &d.DeviceType, &d.Name, &d.Status, &configJSON, &lastSeen, &d.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query device: %w", err)
	}

	if configJSON != nil {
		if err := json.Unmarshal(configJSON, &d.Config); err != nil {
			return nil, fmt.Errorf("unmarshal config: %w", err)
		}
	}
	if d.Config == nil {
		d.Config = make(map[string]interface{})
	}
	d.LastSeen = lastSeen
	return &d, nil
}

// UpdateDevice updates a device's name, device_type, status, and config.
func (s *PostgresStore) UpdateDevice(ctx context.Context, d *model.Device) error {
	configJSON, err := json.Marshal(d.Config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	tag, err := s.pool.Exec(ctx,
		`UPDATE devices SET device_type = $2, name = $3, status = $4, config = $5
		 WHERE device_id = $1`,
		d.DeviceID, d.DeviceType, d.Name, d.Status, configJSON,
	)
	if err != nil {
		return fmt.Errorf("update device: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("device not found")
	}
	return nil
}

// ListDevices returns all devices.
func (s *PostgresStore) ListDevices(ctx context.Context) ([]model.Device, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT device_id, device_type, name, status, config, last_seen, created_at
		 FROM devices ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query devices: %w", err)
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		var configJSON []byte
		var lastSeen *string

		if err := rows.Scan(&d.DeviceID, &d.DeviceType, &d.Name, &d.Status, &configJSON, &lastSeen, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}

		if configJSON != nil {
			if err := json.Unmarshal(configJSON, &d.Config); err != nil {
				return nil, fmt.Errorf("unmarshal config: %w", err)
			}
		}
		if d.Config == nil {
			d.Config = make(map[string]interface{})
		}
		d.LastSeen = lastSeen
		devices = append(devices, d)
	}
	if devices == nil {
		devices = []model.Device{}
	}
	return devices, nil
}
