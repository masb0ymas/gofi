# Integrating and extending (library mode)

Source pages: https://www.authula.dev/docs/integrations/adapters , /custom-routes , /middleware , /hooks , /service-hooks , https://www.authula.dev/docs/guides/library-mode/programmatic-apis
(The Middleware page was not captured; fetch it if the user asks about middleware specifically.)

## Contents

- Router adapters
- Custom routes and RequestContext
- Request hooks
- Service hooks
- Programmatic API

## Router adapters

`auth.Handler()` is a plain `http.Handler`, so most routers need no adapter. Only Fiber has an official adapter.

```go
// net/http
http.Handle("/api/auth/", auth.Handler())

// Chi (Authula already uses Chi internally)
r := chi.NewRouter()
r.Handle("/auth/*", auth.Handler())     // use your base path

// Echo
e := echo.New()
e.Any("/api/auth/*", echo.WrapHandler(auth.Handler()))

// Fiber v3 (official adapter)
import fiberadapter "github.com/Authula/authula/adapters/fiber"
app.Use("/api/auth", fiberadapter.New(fiberadapter.Config{Handler: auth.Handler()}))
```

Caveat: routes you define on your own router outside the Authula handler **do not run Authula hooks**. To get hooks on your own endpoints, register them via `auth.RegisterCustomRoute(...)`.

## Custom routes

```go
auth.RegisterCustomRoute(authulamodels.Route{
    Method: "GET",
    Path:   "/api/health",               // NOT prefixed with base_path
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        reqCtx, _ := authulamodels.GetRequestContext(r.Context())
        reqCtx.SetJSONResponse(http.StatusOK, map[string]any{"status": "ok"})
    }),
    // Middleware: []func(http.Handler) http.Handler{},
    // Metadata: map[string]any{"plugins": []string{"session.auth"}},
})
```

Related methods on `Auth`:

- `RegisterRoute` / `RegisterRoutes`: registered **with** the base path prefix
- `RegisterCustomRoute` / `RegisterCustomRoutes`: **without** the prefix (application routes)
- `RegisterCustomRouteGroup(group)`: shared prefix and metadata for a group
- `RegisterMiddleware(...)`: call **before** `Handler()`
- `RegisterHook` / `RegisterHooks`

The `Route` struct: `Method`, `Path`, `Handler http.Handler`, `Middleware`, `Metadata map[string]any`. `Metadata["plugins"]` lists capabilities to run for the route (for example `"csrf.protect"`, `"session.auth"`).

Custom routes are **unauthenticated by default**. Add auth through `Metadata["plugins"]` or middleware.

`RequestContext` (from `authulamodels.GetRequestContext(r.Context())`) carries request data through the lifecycle: `Request`, `ResponseWriter`, `Path`, `Method`, `Headers`, `Actor`, `ClientIP`, `Values`, `Route`, `Handled`. Helpers: `SetJSONResponse(status, payload)`, `SetResponse(status, headers, body)`, `SetActorInContext(actor)`. Setting `Handled = true` short-circuits remaining hooks in the stage. To read the authenticated actor in your own handlers: `auth.GetActorFromRequest(req)` or `auth.GetActorFromContext(ctx)`.

## Request hooks

Four stages:

| Stage            | Runs                                               | Typical use                                   |
| ---------------- | -------------------------------------------------- | --------------------------------------------- |
| `HookOnRequest`  | Very start of every request, before route matching | global logging, metrics                       |
| `HookBefore`     | After route matching, before the handler           | auth, authorization, validation               |
| `HookAfter`      | After the handler, before the response is sent     | modify response, headers                      |
| `HookOnResponse` | After the response was written                     | analytics, audit (cannot change the response) |

```go
type Hook struct {
    Stage    HookStage
    PluginID string       // optional: only run if the route's metadata lists this plugin
    Matcher  HookMatcher  // optional: func(*RequestContext) bool
    Handler  HookHandler  // func(*RequestContext) error
    Order    int          // lower runs first; local to each PluginID group
    Async    bool         // background goroutine; side effects only
}
```

Rules:

- Hooks without `PluginID` run for all routes.
- A handler error is logged and, in the default mode, the next hook still runs. Error modes via `RouterOptions.HookErrorMode`: `error-log-continue` (default), `error-log-fail-fast` (stop remaining hooks in the stage), `error-silent`.
- `RouterOptions.AsyncHookTimeout` defaults to 30 s.
- Share data across hooks with `reqCtx.Values`; avoid storing sensitive data there.
- To block a request: `reqCtx.SetJSONResponse(401, ...)`, then `reqCtx.Handled = true`, then `return nil`.

```go
auth.RegisterHook(models.Hook{
    Stage: models.HookBefore,
    Matcher: func(rc *models.RequestContext) bool { return strings.HasPrefix(rc.Path, "/admin") },
    Handler: func(rc *models.RequestContext) error {
        if rc.Actor == nil || rc.Actor.Type != models.ActorUser {
            rc.SetJSONResponse(http.StatusUnauthorized, map[string]any{"message": "Authentication required"})
            rc.Handled = true
        }
        return nil
    },
})
```

Many common needs (session, JWT, CSRF, rate limit, RBAC) already exist as plugins, so check `references/plugins.md` before writing a hook.

## Service hooks

Run custom logic around core entity operations. Signature: `type ServiceHook[T any] func(*T) error`. Entities: Users, Accounts, Sessions (before/after create and update), Verifications (before/after create only).

- Before hooks may mutate the entity before it is stored.
- After hooks are for side effects.
- Returning an error aborts the operation.

```go
config := authulaconfig.NewConfig(
    authulaconfig.WithCoreServiceHooks(authulamodels.CoreServiceHooksConfig{
        Users: &authulamodels.ServiceHooks[authulamodels.User]{},
    }),
)
config.CoreServiceHooks.Users.RegisterAfterCreate(func(u *authulamodels.User) error {
    log.Printf("new user %s <%s>", u.ID, u.Email)
    return nil
})
```

The full config struct is at https://www.authula.dev/docs/reference/config (source: `models/config.go` in the repo).

## Programmatic API (no HTTP round trip)

In library mode keep the `Auth` handle:

```go
me, err := auth.Api.GetMe(ctx, "user-id")
res, err := auth.Api.SignOut(ctx, "user-id", nil /*sessionID*/, false /*signOutAll*/)
```

Plugins expose their own `Api` field. Two ways to reach it:

```go
// 1) keep the plugin instance you created
ep := emailpasswordplugin.New(cfg)
auth := authula.New(&authula.AuthConfig{Config: c, Plugins: []authulamodels.Plugin{ep}})
result, err := ep.Api.SignIn(ctx, email, password, nil, nil, nil) // callbackURL, ip, userAgent

// 2) look it up from the registry after init
p, ok := auth.PluginRegistry.GetPlugin(authulamodels.PluginOrganizations.String()).(*organizations.OrganizationsPlugin)
```

Use `auth.Api` for core functionality, a plugin's `Api` for plugin features, and lower-level helpers only for custom routing. The `Auth` type also offers `CoreServices()`, `DB()` (Bun), `EventBus()`, `Router()`.
