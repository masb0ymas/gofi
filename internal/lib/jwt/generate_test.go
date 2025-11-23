package jwt

import (
	"strings"
	"testing"
	"time"

	"gofi/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.ConfigApp
		payload       *JWTPayload
		wantErr       bool
		validateToken bool
	}{
		{
			name: "Valid token generation with 7 days expiry",
			config: &config.ConfigApp{
				Name:      "gofi",
				JWTSecret: "test-secret-key-12345",
			},
			payload: &JWTPayload{
				UID:       uuid.New().String(),
				ExpiresAt: "7",
			},
			wantErr:       false,
			validateToken: true,
		},
		{
			name: "Valid token generation with 1 day expiry",
			config: &config.ConfigApp{
				Name:      "gofi",
				JWTSecret: "test-secret-key-12345",
			},
			payload: &JWTPayload{
				UID:       uuid.New().String(),
				ExpiresAt: "1",
			},
			wantErr:       false,
			validateToken: true,
		},
		{
			name: "Valid token generation with 30 days expiry",
			config: &config.ConfigApp{
				Name:      "gofi",
				JWTSecret: "test-secret-key-12345",
			},
			payload: &JWTPayload{
				UID:       uuid.New().String(),
				ExpiresAt: "30",
			},
			wantErr:       false,
			validateToken: true,
		},
		{
			name: "Invalid ExpiresAt - not a number",
			config: &config.ConfigApp{
				Name:      "gofi",
				JWTSecret: "test-secret-key-12345",
			},
			payload: &JWTPayload{
				UID:       uuid.New().String(),
				ExpiresAt: "invalid",
			},
			wantErr:       true,
			validateToken: false,
		},
		{
			name: "Empty UID",
			config: &config.ConfigApp{
				Name:      "gofi",
				JWTSecret: "test-secret-key-12345",
			},
			payload: &JWTPayload{
				UID:       "",
				ExpiresAt: "7",
			},
			wantErr:       false,
			validateToken: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			j := &JWT{config: tc.config}

			token, expiresAt, err := j.Generate(tc.payload)

			if (err != nil) != tc.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if tc.wantErr {
				return
			}

			// Verify token is not empty
			if token == "" {
				t.Error("Generate() returned empty token")
			}

			// Verify token format (should have 3 parts separated by dots)
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				t.Errorf("Generate() token has incorrect format, got %d parts, want 3", len(parts))
			}

			// Verify expiresAt is in the future
			if expiresAt <= time.Now().Unix() {
				t.Error("Generate() expiresAt should be in the future")
			}

			// Validate token if required
			if tc.validateToken {
				parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
					return []byte(tc.config.JWTSecret), nil
				})
				if err != nil {
					t.Errorf("Generate() produced invalid token: %v", err)
				}

				if !parsedToken.Valid {
					t.Error("Generate() produced invalid token")
				}

				// Verify claims
				if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
					// Check UID
					if uid, ok := claims["uid"].(string); !ok || uid != tc.payload.UID {
						t.Errorf("Generate() token UID = %v, want %v", claims["uid"], tc.payload.UID)
					}

					// Check issuer
					if iss, ok := claims["iss"].(string); !ok || iss != tc.config.Name {
						t.Errorf("Generate() token issuer = %v, want %v", claims["iss"], tc.config.Name)
					}

					// Check expiration
					if exp, ok := claims["exp"].(float64); !ok || int64(exp) != expiresAt {
						t.Errorf("Generate() token exp = %v, want %v", claims["exp"], expiresAt)
					}
				} else {
					t.Error("Generate() token claims are invalid")
				}
			}
		})
	}
}

func TestGenerateTokenExpiry(t *testing.T) {
	j := &JWT{
		config: &config.ConfigApp{
			Name:      "gofi",
			JWTSecret: "test-secret-key-12345",
		},
	}

	tests := []struct {
		name      string
		expiresAt string
		wantDays  int
	}{
		{"1 day expiry", "1", 1},
		{"7 days expiry", "7", 7},
		{"30 days expiry", "30", 30},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := &JWTPayload{
				UID:       uuid.New().String(),
				ExpiresAt: tc.expiresAt,
			}

			_, expiresAt, err := j.Generate(payload)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			// Calculate expected expiry time
			expectedExpiry := time.Now().Add(time.Hour * 24 * time.Duration(tc.wantDays)).Unix()

			// Allow 2 seconds tolerance for test execution time
			diff := expiresAt - expectedExpiry
			if diff < -2 || diff > 2 {
				t.Errorf("Generate() expiresAt = %v, want approximately %v (diff: %v seconds)", expiresAt, expectedExpiry, diff)
			}
		})
	}
}
