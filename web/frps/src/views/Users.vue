<template>
  <div class="users-page">
    <div class="page-header">
      <h2 class="page-title">{{ t('server.users.title') }}</h2>
      <ActionButton size="small" @click="openCreate">{{ t('server.users.add') }}</ActionButton>
    </div>
    <p class="hint">{{ t('server.users.usersHint') }}</p>

    <el-table :data="rows" class="users-table">
      <el-table-column prop="user" :label="t('server.users.user')" />
      <el-table-column :label="t('server.users.role')" width="140">
        <template #default="{ row }"><el-tag>{{ roleLabel(row.role) }}</el-tag></template>
      </el-table-column>
      <el-table-column :label="t('server.users.actions')" width="180">
        <template #default="{ row }">
          <template v-if="!row.legacy">
            <el-button link type="primary" @click="openEdit(row)">{{ t('common.edit') }}</el-button>
            <el-button link type="danger" @click="removeRow(row)">{{ t('common.delete') }}</el-button>
          </template>
          <el-tag v-else type="info">{{ t('server.users.legacy') }}</el-tag>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editing >= 0 ? t('server.users.editTitle') : t('server.users.addTitle')" width="460px">
      <el-form label-position="top">
        <el-form-item :label="t('server.users.user')">
          <el-input v-model="form.user" :disabled="editing >= 0 && rows[editing]?.legacy" />
        </el-form-item>
        <el-form-item :label="t('server.users.password')">
          <el-input v-model="form.password" type="password" show-password :placeholder="t('server.users.passwordPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('server.users.role')">
          <el-select v-model="form.role">
            <el-option :label="t('server.users.admin')" value="admin" />
            <el-option :label="t('server.users.viewer')" value="viewer" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <ActionButton variant="outline" @click="dialogVisible = false">{{ t('common.cancel') }}</ActionButton>
        <ActionButton :loading="saving" @click="save">{{ t('server.users.save') }}</ActionButton>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import ActionButton from '@shared/components/ActionButton.vue'
import { getWebUsers, updateWebUsers, type WebUser } from '../api/system'

const { t } = useI18n()
const rows = ref<WebUser[]>([])
const pwMap = reactive<Record<string, string>>({})
const dialogVisible = ref(false)
const saving = ref(false)
const editing = ref(-1)
const form = reactive({ user: '', password: '', role: 'viewer' })

const roleLabel = (role: string) => (role === 'admin' ? t('server.users.admin') : t('server.users.viewer'))

const fetchData = async () => {
  try {
    const data = await getWebUsers()
    rows.value = data.items || []
  } catch (e: any) {
    ElMessage.error(t('server.users.fetchFailed') + (e?.message || ''))
  }
}

const buildPayload = () =>
  rows.value
    .filter((r) => !r.legacy)
    .map((r) => ({ user: r.user, role: r.role, password: pwMap[r.user] || '' }))

const openCreate = () => {
  editing.value = -1
  Object.assign(form, { user: '', password: '', role: 'viewer' })
  dialogVisible.value = true
}
const openEdit = (row: WebUser) => {
  editing.value = rows.value.findIndex((r) => r.user === row.user)
  Object.assign(form, { user: row.user, password: '', role: row.role })
  dialogVisible.value = true
}

const save = async () => {
  if (!form.user.trim()) {
    ElMessage.warning(t('server.users.userRequired'))
    return
  }
  saving.value = true
  try {
    if (editing.value < 0) {
      rows.value.push({ user: form.user, role: form.role })
      pwMap[form.user] = form.password
    } else {
      const old = rows.value[editing.value]
      if (form.password) pwMap[form.user] = form.password
      if (form.user !== old.user) {
        if (!pwMap[form.user] && pwMap[old.user]) pwMap[form.user] = pwMap[old.user]
        delete pwMap[old.user]
      }
      rows.value[editing.value] = { user: form.user, role: form.role, legacy: old.legacy }
    }
    await updateWebUsers(buildPayload())
    ElMessage.success(t('server.users.saved'))
    dialogVisible.value = false
    await fetchData()
  } catch (e: any) {
    ElMessage.error(t('server.users.saveFailed') + (e?.message || ''))
  } finally {
    saving.value = false
  }
}

const removeRow = async (row: WebUser) => {
  try {
    await ElMessageBox.confirm(t('server.users.deleteConfirm', { user: row.user }), t('server.users.title'), { type: 'warning' })
  } catch {
    return
  }
  try {
    rows.value = rows.value.filter((r) => r.user !== row.user)
    delete pwMap[row.user]
    await updateWebUsers(buildPayload())
    ElMessage.success(t('server.users.deleted'))
    await fetchData()
  } catch (e: any) {
    ElMessage.error(t('server.users.saveFailed') + (e?.message || ''))
  }
}

fetchData()
</script>

<style scoped>
.users-page { display: flex; flex-direction: column; gap: 16px; }
.page-header { display: flex; align-items: center; justify-content: space-between; }
.page-title { font-size: 20px; font-weight: 600; margin: 0; }
.hint { color: var(--el-text-color-secondary); font-size: 12px; margin: 0; }
</style>
