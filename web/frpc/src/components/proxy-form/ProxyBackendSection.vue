<template>
  <!-- Backend Mode -->
  <template v-if="!readonly">
    <el-form-item :label="t('client.form.backendMode')">
      <el-radio-group v-model="backendMode">
        <el-radio value="direct">{{ t('common.direct') }}</el-radio>
        <el-radio value="plugin">{{ t('client.form.plugin') }}</el-radio>
      </el-radio-group>
    </el-form-item>
  </template>

  <!-- Direct mode -->
  <template v-if="backendMode === 'direct'">
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.localIP')" type="text" v-model="form.localIP" placeholder="127.0.0.1" :readonly="readonly" />
      <ConfigField :label="t('client.form.localPort')" type="number" v-model="form.localPort" :min="0" :max="65535" prop="localPort" :readonly="readonly" />
    </div>
  </template>

  <!-- Plugin mode -->
  <template v-else>
    <div class="field-row two-col">
      <ConfigField :label="t('client.form.pluginType')" type="select" v-model="form.pluginType"
        :options="PLUGIN_LIST.map((p) => ({ label: p, value: p }))" :readonly="readonly" />
      <div></div>
    </div>

    <template v-if="['http2https', 'https2http', 'https2https', 'http2http', 'tls2raw'].includes(form.pluginType)">
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.localAddress')" type="text" v-model="form.pluginConfig.localAddr" placeholder="127.0.0.1:8080" :readonly="readonly" />
        <ConfigField v-if="['http2https', 'https2http', 'https2https', 'http2http'].includes(form.pluginType)"
          :label="t('client.form.hostHeaderRewrite')" type="text" v-model="form.pluginConfig.hostHeaderRewrite" :readonly="readonly" />
        <div v-else></div>
      </div>
    </template>
    <template v-if="['http2https', 'https2http', 'https2https', 'http2http'].includes(form.pluginType)">
      <ConfigField :label="t('client.form.requestHeaders')" type="kv" v-model="pluginRequestHeaders"
        key-placeholder="Header" value-placeholder="Value" :readonly="readonly" />
    </template>
    <template v-if="['https2http', 'https2https', 'tls2raw'].includes(form.pluginType)">
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.certificatePath')" type="text" v-model="form.pluginConfig.crtPath" placeholder="/path/to/cert.pem" :readonly="readonly" />
        <ConfigField :label="t('client.form.keyPath')" type="text" v-model="form.pluginConfig.keyPath" placeholder="/path/to/key.pem" :readonly="readonly" />
      </div>
    </template>
    <template v-if="['https2http', 'https2https'].includes(form.pluginType)">
      <ConfigField :label="t('client.form.enableHTTP2')" type="switch" v-model="form.pluginConfig.enableHTTP2" :readonly="readonly" />
    </template>
    <template v-if="form.pluginType === 'http_proxy'">
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.hTTPUser')" type="text" v-model="form.pluginConfig.httpUser" :readonly="readonly" />
        <ConfigField :label="t('client.form.hTTPPassword')" type="password" v-model="form.pluginConfig.httpPassword" :readonly="readonly" />
      </div>
    </template>
    <template v-if="form.pluginType === 'socks5'">
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.username')" type="text" v-model="form.pluginConfig.username" :readonly="readonly" />
        <ConfigField :label="t('client.form.password')" type="password" v-model="form.pluginConfig.password" :readonly="readonly" />
      </div>
    </template>
    <template v-if="form.pluginType === 'static_file'">
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.localPath')" type="text" v-model="form.pluginConfig.localPath" placeholder="/path/to/files" :readonly="readonly" />
        <ConfigField :label="t('client.form.stripPrefix')" type="text" v-model="form.pluginConfig.stripPrefix" :readonly="readonly" />
      </div>
      <div class="field-row two-col">
        <ConfigField :label="t('client.form.hTTPUser')" type="text" v-model="form.pluginConfig.httpUser" :readonly="readonly" />
        <ConfigField :label="t('client.form.hTTPPassword')" type="password" v-model="form.pluginConfig.httpPassword" :readonly="readonly" />
      </div>
    </template>
    <template v-if="form.pluginType === 'unix_domain_socket'">
      <ConfigField :label="t('client.form.unixSocketPath')" type="text" v-model="form.pluginConfig.unixPath" placeholder="/tmp/socket.sock" :readonly="readonly" />
    </template>
  </template>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import type { ProxyFormData } from '../../types'
import ConfigField from '../ConfigField.vue'

const PLUGIN_LIST = [
  'http2https', 'http_proxy', 'https2http', 'https2https', 'http2http',
  'socks5', 'static_file', 'unix_domain_socket', 'tls2raw', 'virtual_net',
]

const props = withDefaults(defineProps<{
  modelValue: ProxyFormData
  readonly?: boolean
}>(), { readonly: false })

const emit = defineEmits<{ 'update:modelValue': [value: ProxyFormData] }>()

const form = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const backendMode = ref<'direct' | 'plugin'>(form.value.pluginType ? 'plugin' : 'direct')
const isHydrating = ref(false)

const pluginRequestHeaders = computed({
  get() {
    const set = form.value.pluginConfig?.requestHeaders?.set
    if (!set || typeof set !== 'object') return []
    return Object.entries(set).map(([key, value]) => ({ key, value: String(value) }))
  },
  set(val: Array<{ key: string; value: string }>) {
    if (!form.value.pluginConfig) form.value.pluginConfig = {}
    if (val.length === 0) {
      delete form.value.pluginConfig.requestHeaders
    } else {
      form.value.pluginConfig.requestHeaders = {
        set: Object.fromEntries(val.map((e) => [e.key, e.value])),
      }
    }
  },
})

watch(() => form.value.pluginType, (newType, oldType) => {
  if (isHydrating.value) return
  if (!oldType || !newType || newType === oldType) return
  if (form.value.pluginConfig && Object.keys(form.value.pluginConfig).length > 0) {
    form.value.pluginConfig = {}
  }
})

watch(backendMode, (mode) => {
  if (mode === 'direct') {
    form.value.pluginType = ''
    form.value.pluginConfig = {}
  } else if (!form.value.pluginType) {
    form.value.pluginType = 'http2https'
  }
})

const hydrate = () => {
  isHydrating.value = true
  backendMode.value = form.value.pluginType ? 'plugin' : 'direct'
  nextTick(() => { isHydrating.value = false })
}

watch(() => props.modelValue, () => { hydrate() })
onMounted(() => { hydrate() })

const { t } = useI18n()
</script>

<style scoped lang="scss">
@use '@/assets/css/form-layout';
</style>