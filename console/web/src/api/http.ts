export async function apiGet<T>(url: string): Promise<T> {
  const res = await fetch(url, { credentials: 'include' })
  if (!res.ok) throw new Error('HTTP ' + res.status)
  return res.json() as Promise<T>
}

export async function apiSend<T>(method: string, url: string, body?: unknown): Promise<T> {
  const res = await fetch(url, {
    method,
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  if (!res.ok) throw new Error('HTTP ' + res.status)
  return (res.status === 204 ? ({} as T) : res.json()) as Promise<T>
}
