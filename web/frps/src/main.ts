import { createApp } from 'vue'
import 'element-plus/theme-chalk/dark/css-vars.css'
import App from './App.vue'
import router from './router'
import { createI18nForApp } from '@shared/i18n'

import './assets/css/var.css'
import './assets/css/dark.css'

const app = createApp(App)

app.use(createI18nForApp('frps'))
app.use(router)

app.mount('#app')
