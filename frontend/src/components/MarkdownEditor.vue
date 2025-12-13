<script setup>
import { ref, computed, onMounted } from 'vue'
import { MdEditor, config } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { Emoji, Mark, ExportPDF } from '@vavt/v3-extension'
import '@vavt/v3-extension/lib/asset/style.css'
import MarkExtension from 'markdown-it-mark'

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
  }
})

const emit = defineEmits(['update:modelValue', 'htmlChanged'])

// 编辑器引用
const editorRef = ref(null)

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

// 在组件挂载后，隐藏图片按钮下拉菜单中的上传和裁剪选项
onMounted(() => {
  // 延迟执行，确保DOM已渲染
  const hideImageOptions = () => {
    // 查找所有图片按钮的下拉菜单项
    // md-editor-v3 的下拉菜单项类名为 md-editor-menu-item，图片特定的还有 md-editor-menu-item-image
    const dropdownItems = document.querySelectorAll('.md-editor-menu-item')
    dropdownItems.forEach((item, index) => {
      const text = item.textContent?.trim() || ''
      // 隐藏包含“上传”或“裁剪”的选项
      if (text.includes('上传') || text.includes('裁剪')) {
        item.style.display = 'none'
        console.log(`隐藏选项 ${index}:`, text)
      }
    })
  }
  
  // 多次尝试确保隐藏成功
  setTimeout(hideImageOptions, 100)
  setTimeout(hideImageOptions, 500)
  setTimeout(hideImageOptions, 1000)
})


// 暴露 onHtmlChanged 事件
const onHtmlChanged = (html) => {
  emit('htmlChanged', html)
}

</script>

<template>
  <div class="markdown-editor">
    <!-- 编辑器容器 -->
    <div class="md-editor-container" :class="{ 'loading': isLoading }">
      <MdEditor
        ref="editorRef"
        v-model="content"
        :theme="theme"
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
/* 多种选择器组合，确保全面覆盖 */
:deep(.md-editor .md-editor-dropdown .md-editor-dropdown-item:nth-child(2)),
:deep(.md-editor .md-editor-dropdown .md-editor-dropdown-item:nth-child(3)) {
  display: none !important;
  height: 0 !important;
  padding: 0 !important;
  margin: 0 !important;
  overflow: hidden !important;
  visibility: hidden !important;
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