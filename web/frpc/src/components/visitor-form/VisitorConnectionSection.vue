<template>
  <ConfigSection :title="t('client.form.connection')" :readonly="readonly">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.serverName')" type="text" v-model="form.serverName" prop="serverName"
        placeholder="Name of the proxy to visit" :readonly="readonly" />
      <ConfigField :label="t('client.form.serverUser')" type="text" v-model="form.serverUser"
        placeholder="Leave empty for same user" :readonly="readonly" />
    </div>
    <ConfigField :label="t('client.form.secretKey')" type="password" v-model="form.secretKey"
      placeholder="Shared secret" :readonly="readonly" />
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.bindAddress')" type="text" v-model="form.bindAddr"
        placeholder="127.0.0.1" :readonly="readonly" />
      <ConfigField :label="t('client.form.bindPort')" type="number" v-model="form.bindPort"
        :min="bindPortMin" :max="65535" prop="bindPort" :readonly="readonly"
        :tip="bindPortTip" />
    </div>
  </ConfigSection>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'
import type { VisitorFormData } from '../../types'
import ConfigSection from '../ConfigSection.vue'
import ConfigField from '../ConfigField.vue'

const props = withDefaults(defineProps<{
  modelValue: VisitorFormData
  readonly?: boolean
}>(), { readonly: false })

const emit = defineEmits<{ 'update:modelValue': [value: VisitorFormData] }>()

const form = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const bindPortMin = computed(() => (form.value.type === 'sudp' ? 1 : undefined))
const bindPortTip = computed(() => form.value.type === 'sudp'
  ? ''
  : 'Use -1 to skip the local listener when connections come from another visitor or plugin.')

const { t } = useI18n()
</script>

<style scoped lang="scss">
@use '@/assets/css/form-layout';
</style>