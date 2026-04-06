# SRS Rules для sing-box

Автоматическая конвертация баз IP-адресов и доменов (geoip/geosite) в два файла `.srs` для sing-box каждые 2 часа.

Файлы `.srs` публикуются в релизе `srs-latest` и дополнительно доступны как артефакты в разделе Actions (хранятся 1 день).

## Источники данных

- **geoip.dat**: [runetfreedom/russia-blocked-geoip](https://github.com/runetfreedom/russia-blocked-geoip)
- **geosite.dat**: [runetfreedom/russia-blocked-geosite](https://github.com/runetfreedom/russia-blocked-geosite)
- **Конвертер**: локальная утилита в `tools/geodat2srs`, основанная на [runetfreedom/geodat2srs](https://github.com/runetfreedom/geodat2srs)

## Автоматическая сборка

GitHub Actions workflow запускается каждые 2 часа и:
1. Скачивает актуальные версии `geoip.dat` и `geosite.dat`
2. Запускает локальную утилиту `tools/geodat2srs`
3. Собирает ровно два файла: `geoip.srs` и `geosite.srs`
4. Публикует результат в release `srs-latest`
5. Сохраняет артефакты (1 день)

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
        "rule_set": "geoip-all",
        "outbound": "block"
      },
      {
        "rule_set": "geosite-all",
        "outbound": "block"
      }
    ],
    "rule_set": [
      {
        "tag": "geoip-all",
        "type": "local",
        "path": "rules/geoip.srs"
      },
      {
        "tag": "geosite-all",
        "type": "local",
        "path": "rules/geosite.srs"
      }
    ]
  }
}
```
