#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
PROTO_DIR="${ROOT_DIR}/proto"

export PATH="/opt/homebrew/bin:${PATH}:${HOME}/go/bin:${ROOT_DIR}/node_modules/.bin"

cd "${PROTO_DIR}"

echo "Generating Go server code..."
buf generate --template buf.gen.go.yaml

echo "Generating TypeScript SDK..."
buf generate --template buf.gen.yaml

echo "Proto generation complete."
