#!/bin/bash

# Скрипт для конвертации набора geoip/geosite dat в SRS файлы
# Список источников задается аргументами запуска локальной утилиты

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="$PROJECT_DIR/rules"
CONVERTER_DIR="$PROJECT_DIR/tools/geodat2srs"

echo "=== Build SRS Files ==="

mkdir -p "$OUTPUT_DIR"

echo "[1/3] Cleaning old SRS files..."
rm -f "$OUTPUT_DIR"/*.srs

echo "[2/3] Running bundled geodat2srs..."
cd "$CONVERTER_DIR"

echo "[3/3] Converting to SRS format..."
go run . \
  -output-dir "$OUTPUT_DIR" \
  -source "geoip:geoip=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/geoip.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/geoip.dat" \
  -source "geoip:geoip-asn=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/geoip-asn.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/geoip-asn.dat" \
  -source "geoip:geoip-ru-only=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/geoip-ru-only.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/geoip-ru-only.dat" \
  -source "geoip:ru-blocked=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/ru-blocked.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/ru-blocked.dat" \
  -source "geoip:ru-blocked-community=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/ru-blocked-community.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/ru-blocked-community.dat" \
  -source "geoip:re-filter=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/re-filter.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/re-filter.dat" \
  -source "geoip:private=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/private.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geoip@release/private.dat" \
  -source "geosite:geosite=https://raw.githubusercontent.com/runetfreedom/russia-blocked-geosite/release/geosite.dat,https://cdn.jsdelivr.net/gh/runetfreedom/russia-blocked-geosite@release/geosite.dat"

echo ""
echo "=== Conversion Complete ==="
echo "SRS files saved to: $OUTPUT_DIR"
ls -la "$OUTPUT_DIR"
