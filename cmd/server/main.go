package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JerkyTreats/llm/internal/api/handler"
	"github.com/JerkyTreats/llm/internal/config"
	"github.com/JerkyTreats/llm/internal/logging"
)

func main() {
	logging.Info("Starting LLM API server")

	// Initialize handler registry
	handlerRegistry, err := handler.NewHandlerRegistry()
	if err != nil {
		logging.Error("Failed to initialize handler registry: %v", err)
		os.Exit(1)
	}

	// Get server configuration
	port := config.GetString("server_port")
	if port == "" {
		port = "8080" // Default port
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handlerRegistry.GetServeMux(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logging.Info("LLM API server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logging.Info("Shutting down LLM API server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logging.Error("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logging.Info("LLM API server stopped")
}
