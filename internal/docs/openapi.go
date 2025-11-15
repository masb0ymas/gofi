package docs

import (
	"fmt"
	"gofi/internal/app"

	"github.com/gofiber/fiber/v2"
)

// OpenAPIGenerator generates OpenAPI specification
type OpenAPIGenerator struct {
	Title       string
	Version     string
	Description string
	ServerURL   string
}

// NewOpenAPIGenerator creates a new OpenAPI generator
func NewOpenAPIGenerator(opts OpenAPIGenerator) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		Title:       opts.Title,
		Version:     opts.Version,
		Description: opts.Description,
		ServerURL:   opts.ServerURL,
	}
}

// GenerateSpec generates the OpenAPI specification
func (g *OpenAPIGenerator) GenerateSpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       g.Title,
			"version":     g.Version,
			"description": g.Description,
		},
		"servers": []map[string]interface{}{
			{
				"url":         g.ServerURL,
				"description": "API Server",
			},
		},
		"paths":      g.generatePaths(),
		"components": g.generateComponents(),
		"tags":       g.generateTags(),
	}
}

func (g *OpenAPIGenerator) generatePaths() map[string]interface{} {
	return map[string]interface{}{
		// Role
		"/v1/roles": map[string]interface{}{
			"get":  g.generateGetRoles(),
			"post": g.generateCreateRole(),
		},
		"/v1/roles/:roleID": map[string]interface{}{
			"get":    g.generateGetRoleById(),
			"put":    g.generateUpdateRole(),
			"delete": g.generateForceDeleteRoleById(),
		},
		"/v1/roles/:roleID/soft-delete": map[string]interface{}{
			"delete": g.generateSoftDeleteRoleById(),
		},
		"/v1/roles/:roleID/restore": map[string]interface{}{
			"put": g.generateRestoreRoleById(),
		},
	}
}

func (g *OpenAPIGenerator) generateComponents() map[string]interface{} {
	return map[string]interface{}{
		"schemas": map[string]interface{}{
			// Header Authorization
			"AuthorizationHeader": g.generateAuthorizationHeader(),
			// Pagination
			"OffsetQuery": g.generateOffsetQuery(),
			"LimitQuery":  g.generateLimitQuery(),
			// Role
			"Role":              g.generateRoleModel(),
			"CreateRoleRequest": g.generateCreateRoleRequest(),
			// Error
			"UnauthorizedError":   g.generateErrorUnauthorized(),
			"NotFoundError":       g.generateErrorNotFound(),
			"InternalServerError": g.generateErrorInternalServer(),
		},
	}
}

func (g *OpenAPIGenerator) generateTags() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "Roles",
			"description": "Role management operations",
		},
	}
}

// GetScalarHTML returns the HTML for Scalar API documentation
func GetScalarHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
    </style>
</head>
<body>
    <script
        id="api-reference"
        data-url="/openapi.json"
        data-configuration='{
            "theme": "purple",
            "layout": "modern",
            "showSidebar": true,
            "hideDownloadButton": false,
            "searchHotKey": "k",
            "customCss": ".scalar-api-reference { height: 100vh; }",
            "metadata": {
                "title": "Gofi API Documentation",
                "description": "Complete API reference for the Gofi application"
            }
        }'></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
}

func SetupDocsRoutes(server *fiber.App, app *app.Application) {
	// Create OpenAPI generator
	generator := NewOpenAPIGenerator(
		OpenAPIGenerator{
			Title:       "Gofi API Documentation",
			Version:     "1.0.0",
			Description: "Complete API documentation for the Gofi application with user management endpoints",
			ServerURL:   fmt.Sprintf("http://localhost:%d", app.Config.App.Port),
		},
	)

	// OpenAPI JSON endpoint
	server.Get("/openapi.json", func(c *fiber.Ctx) error {
		return c.JSON(generator.GenerateSpec())
	})

	// Scalar docs endpoint
	server.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(GetScalarHTML())
	})
}
