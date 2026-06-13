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

| 层级 | 技术 |
|------|------|
| Mobile | Expo (React Native) + TypeScript + React Navigation + Zustand |
| Backend | Go 1.26 + ConnectRPC + pgx/v5 + golang-migrate + go-openai |
| API | Protocol Buffers + buf |
| Database | PostgreSQL |
| Storage | S3 / S3-compatible (dev: MinIO, prod: AWS S3 / R2) |
| AI | OpenAI / Claude / Gemini（统一抽象） |
| Auth | Clerk（推荐）或 Supabase Auth |
| Observability | zerolog + OpenTelemetry + Sentry |

## Key Principles

1. **Proto First**：所有接口先在 `proto/still/v1/` 中定义，再生成 Go Server 与 TS Client。
2. **Mono Service V1**：`apps/backend` 是一个单体服务，不要微服务。
3. **Mood Dictionary**：AI 必须从固定 Mood 词库中选择，不能自由生成。
4. **Quiet UI**：安静、轻盈、留白。避免科技感、AI 感、社交媒体感。
5. **User is Author**：AI 只是辅助，用户可编辑 mood/title/description。

## Development Commands

```bash
# 查看所有可用命令
make help

# 一键：安装依赖、生成环境文件、启动基础设施
make install
make env        # 然后编辑 apps/backend/.env.development 填入 OPENAI_API_KEY
make infra

# 启动后端（端口 8080）和移动端
make backend    # 读取 apps/backend/.env.development
make mobile     # 读取 apps/mobile/.env.development

# 其他常用命令
make proto      # 生成 Proto SDK
make build      # 构建所有 workspace
make lint       # 运行 linter
make test       # 运行测试
make migrate    # 数据库迁移 up
make clean      # 停止基础设施并清理卷
```

## Code Style

- Go：标准 Go 风格，`gofmt` + `go vet`。
- TypeScript：严格模式，`strict: true`。
- 文件命名：kebab-case 用于配置/脚本，camelCase/PascalCase 用于代码。

## When Modifying This File

如果你改变了项目结构、技术栈、构建命令或工作区配置，请同步更新本文件。
