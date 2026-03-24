# Cron

配置 cron 任务至少需要对该仓库拥有推送权限。

## 添加新的 cron 任务

1. 要创建新的 cron 任务，请修改你的 pipeline 配置文件，并为所有希望由该 cron 任务触发的步骤添加事件过滤器：

   ```diff
    steps:
      - name: sync_locales
        image: weblate_sync
        settings:
          url: example.com
          token:
            from_secret: weblate_token
   +    when:
   +      event: cron
   +      cron: "name of the cron job" # 如果只想由特定 cron 任务执行此步骤，请填写任务名称
   ```

2. 在仓库设置中创建新的 cron 任务：

   ![cron settings](./cron-settings.png)

   支持的调度语法请参阅 <https://pkg.go.dev/github.com/gdgvda/cron#hdr-CRON_Expression_Format>。如需了解 cron 语法的基础知识，<https://it-tools.tech/crontab-generator> 是一个很好的入门和实验平台。

   示例：`@every 5m`、`@daily`、`30 * * * *` ...
