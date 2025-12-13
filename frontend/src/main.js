import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

// 导入TDesign
import TDesign from 'tdesign-vue-next'
import 'tdesign-vue-next/es/style/index.css'

// 导入md-editor-v3样式
import 'md-editor-v3/lib/style.css'

const app = createApp(App)

// 添加全局错误处理，方便调试白屏问题
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue Error:', err)
  console.error('Info:', info)
  document.body.innerHTML = `<div style="padding: 20px; color: red; background: #fff; z-index: 9999; position: fixed; top: 0; left: 0; width: 100%; height: 100%; overflow: auto;">
    <h1>Runtime Error</h1>
    <pre style="white-space: pre-wrap;">${err.toString()}</pre>
    <pre style="white-space: pre-wrap;">${err.stack}</pre>
  </div>`
}

window.addEventListener('error', (event) => {
  // 忽略 ResizeObserver 循环错误，这是浏览器的一个已知良性警告
  if (event.message === 'ResizeObserver loop completed with undelivered notifications.' || 
      event.message === 'ResizeObserver loop limit exceeded') {
    return
  }

  console.error('Global Error:', event.error)
  document.body.innerHTML = `<div style="padding: 20px; color: red; background: #fff; z-index: 9999; position: fixed; top: 0; left: 0; width: 100%; height: 100%; overflow: auto;">
    <h1>Global Error</h1>
    <pre style="white-space: pre-wrap;">${event.error?.toString() || event.message}</pre>
  </div>`
})

app.use(TDesign)
app.mount('#app')
