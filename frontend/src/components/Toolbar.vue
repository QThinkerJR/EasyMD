<script setup>
import { computed, ref, onMounted } from 'vue'
import { Button, Space, Tooltip, Dialog } from 'tdesign-vue-next'
import { GetAppVersion } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { t } from '../i18n'

const props = defineProps({
  currentFilePath: {
    type: String,
    default: ''
  },
  isModified: {
    type: Boolean,
    default: false
  },
  theme: {
    type: String,
    default: 'light'
  }
})

const emit = defineEmits([
  'new-file',
  'open-file',
  'save-file',
  'save-file-as',
  'export-html',
  'export-pdf',
  'toggle-theme'
])

const showAbout = ref(false)
const appVersion = ref('')

onMounted(async () => {
  try {
    appVersion.value = await GetAppVersion()
  } catch (e) {
    console.error('Failed to get app version:', e)
    appVersion.value = 'Unknown'
  }
})

// 获取文件名
const fileName = computed(() => {
  if (!props.currentFilePath) return t('untitled')
  return props.currentFilePath.split(/[/\\]/).pop()
})

// 获取保存按钮文本
const saveButtonText = computed(() => {
  return props.currentFilePath ? t('save') : t('save')
})
</script>

<template>
  <div class="toolbar" :class="{ 'dark': theme === 'dark' }">
    <div class="toolbar-left">
      <Space>
        <Tooltip :content="t('newFile')">
          <Button variant="text" @click="$emit('new-file')">
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
              </svg>
            </template>
            {{ t('new') }}
          </Button>
        </Tooltip>
        
        <Tooltip :content="t('openFile')">
          <Button variant="text" @click="$emit('open-file')">
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
              </svg>
            </template>
            {{ t('open') }}
          </Button>
        </Tooltip>
        
        <Tooltip :content="saveButtonText">
          <Button 
            variant="text" 
            @click="$emit('save-file')"
            :disabled="!isModified && !currentFilePath"
          >
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M17 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V7l-4-4zm-5 16c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3zm3-10H5V5h10v4z"/>
              </svg>
            </template>
            {{ saveButtonText }}
          </Button>
        </Tooltip>
        
        <Tooltip :content="t('saveAs')">
          <Button variant="text" @click="$emit('save-file-as')">
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M17 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V7l-4-4zm-5 16c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3zm3-10H5V5h10v4z"/>
              </svg>
            </template>
            {{ t('saveAs') }}
          </Button>
        </Tooltip>
        
        <Tooltip :content="t('exportHtml')">
          <Button variant="text" @click="$emit('export-html')">
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z"/>
              </svg>
            </template>
            {{ t('exportHtml') }}
          </Button>
        </Tooltip>

        <Tooltip :content="t('about')">
          <Button variant="text" @click="showAbout = true">
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M11 7h2v2h-2zm0 4h2v6h-2zm1-9C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"/>
              </svg>
            </template>
            {{ t('about') }}
          </Button>
        </Tooltip>
      </Space>
    </div>
    
    <div class="toolbar-right">
      <Space>
        <Tooltip :content="theme === 'light' ? t('themeTooltipDark') : t('themeTooltipLight')">
          <Button variant="text" @click="$emit('toggle-theme')">
            <template #icon>
              <!-- Sun Icon for Light Mode -->
              <svg v-if="theme === 'light'" viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M12 7c-2.76 0-5 2.24-5 5s2.24 5 5 5 5-2.24 5-5-2.24-5-5-5zM2 13h2c.55 0 1-.45 1-1s-.45-1-1-1H2c-.55 0-1 .45-1 1s.45 1 1 1zm18 0h2c.55 0 1-.45 1-1s-.45-1-1-1h-2c-.55 0-1 .45-1 1s.45 1 1 1zM11 2v2c0 .55.45 1 1 1s1-.45 1-1V2c0-.55-.45-1-1-1s-1 .45-1 1zm0 18v2c0 .55.45 1 1 1s1-.45 1-1v-2c0-.55-.45-1-1-1s-1 .45-1 1zM5.99 4.58c-.39-.39-1.03-.39-1.41 0-.39.39-.39 1.03 0 1.41l1.06 1.06c.39.39 1.03.39 1.41 0s.39-1.03 0-1.41L5.99 4.58zm12.37 12.37c-.39-.39-1.03-.39-1.41 0-.39.39-.39 1.03 0 1.41l1.06 1.06c.39.39 1.03.39 1.41 0 .39-.39.39-1.03 0-1.41l-1.06-1.06zm1.06-10.96c.39-.39.39-1.03 0-1.41-.39-.39-1.03-.39-1.41 0l-1.06 1.06c-.39.39-.39 1.03 0 1.41s1.03.39 1.41 0l1.06-1.06zM7.05 18.36c.39-.39.39-1.03 0-1.41-.39-.39-1.03-.39-1.41 0l-1.06 1.06c-.39.39-.39 1.03 0 1.41s1.03.39 1.41 0l1.06-1.06z"/>
              </svg>
              <!-- Moon Icon for Dark Mode -->
              <svg v-else viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M12 3c-4.97 0-9 4.03-9 9s4.03 9 9 9 9-4.03 9-9c0-.46-.04-.92-.1-1.36-.98 1.37-2.58 2.26-4.4 2.26-2.98 0-5.4-2.42-5.4-5.4 0-1.81.89-3.42 2.26-4.4C12.92 3.04 12.46 3 12 3z"/>
              </svg>
            </template>
          </Button>
        </Tooltip>
      </Space>
      <div class="file-info">
        <span class="file-name">{{ fileName }}</span>
        <span v-if="isModified" class="modified-indicator">*</span>
      </div>
    </div>

    <Dialog
      v-model:visible="showAbout"
      :header="t('aboutTitle')"
      :footer="false"
      width="400px"
      :z-index="99999"
    >
      <div class="about-content">
        <p><strong>{{ t('softwareName') }}:</strong> EasyMD</p>
        <p><strong>{{ t('version') }}:</strong> {{ appVersion }}</p>
        <p><strong>{{ t('author') }}:</strong> QThinker ( <a href="#" @click.prevent="BrowserOpenURL('mailto:ray@qthinker.net')">ray@qthinker.net</a> )</p>
        <div class="libs-section">
          <p><strong>{{ t('thirdPartyLibs') }}:</strong></p>
          <div class="lib-group">
            <p class="lib-title">Go:</p>
            <ul>
              <li>github.com/wailsapp/wails/v2</li>
              <li>github.com/yuin/goldmark</li>
            </ul>
          </div>
          <div class="lib-group">
            <p class="lib-title">Node.js:</p>
            <ul>
              <li>tdesign-vue-next</li>
              <li>md-editor-v3</li>
            </ul>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #e0e0e0;
  min-height: 48px;
}

.toolbar-left {
  display: flex;
  align-items: center;
}

.toolbar-right {
  display: flex;
  align-items: center;
}

.file-info {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #666;
}

.file-name {
  margin-right: 4px;
}

.modified-indicator {
  color: #ff4d4f;
  font-weight: bold;
}

.about-content {
  font-size: 14px;
  line-height: 1.6;
  color: #333;
}

.libs-section {
  margin-top: 16px;
  border-top: 1px solid #eee;
  padding-top: 10px;
}

.lib-group {
  margin-top: 8px;
}

.lib-title {
  font-weight: bold;
  margin-bottom: 4px;
  color: #555;
}

.lib-group ul {
  padding-left: 0;
  margin: 0;
  list-style-type: none;
}

.lib-group li {
  color: #666;
  font-family: monospace;
}

.toolbar.dark {
  background-color: #1e1e1e;
  border-bottom: 1px solid #2d2d2d;
  color: #ffffff;
}

.toolbar.dark .file-info {
  color: #d0d0d0;
}

.toolbar.dark :deep(.t-button--variant-text) {
  color: #ffffff;
}

.toolbar.dark :deep(.t-button--variant-text:hover) {
  background-color: rgba(255, 255, 255, 0.1);
}
</style>