## Custom Build v0.73.0

* 服务端/客户端 Web 面板中英双语 i18n（默认中文，可切换，可扩展）。
* frpc 多服务器配置与故障转移（面板 Servers 管理页）。
* frps 面板增强：操作审计日志、客户端强制下线、代理批量删除、实时监控告警、告警阈值热更新。
* frps 通用配置热重载（API/SIGHUP，热应用鉴权/告警等；监听端口需重启）。
* 配置驱动的多用户 RBAC（admin/viewer）与面板用户管理页。

---

## Features

* frpc can verify server certificates using system CAs, optionally combined with CAs from `trustedCaFile`.
* `frpc reload` now refreshes certificate files for unchanged `https2http` and `https2https` proxies without interrupting existing connections. Failed refreshes keep the previous certificate and report the error in proxy status and logs while reload still succeeds.
* Added `metricsServer` to both frps and frpc, a separate listener for health checks, metrics and optional profiling that does not expose the dashboard or the admin API. It is disabled unless `metricsServer.port` is set, defaults to `127.0.0.1`, and has its own credentials and TLS. Setting the port serves `/metrics` on both: frps exposes runtime and business metrics, so `enablePrometheus` is not needed there, and frpc exposes Go runtime metrics only, without linking an extra metrics library.

  The `/healthz`, `/metrics` and pprof endpoints on `webServer`, `webServer.pprofEnable` and `enablePrometheus` still work but are deprecated: configure `metricsServer` instead. They will continue to work until v0.80.0 is released, and may be removed in v0.80.0 or later.

## Fixes

* Fixed proxy and visitor configuration changes being lost when reloading frpc during reconnection.
* Updated yamux to fix a 100% CPU busy loop after a session shuts down and a session close that could hang when the underlying connection (for example WebSocket) blocks on a write, which could also delay reconnection.
* frps now closes the multiplexed session of a replaced or failed control connection immediately, releasing its streams and work connections instead of keeping them until keepalive times out.
* Fixed a data race when frpc updates log prefixes while other goroutines are logging, which could crash the process.
* Fixed HTTPS and tcpmux route listeners leaving a pending connection open when the listener is removed during connection handoff; the connection is now closed instead of leaving the client to wait for a read timeout.
