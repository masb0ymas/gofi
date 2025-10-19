package docs

func (g *OpenAPIGenerator) generateMetadataResponse() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"current_page": map[string]interface{}{
				"type":    "integer",
				"example": 1,
			},
			"page_size": map[string]interface{}{
				"type":    "integer",
				"example": 10,
			},
			"total": map[string]interface{}{
				"type":    "integer",
				"example": 99,
			},
			"total_pages": map[string]interface{}{
				"type":    "integer",
				"example": 9,
			},
		},
	}
}
