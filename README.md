
# Chirpy

Chirpy is a small X-like http-server built with Go, and backed by PostgreSQL. It creates chirps and users using a RESTful API, and authenticates users using JWTs and refresh tokens.

## Setup

### Prerequisites

- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)

Optional

- [goose](https://github.com/pressly/goose) - if you want to migrate the sql files yourself
- [sqlc](https://github.com/sqlc-dev/sqlc) - generate go code from sql queries

Built with [jwt-go](https://github.com/golang-jwt/jwt) and [argon2id](https://github.com/alexedwards/argon2id) - these are pulled in automatically as Go module dependencies, no separate install needed.

### Installation

```bash
git clone https://github.com/reconfirmok/chirpy.git
```

### Configuration

Create a .env file on the root with these environment variables, **Don't forget to change their values**

```env
DB_URL=""
PLATFORM=""
JWT_SECRET=""
POLKA_KEY=""
```

- DB_URL: postgresql connection string e.g.
    
    ```
    postgres://username:password@host:port/database
    ```
    
- PLATFORM: set value to "dev" to enable critical endpoints, like deleting actions, that should never be invoked in a prod environment
  
- JWT_SECRET: the secret to encrypt your users tokens.
    
- POLKA_KEY: the API key from the webhook service.
    

## API Reference

### Health & admin

| Method | Path             | Description                                                                    |
| ------ | ---------------- | ------------------------------------------------------------------------------ |
| GET    | `/api/healthz`   | Check server running status                                                    |
| GET    | `/admin/metrics` | HTML page showing the fileserver hit count                                     |
| POST   | `/admin/reset`   | Delete server db, .env PLATFORM variable must be "dev" to invoke this endpoint |

### Users

|Method|Path|Description|
|---|---|---|
|POST|`/api/users`|Create a user, must provide email and password on the json body|
|PUT|`/api/users`|Update user data, Must provide access token on the header|
|POST|`/api/login`|Login using user credentials|
|POST|`/api/refresh`|Exchange a valid refresh token for a new access token|
|POST|`/api/revoke`|Revoke a refresh token

### Chirps

|Method|Path|Description|
|---|---|---|
|POST|`/api/chirps`|Create a chirp, must provide access token on the header|
|GET|`/api/chirps`|Get all chirps on the db, with optional queries, `author_id` to get chirps by id, `sort` to sort queries in ASC or DESC order|
|GET|`/api/chirps/{chirpID}`|Get a single chirp by id|
|DELETE|`/api/chirps/{chirpID}`|Delete a chirp, must provide access token on the header, only the author can delete their own chirp|

### Webhooks

|Method|Path|Description|
|---|---|---|
|POST|`/api/polka/webhooks`|3rd-party payment service to mark chirpy users as red "premium user"|


## Authentication

- **Access tokens** are JWTs (HS256), sent as `Authorization: Bearer <token>`, and expire after 1 hour.
- **Refresh tokens** are opaque random hex strings stored in the `refresh_tokens` table with a 60-day expiry. they're also sent as `Authorization: Bearer <token>` to `/api/refresh` and `/api/revoke`.
- **Webhook requests** authenticate with `Authorization: ApiKey <key>`, matched against `POLKA_KEY`.


## Notes

- Chirp validation rejects chirps over 140 characters and censors a small fixed list of words (replaced with `****`).
