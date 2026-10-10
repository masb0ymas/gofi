package middlewares

import (
	"strings"
	"time"

	authulamodels "github.com/Authula/authula/models"
	authulaservices "github.com/Authula/authula/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"gofi/internal/config"
	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
)

// bearerToken returns the token of an "Authorization: Bearer <token>" header,
// or "" when the header is absent or uses another scheme.
func bearerToken(header string) string {
	scheme, token, _ := strings.Cut(header, " ")
	if !strings.EqualFold(scheme, "Bearer") {
		return ""
	}
	return strings.TrimSpace(token)
}

// Authorization requires a valid JWT access token or a live Authula session
// cookie and stores the signed-in user's id for handlers (see
// lib.ContextGetUID). It mirrors the "bearer.auth.optional" + "session.auth"
// checks for routes outside the Authula handler; sliding session renewal
// still happens on Authula's own routes (e.g. /me).
func (m *Middlewares) Authorization(c fiber.Ctx) error {
	userID := m.bearerUserID(c)
	if userID == "" {
		userID = m.sessionUserID(c)
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return apperr.ErrUnauthorized
	}
	lib.ContextSetUID(c, uid)

	return c.Next()
}

// bearerUserID returns the user a valid access token belongs to, or "". The
// jwt plugin verifies the signature and expiry and that the token's session
// still exists, so signing out revokes it.
func (m *Middlewares) bearerUserID(c fiber.Ctx) string {
	token := bearerToken(c.Get(fiber.HeaderAuthorization))
	if token == "" {
		return ""
	}

	jwt, ok := m.app.Auth.ServiceRegistry.Get(authulamodels.ServiceJWT.String()).(authulaservices.JWTService)
	if !ok {
		return ""
	}
	actor, err := jwt.ValidateToken(c.Context(), token)
	if err != nil || actor.Type != authulamodels.ActorUser {
		return ""
	}
	return actor.ID
}

// sessionUserID returns the user of a live session cookie, or "".
func (m *Middlewares) sessionUserID(c fiber.Ctx) string {
	token := c.Cookies(config.SessionCookieName)
	if token == "" {
		return ""
	}

	services := m.app.Auth.CoreServices()
	session, err := services.SessionService.GetByToken(c.Context(), services.TokenService.Hash(token))
	if err != nil || session == nil || session.ExpiresAt.Before(time.Now().UTC()) {
		return ""
	}
	return session.UserID
}
