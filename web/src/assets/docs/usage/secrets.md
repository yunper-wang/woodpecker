# Secrets

Woodpecker 提供了在中央 secret 存储中保存命名变量的能力。
这些 secret 可以通过 `from_secret` 关键字安全地传递给 pipeline 的各个步骤。

secret 共有三个级别。如果在多个级别定义了同名 secret，则按以下优先级处理（后者优先）：

1. **仓库级 secret（Repository secrets）**：对该仓库的所有 pipeline 可用。
1. **组织级 secret（Organization secrets）**：对该组织的所有 pipeline 可用。
1. **全局 secret（Global secrets）**：只有实例管理员才能设置。
   全局 secret 对整个 Woodpecker 实例的**所有** pipeline 可用，因此需谨慎使用。

除了原生的 secret 集成方式外，也可以在 pipeline 步骤中直接与外部 secret 提供商交互来使用外部 secret。访问这些提供商所需的凭据可通过 Woodpecker secret 进行配置，从而实现从各自外部来源检索 secret。

:::warning
Woodpecker 可以对来自自身 secret 存储的 secret 进行脱敏处理，但无法对外部 secret 提供同等保护。因此，这些外部 secret 可能会在 pipeline 日志中暴露。
:::

## 使用方式

你可以通过 `from_secret` 语法将 Woodpecker secret 设置为某个配置项或环境变量的值。

以下示例将名为 `secret_token` 的 secret 通过环境变量 `TOKEN_ENV` 传递给步骤：

```diff
 steps:
   - name: 'step name'
     image: registry/repo/image:tag
     commands:
+      - echo "The secret is $TOKEN_ENV"
+    environment:
+      TOKEN_ENV:
+        from_secret: secret_token
```

同样的语法也可用于将 secret 传递给插件的 settings。
名为 `secret_token` 的 secret 被赋值给 setting `TOKEN`，该 setting 在插件中以环境变量 `PLUGIN_TOKEN` 的形式可用（详见[插件](./51-plugins/20-creating-plugins.md#settings)文档）。
`PLUGIN_TOKEN` 随后由插件内部使用，并在执行时生效。

```diff
 steps:
   - name: 'step name'
     image: registry/repo/image:tag
+    settings:
+      TOKEN:
+        from_secret: secret_token
```

### 转义 secret

请注意，参数表达式会被预处理，即在 pipeline 启动前就会被求值。
如果需要在表达式中使用 secret，必须通过 `$$` 进行适当转义，以确保正确处理。

```diff
 steps:
   - name: docker
     image: docker
     commands:
-      - echo ${TOKEN_ENV}
+      - echo $${TOKEN_ENV}
     environment:
       TOKEN_ENV:
         from_secret: secret_token
```

### 事件过滤

默认情况下，secret 不会暴露给 pull request。
但你可以通过创建 secret 并启用 `pull_request` 事件类型来改变此行为。
这可以通过 UI 或 CLI 进行配置。

:::warning
在 pull request 中暴露 secret 时需格外谨慎。
如果你的仓库是公开的且接受所有人提交 pull request，你的 secret 可能面临风险。
恶意用户可能借此机会暴露你的 secret 或将其传输到外部位置。
:::

### 插件过滤

为防止 secret 被恶意用户滥用，你可以将 secret 限制为只能由特定插件列表使用。
启用后，这些 secret 对其他任何插件都不可用。
插件的优势在于它们无法执行任意命令，因此无法泄露 secret。

:::tip
如果你指定了 tag，过滤器会将其纳入考量。
但是，如果同一镜像在列表中出现多次，权限最低的条目将优先生效。
例如，没有 tag 的镜像将允许所有 tag，即使列表中还有另一个带有 tag 的条目。
:::

![plugins filter](./secrets-plugins-filter.png)

## CLI

除了 UI，还可以使用 CLI 管理 secret。

使用默认设置创建 secret。
该 secret 对 pipeline 中的所有镜像以及所有 `push`、`tag` 和 `deployment` 事件可用（不包括 `pull_request` 事件）。

```bash
woodpecker-cli repo secret add \
  --repository octocat/hello-world \
  --name aws_access_key_id \
  --value <value>
```

创建 secret 并将其限制为单个镜像：

```diff
 woodpecker-cli secret add \
   --repository octocat/hello-world \
+  --image woodpeckerci/plugin-s3 \
   --name aws_access_key_id \
   --value <value>
```

创建 secret 并将其限制为一组镜像：

```diff
 woodpecker-cli repo secret add \
   --repository octocat/hello-world \
+  --image woodpeckerci/plugin-s3 \
+  --image woodpeckerci/plugin-docker-buildx \
   --name aws_access_key_id \
   --value <value>
```

创建 secret 并为多个钩子事件启用它：

```diff
 woodpecker-cli repo secret add \
   --repository octocat/hello-world \
   --image woodpeckerci/plugin-s3 \
+  --event pull_request \
+  --event push \
+  --event tag \
   --name aws_access_key_id \
   --value <value>
```

可以使用 `@` 语法从文件加载 secret。
建议使用此方法从文件加载 secret，因为它能确保保留换行符（这对 SSH 密钥等尤为重要）：

```diff
 woodpecker-cli repo secret add \
   -repository octocat/hello-world \
   -name ssh_key \
+  -value @/root/ssh/id_rsa
```
