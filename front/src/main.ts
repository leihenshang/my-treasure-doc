import { createApp } from 'vue'
import App from './App.vue'
import './index.css'
import Layui from '@layui/layui-vue'
import '@layui/layui-vue/lib/index.css'
import {router} from './router';

let app = createApp(App)
app.use(router)
app.use(Layui)
app.mount('#app')
