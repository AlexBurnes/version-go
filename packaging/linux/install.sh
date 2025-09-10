#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-/usr/local/bin}"
SELF=$(readlink -f "$0" || echo "$0")
BASE_DIR=$(dirname "$SELF")

echo "[*] Target bin dir: $APP_DIR"

install_one() {
  local src="$1"
  local dst="$APP_DIR/$(basename "$src")"
  if command -v sudo >/dev/null 2>&1; then
    sudo install -m 0755 "$src" "$dst"
  else
    install -m 0755 "$src" "$dst"
  fi
  echo "  + $(basename "$src") -> $dst"
}

# Бинарь называется 'version'
if [[ -f "$BASE_DIR/version" ]]; then
  install_one "$BASE_DIR/version"
else
  echo "version binary not found near install.sh"; exit 1
fi

echo "[✓] Done."