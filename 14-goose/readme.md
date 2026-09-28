# 14 - Goose Migrations + Passwordless Magic-Link Auth

A REST API demonstrating [goose](https://github.com/pressly/goose) database migrations against Postgres (Neon), wired into a passwordless "magic link" authentication flow built with [gin](https://github.com/gin-gonic/gin).

## Prerequisites

- Go 1.27+
- The `goose` CLI:
  ```powershell
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```
- A Postgres connection string (this lesson uses a free [Neon](https://neon.tech) database)

## Setup

1. Copy your Postgres connection string and a random JWT secret into `.env` (already gitignored):
   ```
   GOOSE_DRIVER=postgres
   GOOSE_DBSTRING=postgresql://user:pass@host/dbname?sslmode=require

   APP_PORT=8080
   APP_BASE_URL=http://localhost:8080

   JWT_SECRET=<a long random string>
   ACCESS_TOKEN_TTL_MINUTES=15
   REFRESH_TOKEN_TTL_HOURS=168
   MAGIC_LINK_TTL_MINUTES=15

   COOKIE_SECURE=false
   COOKIE_DOMAIN=

   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=you@example.com
   SMTP_PASSWORD=<gmail app password>
   SMTP_FROM=you@example.com
   ```
   Generate a `JWT_SECRET` quickly with:
   ```powershell
   openssl rand -base64 48
   ```
   For `SMTP_PASSWORD` on Gmail/Google Workspace: enable 2-Step Verification on the account, then create one at [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords) and paste the 16-character code with no spaces. Alternatively, point `SMTP_*` at a local catcher like [MailHog](https://github.com/mailhog/MailHog) or [Mailtrap](https://mailtrap.io) during development.

2. Install Go dependencies:
   ```powershell
   go mod tidy
   ```

## Database schema

Three tables, defined in `migrations/`:

| Table | Purpose |
|---|---|
| `users` | `id`, `email`, `created_at` |
| `magic_link_tokens` | one-time sign-in tokens: `id`, `user_id`, `token_hash`, `expires_at`, `used_at`, `created_at` |
| `sessions` | refresh-token sessions: `id`, `user_id`, `refresh_token_hash`, `expires_at`, `created_at` |

Only **hashes** of tokens are ever stored — raw tokens exist only in the email link / client cookie, never in the database.

## Running migrations

`goose` reads `GOOSE_DRIVER` / `GOOSE_DBSTRING` straight out of `.env`, so you don't need to pass them on the command line.

```powershell
# Apply all pending migrations
goose -dir migrations up

# Check what's applied vs. pending
goose -dir migrations status

# Roll back the most recent migration
goose -dir migrations down

# Roll back everything
goose -dir migrations reset
```

To add a new migration:
```powershell
goose -dir migrations create add_something_to_users sql
```
This creates a timestamped file in `migrations/` with `-- +goose Up` / `-- +goose Down` sections for you to fill in. Run `goose -dir migrations up` afterward to apply it.

## Running the application

```powershell
go run .
```
You should see gin list all 6 routes and start listening on the port from `APP_PORT` (default `8080`).

## API

```
POST /auth/magic-link   { "email": "..." }   -> emails a one-time sign-in link
GET  /auth/verify?token=...                   -> read-only check that a token is still valid
POST /auth/exchange     { "token": "..." }    -> consumes the token, starts a session (sets cookies)
POST /auth/refresh                            -> rotates the session using the refresh cookie
POST /auth/logout                             -> ends the session, clears cookies
GET  /users/me                                -> returns the signed-in user (requires session cookies)
```

## Testing the flow end-to-end

Ready-made requests are in `api-test/*.http` (use a REST Client editor extension), or use `curl` with a cookie jar:

```powershell
# 1. Request a magic link
curl -X POST http://localhost:8080/auth/magic-link `
  -H "Content-Type: application/json" `
  -d '{\"email\":\"you@example.com\"}'

# 2. Check your inbox, copy the `token` query param from the link, then confirm it's valid
curl "http://localhost:8080/auth/verify?token=<TOKEN>"

# 3. Exchange it for a session (this is what actually logs you in)
curl -i -c cookies.txt -X POST http://localhost:8080/auth/exchange `
  -H "Content-Type: application/json" `
  -d '{\"token\":\"<TOKEN>\"}'

# 4. Call the protected route
curl -b cookies.txt http://localhost:8080/users/me

# 5. Rotate the session
curl -i -b cookies.txt -c cookies.txt -X POST http://localhost:8080/auth/refresh

# 6. Log out
curl -b cookies.txt -c cookies.txt -X POST http://localhost:8080/auth/logout
```

A magic-link token is single-use (re-exchanging it returns `400`), and refresh tokens rotate on every `/auth/refresh` call (the previous cookie stops working once a new one is issued).
