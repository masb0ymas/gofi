# Plugins: index, key details, live doc URLs

Every plugin page follows the same layout: Overview, Configuration (standalone TOML and library Go), API Reference, Database Schema, Plugin Capabilities, Security Recommendations, Client Plugin (TypeScript SDK). **For anything beyond what is below, fetch the plugin's page.**

## Contents

- Plugin index
- Session
- Email and Email & Password
- OAuth2
- TOTP
- Other core pages and SDK

## Plugin index

Base URL: `https://www.authula.dev/docs/plugins/`

| Plugin            | Page                                                            | What it is                                                                             | Captured here? |
| ----------------- | --------------------------------------------------------------- | -------------------------------------------------------------------------------------- | -------------- |
| Session           | `session`                                                       | Cookie sessions, sliding renewal; capabilities `session.auth`, `session.auth.optional` | yes            |
| Email & Password  | `email-password`                                                | Sign-up, sign-in, verification, reset, email change                                    | yes            |
| Email             | `email`                                                         | Transactional email with multiple providers and automatic failover                     | summary only   |
| OAuth2            | `oauth2` (+ `oauth2/discord`, `oauth2/github`, `oauth2/google`) | Social login                                                                           | yes            |
| TOTP              | `totp`                                                          | Authenticator-app 2FA, backup codes, trusted devices                                   | yes            |
| Magic Link        | `magic-link`                                                    | Passwordless email links                                                               | not captured   |
| JWT               | `jwt`                                                           | JWT auth (Ed25519, key rotation per security page)                                     | not captured   |
| Bearer            | `bearer`                                                        | Bearer-token auth                                                                      | not captured   |
| API Key           | `api-key`                                                       | API key auth                                                                           | not captured   |
| CSRF              | `csrf`                                                          | Double-submit-cookie CSRF protection; capability `csrf.protect`                        | summary only   |
| Rate Limit        | `ratelimit` (note: URL has no hyphen; config ID is `ratelimit`) | Rate limiting                                                                          | not captured   |
| Secondary Storage | `secondary-storage`                                             | In-memory, DB, or Redis key-value store for rate-limit counters and similar data       | summary only   |
| Access Control    | `access-control`                                                | Roles and permissions; used together with route mapping `permissions`                  | not captured   |
| Admin             | `admin`                                                         | Admin endpoints                                                                        | not captured   |
| Organizations     | `organizations`                                                 | Multi-tenant orgs                                                                      | not captured   |

For "not captured" or "summary only" plugins, run `web_fetch` on the page before giving config keys, endpoints, or capability names. Do not guess them.

Library import pattern (from the verified examples):

```go
import (
    sessionplugin "github.com/Authula/authula/plugins/session"
    emailpasswordplugin "github.com/Authula/authula/plugins/email-password"
    emailpasswordplugintypes "github.com/Authula/authula/plugins/email-password/types"
    oauth2plugin "github.com/Authula/authula/plugins/oauth2"
    oauth2plugintypes "github.com/Authula/authula/plugins/oauth2/types"
    totpplugin "github.com/Authula/authula/plugins/totp"
    totpplugintypes "github.com/Authula/authula/plugins/totp/types"
)
```

## Session

Provides the actual cookie-session implementation (the global `[session]` config only sets its properties).

- Capabilities: `session.auth` (401 when the cookie is missing or invalid) and `session.auth.optional` (continue anonymously). In Go: `sessionplugin.HookIDSessionAuth.String()` and `sessionplugin.HookIDSessionAuthOptional.String()`.
- Also has global hooks that issue the session cookie after successful authentication and remove it on sign-out.
- No HTTP endpoints and no database tables of its own.
- Works alongside other auth plugins such as JWT without conflict.
- Enable: `[plugins.session] enabled = true` or `sessionplugin.New(sessionplugin.SessionPluginConfig{Enabled: true})`.
- Session lifetime, cookie name, `secure`, `same_site`, `max_sessions_per_user`, cleanup: global `[session]` table (see `references/setup.md`). Docs examples show differing `expires_in` values (24h in one page, 30m in another), so always set it explicitly.
- Tokens are hashed before storage; use `secure = true` in production.

## Email and Email & Password

Email plugin (needed for verification, reset, and change emails): `[plugins.email] enabled = true, provider = "smtp", from_address = "..."`, with SMTP credentials via env vars.

Email & Password plugin `[plugins.email_password]` keys: `enabled`, `min_password_length` (8 recommended), `max_password_length` (128), `disable_sign_up`, `require_email_verification`, `auto_sign_in`, `send_email_on_sign_up`, `send_email_on_sign_in`, `email_verification_expires_in` (24h), `password_reset_expires_in` (1h), `request_email_change_expires_in` (1h).

Library-only optional callbacks to override email sending: `SendEmailVerification`, `SendPasswordResetEmail`, `SendChangedPasswordEmail`, `SendRequestEmailChangeEmail`, `SendChangedEmailToOldEmail`, `SendChangedEmailToNewEmail`.

Endpoints (under base path):

| Method | Path                                      |
| ------ | ----------------------------------------- |
| POST   | `/email-password/sign-up`                 |
| POST   | `/email-password/sign-in`                 |
| GET    | `/email-password/verify-email`            |
| POST   | `/email-password/send-email-verification` |
| POST   | `/email-password/request-password-reset`  |
| POST   | `/email-password/change-password`         |
| POST   | `/email-password/request-email-change`    |

No tables and no capabilities of its own. Passwords are hashed with Argon2. Recommendations: `require_email_verification = true` in production, short reset expiry, consider the rate-limit plugin on sign-in and reset endpoints. Typical mapping: sign-in, sign-up, verify-email use `session.auth.optional`; the rest use `session.auth` (see `references/setup.md`).

## OAuth2

Providers: Discord, GitHub, Google.

```toml
[plugins.oauth2]
enabled = true

[plugins.oauth2.providers.google]
enabled = true
client_id = "your-client-id"
client_secret = "..."          # prefer env vars (GOOGLE_CLIENT_SECRET etc.)
redirect_url = "http://localhost:8080/auth/oauth2/callback/google"
scopes = []
```

```go
oauth2plugin.New(oauth2plugintypes.OAuth2PluginConfig{
    Enabled: true,
    Providers: map[string]oauth2plugintypes.ProviderConfig{
        "google": {
            Enabled:      true,
            ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
            ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
            RedirectURL:  "https://yourdomain.com/auth/oauth2/callback/google",
        },
    },
})
```

Endpoints: `GET /oauth2/authorize/{provider}` starts the flow; `GET /oauth2/callback/{provider}` handles the callback. The redirect URL includes the configured base path (as in the examples above) and must match what you registered with the provider; use HTTPS in production. No tables, no capabilities. Request minimal scopes. Provider setup guides: `oauth2/discord`, `oauth2/github`, `oauth2/google`.

## TOTP

```toml
[plugins.totp]
enabled = true
skip_verification_on_enable = false   # keep false: require password when enabling
backup_code_count = 10
trusted_device_duration = "720h"
trusted_devices_auto_cleanup = true
trusted_devices_cleanup_interval = "1h"
pending_token_expiry = "5m"
secure_cookie = true
same_site = "lax"
```

Requires `app_name` to be set (used in the `otpauth://` URI).

| Method | Path                          | Auth                                               |
| ------ | ----------------------------- | -------------------------------------------------- |
| POST   | `/totp/enable`                | session; returns `totp_uri` and `backup_codes`     |
| POST   | `/totp/disable`               | session                                            |
| GET    | `/totp/get-uri`               | session                                            |
| POST   | `/totp/generate-backup-codes` | session                                            |
| POST   | `/totp/verify`                | pending-token cookie; body `{code, trust_device?}` |
| POST   | `/totp/verify-backup-code`    | pending-token cookie; single-use codes             |

Capability `totp.intercept`: during sign-in it sets a pending-token cookie and returns a JSON response telling the client to redirect to the TOTP verification step. Do **not** put `session.auth` on `/totp/verify` or `/totp/verify-backup-code`. Tables: `totp` (encrypted secret, hashed backup codes) and `trusted_devices`. Events: `totp.enabled`, `totp.disabled`, `totp.verified`, `totp.backup_code_used`, `totp.device_trusted`.

## Other core pages and SDK

- Config reference: https://www.authula.dev/docs/reference/config (points to `models/config.go` in the repo; each plugin's config lives in its folder under `plugins/`)
- Database schema: https://www.authula.dev/docs/reference/database-schema
- Email templates: https://www.authula.dev/docs/reference/email-templates
- SDKs: https://www.authula.dev/docs/reference/sdks (Node.js/TypeScript SDK: https://github.com/Authula/authula-node-sdk, `createClient({ url, plugins: [new EmailPasswordPlugin(), ...] })` from `authula` and `authula/plugins`)
- Playground (Next.js + React Router frontends with a Go/Echo backend): https://github.com/Authula/authula-playground
- Architecture and comparison with other solutions: https://www.authula.dev/docs/introduction/architecture , https://www.authula.dev/docs/introduction/comparison
- Go API reference: https://pkg.go.dev/github.com/Authula/authula
- Discord community: linked from the GitHub README
