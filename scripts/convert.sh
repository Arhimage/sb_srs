#!/bin/bash

# Скрипт для конвертации geoip.dat и geosite.dat в формат SRS
# Использует утилиту geodat2srs из https://github.com/runetfreedom/geodat2srs

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
TEMP_DIR="/tmp/geodat2srs_build"
OUTPUT_DIR="$PROJECT_DIR/rules"

# URLs для скачивания
GEOIP_URL="https://raw.githubusercontent.com/runetfreedom/russia-blocked-geoip/release/geoip.dat"
GEOSITE_URL="https://raw.githubusercontent.com/runetfreedom/russia-blocked-geosite/release/geosite.dat"
GEODAT2SRS_REPO="https://github.com/runetfreedom/geodat2srs.git"

echo "=== Build SRS Files ==="

# Очистка временной директории
rm -rf "$TEMP_DIR"
mkdir -p "$TEMP_DIR" "$OUTPUT_DIR"

# Скачивание .dat файлов
echo "[1/5] Downloading geoip.dat..."
curl -L -o "$TEMP_DIR/geoip.dat" "$GEOIP_URL"

echo "[2/5] Downloading geosite.dat..."
curl -L -o "$TEMP_DIR/geosite.dat" "$GEOSITE_URL"

# Клонирование и компиляция geodat2srs
echo "[3/5] Cloning geodat2srs..."
git clone "$GEODAT2SRS_REPO" "$TEMP_DIR/geodat2srs"

echo "[4/5] Building geodat2srs..."
cd "$TEMP_DIR/geodat2srs"
go build -o "$TEMP_DIR/geodat2srs" .

# Конвертация
echo "[5/5] Converting to SRS format..."
"$TEMP_DIR/geodat2srs" geoip -i "$TEMP_DIR/geoip.dat" -o "$OUTPUT_DIR" --prefix "geoip-"
"$TEMP_DIR/geodat2srs" geosite -i "$TEMP_DIR/geosite.dat" -o "$OUTPUT_DIR" --prefix "geosite-"

# Очистка
rm -rf "$TEMP_DIR"

echo ""
echo "=== Conversion Complete ==="
echo "SRS files saved to: $OUTPUT_DIR"
ls -la "$OUTPUT_DIR"
