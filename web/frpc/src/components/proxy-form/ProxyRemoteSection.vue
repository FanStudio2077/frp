<template>
  <template v-if="['tcp', 'udp'].includes(form.type)">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.remotePort')" type="number" v-model="form.remotePort"
        :min="0" :max="65535" prop="remotePort" tip="Use 0 for random port assignment" :readonly="readonly" />
      <div></div>
    </div>
  </template>
  <template v-if="['http', 'https', 'tcpmux'].includes(form.type)">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.customDomains')" type="tags" v-model="form.customDomains"
        prop="customDomains" placeholder="example.com" :readonly="readonly" />
      <ConfigField v-if="form.type !== 'tcpmux'" :label="t('client.form.subdomain')" type="text"
        v-model="form.subdomain" placeholder="test" :readonly="readonly" />
      <ConfigField v-if="form.type === 'tcpmux'" :label="t('client.form.multiplexer')" type="select"
        v-model="form.multiplexer" :options="[{ label: 'HTTP CONNECT', value: 'httpconnect' }]" :readonly="readonly" />
    </div>
  </template>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'
import type { ProxyFormData } from '../../types'
import ConfigField from '../ConfigField.vue'

const props = withDefaults(defineProps<{
  modelValue: ProxyFormData
  readonly?: boolean
}>(), { readonly: false })

const emit = defineEmits<{ 'update:modelValue': [value: ProxyFormData] }>()

const form = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const { t } = useI18n()
</script>

<style scoped lang="scss">
@use '@/assets/css/form-layout';
</style>