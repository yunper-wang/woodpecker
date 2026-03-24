# 欢迎使用 Woodpecker CI

Woodpecker 是一个社区驱动的持续集成（CI）系统，专为自托管而设计。

它是 [Drone CI](https://github.com/drone/drone) 的一个分支，专注于简洁性与可靠性。

![Woodpecker CI](./woodpecker.png)

## 工作原理

Woodpecker 由两个主要组件构成：

- **Server**：负责管理整个系统，包括与代码托管平台（forge）通信、存储数据、调度任务等。
- **Agent**：负责实际执行 pipeline，接收 server 分发的任务并在容器中运行。

当代码发生变更时，forge 通过 webhook 通知 server，server 随即将构建任务调度给可用的 agent 执行。

## 支持的代码托管平台（Forge）

Woodpecker 支持以下代码托管平台：

- GitHub
- Gitea / Forgejo
- GitLab
- Bitbucket
- Bitbucket Data Center

## 下一步

- [快速入门](../usage/intro)
- 了解 [Pipeline 语法](../usage/workflow-syntax)
