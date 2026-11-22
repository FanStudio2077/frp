<template>
  <ConfigSection :title="t('client.form.transport')" collapsible :readonly="readonly"
    :has-value="form.useEncryption || form.useCompression || !!form.bandwidthLimit || (!!form.bandwidthLimitMode && form.bandwidthLimitMode !== 'client') || !!form.proxyProtocolVersion">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.useEncryption')" type="switch" v-model="form.useEncryption" :readonly="readonly" />
      <ConfigField :label="t('client.form.useCompression')" type="switch" v-model="form.useCompression" :readonly="readonly" />
    </div>
    <div class="field-row three-col">
      <ConfigField :label="t('client.form.bandwidthLimit')" type="text" v-model="form.bandwidthLimit" placeholder="1MB" tip="e.g., 1MB, 500KB" :readonly="readonly" />
      <ConfigField :label="t('client.form.bandwidthLimitMode')" type="select" v-model="form.bandwidthLimitMode"
        :options="[{ label: 'Client', value: 'client' }, { label: 'Server', value: 'server' }]" :readonly="readonly" />
      <ConfigField :label="t('client.form.proxyProtocolVersion')" type="select" v-model="form.proxyProtocolVersion"
        :options="[{ label: 'None', value: '' }, { label: 'v1', value: 'v1' }, { label: 'v2', value: 'v2' }]" :readonly="readonly" />
    </div>
  </ConfigSection>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'
import type { ProxyFormData } from '../../types'
import ConfigSection from '../ConfigSection.vue'
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