# Setup: standalone and library mode

Source pages: https://www.authula.dev/docs/get-started/installation and https://www.authula.dev/docs/get-started/basic-usage

## Contents

- Standalone mode (config.toml, .env, Docker)
- Library mode (install, config, database, plugins, mount)
- Working example: email/password + sessions
- Testing with curl

## Standalone mode

### config.toml (top-level keys and sections)

```toml
app_name = "Authula"
base_url = "http://localhost:8080"      # base of API, redirects, emails
base_path = "/auth"                     # API prefix, e.g. "/api/auth"
# secret = ""                           # prefer AUTHULA_SECRET env var
disabled_paths = []                     # MUST be top-level, before any [table]

[database]
provider = "postgres"                   # sqlite | postgres | mysql
# url = "..."                           # prefer AUTHULA_DATABASE_URL
max_open_conns = 25
max_idle_conns = 5
conn_max_lifetime = "10m"

[logger]
level = "debug"                         # debug | info | warn | error

[session]
cookie_name = "authula.session_token"
expires_in = "24h"
update_age = "5m"
cookie_max_age = "24h"
secure = false                          # true in production (HTTPS)
http_only = true
same_site = "lax"                       # lax | strict | none
auto_cleanup = true
cleanup_interval = "1m"
max_sessions_per_user = 5               # 0 = unlimited

[verification]
auto_cleanup = true
cleanup_interval = "1m"

[security]
trusted_origins = ["http://localhost:3000"]
trusted_headers = []                    # e.g. ["X-Real-IP", "CF-Connecting-IP"]
trusted_proxies = []
[security.cors]
allow_credentials = true
allowed_origins = ["http://localhost:3000"]   # exact origins when using sessions
allowed_methods = ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
allowed_headers = ["Authorization", "Content-Type", "Set-Cookie", "Cookie"]
exposed_headers = []
max_age = "24h"

[event_bus]
prefix = ""
max_concurrent_handlers = 100
context_timeout = "5s"
provider = "sqlite"   # gochannel | sqlite | postgres | redis | kafka | nats | rabbitmq
[event_bus.go_channel]
buffer_size = 100

[plugins]
# per-plugin tables go here, see references/plugins.md
```

Event bus notes: `gochannel` is in-memory (events lost on restart, fine for dev). Redis, Kafka, NATS, RabbitMQ enable distributed handling across instances. For sqlite the file path is set in `[event_bus.sqlite] db_path`.

### Environment variables (`.env`)

```
AUTHULA_CONFIG_PATH=       # optional; defaults to ./config.toml
AUTHULA_BASE_URL=          # optional; can live in config.toml
AUTHULA_SECRET=            # required if not in config.toml (openssl rand -hex 32)
AUTHULA_DATABASE_URL=      # required if not in config.toml
# sqlite:     auth.db
# postgres:   postgresql://user:pass@localhost:5432/authula?sslmode=disable
# mysql:      user:pass@tcp(localhost:3306)/authula
GO_ENV=development
PORT=8080

# OAuth2 plugin
DISCORD_CLIENT_ID= / DISCORD_CLIENT_SECRET=
GITHUB_CLIENT_ID=  / GITHUB_CLIENT_SECRET=
GOOGLE_CLIENT_ID=  / GOOGLE_CLIENT_SECRET=
# Depending on providers used
POSTGRES_URL= REDIS_URL= KAFKA_BROKERS= NATS_URL= RABBITMQ_URL= EVENT_BUS_CONSUMER_GROUP=
```

### Docker

```bash
docker container run -itd \
  -p 8080:8080 \
  -v ./config.toml:/home/appuser/config.toml \
  --env-file <path-to-env-file> \
  ghcr.io/authula/authula:v1.46.0
```

Point `--env-file` at the `.env` file you created (for example `./.env`) and pin the image tag (or digest) you intend to run.

## Library mode

```bash
go get github.com/Authula/authula
```

```go
import (
    authula "github.com/Authula/authula"
    authulaconfig "github.com/Authula/authula/config"
    authulamodels "github.com/Authula/authula/models"
)

config := authulaconfig.NewConfig(
    // Same options as config.toml, as functional options. Known ones include:
    // WithAppName, WithBasePath, WithDatabase, WithLogger, WithSession,
    // WithVerification, WithSecurity, WithEventBus, WithRouteMappings,
    // WithCoreServiceHooks. Full list: config/options.go in the repo.
)
auth := authula.New(&authula.AuthConfig{
    Config:  config,
    Plugins: []authulamodels.Plugin{ /* plugin instances */ },
})
```

Database (`Provider`: `"sqlite"`, `"postgres"`, `"mysql"`):

```go
authulaconfig.WithDatabase(authulamodels.DatabaseConfig{Provider: "sqlite", URL: "auth.db"})
```

Migrations: core and plugin migrations run automatically when `authula.New` is called (custom migrator on Bun + raw SQL). Core tables are documented at https://www.authula.dev/docs/reference/database-schema. The `Auth` type also exposes `RunCoreMigrations`, `DropCoreMigrations`, `Migrator()`, `MigrationManager()`.

Mounting the handler (net/http):

```go
log.Fatal(http.ListenAndServe(":8080", auth.Handler()))
// or under a prefix:
http.Handle("/api/auth/", auth.Handler())
```

For Chi, Echo, Fiber see `references/extending.md`.

## Working example: email/password + sessions

Standalone additions on top of the base config:

```toml
[plugins.email]
enabled = true
provider = "smtp"
from_address = "noreply@example.com"

[plugins.email_password]
enabled = true
min_password_length = 8
max_password_length = 128
disable_sign_up = false
require_email_verification = true
auto_sign_in = true
send_email_on_sign_up = true
send_email_on_sign_in = false
email_verification_expires_in = "24h"
password_reset_expires_in = "1h"
request_email_change_expires_in = "1h"

[plugins.session]
enabled = true

[[route_mappings]]
paths = ["GET:/me", "POST:/sign-out"]
plugins = ["session.auth"]

[[route_mappings]]
paths = ["POST:/email-password/sign-in", "POST:/email-password/sign-up", "GET:/email-password/verify-email"]
plugins = ["session.auth.optional"]

[[route_mappings]]
paths = [
  "POST:/email-password/send-email-verification",
  "POST:/email-password/request-password-reset",
  "POST:/email-password/change-password",
  "POST:/email-password/request-email-change",
]
plugins = ["session.auth"]
```

SMTP env vars used by the email plugin in the example: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`.

Library-mode equivalent: pass `emailplugin.New(emailplugintypes.EmailPluginConfig{Enabled: true, Provider: emailplugintypes.ProviderSMTP, FromAddress: "..."})`, `emailpasswordplugin.New(emailpasswordplugintypes.EmailPasswordPluginConfig{...})` and `sessionplugin.New(sessionplugin.SessionPluginConfig{Enabled: true})` in `Plugins`, and the same mappings through `authulaconfig.WithRouteMappings(...)`.

## Testing with curl (base path `/api/auth`)

```bash
curl -X POST http://localhost:8080/api/auth/email-password/sign-up \
  -H "Content-Type: application/json" \
  --json '{"name":"John Doe","email":"john.doe@example.com","password":"password123"}'

curl -X POST http://localhost:8080/api/auth/email-password/sign-in \
  -H "Content-Type: application/json" \
  --json '{"email":"john.doe@example.com","password":"password123"}'

curl http://localhost:8080/api/auth/me -H "Cookie: authula.session_token=<token>"
```

`--json` requires curl 7.82+; on older curl use `-d` instead.

Sign-up body fields: `name`, `email`, `password`, optional `image`, `metadata`, `callback_url`.
