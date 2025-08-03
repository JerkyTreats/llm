package docs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/JerkyTreats/llm/internal/config"
	"github.com/JerkyTreats/llm/internal/logging"
)

// DocsHandler serves Swagger UI and OpenAPI specifications
type DocsHandler struct {
	swaggerConfig SwaggerConfig
}

// SwaggerConfig represents the swagger configuration
type SwaggerConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Path      string `yaml:"path"`
	SpecPath  string `yaml:"spec_path"`
	UITitle   string `yaml:"ui.title"`
	Theme     string `yaml:"ui.theme"`
}

// NewDocsHandler creates a new documentation handler
func NewDocsHandler() (*DocsHandler, error) {
	swaggerConfig := SwaggerConfig{
		Enabled:  true,
		Path:     "/swagger",
		SpecPath: "/docs/openapi.yaml",
		UITitle:  "LLM API Documentation",
		Theme:    "dark",
	}

	return &DocsHandler{
		swaggerConfig: swaggerConfig,
	}, nil
}

// ServeSwaggerUI serves the Swagger UI interface
func (h *DocsHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logging.Debug("Serving Swagger UI for path: %s", r.URL.Path)

	// Generate Swagger UI HTML
	html := h.generateSwaggerHTML(r)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// ServeOpenAPISpec serves the OpenAPI specification file
func (h *DocsHandler) ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logging.Debug("Serving OpenAPI spec for path: %s", r.URL.Path)

	// Find the OpenAPI spec file
	specPath := "docs/api/openapi.yaml"
	
	// Check if file exists
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		logging.Warn("OpenAPI spec file not found: %s", specPath)
		http.Error(w, "OpenAPI specification not found", http.StatusNotFound)
		return
	}

	// Read and serve the file
	content, err := os.ReadFile(specPath)
	if err != nil {
		logging.Error("Failed to read OpenAPI spec: %v", err)
		http.Error(w, "Failed to read OpenAPI specification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow CORS for Swagger UI
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// generateSwaggerHTML generates the Swagger UI HTML page
func (h *DocsHandler) generateSwaggerHTML(r *http.Request) string {
	// Determine the current protocol from the request
	var baseURL string
	
	// Log for debugging to understand what's happening with the request
	logging.Debug("Swagger HTML generation - Host: %s, TLS: %v, URL: %s, X-Forwarded-Proto: %s", 
		r.Host, r.TLS != nil, r.URL.String(), r.Header.Get("X-Forwarded-Proto"))
	
	// Use request host, but provide fallback if empty
	host := r.Host
	if host == "" {
		// Fallback: construct from server config
		logging.Warn("Request Host header is empty, falling back to server config")
		serverHost := config.GetString("server.host")
		serverPort := config.GetInt("server.port")
		
		if serverHost == "0.0.0.0" || serverHost == "" {
			host = fmt.Sprintf("localhost:%d", serverPort)
		} else {
			host = fmt.Sprintf("%s:%d", serverHost, serverPort)
		}
		logging.Debug("Using fallback host: %s", host)
	}
	
	// Determine if request was made over HTTPS
	// Check both direct TLS and common proxy headers
	isHTTPS := r.TLS != nil ||
		r.Header.Get("X-Forwarded-Proto") == "https" ||
		r.Header.Get("X-Forwarded-Scheme") == "https" ||
		strings.ToLower(r.Header.Get("X-Forwarded-Ssl")) == "on"
	
	logging.Debug("HTTPS detection - TLS: %v, X-Forwarded-Proto: %s, X-Forwarded-Scheme: %s, X-Forwarded-Ssl: %s, Final isHTTPS: %v",
		r.TLS != nil, r.Header.Get("X-Forwarded-Proto"), r.Header.Get("X-Forwarded-Scheme"), 
		r.Header.Get("X-Forwarded-Ssl"), isHTTPS)
	
	if isHTTPS {
		// Request came via HTTPS, use HTTPS for spec URL
		baseURL = fmt.Sprintf("https://%s", host)
	} else {
		// Request came via HTTP, use HTTP for spec URL
		baseURL = fmt.Sprintf("http://%s", host)
	}
	
	logging.Debug("Swagger using base URL: %s", baseURL)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
        %s
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '%s/docs/openapi.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                docExpansion: "list",
                defaultModelsExpandDepth: 3,
                defaultModelExpandDepth: 3,
                displayRequestDuration: true,
                filter: true,
                tryItOutEnabled: true,
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
                onComplete: function() {
                    console.log('Swagger UI loaded successfully');
                },
                onFailure: function(error) {
                    console.error('Swagger UI failed to load:', error);
                }
            });
        };
    </script>
</body>
</html>`, h.swaggerConfig.UITitle, h.getThemeCSS(), baseURL)
}

// getThemeCSS returns CSS for the configured theme
func (h *DocsHandler) getThemeCSS() string {
	if strings.ToLower(h.swaggerConfig.Theme) == "dark" {
		return `
        body {
            background: #1f1f1f !important;
        }
        .swagger-ui {
            filter: invert(88%) hue-rotate(180deg);
        }
        .swagger-ui .microlight {
            filter: invert(100%) hue-rotate(180deg);
        }`
	}
	return ""
}

// ServeDocs handles requests to the docs directory (for static files if needed)
func (h *DocsHandler) ServeDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Remove /docs prefix and get the requested file path
	requestPath := strings.TrimPrefix(r.URL.Path, "/docs")
	if requestPath == "" || requestPath == "/" {
		// Redirect to swagger UI
		http.Redirect(w, r, "/swagger", http.StatusFound)
		return
	}

	// Security check: prevent directory traversal
	if strings.Contains(requestPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Construct file path
	filePath := filepath.Join("docs", strings.TrimPrefix(requestPath, "/"))
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}