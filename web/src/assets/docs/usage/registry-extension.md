# Registry Extension

Woodpecker 使用 registry extension 来获取 registry 凭证。你可以在仓库设置的 extensions 标签页中配置一个 HTTP 端点。

使用此类 extension 在以下场景中非常有用：

- 集中管理 registry 凭证
- 使用外部存储来存放凭证
- 动态管理 Woodpecker 应使用的凭证

## 安全性

:::warning
由于 Woodpecker 会传递 token 等私密信息，并会执行返回的配置，因此对外部 extension 进行安全防护极为重要。Woodpecker 会对每个请求进行签名。详情请参阅[安全章节](./extensions.md#安全性)。
:::

## 全局配置

除了可以按仓库配置 extension 之外，你还可以在 Woodpecker server 配置中设置全局端点。如果你希望对所有仓库使用该 extension，这会非常有用。如果你与他人共享 Woodpecker server，请谨慎操作，因为他们也会使用你的 registry extension。

如果全局 extension 和仓库级 extension 都返回了同一 registry 的凭证，将使用仓库 extension 返回的凭证。

```ini title="Server"
WOODPECKER_REGISTRY_SERVICE_ENDPOINT=https://example.com/ciconfig
```

## 工作原理

当 pipeline 被触发时，Woodpecker 会从你的服务获取凭证。如果获取失败，则使用直接在 Woodpecker 中配置的凭证作为备用。

### Request

extension 接收到的 HTTP POST 请求包含以下 JSON payload：

```ts
class Request {
  repo: Repo;
  pipeline: Pipeline;
}
```

查看以下模型以获取更多信息：

- [repo model](https://github.com/woodpecker-ci/woodpecker/blob/main/server/model/repo.go)
- [pipeline model](https://github.com/woodpecker-ci/woodpecker/blob/main/server/model/pipeline.go)

请求示例：

```json
// Please check the latest structure in the models mentioned above.
// This example is likely outdated.

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
  }
}
```

### Response

extension 应返回一个包含 registry 凭证的 JSON payload。
如果 extension 希望保留现有配置，可以返回 HTTP 状态码 `204 No Content`。

```ts
class Response {
  registries: {
    address: string; // the docker registry address
    username: string; // registry username
    password: string; // registry password
  }[];
}
```

响应示例：

```json
{
  "registries": [
    {
      "address": "docker.io",
      "username": "woodpecker-bot",
      "password": "your-pass-word-123"
    }
  ]
}
```
