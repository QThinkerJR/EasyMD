<script setup>
import { reactive, ref, watch, onMounted, onUnmounted } from 'vue'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import MarkdownEditor from './components/MarkdownEditor.vue'
import Toolbar from './components/Toolbar.vue'
import { OpenFile, ReadFile, SaveFile, SaveFileAs, ExportToHTMLAs, SetLanguage } from '../wailsjs/go/main/App.js'
import { OnFileDrop, OnFileDropOff } from '../wailsjs/runtime/runtime.js'
import { detectLanguage, t, currentLocale, state } from './i18n'

// 初始化语言
detectLanguage()
// 同步语言设置到后端
SetLanguage(state.locale)

// 应用状态
const appState = reactive({
  currentFilePath: '',
  markdownContent: '',
  htmlContent: '',
  isModified: false,
  isLoading: false,
  theme: 'light'
})

// 编辑器引用
const editorRef = ref(null)

// 切换主题
const toggleTheme = () => {
  appState.theme = appState.theme === 'light' ? 'dark' : 'light'
  // 保存主题设置到本地存储
  localStorage.setItem('theme', appState.theme)
  
  // 更新文档 body 的 class，用于全局样式控制
  if (appState.theme === 'dark') {
    document.body.classList.add('theme-dark')
  } else {
    document.body.classList.remove('theme-dark')
  }
}

// 初始化主题
const initTheme = () => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    appState.theme = savedTheme
    if (savedTheme === 'dark') {
      document.body.classList.add('theme-dark')
    }
  } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    appState.theme = 'dark'
    document.body.classList.add('theme-dark')
  }
}

// 在组件挂载前初始化主题
initTheme()

// 处理文件拖拽
const handleFileDrop = async (x, y, paths) => {
  console.log('触发文件拖拽事件', x, y, paths)
  if (paths && paths.length > 0) {
    const filePath = paths[0]
    // 只处理 .md 文件
    if (filePath.toLowerCase().endsWith('.md')) {
      try {
        appState.isLoading = true
        console.log('准备读取文件:', filePath)
        
        // 确保 ReadFile 已定义
        if (typeof ReadFile !== 'function') {
          throw new Error('ReadFile function is not defined. Please restart the application.')
        }

        const content = await ReadFile(filePath)
        
        if (content !== undefined) {
          appState.markdownContent = content
          appState.currentFilePath = filePath
          appState.isModified = false
          console.log('拖拽打开文件成功，路径已更新为:', appState.currentFilePath)
          MessagePlugin.success({ content: t('openSuccess') + ': ' + filePath, placement: 'center' })
        } else {
           console.error('读取文件内容为空或 undefined')
           MessagePlugin.warning({ content: t('fileContentEmpty'), placement: 'center' })
        }
      } catch (error) {
        console.error('拖拽打开文件失败:', error)
        MessagePlugin.error({ content: t('openFailed') + ': ' + error.message, placement: 'center' })
      } finally {
        appState.isLoading = false
      }
    } else {
      MessagePlugin.warning({ content: t('onlyMarkdown'), placement: 'center' })
    }
  } else {
    console.log('拖拽事件未包含有效文件路径')
  }
}

// 全局拖拽处理（DOM 层面）
const handleGlobalDrop = (e) => {
  // 阻止默认行为，防止浏览器打开文件或编辑器插入内容
  e.preventDefault()
  
  // 注意：不要调用 stopPropagation()，否则可能导致 Wails 的 OnFileDrop 监听器无法接收事件
  // e.stopPropagation()
  
  console.log('DOM Drop Event Triggered', e)
  
  // 尝试直接从 DOM 事件获取文件
  if (e.dataTransfer && e.dataTransfer.files && e.dataTransfer.files.length > 0) {
    const file = e.dataTransfer.files[0]
    console.log('Drop File:', file)
    
    // 某些 WebView (如 Electron) 会在 File 对象上提供 path 属性
    // Wails 在某些配置下可能不会在 File 对象上暴露 path
    if (file.path) {
       console.log('Found path in DOM event:', file.path)
       handleFileDrop(0, 0, [file.path])
       return
    }
  }
  
  console.log('Waiting for Wails OnFileDrop to trigger...')
}

onMounted(() => {
  // 添加全局捕获监听器，确保在任何子组件之前拦截
  window.addEventListener('drop', handleGlobalDrop, true) // true = capture
  window.addEventListener('dragover', (e) => e.preventDefault(), true)
  
  // 注册 Wails 文件拖拽监听
  // 第二个参数为 false，表示在整个窗口范围内都响应拖拽事件，不需要特定的 CSS 标记
  OnFileDrop(handleFileDrop, false)
})

onUnmounted(() => {
  window.removeEventListener('drop', handleGlobalDrop, true)
  // 移除文件拖拽监听
  OnFileDropOff()
})

// 监听内容变化，标记为已修改
watch(() => appState.markdownContent, () => {
  if (!appState.isLoading) {
    appState.isModified = true
  }
})

// 监听 HTML 内容变化
const onHtmlChanged = (html) => {
  appState.htmlContent = html
}

// 打开文件
const openFile = async () => {
  try {
    appState.isLoading = true
    console.log('开始打开文件...')
    const result = await OpenFile()
    console.log('打开文件结果:', result)
    
    // 检查返回结果的结构
    if (result) {
      // 如果返回的是数组 [content, filePath]
      if (Array.isArray(result) && result.length === 2) {
        const [content, filePath] = result
        console.log('文件内容:', content)
        console.log('文件路径:', filePath)
        
        if (content && filePath) {
          appState.markdownContent = content
          appState.currentFilePath = filePath
          appState.isModified = false
          console.log('文件内容已设置到编辑器')
        }
      }
      // 如果返回的是对象 {content: string, path: string}
      else if (typeof result === 'object' && result.content && result.path) {
        console.log('文件内容:', result.content)
        console.log('文件路径:', result.path)
        
        appState.markdownContent = result.content
        appState.currentFilePath = result.path
        appState.isModified = false
        console.log('文件内容已设置到编辑器')
      }
      // 如果返回的是字符串（只有内容）
      else if (typeof result === 'string') {
        console.log('文件内容:', result)
        appState.markdownContent = result
        appState.currentFilePath = ''
        appState.isModified = false
        console.log('文件内容已设置到编辑器')
      }
    } else {
      console.log('用户取消了文件选择或文件为空')
    }
  } catch (error) {
    console.error('打开文件失败:', error)
  } finally {
    appState.isLoading = false
  }
}

// 保存文件
const saveFile = async () => {
  try {
    if (appState.currentFilePath) {
      await SaveFile(appState.markdownContent, appState.currentFilePath)
      appState.isModified = false
      MessagePlugin.success({ content: t('saveSuccess'), placement: 'center' })
    } else {
      // 如果没有当前文件路径，则另存为
      const filePath = await SaveFileAs(appState.markdownContent)
      if (filePath) {
        appState.currentFilePath = filePath
        appState.isModified = false
        MessagePlugin.success({ content: t('saveSuccess'), placement: 'center' })
      }
    }
  } catch (error) {
    console.error('保存文件失败:', error)
    MessagePlugin.error({ content: t('saveFailed'), placement: 'center' })
  }
}

// 另存为
const saveFileAs = async () => {
  try {
    const filePath = await SaveFileAs(appState.markdownContent)
    if (filePath) {
      appState.currentFilePath = filePath
      appState.isModified = false
      MessagePlugin.success({ content: t('saveSuccess'), placement: 'center' })
    }
  } catch (error) {
    console.error('另存为失败:', error)
    MessagePlugin.error({ content: t('saveFailed'), placement: 'center' })
  }
}

// 导出为HTML
const exportToHTML = async () => {
  try {
    // 使用 Markdown 原始内容进行导出，而不是渲染后的 HTML
    // 这样后端可以正确处理 Mermaid 和 ECharts 代码块
    const result = await ExportToHTMLAs(appState.markdownContent)
    if (result) {
      MessagePlugin.success({ content: t('exportSuccess'), placement: 'center' })
    }
  } catch (error) {
    console.error('导出HTML失败:', error)
    MessagePlugin.error({ content: t('exportFailed'), placement: 'center' })
  }
}

// 导出为PDF
const exportToPDF = () => {
  if (editorRef.value) {
    editorRef.value.exportPDF()
  }
}

// 新建文件
const newFile = () => {
  if (appState.isModified) {
    // 使用 TDesign 的 Dialog 组件显示确认对话框
    const dialogInstance = DialogPlugin.confirm({
      header: t('prompt'),
      body: t('unsavedWarning'),
      confirmBtn: t('continue'),
      cancelBtn: t('cancel'),
      onConfirm: () => {
        appState.markdownContent = ''
        appState.currentFilePath = ''
        appState.isModified = false
        dialogInstance.hide()
      },
      onCancel: () => {
        dialogInstance.hide()
      }
    })
  } else {
    appState.markdownContent = ''
    appState.currentFilePath = ''
    appState.isModified = false
  }
}

// 获取窗口标题
const getWindowTitle = () => {
  const fileName = appState.currentFilePath ? 
    appState.currentFilePath.split(/[/\\]/).pop() : t('untitled')
  const modified = appState.isModified ? ' *' : ''
  return `${fileName}${modified} - ${t('editorTitle')}`
}
</script>

<template>
  <div class="app-container">
    <Toolbar 
      :current-file-path="appState.currentFilePath"
      :is-modified="appState.isModified"
      :theme="appState.theme"
      @new-file="newFile"
      @open-file="openFile"
      @save-file="saveFile"
      @save-file-as="saveFileAs"
      @export-html="exportToHTML"
      @export-pdf="exportToPDF"
      @toggle-theme="toggleTheme"
    />
    
    <div class="editor-container">
      <MarkdownEditor 
        ref="editorRef"
        v-model="appState.markdownContent"
        :is-loading="appState.isLoading"
        :theme="appState.theme"
        :language="currentLocale"
        @html-changed="onHtmlChanged"
      />
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
}

body.theme-dark {
  background-color: #0d1117;
  color: #e6edf3;
}

#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.editor-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>
