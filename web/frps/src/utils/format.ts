function currentLang(): string {
  if (typeof document !== 'undefined' && document.documentElement.lang) {
    return document.documentElement.lang
  }
  return 'zh-CN'
}

export function formatDistanceToNow(date: Date): string {
  const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
  const rtf = new Intl.RelativeTimeFormat(currentLang(), { numeric: 'auto' })
  const units: Array<['year' | 'month' | 'day' | 'hour' | 'minute' | 'second', number]> = [
    ['year', 31536000],
    ['month', 2592000],
    ['day', 86400],
    ['hour', 3600],
    ['minute', 60],
    ['second', 1],
  ]
  for (const [unit, secs] of units) {
    if (seconds >= secs) {
      return rtf.format(-Math.floor(seconds / secs), unit)
    }
  }
  return rtf.format(0, 'second')
}

export function formatUnixSeconds(seconds?: number): string {
  if (seconds == null || !Number.isFinite(seconds) || seconds <= 0) return ''

  const date = new Date(seconds * 1000)
  const pad = (value: number) => value.toString().padStart(2, '0')
  return `${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

export function formatFileSize(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes < 0) return '0 B'
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const unit = sizes[i] || sizes[sizes.length - 1]
  const val = bytes / Math.pow(k, i)

  return parseFloat(val.toFixed(2)) + ' ' + unit
}