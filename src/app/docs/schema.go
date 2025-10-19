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

func (g *OpenAPIGenerator) generatePageQuery() map[string]interface{} {
	return map[string]interface{}{
		"name":        "page",
		"in":          "query",
		"required":    false,
		"description": "Page number",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}

func (g *OpenAPIGenerator) generatePageSizeQuery() map[string]interface{} {
	return map[string]interface{}{
		"name":        "page_size",
		"in":          "query",
		"required":    false,
		"description": "Page size",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}

func (g *OpenAPIGenerator) generateFilteredQuery() map[string]interface{} {
	return map[string]interface{}{
		"name":        "filtered",
		"in":          "query",
		"required":    false,
		"description": "Query Filter",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}

func (g *OpenAPIGenerator) generateSortedQuery() map[string]interface{} {
	return map[string]interface{}{
		"name":        "sorted",
		"in":          "query",
		"required":    false,
		"description": "Query Sort",
		"schema": map[string]interface{}{
			"type": "string",
		},
	}
}
