<template>
  <ConfigSection :title="t('client.form.hTTPOptions')" collapsible :readonly="readonly"
    :has-value="form.locations.length > 0 || !!form.hostHeaderRewrite || form.requestHeaders.length > 0 || form.responseHeaders.length > 0">
    <ConfigField :label="t('client.form.locations')" type="tags" v-model="form.locations" placeholder="/path" :readonly="readonly" />
    <ConfigField :label="t('client.form.hostHeaderRewrite')" type="text" v-model="form.hostHeaderRewrite" :readonly="readonly" />
    <ConfigField :label="t('client.form.requestHeaders')" type="kv" v-model="form.requestHeaders" key-placeholder="Header" value-placeholder="Value" :readonly="readonly" />
    <ConfigField :label="t('client.form.responseHeaders')" type="kv" v-model="form.responseHeaders" key-placeholder="Header" value-placeholder="Value" :readonly="readonly" />
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