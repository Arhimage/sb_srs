# SRS Rules for sing-box

Автоматическая конвертация набора `geoip/geosite` `.dat` файлов в `.srs` для `sing-box` каждые 2 часа.

Сгенерированные `.srs` публикуются в релизе `srs-latest` и дополнительно доступны как GitHub Actions artifacts на 1 день.

## Sources

- `geoip.dat` family: [runetfreedom/russia-blocked-geoip](https://github.com/runetfreedom/russia-blocked-geoip)
- `geosite.dat`: [runetfreedom/russia-blocked-geosite](https://github.com/runetfreedom/russia-blocked-geosite)
- локальный конвертер: `tools/geodat2srs`

## Release Links

### IP Rules

- `geoip`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geoip.srs`
- `geoip-asn`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geoip-asn.srs`
- `geoip-ru-only`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geoip-ru-only.srs`
- `ru-blocked`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/ru-blocked.srs`
- `ru-blocked-community`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/ru-blocked-community.srs`
- `re-filter`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/re-filter.srs`
- `private`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/private.srs`

### Site Rules

- `geosite-ru-blocked`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-ru-blocked.srs`
- `geosite-ru-blocked-all`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-ru-blocked-all.srs`
- `geosite-category-ads-all`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-category-ads-all.srs`
- `geosite-openai`
  `https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-openai.srs`

Остальные site-категории публикуются в том же формате:
`https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-<category>.srs`

## Build Flow

1. Запускается `scripts/convert.sh`.
2. Список источников передаётся аргументами в `tools/geodat2srs`.
3. Для `geoip` создаётся один `.srs` на источник.
4. `geosite.dat` делится на отдельные `geosite-<category>.srs`.
5. Результат публикуется в release `srs-latest`.

## sing-box Example

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
      },
      {
        "rule_set": "ads-sites",
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
        "url": "https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-ru-blocked.srs"
      },
      {
        "tag": "ads-sites",
        "type": "remote",
        "format": "binary",
        "url": "https://github.com/Arhimage/sb_srs/releases/latest/download/geosite-category-ads-all.srs"
      }
    ]
  }
}
```
