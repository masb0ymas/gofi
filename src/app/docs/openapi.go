package docs

import "github.com/gofiber/fiber/v2"

// OpenAPIGenerator generates OpenAPI specification
type OpenAPIGenerator struct {
	Title       string
	Version     string
	Description string
	ServerURL   string
}

// NewOpenAPIGenerator creates a new OpenAPI generator
func NewOpenAPIGenerator(title, version, description, serverURL string) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		Title:       title,
		Version:     version,
		Description: description,
		ServerURL:   serverURL,
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
		// Auth
		"/v1/auth/sign-up": map[string]interface{}{
			"post": g.generateAuthSignUp(),
		},
		"/v1/auth/sign-in": map[string]interface{}{
			"post": g.generateAuthSignIn(),
		},
		"/v1/auth/verify": map[string]interface{}{
			"get": g.generateAuthVerify(),
		},
		"/v1/auth/verify-session": map[string]interface{}{
			"get": g.generateAuthVerifySession(),
		},
		"/v1/auth/sign-out": map[string]interface{}{
			"post": g.generateAuthSignOut(),
		},

		// Role
		"/v1/role": map[string]interface{}{
			"get":  g.generateGetRoles(),
			"post": g.generateCreateRole(),
		},
		"/v1/role/:id": map[string]interface{}{
			"get": g.generateGetRoleById(),
			"put": g.generateUpdateRole(),
		},
		"/v1/role/restore/:id": map[string]interface{}{
			"put": g.generateRestoreRoleById(),
		},
		"/v1/role/soft-delete/:id": map[string]interface{}{
			"delete": g.generateSoftDeleteRoleById(),
		},
		"/v1/role/force-delete/:id": map[string]interface{}{
			"delete": g.generateForceDeleteRoleById(),
		},
	}
}

func (g *OpenAPIGenerator) generateComponents() map[string]interface{} {
	return map[string]interface{}{
		"schemas": map[string]interface{}{
			// Header Authorization
			"AuthorizationHeader": g.generateAuthorizationHeader(),
			// Pagination
			"PageQuery":     g.generatePageQuery(),
			"PageSizeQuery": g.generatePageSizeQuery(),
			"FilteredQuery": g.generateFilteredQuery(),
			"SortedQuery":   g.generateSortedQuery(),
			// Auth
			"AuthSignUpRequest":        g.generateAuthSignUpRequest(),
			"AuthSignInRequest":        g.generateAuthSignInRequest(),
			"AuthResetPasswordRequest": g.generateAuthResetPasswordRequest(),
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
			"name":        "Auth",
			"description": "Auth endpoints",
		},
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

func SetupDocsRoutes(app *fiber.App) {
	// Create OpenAPI generator
	generator := NewOpenAPIGenerator(
		"Gofi API Documentation",
		"1.0.0",
		"Complete API documentation for the Gofi application with user management endpoints",
		"http://localhost:8000",
	)

	// OpenAPI JSON endpoint
	app.Get("/openapi.json", func(c *fiber.Ctx) error {
		return c.JSON(generator.GenerateSpec())
	})

	// Scalar docs endpoint
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(GetScalarHTML())
	})
}
