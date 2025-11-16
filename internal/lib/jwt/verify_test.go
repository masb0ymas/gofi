package jwt

import (
	"strings"
	"testing"
	"time"

	"gofi/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestVerify(t *testing.T) {
	secret := "test-secret-key-12345"
	appName := "gofi"

	j := &JWT{
		config: &config.ConfigApp{
			Name:      appName,
			JWTSecret: secret,
		},
	}

	// Helper to create a valid token
	createValidToken := func(uid string, expiresIn time.Duration) string {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(expiresIn).Unix(),
			"iss": appName,
			"uid": uid,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte(secret))
		return signedToken
	}

	// Helper to create a token with wrong secret
	createTokenWithWrongSecret := func(uid string) string {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
			"iss": appName,
			"uid": uid,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte("wrong-secret"))
		return signedToken
	}

	tests := []struct {
		name        string
		token       string
		wantClaims  *JWTClaims
		wantErr     bool
		errContains string
	}{
		{
			name:  "Valid token with all claims",
			token: createValidToken(uuid.New().String(), time.Hour*24),
			wantClaims: &JWTClaims{
				Iss: appName,
			},
			wantErr: false,
		},
		{
			name:  "Valid token with long expiry",
			token: createValidToken(uuid.New().String(), time.Hour*24*30),
			wantClaims: &JWTClaims{
				Iss: appName,
			},
			wantErr: false,
		},
		{
			name:  "Valid token with short expiry",
			token: createValidToken(uuid.New().String(), time.Minute*5),
			wantClaims: &JWTClaims{
				Iss: appName,
			},
			wantErr: false,
		},
		{
			name:        "Expired token",
			token:       createValidToken(uuid.New().String(), -time.Hour),
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
		{
			name:        "Token with wrong secret",
			token:       createTokenWithWrongSecret(uuid.New().String()),
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
		{
			name:        "Malformed token - not a JWT",
			token:       "not.a.valid.jwt.token",
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
		{
			name:        "Malformed token - missing parts",
			token:       "header.payload",
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
		{
			name:        "Empty token",
			token:       "",
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
		{
			name:        "Token with invalid base64",
			token:       "invalid!!!.base64!!!.encoding!!!",
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
		{
			name:        "Random string",
			token:       "random-string-not-a-token",
			wantClaims:  nil,
			wantErr:     true,
			errContains: "Token verification failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := j.Verify(tc.token)

			if (err != nil) != tc.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if tc.wantErr {
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("Verify() error = %v, should contain %v", err.Error(), tc.errContains)
				}
				return
			}

			if claims == nil {
				t.Error("Verify() returned nil claims for valid token")
				return
			}

			// Verify issuer
			if claims.Iss != tc.wantClaims.Iss {
				t.Errorf("Verify() claims.Iss = %v, want %v", claims.Iss, tc.wantClaims.Iss)
			}

			// Verify UID is not empty
			if claims.UID == "" {
				t.Error("Verify() claims.UID is empty")
			}

			// Verify expiration is in the future
			if claims.Exp <= time.Now().Unix() {
				t.Error("Verify() claims.Exp should be in the future")
			}
		})
	}
}

func TestVerifyClaimsExtraction(t *testing.T) {
	secret := "test-secret-key-12345"
	appName := "gofi"
	uid := uuid.New().String()

	j := &JWT{
		config: &config.ConfigApp{
			Name:      appName,
			JWTSecret: secret,
		},
	}

	// Create a token with specific claims
	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"exp": expiresAt,
		"iss": appName,
		"uid": uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Verify the token
	verifiedClaims, err := j.Verify(signedToken)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	// Check all claims are correctly extracted
	if verifiedClaims.UID != uid {
		t.Errorf("Verify() claims.UID = %v, want %v", verifiedClaims.UID, uid)
	}

	if verifiedClaims.Iss != appName {
		t.Errorf("Verify() claims.Iss = %v, want %v", verifiedClaims.Iss, appName)
	}

	if verifiedClaims.Exp != expiresAt {
		t.Errorf("Verify() claims.Exp = %v, want %v", verifiedClaims.Exp, expiresAt)
	}
}

func TestVerifyWithDifferentSecrets(t *testing.T) {
	appName := "gofi"
	uid := uuid.New().String()

	tests := []struct {
		name         string
		tokenSecret  string
		verifySecret string
		wantErr      bool
	}{
		{
			name:         "Matching secrets",
			tokenSecret:  "secret-123",
			verifySecret: "secret-123",
			wantErr:      false,
		},
		{
			name:         "Different secrets",
			tokenSecret:  "secret-123",
			verifySecret: "secret-456",
			wantErr:      true,
		},
		{
			name:         "Empty verify secret",
			tokenSecret:  "secret-123",
			verifySecret: "",
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create token with tokenSecret
			claims := jwt.MapClaims{
				"exp": time.Now().Add(time.Hour).Unix(),
				"iss": appName,
				"uid": uid,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, _ := token.SignedString([]byte(tc.tokenSecret))

			// Verify with verifySecret
			j := &JWT{
				config: &config.ConfigApp{
					Name:      appName,
					JWTSecret: tc.verifySecret,
				},
			}

			_, err := j.Verify(signedToken)

			if (err != nil) != tc.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestVerifyIntegrationWithGenerate(t *testing.T) {
	j := &JWT{
		config: &config.ConfigApp{
			Name:      "gofi",
			JWTSecret: "test-secret-key-12345",
		},
	}

	uid := uuid.New().String()

	// Generate a token
	payload := &JWTPayload{
		UID:       uid,
		ExpiresAt: "7",
	}

	token, expiresAt, err := j.Generate(payload)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify the generated token
	claims, err := j.Verify(token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	// Validate claims
	if claims.UID != uid {
		t.Errorf("Verify() claims.UID = %v, want %v", claims.UID, uid)
	}

	if claims.Iss != j.config.Name {
		t.Errorf("Verify() claims.Iss = %v, want %v", claims.Iss, j.config.Name)
	}

	if claims.Exp != expiresAt {
		t.Errorf("Verify() claims.Exp = %v, want %v", claims.Exp, expiresAt)
	}
}

func TestVerifyErrorMessages(t *testing.T) {
	j := &JWT{
		config: &config.ConfigApp{
			Name:      "gofi",
			JWTSecret: "test-secret",
		},
	}

	tests := []struct {
		name           string
		token          string
		wantErrContain string
	}{
		{
			name:           "Empty token",
			token:          "",
			wantErrContain: "Token verification failed",
		},
		{
			name:           "Malformed token",
			token:          "not-a-jwt",
			wantErrContain: "Token verification failed",
		},
		{
			name:           "Invalid signature",
			token:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.invalid_signature",
			wantErrContain: "Token verification failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := j.Verify(tc.token)

			if err == nil {
				t.Error("Verify() expected error but got nil")
				return
			}

			if !strings.Contains(err.Error(), tc.wantErrContain) {
				t.Errorf("Verify() error = %v, should contain %v", err.Error(), tc.wantErrContain)
			}
		})
	}
}
