#!/bin/bash

# Скрипт для конвертации geoip.dat и geosite.dat в два итоговых SRS файла
# Использует локальную утилиту на основе https://github.com/runetfreedom/geodat2srs

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
TEMP_DIR="/tmp/geodat2srs_build"
OUTPUT_DIR="$PROJECT_DIR/rules"
CONVERTER_DIR="$PROJECT_DIR/tools/geodat2srs"

# URLs для скачивания
GEOIP_URL="https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/geoip.dat"
GEOSITE_URL="https://raw.githubusercontent.com/runetfreedom/russia-blocked-geosite/release/geosite.dat"

echo "=== Build SRS Files ==="

# Очистка временной директории
rm -rf "$TEMP_DIR"
mkdir -p "$TEMP_DIR" "$OUTPUT_DIR"

# Скачивание .dat файлов
echo "[1/5] Downloading geoip.dat..."
curl -L -o "$TEMP_DIR/geoip.dat" "$GEOIP_URL"

echo "[2/5] Downloading geosite.dat..."
curl -L -o "$TEMP_DIR/geosite.dat" "$GEOSITE_URL"

# Очистка старых результатов и конвертация
echo "[3/5] Cleaning old SRS files..."
rm -f "$OUTPUT_DIR"/*.srs

echo "[4/5] Running bundled geodat2srs..."
cd "$CONVERTER_DIR"

echo "[5/5] Converting to SRS format..."
go run . -geoip "$TEMP_DIR/geoip.dat" -geosite "$TEMP_DIR/geosite.dat" -output-dir "$OUTPUT_DIR"

# Очистка
rm -rf "$TEMP_DIR"

echo ""
echo "=== Conversion Complete ==="
echo "SRS files saved to: $OUTPUT_DIR"
ls -la "$OUTPUT_DIR"
