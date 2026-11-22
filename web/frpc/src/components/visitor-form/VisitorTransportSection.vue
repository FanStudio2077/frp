<template>
  <ConfigSection :title="t('client.form.transportOptions')" collapsible :readonly="readonly"
    :has-value="form.useEncryption || form.useCompression">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.useEncryption')" type="switch" v-model="form.useEncryption" :readonly="readonly" />
      <ConfigField :label="t('client.form.useCompression')" type="switch" v-model="form.useCompression" :readonly="readonly" />
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

const { t } = useI18n()
</script>

<style scoped lang="scss">
@use '@/assets/css/form-layout';
</style>