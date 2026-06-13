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
└── docs/                # PRD, architecture, and UI design docs
```

## Quick Start

### Prerequisites

- Node.js >= 20
- Go >= 1.26
- PostgreSQL 15+ (for backend)
- buf CLI (for protocol buffer generation)

### Install

```bash
yarn install
```

### Run Backend

```bash
cd apps/backend
go run ./cmd/server
```

### Run Mobile

```bash
cd apps/mobile
yarn start
```

## License

MIT
