<template>
  <div class="servers-page">
    <div class="page-header">
      <h2 class="page-title">{{ t('client.servers.title') }}</h2>
      <div class="header-actions">
        <ActionButton variant="outline" size="small" @click="fetchData">{{ t('client.servers.refresh') }}</ActionButton>
        <ActionButton size="small" @click="openCreate">{{ t('client.servers.add') }}</ActionButton>
      </div>
    </div>

    <el-table v-loading="loading" :data="servers" class="servers-table">
      <el-table-column prop="name" :label="t('client.servers.name')" />
      <el-table-column :label="t('client.servers.address')">
        <template #default="{ row }">{{ row.serverAddr }}:{{ row.serverPort }}</template>
      </el-table-column>
      <el-table-column prop="priority" :label="t('client.servers.priority')" width="90" />
      <el-table-column :label="t('client.servers.enabled')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? t('client.servers.on') : t('client.servers.off') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('client.servers.status')" width="130">
        <template #default="{ row }">
          <el-tag v-if="row.active" type="success">{{ t('client.servers.active') }}</el-tag>
          <el-button v-else link type="primary" @click="activate(row.name)">{{ t('client.servers.activate') }}</el-button>
        </template>
      </el-table-column>
      <el-table-column :label="t('client.servers.actions')" width="160">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">{{ t('common.edit') }}</el-button>
          <el-button link type="danger" @click="removeRow(row)">{{ t('common.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingIndex >= 0 ? t('client.servers.editTitle') : t('client.servers.addTitle')" width="480px">
      <el-form :model="form" label-position="top">
        <el-form-item :label="t('client.servers.name')">
          <el-input v-model="form.name" :placeholder="t('client.servers.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('client.servers.address')">
          <el-input v-model="form.serverAddr" placeholder="1.2.3.4" />
        </el-form-item>
        <el-form-item :label="t('client.servers.port')">
          <el-input-number v-model="form.serverPort" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item :label="t('client.servers.priority')">
          <el-input-number v-model="form.priority" :min="0" />
        </el-form-item>
        <el-form-item :label="t('client.servers.enabled')">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <ActionButton variant="outline" @click="dialogVisible = false">{{ t('common.cancel') }}</ActionButton>
        <ActionButton :loading="saving" @click="save">{{ t('client.servers.save') }}</ActionButton>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import ActionButton from '@shared/components/ActionButton.vue'
import { getServers, putServers, activateServer, type ServerEntry } from '../api/servers'

const { t } = useI18n()
const servers = ref<ServerEntry[]>([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const editingIndex = ref(-1)
const form = reactive({ name: '', serverAddr: '', serverPort: 7000, priority: 0, enabled: true })

const toPayload = (s: ServerEntry) => ({ name: s.name, serverAddr: s.serverAddr, serverPort: s.serverPort, priority: s.priority, enabled: s.enabled })

const fetchData = async () => {
  loading.value = true
  try {
    const data = await getServers()
    servers.value = data.servers || []
  } catch (e: any) {
    ElMessage.error(t('client.servers.fetchFailed') + (e?.message || ''))
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editingIndex.value = -1
  Object.assign(form, { name: '', serverAddr: '', serverPort: 7000, priority: servers.value.length, enabled: true })
  dialogVisible.value = true
}

const openEdit = (row: ServerEntry) => {
  editingIndex.value = servers.value.findIndex((s) => s.name === row.name)
  Object.assign(form, { name: row.name, serverAddr: row.serverAddr, serverPort: row.serverPort, priority: row.priority, enabled: row.enabled })
  dialogVisible.value = true
}

const save = async () => {
  saving.value = true
  try {
    const list = servers.value.map(toPayload)
    const entry = { name: form.name, serverAddr: form.serverAddr, serverPort: form.serverPort, priority: form.priority, enabled: form.enabled }
    if (editingIndex.value >= 0) list[editingIndex.value] = entry
    else list.push(entry)
    await putServers(list)
    ElMessage.success(t('client.servers.saved'))
    dialogVisible.value = false
    await fetchData()
  } catch (e: any) {
    ElMessage.error(t('client.servers.saveFailed') + (e?.message || ''))
  } finally {
    saving.value = false
  }
}

const removeRow = async (row: ServerEntry) => {
  try {
    await ElMessageBox.confirm(t('client.servers.deleteConfirm', { name: row.name }), t('client.servers.deleteTitle'), { type: 'warning' })
  } catch {
    return
  }
  try {
    await putServers(servers.value.filter((s) => s.name !== row.name).map(toPayload))
    ElMessage.success(t('client.servers.deleted'))
    await fetchData()
  } catch (e: any) {
    ElMessage.error(t('client.servers.saveFailed') + (e?.message || ''))
  }
}

const activate = async (name: string) => {
  try {
    await activateServer(name)
    ElMessage.success(t('client.servers.activated'))
    await fetchData()
  } catch (e: any) {
    ElMessage.error(t('client.servers.activateFailed') + (e?.message || ''))
  }
}

fetchData()
</script>

<style scoped>
.servers-page { display: flex; flex-direction: column; gap: 16px; }
.page-header { display: flex; align-items: center; justify-content: space-between; }
.page-title { font-size: 20px; font-weight: 600; margin: 0; }
.header-actions { display: flex; gap: 8px; }
</style>
