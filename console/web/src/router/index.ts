import { createRouter, createWebHashHistory } from 'vue-router'
import Overview from '../views/Overview.vue'
import Servers from '../views/Servers.vue'
import Nodes from '../views/Nodes.vue'
import Proxies from '../views/Proxies.vue'
import Visitors from '../views/Visitors.vue'
import Users from '../views/Users.vue'
import Audit from '../views/Audit.vue'
import Alerts from '../views/Alerts.vue'
import Settings from '../views/Settings.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/overview' },
    { path: '/overview', name: 'Overview', component: Overview },
    { path: '/servers', name: 'Servers', component: Servers },
    { path: '/nodes', name: 'Nodes', component: Nodes },
    { path: '/proxies', name: 'Proxies', component: Proxies },
    { path: '/visitors', name: 'Visitors', component: Visitors },
    { path: '/users', name: 'Users', component: Users },
    { path: '/audit', name: 'Audit', component: Audit },
    { path: '/alerts', name: 'Alerts', component: Alerts },
    { path: '/settings', name: 'Settings', component: Settings },
  ],
})

export default router
