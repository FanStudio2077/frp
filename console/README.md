# frp 管理控制台（console）

长期产品化的 frp 管理控制台：**一个面板**统一管理多台 frps、多台客户端节点（Agent），新建代理时选择「客户端节点 + 目标 frps 服务器」。frp 作为数据面复用现有镜像，**不修改 frp 内核**。

## 目录
```
console/
  backend/   # Go 控制面后端（独立 module: github.com/FanStudio2077/frp-console）
  web/       # Vue3 + Element Plus + Vite 前端
  Dockerfile.console / Dockerfile.agent
  docker-compose.yml
  HANDOFF_PROMPT.md   # 交接给下一个 AI 的完整任务书
```

## 后端
```bash
cd console/backend
go build ./...
go test ./...
go run ./cmd/console   # 默认 :8080，GET /healthz、GET /api/ping
```

## 前端
```bash
cd console/web
npm install
npm run build
npm run test:unit
npm run dev
```

## Docker
```bash
cd console
docker compose build console
docker compose up -d console
```

## 当前状态
脚手架：后端可编译（健康检查/路由骨架/数据模型/迁移 SQL）；前端可构建（占位页面、侧边栏、中英切换、i18n 骨架、smoke 测试）。
业务实现请依据 [HANDOFF_PROMPT.md](HANDOFF_PROMPT.md) 的 P1–P6 推进。
