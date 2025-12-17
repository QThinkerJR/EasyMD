<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { MdEditor, config } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { Emoji, Mark, ExportPDF } from '@vavt/v3-extension'
import '@vavt/v3-extension/lib/asset/style.css'
import MarkExtension from 'markdown-it-mark'
import { t } from '../i18n'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime.js'

// 配置 markdown-it 扩展
config({
  markdownItConfig(md) {
    md.use(MarkExtension)
  }
})

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  theme: {
    type: String,
    default: 'light'
  },
  language: {
    type: String,
    default: 'zh-CN'
  }
})

const emit = defineEmits(['update:modelValue', 'htmlChanged'])

// 编辑器引用
const editorRef = ref(null)
const containerRef = ref(null)

// 处理链接点击
const handleLinkClick = (e) => {
  const target = e.target.closest('a')
  if (target) {
    const href = target.getAttribute('href')
    // 如果是外部链接（不以 # 开头），则拦截并使用系统浏览器打开
    if (href && !href.startsWith('#')) {
      e.preventDefault()
      BrowserOpenURL(href)
    }
  }
}

// 编辑器内容
const content = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 编辑器工具栏配置 - 使用官方扩展库的 Emoji 和 Mark 组件
const toolbars = [
  'bold',
  'italic',
  'underline',
  'strikeThrough',
  'title',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  'task',
  '-',
  'codeRow',
  'code',
  'link',
  'image',
  'table',
  'mermaid',
  'katex',
  'revoke',
  'next',
  '-',
  0,  // Mark按钮（官方扩展）
  1,  // Emoji按钮（官方扩展）
  2,  // ExportPDF按钮（官方扩展）
  '=',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'previewOnly',
  'catalog'
]

// 观察者实例
let observer = null

// 检查并隐藏菜单项
const checkAndHide = (item) => {
  const text = item.textContent?.trim() || ''
  const title = item.getAttribute('title') || ''
  // 隐藏包含“上传”或“裁剪”的选项 (中英文)
  const keywords = ['上传', '裁剪', 'Upload', 'Clip']
  if (keywords.some(k => text.includes(k) || title.includes(k))) {
    item.style.display = 'none'
  }
}

// 在组件挂载后，隐藏图片按钮下拉菜单中的上传和裁剪选项
onMounted(() => {
  // 添加链接点击监听
  if (containerRef.value) {
    containerRef.value.addEventListener('click', handleLinkClick)
  }

  // 使用 MutationObserver 监听 DOM 变化，处理动态生成的菜单项
  observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.addedNodes.length) {
        mutation.addedNodes.forEach((node) => {
          if (node.nodeType === 1) { // Element
            if (node.classList && node.classList.contains('md-editor-menu-item')) {
              checkAndHide(node)
            } else if (node.querySelectorAll) {
              const items = node.querySelectorAll('.md-editor-menu-item')
              items.forEach(checkAndHide)
            }
          }
        })
      }
    })
  })
  
  // 观察 body，因为 dropdown 可能被 teleport 到 body
  if (document.body) {
    observer.observe(document.body, { childList: true, subtree: true })
  }
  
  // 初始检查
  document.querySelectorAll('.md-editor-menu-item').forEach(checkAndHide)
})

// 组件卸载时移除监听和观察者
onUnmounted(() => {
  if (containerRef.value) {
    containerRef.value.removeEventListener('click', handleLinkClick)
  }
  if (observer) {
    observer.disconnect()
  }
})

// 导出 PDF 方法
const exportPDF = () => {
  // 尝试查找导出PDF按钮并点击
  // 首先尝试通过 title 查找
  let btn = document.querySelector('.md-editor-toolbar-item[title*="PDF"]')
  
  if (!btn) {
    // 尝试查找 svg 的 data-def 属性或其他特征
    // 这里做一个简单的遍历查找
    const items = document.querySelectorAll('.md-editor-toolbar-item')
    for (const item of items) {
      if (item.getAttribute('title')?.includes('PDF')) {
        btn = item
        break
      }
    }
  }

  if (btn) {
    btn.click()
  } else {
    console.warn('PDF export button not found in toolbar')
    // 备用方案：如果按钮有特定的索引位置，可以尝试硬编码索引
    // 但由于有隐藏元素和分隔符，这不可靠
  }
}

// 暴露 onHtmlChanged 事件
const onHtmlChanged = (html) => {
  emit('htmlChanged', html)
}

defineExpose({
  exportPDF
})
</script>

<template>
  <div class="markdown-editor" ref="containerRef">
    <!-- 编辑器容器 -->
    <div class="md-editor-container" :class="{ 'loading': isLoading }">
      <MdEditor
        ref="editorRef"
        v-model="content"
        :theme="theme"
        :language="language"
        :preview="true"
        :toolbars="toolbars"
        :toolbarsExclude="[]"
        class="md-editor"
        @onHtmlChanged="onHtmlChanged"
      >
        <template #defToolbars>
          <Mark>
            <template #trigger>
              <svg class="md-editor-icon" viewBox="0 0 24 24" width="24" height="24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z" fill="currentColor"/>
              </svg>
            </template>
          </Mark>
          <Emoji>
            <template #trigger>
              <svg class="md-editor-icon" viewBox="0 0 24 24" width="24" height="24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z" fill="currentColor"/>
              </svg>
            </template>
          </Emoji>
          <ExportPDF :modelValue="content" />
        </template>
      </MdEditor>
    </div>
  </div>
</template>

<style scoped>
.markdown-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.md-editor-container {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.md-editor-container.loading {
  opacity: 0.6;
  pointer-events: none;
}

.md-editor {
  height: 100%;
}

/* 自定义md-editor-v3样式 */
:deep(.md-editor) {
  height: 100%;
}

/* 强制内容左对齐 - 针对 CodeMirror 编辑器 */
:deep(.cm-editor),
:deep(.cm-scroller),
:deep(.cm-content),
:deep(.cm-line) {
  text-align: left !important;
  font-family: Arial, "Microsoft YaHei", sans-serif !important;
}

/* 强制预览区域左对齐 */
:deep(.md-editor-preview-wrapper) {
  padding: 0 20px !important;
  text-align: left !important;
}

:deep(.md-editor-preview) {
  font-family: Arial, "Microsoft YaHei", sans-serif !important;
  line-height: 1.6;
  color: #333;
  text-align: left !important;
  width: 100% !important;
  max-width: none !important;
  margin: 0 !important;
  padding: 0 !important;
}

/* 适配暗黑模式 */
:deep(.md-editor-dark .md-editor-preview) {
  color: #e6edf3;
}

:deep(.md-editor-preview h1),
:deep(.md-editor-preview h2),
:deep(.md-editor-preview h3),
:deep(.md-editor-preview h4),
:deep(.md-editor-preview h5),
:deep(.md-editor-preview h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

:deep(.md-editor-preview code) {
  background-color: #f6f8fa;
  border-radius: 3px;
  padding: 0.2em 0.4em;
  font-size: 85%;
}

:deep(.md-editor-preview pre) {
  background-color: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  overflow: auto;
}

:deep(.md-editor-preview blockquote) {
  border-left: 4px solid #dfe2e5;
  padding-left: 16px;
  color: #6a737d;
}

:deep(.md-editor-preview table) {
  border-collapse: collapse;
  width: 100%;
}

:deep(.md-editor-preview th),
:deep(.md-editor-preview td) {
  border: 1px solid #dfe2e5;
  padding: 6px 13px;
}

:deep(.md-editor-preview img) {
  max-width: 100%;
}

/* 隐藏图片下拉菜单中的上传图片和裁剪上传选项 */
/* 通过 title 属性精准匹配中文和英文 */
:deep(.md-editor-menu-item[title*="上传"]),
:deep(.md-editor-menu-item[title*="裁剪"]),
:deep(.md-editor-menu-item[title*="Upload"]),
:deep(.md-editor-menu-item[title*="Clip"]) {
  display: none !important;
}

/* 修复代码块复制按钮遮挡问题 */
:deep(.md-editor-code-head),
:deep(.md-editor-copy-button) {
  z-index: 10 !important;
}

/* 调整 emoji 表格样式，使其更大且居中 */
:deep(.emojis li) {
  width: 42px !important;
  height: 42px !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  font-size: 24px !important;
  border-color: #f0f0f0 !important;
  padding: 0 !important;
}

:deep(.emojis) {
  width: auto !important;
  max-width: 440px !important;
}
</style>

<style>
/* 针对 ExportPDF 弹窗的全局样式覆盖 - 强制左对齐 */
.export-pdf-modal .md-editor-preview-wrapper,
.export-pdf-content .md-editor-preview-wrapper,
#export-pdf-preview {
  display: block !important;
  text-align: left !important;
  padding: 0 20px !important;
}

.export-pdf-modal .md-editor-preview,
.export-pdf-content .md-editor-preview {
  text-align: left !important;
  margin: 0 !important;
  width: 100% !important;
  max-width: 100% !important;
}

/* 覆盖模态框主体的对齐方式 */
.export-pdf-modal .md-editor-modal-body {
  text-align: left !important;
  display: block !important;
}

/* 打印时的样式 */
@media print {
  .export-pdf-content .md-editor-preview-wrapper,
  .export-pdf-content .md-editor-preview,
  #export-pdf-preview {
    text-align: left !important;
    margin: 0 !important;
    width: 100% !important;
    max-width: none !important;
  }
  
  /* 确保每一页的内容也是左对齐 */
  .md-editor-preview h1, 
  .md-editor-preview h2, 
  .md-editor-preview h3, 
  .md-editor-preview p,
  .md-editor-preview ul,
  .md-editor-preview ol,
  .md-editor-preview table {
    text-align: left !important;
  }
}
</style>