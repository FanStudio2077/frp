import { http, buildQueryString } from './http'

export interface AuditEntry {
  time: string
  user: string
  method: string
  path: string
  status: number
}

export interface AuditLogsResp {
  total: number
  items: AuditEntry[]
}

export const getAuditLogs = (params: { limit?: number; offset?: number; method?: string }) =>
  http.getV2<AuditLogsResp>('../api/v2/audit/logs' + buildQueryString(params))

export const clearAuditLogs = () => http.delete<void>('../api/v2/audit/logs')
