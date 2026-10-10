# GoFi

A Go API starter on Fiber and PostgreSQL, with authentication by
[Authula](https://github.com/Authula/authula) embedded as a library.

## Run

Requires Go 1.27.2+ and PostgreSQL 18 (the migrations use `uuidv7()`).

```sh
cp .env.example .env   # then fill in APP_SECRET, RESEND_API_KEY, SUPER_USER
make run
```

`make` exports `.env` to the binaries; every setting is an environment
variable (or the flag of the same name, see `go run ./cmd/api -h`). Nothing
secret is passed as a process argument.

| Variable | Notes |
|---|---|
| `ENV` | `development`, `staging` or `production`. Unset means `production`. Only `development` drops `Secure` from the session cookie and exempts loopback from rate limits. |
| `DEBUG` | Logs every query. Refused in production: auth queries are logged with their values. |
| `APP_SECRET` | At least 32 characters (`openssl rand -base64 32`). |
| `DB_DSN` | PostgreSQL URL. `AUTHULA_DATABASE_URL` must not be set. |
| `MIGRATE_ON_BOOT` | Apply the embedded migrations before serving. |
| `SEED_ON_BOOT` | Run the seeders before serving; needs `SUPER_USER`. |
| `SUPER_USER` | `email:password`, password of 12 to 128 characters. |
| `SENTRY_DSN` | Optional; panics are reported when set. |

## Database

```sh
make db/migrations/up            # apply migrations
make db/migrations/up/seed       # ... and seed the categories
make db/migrations/force version=N   # after a failed migration
make db/migrations/refresh       # roll back everything and reapply
```

`down` and `refresh` delete all data, users included, so they refuse a
database that is not on localhost unless `--allow-remote` is passed to
`cmd/migrate`.

Authula keeps its tables in the `authula` schema and creates them on boot. `public.users` mirrors `authula.users` (id, name, email, image) and
adds the application `role`. The mirror is updated by a hook on every user
create and update and repaired on every boot.

## Super user

With `SEED_ON_BOOT=true` the API creates the `SUPER_USER` account through
Authula, with a verified email, and gives it the `super_user` role. Only that
role may create, update or delete categories; everyone can read them.

`SUPER_USER` is the source of truth for the password: if the account's
password differs, the next seeded boot resets it. `cmd/migrate --seed` does
not create the super user, because that needs a running Authula.

## Authentication

Authula is mounted at `/v1/auth`. Requests take a JSON body.

| Sign-in | Endpoint that returns the session and tokens |
|---|---|
| Email & password | `POST /v1/auth/email-password/sign-in` |
| Magic link | `POST /v1/auth/magic-link/exchange` |
| Google | the OAuth callback, when it answers with JSON |

A successful sign-in sets the session cookie and returns `access_token`,
`refresh_token`, `expires_at` and `expires_in` next to `user` and `session`.
The API accepts either the cookie or `Authorization: Bearer <access_token>`.

- The access token lasts 15 minutes; renew it with
  `POST /v1/auth/token/refresh` and `{"refresh_token": "..."}`.
- Signing out revokes the tokens of that session.
- Signing in again while a session is still active returns no new tokens; use
  the refresh endpoint.
- A Google sign-in started with `redirect_to` answers with a redirect, which
  carries the cookie but no tokens.

## Checks

```sh
go test ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

CI runs both on every push and pull request.
