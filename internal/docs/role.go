package docs

// Model
func (g *OpenAPIGenerator) generateRoleModel() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"example":     "019aac2a-338a-79d1-960d-c80d1fa9e5b8",
				"description": "Unique identifier for the role",
			},
			"created_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"example":     "2025-08-27T15:38:30.383Z",
				"description": "Timestamp when the role was created",
			},
			"updated_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"example":     "2025-08-27T15:38:30.383Z",
				"description": "Timestamp when the role was updated",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"example":     "Guest",
				"description": "Name of the role",
			},
		},
		"required": []string{"id", "name"},
	}
}

func (g *OpenAPIGenerator) generateCreateRoleRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"example":     "Guest",
				"description": "Name of the role",
			},
		},
		"required": []string{"name"},
	}
}

// Find All
func (g *OpenAPIGenerator) generateGetRoles() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Get Roles",
		"description": "Retrieve all roles",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
			{"$ref": "#/components/schemas/OffsetQuery"},
			{"$ref": "#/components/schemas/LimitQuery"},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"meta": g.generateMetadataResponse(),
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been received",
				},
				"data": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/components/schemas/Role",
					},
				},
			},
		}),
	}
}

// Find By ID
func (g *OpenAPIGenerator) generateGetRoleById() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Get Role by ID",
		"description": "Retrieve a specific role by their ID",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
			{
				"name":        "roleID",
				"in":          "path",
				"required":    true,
				"description": "Role ID",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been received",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/Role",
				},
			},
		}),
	}
}

// Create
func (g *OpenAPIGenerator) generateCreateRole() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Create Role",
		"description": "Create a new role",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
		},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/CreateRoleRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been created",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/Role",
				},
			},
		}),
	}
}

// Update
func (g *OpenAPIGenerator) generateUpdateRole() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Update Role",
		"description": "Update an existing role",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
			{
				"name":        "roleID",
				"in":          "path",
				"required":    true,
				"description": "Role ID",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/CreateRoleRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been updated",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/Role",
				},
			},
		}),
	}
}

// Restore
func (g *OpenAPIGenerator) generateRestoreRoleById() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Restore Role by ID",
		"description": "Restore a specific role by their ID",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
			{
				"name":        "roleID",
				"in":          "path",
				"required":    true,
				"description": "Role ID",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been restored",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/Role",
				},
			},
		}),
	}
}

// Soft Delete
func (g *OpenAPIGenerator) generateSoftDeleteRoleById() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Soft Delete Role by ID",
		"description": "Soft delete a specific role by their ID",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
			{
				"name":        "roleID",
				"in":          "path",
				"required":    true,
				"description": "Role ID",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been soft deleted",
				},
			},
		}),
	}
}

// Force Delete
func (g *OpenAPIGenerator) generateForceDeleteRoleById() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Force Delete Role by ID",
		"description": "Force delete a specific role by their ID",
		"tags":        []string{"Roles"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
			{
				"name":        "roleID",
				"in":          "path",
				"required":    true,
				"description": "Role ID",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been permanently deleted",
				},
			},
		}),
	}
}
