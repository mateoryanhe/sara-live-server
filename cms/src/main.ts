import './assets/main.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import {restoreAuthSession} from '@/utils/auth'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/message-box/style/css'
import 'element-plus/es/components/overlay/style/css'

restoreAuthSession()

const app = createApp(App)

app.use(router)
app.mount('#app')
