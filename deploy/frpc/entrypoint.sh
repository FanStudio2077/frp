#!/bin/sh
set -e

CFG="${FRP_CONFIG_PATH:-/etc/frp/frpc.toml}"

if [ -n "$FRP_CONFIG_B64" ]; then
  # Full config provided base64-encoded
  printf '%s' "$FRP_CONFIG_B64" | base64 -d > "$CFG"
elif [ -n "$FRP_SERVER_ADDR" ]; then
  mkdir -p "$(dirname "$CFG")"
  {
    printf 'serverAddr = "%s"\n' "$FRP_SERVER_ADDR"
    printf 'serverPort = %s\n' "${FRP_SERVER_PORT:-7000}"
    if [ -n "$FRP_TOKEN" ]; then
      printf 'auth.method = "token"\n'
      printf 'auth.token = "%s"\n' "$FRP_TOKEN"
    fi
    printf 'webServer.addr = "0.0.0.0"\n'
    printf 'webServer.port = %s\n' "${FRP_WEB_PORT:-7400}"
    if [ -n "$FRP_WEB_USER" ]; then printf 'webServer.user = "%s"\n' "$FRP_WEB_USER"; fi
    if [ -n "$FRP_WEB_PASS" ]; then printf 'webServer.password = "%s"\n' "$FRP_WEB_PASS"; fi
    if [ -n "$FRP_SERVERS" ]; then printf '%b\n' "$FRP_SERVERS"; fi
    if [ -n "$FRP_PROXIES" ]; then printf '%b\n' "$FRP_PROXIES"; fi
  } > "$CFG"
fi

if [ ! -f "$CFG" ]; then
  echo "ERROR: no frpc config. Set FRP_SERVER_ADDR/FRP_TOKEN, or mount $CFG." >&2
  exit 1
fi

echo "----- generated frpc.toml -----"
cat "$CFG"
echo "-------------------------------"
exec /usr/bin/frpc -c "$CFG"
