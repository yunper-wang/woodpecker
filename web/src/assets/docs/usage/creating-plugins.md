# 创建 Plugin

创建一个新 plugin 非常简单：构建一个以你的 plugin 逻辑作为 ENTRYPOINT 的 Docker 容器即可。

## Settings

为了让用户能够配置 plugin 的行为，你应该使用 `settings:`。

这些 settings 会以带有 `PLUGIN_` 前缀的大写环境变量形式传递给你的 plugin。
使用像 `url` 这样的 setting 会生成名为 `PLUGIN_URL` 的环境变量。

像 `-` 这样的字符会被转换为下划线（`_`）。`some_String` 会变成 `PLUGIN_SOME_STRING`。
不遵循驼峰命名规则，`anInt` 会变成 `PLUGIN_ANINT`。 <!-- cspell:ignore ANINT -->

### 基础 settings

任何基础 YAML 类型（标量）都会被转换为字符串：

| Setting              | 环境变量值                    |
| -------------------- | ---------------------------- |
| `some-bool: false`   | `PLUGIN_SOME_BOOL="false"`   |
| `some_String: hello` | `PLUGIN_SOME_STRING="hello"` |
| `anInt: 3`           | `PLUGIN_ANINT="3"`           |

### 复杂 settings

也可以使用复杂类型的 settings，例如：

```yaml
steps:
  - name: plugin
    image: foo/plugin
    settings:
      complex:
        abc: 2
        list:
          - 2
          - 3
```

这类值会被转换为 JSON 后传递给你的 plugin。在上面的例子中，环境变量 `PLUGIN_COMPLEX` 的值为 `{"abc": "2", "list": [ "2", "3" ]}`。

### Secrets

Secrets 也应通过 settings 传递。因此，用户应使用 [`from_secret`](../40-secrets.md#usage)。

## Plugin 库

对于 Go，我们提供了一个 plugin 库，方便你访问内部环境变量和 settings。参见 <https://codeberg.org/woodpecker-plugins/go-plugin>。

## 元数据

在你的文档中，可以使用 Markdown 头部来定义 plugin 的元数据。这些数据会被 [我们的 plugin index](/plugins) 使用。

支持的元数据字段：

- `name`：plugin 的完整名称
- `icon`：你的 plugin 图标的 URL
- `description`：简短描述其功能
- `author`：你的名字
- `tags`：关键词列表（例如 clone plugin 使用 `[git, clone]`）
- `containerImage`：容器镜像的名称
- `containerImageUrl`：容器镜像的链接
- `url`：你的 plugin 的主页或仓库地址

如果你想让 plugin 出现在 index 中，应尽量填写所有字段，但只有 `name` 是必填的。

## 示例 Plugin

下面提供了一个简短的教程，演示如何使用简单的 shell 脚本创建 Woodpecker webhook plugin，
在构建 pipeline 中发起 HTTP 请求。

### 终端用户看到的效果

下面的例子演示了如何在 YAML 文件中配置一个 webhook plugin：

```yaml
steps:
  - name: webhook
    image: foo/webhook
    settings:
      url: https://example.com
      method: post
      body: |
        hello world
```

### 编写逻辑

创建一个简单的 shell 脚本，使用 YAML 配置参数调用 curl，这些参数以大写并带 `PLUGIN_` 前缀的环境变量形式传递给脚本。

```bash
#!/bin/sh

curl \
  -X ${PLUGIN_METHOD} \
  -d ${PLUGIN_BODY} \
  ${PLUGIN_URL}
```

### 打包

创建一个 Dockerfile，将你的 shell 脚本添加到镜像中，并配置镜像以你的 shell 脚本作为主要 entrypoint 执行。

```dockerfile
# please pin the version, e.g. alpine:3.19
FROM alpine
ADD script.sh /bin/
RUN chmod +x /bin/script.sh
RUN apk -Uuv add curl ca-certificates
ENTRYPOINT /bin/script.sh
```

构建并发布你的 plugin 到 Docker registry。发布后，你的 plugin 就可以与更广泛的 Woodpecker 社区共享了。

```shell
docker build -t foo/webhook .
docker push foo/webhook
```

在命令行本地执行你的 plugin 以验证其是否正常工作：

```shell
docker run --rm \
  -e PLUGIN_METHOD=post \
  -e PLUGIN_URL=https://example.com \
  -e PLUGIN_BODY="hello world" \
  foo/webhook
```

## 最佳实践

- 为不同架构构建你的 plugin，以便更多用户可以使用它。
  至少应支持 `amd64` 和 `arm64`。
- 为使用 `local` backend 的用户提供二进制文件。
  这些文件也应为不同的操作系统/架构构建。
- 尽可能使用[内置环境变量](../50-environment.md#built-in-environment-variables)。
- 除 settings（和内部环境变量）外，不要使用任何其他配置方式。这意味着：不要要求使用 [`environment`](../50-environment.md)，也不要要求特定的 secret 名称。
- 添加 `docs.md` 文件，列出所有 settings 和 plugin 元数据（[示例](https://github.com/woodpecker-ci/plugin-git/blob/main/docs.md)）。
- 使用你的 `docs.md` 将 plugin 添加到 [plugin index](/plugins)（[index 中的上述示例](https://woodpecker-ci.org/plugins/git-clone)）。
