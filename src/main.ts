import { createApp } from 'vue'
import App from './App.vue'
import 'element-plus/dist/index.css'
import '@/assets/fonts/iconfont.css'

// dev 环境下引入 mock 模块
if (import.meta.env.DEV) {
  await import('@/mock')
}

// 挂载 Vue 实例
createApp(App).mount('#app')
