package docs

func (g *OpenAPIGenerator) generateMetadataResponse() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"total": map[string]interface{}{
				"type":    "integer",
				"example": 10,
			},
		},
	}
}
