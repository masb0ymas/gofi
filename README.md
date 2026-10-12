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

| Variable                                 | Notes                                                                                                                                                                |
| ---------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ENV`                                    | `development`, `staging` or `production`. Unset means `production`. Only `development` drops `Secure` from the session cookie and exempts loopback from rate limits. |
| `DEBUG`                                  | Logs every query. Refused in production: auth queries are logged with their values.                                                                                  |
| `APP_SECRET`                             | At least 32 characters (`openssl rand -base64 32`).                                                                                                                  |
| `DB_DSN`                                 | PostgreSQL URL. `AUTHULA_DATABASE_URL` must not be set.                                                                                                              |
| `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS` | The whole replica's budget: Authula's pool takes half, the API's the rest.                                                                                           |
| `TRUSTED_PROXIES`                        | IPs or CIDRs of the reverse proxy. The client IP (rate limits) is the rightmost `X-Forwarded-For` entry that is not one of them.                                     |
| `MIGRATE_ON_BOOT`                        | Apply the embedded migrations before serving.                                                                                                                        |
| `SEED_ON_BOOT`                           | Run the seeders before serving; needs `SUPER_USER`.                                                                                                                  |
| `SUPER_USER`                             | `email:password`, password of 12 to 128 characters.                                                                                                                  |
| `SENTRY_DSN`                             | Optional; panics are reported when set.                                                                                                                              |

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
password differs, the next seeded boot resets it and ends every session of
that account, so rotating a leaked password locks the holder out. An account
someone registered with that email before the first seed is taken over the
same way: verified, reset and signed out. `cmd/migrate --seed` does not create
the super user, because that needs a running Authula.

## Authentication

Authula is mounted at `/v1/auth`. Requests take a JSON body.

| Sign-in          | Endpoint that returns the session and tokens  |
| ---------------- | --------------------------------------------- |
| Email & password | `POST /v1/auth/email-password/sign-in`        |
| Magic link       | `POST /v1/auth/magic-link/exchange`           |
| Google           | the OAuth callback, when it answers with JSON |

A successful sign-in sets the session cookie and returns `access_token`,
`refresh_token`, `expires_at` and `expires_in` next to `user` and `session`.
The API accepts either the cookie or `Authorization: Bearer <access_token>`.

- `POST /v1/auth/email-password/sign-up` creates the account and sends the
  verification email; it does not sign in. Password sign-in is refused, with
  the same 401 as a wrong password, until the email is verified, either from
  that email or by signing in with a magic link.
- The access token lasts 15 minutes and the refresh token 24 hours; renew
  them with `POST /v1/auth/token/refresh` and `{"refresh_token": "..."}`.
- `POST /v1/auth/token` returns a new pair for the session in the cookie. It
  refuses a bearer token, so a stolen access token cannot outlive its 15
  minutes. Use it when a session has no tokens: after a Google sign-in started
  with `redirect_to`, which answers with a redirect, or once the refresh token
  has expired while the 7-day session is still alive. A client without the
  cookie signs in again once its refresh token expires.
- Signing out revokes the tokens of that session, and an access token stops
  working when its session expires. A password reset ends every session of
  the account.
- A state-changing request that a browser sends from another origin is
  refused with 403 unless the origin is in `CORS_ALLOWED_ORIGINS`. Clients
  that are not browsers are unaffected.

## Idempotent writes

A write to `/v1/categories` accepts an idempotency key, so a retried request
does not apply twice. Send the key on the first attempt and repeat it verbatim
on every retry:

```sh
curl -X POST http://localhost:8080/v1/categories \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -H "X-Idempotency-Key: 8e03978e-40d5-43e8-bc93-6894a57f9324" \
  -d '{"name":"Books"}'
```

- The key is a UUID. `Idempotency-Key`, the name the IETF draft standardises,
  is accepted as well. A missing key means the write runs normally, and a key
  that is not a UUID is refused with 400.
- The first request runs the handler and stores the answer for 30 minutes; a
  retry within that window is answered from the store, with the same status and
  body and `Idempotent-Replayed: true` (browser clients read it because it is
  in the CORS expose list). A new key starts a new write.
- The key is scoped to the caller and the route, so the same key cannot replay
  another user's response or a response from a different endpoint.
- Only successful responses are stored. A retry after 401, 403, 404, 409 or 422
  runs the handler again, so a failed write stays retryable.
- Safe methods (`GET`, `HEAD`, `OPTIONS`) are never deduplicated, with or
  without a key.

`ponytail`: the store and its lock are in this process, so each replica
deduplicates on its own. Give the middleware shared storage (Redis) when there
is more than one instance; the keys are already opaque, so nothing else changes.

## Docker

```sh
docker build -t gofi .
docker run --env-file .env -p 8080:8080 gofi
```

The image holds two binaries, `/app/api` (the entrypoint) and `/app/migrate`:

```sh
docker run --env-file .env --entrypoint /app/migrate gofi up
```

## Release

Releases are cut from `main` with [release-it](https://github.com/release-it/release-it),
run by GitHub Actions, in three steps:

1. **Release Prepare** (`release-prepare.yml`), started by hand from the
   Actions tab. It runs the checks below, then release-it bumps the version in
   `package.json` and updates `CHANGELOG.md`, and the workflow opens a
   `chore: release vX.Y.Z` pull request from a `release/vX.Y.Z` branch.
   `increment: auto` derives the bump from the conventional commits.
2. **Release Publish** (`release-publish.yml`) runs when that pull request is
   merged into `main`. release-it tags the merge commit, creates the GitHub
   release and the workflow starts the image build.
3. **Release** (`release.yml`) builds the image for the tag and pushes it to
   `ghcr.io/masb0ymas/gofi` as `vX.Y.Z`, `vX.Y` and `latest`. Run it by hand to
   publish an image of any ref under its short commit SHA.

`package.json` exists only to hold the version and the release-it commands;
it has no dependencies. The repository setting "Allow GitHub Actions to create
and approve pull requests" must be on for step 1.

## Checks

```sh
go test ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
```

`cmd/api` has an integration test that boots the API against PostgreSQL 18:
migrations up and down, the seeders, the users mirror, sign-up and sign-in,
tokens, the role check, session revocation and the rate limit behind a proxy.
It is skipped unless `TEST_DB_DSN` points at a database it may wipe:

```sh
docker run -d --rm --name gofi-test-pg -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gofi_test -p 55432:5432 postgres:18
TEST_DB_DSN='postgres://postgres:postgres@localhost:55432/gofi_test?sslmode=disable' go test ./cmd/api
```

The release workflow runs all of it, the integration test included against a
`postgres:18` service, before it opens the release pull request; there is no
workflow that runs on every push. Its third-party action is pinned to a
commit and govulncheck to a version, because the job can write to the
repository.
