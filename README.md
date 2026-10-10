# Short Code

A small HTTP short-link service backed by MySQL 8.0+.

## Run

The service initializes the schema in `db/init.sql` on startup. The MySQL account must be able to create the `short_url` database and its table.

```sh
export MYSQL_USER=root
export MYSQL_PASSWORD='your-password'
export MYSQL_ADDR=127.0.0.1:3306
export HTTP_ADDR=:8080
go run ./cmd/server
```

Alternatively, set `SHORT_CODE_DB_DSN` to a complete MySQL DSN. `SHORT_CODE_SCHEMA` can point to the schema file if the process is not started from the repository root.

## API

Create a short link; `short_code` is optional and generated when omitted. `expires_at` is an optional RFC 3339 timestamp.

```sh
curl -i -X POST http://localhost:8080/api/short-links \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://example.com/a/long/path"}'
```

The response contains `short_code` and `short_url`. Visiting the returned URL redirects with HTTP 302. Disabled, deleted, expired, and unknown links return 404. `GET /healthz` returns `200 OK` when the HTTP process is running.

Run tests with:

```sh
go test ./...
```
