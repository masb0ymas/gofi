package docs

func (g *OpenAPIGenerator) generateAuthSignUpRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"first_name": map[string]interface{}{
				"type":        "string",
				"example":     "John",
				"description": "First name of the user",
			},
			"last_name": map[string]interface{}{
				"type":        "string",
				"example":     "Doe",
				"description": "Last name of the user",
			},
			"email": map[string]interface{}{
				"type":        "string",
				"example":     "your@email.com",
				"description": "Email of the user",
			},
			"password": map[string]interface{}{
				"type":        "string",
				"format":      "password",
				"example":     "********",
				"description": "Password of the user",
			},
			"phone": map[string]interface{}{
				"type":        "string",
				"example":     "+628123456789",
				"description": "Phone number of the user",
			},
		},
		"required": []string{"first_name", "email", "password"},
	}
}

func (g *OpenAPIGenerator) generateAuthSignInRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"email": map[string]interface{}{
				"type":        "string",
				"example":     "your@email.com",
				"description": "Email of the user",
			},
			"password": map[string]interface{}{
				"type":        "string",
				"format":      "password",
				"example":     "********",
				"description": "Password of the user",
			},
		},
		"required": []string{"email", "password"},
	}
}

func (g *OpenAPIGenerator) generateAuthSignInResponse() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"display_name": map[string]interface{}{
				"type":    "string",
				"example": "John Doe",
			},
			"email": map[string]interface{}{
				"type":    "string",
				"example": "your@email.com",
			},
			"uid": map[string]interface{}{
				"type":    "string",
				"example": "your-uid",
			},
			"is_admin": map[string]interface{}{
				"type":    "boolean",
				"example": true,
			},
			"access_token": map[string]interface{}{
				"type":    "string",
				"example": "your-token",
			},
			"refresh_token": map[string]interface{}{
				"type":    "string",
				"example": "your-token",
			},
		},
	}
}

func (g *OpenAPIGenerator) generateAuthVerifySessionResponse() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"example":     "019aac2a-338a-79d1-960d-c80d1fa9e5b8",
				"description": "Unique identifier for the user",
			},
			"created_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"example":     "2025-08-27T15:38:30.383Z",
				"description": "Timestamp when the user was created",
			},
			"updated_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"example":     "2025-08-27T15:38:30.383Z",
				"description": "Timestamp when the user was updated",
			},
			"first_name": map[string]interface{}{
				"type":    "string",
				"example": "John",
			},
			"last_name": map[string]interface{}{
				"type":    "string",
				"example": "Doe",
			},
			"email": map[string]interface{}{
				"type":    "string",
				"example": "your@email.com",
			},
			"active_at": map[string]interface{}{
				"type":        "string",
				"format":      "date-time",
				"example":     "2025-08-27T15:38:30.383Z",
				"description": "Timestamp when the user was active",
			},
			"role_id": map[string]interface{}{
				"type":        "string",
				"example":     "019aac2a-338a-79d1-960d-c80d1fa9e5b8",
				"description": "Unique identifier for the role",
			},
		},
	}
}

func (g *OpenAPIGenerator) generateAuthVerifyRegistrationRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"token": map[string]interface{}{
				"type":        "string",
				"example":     "your-token",
				"description": "Verification token",
			},
		},
		"required": []string{"token"},
	}
}

// Sign Up
func (g *OpenAPIGenerator) generateAuthSignUp() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Sign Up",
		"description": "Sign Up for an account",
		"tags":        []string{"Auth"},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/SignUpRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Sign Up successfully",
				},
			},
		}),
	}
}

// Sign In
func (g *OpenAPIGenerator) generateAuthSignIn() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Sign In",
		"description": "Sign In to an account",
		"tags":        []string{"Auth"},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/SignInRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Sign In successfully",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/SignInResponse",
				},
			},
		}),
	}
}

// Verify Registration
func (g *OpenAPIGenerator) generateAuthVerifyRegistration() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Verify Registration",
		"description": "Verify Registration for an account",
		"tags":        []string{"Auth"},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/VerifyRegistrationRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Registration verified successfully",
				},
			},
		}),
	}
}

// Verify Session
func (g *OpenAPIGenerator) generateAuthVerifySession() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Verify Session",
		"description": "Verify Session for an account",
		"tags":        []string{"Auth"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Verify session successfully",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/VerifySessionResponse",
				},
			},
		}),
	}
}

// Sign Out
func (g *OpenAPIGenerator) generateAuthSignOut() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Sign Out",
		"description": "Sign Out from an account",
		"tags":        []string{"Auth"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Sign out successfully",
				},
			},
		}),
	}
}
