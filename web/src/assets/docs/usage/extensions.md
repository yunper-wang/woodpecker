# Extensions

Woodpecker 允许你通过使用预定义的 HTTP 端点，将内部逻辑替换为外部 extension。

目前支持以下类型的 extension：

- [Configuration extension](./configuration-extension.md)：动态修改或生成 pipeline 配置。
- [Registry extension](./registry-extension.md)：从 extension 获取 registry 凭证。

## 安全性

:::warning
你需要信任这些 extension，因为它们会接收 secret 和 token 等私密信息，
并可能返回有害数据，例如可能被执行的恶意 pipeline 配置。
:::

为防止你的 extension 遭受此类攻击，Woodpecker 使用 [HTTP signatures](https://tools.ietf.org/html/draft-cavage-http-signatures) 对所有 HTTP 请求进行签名。Woodpecker 为此使用一对公私 ed25519 密钥对。
要验证请求，你的 extension 必须使用公钥并借助类似 [httpsign](https://github.com/yaronf/httpsign) 的库来验证所有请求的签名。
你可以通过访问 `http://my-woodpecker.tld/api/signature/public-key` 或在 Woodpecker UI 中进入仓库设置并打开 extensions 页面来获取 Woodpecker 公钥。

## 示例 Extension

提供 config 和 secrets extension 端点的简单示例服务可在此处找到：[https://github.com/woodpecker-ci/example-extensions](https://github.com/woodpecker-ci/example-extensions)

## 配置

为防止 extension 调用本地服务，默认情况下只允许访问外部主机/IP 地址。你可以通过设置 `WOODPECKER_EXTENSIONS_ALLOWED_HOSTS` 环境变量来修改此行为。可以使用逗号分隔的列表，包含以下内容：

- 内置网络：
  - `loopback`：IPv4 的 127.0.0.0/8 和 IPv6 的 ::1/128，包含 localhost。
  - `private`：RFC 1918（10.0.0.0/8、172.16.0.0/12、192.168.0.0/16）和 RFC 4193（FC00::/7），也称为局域网/内网。
  - `external`：有效的非私有单播 IP，可以访问公共互联网上的所有主机。
  - `*`：允许所有主机。
- CIDR 列表：IPv4 用 `1.2.3.0/8`，IPv6 用 `2001:db8::/32`
- （通配符）主机：`example.com`、`*.example.com`、`192.168.100.*`
