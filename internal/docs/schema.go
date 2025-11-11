package docs

func (g *OpenAPIGenerator) generateAuthorizationHeader() map[string]interface{} {
	return map[string]interface{}{
		"name":        "Authorization",
		"in":          "header",
		"required":    true,
		"description": "Bearer Token",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}

func (g *OpenAPIGenerator) generateOffsetQuery() map[string]interface{} {
	return map[string]interface{}{
		"name":        "offset",
		"in":          "query",
		"required":    false,
		"description": "Offset number",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}

func (g *OpenAPIGenerator) generateLimitQuery() map[string]interface{} {
	return map[string]interface{}{
		"name":        "limit",
		"in":          "query",
		"required":    false,
		"description": "Limit number",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}
