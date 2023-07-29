import { createApp } from 'vue'
import App from './App.vue'

// https://element-plus.org/en-US/guide/quickstart.html#full-import
// https://element-plus.org/zh-CN/guide/installation.html#hello-world
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)
app.use(ElementPlus)

app.mount('#app')
