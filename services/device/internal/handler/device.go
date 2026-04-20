package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/model"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/store"
)

// DeviceHandler handles device CRUD and shadow endpoints.
type DeviceHandler struct {
	pg    *store.PostgresStore
	dynamo *store.DynamoStore
}

// NewDeviceHandler creates a new DeviceHandler.
func NewDeviceHandler(pg *store.PostgresStore, dynamo *store.DynamoStore) *DeviceHandler {
	return &DeviceHandler{pg: pg, dynamo: dynamo}
}

// CreateDevice handles POST /devices.
func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceType string                 `json:"device_type"`
		Name       string                 `json:"name"`
		Config     map[string]interface{} `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DeviceType == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "device_type and name are required")
		return
	}

	device := &model.Device{
		DeviceID:   uuid.New().String(),
		DeviceType: req.DeviceType,
		Name:       req.Name,
		Status:     "offline",
		Config:     req.Config,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	if device.Config == nil {
		device.Config = make(map[string]interface{})
	}

	if err := h.pg.CreateDevice(r.Context(), device); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create device")
		return
	}

	// Initialize device shadow in DynamoDB
	shadow := &model.DeviceShadow{
		DeviceID:      device.DeviceID,
		DesiredState:  make(map[string]interface{}),
		ReportedState: make(map[string]interface{}),
		UpdatedAt:     device.CreatedAt,
	}
	if err := h.dynamo.PutShadow(r.Context(), shadow); err != nil {
		// Log but don't fail the request — shadow can be created lazily
		_ = err
	}

	writeJSON(w, http.StatusCreated, device)
}

// GetDevice handles GET /devices/{id}.
func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	device, err := h.pg.GetDevice(r.Context(), deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get device")
		return
	}
	if device == nil {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	// Merge shadow state from DynamoDB
	shadow, err := h.dynamo.GetShadow(r.Context(), deviceID)
	if err != nil {
		// Return device without shadow on DynamoDB error
		writeJSON(w, http.StatusOK, device)
		return
	}

	// Return device with shadow merged
	response := struct {
		*model.Device
		Shadow *model.DeviceShadow `json:"shadow,omitempty"`
	}{
		Device: device,
		Shadow: shadow,
	}

	writeJSON(w, http.StatusOK, response)
}

// UpdateDevice handles PUT /devices/{id}.
func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	// Check device exists
	existing, err := h.pg.GetDevice(r.Context(), deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get device")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	var req struct {
		DeviceType *string                `json:"device_type"`
		Name       *string                `json:"name"`
		Status     *string                `json:"status"`
		Config     map[string]interface{} `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Apply partial updates
	if req.DeviceType != nil {
		existing.DeviceType = *req.DeviceType
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Config != nil {
		existing.Config = req.Config
	}

	if err := h.pg.UpdateDevice(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update device")
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

// ListDevices handles GET /devices.
func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.pg.ListDevices(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}

	writeJSON(w, http.StatusOK, devices)
}

// GetShadow handles GET /devices/{id}/shadow.
func (h *DeviceHandler) GetShadow(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	shadow, err := h.dynamo.GetShadow(r.Context(), deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get shadow")
		return
	}
	if shadow == nil {
		writeError(w, http.StatusNotFound, "device shadow not found")
		return
	}

	writeJSON(w, http.StatusOK, shadow)
}

// UpdateShadow handles PUT /devices/{id}/shadow.
func (h *DeviceHandler) UpdateShadow(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	var req struct {
		DesiredState map[string]interface{} `json:"desired_state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DesiredState == nil {
		writeError(w, http.StatusBadRequest, "desired_state is required")
		return
	}

	if err := h.dynamo.UpdateDesiredState(r.Context(), deviceID, req.DesiredState); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update shadow")
		return
	}

	// Return updated shadow
	shadow, err := h.dynamo.GetShadow(r.Context(), deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get updated shadow")
		return
	}

	writeJSON(w, http.StatusOK, shadow)
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
