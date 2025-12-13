<script setup>
import { reactive, ref, watch } from 'vue'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import MarkdownEditor from './components/MarkdownEditor.vue'
import Toolbar from './components/Toolbar.vue'
import { OpenFile, SaveFile, SaveFileAs, ExportToHTMLAs } from '../wailsjs/go/main/App.js'

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
      MessagePlugin.success({ content: '保存成功', placement: 'center' })
    } else {
      // 如果没有当前文件路径，则另存为
      const filePath = await SaveFileAs(appState.markdownContent)
      if (filePath) {
        appState.currentFilePath = filePath
        appState.isModified = false
        MessagePlugin.success({ content: '保存成功', placement: 'center' })
      }
    }
  } catch (error) {
    console.error('保存文件失败:', error)
  }
}

// 另存为
const saveFileAs = async () => {
  try {
    const filePath = await SaveFileAs(appState.markdownContent)
    if (filePath) {
      appState.currentFilePath = filePath
      appState.isModified = false
      MessagePlugin.success({ content: '保存成功', placement: 'center' })
    }
  } catch (error) {
    console.error('另存为失败:', error)
  }
}

// 导出为HTML
const exportToHTML = async () => {
  try {
    // 使用 Markdown 原始内容进行导出，而不是渲染后的 HTML
    // 这样后端可以正确处理 Mermaid 和 ECharts 代码块
    await ExportToHTMLAs(appState.markdownContent)
    MessagePlugin.success({ content: '导出成功', placement: 'center' })
  } catch (error) {
    console.error('导出HTML失败:', error)
    MessagePlugin.error({ content: '导出HTML失败', placement: 'center' })
  }
}

// 新建文件
const newFile = () => {
  if (appState.isModified) {
    // 使用 TDesign 的 Dialog 组件显示确认对话框
    const dialogInstance = DialogPlugin.confirm({
      header: '提示',
      body: '当前文件尚未保存，是否继续新建？未保存的内容将会丢失。',
      confirmBtn: '继续',
      cancelBtn: '取消',
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
    appState.currentFilePath.split(/[/\\]/).pop() : '未命名'
  const modified = appState.isModified ? ' *' : ''
  return `${fileName}${modified} - Markdown编辑器`
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
      @toggle-theme="toggleTheme"
    />
    
    <div class="editor-container">
      <MarkdownEditor 
        ref="editorRef"
        v-model="appState.markdownContent"
        :is-loading="appState.isLoading"
        :theme="appState.theme"
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
