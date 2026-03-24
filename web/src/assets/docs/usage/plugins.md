# Plugins

Plugins 是 pipeline step，执行预定义的任务，并作为 step 配置在你的 pipeline 中。
Plugins 可用于部署代码、发布 artifacts、发送通知等。

它们会从 agent 配置的默认容器 registry 自动拉取。

```dockerfile title="Dockerfile"
FROM cloud/kubectl
COPY deploy /usr/local/deploy
ENTRYPOINT ["/usr/local/deploy"]
```

```bash title="deploy"
kubectl apply -f $PLUGIN_TEMPLATE
```

```yaml title=".woodpecker.yaml"
steps:
  - name: deploy-to-k8s
    image: cloud/my-k8s-plugin
    settings:
      template: config/k8s/service.yaml
```

使用 Prettier 和 S3 plugin 的示例 pipeline：

```yaml
steps:
  - name: build
    image: golang
    commands:
      - go build
      - go test

  - name: prettier
    image: woodpeckerci/plugin-prettier

  - name: publish
    image: woodpeckerci/plugin-s3
    settings:
      bucket: my-bucket-name
      source: some-file-name
      target: /target/some-file
```

## Plugin 隔离

Plugins 本质上就是 pipeline step。它们共享构建工作区（以 volume 方式挂载），因此可以访问你的源代码树。
普通 step 可以执行任意代码，而 plugin 应只允许 plugin 作者预定义的功能。

因此存在一些限制。工作区根目录始终挂载在 `/woodpecker`，但工作目录会动态调整，
作为 plugin 的使用者，你无需关心这些细节。此外，plugin 不能与 `commands`
或 `entrypoint` 一起使用，否则会报错。使用 `environment` 是可以的，但这种情况下，该 plugin 在内部将不再被当作 plugin 处理，
容器将无法再通过 plugin 过滤器访问 secret，也不会在没有明确定义的情况下具有特权。

## 查找 Plugin

对于官方 plugin，可以使用 Woodpecker plugin index：

- [Official Woodpecker Plugins](https://woodpecker-ci.org/plugins)

:::tip
还有其他 plugin 列表提供更多 plugin。请注意，[Drone](https://www.drone.io/) plugin 通常是兼容的，但可能需要一些调整和修改。

- [Drone Plugins](http://plugins.drone.io)
- [Geeklab Woodpecker Plugins](https://woodpecker-plugins.geekdocs.de/)
- [Woodpecker Community Plugins](https://codeberg.org/woodpecker-community)

:::
