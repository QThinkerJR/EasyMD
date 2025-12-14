import { reactive, computed } from 'vue'

const messages = {
  'zh-CN': {
    new: '新建',
    newFile: '新建文件',
    open: '打开',
    openFile: '打开文件',
    save: '保存',
    saveAs: '另存为',
    exportHtml: '导出 HTML',
    toggleTheme: '切换主题',
    about: '关于',
    version: '版本',
    untitled: '未命名',
    fileContentEmpty: '文件内容为空',
    openSuccess: '打开成功',
    openFailed: '打开文件失败',
    saveSuccess: '保存成功',
    saveFailed: '保存文件失败',
    onlyMarkdown: '只支持 Markdown 文件',
    upload: '上传',
    clip: '裁剪',
    themeDark: '深色模式',
    themeLight: '浅色模式',
    aboutTitle: '关于 EasyMD',
    aboutDesc: '一个简单、优雅的 Markdown 编辑器',
    githubLink: 'GitHub 仓库',
    dragDropFile: '拖拽文件到此处打开',
    prompt: '提示',
    unsavedWarning: '当前文件尚未保存，是否继续新建？未保存的内容将会丢失。',
    continue: '继续',
    cancel: '取消',
    editorTitle: 'Markdown编辑器',
    exportSuccess: '导出成功',
    exportFailed: '导出HTML失败',
    softwareName: '软件名称',
    author: '作者',
    thirdPartyLibs: '使用的第三方库',
    themeTooltipLight: '切换到亮色模式',
    themeTooltipDark: '切换到暗色模式'
  },
  'en-US': {
    new: 'New',
    newFile: 'New File',
    open: 'Open',
    openFile: 'Open File',
    save: 'Save',
    saveAs: 'Save As',
    exportHtml: 'Export HTML',
    toggleTheme: 'Toggle Theme',
    about: 'About',
    version: 'Version',
    untitled: 'Untitled',
    fileContentEmpty: 'File content is empty',
    openSuccess: 'Opened successfully',
    openFailed: 'Failed to open file',
    saveSuccess: 'Saved successfully',
    saveFailed: 'Failed to save file',
    onlyMarkdown: 'Only Markdown files are supported',
    upload: 'Upload',
    clip: 'Clip',
    themeDark: 'Dark Mode',
    themeLight: 'Light Mode',
    aboutTitle: 'About EasyMD',
    aboutDesc: 'A simple and elegant Markdown editor',
    githubLink: 'GitHub Repository',
    dragDropFile: 'Drag and drop file here to open',
    prompt: 'Prompt',
    unsavedWarning: 'The current file has not been saved. Do you want to continue creating a new one? Unsaved content will be lost.',
    continue: 'Continue',
    cancel: 'Cancel',
    editorTitle: 'Markdown Editor',
    exportSuccess: 'Exported successfully',
    exportFailed: 'Failed to export HTML',
    softwareName: 'Software Name',
    author: 'Author',
    thirdPartyLibs: 'Third-party Libraries',
    themeTooltipLight: 'Switch to Light Mode',
    themeTooltipDark: 'Switch to Dark Mode'
  }
}

const state = reactive({
  locale: 'zh-CN'
})

// 检测语言
const detectLanguage = () => {
  const lang = navigator.language
  // 如果语言以 zh 开头 (zh, zh-CN, zh-TW 等)，使用中文
  // 否则使用英文
  if (lang.startsWith('zh')) {
    state.locale = 'zh-CN'
  } else {
    state.locale = 'en-US'
  }
  console.log('Detected language:', lang, 'Set locale to:', state.locale)
}

// 翻译函数
const t = (key) => {
  return messages[state.locale][key] || key
}

// 获取当前语言
const currentLocale = computed(() => state.locale)

export { state, detectLanguage, t, currentLocale }
