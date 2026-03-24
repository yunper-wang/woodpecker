# Services

Woodpecker 在 YAML 文件中提供了 `services` 部分，用于定义 service 容器。
下面的配置组合了数据库和缓存容器。

Service 通过自定义 hostname 访问。
在以下示例中，MySQL service 被分配了 hostname `database`，可通过 `database:3306` 访问。

```yaml
steps:
  - name: build
    image: golang
    commands:
      - go build
      - go test

services:
  - name: database
    image: mysql

  - name: cache
    image: redis
```

你也可以显式定义端口和协议：

```yaml
services:
  - name: database
    image: mysql
    ports:
      - 3306

  - name: wireguard
    image: wg
    ports:
      - 51820/udp
```

## 停止

不再需要的 service 会收到 **SIGTERM** 信号。如果未响应，将强制以 **SIGKILL** 终止。
如果某些 service 无法正常关闭但这无关紧要，你可以直接忽略该错误：

```diff
 services:
   - name: database
     image: mysql
+    failure: ignore # we don't care how mysql exits
     ports:
       - 3306
```

## 配置

Service 容器通常通过环境变量来自定义启动参数，例如默认的用户名、密码和端口。请参阅官方镜像文档了解更多详情。

```diff
 services:
   - name: database
     image: mysql
+    environment:
+      MYSQL_DATABASE: test
+      MYSQL_ALLOW_EMPTY_PASSWORD: yes

   - name: cache
     image: redis
```

## 分离模式（Detachment）

Service 和长时间运行的容器也可以通过 `detach` 参数包含在 pipeline 的 `steps` 部分中，而不会阻塞其他步骤。当需要显式控制启动顺序时应使用此方式。

```diff
 steps:
   - name: build
     image: golang
     commands:
       - go build
       - go test

   - name: database
     image: redis
+    detach: true

   - name: test
     image: golang
     commands:
       - go test
```

分离步骤的容器将在 pipeline 结束时终止。

## 初始化

Service 容器需要一些时间来初始化并开始接受连接。如果无法连接到 service，可能需要等待几秒钟或实现退避重试逻辑。

```diff
 steps:
   - name: test
     image: golang
     commands:
+      - sleep 15
       - go get
       - go test

 services:
   - name: database
     image: mysql
```

## 完整 Pipeline 示例

```yaml
services:
  - name: database
    image: mysql
    environment:
      MYSQL_DATABASE: test
      MYSQL_ROOT_PASSWORD: example
steps:
  - name: get-version
    image: ubuntu
    commands:
      - ( apt update && apt dist-upgrade -y && apt install -y mysql-client 2>&1 )> /dev/null
      - sleep 30s # need to wait for mysql-server init
      - echo 'SHOW VARIABLES LIKE "version"' | mysql -u root -h database test -p example
```
