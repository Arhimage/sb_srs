# SRS Rules для sing-box

Автоматическая конвертация баз IP-адресов и доменов (geoip/geosite) в формат `.srs` для sing-box каждые 2 часа.

Файлы `.srs` доступны как артефакты в разделе Actions (хранятся 1 день).

## Источники данных

- **geoip.dat**: [runetfreedom/russia-blocked-geoip](https://github.com/runetfreedom/russia-blocked-geoip)
- **geosite.dat**: [runetfreedom/russia-blocked-geosite](https://github.com/runetfreedom/russia-blocked-geosite)
- **Конвертер**: [runetfreedom/geodat2srs](https://github.com/runetfreedom/geodat2srs)

## Автоматическая сборка

GitHub Actions workflow запускается каждые 2 часа и:
1. Скачивает актуальные версии `geoip.dat` и `geosite.dat`
2. Компилирует утилиту `geodat2srs`
3. Конвертирует файлы в формат `.srs`
4. Сохраняет артефакты (1 день)

## Локальный запуск

Для локальной конвертации (требуется Go):

```bash
chmod +x scripts/convert.sh
./scripts/convert.sh
```

## Использование в sing-box

Файлы из директории `rules/` можно использовать в конфигурации sing-box:

```json
{
  "route": {
    "rules": [
      {
        "rule_set": "geoip-ru",
        "outbound": "block"
      },
      {
        "rule_set": "geosite-category-ads",
        "outbound": "block"
      }
    ],
    "rule_set": [
      {
        "tag": "geoip-ru",
        "type": "local",
        "path": "rules/geoip-ru.srs"
      },
      {
        "tag": "geosite-category-ads",
        "type": "local",
        "path": "rules/geosite-category-ads.srs"
      }
    ]
  }
}
```
