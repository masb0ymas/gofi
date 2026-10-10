---
name: authula
description: Guide for building authentication with Authula, the open-source, plugin-based Go auth framework (github.com/Authula/authula, docs at authula.dev/docs). Use this skill whenever the user mentions Authula, or wants to add or configure auth in a Go backend using its plugins (email-password, oauth2, session, jwt, totp, magic-link, csrf, rate-limit, organizations, access-control, admin, api-key, bearer, secondary-storage), config.toml route_mappings, hooks, service hooks, custom routes, or a custom Authula plugin. Also use it for running Authula as a standalone auth server (Docker, config.toml) or embedding it as a Go library. Trigger even if the user only says "Authula config", "session.auth", "authula plugin", or shows code importing github.com/Authula/authula.
---

# Authula

Authula is an open-source authentication framework for Go. It runs in two modes:

- **Library mode**: embed it in a Go app; `auth.Handler()` returns a standard `http.Handler`.
- **Standalone mode**: run it as an auth server (Docker image `ghcr.io/authula/authula:latest`) configured with `config.toml`, usable from any tech stack (a Node.js SDK exists).

Everything is delivered through **plugins**. Each plugin exposes **capabilities** (for example `session.auth`), and every route declares which capabilities apply to it through **route mappings**. Nothing is implicit. Under the hood it uses Chi, Bun ORM, and Watermill (event bus).

The docs move fast (the module was at v1.46.0 in Sept 2026). This skill is a condensed snapshot. **When precision matters (exact config fields, a plugin's endpoints, a signature), fetch the live page** listed in `references/plugins.md` instead of trusting memory.

## First, work out the situation

Before writing code, establish (ask only what you can't infer):

1. **Mode**: library (Go app) or standalone server?
2. **Auth methods needed**: email/password, OAuth2 (Discord/GitHub/Google), magic link, TOTP 2FA, JWT/bearer, API keys?
3. **Database**: SQLite, PostgreSQL, or MySQL.
4. **Router** (library mode): net/http, Chi, Echo, Fiber.

## Where to read next

Read only the file that matches the task:

| Task                                                                    | File                                 |
| ----------------------------------------------------------------------- | ------------------------------------ |
| Install, `config.toml`, `.env`, Docker, first working setup, curl tests | `references/setup.md`                |
| Protecting routes, route mappings, disabled paths, security model       | `references/routing-and-security.md` |
| Router adapters, custom routes, hooks, service hooks, programmatic API  | `references/extending.md`            |
| Writing your own plugin                                                 | `references/custom-plugins.md`       |
| Which plugin does what, per-plugin config/endpoints, live doc URLs      | `references/plugins.md`              |

## Core mental model (keep in mind for every answer)

1. **Enable a plugin** (config `[plugins.<id>] enabled = true`, or pass `<plugin>.New(...)` in `AuthConfig.Plugins` in library mode).
2. **Protect routes with route mappings.** A plugin being enabled does not protect anything by itself. Routes are protected only when a mapping attaches a capability such as `session.auth` (401 if no valid session) or `session.auth.optional` (continue anonymously).
3. **Route paths in mappings are relative to `base_path`.** `GET:/me` with base path `/api/auth` is served at `/api/auth/me`.
4. **Secrets belong in env vars**, not in `config.toml`: `AUTHULA_SECRET` (generate with `openssl rand -hex 32`), `AUTHULA_DATABASE_URL`, OAuth client secrets.
5. **Migrations run automatically** when the Authula instance is created (core and plugin tables).

## Minimal library-mode skeleton

```go
config := authulaconfig.NewConfig(
    authulaconfig.WithAppName("MyApp"),
    authulaconfig.WithBasePath("/api/auth"),
    authulaconfig.WithDatabase(authulamodels.DatabaseConfig{Provider: "sqlite", URL: "auth.db"}),
    authulaconfig.WithRouteMappings([]authulamodels.RouteMapping{
        {Paths: []string{"GET:/me", "POST:/sign-out"}, Plugins: []string{"session.auth"}},
    }),
)
auth := authula.New(&authula.AuthConfig{
    Config: config,
    Plugins: []authulamodels.Plugin{
        sessionplugin.New(sessionplugin.SessionPluginConfig{Enabled: true}),
        // emailpasswordplugin.New(...), emailplugin.New(...), ...
    },
})
log.Fatal(http.ListenAndServe(":8080", auth.Handler()))
```

## Pitfalls to warn about

- **CORS + sessions**: with the session plugin and `allow_credentials = true`, `allowed_origins` must list exact origins. `["*"]` is rejected by browsers and fails.
- **CSRF plugin**: if used, the CSRF token header must also be listed in `security.cors.allowed_headers`.
- **TOTP**: keep `/totp/verify` and `/totp/verify-backup-code` reachable without `session.auth`; they rely on a pending-token cookie. Also set `app_name`, since it goes into the authenticator URI.
- **Custom routes are not authenticated by default.** Add capabilities via `Metadata: {"plugins": [...]}` or middleware.
- **Routes registered outside the Authula handler do not get Authula hooks.** Use `auth.RegisterCustomRoute()` if hooks must apply.
- **`disabled_paths` must be top-level** in `config.toml`, not under a table such as `[plugins]`.
- **Trusted proxies**: Authula ignores `X-Forwarded-For` unless the proxy CIDR is configured in `security.trusted_proxies`. Behind a load balancer, rate limiting sees the wrong IP without it.
- **Async hooks are for side effects only** (logging, analytics). Never use them for auth, CSRF, or rate limiting.
- **Docs have small inconsistencies** (for example `authulaconfig.New` vs `NewConfig`, or `authula.New(config)` vs `authula.New(&AuthConfig{...})`). The verified forms are `authulaconfig.NewConfig(...)` and `authula.New(&authula.AuthConfig{...})`. If a compile error appears, check the current pkg.go.dev page for `github.com/Authula/authula`.
- **Secure cookie settings**: the sample configs use `secure = false` for local development. Tell users to set `secure = true` behind HTTPS in production.

## How to answer

- Match the user's mode (TOML for standalone, Go options for library) and give only that one.
- Show the complete wiring: plugin enabled **and** route mapping **and** required env vars.
- Note anything that differs between the snapshot and the live docs, and cite the page URL when you fetched it.
