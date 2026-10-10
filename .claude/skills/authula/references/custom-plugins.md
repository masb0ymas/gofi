# Writing a custom Authula plugin (Go)

Source: https://www.authula.dev/docs/guides/server-plugins . For TypeScript SDK client plugins see https://www.authula.dev/docs/guides/client-plugins (not captured here).

## Contents

- Required interface
- Optional capabilities
- Step-by-step
- Config and best practices

## Required interface

```go
type Plugin interface {
    Metadata() PluginMetadata   // ID, Version, Description
    Config() any                // typed config struct
    Init(ctx *PluginContext) error
    Close() error
}

type PluginContext struct {
    DB              bun.IDB
    Logger          Logger
    EventBus        EventBus
    ServiceRegistry ServiceRegistry   // Register(name, svc) / Get(name)
    GetConfig       func() *Config
}
```

Plugins talk to each other through the service registry, for example `ctx.ServiceRegistry.Get(models.ServiceUser.String()).(services.UserService)`.

## Optional capabilities (implement the interface to opt in)

| Interface              | Method(s)                                                                    | Adds                                                                         |
| ---------------------- | ---------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `PluginWithMigrations` | `Migrations(provider string) []migrations.Migration`, `DependsOn() []string` | DB tables, per provider (sqlite/postgres/mysql) via `migrations.ForProvider` |
| `PluginWithRoutes`     | `Routes() []Route`                                                           | HTTP endpoints                                                               |
| `PluginWithMiddleware` | `Middleware() []func(http.Handler) http.Handler`                             | global middleware                                                            |
| `AuthMethodProvider`   | `AuthMiddleware()`, `OptionalAuthMiddleware()`                               | a new auth method (reject vs allow anonymous)                                |
| `PluginWithHooks`      | `Hooks() []Hook`                                                             | lifecycle hooks (see `references/extending.md`)                              |

Hook execution rules: hooks without `PluginID` run everywhere; with `PluginID` they run only when the route's `Metadata["plugins"]` lists it. Sorted by `PluginID` group, then `Order`. `Matcher` returning false skips the hook. `Handled = true` stops later hooks in that stage.

## Step by step

```go
package myplugin

type MyPluginConfig struct {
    Enabled   bool   `json:"enabled" toml:"enabled"`
    SomeField string `json:"some_field" toml:"some_field"`
}

type MyPlugin struct {
    config      MyPluginConfig
    logger      models.Logger
    ctx         *models.PluginContext
    userService services.UserService
}

func New(cfg MyPluginConfig) *MyPlugin { return &MyPlugin{config: cfg} }

func (p *MyPlugin) Metadata() models.PluginMetadata {
    return models.PluginMetadata{ID: "my_plugin", Version: "1.0.0", Description: "Custom functionality"}
}
func (p *MyPlugin) Config() any { return p.config }

func (p *MyPlugin) Init(ctx *models.PluginContext) error {
    p.ctx, p.logger = ctx, ctx.Logger
    // merge values from config.toml [plugins.my_plugin]
    if err := util.LoadPluginConfig(ctx.GetConfig(), p.Metadata().ID, &p.config); err != nil {
        return err
    }
    svc, ok := ctx.ServiceRegistry.Get(models.ServiceUser.String()).(services.UserService)
    if !ok {
        return fmt.Errorf("user service not available")
    }
    p.userService = svc
    return nil
}
func (p *MyPlugin) Close() error { return nil }

func (p *MyPlugin) Routes() []models.Route {
    return []models.Route{{
        Method: "GET", Path: "/api/my-plugin/info",
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            util.JSONResponse(w, http.StatusOK, map[string]string{"plugin_id": p.Metadata().ID})
        }),
    }}
}
```

Register it:

```go
auth := authula.New(&authula.AuthConfig{
    Config: authulaconfig.NewConfig(),
    Plugins: []models.Plugin{myplugin.New(myplugin.MyPluginConfig{SomeField: "value"})},
})
```

Note: the docs' imports use `github.com/Authula/authula/internal/util`, which is an `internal` package. Go forbids importing it from outside the Authula module, so in your own module write small local helpers (JSON response, config loading) or check whether a public helper exists in the current release.

## Config

```toml
[plugins.my_plugin]
enabled = true
some_field = "value"
```

Plugin IDs: lowercase with underscores.

## Best practices

1. Validate config in `Init` and return an error when invalid.
2. Release resources in `Close`.
3. Use the injected logger.
4. `Async: true` only for side effects (logging, analytics, events, webhooks). Never for auth, CSRF, or rate limiting.
5. Set the authenticated identity with `reqCtx.SetActorInContext(actor)`.
6. Keep responses consistent (JSON helper).
7. Before writing a plugin, check whether an existing plugin or a plain hook already covers the need.
