import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import i18n from './i18n'
import App from './App.vue'
import router from './router'
import api from './api'
import './style.css'

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(i18n)
app.config.globalProperties.$api = api
app.mount('#app')