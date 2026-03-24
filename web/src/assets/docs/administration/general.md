# 概述

Woodpecker 由核心组件（`server` 和 `agent`）以及一个可选组件（`autoscaler`）构成。

**server** 提供用户界面，处理发送到底层 forge 的 webhook 请求，提供 API 服务，并解析来自 YAML 文件的 pipeline 配置。

**agent** 通过特定的 [backend](../../20-usage/15-terminology/index.md)（Docker、Kubernetes、local）执行 [workflow](../../20-usage/15-terminology/index.md)，并通过 GRPC 连接到 server。可以同时运行多个 agent，从而针对单个实例精细调整任务限制、backend 选择以及其他 agent 相关设置。

**autoscaler** 允许在所选云服务商上启动新的虚拟机来处理待执行的构建任务。构建完成后，虚拟机会被销毁（在短暂的过渡时间之后）。

:::tip
你可以添加更多 agent 来增加并行 workflow 的数量，或者为 agent 设置 [`WOODPECKER_MAX_WORKFLOWS=1`](./agent.md#max_workflows) 环境变量来增加每个 agent 的并行 workflow 数量。
:::

## 数据库

Woodpecker 默认使用 SQLite 数据库，无需安装或配置。对于较大规模的实例，建议使用 Postgres 或 MariaDB 实例。详情请参阅[数据库设置](./server.md#databases)页面。

## Forge

没有代码的 CI/CD 系统无从谈起。通过将 Woodpecker 连接到你的 [forge](../../20-usage/15-terminology/index.md)，你可以在 push 或 pull request 等事件上启动 pipeline。Woodpecker 还会使用你的 forge 进行身份验证，并将 pipeline 的状态反馈回去。详情请参阅 [forge 设置](./10-configuration/12-forges/11-overview.md)页面。

## 容器镜像

:::info
不存在 `latest` 标签，以防止意外升级主版本。请使用 SemVer 标签或滚动的主/次版本标签。也可以使用 `next` 标签来获取从 `main` 分支构建的滚动版本。
:::

- `vX.Y.Z`：特定版本的 SemVer 标签，无 entrypoint shell（scratch 镜像）
  - `vX.Y`
  - `vX`
- `vX.Y.Z-alpine`：特定版本的 SemVer 标签，Server 和 CLI 以 rootless 方式运行（自 v3.0 起）。
  - `vX.Y-alpine`
  - `vX-alpine`
- `next`：从 `main` 分支构建
- `pull_<PR_ID>`：从 Pull Request 分支构建的镜像。

镜像推送至 DockerHub 和 Quay。

- woodpecker-server（[DockerHub](https://hub.docker.com/r/woodpeckerci/woodpecker-server) 或 [Quay](https://quay.io/repository/woodpeckerci/woodpecker-server)）
- woodpecker-agent（[DockerHub](https://hub.docker.com/r/woodpeckerci/woodpecker-agent) 或 [Quay](https://quay.io/repository/woodpeckerci/woodpecker-agent)）
- woodpecker-cli（[DockerHub](https://hub.docker.com/r/woodpeckerci/woodpecker-cli) 或 [Quay](https://quay.io/repository/woodpeckerci/woodpecker-cli)）
- woodpecker-autoscaler（[DockerHub](https://hub.docker.com/r/woodpeckerci/autoscaler)）
