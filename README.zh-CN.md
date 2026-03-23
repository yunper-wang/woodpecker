# Woodpecker CI

<p align="center">
  <a href="https://github.com/woodpecker-ci/woodpecker/">
    <img alt="Woodpecker" src="docs/static/img/logo.svg" width="220"/>
  </a>
</p>

<p align="center">
  <a href="https://godoc.org/go.woodpecker-ci.org/woodpecker"><img src="https://godoc.org/go.woodpecker-ci.org/woodpecker?status.svg" alt="GoDoc"/></a>
  <a href="https://goreportcard.com/report/go.woodpecker-ci.org/woodpecker/v3"><img src="https://goreportcard.com/badge/go.woodpecker-ci.org/woodpecker/v3" alt="Go Report Card"/></a>
  <a href="https://github.com/woodpecker-ci/woodpecker/blob/main/LICENSE"><img src="https://img.shields.io/github/license/woodpecker-ci/woodpecker" alt="License"/></a>
</p>

Woodpecker 是一个基于 [Drone](https://github.com/drone/drone) 构建的社区驱动的持续集成（CI/CD）系统，轻量、可扩展、易于自托管。

> 「这只小鸟知道你的秘密」

---

## 目录

- [架构概览](#架构概览)
- [主要特性](#主要特性)
- [本分支新增功能](#本分支新增功能)
- [快速开始](#快速开始)
- [流水线配置](#流水线配置)
- [Matrix 矩阵构建](#matrix-矩阵构建)
- [Secret 密钥管理](#secret-密钥管理)
- [Artifact 产物存储](#artifact-产物存储)
- [审批门控](#审批门控)
- [插件系统](#插件系统)
- [支持的代码托管平台](#支持的代码托管平台)
- [执行后端](#执行后端)
- [安装部署](#安装部署)
- [构建与开发](#构建与开发)
- [测试](#测试)
- [项目结构](#项目结构)
- [许可证](#许可证)

---

## 架构概览

Woodpecker 由三个独立的可执行文件组成，共享同一个 monorepo：

```
┌─────────────────────────────────────────────────────────────────┐
│                         Woodpecker CI                           │
│                                                                 │
│   VCS Webhook ──► server/api ──► forge（解析 Webhook）          │
│                              ──► server/pipeline（编译调度）    │
│                              ──► gRPC ──► agent/runner          │
│                                          └──► backend（执行）   │
│                                          └──► gRPC 日志回传     │
└─────────────────────────────────────────────────────────────────┘
```

| 组件 | 路径 | 职责 |
|------|------|------|
| **server** | `cmd/server` | 提供 Web UI、REST API、处理 VCS Webhook、调度流水线、gRPC 服务端 |
| **agent** | `cmd/agent` | 通过 gRPC 向 server 轮询任务，通过 backend 执行流水线步骤，回传日志 |
| **cli** | `cmd/cli` | 面向用户的命令行客户端，支持触发流水线、查看日志、管理密钥等操作 |

### 核心包说明

| 包路径 | 功能 |
|--------|------|
| `server/forge/` | VCS 集成层，实现 `forge.Forge` 接口；子目录按平台区分：github、gitea、forgejo、gitlab、bitbucket |
| `server/model/` | 全部数据模型：Repo、Pipeline、Workflow、Step、User、Secret、Artifact、Approval 等 |
| `server/store/` | 数据库抽象层；`datastore/` 基于 XORM 实现 SQLite/PostgreSQL/MySQL 支持 |
| `server/pipeline/` | 流水线生命周期：编译、调度、事件处理、步骤更新、审批 |
| `server/api/` | REST API handlers（Gin 框架） |
| `server/router/` | HTTP 路由注册 |
| `server/rpc/` | gRPC 服务端实现（agent ↔ server 通信协议） |
| `pipeline/frontend/` | `.woodpecker.yml` 解析与编译，生成内部步骤图 |
| `pipeline/backend/` | 执行后端：docker、kubernetes、local |
| `agent/runner/` | agent 端流水线执行循环 |
| `web/src/` | Vue 3 + TypeScript 前端（Vite、Pinia、WindiCSS） |

---

## 主要特性

- **多平台 Forge 支持**：GitHub、Gitea、Forgejo、GitLab、Bitbucket、Bitbucket Data Center
- **YAML 流水线配置**：语法简洁，支持多 workflow 文件（`.woodpecker/` 目录）
- **Matrix 矩阵构建**：笛卡尔积自动展开，支持 `include`/`exclude` 精确控制
- **多执行后端**：Docker、Kubernetes、本地（local）；支持自定义 backend
- **丰富的插件生态**：通过容器镜像复用 CI 逻辑，支持 S3、Slack、Kubernetes 部署等
- **三级 Secret 管理**：仓库级 → 组织级 → 全局级，按需隔离
- **Registry 配置**：私有镜像仓库凭证管理
- **Cron 定时触发**：支持类 cron 表达式的周期性流水线
- **路径过滤（`when.path`）**：仅在特定文件变更时触发步骤
- **流水线审批门控**：可配置多级审批，防止不受信任的代码执行
- **Artifact 产物存储**：上传、下载、列出、删除流水线产物
- **徽章支持**：SVG 状态徽章和 CruiseControl XML 格式
- **SSE 日志流**：实时推送步骤日志到前端
- **Autoscaler**：按需启动云虚拟机执行构建，完成后自动销毁

---

## 本分支新增功能

本分支在上游 Woodpecker CI 基础上新增了以下功能与修复：

### 1. Matrix `exclude` 支持

在 matrix 定义中可使用 `exclude` 排除特定组合，精确控制需要构建的组合集：

```yaml
matrix:
  go_version:
    - "1.21"
    - "1.22"
  os:
    - linux
    - windows
  exclude:
    - go_version: "1.21"
      os: windows
```

上述配置将跳过 `go_version=1.21 + os=windows` 这一组合。超出上限时会返回明确的错误（而非静默截断）。

### 2. Matrix 上限可配置

通过环境变量或启动参数自定义 matrix 最大规模，适用于大型构建矩阵场景：

```bash
# 最大排列数，默认 25
WOODPECKER_MAX_MATRIX_AXIS=50

# 最大变量维度数，默认 10
WOODPECKER_MAX_MATRIX_TAGS=15
```

或通过启动参数传入：

```bash
woodpecker-server --max-matrix-axis 50 --max-matrix-tags 15
```

### 3. `when.path` 支持 `tag` 和 `release` 事件

路径过滤条件（`when.path`）现在对 `tag` 和 `release` 事件同样生效。GitHub、Gitea、Forgejo 三个 forge 均已同步支持，触发 tag/release 时会自动获取变更文件列表：

```yaml
steps:
  - name: deploy
    image: alpine
    commands:
      - ./deploy.sh
    when:
      event: [tag, release]
      path:
        include:
          - "src/**"
          - "go.mod"
        exclude:
          - "docs/**"
```

### 4. 流水线 Artifact 存储

工作流可以上传构建产物（artifact），供后续步骤、其他流水线或用户手动下载：

```bash
# 上传 artifact（需 push 权限）
POST /api/repos/{repo_id}/pipelines/{number}/artifacts?name=build.tar.gz&workflow_id=1

# 列出 artifact（需 pull 权限）
GET  /api/repos/{repo_id}/pipelines/{number}/artifacts

# 下载 artifact（需 pull 权限）
GET  /api/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id}

# 删除 artifact（需 push 权限）
DELETE /api/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id}
```

Artifact 数据存储于数据库，支持 SQLite/PostgreSQL/MySQL。

### 5. 多级审批门控

支持配置需要多名用户批准才能放行被阻塞的流水线。适用于高安全要求的部署流程：

在项目设置中配置 `required_approvals`（默认为 1）。当流水线处于 `blocked` 状态时，需要累计达到指定数量的不同用户批准后方可继续执行。

```bash
# 查看当前流水线审批记录
GET  /api/repos/{repo_id}/pipelines/{number}/approvals

# 提交审批（approve / decline）
POST /api/repos/{repo_id}/pipelines/{number}/approve
POST /api/repos/{repo_id}/pipelines/{number}/decline
```

每条审批记录包含：审批人用户名、用户 ID、操作类型（approve/decline）、时间戳。同一用户的重复审批只计一次。

### 6. Bug 修复：`when.status` 逻辑修复

修复了 `IncludesStatusSuccess` 函数中遗留的调试日志（`fmt.Println`）及重复调用 `c.Match()` 导致的副作用问题，使步骤状态过滤行为更加可靠。

---

## 快速开始

### 使用 Docker Compose 部署（推荐）

创建 `docker-compose.yaml`：

```yaml
services:
  woodpecker-server:
    image: woodpeckerci/woodpecker-server:v3
    ports:
      - 8000:8000
    volumes:
      - woodpecker-server-data:/var/lib/woodpecker/
    environment:
      - WOODPECKER_OPEN=true
      - WOODPECKER_HOST=${WOODPECKER_HOST}          # 如 https://ci.example.com
      - WOODPECKER_GITHUB=true
      - WOODPECKER_GITHUB_CLIENT=${WOODPECKER_GITHUB_CLIENT}
      - WOODPECKER_GITHUB_SECRET=${WOODPECKER_GITHUB_SECRET}
      - WOODPECKER_AGENT_SECRET=${WOODPECKER_AGENT_SECRET}

  woodpecker-agent:
    image: woodpeckerci/woodpecker-agent:v3
    command: agent
    restart: always
    depends_on:
      - woodpecker-server
    volumes:
      - woodpecker-agent-config:/etc/woodpecker
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - WOODPECKER_SERVER=woodpecker-server:9000
      - WOODPECKER_AGENT_SECRET=${WOODPECKER_AGENT_SECRET}

volumes:
  woodpecker-server-data:
  woodpecker-agent-config:
```

生成共享密钥：

```bash
openssl rand -hex 32
```

启动服务：

```bash
docker compose up -d
```

访问 `http://localhost:8000`，使用 VCS 账户登录。

### 配置 Forge（代码托管平台）

#### GitHub

在 GitHub 创建 OAuth App（Settings → Developer settings → OAuth Apps），回调地址填写 `https://<your-host>/authorize`，然后配置：

```bash
WOODPECKER_GITHUB=true
WOODPECKER_GITHUB_CLIENT=<Client ID>
WOODPECKER_GITHUB_SECRET=<Client Secret>
```

#### Gitea / Forgejo

```bash
WOODPECKER_GITEA=true
WOODPECKER_GITEA_URL=https://gitea.example.com
WOODPECKER_GITEA_CLIENT=<OAuth2 Client ID>
WOODPECKER_GITEA_SECRET=<OAuth2 Client Secret>
```

---

## 流水线配置

Woodpecker 在仓库的以下路径查找流水线配置（按优先级）：

1. `.woodpecker/` 目录下所有 `.yaml`/`.yml` 文件（支持多 workflow）
2. `.woodpecker.yaml`
3. `.woodpecker.yml`

### 基础示例

```yaml
# .woodpecker/ci.yaml
when:
  - event: push
    branch: main

steps:
  - name: build
    image: golang:1.22
    commands:
      - go build ./...
      - go test ./...

  - name: docker-build
    image: woodpeckerci/plugin-docker-buildx
    settings:
      repo: myorg/myimage
      tags: latest
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password
    when:
      event: push
      branch: main
```

### `when` 条件过滤

`when` 支持在 workflow 级别或 step 级别进行条件过滤：

```yaml
when:
  event: [push, pull_request]   # 事件类型
  branch: main                  # 分支名（支持 glob）
  status: [success, failure]    # 上一步状态
  path:                         # 文件路径变更（tag/release 事件同样支持）
    include:
      - "src/**"
      - "go.mod"
    exclude:
      - "docs/**"
```

### 并行步骤

使用 `depends_on` 控制步骤依赖关系，实现并行执行：

```yaml
steps:
  - name: build-frontend
    image: node:20
    commands:
      - npm ci && npm run build

  - name: build-backend
    image: golang:1.22
    commands:
      - go build ./...

  - name: deploy
    image: alpine
    commands:
      - ./deploy.sh
    depends_on:
      - build-frontend
      - build-backend
```

### 服务容器（Services）

```yaml
steps:
  - name: test
    image: golang:1.22
    commands:
      - go test -tags integration ./...
    environment:
      DATABASE_URL: postgres://postgres:password@postgres:5432/testdb

services:
  - name: postgres
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: password
      POSTGRES_DB: testdb
```

---

## Matrix 矩阵构建

Matrix 构建会为每种参数组合自动创建独立的流水线实例。

### 基础矩阵

```yaml
matrix:
  go_version:
    - "1.21"
    - "1.22"
  os:
    - linux
    - windows

steps:
  - name: test
    image: golang:${go_version}
    commands:
      - go test ./...
```

上述配置将生成 4 个并行流水线：`1.21+linux`、`1.21+windows`、`1.22+linux`、`1.22+windows`。

### include：指定特定组合

```yaml
matrix:
  include:
    - go_version: "1.21"
      os: linux
    - go_version: "1.22"
      os: linux
    - go_version: "1.22"
      os: windows
```

### exclude：排除特定组合（本分支新增）

```yaml
matrix:
  go_version:
    - "1.21"
    - "1.22"
  os:
    - linux
    - windows
  exclude:
    - go_version: "1.21"
      os: windows
```

### 配置上限（本分支新增）

```bash
# 最大排列数（默认 25）
WOODPECKER_MAX_MATRIX_AXIS=50

# 最大变量维度数（默认 10）
WOODPECKER_MAX_MATRIX_TAGS=15
```

---

## Secret 密钥管理

Woodpecker 提供三级密钥管理，优先级从低到高：

| 级别 | 作用域 | 管理权限 |
|------|--------|----------|
| 全局 Secret | 整个实例所有流水线 | 仅管理员 |
| 组织 Secret | 该组织所有仓库 | 组织管理员 |
| 仓库 Secret | 该仓库所有流水线 | 仓库拥有者 |

### 在流水线中使用 Secret

```yaml
steps:
  - name: deploy
    image: alpine
    environment:
      AWS_ACCESS_KEY:
        from_secret: aws_access_key
      AWS_SECRET_KEY:
        from_secret: aws_secret_key
    commands:
      - ./deploy.sh

  - name: notify
    image: woodpeckerci/plugin-slack
    settings:
      webhook:
        from_secret: slack_webhook
      message: "Build completed"
```

> 注意：默认情况下，Secret 不会暴露给 Pull Request 流水线，需在 Secret 设置中显式开启。

---

## Artifact 产物存储

本分支新增了流水线 Artifact 功能，支持上传、下载、管理构建产物。

### API 接口

```bash
# 上传 artifact（需 push 权限）
POST /api/repos/{repo_id}/pipelines/{number}/artifacts?name=build.tar.gz&workflow_id=1

# 列出 artifact（需 pull 权限）
GET  /api/repos/{repo_id}/pipelines/{number}/artifacts

# 下载 artifact（需 pull 权限）
GET  /api/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id}

# 删除 artifact（需 push 权限）
DELETE /api/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id}
```

### 数据模型

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 主键 |
| repo_id | int64 | 所属仓库 |
| pipeline_id | int64 | 所属流水线 |
| workflow_id | int64 | 所属 workflow |
| name | string | artifact 文件名 |
| file_size | int64 | 文件大小（字节） |
| data | []byte | 文件内容 |
| created_at | int64 | 创建时间（Unix 时间戳） |

---

## 审批门控

Woodpecker 支持对不受信任的代码（如 fork PR）要求审批后才执行流水线，防止密钥泄露和恶意代码运行。

### 配置方式

在项目设置（Project Settings → Require approval for）中配置审批策略：

- `None`：所有流水线直接执行
- `Forked repositories`（默认）：来自 fork 的 PR 需要审批
- `All repositories`：所有 PR 均需审批

### 多级审批（本分支新增）

在仓库配置中设置 `required_approvals` 字段，要求多名用户批准才能放行：

```bash
# 查看审批记录
GET  /api/repos/{repo_id}/pipelines/{number}/approvals

# 提交审批
POST /api/repos/{repo_id}/pipelines/{number}/approve

# 拒绝流水线
POST /api/repos/{repo_id}/pipelines/{number}/decline
```

同一用户的重复审批仅计一次，需要达到 `required_approvals` 指定的不同用户数量后流水线方可继续执行。

---

## 插件系统

Woodpecker 插件是普通的容器镜像，通过环境变量接收参数并执行特定任务。

```yaml
steps:
  - name: upload-to-s3
    image: woodpeckerci/plugin-s3
    settings:
      bucket: my-bucket
      access_key:
        from_secret: aws_access_key
      secret_key:
        from_secret: aws_secret_key
      source: dist/**/*
      target: /releases/

  - name: send-notification
    image: woodpeckerci/plugin-slack
    settings:
      webhook:
        from_secret: slack_webhook
      channel: "#ci"
```

常用官方插件：

| 插件 | 功能 |
|------|------|
| `woodpeckerci/plugin-s3` | S3 文件上传/下载 |
| `woodpeckerci/plugin-docker-buildx` | Docker 镜像构建与推送 |
| `woodpeckerci/plugin-slack` | Slack 通知 |
| `woodpeckerci/plugin-helm` | Kubernetes Helm 部署 |
| `woodpeckerci/plugin-git` | 高级 Git 操作 |
| `woodpeckerci/plugin-prettier` | 代码格式检查 |

更多插件参见 [Woodpecker 插件中心](https://woodpecker-ci.org/plugins)。

---

## 支持的代码托管平台

| Forge | 说明 |
|-------|------|
| GitHub | 支持 GitHub.com 和 GitHub Enterprise |
| Gitea | 轻量级自托管 Git 服务 |
| Forgejo | Gitea 社区 fork |
| GitLab | 支持 GitLab.com 和自托管版本 |
| Bitbucket Cloud | Atlassian Bitbucket 云服务 |
| Bitbucket Data Center | Bitbucket 企业自托管版 |

---

## 执行后端

| 后端 | 说明 |
|------|------|
| **Docker**（默认） | 每个步骤在独立容器中运行，需要 Docker daemon |
| **Kubernetes** | 在 Kubernetes 集群中以 Pod 形式运行步骤 |
| **Local** | 直接在 agent 宿主机上运行命令（无容器隔离） |
| **自定义后端** | 实现 `pipeline/backend/backend.go` 中的 `Backend` 接口扩展 |

---

## 安装部署

### 数据库

Woodpecker 默认使用 SQLite（无需额外配置），生产环境建议使用 PostgreSQL：

```bash
WOODPECKER_DATABASE_DRIVER=postgres
WOODPECKER_DATABASE_DATASOURCE=postgres://user:password@localhost:5432/woodpecker?sslmode=disable
```

支持：`sqlite3`、`postgres`、`mysql`。

### 重要环境变量

```bash
# 服务器公开地址（必填）
WOODPECKER_HOST=https://ci.example.com

# agent 与 server 共享密钥（必填）
WOODPECKER_AGENT_SECRET=<用 openssl rand -hex 32 生成>

# 是否开放注册
WOODPECKER_OPEN=true

# Matrix 构建上限（本分支新增）
WOODPECKER_MAX_MATRIX_AXIS=25
WOODPECKER_MAX_MATRIX_TAGS=10
```

### Helm Chart 部署

```bash
helm repo add woodpecker https://woodpecker-ci.org/helm
helm install woodpecker woodpecker/woodpecker \
  --set server.env.WOODPECKER_HOST=https://ci.example.com \
  --set server.env.WOODPECKER_GITHUB=true
```

---

## 构建与开发

### 环境要求

- Go 1.25+（server 需要 `CGO_ENABLED=1` 支持 SQLite；agent/CLI 使用 `CGO_ENABLED=0`）
- pnpm（前端依赖管理）
- Node.js 20+
- Make

### 构建

```bash
# 构建全部二进制（会先构建前端）
make build

# 单独构建
make build-server   # -> dist/woodpecker-server
make build-agent    # -> dist/woodpecker-agent
make build-cli      # -> dist/woodpecker-cli

# 仅构建前端
make build-ui

# 前端热重载开发
cd web && pnpm install --frozen-lockfile && pnpm dev
```

### 代码生成

```bash
# 生成 mock、OpenAPI spec
make generate
```

---

## 项目结构

```
woodpecker/
├── cmd/
│   ├── server/          # server 入口，flags、setup、启动
│   ├── agent/           # agent 入口
│   └── cli/             # CLI 入口
├── server/
│   ├── api/             # REST API handlers（Gin）
│   ├── forge/           # VCS 集成（github/gitea/forgejo/gitlab/bitbucket）
│   ├── model/           # 数据模型（Repo/Pipeline/Workflow/Step/User/Secret/Artifact/Approval）
│   ├── pipeline/        # 流水线生命周期（编译/调度/审批/取消）
│   ├── router/          # HTTP 路由注册
│   ├── rpc/             # gRPC 服务端（agent ↔ server）
│   └── store/           # 数据库抽象层（datastore/ 为 XORM 实现）
├── agent/
│   └── runner/          # agent 端流水线执行循环
├── pipeline/
│   ├── backend/         # 执行后端（docker/kubernetes/local）
│   ├── frontend/yaml/   # .woodpecker.yml 解析与编译
│   └── rpc/proto/       # gRPC proto 定义
├── web/src/             # Vue 3 + TypeScript 前端
└── docs/                # 官方文档（Docusaurus）
```

---

## 许可证

Apache License 2.0，详见 [LICENSE](LICENSE)。

---

## 相关链接

- 上游项目：https://github.com/woodpecker-ci/woodpecker
- 官方文档：https://woodpecker-ci.org
- 插件中心：https://woodpecker-ci.org/plugins
- 本分支仓库：https://github.com/yunper-wang/woodpecker
