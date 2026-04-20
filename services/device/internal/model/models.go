package model

// Device represents a registered IoT device stored in PostgreSQL.
type Device struct {
	DeviceID   string                 `json:"device_id"`
	DeviceType string                 `json:"device_type"`
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	Config     map[string]interface{} `json:"config"`
	LastSeen   *string                `json:"last_seen,omitempty"`
	CreatedAt  string                 `json:"created_at"`
}

// DeviceShadow represents the desired and reported state of a device stored in DynamoDB.
type DeviceShadow struct {
	DeviceID      string                 `json:"device_id"`
	DesiredState  map[string]interface{} `json:"desired_state"`
	ReportedState map[string]interface{} `json:"reported_state"`
	UpdatedAt     string                 `json:"updated_at"`
}

// CommandRecord represents a command issued to a device stored in DynamoDB.
type CommandRecord struct {
	CommandID string            `json:"command_id"`
	DeviceID  string            `json:"device_id"`
	Action    string            `json:"action"`
	Params    map[string]string `json:"params"`
	Status    string            `json:"status"` // pending, acknowledged, timed_out, failed
	CreatedAt string            `json:"created_at"`
	AckedAt   string            `json:"acked_at,omitempty"`
}

// TelemetryRecord represents telemetry data from a device stored in DynamoDB.
type TelemetryRecord struct {
	DeviceID   string                 `json:"device_id"`
	Timestamp  string                 `json:"timestamp"`
	DeviceType string                 `json:"device_type"`
	Metrics    map[string]interface{} `json:"metrics"`
}
