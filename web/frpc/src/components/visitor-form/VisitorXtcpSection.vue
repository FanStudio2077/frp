<template>
  <!-- XTCP Options -->
  <ConfigSection :title="t('client.form.xTCPOptions')" collapsible :readonly="readonly"
    :has-value="form.protocol !== 'quic' || form.keepTunnelOpen || form.maxRetriesAnHour != null || form.minRetryInterval != null || !!form.fallbackTo || form.fallbackTimeoutMs != null">
    <ConfigField :label="t('client.form.protocol')" type="select" v-model="form.protocol"
      :options="[{ label: 'QUIC', value: 'quic' }, { label: 'KCP', value: 'kcp' }]" :readonly="readonly" />
    <ConfigField :label="t('client.form.keepTunnelOpen')" type="switch" v-model="form.keepTunnelOpen" :readonly="readonly" />
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.maxRetriesPerHour')" type="number" v-model="form.maxRetriesAnHour" :min="0" :readonly="readonly" />
      <ConfigField :label="t('client.form.minRetryIntervalS')" type="number" v-model="form.minRetryInterval" :min="0" :readonly="readonly" />
    </div>
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.fallbackTo')" type="text" v-model="form.fallbackTo" placeholder="Fallback visitor name" :readonly="readonly" />
      <ConfigField :label="t('client.form.fallbackTimeoutMs')" type="number" v-model="form.fallbackTimeoutMs" :min="0" :readonly="readonly" />
    </div>
  </ConfigSection>

  <!-- NAT Traversal -->
  <ConfigSection :title="t('client.form.nATTraversal')" collapsible :readonly="readonly"
    :has-value="form.natTraversalDisableAssistedAddrs">
    <ConfigField :label="t('client.form.disableAssistedAddresses')" type="switch" v-model="form.natTraversalDisableAssistedAddrs"
      tip="Only use STUN-discovered public addresses" :readonly="readonly" />
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

const { t } = useI18n()
</script>

<style scoped lang="scss">
@use '@/assets/css/form-layout';
</style>