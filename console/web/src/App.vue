<template>
  <div class="layout">
    <header class="header">
      <strong>{{ t('brand') }}</strong>
      <el-select v-model="locale" size="small" style="width: 130px" @change="onLocale">
        <el-option label="简体中文" value="zh-CN" />
        <el-option label="English" value="en" />
      </el-select>
    </header>
    <div class="body">
      <aside class="sidebar">
        <router-link v-for="item in navs" :key="item.path" :to="item.path" class="nav-link">{{ t(item.key) }}</router-link>
      </aside>
      <main class="content"><router-view /></main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const navs = [
  { path: '/overview', key: 'nav.overview' },
  { path: '/servers', key: 'nav.servers' },
  { path: '/nodes', key: 'nav.nodes' },
  { path: '/proxies', key: 'nav.proxies' },
  { path: '/visitors', key: 'nav.visitors' },
  { path: '/users', key: 'nav.users' },
  { path: '/audit', key: 'nav.audit' },
  { path: '/alerts', key: 'nav.alerts' },
  { path: '/settings', key: 'nav.settings' },
]
const onLocale = (v: string) => {
  locale.value = v
  if (typeof localStorage !== 'undefined') localStorage.setItem('console-locale', v)
}
</script>

<style>
body { margin: 0; font-family: system-ui, -apple-system, Segoe UI, Arial, sans-serif; }
.layout { height: 100vh; display: flex; flex-direction: column; }
.header { height: 52px; display: flex; align-items: center; justify-content: space-between; padding: 0 16px; border-bottom: 1px solid #e4e7ed; }
.body { flex: 1; display: flex; min-height: 0; }
.sidebar { width: 180px; border-right: 1px solid #e4e7ed; padding: 12px 8px; display: flex; flex-direction: column; gap: 4px; }
.nav-link { display: block; padding: 8px 12px; border-radius: 6px; color: #303133; text-decoration: none; }
.nav-link.router-link-active { background: #eef2ff; color: #3b82f6; }
.content { flex: 1; padding: 24px; overflow: auto; }
</style>
