#!/usr/bin/env bash
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$DIR"
if [[ ! -f config.json ]]; then
  echo "missing config.json (copy from config.example.json)"
  exit 1
fi
exec ./live-hot-restart-gw
