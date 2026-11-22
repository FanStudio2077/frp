<template>
  <div class="audit-page">
    <div class="page-header">
      <h2 class="page-title">{{ t('server.audit.title') }}</h2>
      <div class="header-actions">
        <el-select v-model="methodFilter" class="method-select" @change="onFilterChange">
          <el-option :label="t('server.audit.all')" value="" />
          <el-option label="POST" value="POST" />
          <el-option label="PUT" value="PUT" />
          <el-option label="DELETE" value="DELETE" />
          <el-option label="PATCH" value="PATCH" />
        </el-select>
        <ActionButton variant="outline" size="small" @click="fetchData">{{ t('server.audit.refresh') }}</ActionButton>
        <ActionButton variant="outline" size="small" danger @click="clearVisible = true">{{ t('server.audit.clear') }}</ActionButton>
      </div>
    </div>

    <el-table v-loading="loading" :data="entries" class="audit-table">
      <el-table-column :label="t('server.audit.time')" width="200">
        <template #default="{ row }">{{ formatTime(row.time) }}</template>
      </el-table-column>
      <el-table-column prop="user" :label="t('server.audit.user')" width="140" />
      <el-table-column prop="method" :label="t('server.audit.method')" width="110" />
      <el-table-column prop="path" :label="t('server.audit.path')" />
      <el-table-column prop="status" :label="t('server.audit.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status < 400 ? 'success' : 'danger'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <template #empty>{{ t('server.audit.empty') }}</template>
    </el-table>

    <div v-if="total > 0" class="pagination-section">
      <ElPagination
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @current-change="onPageChange"
        @size-change="onPageSizeChange"
      />
    </div>

    <ConfirmDialog
      v-model="clearVisible"
      :title="t('server.audit.clear')"
      :message="t('server.audit.clearConfirm')"
      :confirm-text="t('server.audit.clear')"
      danger
      @confirm="doClear"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElPagination } from 'element-plus'
import ActionButton from '@shared/components/ActionButton.vue'
import ConfirmDialog from '@shared/components/ConfirmDialog.vue'
import { getAuditLogs, clearAuditLogs, type AuditEntry } from '../api/audit'

const { t, locale } = useI18n()
const entries = ref<AuditEntry[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const methodFilter = ref('')
const clearVisible = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const data = await getAuditLogs({ limit: pageSize.value, offset: (page.value - 1) * pageSize.value, method: methodFilter.value })
    entries.value = data.items || []
    total.value = data.total || 0
  } catch (e: any) {
    ElMessage.error(t('server.audit.fetchFailed') + (e?.message || ''))
  } finally {
    loading.value = false
  }
}

const formatTime = (s: string) => {
  const d = new Date(s)
  return isNaN(d.getTime()) ? s : d.toLocaleString(locale.value)
}

const onFilterChange = () => {
  page.value = 1
  fetchData()
}
const onPageChange = (v: number) => {
  page.value = v
  fetchData()
}
const onPageSizeChange = (v: number) => {
  pageSize.value = v
  page.value = 1
  fetchData()
}
const doClear = async () => {
  try {
    await clearAuditLogs()
    ElMessage.success(t('server.audit.cleared'))
    clearVisible.value = false
    page.value = 1
    await fetchData()
  } catch (e: any) {
    ElMessage.error(t('server.audit.fetchFailed') + (e?.message || ''))
  }
}

fetchData()
</script>

<style scoped>
.audit-page { display: flex; flex-direction: column; gap: 16px; }
.page-header { display: flex; align-items: center; justify-content: space-between; }
.page-title { font-size: 20px; font-weight: 600; margin: 0; }
.header-actions { display: flex; align-items: center; gap: 8px; }
.method-select { width: 140px; }
</style>
