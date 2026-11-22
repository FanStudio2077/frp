# Docker Compose 部署（v0.73.0）

镜像（公开，可直接拉取）：
- `ghcr.io/fanstudio2077/frps:v0.73.0`（服务端，部署在公网服务器）
- `ghcr.io/fanstudio2077/frpc:v0.73.0`（客户端，部署在内网机器）

> 当前镜像为 **linux/amd64** 单架构。

## 服务端 frps
```bash
cd deploy/frps
# 编辑 frps.toml：auth.token、webServer 密码、subDomainHost 等
docker compose up -d
docker compose logs -f
```
- 管理面板：`http://<服务器IP>:7500`（账号见 frps.toml 的 webServer.user/password）
- 需放行：7000/tcp+udp、80、443，以及各代理的 remotePort
- 热重载配置：`docker kill -s HUP frps`（或面板「系统设置 → 重新加载配置」）
- 审计/RBAC/告警：面板「操作审计 / 用户管理 / 系统设置」

## 客户端 frpc
```bash
cd deploy/frpc
# 编辑 frpc.toml：serverAddr/serverPort、auth.token、servers、[[proxies]]
docker compose up -d
docker compose logs -f
```
- 管理面板：`http://<本机IP>:7400`
- 多服务器：改 `servers` 列表，或在面板「服务器」页维护

## 更新版本
```bash
docker compose pull && docker compose up -d
```
