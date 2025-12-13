<script setup>
import { ref, computed, watch } from 'vue'
import { Button, Space } from 'tdesign-vue-next'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  isLoading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

// 分屏模式：'edit', 'preview', 'split'
const viewMode = ref('split')

// 编辑器内容
const content = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 引用编辑器和预览区域
const editorTextarea = ref(null)
const previewContent = ref(null)

// 同步滚动标志，防止循环触发
let isSyncing = false

// 简单的Markdown到HTML转换
const markdownToHtml = (markdown) => {
  if (!markdown) return ''
  
  let html = markdown
  
  // 先处理代码块（避免与其他规则冲突）
  html = html.replace(/```(mermaid|gantt)?\n([\s\S]*?)```/g, (match, type, code) => {
    if (type === 'mermaid') {
      return `<div class="mermaid-diagram">${code.trim()}</div>`
    } else if (type === 'gantt') {
      return `<div class="gantt-diagram">${code.trim()}</div>`
    } else {
      return `<pre class="code-block"><code class="language-${type || 'text'}">${code.trim()}</code></pre>`
    }
  })
  
  // 处理表格
  html = html.replace(/\|(.+)\|\n\|[-\s|]+\|((?:\n\|.+\|)*)/g, (match, header, body) => {
    const headerCells = header.split('|').filter(cell => cell.trim()).map(cell => `<th>${cell.trim()}</th>`).join('')
    const rows = body.trim().split('\n').map(row => {
      const cells = row.split('|').filter(cell => cell.trim()).map(cell => `<td>${cell.trim()}</td>`).join('')
      return `<tr>${cells}</tr>`
    }).join('')
    return `<table class="markdown-table"><thead><tr>${headerCells}</tr></thead><tbody>${rows}</tbody></table>`
  })
  
  // 处理数学公式（块级）
  html = html.replace(/\$\$([^$]+)\$\$/g, '<div class="math-formula">$1</div>')
  
  // 处理数学公式（行内）
  html = html.replace(/\$([^$]+)\$/g, '<span class="inline-math">$1</span>')
  
  // 处理行内代码
  html = html.replace(/`([^`]+)`/g, '<code class="inline-code">$1</code>')
  
  // 处理标题
  html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>')
  
  // 处理文本格式
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')
  
  // 处理链接
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank">$1</a>')
  
  // 处理列表
  html = html.replace(/^\d+\.\s+(.+)$/gim, '<li class="ordered-list-item">$1</li>')
  html = html.replace(/^[-*+]\s+(.+)$/gim, '<li class="unordered-list-item">$1</li>')
  
  // 处理引用
  html = html.replace(/^>\s+(.+)$/gim, '<blockquote>$1</blockquote>')
  
  // 处理段落
  html = html.replace(/\n\n/gim, '</p><p>')
  html = html.replace(/^/, '<p>')
  html = html.replace(/$/, '</p>')
  
  // 清理空段落和表格周围的段落标签
  html = html.replace(/<p><\/p>/g, '')
  html = html.replace(/<p>\s*<\/p>/g, '')
  html = html.replace(/<p>(<table[\s\S]*?<\/table>)<\/p>/g, '$1')
  html = html.replace(/<p>(<div class="(?:mermaid-diagram|gantt-diagram|math-formula)[\s\S]*?<\/div>)<\/p>/g, '$1')
  html = html.replace(/<p>(<pre class="code-block"[\s\S]*?<\/pre>)<\/p>/g, '$1')
  
  return html
}

// 计算预览HTML
const previewHtml = computed(() => {
  return markdownToHtml(content.value)
})

// 切换视图模式
const setViewMode = (mode) => {
  viewMode.value = mode
}

// 同步滚动处理
const handleEditorScroll = () => {
  if (isSyncing || viewMode.value !== 'split') return
  
  const editor = editorTextarea.value
  const preview = previewContent.value
  
  if (!editor || !preview) return
  
  isSyncing = true
  
  // 计算滚动比例
  const editorScrollRatio = editor.scrollTop / (editor.scrollHeight - editor.clientHeight)
  
  // 同步预览区域滚动
  const previewScrollTop = editorScrollRatio * (preview.scrollHeight - preview.clientHeight)
  preview.scrollTop = previewScrollTop
  
  setTimeout(() => {
    isSyncing = false
  }, 10)
}

const handlePreviewScroll = () => {
  if (isSyncing || viewMode.value !== 'split') return
  
  const editor = editorTextarea.value
  const preview = previewContent.value
  
  if (!editor || !preview) return
  
  isSyncing = true
  
  // 计算滚动比例
  const previewScrollRatio = preview.scrollTop / (preview.scrollHeight - preview.clientHeight)
  
  // 同步编辑器滚动
  const editorScrollTop = previewScrollRatio * (editor.scrollHeight - editor.clientHeight)
  editor.scrollTop = editorScrollTop
  
  setTimeout(() => {
    isSyncing = false
  }, 10)
}
</script>

<template>
  <div class="markdown-editor">
    <!-- 视图模式切换按钮 -->
    <div class="view-mode-switcher">
      <Space>
        <Button 
          :variant="viewMode === 'edit' ? 'base' : 'text'"
          size="small"
          @click="setViewMode('edit')"
        >
          编辑
        </Button>
        <Button 
          :variant="viewMode === 'split' ? 'base' : 'text'"
          size="small"
          @click="setViewMode('split')"
        >
          分屏
        </Button>
        <Button 
          :variant="viewMode === 'preview' ? 'base' : 'text'"
          size="small"
          @click="setViewMode('preview')"
        >
          预览
        </Button>
      </Space>
    </div>
    
    <!-- 编辑器容器 -->
    <div class="editor-container" :class="{ 'loading': isLoading }">
      <div class="editor-wrapper" v-show="viewMode === 'edit' || viewMode === 'split'">
        <!-- 编辑模式 -->
        <div
          class="editor-pane"
          :class="{ 'full-width': viewMode === 'edit' }"
        >
          <textarea
            ref="editorTextarea"
            v-model="content"
            class="markdown-textarea"
            placeholder="在此输入Markdown内容..."
            @scroll="handleEditorScroll"
          />
        </div>
        
        <!-- 预览模式 -->
        <div
          v-show="viewMode === 'preview' || viewMode === 'split'"
          class="preview-pane"
          :class="{ 'full-width': viewMode === 'preview' }"
        >
          <div
            ref="previewContent"
            class="preview-content"
            v-html="previewHtml"
            @scroll="handlePreviewScroll"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.markdown-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.view-mode-switcher {
  padding: 8px 16px;
  background-color: #fafafa;
  border-bottom: 1px solid #e0e0e0;
}

.editor-container {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
}

.editor-container.loading {
  opacity: 0.6;
  pointer-events: none;
}

.editor-wrapper {
  display: flex;
  width: 100%;
  height: 100%;
}

.editor-pane, .preview-pane {
  height: 100%;
  overflow: hidden;
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.editor-pane {
  border-right: 2px solid #e0e0e0;
}

.editor-pane.full-width, .preview-pane.full-width {
  flex: 1 1 100%;
  border-right: none;
}

.preview-pane {
  background-color: #fff;
}

.markdown-textarea {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  padding: 20px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  background-color: #fff;
  color: #333;
  box-sizing: border-box;
}

.preview-content {
  height: 100%;
  overflow: auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  line-height: 1.6;
  color: #333;
  text-align: left;
}

.preview-content h1,
.preview-content h2,
.preview-content h3,
.preview-content h4,
.preview-content h5,
.preview-content h6 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.preview-content code {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-size: 85%;
}

.preview-content p {
  margin-bottom: 16px;
  text-align: left;
}

.preview-content a {
  color: #1890ff;
  text-decoration: none;
}

.preview-content a:hover {
  text-decoration: underline;
}

/* 代码块样式 */
.preview-content .code-block {
  background-color: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  margin: 16px 0;
  overflow-x: auto;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.45;
}

.preview-content .inline-code {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 85%;
}

/* 数学公式样式 */
.preview-content .math-formula {
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 16px;
  margin: 16px 0;
  text-align: center;
  font-family: 'Cambria Math', 'Times New Roman', serif;
  font-size: 16px;
}

.preview-content .inline-math {
  font-family: 'Cambria Math', 'Times New Roman', serif;
  font-style: italic;
}

/* 列表样式 */
.preview-content .ordered-list-item,
.preview-content .unordered-list-item {
  margin: 4px 0;
  padding-left: 20px;
  position: relative;
}

.preview-content .ordered-list-item::before {
  content: attr(data-counter) ".";
  position: absolute;
  left: 0;
  font-weight: bold;
}

.preview-content .unordered-list-item::before {
  content: "•";
  position: absolute;
  left: 0;
}

/* 引用样式 */
.preview-content blockquote {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  margin: 16px 0;
  color: #6a737d;
  font-style: italic;
}

/* 表格样式 */
.preview-content .markdown-table {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
  border: 1px solid #dfe2e5;
  border-radius: 6px;
  overflow: hidden;
}

.preview-content .markdown-table th,
.preview-content .markdown-table td {
  border: 1px solid #dfe2e5;
  padding: 8px 12px;
  text-align: left;
}

.preview-content .markdown-table th {
  background-color: #f6f8fa;
  font-weight: 600;
}

.preview-content .markdown-table tr:nth-child(even) {
  background-color: #f9f9f9;
}

.preview-content .markdown-table tr:hover {
  background-color: #f5f5f5;
}

/* 流程图样式 */
.preview-content .mermaid-diagram,
.preview-content .gantt-diagram {
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  padding: 16px;
  margin: 16px 0;
  text-align: center;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  overflow-x: auto;
  white-space: pre;
  color: #333;
}
</style>