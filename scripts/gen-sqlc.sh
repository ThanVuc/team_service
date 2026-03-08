#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

SQLC_DIR="$ROOT_DIR/internal/infrastructure/persistence/db"
GEN_DIR="$ROOT_DIR/internal/infrastructure/persistence/db/database"

echo "==> Cleaning generated sqlc directory: $GEN_DIR"
rm -rf "$GEN_DIR"
mkdir -p "$GEN_DIR"

echo "==> Running sqlc generate"
cd "$SQLC_DIR"
sqlc generate

echo "==> sqlc generation complete"
