package handler

import (
	"net/http"

	"github.com/JerkyTreats/llm/internal/api/types"
	"github.com/JerkyTreats/llm/internal/docs"
	"github.com/JerkyTreats/llm/internal/logging"
)

// HandlerRegistry manages all HTTP handlers for the application
type HandlerRegistry struct {
	healthHandler *HealthHandler
	docsHandler   *docs.DocsHandler
	mux           *http.ServeMux
}

// NewHandlerRegistry creates a new handler registry with all handlers initialized
func NewHandlerRegistry() (*HandlerRegistry, error) {
	logging.Info("Initializing handler registry with all application handlers")

	// Initialize health handler
	healthHandler, err := NewHealthHandler()
	if err != nil {
		return nil, err
	}

	// Initialize docs handler
	docsHandler, err := docs.NewDocsHandler()
	if err != nil {
		return nil, err
	}

	registry := &HandlerRegistry{
		healthHandler: healthHandler,
		docsHandler:   docsHandler,
		mux:           http.NewServeMux(),
	}

	registry.RegisterHandlers(registry.mux)
	logging.Info("Handler registry initialized successfully with all handlers")

	return registry, nil
}

// RegisterHandlers registers all application handlers using the RouteInfo registry
func (hr *HandlerRegistry) RegisterHandlers(mux *http.ServeMux) {
	logging.Info("Registering all application handlers from RouteInfo registry")

	// Update RouteInfo registry with actual handler functions
	hr.updateRouteHandlers()

	// Register all routes from the central registry
	routes := GetRegisteredRoutes()
	for _, route := range routes {
		if route.Handler != nil {
			mux.HandleFunc(route.Path, route.Handler)
			logging.Debug("Registered %s %s from %s module", route.Method, route.Path, route.Module)
		} else {
			logging.Warn("Skipping route %s %s - handler is nil", route.Method, route.Path)
		}
	}

	logging.Info("Successfully registered %d handlers from RouteInfo registry", len(routes))
}

// GetServeMux returns the internal ServeMux with all handlers registered
func (hr *HandlerRegistry) GetServeMux() *http.ServeMux {
	return hr.mux
}

// GetHealthHandler returns the health handler instance for direct access if needed
func (hr *HandlerRegistry) GetHealthHandler() *HealthHandler {
	return hr.healthHandler
}

// updateRouteHandlers updates the RouteInfo registry with actual handler function references
func (hr *HandlerRegistry) updateRouteHandlers() {
	routes := GetRegisteredRoutes()
	for i, route := range routes {
		switch route.Path {
		case "/health":
			if hr.healthHandler != nil {
				routes[i].Handler = hr.healthHandler.ServeHTTP
			}
		case "/swagger":
			if hr.docsHandler != nil {
				routes[i].Handler = hr.docsHandler.ServeSwaggerUI
			}
		case "/docs/openapi.yaml":
			if hr.docsHandler != nil {
				routes[i].Handler = hr.docsHandler.ServeOpenAPISpec
			}
		case "/docs":
			if hr.docsHandler != nil {
				routes[i].Handler = hr.docsHandler.ServeDocs
			}
		}
	}

	// Update the global registry with the handler references
	types.UpdateRouteRegistry(routes)
}
