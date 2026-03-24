# Volumes

Woodpecker 支持在 YAML 中定义 Docker volumes。你可以使用此参数将主机上的文件或文件夹挂载到容器中。

:::note
Volumes 仅对受信任的仓库可用，出于安全原因应仅在私有环境中使用。请参阅[项目设置](./75-project-settings.md#trusted)以启用受信任模式。
:::

```diff
 steps:
   - name: build
     image: docker
     commands:
       - docker build --rm -t octocat/hello-world .
       - docker run --rm octocat/hello-world --test
       - docker push octocat/hello-world
       - docker rmi octocat/hello-world
     volumes:
+      - /var/run/docker.sock:/var/run/docker.sock
```

如果你使用 Docker backend，还可以使用命名 volume，例如 `some_volume_name:/var/run/volume`。

请注意，Woodpecker 在主机上挂载 volume，因此配置 volume 时必须使用绝对路径。使用相对路径将导致报错。

```diff
-volumes: [ ./certs:/etc/ssl/certs ]
+volumes: [ /etc/ssl/certs:/etc/ssl/certs ]
```
