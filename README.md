# SRS Rules для sing-box

Автоматическая конвертация набора geoip/geosite `.dat` файлов в `.srs` для sing-box каждые 2 часа.

Файлы `.srs` публикуются в релизе `srs-latest` и дополнительно доступны как артефакты в разделе Actions (хранятся 1 день).

## Источники данных

- **geoip.dat family**: [runetfreedom/russia-blocked-geoip](https://github.com/runetfreedom/russia-blocked-geoip)
- **geosite.dat**: [runetfreedom/russia-blocked-geosite](https://github.com/runetfreedom/russia-blocked-geosite)
- **Конвертер**: локальная утилита в `tools/geodat2srs`, основанная на [runetfreedom/geodat2srs](https://github.com/runetfreedom/geodat2srs)

Доступные категории `geosite.dat`:
- Все категории из `@v2fly/domain-list-community`, включая `google`, `discord`, `youtube`, `twitter`, `meta`, `openai` и другие
- `geosite:ru-blocked` — заблокированные в России домены (`antifilter-download-community` + `re:filter`)
- `geosite:ru-blocked-all` — все известные заблокированные в России домены (`antifilter-download` + `antifilter-download-community` + `re:filter`)
- `geosite:ru-available-only-inside` — домены, доступные только внутри России
- `geosite:antifilter-download` — все домены из `antifilter.download`
- `geosite:antifilter-download-community` — все домены из `community.antifilter.download`
- `geosite:refilter` — все домены из `re:filter`
- `geosite:category-ads-all` — все рекламные домены
- `geosite:win-spy` — домены, используемые Windows для слежки и сбора аналитики
- `geosite:win-update` — домены, используемые Windows для обновлений
- `geosite:win-extra` — прочие домены, используемые Windows

Ссылки на последние `.srs` из релиза:
- `geoip.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/geoip.srs`
- `geoip-asn.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/geoip-asn.srs`
- `geoip-ru-only.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/geoip-ru-only.srs`
- `ru-blocked.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/ru-blocked.srs`
- `ru-blocked-community.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/ru-blocked-community.srs`
- `re-filter.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/re-filter.srs`
- `private.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/private.srs`
- `geosite.srs`: `https://github.com/Arhimage/sb_srs/releases/latest/download/geosite.srs`

## Автоматическая сборка

GitHub Actions workflow запускается каждые 2 часа и:
1. Запускает `scripts/convert.sh`
2. Передаёт список источников аргументами в локальную утилиту `tools/geodat2srs`
3. Собирает набор `.srs` файлов с именами источников
4. Публикует результат в release `srs-latest`
5. Сохраняет артефакты (1 день)

## Использование в sing-box

Файлы из директории `rules/` или из release можно использовать в конфигурации sing-box:

```json
{
  "route": {
    "rules": [
      {
        "rule_set": "ru-blocked-ips",
        "outbound": "block"
      },
      {
        "rule_set": "ru-blocked-sites",
        "outbound": "block"
      }
    ],
    "rule_set": [
      {
        "tag": "ru-blocked-ips",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/Arhimage/sb_srs/releases/latest/download/ru-blocked.srs"
      },
      {
        "tag": "ru-blocked-sites",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/Arhimage/sb_srs/releases/latest/download/geosite.srs"
      }
    ]
  }
}
```
