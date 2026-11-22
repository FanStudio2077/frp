import { describe, it, expect } from 'vitest'
import zhCN from '../src/i18n/zh-CN'

describe('i18n scaffold', () => {
  it('has navigation labels', () => {
    expect(zhCN.nav.overview).toBeTruthy()
    expect(zhCN.nav.servers).toBe('服务器')
  })
})
