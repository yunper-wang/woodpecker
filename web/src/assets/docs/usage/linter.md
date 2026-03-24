# Linter

Woodpecker 会自动对你的 workflow 文件进行 lint 检查，检测错误、已废弃的用法和不良习惯。任何 pipeline 的错误和警告都会在 UI 中显示。

![errors and warnings in UI](./linter-warnings-errors.png)

## 通过 CLI 运行 linter

你也可以通过 CLI 手动运行 linter：

```shell
woodpecker-cli lint <workflow files>
```

## 不良习惯警告

如果你的配置中包含某些不良习惯，Woodpecker 会发出警告。

### 为所有步骤添加事件过滤器

`when` 块中的所有条目都应包含 `event` 过滤器，确保没有步骤在所有事件上都运行。这样做的原因是：如果未来新增了事件类型，你的步骤可能并不应该在那些事件上运行。

**不正确**的配置示例：

```yaml
when:
  - branch: main
  - event: tag
```

这会触发警告，因为第一个条目（`branch: main`）没有使用事件过滤。

```yaml
steps:
  - name: test
    when:
      branch: main

  - name: deploy
    when:
      event: tag
```

**正确**的配置示例：

```yaml
when:
  - branch: main
    event: push
  - event: tag
```

```yaml
steps:
  - name: test
    when:
      event: [tag, push]

  - name: deploy
    when:
      - event: tag
```
