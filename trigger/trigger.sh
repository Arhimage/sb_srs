#!/bin/sh

CRON="${CRON_SCHEDULE:-0 0 * * *}"
GITHUB_API="https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/actions/workflows/build-srs.yml/dispatches"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

trigger() {
  log "Triggering workflow..."
  RESP=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    -d '{"ref":"main"}' \
    "$GITHUB_API")

  if [ "$RESP" = "204" ]; then
    log "OK — workflow dispatched."
  else
    log "ERROR — HTTP $RESP"
  fi
}

trigger

log "Cron schedule: $CRON"
echo "$CRON /trigger.sh >> /var/log/cron.log 2>&1" | crontab -
crond -f
