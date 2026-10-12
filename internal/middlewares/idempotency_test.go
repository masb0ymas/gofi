package middlewares

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
)

func TestScopedIdempotencyKey(t *testing.T) {
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	other := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	clientKey := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	otherKey := uuid.MustParse("44444444-4444-4444-4444-444444444444")

	base := scopedIdempotencyKey(uid, fiber.MethodPost, "/v1/categories", clientKey)

	if _, err := uuid.Parse(base); err != nil {
		t.Fatalf("derived key %q is not a UUID: %v", base, err)
	}
	// A retry of the same request must land on the same entry.
	if got := scopedIdempotencyKey(uid, fiber.MethodPost, "/v1/categories", clientKey); got != base {
		t.Errorf("same request derived %q, want %q", got, base)
	}

	// Each component is part of the scope: reusing a key must not reach an
	// entry another caller or route owns.
	cases := map[string]string{
		"other caller": scopedIdempotencyKey(other, fiber.MethodPost, "/v1/categories", clientKey),
		"other method": scopedIdempotencyKey(uid, fiber.MethodPut, "/v1/categories", clientKey),
		"other route":  scopedIdempotencyKey(uid, fiber.MethodPost, "/v1/categories/1", clientKey),
		"other key":    scopedIdempotencyKey(uid, fiber.MethodPost, "/v1/categories", otherKey),
	}
	for name, got := range cases {
		if got == base {
			t.Errorf("%s derived the same key as the base request: %q", name, got)
		}
	}
}

// idempotencyApp mounts Idempotency() behind a handler that sets the caller id
// and counts how many times the write handler actually ran.
func idempotencyApp(t *testing.T, uid uuid.UUID, runs *atomic.Int32) *fiber.App {
	t.Helper()

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil)))})
	app.Use(func(c fiber.Ctx) error {
		lib.ContextSetUID(c, uid)
		return c.Next()
	})
	app.Post("/write", Idempotency(), func(c fiber.Ctx) error {
		if c.Get("Fail") != "" {
			return apperr.New(apperr.KindConflict, "conflict")
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"run": runs.Add(1)})
	})
	app.Get("/write", Idempotency(), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"run": runs.Add(1)})
	})
	app.Post("/other", Idempotency(), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"run": runs.Add(1)})
	})
	return app
}

// call sends one request with the given idempotency header name and value.
func call(t *testing.T, app *fiber.App, method, path, header, key string, extra ...string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, path, http.NoBody)
	if header != "" {
		req.Header.Set(header, key)
	}
	for i := 0; i+1 < len(extra); i += 2 {
		req.Header.Set(extra[i], extra[i+1])
	}
	resp, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func body(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestIdempotencyReplaysWrite(t *testing.T) {
	const key = "8e03978e-40d5-43e8-bc93-6894a57f9324"

	var runs atomic.Int32
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	app := idempotencyApp(t, uid, &runs)

	first := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, key)
	firstBody := body(t, first)
	if first.StatusCode != fiber.StatusCreated {
		t.Fatalf("first write: %d %s", first.StatusCode, firstBody)
	}
	if first.Header.Get(HeaderIdempotentReplayed) != "" {
		t.Error("first write was marked as replayed")
	}

	// The retry is answered from the cache: same status and body, one run.
	retry := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, key)
	retryBody := body(t, retry)
	if retry.StatusCode != first.StatusCode || retryBody != firstBody {
		t.Errorf("retry answered %d %s, want %d %s", retry.StatusCode, retryBody, first.StatusCode, firstBody)
	}
	if retry.Header.Get(HeaderIdempotentReplayed) != "true" {
		t.Errorf("retry %s = %q, want \"true\"", HeaderIdempotentReplayed, retry.Header.Get(HeaderIdempotentReplayed))
	}
	if retry.Header.Get(fiber.HeaderContentType) == "" {
		t.Error("retry dropped the Content-Type of the stored response")
	}
	if got := runs.Load(); got != 1 {
		t.Errorf("handler ran %d times, want 1", got)
	}

	// A new key is a new request.
	if resp := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, "00000000-0000-0000-0000-000000000001"); resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("write with a fresh key: %d", resp.StatusCode)
	} else {
		body(t, resp)
	}
	if got := runs.Load(); got != 2 {
		t.Errorf("handler ran %d times, want 2", got)
	}
}

func TestIdempotencyScope(t *testing.T) {
	const key = "8e03978e-40d5-43e8-bc93-6894a57f9324"

	var runs atomic.Int32
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	app := idempotencyApp(t, uid, &runs)

	body(t, call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, key))

	// Same key, another route: not this route's cached answer.
	body(t, call(t, app, fiber.MethodPost, "/other", HeaderIdempotencyKey, key))
	if got := runs.Load(); got != 2 {
		t.Errorf("handler ran %d times across two routes, want 2", got)
	}

	// Same key, another caller: their own entry, not the first caller's.
	otherApp := idempotencyApp(t, uuid.MustParse("22222222-2222-2222-2222-222222222222"), &runs)
	body(t, call(t, otherApp, fiber.MethodPost, "/write", HeaderIdempotencyKey, key))
	if got := runs.Load(); got != 3 {
		t.Errorf("handler ran %d times across two callers, want 3", got)
	}
}

func TestIdempotencyHeaderHandling(t *testing.T) {
	var runs atomic.Int32
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	app := idempotencyApp(t, uid, &runs)

	// The draft's header name works too, and shares the scope with the
	// Fiber-default one.
	const key = "8e03978e-40d5-43e8-bc93-6894a57f9324"
	body(t, call(t, app, fiber.MethodPost, "/write", headerIdempotencyKeyDraft, key))
	replay := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, key)
	if replay.Header.Get(HeaderIdempotentReplayed) != "true" {
		t.Error("a key sent under the draft header name did not reach the same cache entry")
	}
	body(t, replay)
	if got := runs.Load(); got != 1 {
		t.Errorf("handler ran %d times, want 1", got)
	}

	// No key: every request runs the handler.
	body(t, call(t, app, fiber.MethodPost, "/write", "", ""))
	body(t, call(t, app, fiber.MethodPost, "/write", "", ""))
	if got := runs.Load(); got != 3 {
		t.Errorf("handler ran %d times without a key, want 3", got)
	}

	// A malformed key is the client's mistake, not a server error.
	if resp := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, "not-a-uuid"); resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("malformed key: %d, want 400", resp.StatusCode)
	} else {
		body(t, resp)
	}

	// Safe methods are not deduplicated, even with a key.
	body(t, call(t, app, fiber.MethodGet, "/write", HeaderIdempotencyKey, key))
	body(t, call(t, app, fiber.MethodGet, "/write", HeaderIdempotencyKey, key))
	if got := runs.Load(); got != 5 {
		t.Errorf("handler ran %d times, want 5 (two GETs must not be deduplicated)", got)
	}
}

func TestIdempotencyDoesNotCacheErrors(t *testing.T) {
	const key = "8e03978e-40d5-43e8-bc93-6894a57f9324"

	var runs atomic.Int32
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	app := idempotencyApp(t, uid, &runs)

	// A refused write is not stored: the retry must be able to succeed.
	failed := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, key, "Fail", "1")
	if failed.StatusCode != fiber.StatusConflict {
		t.Fatalf("failing write: %d, want 409", failed.StatusCode)
	}
	body(t, failed)

	retry := call(t, app, fiber.MethodPost, "/write", HeaderIdempotencyKey, key)
	if retry.StatusCode != fiber.StatusCreated {
		t.Errorf("retry after a failure: %d, want 201", retry.StatusCode)
	}
	body(t, retry)
	if got := runs.Load(); got != 1 {
		t.Errorf("handler ran %d times, want 1 (the failed run must not have been cached)", got)
	}
}

func BenchmarkIdempotencyScopedKey(b *testing.B) {
	uid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	clientKey := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	b.ReportAllocs()
	for b.Loop() {
		_ = scopedIdempotencyKey(uid, fiber.MethodPost, "/v1/categories", clientKey)
	}
}
