# Still

A photo mood app — express what lingers.

> Still doesn't tell you how you feel. It helps you express it.

## What is Still?

Still is a quiet, photo-based community for emotional expression.

1. Upload a photo.
2. AI reads the image and suggests a mood, title, and short description.
3. Edit if you want — the user is always the author.
4. Share to the community.
5. Others can **Resonate** with your moment.

## Project Structure

```text
still/
├── apps/
│   ├── mobile/          # Expo React Native app
│   └── backend/         # Go + ConnectRPC API server
├── proto/
│   └── still/v1/        # Protocol buffer definitions
├── packages/
│   ├── design-system/   # Shared UI tokens and components
│   ├── generated-sdk/   # Generated TypeScript client SDK
│   └── shared-types/    # Shared TypeScript constants and types
├── scripts/             # Development and codegen scripts
├── docs/                # PRD, architecture, and UI design docs
└── Makefile             # Common development tasks
```

## Tech Stack

| Layer           | Tech                                                                                     |
| --------------- | ---------------------------------------------------------------------------------------- |
| Mobile          | Expo SDK 54 (React Native 0.81) + TypeScript + React Navigation + Zustand                |
| Backend         | Go + ConnectRPC + pgx/v5 + golang-migrate                                                |
| API             | Protocol Buffers + buf                                                                   |
| Database        | PostgreSQL                                                                               |
| Storage         | S3 / S3-compatible (dev: MinIO, prod: AWS S3 / Cloudflare R2)                            |
| AI              | CloudWeGo Eino + multi-provider LLM layer (OpenAI, DeepSeek, Moonshot, Qwen)             |
| Auth            | Clerk                                                                                    |
| Observability   | zerolog + OpenTelemetry + Sentry                                                         |
| CI / Deployment | GitHub Actions + Docker + Fly.io / Railway / Render configs                              |

## Quick Start

### Prerequisites

- Node.js >= 20
- Go >= 1.26
- Docker + Docker Compose (for PostgreSQL and MinIO)
- buf CLI (for protocol buffer generation)
- OpenAI API key

### 1. Install Dependencies

```bash
make install
```

### 2. Configure Environment

```bash
make env
```

This creates dev env files from `.env.example`:

- `apps/backend/.env.development`
- `apps/mobile/.env.development`

Edit `apps/backend/.env.development` and set your `LLM_API_KEY` (legacy `OPENAI_API_KEY` still works when using OpenAI).

### 3. Start Infrastructure

```bash
make infra
```

Starts PostgreSQL and MinIO (S3-compatible storage for local dev).

### 4. Start Backend

```bash
make backend
```

The backend loads `apps/backend/.env.development` and runs on http://localhost:8080.

### 5. Start Mobile App

```bash
make mobile
```

The mobile app loads `apps/mobile/.env.development` and connects to the local backend.

## Common Commands

```bash
make help          # Show all available commands
make build         # Build all workspaces
make lint          # Run linters
make test          # Run tests
make proto         # Regenerate Go + TypeScript SDK from proto files
make migrate       # Run database migrations up
make migrate-down  # Rollback one migration
make infra-down    # Stop local infrastructure
make clean         # Stop infrastructure and remove volumes
```

## Production

1. Fill `apps/backend/.env.production` with real credentials:
   - PostgreSQL `DATABASE_URL`
   - `LLM_PROVIDER` and `LLM_API_KEY` (see `.env.example` for supported providers)
   - `CLERK_SECRET_KEY`
   - S3 credentials (`S3_ENDPOINT`, `S3_REGION`, `S3_BUCKET`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`)
   - Optional: `SENTRY_DSN`, `OTEL_EXPORTER_OTLP_ENDPOINT`
2. Set `apps/mobile/.env.production` `EXPO_PUBLIC_API_BASE_URL` to your production API domain.
3. Build and deploy the backend and mobile app as usual.

See `docs/` for PRD, architecture decisions, and deployment guides.

## Contributing

We welcome contributions to the Community Edition. Please read
[CONTRIBUTING.md](CONTRIBUTING.md) and [CLA.md](CLA.md) before opening a PR.

## Commercial Edition

This repository is the open-source **Community Edition**. For the official hosted
service, premium AI features, and enterprise support, see
[COMMERCIAL.md](COMMERCIAL.md) or visit https://still.app (placeholder).

## License

[Apache License 2.0](LICENSE)
