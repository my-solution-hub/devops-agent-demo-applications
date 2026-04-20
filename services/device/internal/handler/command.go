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

// CommandHandler handles device command endpoints.
type CommandHandler struct {
	dynamo *store.DynamoStore
}

// NewCommandHandler creates a new CommandHandler.
func NewCommandHandler(dynamo *store.DynamoStore) *CommandHandler {
	return &CommandHandler{dynamo: dynamo}
}

// SubmitCommand handles POST /devices/{id}/commands.
func (h *CommandHandler) SubmitCommand(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	var req struct {
		Action string            `json:"action"`
		Params map[string]string `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Action == "" {
		writeError(w, http.StatusBadRequest, "action is required")
		return
	}
	if req.Params == nil {
		req.Params = make(map[string]string)
	}

	cmd := &model.CommandRecord{
		CommandID: uuid.New().String(),
		DeviceID:  deviceID,
		Action:    req.Action,
		Params:    req.Params,
		Status:    "pending",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if err := h.dynamo.CreateCommand(r.Context(), cmd); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create command")
		return
	}

	// Update desired state to reflect the command
	desiredState := map[string]interface{}{
		"action": req.Action,
	}
	for k, v := range req.Params {
		desiredState[k] = v
	}
	if err := h.dynamo.UpdateDesiredState(r.Context(), deviceID, desiredState); err != nil {
		// Command was created, log shadow update failure but don't fail the request
		_ = err
	}

	writeJSON(w, http.StatusCreated, cmd)
}

// GetCommand handles GET /devices/{id}/commands/{cmd_id}.
func (h *CommandHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	commandID := chi.URLParam(r, "cmd_id")

	cmd, err := h.dynamo.GetCommand(r.Context(), commandID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get command")
		return
	}
	if cmd == nil {
		writeError(w, http.StatusNotFound, "command not found")
		return
	}

	writeJSON(w, http.StatusOK, cmd)
}

// ListCommands handles GET /devices/{id}/commands.
func (h *CommandHandler) ListCommands(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	status := r.URL.Query().Get("status")

	commands, err := h.dynamo.ListCommands(r.Context(), deviceID, status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list commands")
		return
	}

	writeJSON(w, http.StatusOK, commands)
}

// AcknowledgeCommand handles POST /devices/{id}/commands/{cmd_id}/ack.
func (h *CommandHandler) AcknowledgeCommand(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	commandID := chi.URLParam(r, "cmd_id")

	// Verify the command exists and belongs to this device
	cmd, err := h.dynamo.GetCommand(r.Context(), commandID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get command")
		return
	}
	if cmd == nil {
		writeError(w, http.StatusNotFound, "command not found")
		return
	}
	if cmd.DeviceID != deviceID {
		writeError(w, http.StatusBadRequest, "command does not belong to this device")
		return
	}
	if cmd.Status != "pending" {
		writeError(w, http.StatusBadRequest, "command is not in pending status")
		return
	}

	if err := h.dynamo.AcknowledgeCommand(r.Context(), commandID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to acknowledge command")
		return
	}

	// Update reported state to reflect the acknowledged command
	reportedState := map[string]interface{}{
		"action": cmd.Action,
		"status": "completed",
	}
	for k, v := range cmd.Params {
		reportedState[k] = v
	}
	if err := h.dynamo.UpdateReportedState(r.Context(), deviceID, reportedState); err != nil {
		// Ack succeeded, log shadow update failure but don't fail the request
		_ = err
	}

	// Return updated command
	updatedCmd, err := h.dynamo.GetCommand(r.Context(), commandID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get updated command")
		return
	}

	writeJSON(w, http.StatusOK, updatedCmd)
}
