# Registries

Woodpecker 提供了在仓库设置中添加容器 registry 的能力。添加 registry 后，在 pipeline 的步骤中使用私有镜像时，Woodpecker 可以对该容器 registry 进行身份验证并拉取私有镜像。使用 registry 凭据还可以帮助你避免从公共 registry 拉取镜像时遭遇限速。

## 来自私有 registry 的镜像

要拉取 YAML 配置文件中定义的私有容器镜像，你必须在 UI 中提供 registry 凭据。

这些凭据不会暴露给你的步骤，这意味着它们不能用于推送镜像，并且可以安全地用于 pull request 等场景。向 registry 推送镜像仍需为相应的插件单独配置凭据。

使用私有镜像的配置示例：

```diff
 steps:
   - name: build
+    image: gcr.io/custom/golang
     commands:
       - go build
       - go test
```

Woodpecker 会将 registry 的 hostname 与 YAML 中的每个镜像进行匹配。如果 hostname 匹配，则使用 registry 凭据向你的 registry 进行身份验证并拉取镜像。请注意，registry 凭据由 Woodpecker agent 使用，不会暴露给你的构建容器。

Registry hostname 示例：

- 镜像 `gcr.io/foo/bar` 的 hostname 为 `gcr.io`
- 镜像 `foo/bar` 的 hostname 为 `docker.io`
- 镜像 `qux.com:8000/foo/bar` 的 hostname 为 `qux.com:8000`

Registry hostname 匹配逻辑示例：

- Hostname `gcr.io` 匹配镜像 `gcr.io/foo/bar`
- Hostname `docker.io` 匹配 `golang`
- Hostname `docker.io` 匹配 `library/golang`
- Hostname `docker.io` 匹配 `bradrydzewski/golang`
- Hostname `docker.io` 匹配 `bradrydzewski/golang:latest`

## 全局 registry 支持

要使私有 registry 在全局范围内可用，请查阅 [server 配置文档](../30-administration/10-configuration/10-server.md#docker_config)。

## GCR registry 支持

有关配置 Google Container Registry 访问权限的具体详情，请查阅[此处](https://cloud.google.com/container-registry/docs/advanced-authentication#using_a_json_key_file)的文档。

## 本地镜像

:::warning
此功能需要特权权限，仅管理员可用。此外，此方式仅在使用单个 agent 时有效。
:::

可以通过将 docker socket 挂载为 volume 来构建本地镜像。

在项目根目录下有一个 `Dockerfile` 时：

```yaml
steps:
  - name: build-image
    image: docker
    commands:
      - docker build --rm -t local/project-image .
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  - name: build-project
    image: local/project-image
    commands:
      - ./build.sh
```
