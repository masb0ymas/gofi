package docs

type Response struct {
	Properties map[string]interface{}
}

func (g *OpenAPIGenerator) generateResponse(response Response) map[string]interface{} {
	return map[string]interface{}{
		"200": map[string]interface{}{
			"description": "Successful response",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type":       "object",
						"properties": response.Properties,
					},
				},
			},
		},
		"201": map[string]interface{}{
			"description": "Created response",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type":       "object",
						"properties": response.Properties,
					},
				},
			},
		},
		"401": map[string]interface{}{
			"description": "Unauthorized",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/UnauthorizedError",
					},
				},
			},
		},
		"403": map[string]interface{}{
			"description": "Forbidden",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/ForbiddenError",
					},
				},
			},
		},
		"404": map[string]interface{}{
			"description": "Not found",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/NotFoundError",
					},
				},
			},
		},
		"500": map[string]interface{}{
			"description": "Internal Server Error",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/InternalServerError",
					},
				},
			},
		},
	}
}
