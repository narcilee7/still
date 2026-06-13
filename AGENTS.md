# Still — Agent Guide

## Project Overview

Still 是一个通过照片表达当下情绪的美学社区。

- **产品定位**：介于 Instagram / VSCO / Tumblr 之间的轻量情绪表达产品。
- **技术定位**：AI Native Consumer App，核心链路是「照片 → AI 情绪解读 → 社区分发」。

## Repository Layout

Monorepo，使用 Yarn workspaces。

```text
still/
├── apps/
│   ├── mobile/          # Expo React Native (TypeScript): navigation / screens / store / services / data
│   └── backend/         # Go + ConnectRPC 单体服务：internal/{db,repository,service} / cmd/server
├── proto/
│   └── still/v1/        # API 协议定义（Proto First）
├── packages/
│   ├── design-system/   # 共享 UI tokens / theme
│   ├── generated-sdk/   # buf 生成的 TypeScript SDK
│   └── shared-types/    # 纯 TypeScript 共享类型/常量
├── scripts/             # 开发/生成/迁移脚本
└── docs/                # PRD / Tech / UI 文档
```

## Tech Stack

| 层级          | 技术                                                                              |
| ------------- | --------------------------------------------------------------------------------- |
| Mobile        | Expo SDK 54 (React Native 0.81) + TypeScript + React Navigation + Zustand + Clerk |
| Backend       | Go 1.26 + ConnectRPC + pgx/v5 + golang-migrate + go-openai                        |
| API           | Protocol Buffers + buf                                                            |
| Database      | PostgreSQL                                                                        |
| Storage       | S3 / S3-compatible (dev: MinIO, prod: AWS S3 / R2)                                |
| AI            | OpenAI / Claude / Gemini（统一抽象）                                              |
| Auth          | Clerk（JWT 验证，后端 `internal/auth`）                                           |
| Observability | zerolog + OpenTelemetry + Sentry                                                  |

## Key Principles

1. **Proto First**：所有接口先在 `proto/still/v1/` 中定义，再生成 Go Server 与 TS Client。
2. **Mono Service V1**：`apps/backend` 是一个单体服务，不要微服务。
3. **Mood Dictionary**：AI 必须从固定 Mood 词库中选择，不能自由生成。
4. **Quiet UI**：安静、轻盈、留白。避免科技感、AI 感、社交媒体感。
5. **User is Author**：AI 只是辅助，用户可编辑 mood/title/description。
6. **Auth Required**：V1 除 `/health` 外，所有 RPC 都需要 Clerk Bearer Token；未登录用户进入 `AuthStack` 登录。

## Development Commands

```bash
# 查看所有可用命令
make help

# 一键：安装依赖、生成环境文件、启动基础设施
make install
make env        # 然后编辑 apps/backend/.env.development 填入 OPENAI_API_KEY 和 CLERK_SECRET_KEY
                # 编辑 apps/mobile/.env.development 填入 EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY
make infra-up   # 启动本地 PostgreSQL + MinIO

# 启动后端（端口 8080）和移动端
make backend    # 读取 apps/backend/.env.development
make mobile     # 读取 apps/mobile/.env.development

# Docker / 部署
make docker-build      # 构建后端 Docker 镜像
make docker-prod-up    # 本地启动类生产容器栈

# 其他常用命令
make proto      # 生成 Proto SDK
make build      # 构建所有 workspace
make lint       # 运行 linter
make test       # 运行测试
make migrate    # 数据库迁移 up
make clean      # 停止基础设施并清理卷
```

## Authentication

- 后端通过 `internal/auth` 验证 Clerk JWT，并将 `clerk_user_id` 写入请求上下文。
- `users.clerk_user_id` 映射到 Clerk 的 `user_xxx`，内部仍使用 UUID 主键。
- `UserService.GetMe` 会根据 Clerk ID 自动创建/读取本地用户记录。
- 移动端使用 `@clerk/expo` + `expo-secure-store` 管理会话，`TokenBridge` 把 `getToken()` 注入 API transport 的 `Authorization` 头。

## Deployment

- 后端镜像：`apps/backend/Dockerfile`（多阶段构建，Alpine 运行）。
- 本地类生产栈：`docker-compose.prod.yml`。
- 平台配置：`fly.toml`、`render.yaml`、`railway.toml`。
- 详细说明见 `docs/Tech/deployment.md`。

## Code Style

- Go：标准 Go 风格，`gofmt` + `go vet`。
- TypeScript：严格模式，`strict: true`。
- 文件命名：kebab-case 用于配置/脚本，camelCase/PascalCase 用于代码。

## When Modifying This File

如果你改变了项目结构、技术栈、构建命令或工作区配置，请同步更新本文件。
