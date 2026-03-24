# 环境变量

Woodpecker 支持向 pipeline 的各个步骤传递自定义环境变量。注意，自定义变量不能覆盖已有的内置变量。以下是包含自定义环境变量的 pipeline 步骤示例：

```diff
 steps:
   - name: build
     image: golang
+    environment:
+      CGO: 0
+      GOOS: linux
+      GOARCH: amd64
     commands:
       - go build
       - go test
```

请注意，`environment` 部分无法展开环境变量。如果需要展开变量，应在 `commands` 部分使用 export 导出。

```diff
 steps:
   - name: build
     image: golang
-    environment:
-      - PATH=$PATH:/go
     commands:
+      - export PATH=$PATH:/go
       - go build
       - go test
```

:::warning
`${variable}` 表达式会被预处理。如果不希望预处理器对你的表达式求值，必须对其进行转义：
:::

```diff
 steps:
   - name: build
     image: golang
     commands:
-      - export PATH=${PATH}:/go
+      - export PATH=$${PATH}:/go
       - go build
       - go test
```

## 内置环境变量

以下是所有可在 pipeline 容器中使用的环境变量参考列表。这些变量在运行时注入到 pipeline 步骤和插件容器中。

| 名称                               | 描述                                                                                                                   | 示例                                                                                                       |
| ---------------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| `CI`                               | CI 环境名称                                                                                                             | `woodpecker`                                                                                               |
|                                    | **仓库**                                                                                                               |                                                                                                            |
| `CI_REPO`                          | 仓库全名 `<owner>/<name>`                                                                                               | `john-doe/my-repo`                                                                                         |
| `CI_REPO_OWNER`                    | 仓库所有者                                                                                                              | `john-doe`                                                                                                 |
| `CI_REPO_NAME`                     | 仓库名称                                                                                                                | `my-repo`                                                                                                  |
| `CI_REPO_REMOTE_ID`                | 仓库在 forge 中的远程 ID（UID）                                                                                          | `82`                                                                                                       |
| `CI_REPO_URL`                      | 仓库 Web URL                                                                                                            | `https://git.example.com/john-doe/my-repo`                                                                 |
| `CI_REPO_CLONE_URL`                | 仓库 clone URL                                                                                                          | `https://git.example.com/john-doe/my-repo.git`                                                             |
| `CI_REPO_CLONE_SSH_URL`            | 仓库 SSH clone URL                                                                                                      | `git@git.example.com:john-doe/my-repo.git`                                                                 |
| `CI_REPO_DEFAULT_BRANCH`           | 仓库默认分支                                                                                                             | `main`                                                                                                     |
| `CI_REPO_PRIVATE`                  | 仓库是否为私有                                                                                                           | `true`                                                                                                     |
| `CI_REPO_TRUSTED_NETWORK`          | 仓库是否拥有可信网络访问权限                                                                                              | `false`                                                                                                    |
| `CI_REPO_TRUSTED_VOLUMES`          | 仓库是否拥有可信 volumes 访问权限                                                                                         | `false`                                                                                                    |
| `CI_REPO_TRUSTED_SECURITY`         | 仓库是否拥有可信安全访问权限                                                                                              | `false`                                                                                                    |
|                                    | **当前 Commit**                                                                                                         |                                                                                                            |
| `CI_COMMIT_SHA`                    | commit SHA                                                                                                              | `eba09b46064473a1d345da7abf28b477468e8dbd`                                                                 |
| `CI_COMMIT_REF`                    | commit ref                                                                                                              | `refs/heads/main`                                                                                          |
| `CI_COMMIT_REFSPEC`                | commit ref spec                                                                                                         | `issue-branch:main`                                                                                        |
| `CI_COMMIT_BRANCH`                 | commit 分支（pull request 时等同于目标分支）                                                                              | `main`                                                                                                     |
| `CI_COMMIT_SOURCE_BRANCH`          | commit 源分支（仅在 pull request 事件时设置）                                                                             | `issue-branch`                                                                                             |
| `CI_COMMIT_TARGET_BRANCH`          | commit 目标分支（仅在 pull request 事件时设置）                                                                           | `main`                                                                                                     |
| `CI_COMMIT_TAG`                    | commit tag 名称（非 `tag` 事件时为空）                                                                                    | `v1.10.3`                                                                                                  |
| `CI_COMMIT_PULL_REQUEST`           | commit pull request 编号（仅在 pull request 事件时设置）                                                                  | `1`                                                                                                        |
| `CI_COMMIT_PULL_REQUEST_LABELS`    | pull request 上的标签（仅在 pull request 事件时设置）                                                                     | `server`                                                                                                   |
| `CI_COMMIT_PULL_REQUEST_MILESTONE` | pull request 所属里程碑（仅在 `pull_request` 和 `pull_request_closed` 事件时设置）                                        | `summer-sprint`                                                                                            |
| `CI_COMMIT_MESSAGE`                | commit 消息                                                                                                             | `Initial commit`                                                                                           |
| `CI_COMMIT_AUTHOR`                 | commit 作者用户名                                                                                                        | `john-doe`                                                                                                 |
| `CI_COMMIT_AUTHOR_EMAIL`           | commit 作者邮箱地址                                                                                                      | `john-doe@example.com`                                                                                     |
| `CI_COMMIT_PRERELEASE`             | 是否为预发布版本（非 `release` 事件时为空）                                                                               | `false`                                                                                                    |
|                                    | **当前 pipeline**                                                                                                       |                                                                                                            |
| `CI_PIPELINE_NUMBER`               | pipeline 编号                                                                                                           | `8`                                                                                                        |
| `CI_PIPELINE_PARENT`               | 父 pipeline 编号                                                                                                        | `0`                                                                                                        |
| `CI_PIPELINE_EVENT`                | pipeline 事件（参见 [`event`](../20-usage/20-workflow-syntax.md#event)）                                                  | `push`, `pull_request`, `pull_request_closed`, `pull_request_metadata`, `tag`, `release`, `manual`, `cron` |
| `CI_PIPELINE_EVENT_REASON`         | 触发 `pull_request_metadata` 事件的具体原因，与 forge 实例相关，可能发生变化                                               | `label_updated`, `milestoned`, `demilestoned`, `assigned`, `edited`, ...                                   |
| `CI_PIPELINE_URL`                  | pipeline 在 Web UI 中的链接                                                                                             | `https://ci.example.com/repos/7/pipeline/8`                                                                |
| `CI_PIPELINE_FORGE_URL`            | forge Web UI 中触发该 pipeline 的 commit 或 tag 的链接                                                                   | `https://git.example.com/john-doe/my-repo/commit/eba09b46064473a1d345da7abf28b477468e8dbd`                 |
| `CI_PIPELINE_DEPLOY_TARGET`        | `deployment` 事件的 pipeline 部署目标                                                                                    | `production`                                                                                               |
| `CI_PIPELINE_DEPLOY_TASK`          | `deployment` 事件的 pipeline 部署任务                                                                                    | `migration`                                                                                                |
| `CI_PIPELINE_CREATED`              | pipeline 创建时间（UNIX 时间戳）                                                                                          | `1722617519`                                                                                               |
| `CI_PIPELINE_STARTED`              | pipeline 启动时间（UNIX 时间戳）                                                                                          | `1722617519`                                                                                               |
| `CI_PIPELINE_FILES`                | 变更文件列表（非 `push` 或 `pull_request` 事件时为空；超过 500 个文件时为未定义）                                           | `[]`, `[".woodpecker.yml","README.md"]`                                                                   |
| `CI_PIPELINE_AUTHOR`               | pipeline 触发者用户名                                                                                                    | `octocat`                                                                                                  |
| `CI_PIPELINE_AVATAR`               | pipeline 触发者头像                                                                                                      | `https://git.example.com/avatars/5dcbcadbce6f87f8abef`                                                     |
|                                    | **当前 workflow**                                                                                                       |                                                                                                            |
| `CI_WORKFLOW_NAME`                 | workflow 名称                                                                                                           | `release`                                                                                                  |
|                                    | **当前步骤**                                                                                                            |                                                                                                            |
| `CI_STEP_NAME`                     | 步骤名称                                                                                                                | `build package`                                                                                            |
| `CI_STEP_NUMBER`                   | 步骤编号                                                                                                                | `0`                                                                                                        |
| `CI_STEP_STARTED`                  | 步骤启动时间（UNIX 时间戳）                                                                                               | `1722617519`                                                                                               |
| `CI_STEP_URL`                      | 步骤在 UI 中的 URL                                                                                                      | `https://ci.example.com/repos/7/pipeline/8`                                                                |
|                                    | **上一个 Commit**                                                                                                       |                                                                                                            |
| `CI_PREV_COMMIT_SHA`               | 上一个 commit SHA                                                                                                       | `15784117e4e103f36cba75a9e29da48046eb82c4`                                                                 |
| `CI_PREV_COMMIT_REF`               | 上一个 commit ref                                                                                                       | `refs/heads/main`                                                                                          |
| `CI_PREV_COMMIT_REFSPEC`           | 上一个 commit ref spec                                                                                                  | `issue-branch:main`                                                                                        |
| `CI_PREV_COMMIT_BRANCH`            | 上一个 commit 分支                                                                                                       | `main`                                                                                                     |
| `CI_PREV_COMMIT_SOURCE_BRANCH`     | 上一个 commit 源分支（仅在 pull request 事件时设置）                                                                      | `issue-branch`                                                                                             |
| `CI_PREV_COMMIT_TARGET_BRANCH`     | 上一个 commit 目标分支（仅在 pull request 事件时设置）                                                                    | `main`                                                                                                     |
| `CI_PREV_COMMIT_URL`               | 上一个 commit 在 forge 中的链接                                                                                          | `https://git.example.com/john-doe/my-repo/commit/15784117e4e103f36cba75a9e29da48046eb82c4`                 |
| `CI_PREV_COMMIT_MESSAGE`           | 上一个 commit 消息                                                                                                       | `test`                                                                                                     |
| `CI_PREV_COMMIT_AUTHOR`            | 上一个 commit 作者用户名                                                                                                  | `john-doe`                                                                                                 |
| `CI_PREV_COMMIT_AUTHOR_EMAIL`      | 上一个 commit 作者邮箱地址                                                                                                | `john-doe@example.com`                                                                                     |
|                                    | **上一个 pipeline**                                                                                                     |                                                                                                            |
| `CI_PREV_PIPELINE_NUMBER`          | 上一个 pipeline 编号                                                                                                    | `7`                                                                                                        |
| `CI_PREV_PIPELINE_PARENT`          | 上一个 pipeline 的父 pipeline 编号                                                                                       | `0`                                                                                                        |
| `CI_PREV_PIPELINE_EVENT`           | 上一个 pipeline 事件（参见 [`event`](../20-usage/20-workflow-syntax.md#event)）                                           | `push`, `pull_request`, `pull_request_closed`, `pull_request_metadata`, `tag`, `release`, `manual`, `cron` |
| `CI_PREV_PIPELINE_EVENT_REASON`    | 上一个 pipeline 触发 `pull_request_metadata` 事件的具体原因，与 forge 实例相关，可能发生变化                               | `label_updated`, `milestoned`, `demilestoned`, `assigned`, `edited`, ...                                   |
| `CI_PREV_PIPELINE_URL`             | 上一个 pipeline 在 CI 中的链接                                                                                           | `https://ci.example.com/repos/7/pipeline/7`                                                                |
| `CI_PREV_PIPELINE_FORGE_URL`       | 上一个 pipeline 在 forge 中的事件链接                                                                                    | `https://git.example.com/john-doe/my-repo/commit/15784117e4e103f36cba75a9e29da48046eb82c4`                 |
| `CI_PREV_PIPELINE_DEPLOY_TARGET`   | 上一个 pipeline 的 `deployment` 事件部署目标                                                                             | `production`                                                                                               |
| `CI_PREV_PIPELINE_DEPLOY_TASK`     | 上一个 pipeline 的 `deployment` 事件部署任务                                                                             | `migration`                                                                                                |
| `CI_PREV_PIPELINE_STATUS`          | 上一个 pipeline 状态                                                                                                    | `success`, `failure`                                                                                       |
| `CI_PREV_PIPELINE_CREATED`         | 上一个 pipeline 创建时间（UNIX 时间戳）                                                                                   | `1722610173`                                                                                               |
| `CI_PREV_PIPELINE_STARTED`         | 上一个 pipeline 启动时间（UNIX 时间戳）                                                                                   | `1722610173`                                                                                               |
| `CI_PREV_PIPELINE_FINISHED`        | 上一个 pipeline 完成时间（UNIX 时间戳）                                                                                   | `1722610383`                                                                                               |
| `CI_PREV_PIPELINE_AUTHOR`          | 上一个 pipeline 触发者用户名                                                                                             | `octocat`                                                                                                  |
| `CI_PREV_PIPELINE_AVATAR`          | 上一个 pipeline 触发者头像                                                                                               | `https://git.example.com/avatars/5dcbcadbce6f87f8abef`                                                     |
|                                    | &emsp;                                                                                                                  |                                                                                                            |
| `CI_WORKSPACE`                     | 源代码被 clone 到的工作区路径                                                                                            | `/woodpecker/src/git.example.com/john-doe/my-repo`                                                         |
|                                    | **系统**                                                                                                                |                                                                                                            |
| `CI_SYSTEM_NAME`                   | CI 系统名称                                                                                                             | `woodpecker`                                                                                               |
| `CI_SYSTEM_URL`                    | CI 系统链接                                                                                                             | `https://ci.example.com`                                                                                   |
| `CI_SYSTEM_HOST`                   | CI server 主机名                                                                                                        | `ci.example.com`                                                                                           |
| `CI_SYSTEM_VERSION`                | server 版本                                                                                                             | `2.7.0`                                                                                                    |
|                                    | **Forge**                                                                                                               |                                                                                                            |
| `CI_FORGE_TYPE`                    | forge 名称                                                                                                              | `bitbucket` , `bitbucket_dc` , `forgejo` , `gitea` , `github` , `gitlab`                                   |
| `CI_FORGE_URL`                     | 已配置的 forge 根 URL                                                                                                   | `https://git.example.com`                                                                                  |
|                                    | **内部变量** - 请勿使用！                                                                                               |                                                                                                            |
| `CI_SCRIPT`                        | 内部脚本路径，用于调用 pipeline 步骤命令。                                                                               |                                                                                                            |
| `CI_NETRC_USERNAME`                | 私有仓库 clone 所需的凭据。（仅对特定镜像可用）                                                                           |                                                                                                            |
| `CI_NETRC_PASSWORD`                | 私有仓库 clone 所需的凭据。（仅对特定镜像可用）                                                                           |                                                                                                            |
| `CI_NETRC_MACHINE`                 | 私有仓库 clone 所需的凭据。（仅对特定镜像可用）                                                                           |                                                                                                            |

## 全局环境变量

如果希望特定环境变量在所有 pipeline 中均可用，可以在 Woodpecker server 上使用 `WOODPECKER_ENVIRONMENT` 配置项。注意，这些变量不能覆盖已有的内置变量。

```ini
WOODPECKER_ENVIRONMENT=first_var:value1,second_var:value2
```

例如，可以用来统一管理多个项目所使用的镜像 tag：

```ini
WOODPECKER_ENVIRONMENT=GOLANG_VERSION:1.18
```

```diff
 steps:
   - name: build
-    image: golang:1.18
+    image: golang:${GOLANG_VERSION}
     commands:
       - [...]
```

## 字符串替换

Woodpecker 支持在运行时对环境变量进行字符串替换，使我们能够在 pipeline 配置中使用动态的 settings、命令和过滤器。

commit 替换示例：

```diff
 steps:
   - name: s3
     image: woodpeckerci/plugin-s3
     settings:
+      target: /target/${CI_COMMIT_SHA}
```

tag 替换示例：

```diff
 steps:
   - name: s3
     image: woodpeckerci/plugin-s3
     settings:
+      target: /target/${CI_COMMIT_TAG}
```

## 字符串操作

Woodpecker 还模拟了 bash 的字符串操作，使我们能够在替换前对字符串进行处理。典型用途包括截取子串、去除前缀或后缀等。

| 操作                 | 描述                         |
| -------------------- | ---------------------------- |
| `${param}`           | 参数替换                     |
| `${param,}`          | 参数替换，首字母转小写        |
| `${param,,}`         | 参数替换，全部转小写          |
| `${param^}`          | 参数替换，首字母转大写        |
| `${param^^}`         | 参数替换，全部转大写          |
| `${param:pos}`       | 参数替换，截取子串            |
| `${param:pos:len}`   | 参数替换，截取指定长度子串    |
| `${param=default}`   | 参数替换，设置默认值          |
| `${param##prefix}`   | 参数替换，删除前缀            |
| `${param%%suffix}`   | 参数替换，删除后缀            |
| `${param/old/new}`   | 参数替换，查找并替换          |

截取子串的变量替换示例：

```diff
 steps:
   - name: s3
     image: woodpeckerci/plugin-s3
     settings:
+      target: /target/${CI_COMMIT_SHA:0:8}
```

删除 `v.1.0.0` 中 `v` 前缀的变量替换示例：

```diff
 steps:
   - name: s3
     image: woodpeckerci/plugin-s3
     settings:
+      target: /target/${CI_COMMIT_TAG##v}
```

## `pull_request_metadata` 事件的具体原因值

对于 `pull_request_metadata` 事件，检测到元数据变更的具体原因通过 `CI_PIPELINE_EVENT_REASON` 传递。

**GitLab** 将元数据更新合并到一个 webhook 中。事件原因以 `,` 分隔为列表。

:::note
事件原因值与 forge 相关，可能在版本之间发生变化。
:::

| 事件                   | GitHub             | Gitea              | Forgejo            | GitLab             | Bitbucket | Bitbucket Datacenter | 描述                                                                           |
| -------------------- | ------------------ | ------------------ | ------------------ | ------------------ | --------- | -------------------- | ------------------------------------------------------------------------------ |
| `assigned`           | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:       | :x:                  | Pull request 被分配给某用户                                                     |
| `converted_to_draft` | :white_check_mark: | :x:                | :x:                | :x:                | :x:       | :x:                  | Pull request 被转为草稿                                                         |
| `demilestoned`       | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:       | :x:                  | Pull request 从里程碑中移除                                                     |
| `description_edited` | :x:                | :x:                | :x:                | :white_check_mark: | :x:       | :x:                  | 描述被编辑                                                                      |
| `edited`             | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                | :x:       | :x:                  | Pull request 的标题、正文被编辑，或基础分支被更改                                |
| `label_added`        | :x:                | :x:                | :x:                | :white_check_mark: | :x:       | :x:                  | Pull request 原本无标签，现在添加了标签                                          |
| `label_cleared`      | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:       | :x:                  | 所有标签被移除                                                                  |
| `label_updated`      | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:       | :x:                  | 新增标签或标签发生变更                                                           |
| `locked`             | :white_check_mark: | :x:                | :x:                | :x:                | :x:       | :x:                  | Pull request 的对话被锁定                                                       |
| `milestoned`         | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:       | :x:                  | Pull request 被添加到里程碑                                                     |
| `ready_for_review`   | :white_check_mark: | :x:                | :x:                | :x:                | :x:       | :x:                  | 草稿 pull request 被标记为可供审阅                                               |
| `review_requested`   | :x:                | :x:                | :x:                | :white_check_mark: | :x:       | :x:                  | 请求了新的审阅                                                                  |
| `title_edited`       | :x:                | :x:                | :x:                | :white_check_mark: | :x:       | :x:                  | 标题被编辑                                                                      |
| `unassigned`         | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:       | :x:                  | 用户被从 pull request 中取消分配                                                 |
| `unlabeled`          | :white_check_mark: | :x:                | :x:                | :x:                | :x:       | :x:                  | 标签从 pull request 中被移除                                                    |
| `unlocked`           | :white_check_mark: | :x:                | :x:                | :x:                | :x:       | :x:                  | Pull request 的对话被解锁                                                       |

**Bitbucket** 和 **Bitbucket Datacenter** [目前暂不支持](https://github.com/woodpecker-ci/woodpecker/pull/5214)。