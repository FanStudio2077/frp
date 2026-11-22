<template>
  <ConfigSection :title="t('client.form.loadBalancer')" collapsible :readonly="readonly" :has-value="!!form.loadBalancerGroup">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.group')" type="text" v-model="form.loadBalancerGroup" placeholder="Group name" :readonly="readonly" />
      <ConfigField :label="t('client.form.groupKey')" type="text" v-model="form.loadBalancerGroupKey" :readonly="readonly" />
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