#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
MIGRATIONS_DIR="${ROOT_DIR}/apps/backend/migrations"

DATABASE_URL="${DATABASE_URL:-postgres://localhost/still?sslmode=disable}"

case "${1:-up}" in
  up)
    echo "Running migrations up..."
    migrate -path "${MIGRATIONS_DIR}" -database "${DATABASE_URL}" up
    ;;
  down)
    echo "Running migrations down..."
    migrate -path "${MIGRATIONS_DIR}" -database "${DATABASE_URL}" down 1
    ;;
  create)
    if [ -z "${2:-}" ]; then
      echo "Usage: $0 create <migration_name>"
      exit 1
    fi
    migrate create -ext sql -dir "${MIGRATIONS_DIR}" -seq "$2"
    ;;
  *)
    echo "Usage: $0 {up|down|create <name>}"
    exit 1
    ;;
esac
