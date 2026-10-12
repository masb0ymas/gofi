package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/google/uuid"

	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
)

// Idempotency guards state-changing routes against duplicate work: a client
// that repeats a write with the same key gets the original response instead of
// a second row. Mount it after Authorization, which supplies the caller id.
//
// The Fiber middleware keys its cache on the raw header value, so two callers
// using the same key would share one entry. This wrapper scopes the key to the
// caller and the route (see scopedIdempotencyKey) before delegating, which also
// keeps the stored key a fixed-width UUID.
//
// ponytail: the response cache and its lock live in this process, so each
// replica deduplicates on its own; move both to shared storage (Redis) once
// there is more than one instance.
func Idempotency() fiber.Handler {
	guard := idempotency.New(idempotency.Config{
		KeyHeader: HeaderIdempotencyKey,
		// Keep only the content type: fasthttp recomputes Content-Length and
		// Date when it writes the replayed response, and a stored copy of
		// either would be stale.
		KeepResponseHeaders: []string{fiber.HeaderContentType},
		Lifetime:            idempotencyLifetime,
	})

	return func(c fiber.Ctx) error {
		// The IETF draft names the header Idempotency-Key; Fiber defaults to
		// X-Idempotency-Key. Accept either so a client following the draft is
		// not silently retried into a duplicate write.
		header := HeaderIdempotencyKey
		raw := c.Get(header)
		if raw == "" {
			header = headerIdempotencyKeyDraft
			raw = c.Get(header)
		}
		if raw == "" {
			return c.Next()
		}

		// Normalise before deriving, so "{uuid}" and the bare uuid share one
		// cache entry. The derived key is a UUID, which satisfies the
		// middleware's own 36-character check.
		clientKey, err := uuid.Parse(raw)
		if err != nil {
			return apperr.New(apperr.KindBadRequest, "%s must be a UUID", header)
		}

		uid, err := lib.ContextGetUID(c)
		if err != nil {
			return apperr.ErrUnauthorized
		}

		// The middleware reads this header itself; rewriting it is how the
		// scope reaches the cache. No handler reads it downstream.
		c.Request().Header.Set(HeaderIdempotencyKey, scopedIdempotencyKey(uid, c.Method(), c.Path(), clientKey))

		if err := guard(c); err != nil {
			return err
		}

		// Only a replayed response is marked; the stored copy is written
		// before this runs, so the marker never enters the cache.
		if idempotency.IsFromCache(c) {
			c.Set(HeaderIdempotentReplayed, "true")
		}
		return nil
	}
}

// Header names and the retry window.
const (
	// HeaderIdempotencyKey carries the client-generated key on a write.
	HeaderIdempotencyKey = "X-Idempotency-Key"
	// headerIdempotencyKeyDraft is the name the IETF draft standardises.
	headerIdempotencyKeyDraft = "Idempotency-Key"
	// HeaderIdempotentReplayed marks a response served from the cache.
	HeaderIdempotentReplayed = "Idempotent-Replayed"
	// idempotencyLifetime is how long a key stays replayable; it has to outlast
	// the client's retry budget, and memory grows with the traffic inside it.
	idempotencyLifetime = 30 * time.Minute
)

// idempotencyNamespace is the fixed UUIDv5 namespace for derived keys; changing
// it invalidates every key in flight.
var idempotencyNamespace = uuid.MustParse("a814bdb8-df9a-4f80-80c5-cf636f828a6f")

// scopedIdempotencyKey derives the storage key from the caller, the route and
// the client's key. Every component is included so a key reused by another user
// or on another route cannot replay this response; a retry of the same request
// repeats all three and still hits the entry.
func scopedIdempotencyKey(uid uuid.UUID, method, path string, clientKey uuid.UUID) string {
	scope := fmt.Appendf(nil, "%s\x00%s\x00%s\x00%s", uid, method, path, clientKey)
	return uuid.NewSHA1(idempotencyNamespace, scope).String()
}
