package jwt

import (
	"io"
	"net/http/httptest"
	"testing"

	"gofi/internal/config"

	"github.com/gofiber/fiber/v2"
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name      string
		setupCtx  func(*fiber.Ctx)
		wantToken string
		wantErr   bool
	}{
		{
			name: "Extract token from query parameter",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().SetRequestURI("/?token=query-token-123")
			},
			wantToken: "query-token-123",
			wantErr:   false,
		},
		{
			name: "Extract token from cookie",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.SetCookie("token", "cookie-token-456")
			},
			wantToken: "cookie-token-456",
			wantErr:   false,
		},
		{
			name: "Extract token from Authorization header with Bearer",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer header-token-789")
			},
			wantToken: "header-token-789",
			wantErr:   false,
		},
		{
			name: "Query parameter takes precedence over cookie",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().SetRequestURI("/?token=query-token")
				c.Request().Header.SetCookie("token", "cookie-token")
			},
			wantToken: "query-token",
			wantErr:   false,
		},
		{
			name: "Query parameter takes precedence over header",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().SetRequestURI("/?token=query-token")
				c.Request().Header.Set("Authorization", "Bearer header-token")
			},
			wantToken: "query-token",
			wantErr:   false,
		},
		{
			name: "Cookie takes precedence over header",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.SetCookie("token", "cookie-token")
				c.Request().Header.Set("Authorization", "Bearer header-token")
			},
			wantToken: "cookie-token",
			wantErr:   false,
		},
		{
			name: "Invalid Authorization header - missing Bearer prefix",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "token-without-bearer")
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Invalid Authorization header - only Bearer without token",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer ")
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Invalid Authorization header - Bearer with empty token",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer  ")
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Invalid Authorization header - wrong prefix",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Basic token-123")
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Invalid Authorization header - too many parts",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer token extra-part")
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "No token provided anywhere",
			setupCtx:  func(c *fiber.Ctx) {},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Empty query parameter is ignored",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().SetRequestURI("/?token=")
				c.Request().Header.Set("Authorization", "Bearer fallback-token")
			},
			wantToken: "fallback-token",
			wantErr:   false,
		},
		{
			name: "Token with special characters in query",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().SetRequestURI("/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			wantErr:   false,
		},
		{
			name: "Token with special characters in cookie",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.SetCookie("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			wantErr:   false,
		},
		{
			name: "Token with special characters in Authorization header",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			wantErr:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			j := &JWT{
				config: &config.ConfigApp{
					Name:      "gofi",
					JWTSecret: "test-secret",
				},
			}

			var gotToken string
			var gotErr error

			app.Get("/test", func(c *fiber.Ctx) error {
				tc.setupCtx(c)
				gotToken, gotErr = j.ExtractToken(c)
				return c.SendString("ok")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to perform test request: %v", err)
			}
			defer resp.Body.Close()

			// Read response body to ensure request completes
			_, _ = io.ReadAll(resp.Body)

			if (gotErr != nil) != tc.wantErr {
				t.Errorf("ExtractToken() error = %v, wantErr %v", gotErr, tc.wantErr)
				return
			}

			if gotToken != tc.wantToken {
				t.Errorf("ExtractToken() token = %v, want %v", gotToken, tc.wantToken)
			}
		})
	}
}

func TestExtractTokenErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		setupCtx       func(*fiber.Ctx)
		wantErrMessage string
	}{
		{
			name: "No token - token not found error",
			setupCtx: func(c *fiber.Ctx) {
				// No token setup
			},
			wantErrMessage: "token not found",
		},
		{
			name: "Invalid format - missing Bearer",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "just-a-token")
			},
			wantErrMessage: "invalid token format",
		},
		{
			name: "Invalid format - empty token after Bearer",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer ")
			},
			wantErrMessage: "invalid token format",
		},
		{
			name: "Invalid format - too many parts",
			setupCtx: func(c *fiber.Ctx) {
				c.Request().Header.Set("Authorization", "Bearer token1 token2")
			},
			wantErrMessage: "invalid token format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			j := &JWT{
				config: &config.ConfigApp{
					Name:      "gofi",
					JWTSecret: "test-secret",
				},
			}

			var gotErr error

			app.Get("/test", func(c *fiber.Ctx) error {
				tc.setupCtx(c)
				_, gotErr = j.ExtractToken(c)
				return c.SendString("ok")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to perform test request: %v", err)
			}
			defer resp.Body.Close()

			// Read response body to ensure request completes
			_, _ = io.ReadAll(resp.Body)

			if gotErr == nil {
				t.Errorf("ExtractToken() expected error but got nil")
				return
			}

			if gotErr.Error() != tc.wantErrMessage {
				t.Errorf("ExtractToken() error message = %v, want %v", gotErr.Error(), tc.wantErrMessage)
			}
		})
	}
}
