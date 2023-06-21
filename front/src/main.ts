import { createApp } from 'vue'
import App from './App.vue'
import './index.scss'
import { router } from './router'
import naive from 'naive-ui'
import axios from "axios";
import { createPinia } from 'pinia'

// axios.defaults.baseURL = 'http://localhost:2021'

let app = createApp(App)
app.use(router)
app.use(naive)
const pinia = createPinia()
app.use(pinia)

app.mount('#app')
