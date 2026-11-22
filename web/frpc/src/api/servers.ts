import { http } from './http'

export interface ServerEntry {
  name: string
  serverAddr: string
  serverPort: number
  enabled: boolean
  priority: number
  active?: boolean
}

export interface ServersResponse {
  active: string
  servers: ServerEntry[]
}

export const getServers = () => http.get<ServersResponse>('/api/servers')

export const putServers = (servers: Array<Omit<ServerEntry, 'active'>>) =>
  http.put<void>('/api/servers', { servers })

export const activateServer = (name: string) =>
  http.post<void>('/api/servers/' + encodeURIComponent(name) + '/activate')
