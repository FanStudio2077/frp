<template>
  <ConfigSection :title="t('client.form.authentication')" :readonly="readonly">
    <template v-if="['http', 'tcpmux'].includes(form.type)">
      <div class="field-row three-col">
        <ConfigField :label="t('client.form.hTTPUser')" type="text" v-model="form.httpUser" :readonly="readonly" />
        <ConfigField :label="t('client.form.hTTPPassword')" type="password" v-model="form.httpPassword" :readonly="readonly" />
        <ConfigField :label="t('client.form.routeByHTTPUser')" type="text" v-model="form.routeByHTTPUser" :readonly="readonly" />
      </div>
    </template>
    <template v-if="['stcp', 'sudp', 'xtcp'].includes(form.type)">
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.secretKey')" type="password" v-model="form.secretKey" prop="secretKey" :readonly="readonly" />
        <ConfigField :label="t('client.form.allowUsers')" type="tags" v-model="form.allowUsers" placeholder="username" :readonly="readonly" />
      </div>
    </template>
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