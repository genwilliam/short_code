# Short Code 短链接服务

语言 / Language: [English](README.md) | [简体中文](README.zh-CN.md)

一个基于 Go 和 MySQL 8.0+ 的 HTTP 短链接服务。

## 环境要求

- Go 1.25.10 或更高版本
- MySQL 8.0 或更高版本
- MySQL 账号需要有创建 `short_url` 数据库和数据表的权限。服务启动时会执行 `db/init.sql` 初始化结构。

## 启动服务

在项目根目录配置数据库和监听地址后运行：

```sh
export MYSQL_USER=root
export MYSQL_PASSWORD='你的数据库密码'
export MYSQL_ADDR=127.0.0.1:3306
export HTTP_ADDR=:8080
go run ./cmd/server
```

也可以使用完整 DSN 配置数据库连接；设置后 `SHORT_CODE_DB_DSN` 优先于 `MYSQL_USER`、`MYSQL_PASSWORD` 和 `MYSQL_ADDR`：

```sh
export SHORT_CODE_DB_DSN='user:password@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local'
go run ./cmd/server
```

默认情况下，服务从项目根目录读取 `db/init.sql`。如果从其他目录启动，可通过 `SHORT_CODE_SCHEMA` 指定 SQL 文件路径。

## API

### 创建短链接

`POST /api/short-links` 接收 JSON。`original_url` 必填；`short_code` 可选，省略时会生成 8 位短码；`expires_at` 可选，格式为 RFC 3339，且必须是未来时间。

```sh
curl -i -X POST http://localhost:8080/api/short-links \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://example.com/a/long/path"}'
```

指定短码和过期时间的示例：

```sh
curl -i -X POST http://localhost:8080/api/short-links \
  -H 'Content-Type: application/json' \
  -d '{"short_code":"example1","original_url":"https://example.com/page","expires_at":"2030-01-01T00:00:00Z"}'
```

短码允许 4 到 32 个英文字母、数字、连字符或下划线。成功时返回 `201 Created`，响应示例：

```json
{
  "short_code": "example1",
  "short_url": "http://localhost:8080/example1"
}
```

无效请求返回 `400 Bad Request`；短码已存在返回 `409 Conflict`。

### 访问短链接

访问 `GET /{short_code}` 会以 `302 Found` 重定向到原始 URL。短码不存在、已停用、已删除或已过期时返回 `404 Not Found`。

### 健康检查

`GET /healthz` 在 HTTP 服务运行时返回 `200 OK` 和 `ok`。该端点只检查进程是否响应，不检查数据库健康状态。

## 测试

```sh
go test ./...
```

## 部署提示

创建接口目前没有身份验证。公开部署前，建议在服务前增加鉴权和限流，避免未授权创建短链接。
