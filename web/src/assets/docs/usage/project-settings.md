# 项目设置

作为 Woodpecker 项目的所有者，你可以通过 Web 界面修改项目相关设置。

![project settings](./project-settings.png)

## Pipeline 路径

pipeline 配置文件或文件夹的路径。默认留空时，将按以下顺序解析配置文件：`.woodpecker/*.{yaml,yml}` -> `.woodpecker.yaml` -> `.woodpecker.yml`。如果设置了自定义路径，Woodpecker 会尝试从指定位置加载配置，若找不到则报错。要使用[多 workflow](./25-workflows.md) 并自定义路径，需将路径改为以 `/` 结尾的文件夹路径，例如 `.woodpecker/`。

## 仓库 Hooks

你的版本控制系统会通过 webhook 向 Woodpecker 通知事件。如果你希望 pipeline 只在特定 webhook 上运行，可以在此处勾选对应的事件类型。

## 允许 Pull Requests

启用后，Woodpecker 将处理 webhook 的 pull request 事件。如果禁用，pipeline 不会在 pull request 时运行。

## 允许部署

启用后，可以从成功的 pipeline 通过 `deploy` 事件启动新的 pipeline。

:::danger
仅当你信任所有对仓库拥有推送权限的用户时，才应启用此选项。
否则，这些用户将能够窃取仅对 `deploy` 事件可用的 secret。
:::

## 需要审批

为防止恶意 pipeline 提取 secret 或执行有害命令，或防止意外触发 pipeline，你可以要求进行额外的审批流程。根据启用的选项，pipeline 创建后将被挂起，仅在审批通过后才继续执行。默认的限制性设置为 `对 fork 仓库要求审批`。

## 受信任（Trusted）

如果将项目设置为受信任，pipeline 步骤及其底层容器将获得提升的能力，例如挂载 volume。

:::note

只有 server 管理员才能设置此选项。如果你不是 server 管理员，项目设置中将不会显示此选项。

:::

## 自定义受信任 clone 插件

在 clone 过程中，可能需要 Git 凭据（例如用于私有仓库）。
这些凭据通过 [`netrc`](https://everything.curl.dev/usingcurl/netrc.html) 提供。

这些凭据仅注入到环境变量 `WOODPECKER_PLUGINS_TRUSTED_CLONE`（实例级 Woodpecker server 设置）中指定的受信任插件，或在此仓库级设置中声明的插件中。

使用这些凭据可以执行任何 Git