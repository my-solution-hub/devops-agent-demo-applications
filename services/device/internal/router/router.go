package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/handler"
	"github.com/my-solution-hub/devops-agent-demo-applications/services/device/internal/store"
)

// New creates and configures the chi router with all device service routes.
func New(pg *store.PostgresStore, dynamo *store.DynamoStore) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Handlers
	deviceHandler := handler.NewDeviceHandler(pg, dynamo)
	commandHandler := handler.NewCommandHandler(dynamo)

	// Device CRUD routes
	r.Post("/devices", deviceHandler.CreateDevice)
	r.Get("/devices", deviceHandler.ListDevices)
	r.Get("/devices/{id}", deviceHandler.GetDevice)
	r.Put("/devices/{id}", deviceHandler.UpdateDevice)

	// Device shadow routes
	r.Get("/devices/{id}/shadow", deviceHandler.GetShadow)
	r.Put("/devices/{id}/shadow", deviceHandler.UpdateShadow)

	// Device command routes
	r.Post("/devices/{id}/commands", commandHandler.SubmitCommand)
	r.Get("/devices/{id}/commands", commandHandler.ListCommands)
	r.Get("/devices/{id}/commands/{cmd_id}", commandHandler.GetCommand)
	r.Post("/devices/{id}/commands/{cmd_id}/ack", commandHandler.AcknowledgeCommand)

	return r
}
