# 宝塔 Docker Compose 一键部署（v0.73.0）

在宝塔面板「Docker → Compose → 新建」中，把其中一份 YAML 整段粘贴进去即可（配置已内联，无需额外文件）。

## 服务端 frps（部署在公网服务器）
复制 `frps.compose.yml` 内容，按需修改：
- `auth.token` → 你的连接令牌
- `webServer.password` / `webServer.users` → 面板账号
- `subDomainHost` → 你的域名
- 端口：7000(tcp+udp)、80、443、7500（面板）
在宝塔「安全」放行上述端口。

面板地址：`http://服务器IP:7500`

## 客户端 frpc（部署在内网机器）
复制 `frpc.compose.yml` 内容，修改：
- `serverAddr` → 服务器公网 IP
- `auth.token` → 与服务端一致
- `[[proxies]]` → 你要暴露的服务（localPort/remotePort）
面板：`http://本机IP:7400`

> 镜像为 linux/amd64 单架构。若为 ARM 服务器请告知。
