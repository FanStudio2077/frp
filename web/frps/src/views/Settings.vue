<template>
  <div class="settings-page">
    <div class="page-header">
      <h2 class="page-title">{{ t('server.settings.title') }}</h2>
      <ActionButton size="small" :loading="reloading" @click="doReload">{{ t('server.settings.reload') }}</ActionButton>
    </div>
    <el-alert v-if="reloadResult" type="success" :title="t('server.settings.reloaded')" show-icon :closable="true">
      <div>{{ t('server.settings.applied') }}: {{ reloadResult.applied.join(', ') }}</div>
      <div v-if="reloadResult.restartRequired.length">{{ t('server.settings.restartRequired') }}: {{ reloadResult.restartRequired.join(', ') }}</div>
    </el-alert>

    <el-card shadow="hover" class="settings-card">
      <template #header><span class="card-title">{{ t('server.settings.configTitle') }}</span></template>
      <div class="config-grid">
        <div class="config-item" v-for="item in configItems" :key="item.label">
          <span class="config-label">{{ item.label }}</span>
          <span class="config-value">{{ item.value }}</span>
        </div>
      </div>
      <p class="hint">{{ t('server.settings.restartHint') }}</p>
    </el-card>

    <el-card shadow="hover" class="settings-card">
      <template #header><span class="card-title">{{ t('server.settings.alertsTitle') }}</span></template>
      <el-form label-position="top" class="alert-form">
        <el-form-item :label="t('server.settings.enabled')">
          <el-switch v-model="form.enabled" />
        </el-form-item>
        <el-form-item :label="t('server.settings.maxOnlineClients')">
          <el-input-number v-model="form.maxOnlineClients" :min="0" />
        </el-form-item>
        <el-form-item :label="t('server.settings.maxCurConns')">
          <el-input-number v-model="form.maxCurConns" :min="0" />
        </el-form-item>
        <el-form-item :label="t('server.settings.interval')">
          <el-input-number v-model="form.checkIntervalSeconds" :min="1" />
        </el-form-item>
      </el-form>
      <ActionButton :loading="saving" @click="save">{{ t('server.settings.save') }}</ActionButton>
    </el-card>

    <el-card shadow="hover" class="settings-card">
      <template #header>
        <div class="card-header-row">
          <span class="card-title">{{ t('server.settings.currentAlerts') }}</span>
          <ActionButton variant="outline" size="small" danger @click="doClear">{{ t('server.settings.clear') }}</ActionButton>
        </div>
      </template>
      <el-table :data="alerts" size="small">
        <el-table-column :label="t('server.settings.time')" width="180">
          <template #default="{ row }">{{ formatTime(row.time) }}</template>
        </el-table-column>
        <el-table-column prop="metric" :label="t('server.settings.metric')" width="150" />
        <el-table-column prop="value" :label="t('server.settings.value')" width="100" />
        <el-table-column prop="threshold" :label="t('server.settings.threshold')" width="100" />
        <template #empty>{{ t('server.settings.noAlerts') }}</template>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import ActionButton from '@shared/components/ActionButton.vue'
import {
  getSystemConfig,
  getAlerts,
  updateAlerts,
  clearAlerts,
  reloadConfig,
  type ReloadResult,
  type SystemConfig,
  type AlertItem,
  type AlertThresholds,
} from '../api/system'

const { t, locale } = useI18n()
const config = ref<SystemConfig | null>(null)
const alerts = ref<AlertItem[]>([])
const saving = ref(false)
const reloading = ref(false)
const reloadResult = ref<ReloadResult | null>(null)
const form = reactive<AlertThresholds>({ enabled: false, maxOnlineClients: 0, maxCurConns: 0, checkIntervalSeconds: 10 })

const configItems = computed(() => {
  const c = config.value
  if (!c) return []
  return [
    { label: 'bindAddr', value: c.bindAddr },
    { label: 'bindPort', value: c.bindPort },
    { label: 'kcpBindPort', value: c.kcpBindPort },
    { label: 'quicBindPort', value: c.quicBindPort },
    { label: 'vhostHTTPPort', value: c.vhostHTTPPort },
    { label: 'vhostHTTPSPort', value: c.vhostHTTPSPort },
    { label: 'tcpmuxHTTPConnectPort', value: c.tcpmuxHTTPConnectPort },
    { label: 'subdomainHost', value: c.subdomainHost },
    { label: 'maxPortsPerClient', value: c.maxPortsPerClient },
  ]
})

const formatTime = (s: string) => {
  const d = new Date(s)
  return isNaN(d.getTime()) ? s : d.toLocaleString(locale.value)
}

const fetchData = async () => {
  try {
    const [cfg, al] = await Promise.all([getSystemConfig(), getAlerts()])
    config.value = cfg
    Object.assign(form, cfg.alerts)
    alerts.value = al.items || []
  } catch (e: any) {
    ElMessage.error(t('server.settings.fetchFailed') + (e?.message || ''))
  }
}

const save = async () => {
  saving.value = true
  try {
    await updateAlerts({ ...form })
    ElMessage.success(t('server.settings.saved'))
  } catch (e: any) {
    ElMessage.error(t('server.settings.saveFailed') + (e?.message || ''))
  } finally {
    saving.value = false
  }
}

const doReload = async () => {
  reloading.value = true
  try {
    reloadResult.value = await reloadConfig()
    ElMessage.success(t('server.settings.reloaded'))
  } catch (e: any) {
    ElMessage.error(t('server.settings.reloadFailed') + (e?.message || ''))
  } finally {
    reloading.value = false
  }
}

const doClear = async () => {
  try {
    await clearAlerts()
    ElMessage.success(t('server.settings.cleared'))
    alerts.value = []
  } catch (e: any) {
    ElMessage.error(t('server.settings.saveFailed') + (e?.message || ''))
  }
}

fetchData()
</script>

<style scoped>
.settings-page { display: flex; flex-direction: column; gap: 16px; }
.page-header { display: flex; align-items: center; justify-content: space-between; }
.page-title { font-size: 20px; font-weight: 600; margin: 0; }
.card-title { font-size: 16px; font-weight: 500; }
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.config-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 12px; }
.config-item { display: flex; justify-content: space-between; border-bottom: 1px dashed var(--el-border-color-lighter); padding: 6px 0; }
.config-label { color: var(--el-text-color-secondary); }
.config-value { font-weight: 500; }
.hint { color: var(--el-text-color-secondary); font-size: 12px; margin-top: 12px; }
.alert-form { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; }
</style>
