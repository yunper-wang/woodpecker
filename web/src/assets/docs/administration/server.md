---
toc_max_heading_level: 3
---

# Server

## Forge 与用户配置

Woodpecker 没有自己的用户注册功能。用户由你的 [forge](./12-forges/11-overview.md)（通过 OAuth2）提供。注册默认为关闭状态（`WOODPECKER_OPEN=false`）。如果注册开放，任何在已配置 forge 上拥有账号的用户均可登录 Woodpecker。

你也可以限制注册：

- 关闭注册，通过 CLI `woodpecker-cli user` 手动管理用户
- 开放注册，并通过 `WOODPECKER_ADMIN` 设置允许特定管理员用户

  ```ini
  WOODPECKER_OPEN=false
  WOODPECKER_ADMIN=john.smith,jane_doe
  ```

- 开放注册，并通过 `WOODPECKER_ORGS` 设置按组织归属过滤

  ```ini
  WOODPECKER_OPEN=true
  WOODPECKER_ORGS=dolores,dog-patch
  ```

管理员也应在配置中显式设置。

```ini
WOODPECKER_ADMIN=john.smith,jane_doe
```

## 仓库配置

Woodpecker 使用用户在 forge 上的 OAuth 权限进行操作。默认情况下，Woodpecker 会同步用户有权访问的所有仓库。使用 `WOODPECKER_REPO_OWNERS` 变量可以过滤只同步特定 GitHub 用户的仓库。通常应在此填写你公司的 GitHub 名称。

```ini
WOODPECKER_REPO_OWNERS=my_company,my_company_oss_github_user
```

## 数据库

Woodpecker 默认使用内嵌的 SQLite 数据库，无需安装或配置。但你可以将其替换为 MySQL/MariaDB 或 PostgreSQL 数据库。有几点基本事项需要注意：

- Woodpecker 不会自动创建数据库。如果使用 MySQL 或 Postgres 驱动，你需要手动使用 `CREATE DATABASE` 创建数据库。

- Woodpecker 不执行数据归档；这被认为超出了项目范围。Woodpecker 对存储的数据量相当保守，但你应预期数据库日志会显著增加数据库的体积。

- Woodpecker 自动处理数据库迁移，包括初始创建表和索引。新版本的 Woodpecker 会自动升级数据库，除非发行说明中另有说明。

- Woodpecker 不执行数据库备份。应由你所选数据库供应商提供的第三方工具来处理。

### SQLite

默认情况下，Woodpecker 使用存储在 `/var/lib/woodpecker/` 下的 SQLite 数据库。如果使用容器，可以挂载[数据卷](https://docs.docker.com/storage/volumes/#create-and-manage-volumes)来持久化 SQLite 数据库。

```diff title="docker-compose.yaml"
 services:
   woodpecker-server:
     [...]
+    volumes:
+      - woodpecker-server-data:/var/lib/woodpecker/
```

### MySQL/MariaDB

下面的示例演示了 MySQL 数据库配置。有关配置选项和示例，请参阅官方驱动[文档](https://github.com/go-sql-driver/mysql#dsn-data-source-name)。
MySQL/MariaDB 所需的最低版本由 `go-sql-driver/mysql` 决定——详情请参阅[其 README](https://github.com/go-sql-driver/mysql#requirements)。

```ini
WOODPECKER_DATABASE_DRIVER=mysql
WOODPECKER_DATABASE_DATASOURCE=root:password@tcp(1.2.3.4:3306)/woodpecker?parseTime=true
```

### PostgreSQL

下面的示例演示了 Postgres 数据库配置。有关配置选项和示例，请参阅官方驱动[文档](https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING)。
请使用 **11** 或更高版本的 Postgres。

```ini
WOODPECKER_DATABASE_DRIVER=postgres
WOODPECKER_DATABASE_DATASOURCE=postgres://root:password@1.2.3.4:5432/postgres?sslmode=disable
```

## TLS

Woodpecker 支持通过将证书挂载到容器中来配置 SSL。

```ini
WOODPECKER_SERVER_CERT=/etc/certs/woodpecker.example.com/server.crt
WOODPECKER_SERVER_KEY=/etc/certs/woodpecker.example.com/server.key
```

TLS 支持通过 Go 标准库的 [ListenAndServeTLS](https://golang.org/pkg/net/http/#ListenAndServeTLS) 函数提供。

### 容器配置

除了 [docker-compose](../05-installation/10-docker-compose.md) 安装中显示的端口外，还需要暴露 `443` 端口：

```diff title="docker-compose.yaml"
 services:
   woodpecker-server:
     [...]
     ports:
+      - 80:80
+      - 443:443
       - 9000:9000
```

此外，证书和密钥必须被挂载并引用：

```diff title="docker-compose.yaml"
 services:
   woodpecker-server:
     [...]
     environment:
+      - WOODPECKER_SERVER_CERT=/etc/certs/woodpecker.example.com/server.crt
+      - WOODPECKER_SERVER_KEY=/etc/certs/woodpecker.example.com/server.key
     volumes:
+      - /etc/certs/woodpecker.example.com/server.crt:/etc/certs/woodpecker.example.com/server.crt
+      - /etc/certs/woodpecker.example.com/server.key:/etc/certs/woodpecker.example.com/server.key
```

## 反向代理

### Apache

本指南简要介绍如何在 Apache2 Web 服务器后面安装 Woodpecker server。以下是一个示例配置：

<!-- cspell:ignore apacheconf -->

```apacheconf
ProxyPreserveHost On

RequestHeader set X-Forwarded-Proto "https"

ProxyPass / http://127.0.0.1:8000/
ProxyPassReverse / http://127.0.0.1:8000/
```

你必须安装以下 Apache 模块：

- `proxy`
- `proxy_http`

使用 https 时，你必须配置 Apache 设置 `X-Forwarded-Proto`。

```diff
 ProxyPreserveHost On

+RequestHeader set X-Forwarded-Proto "https"

 ProxyPass / http://127.0.0.1:8000/
 ProxyPassReverse / http://127.0.0.1:8000/
```

### Nginx

本指南简要介绍如何在 Nginx Web 服务器后面安装 Woodpecker server。更多高级配置选项请参阅 Nginx 官方[文档](https://docs.nginx.com/nginx/admin-guide)。

示例配置：

```nginx
server {
    listen 80;
    server_name woodpecker.example.com;

    location / {
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Host $http_host;

        proxy_pass http://127.0.0.1:8000;
        proxy_redirect off;
        proxy_http_version 1.1;
        proxy_buffering off;

        chunked_transfer_encoding off;
    }
}
```

你必须配置代理以设置 `X-Forwarded` 代理头：

```diff
 server {
     listen 80;
     server_name woodpecker.example.com;

     location / {
+        proxy_set_header X-Forwarded-For $remote_addr;
+        proxy_set_header X-Forwarded-Proto $scheme;

         proxy_pass http://127.0.0.1:8000;
         proxy_redirect off;
         proxy_http_version 1.1;
         proxy_buffering off;

         chunked_transfer_encoding off;
     }
 }
```

### Caddy

本指南简要介绍如何在 [Caddy web server](https://caddyserver.com/) 后面安装 Woodpecker server。以下是一个 Caddyfile 代理配置示例：

```caddy
# expose WebUI and API
woodpecker.example.com {
  reverse_proxy woodpecker-server:8000
}

# expose gRPC
woodpecker-agent.example.com {
  reverse_proxy h2c://woodpecker-server:9000
}
```

### Tunnelmole

[Tunnelmole](https://github.com/robbie-cahill/tunnelmole-client) 是一个开源隧道工具。

首先[安装 tunnelmole](https://github.com/robbie-cahill/tunnelmole-client#installation)。

安装完成后，运行以下命令启动 tunnelmole：

```bash
tmole 8000
```

它将启动一个隧道并给出如下响应：

```bash
➜  ~ tmole 8000
http://bvdo5f-ip-49-183-170-144.tunnelmole.net is forwarding to localhost:8000
https://bvdo5f-ip-49-183-170-144.tunnelmole.net is forwarding to localhost:8000
```

将 `WOODPECKER_HOST` 设置为 Tunnelmole URL（`xxx.tunnelmole.net`）并启动 server。

### Ngrok

[Ngrok](https://ngrok.com/) 是一个流行的闭源隧道工具。安装 ngrok 后，打开新的控制台并运行以下命令：

```bash
ngrok http 8000
```

将 `WOODPECKER_HOST` 设置为 ngrok URL（通常为 xxx.ngrok.io）并启动 server。

### Traefik

要在 [Traefik](https://traefik.io/) 负载均衡器后面安装 Woodpecker server，必须同时暴露 `http` 和 `gRPC` 端口。以下是一个完整示例，假设你在 docker swarm 中运行 Traefik，并希望进行 TLS 终止以及从 http 到 https 的自动重定向。

<!-- cspell:words redirectscheme certresolver  -->

```yaml
services:
  server:
    image: woodpeckerci/woodpecker-server:latest
    environment:
      - WOODPECKER_OPEN=true
      - WOODPECKER_ADMIN=your_admin_user
      # other settings ...

    networks:
      - dmz # externally defined network, so that traefik can connect to the server
    volumes:
      - woodpecker-server-data:/var/lib/woodpecker/

    deploy:
      labels:
        - traefik.enable=true

        # web server
        - traefik.http.services.woodpecker-service.loadbalancer.server.port=8000

        - traefik.http.routers.woodpecker-secure.rule=Host(`ci.example.com`)
        - traefik.http.routers.woodpecker-secure.tls=true
        - traefik.http.routers.woodpecker-secure.tls.certresolver=letsencrypt
        - traefik.http.routers.woodpecker-secure.entrypoints=web-secure
        - traefik.http.routers.woodpecker-secure.service=woodpecker-service

        - traefik.http.routers.woodpecker.rule=Host(`ci.example.com`)
        - traefik.http.routers.woodpecker.entrypoints=web
        - traefik.http.routers.woodpecker.service=woodpecker-service

        - traefik.http.middlewares.woodpecker-redirect.redirectscheme.scheme=https
        - traefik.http.middlewares.woodpecker-redirect.redirectscheme.permanent=true
        - traefik.http.routers.woodpecker.middlewares=woodpecker-redirect@docker

        #  gRPC service
        - traefik.http.services.woodpecker-grpc.loadbalancer.server.port=9000
        - traefik.http.services.woodpecker-grpc.loadbalancer.server.scheme=h2c

        - traefik.http.routers.woodpecker-grpc-secure.rule=Host(`woodpecker-grpc.example.com`)
        - traefik.http.routers.woodpecker-grpc-secure.tls=true
        - traefik.http.routers.woodpecker-grpc-secure.tls.certresolver=letsencrypt
        - traefik.http.routers.woodpecker-grpc-secure.entrypoints=web-secure
        - traefik.http.routers.woodpecker-grpc-secure.service=woodpecker-grpc

        - traefik.http.routers.woodpecker-grpc.rule=Host(`woodpecker-grpc.example.com`)
        - traefik.http.routers.woodpecker-grpc.entrypoints=web
        - traefik.http.routers.woodpecker-grpc.service=woodpecker-grpc

        - traefik.http.middlewares.woodpecker-grpc-redirect.redirectscheme.scheme=https
        - traefik.http.middlewares.woodpecker-grpc-redirect.redirectscheme.permanent=true
        - traefik.http.routers.woodpecker-grpc.middlewares=woodpecker-grpc-redirect@docker

volumes:
  woodpecker-server-data:
    driver: local

networks:
  dmz:
    external: true
```

## 监控指标

### 端点

Woodpecker 兼容 Prometheus，如果设置了环境变量 `WOODPECKER_PROMETHEUS_AUTH_TOKEN`，则会暴露 `/metrics` 端点。请注意，访问指标端点受到限制，需要上述环境变量中的授权 token。

```yaml
global:
  scrape_interval: 60s

scrape_configs:
  - job_name: 'woodpecker'
    bearer_token: dummyToken...

    static_configs:
      - targets: ['woodpecker.domain.com']
```

### 授权

管理员需要生成用户 API token，并在 Prometheus 配置文件中将其配置为 bearer token。请参阅以下示例：

```diff
 global:
   scrape_interval: 60s

 scrape_configs:
   - job_name: 'woodpecker'
+    bearer_token: dummyToken...

     static_configs:
        - targets: ['woodpecker.domain.com']
```

也可以从文件中读取 token：

```diff
 global:
   scrape_interval: 60s

 scrape_configs:
   - job_name: 'woodpecker'
+    bearer_token_file: /etc/secrets/woodpecker-monitoring-token

     static_configs:
        - targets: ['woodpecker.domain.com']
```

## UI 定制

Woodpecker 支持自定义 JS 和 CSS 文件。这些文件必须存在于 server 的文件系统中。
它们可以内置在 Docker 镜像中，也可以在 Kubernetes 环境中从 ConfigMap 挂载。
这两个配置变量相互独立，可以只存在其中一个文件，也可以两者都存在。

```ini
WOODPECKER_CUSTOM_CSS_FILE=/usr/local/www/woodpecker.css
WOODPECKER_CUSTOM_JS_FILE=/usr/local/www/woodpecker.js
```

以下示例展示如何在 Woodpecker 顶部导航栏中放置横幅消息。

```css title="woodpecker.css"
.banner-message {
  position: absolute;
  width: 280px;
  height: 40px;
  margin-left: 240px;
  margin-top: 5px;
  padding-top: 5px;
  font-weight: bold;
  background: red no-repeat;
  text-align: center;
}
```

```javascript title="woodpecker.js"
// place/copy a minified version of your preferred lightweight JavaScript library here ...
!(function () {
  'use strict';
  function e() {} /*...*/
})();

$().ready(function () {
  $('.app nav img').first().htmlAfter("<div class='banner-message'>This is a demo banner message :)</div>");
});
```

## 环境变量

### LOG_LEVEL

- Name: `WOODPECKER_LOG_LEVEL`
- Default: `info`

配置日志级别。可选值为 `trace`、`debug`、`info`、`warn`、`error`、`fatal`、`panic`、`disabled` 以及空值。

---

### LOG_FILE

- Name: `WOODPECKER_LOG_FILE`
- Default: `stderr`

日志输出目标。可使用 `stdout` 和 `stderr` 作为特殊关键字。

---

### DATABASE_LOG

- Name: `WOODPECKER_DATABASE_LOG`
- Default: `false`

启用数据库引擎（当前为 xorm）的日志记录。

---

### DATABASE_LOG_SQL

- Name: `WOODPECKER_DATABASE_LOG_SQL`
- Default: `false`

启用 SQL 命令日志记录。

---

### DATABASE_MAX_CONNECTIONS

- Name: `WOODPECKER_DATABASE_MAX_CONNECTIONS`
- Default: `100`

xorm 允许创建的最大数据库连接数。

---

### DATABASE_IDLE_CONNECTIONS

- Name: `WOODPECKER_DATABASE_IDLE_CONNECTIONS`
- Default: `2`

xorm 保持打开的数据库连接数。

---

### DATABASE_CONNECTION_TIMEOUT

- Name: `WOODPECKER_DATABASE_CONNECTION_TIMEOUT`
- Default: `3 Seconds`

活跃数据库连接允许保持打开的时间。

---

### DEBUG_PRETTY

- Name: `WOODPECKER_DEBUG_PRETTY`
- Default: `false`

启用格式化的 debug 输出。

---

### DEBUG_NOCOLOR

- Name: `WOODPECKER_DEBUG_NOCOLOR`
- Default: `true`

禁用彩色 debug 输出。

---

### HOST

- Name: `WOODPECKER_HOST`
- Default: none

Server 面向用户的完整 URL，包含主机名、端口（如果不是 HTTP/HTTPS 的默认端口）和路径前缀。

示例：

- `WOODPECKER_HOST=http://woodpecker.example.org`
- `WOODPECKER_HOST=http://example.org/woodpecker`
- `WOODPECKER_HOST=http://example.org:1234/woodpecker`

---

### SERVER_ADDR

- Name: `WOODPECKER_SERVER_ADDR`
- Default: `:8000`

配置 HTTP 监听端口。

---

### SERVER_ADDR_TLS

- Name: `WOODPECKER_SERVER_ADDR_TLS`
- Default: `:443`

启用 SSL 时配置 HTTPS 监听端口。

---

### SERVER_CERT

- Name: `WOODPECKER_SERVER_CERT`
- Default: none

server 用于接受 HTTPS 请求的 SSL 证书路径。

示例：`WOODPECKER_SERVER_CERT=/path/to/cert.pem`

---

### SERVER_KEY

- Name: `WOODPECKER_SERVER_KEY`
- Default: none

server 用于接受 HTTPS 请求的 SSL 证书密钥路径。

示例：`WOODPECKER_SERVER_KEY=/path/to/key.pem`

---

### CUSTOM_CSS_FILE

- Name: `WOODPECKER_CUSTOM_CSS_FILE`
- Default: none

server 用于提供自定义 .CSS 文件的路径，用于 UI 定制。
可用于显示横幅消息、Logo 或特定环境提示（即白标定制）。
文件必须采用 UTF-8 编码，以确保所有特殊字符得以保留。

示例：`WOODPECKER_CUSTOM_CSS_FILE=/usr/local/www/woodpecker.css`

---

### CUSTOM_JS_FILE

- Name: `WOODPECKER_CUSTOM_JS_FILE`
- Default: none

server 用于提供自定义 .JS 文件的路径，用于 UI 定制。
可用于显示横幅消息、Logo 或特定环境提示（即白标定制）。
文件必须采用 UTF-8 编码，以确保所有特殊字符得以保留。

示例：`WOODPECKER_CUSTOM_JS_FILE=/usr/local/www/woodpecker.js`

---

### GRPC_ADDR

- Name: `WOODPECKER_GRPC_ADDR`
- Default: `:9000`

配置 gRPC 监听端口。

---

### GRPC_SECRET

- Name: `WOODPECKER_GRPC_SECRET`
- Default: `secret`

配置 gRPC JWT secret。

---

### GRPC_SECRET_FILE

- Name: `WOODPECKER_GRPC_SECRET_FILE`
- Default: none

从指定文件路径读取 `WOODPECKER_GRPC_SECRET` 的值。

---

### METRICS_SERVER_ADDR

- Name: `WOODPECKER_METRICS_SERVER_ADDR`
- Default: none

配置一个不受保护的 metrics 端点。空值表示完全禁用 metrics 端点。

示例：`:9001`

---

### ADMIN

- Name: `WOODPECKER_ADMIN`
- Default: none

逗号分隔的管理员账号列表。

示例：`WOODPECKER_ADMIN=user1,user2`

---

### ORGS

- Name: `WOODPECKER_ORGS`
- Default: none

逗号分隔的已批准组织列表。

示例：`org1,org2`

---

### REPO_OWNERS

- Name: `WOODPECKER_REPO_OWNERS`
- Default: none

这些所有者的仓库将被允许在 Woodpecker 中使用。

示例：`user1,user2`

---

### OPEN

- Name: `WOODPECKER_OPEN`
- Default: `false`

启用后允许用户注册。

---

### AUTHENTICATE_PUBLIC_REPOS

- Name: `WOODPECKER_AUTHENTICATE_PUBLIC_REPOS`
- Default: `false`

即使仓库是公开的也始终使用身份验证来克隆仓库。如果 forge 要求始终进行身份验证（许多公司使用此方式），则需要启用此选项。

---

### DEFAULT_ALLOW_PULL_REQUESTS

- Name: `WOODPECKER_DEFAULT_ALLOW_PULL_REQUESTS`
- Default: `true`

仓库允许 pull request 的默认设置。

---

### DEFAULT_APPROVAL_MODE

- Name: `WOODPECKER_DEFAULT_APPROVAL_MODE`
- Default: `forks`

仓库审批模式的默认设置。可选值：`none`、`forks`、`pull_requests` 或 `all_events`。

---

### DEFAULT_CANCEL_PREVIOUS_PIPELINE_EVENTS

- Name: `WOODPECKER_DEFAULT_CANCEL_PREVIOUS_PIPELINE_EVENTS`
- Default: `pull_request, push`

当针对相同上下文（tag、branch）创建新 pipeline 时，将被取消的事件名称列表。

---

### DEFAULT_CLONE_PLUGIN

- Name: `WOODPECKER_DEFAULT_CLONE_PLUGIN`
- Default: `docker.io/woodpeckerci/plugin-git`

克隆仓库时默认使用的 Docker 镜像。

该镜像也会被添加到受信任的 clone plugin 列表中。

### DEFAULT_WORKFLOW_LABELS

- Name: `WOODPECKER_DEFAULT_WORKFLOW_LABELS`
- Default: none

你可以指定默认的 label/platform 条件，用于未设置 label 条件的 workflow 的 agent 选择。

示例：`platform=linux/amd64,backend=docker`

### DEFAULT_PIPELINE_TIMEOUT

- Name: `WOODPECKER_DEFAULT_PIPELINE_TIMEOUT`
- Default: 60

pipeline 被终止前的默认时间（分钟）。

### MAX_PIPELINE_TIMEOUT

- Name: `WOODPECKER_MAX_PIPELINE_TIMEOUT`
- Default: 120

可在仓库设置中设置的 pipeline 被终止前的最长时间（分钟）。

---

### SESSION_EXPIRES

- Name: `WOODPECKER_SESSION_EXPIRES`
- Default: `72h`

配置 session 过期时间。
说明：当用户登录 Woodpecker 时，会创建一个临时 session token。
在 session 有效期间（直到过期或退出登录），用户无需重新认证即可登录 Woodpecker。

### PLUGINS_PRIVILEGED

- Name: `WOODPECKER_PLUGINS_PRIVILEGED`
- Default: none

以特权模式运行的 Docker 镜像。除非你清楚自己在做什么，否则请勿修改！

建议同时指定镜像标签，以强制精确匹配。

### PLUGINS_TRUSTED_CLONE

- Name: `WOODPECKER_PLUGINS_TRUSTED_CLONE`
- Default: `docker.io/woodpeckerci/plugin-git,docker.io/woodpeckerci/plugin-git,quay.io/woodpeckerci/plugin-git`

受信任的可在 clone step 中处理 Git 凭证信息的 plugin。
如果 clone step 使用的镜像不在此列表中，Git 凭证将不会被注入，用户需要使用其他方式（例如
