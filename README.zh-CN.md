# Woodpecker CI

<p align="center">
  <a href="https://github.com/woodpecker-ci/woodpecker/">
    <img alt="Woodpecker" src="docs/static/img/logo.svg" width="220"/>
  </a>
</p>

Woodpecker 是一个基于 [Drone](https://github.com/drone/drone) 构建的社区驱动的持续集成（CI）系统。

> 「这只小鸟知道你的秘密」

## 概述

Woodpecker CI 是一个轻量、可扩展的 CI/CD 平台，支持通过 YAML 文件定义流水线，并可运行在 Docker、Kubernetes 等多种后端之上。

## 架构

Woodpecker 由三个独立的二进制组成：

| 组件 | 路径 | 说明 |
|------|------|------|
| **server** | `cmd/server` | HTTP API、Webhook 接收、流水线调度、gRPC 服务端 |
| **agent** | `cmd/agent` | 通过 gRPC 向 server 轮询任务，执行流水线步骤 |
| **cli** | `cmd/cli` | 面向用户的命令行客户端 |

## 主要特性

- 支持 GitHub、Gitea、Forgejo、GitLab、Bitbucket 等多种代码托管平台
- 基于 YAML 定义流水线，语法简洁
- 支持 Matrix 矩阵构建
- 支持 Docker、Kubernetes、本地等多种执行后端
- 插件生态丰富
- 支持 Secret 和 Registry 管理
- 内置流水线审批门控

## 本分支新增功能

本分支在上游基础上新增了以下功能：

### 1. Matrix `exclude` 支持

在 matrix 定义中可以使用 `exclude` 排除特定组合，避免不需要的构建：

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

同时，超出上限时会返回明确错误（而非静默截断）。

### 2. Matrix 上限可配置

通过环境变量或启动参数自定义 matrix 最大规模：

```bash
WOODPECKER_MAX_MATRIX_AXIS=50   # 最大排列数，默认 25
WOODPECKER_MAX_MATRIX_TAGS=15   # 最大变量维度数，默认 10
```

### 3. `when.path` 支持 tag 和 release 事件

路径过滤条件现在也对 `tag` 和 `release` 事件生效：

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
```

GitHub / Gitea / Forgejo 三个 forge 均已同步支持，会在 tag/release 触发时自动获取变更文件列表。

### 4. 流水线 Artifact 存储

工作流可以上传 artifact，供后续下载或共享：

```bash
# 上传 artifact
POST /api/repos/{repo_id}/pipelines/{number}/artifacts?name=build.tar.gz&workflow_id=1

# 列出 artifact
GET /api/repos/{repo_id}/pipelines/{number}/artifacts

# 下载 artifact
GET /api/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id}

# 删除 artifact
DELETE /api/repos/{repo_id}/pipelines/{number}/artifacts/{artifact_id}
```

### 5. 多级审批门控

支持配置需要多人批准才能放行被阻塞的流水线：

在项目设置中配置 `required_approvals`（默认为 1）。当流水线处于 `blocked` 状态时，需要达到指定数量的不同用户批准后方可继续执行。

```bash
# 查看审批记录
GET /api/repos/{repo_id}/pipelines/{number}/approvals

# 审批通过
POST /api/repos/{repo_id}/pipelines/{number}/approve
```

### 6. Bug 修复：`when.status` 逻辑修复

修复了 `IncludesStatusSuccess` 函数中遗留的调试日志输出（`fmt.Println`）及重复调用 `c.Match()` 的问题。

## 构建

```bash
# 构建所有二进制
make build

# 单独构建
make build-server
make build-agent
make build-cli

# 构建前端
make build-ui
```

## 测试

```bash
# 运行所有 Go 测试
make test

# 运行指定包
go test ./pipeline/frontend/yaml/...
go test ./server/pipeline/...
```

## 许可证

Apache License 2.0，详见 [LICENSE](LICENSE)。

## 相关链接

- 上游项目：https://github.com/woodpecker-ci/woodpecker
- 官方文档：https://woodpecker-ci.org
- 插件中心：https://woodpecker-ci.org/plugins
