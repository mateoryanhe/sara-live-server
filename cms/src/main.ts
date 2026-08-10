import './assets/main.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import {restoreAuthSession} from '@/utils/auth'
import {btnPermission} from '@/directives/btn-permission'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/message-box/style/css'
import 'element-plus/es/components/overlay/style/css'

restoreAuthSession()

const app = createApp(App)

app.directive('btn-permission', btnPermission)
app.use(router)
app.mount('#app')
