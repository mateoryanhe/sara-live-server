#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

BUILD_DIR="$(cd .. && pwd)/go-build"
mkdir -p "$BUILD_DIR"

go build -o "$BUILD_DIR/xr-game-server" .
echo "Built: $BUILD_DIR/xr-game-server"
