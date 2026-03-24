# Configuration Extension

Configuration extension 可用于修改或生成 Woodpecker 的 pipeline 配置。你可以在仓库设置的 extensions 标签页中配置一个 HTTP 端点。

使用此类 extension 在以下场景中非常有用：

<!-- cSpell:words templating,Starlark,Jsonnet -->

- 使用 Go templating 等方式对原始配置文件进行预处理
- 将自定义属性转换为 Woodpecker 属性
- 为配置添加默认值，例如默认 step
- 将完全不同格式的配置文件（如 Gitlab CI 配置、Starlark、Jsonnet 等）进行转换
- 将多个仓库的配置集中在一处管理

## 安全性

:::warning
由于 Woodpecker 会传递 token 等私密信息，并会执行返回的配置，因此对外部 extension 进行安全防护极为重要。Woodpecker 会对每个请求进行签名。详情请参阅[安全章节](./extensions.md#安全性)。
:::

## 全局配置

除了可以按仓库配置 extension 之外，你还可以在 Woodpecker server 配置中设置全局端点。如果你希望对所有仓库使用该 extension，这会非常有用。如果你与他人共享 Woodpecker server，请谨慎操作，因为他们也会使用你的 configuration extension。

如果同时配置了全局配置和仓库级配置，且仓库未启用独占设置，则全局配置会在仓库特定的 configuration extension 之前被调用。

```ini title="Server"
WOODPECKER_CONFIG_SERVICE_ENDPOINT=https://example.com/ciconfig
```

## 工作原理

当 pipeline 被触发时，Woodpecker 会从仓库获取 pipeline 配置，然后向已配置的 extension 发送一个 HTTP POST 请求，请求的 JSON payload 包含仓库、pipeline 信息以及从仓库获取的当前配置文件等数据。extension 随后可以返回修改后的或全新的 pipeline 配置（遵循 Woodpecker 官方 YAML 格式）。

你可以启用独占设置（全局或仓库级别均可）。此时 Woodpecker 将只调用你的 extension，不做其他操作。这允许你完全跳过 forge。发送给 extension 的请求将不包含配置文件。

### Request

extension 接收到的 HTTP POST 请求包含以下 JSON payload：

```ts
class Request {
  repo: Repo;
  pipeline: Pipeline;
  netrc: Netrc;
  configuration: {
    // list of configurations. Not send if there was none.
    name: string; // filename of the configuration file
    data: string; // content of the configuration file
  }[];
}
```

查看以下模型以获取更多信息：

- [repo model](https://github.com/woodpecker-ci/woodpecker/blob/main/server/model/repo.go)
- [pipeline model](https://github.com/woodpecker-ci/woodpecker/blob/main/server/model/pipeline.go)
- [netrc model](https://github.com/woodpecker-ci/woodpecker/blob/main/server/model/netrc.go)

:::tip
`netrc` 数据非常强大，因为它包含访问仓库的凭证。你可以用它来克隆仓库，甚至使用 forge（Github 或 Gitlab 等）API 获取仓库的更多信息。
:::

请求示例：

```json
{
  "repo": {
    "id": 100,
    "uid": "",
    "user_id": 0,
    "namespace": "",
    "name": "woodpecker-test-pipeline",
    "slug": "",
    "scm": "git",
    "git_http_url": "",
    "git_ssh_url": "",
    "link": "",
    "default_branch": "",
    "private": true,
    "visibility": "private",
    "active": true,
    "config": "",
    "trusted": false,
    "protected": false,
    "ignore_forks": false,
    "ignore_pulls": false,
    "cancel_pulls": false,
    "timeout": 60,
    "counter": 0,
    "synced": 0,
    "created": 0,
    "updated": 0,
    "version": 0
  },
  "pipeline": {
    "author": "myUser",
    "author_avatar": "https://myforge.com/avatars/d6b3f7787a685fcdf2a44e2c685c7e03",
    "author_email": "my@email.com",
    "branch": "main",
    "changed_files": ["some-filename.txt"],
    "commit": "2fff90f8d288a4640e90f05049fe30e61a14fd50",
    "created_at": 0,
    "deploy_to": "",
    "enqueued_at": 0,
    "error": "",
    "event": "push",
    "finished_at": 0,
    "id": 0,
    "link_url": "https://myforge.com/myUser/woodpecker-testpipe/commit/2fff90f8d288a4640e90f05049fe30e61a14fd50",
    "message": "test old config\n",
    "number": 0,
    "parent": 0,
    "ref": "refs/heads/main",
    "refspec": "",
    "clone_url": "",
    "reviewed_at": 0,
    "reviewed_by": "",
    "sender": "myUser",
    "signed": false,
    "started_at": 0,
    "status": "",
    "timestamp": 1645962783,
    "title": "",
    "updated_at": 0,
    "verified": false
  },
  "configuration": [
    {
      "name": ".woodpecker.yaml",
      "data": "steps:\n  - name: backend\n    image: alpine\n    commands:\n      - echo \"Hello there from Repo (.woodpecker.yaml)\"\n"
    }
  ]
}
```

### Response

extension 应返回一个 JSON payload，包含遵循 Woodpecker 官方 YAML 格式的新配置文件。
如果 extension 希望保留现有配置文件，可以返回 HTTP 状态码 `204 No Content`。

```ts
class Response {
  configs: {
    name: string; // filename of the configuration file
    data: string; // content of the configuration file
  }[];
}
```

响应示例：

```json
{
  "configs": [
    {
      "name": "central-override",
      "data": "steps:\n  - name: backend\n    image: alpine\n    commands:\n      - echo \"Hello there from ConfigAPI\"\n"
    }
  ]
}
```
