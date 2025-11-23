package docs

func (g *OpenAPIGenerator) generateErrorNotFound() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"success": map[string]interface{}{
				"type":    "boolean",
				"example": false,
			},
			"message": map[string]interface{}{
				"type":    "string",
				"example": "Not Found",
			},
		},
		"required": []string{"success", "message"},
	}
}

func (g *OpenAPIGenerator) generateErrorForbidden() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"success": map[string]interface{}{
				"type":    "boolean",
				"example": false,
			},
			"message": map[string]interface{}{
				"type":    "string",
				"example": "Forbidden",
			},
		},
		"required": []string{"success", "message"},
	}
}

func (g *OpenAPIGenerator) generateErrorUnauthorized() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"success": map[string]interface{}{
				"type":    "boolean",
				"example": false,
			},
			"message": map[string]interface{}{
				"type":    "string",
				"example": "Token verification failed: token is malformed: token contains an invalid number of segments. Please ensure your token is valid and not expired.",
			},
		},
		"required": []string{"success", "message"},
	}
}

func (g *OpenAPIGenerator) generateErrorInternalServer() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"success": map[string]interface{}{
				"type":    "boolean",
				"example": false,
			},
			"message": map[string]interface{}{
				"type":    "string",
				"example": "Internal Server Error",
			},
		},
		"required": []string{"success", "message"},
	}
}
