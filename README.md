# Short Code

Language / 语言: [English](README.md) | [简体中文](README.zh-CN.md)

A small HTTP short-link service backed by MySQL 8.0+.

## Run

Copy the example configuration and edit `configs/config.yaml` with your local MySQL credentials and server settings. The real configuration file is ignored by Git; only the example template is tracked.

```sh
cp -n configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml and set database.password.
go run ./cmd/server
```

The service initializes the schema file configured by `database.schema_file` on startup. The MySQL account must be able to create the configured database and its table. By default, the server listens on `:8080`.

To use a configuration file at another path, pass `-config`, for example `go run ./cmd/server -config /path/to/config.yaml`. Relative `database.schema_file` paths are resolved from the configuration file's directory.

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
