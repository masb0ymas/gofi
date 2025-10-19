package docs

func (g *OpenAPIGenerator) generateAuthSignUpRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"fullname": map[string]interface{}{
				"type":        "string",
				"description": "Fullname of the user",
			},
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Email of the user",
			},
			"password": map[string]interface{}{
				"type":        "string",
				"format":      "password",
				"description": "Password of the user",
			},
		},
		"required": []string{"fullname", "email", "password"},
	}
}

func (g *OpenAPIGenerator) generateAuthSignInRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Email of the user",
			},
			"password": map[string]interface{}{
				"type":        "string",
				"format":      "password",
				"description": "Password of the user",
			},
		},
		"required": []string{"email", "password"},
	}
}

func (g *OpenAPIGenerator) generateAuthResetPasswordRequest() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Email of the user",
			},
			"token": map[string]interface{}{
				"type":        "string",
				"description": "Token of the user",
			},
			"password": map[string]interface{}{
				"type":        "string",
				"format":      "password",
				"description": "Password of the user",
			},
		},
		"required": []string{"email", "token", "password"},
	}
}

func (g *OpenAPIGenerator) generateAuthSignUp() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Sign Up",
		"description": "Sign Up",
		"tags":        []string{"Auth"},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/AuthSignUpRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been registered",
				},
			},
		}),
	}
}

func (g *OpenAPIGenerator) generateAuthSignIn() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Sign In",
		"description": "Sign In",
		"tags":        []string{"Auth"},
		"requestBody": map[string]interface{}{
			"required": true,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"$ref": "#/components/schemas/AuthSignInRequest",
					},
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been signed in",
				},
			},
		}),
	}
}

func (g *OpenAPIGenerator) generateAuthVerify() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Verify",
		"description": "Verify",
		"tags":        []string{"Auth"},
		"parameters": []map[string]interface{}{
			{
				"name":        "uid",
				"in":          "query",
				"required":    true,
				"description": "UID",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
			{
				"name":        "token-verify",
				"in":          "query",
				"required":    true,
				"description": "Token Verify",
				"schema": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been verified",
				},
			},
		}),
	}
}

func (g *OpenAPIGenerator) generateAuthVerifySession() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Verify Session",
		"description": "Verify Session",
		"tags":        []string{"Auth"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"read", "write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been verified",
				},
			},
		}),
	}
}

func (g *OpenAPIGenerator) generateAuthSignOut() map[string]interface{} {
	return map[string]interface{}{
		"summary":     "Sign Out",
		"description": "Sign Out",
		"tags":        []string{"Auth"},
		"security": []map[string]interface{}{
			{
				"BearerAuth": []string{"read", "write"},
			},
		},
		"parameters": []map[string]interface{}{
			{"$ref": "#/components/schemas/AuthorizationHeader"},
		},
		"responses": g.generateResponse(Response{
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":    "string",
					"example": "data has been signed out",
				},
			},
		}),
	}
}
