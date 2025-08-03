package handler

import (
	"reflect"

	"github.com/JerkyTreats/template-goapi/internal/api/types"
)

func init() {
	// Register health check endpoint
	types.RegisterRoute(types.RouteInfo{
		Method:       "GET",
		Path:         "/health",
		Handler:      nil, // Will be set during handler initialization
		RequestType:  nil, // GET request has no body
		ResponseType: reflect.TypeOf(HealthResponse{}),
		Module:       "health",
		Summary:      "Health check endpoint returning service status",
	})
}