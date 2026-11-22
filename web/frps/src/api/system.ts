import { http } from './http'

export interface AlertThresholds {
  enabled: boolean
  maxOnlineClients: number
  maxCurConns: number
  checkIntervalSeconds: number
}

export interface AlertItem {
  time: string
  level: string
  metric: string
  value: number
  threshold: number
  message: string
}

export interface SystemConfig {
  bindAddr: string
  bindPort: number
  kcpBindPort: number
  quicBindPort: number
  vhostHTTPPort: number
  vhostHTTPSPort: number
  tcpmuxHTTPConnectPort: number
  subdomainHost: string
  maxPortsPerClient: number
  alerts: AlertThresholds
}

export const getSystemConfig = () => http.getV2<SystemConfig>('../api/v2/system/config')

export const getAlerts = () =>
  http.getV2<{ enabled: boolean; thresholds: AlertThresholds; items: AlertItem[] }>('../api/v2/alerts')

export const updateAlerts = (cfg: AlertThresholds) => http.put<void>('../api/v2/system/alerts', cfg)

export const clearAlerts = () => http.delete<void>('../api/v2/alerts')

export interface WebUser {
  user: string
  role: string
  legacy?: boolean
}

export const getWebUsers = () => http.getV2<{ items: WebUser[] }>('../api/v2/web/users')

export const updateWebUsers = (
  users: Array<{ user: string; role: string; password: string }>,
) => http.put<void>('../api/v2/web/users', { users })

export interface ReloadResult {
  applied: string[]
  restartRequired: string[]
}

export const reloadConfig = () => http.postV2<ReloadResult>('../api/v2/system/reload')
