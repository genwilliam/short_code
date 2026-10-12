# Short Code 短链接服务

语言 / Language: [English](README.md) | [简体中文](README.zh-CN.md)

一个基于 Go 和 MySQL 8.0+ 的 HTTP 短链接服务。

## 环境要求

- Go 1.25.10 或更高版本
- MySQL 8.0 或更高版本
- MySQL 账号需要有创建 `short_url` 数据库和数据表的权限。服务启动时会执行 `db/init.sql` 初始化结构。

## 启动服务

先复制配置模板，再编辑 `configs/config.yaml`，填写本机 MySQL 用户、密码和服务监听地址。真实配置文件已加入 Git 忽略规则，只有不含个人密码的示例模板会被 Git 跟踪。

```sh
cp configs/config.example.yaml configs/config.yaml
# 编辑 configs/config.yaml，填写 database.password 等本机配置。
go run ./cmd/server
```

服务启动时会执行 `database.schema_file` 指定的 SQL 文件。MySQL 账号需要有创建配置中指定的数据库和数据表的权限。默认监听地址为 `:8080`。

也可以通过 `-config` 指定其他配置文件路径，例如：`go run ./cmd/server -config /path/to/config.yaml`。相对的 `database.schema_file` 路径以配置文件所在目录为基准解析。

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
