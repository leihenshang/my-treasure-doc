import { createApp } from 'vue'
import App from './App.vue'
import './index.scss'
import { router } from './router'
import naive from 'naive-ui'
import VueAxiosPlugin from './plugin/axios/axios.ts'
import axios from "axios";
import { createPinia } from 'pinia'

// axios.defaults.baseURL = 'http://localhost:2021'

let app = createApp(App)
app.use(router)
app.use(naive)
app.use(VueAxiosPlugin, {
    // request interceptor handler
    reqHandleFunc: (config: any) => config,
    reqErrorFunc: (error: any) => Promise.reject(error),
    // response interceptor handler
    resHandleFunc: (response: any) => response,
    resErrorFunc: (error: any) => Promise.reject(error)
})

const pinia = createPinia()
app.use(pinia)

app.mount('#app')
