# frp 管理控制台 —— 交接提示词（HANDOFF PROMPT）

> 你是接续开发的 AI。基于本仓库 `console/` 脚手架继续开发一个**长期产品化**的「frp 管理控制台」。先通读本文件再编码。

## 1. 一句话定位
一个 Web 管理控制台（控制面）：用**一个面板**统一管理多台 frps 服务器、多台客户端节点（Agent），新建代理时选择「客户端节点 + 目标 frps 服务器」。frp 作为数据面复用现有镜像，**不要修改 frp 内核**。

## 2. 已有基础（本仓库是 frp 定制 fork）
- frp `0.73.0`，镜像 `ghcr.io/fanstudio2077/frps:v0.73.0`、`ghcr.io/fanstudio2077/frpc:v0.73.0`（公开、linux/amd64）。
- frps/frpc 面板已具备：中英 i18n、操作审计、RBAC(admin/viewer)、告警阈值、配置热重载、frpc Store。控制台要把这些能力集中化。
- frps 管理员 API（Basic Auth）：`/api/v2/system/info`、`/api/v2/clients`、`/api/v2/clients/{key}`、`/api/v2/clients/{key}/disconnect`、`/api/v2/proxies`、`/api/v2/proxies/{name}`、`/api/v2/proxies/{name}/traffic`、`/api/v2/users`、`/api/v2/alerts`、`/api/v2/audit/logs`、`/api/v2/system/config`、`/api/v2/system/reload`、`/api/v2/web/users`、`DELETE /api/v2/proxies`。
- frpc 管理员 API（Basic Auth）：`/api/status`、`/api/config`(GET/PUT)、`/api/reload`、`/api/stop`、`/api/store/proxies`(GET/POST/PUT/DELETE)、`/api/store/visitors`、`/api/proxy/{name}/config`、`/api/visitor/{name}/config`、`/api/servers`。

## 3. 目标与验收标准
1. 一个面板添加 N 台 frps、M 台客户端节点（各装 Agent），展示在线状态。
2. 新建代理可选「节点 + 目标服务器」，保存后该代理真实出现在对应 frps 上（由该节点上绑定该服务器的 frpc 实例承载）。
3. 概览聚合所有服务器的客户端数/代理数/连接数/流量/告警。
4. 登录 + RBAC(admin/operator/viewer)、操作审计、阈值告警、中英双语。
5. 后端 `go test`、前端 `vitest` + 构建通过；`docker compose` 可一键起控制台，Agent 可部署到客户端机器。

## 4. 架构
```
浏览器 --HTTP--> console(控制面: API+DB+Web) <--WS(出站)-- console-agent(每台客户端机器)
                   | 轮询聚合 frps /api/v2/*                      | 通过 Docker API 管理
                   +-- RBAC/审计/告警                            +-- frpc实例(→frps-1)、frpc实例(→frps-2)…（现有镜像）
```
- 单一面板 = console 前端；按代理选服务器 = console 把代理下发到该节点上绑定目标服务器的 frpc 实例，无需改 frp 内核。
- Agent：每台客户端机器一个轻量 Agent（容器，host 网络），主动外联 console 的 WebSocket（NAT 友好），用本机 Docker API 按 (节点,服务器) 起停/更新标准 frpc 容器。
- 动态代理：优先用 frpc Store + admin `/api/store/proxies` 增删改，无需重启。

## 5. 技术栈与目录
- 后端：Go 1.26，标准库 `net/http`（Go 1.22+ 路由模式）；SQLite（`modernc.org/sqlite`）；bcrypt；JWT/会话 Cookie。
- 前端：Vue3 + Element Plus + Vite + vue-i18n；可复用本仓库 `web/shared`（vite alias `@shared → ../../web/shared`）与 frpc 代理表单（`@frpc → ../../web/frpc/src`）。
- 目录：`console/backend`（独立 module `github.com/FanStudio2077/frp-console`）、`console/web`。

## 6. 数据模型（见 `console/backend/internal/store/migrations/0001_init.sql`）
- users(id, username unique, password_hash, role, disabled, created_at)
- servers(id, name, addr, port, token_enc, tls_json, remark, created_at)  # token 加密存储
- nodes(id, name, token_hash, ip, os, arch, version, online, last_seen, created_at)
- instances(id, node_id, server_id, name, admin_port, status, created_at) UNIQUE(node_id,server_id)
- proxies(id, instance_id, name, type, local_ip, local_port, remote_port, subdomain, custom_domains, extra_json, enabled, created_at, updated_at) UNIQUE(instance_id,name)
- visitors(...) 同构；audit_logs / alert_rules / alert_events / settings
- 服务器 token 用主密钥 AES-GCM 加密；API 永不回显明文。

## 7. 后端 API 规范（`/api`，统一 JWT+RBAC）
- 认证：POST /api/auth/login、POST /api/auth/logout、GET /api/auth/me
- 用户：GET/POST/PUT/DELETE /api/users（admin）
- 服务器：GET/POST/PUT/DELETE /api/servers、POST /api/servers/{id}/test
- 节点：GET /api/nodes、GET /api/nodes/{id}、DELETE /api/nodes/{id}、POST /api/nodes/{id}/rotate-token、GET /api/agent/ws
- 实例：GET /api/instances、POST /api/instances/ensure、POST /api/instances/{id}/reload
- 代理/访问者：GET/POST/PUT/DELETE /api/proxies、/api/visitors（body 含 nodeId+serverId）
- 聚合：GET /api/overview、GET /api/servers/{id}/clients|proxies|traffic、GET /api/clients
- 审计：GET /api/audit；告警：GET/POST/PUT/DELETE /api/alert-rules、GET /api/alert-events

## 8. Agent 协议（WebSocket JSON）
消息：register(节点身份/机器信息)、heartbeat、ensure_instance(nodeId,serverId,frpcConfig)、remove_instance、apply_proxy(instanceId,proxyDef)、remove_proxy、instance_status。
流程：Agent 启动 -> 带 token 连 /api/agent/ws -> register -> 周期 heartbeat -> 收到 ensure_instance 用 Docker API 拉起/更新 frpc 容器（host 网络，admin 绑定 127.0.0.1:adminPort，启用 store）-> apply_proxy 时调该 frpc admin /api/store/proxies。

## 9. 关键流程：新建代理
1. 前端选 节点A + 服务器1，填参数 -> POST /api/proxies。
2. 后端入库；确保 instances(A,1) 存在（不存在则经 Agent ensure 拉起 frpc 实例）。
3. 下发 apply_proxy；Agent 调本机 frpc admin POST /api/store/proxies 动态生效。
4. 聚合刷新，代理出现在 server1 的 frps 与统一列表。

## 10. 前端页面（Vue3 + Element Plus + i18n）
侧栏：概览 / 服务器 / 节点 / 代理 / 访问者 / 用户 / 审计 / 告警 / 设置。
- 登录页；401 跳登录。
- 概览：跨服务器聚合卡片 + 告警 + 各服务器状态。
- 服务器：CRUD、连通性测试、详情（客户端/代理）。
- 节点：列表/在线状态/详情（承载的 frpc 实例）。
- 代理：统一列表（按服务器/节点筛选）；新建/编辑弹窗含「客户端节点」+「目标服务器」下拉，类型字段复用 frpc 代理表单。
- 用户(RBAC)/审计/告警/设置。

## 11. 分阶段任务（每阶段可构建/测试/验收）
- P1 后端基座：SQLite store + 迁移执行；config；登录/JWT/RBAC；用户 CRUD；服务器 CRUD + 连通性测试；审计。
- P2 Agent：Agent 容器；注册/心跳/WS 隧道；Docker API 管理 frpc 实例。
- P3 代理/访问者：CRUD + 选择节点/服务器 + 经 Agent 下发到 frpc store；前端统一列表与表单。
- P4 聚合：轮询 frps /api/v2/* 与 frpc 状态；概览与详情页。
- P5 告警 + i18n 完善。
- P6 打包交付：Dockerfile/compose、文档、e2e。

## 12. 构建/测试命令
```bash
cd console/backend && go build ./... && go test ./...
cd console/web && npm install && npm run build && npm run test:unit
# 前端构建若遇 esbuild 临时目录权限：设置 TMP/TEMP 到工作区内目录
```

## 13. 编码/提交约定
- 分支 dev；conventional commits（feat/fix/docs/ci）。
- 后端分层：config/store/model/api/auth/audit/aggregate/agent/agentproto/crypto/genconfig。
- 前端复用 web/shared；i18n 消息放 console/web/src/i18n，默认中文可切英文。
- 不要修改 frp 内核；数据面用现有镜像。token 加密存储、接口脱敏。

## 14. 当前状态（本轮脚手架已完成）
- console/backend：可编译；cmd/console（/healthz、/api/ping）、internal/{config,store,model,api}、migrations/0001_init.sql、占位包 auth/audit/aggregate/agent/agentproto/crypto/genconfig。
- console/web：可构建；侧边栏 + 9 个占位页 + 中英切换 + i18n 骨架 + smoke 测试；vite alias 已配 @shared/@frpc。
- Dockerfile.console、Dockerfile.agent、docker-compose.yml、README.md。
- 尚未实现：真实 API/DB 查询、Agent、聚合、告警。

## 15. 下一步（从这里开始）
1. 实现 internal/store 的 SQLite 打开与迁移执行（modernc.org/sqlite），跑通建表。
2. 实现 internal/auth（bcrypt + JWT 会话 + RBAC 中间件）与 /api/auth/*、/api/users。
3. 实现 /api/servers CRUD + /test（调用 frps /api/v2/system/info，token 加密存储）。
4. 同步补前端登录页、服务器页真实接口。

## 16. 注意事项/风险
- Agent 穿越 NAT：出站 WS。
- 服务器 token 安全：AES-GCM 加密，接口脱敏。
- 规模：聚合缓存 + 可调轮询；后续可换 Postgres（Store 接口已隔离）。
- 代理生效优先用 frpc Store 动态 API；无 Store 时回退生成配置+reload。
