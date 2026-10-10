# Routing and security concepts

Source pages: https://www.authula.dev/docs/concepts/routes , https://www.authula.dev/docs/concepts/security , https://www.authula.dev/docs/concepts/event-bus (event bus page not captured here; fetch it when needed)

## Contents

- Route mappings
- Disabled paths
- Mappings vs disabled paths
- Security model
- Event bus

## Route mappings

Route mappings say which plugin capabilities (and optionally permissions) run for which routes. Standalone: repeated `[[route_mappings]]` tables. Library: `authulaconfig.WithRouteMappings([]authulamodels.RouteMapping{...})`.

Fields:

- `paths`: one or more patterns
- `plugins`: capabilities to run, for example `session.auth`, `session.auth.optional`, `csrf.protect`
- `permissions`: optional permission names (enforced by the Access Control plugin)

Path forms:

- `METHOD:/path` for a specific method
- `/path` for all methods
- `{param}` for dynamic segments, for example `GET:/access-control/roles/{role_id}`
- `/prefix/*` for a whole subtree, for example `/organizations/*`

Behaviour:

- Several mappings can target the same route; Authula **merges** their plugin and permission lists. Use this to separate concerns (session in one mapping, permissions in another).
- Paths are **base-path aware**: with `base_path = "/api/auth"`, `GET:/me` means `GET:/api/auth/me`.
- A route with no mapping is still processed by Authula, just without extra capabilities.

```toml
[[route_mappings]]
paths = [
  "GET:/access-control/users/{user_id}/roles",
  "DELETE:/access-control/users/{user_id}/roles/{role_id}",
]
plugins = ["session.auth"]

[[route_mappings]]
paths = ["/organizations/*"]
plugins = ["session.auth"]
```

## Disabled paths

`disabled_paths` makes Authula skip route registration entirely. Use it to hide plugin endpoints you do not want exposed (typical when embedding as a library and providing your own endpoints), exclude externally managed paths, or turn off groups during testing.

```toml
# top-level only, not inside [plugins] or any other table
disabled_paths = [
  "GET:/email-password/request-email-change",
  "/organizations",
  "/admin/*",
]
```

Same three forms: `METHOD:/path`, `/path`, `/path/*`. Matching is base-path aware.

## Mappings vs disabled paths

- Route should stay available but needs auth or permission checks: **route mapping**.
- Route should not be handled by Authula at all: **disabled path**.

## Security model (summary)

Values below come from the Security concepts page. Where they conflict with defaults in `config.toml` or a plugin page, the config or plugin page wins for actual behaviour.

| Area | What Authula does |
| ---- | ----------------- |
| Password hashing | Argon2id (64 MB memory, 4 threads, 16-byte salts) |
| JWT signatures | Ed25519 |
| Data encryption | ChaCha20-Poly1305-X |
| Sessions | Sliding window with periodic revalidation; fingerprints IP and user agent to detect hijacking; session tokens are hashed in the database |
| JWT | Access/refresh token pair; short-lived access-token (about 15 min), refresh (about 7 days); automatic key rotation about every 30 days with a 1-hour grace period |
| CSRF | Double-submit cookie (24-byte tokens via header) plus Go 1.25 `CrossOriginProtection` (checks `Sec-Fetch-Site`/`Origin`) |
| Rate limiting | In-memory, Redis, or DB backends with failover; `X-RateLimit` headers; proxy-aware client IP |
| Token invalidation | Redis-backed blacklist with TTL; reuse events `TokenReuseRecoveredEvent` (first reuse) and `TokenReuseMaliciousEvent` (repeated) |
| IP handling | Zero-trust: ignores `X-Forwarded-For` unless the source is in `trusted_proxies` |
| Headers and CORS | Strict origin validation, no wildcard with credentials; injects `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy` |
| Config vaulting | Sensitive config keys are detected and encrypted in memory |

Practical rules to give users:

- `secure = true` on cookies in production, and serve over HTTPS.
- Configure `trusted_proxies` when behind a load balancer or Cloudflare (plus `trusted_headers` such as `CF-Connecting-IP`).
- Use the rate-limit plugin on sign-in and password-reset endpoints.
- Keep secrets in environment variables.

## Event bus (Event Driven Architecture)

Internal auth events flow through an event bus (Watermill-based). Providers: `gochannel`, `sqlite`, `postgres`, `redis`, `kafka`, `nats`, `rabbitmq`. Plugins emit events (for example TOTP emits `totp.enabled`, `totp.disabled`, `totp.verified`, `totp.backup_code_used`, `totp.device_trusted`) that your code can subscribe to. In library mode, get the bus with `auth.EventBus()`. For subscription API details, fetch https://www.authula.dev/docs/concepts/event-bus.
