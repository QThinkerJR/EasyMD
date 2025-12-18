<script setup>
import { ref, computed, nextTick } from 'vue'
import { Dialog, Button, Space, MessagePlugin, Loading } from 'tdesign-vue-next'
import { MdPreview, config } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import Footnote from 'markdown-it-footnote'

// 配置 markdown-it 扩展以支持脚注
config({
  markdownItConfig(md) {
    md.use(Footnote)
  }
})

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible'])

const currentTheme = ref('default')
const previewRef = ref(null)
const copying = ref(false)

const themes = [
  { label: '经典绿色', value: 'default', color: '#07c160' },
  { label: '优雅紫韵', value: 'elegant', color: '#722ed1' },
  { label: '科技蓝调', value: 'tech', color: '#165dff' },
  { label: '暖心橙意', value: 'warm', color: '#fa8c16' },
  { label: '简约黑白', value: 'minimal', color: '#333333' },
  { label: '清新薄荷', value: 'mint', color: '#38b2ac' },
  { label: '浪漫粉樱', value: 'sakura', color: '#f687b3' },
  { label: '商务深蓝', value: 'business', color: '#2c5282' },
  { label: '复古棕褐', value: 'vintage', color: '#8b5a3c' }
]

const handleClose = () => {
  emit('update:visible', false)
}

// Inline styles generator
const copyToWeChat = async () => {
  if (!previewRef.value) return
  
  copying.value = true
  
  try {
    // 1. Clone the container
    // We need to get the inner HTML of the MdPreview
    // MdPreview usually renders into a div with class 'md-editor-preview'
    const previewEl = previewRef.value.querySelector('.md-editor-preview') || previewRef.value
    const clone = previewEl.cloneNode(true)
    
    // Remove code block artifacts (copy buttons, language labels, etc.)
    const artifacts = clone.querySelectorAll('.md-editor-code-head, .md-editor-copy-button, .md-editor-code-lang')
    artifacts.forEach(el => el.remove())
    
    // Handle <details> tags (code folding) if they exist
    const details = clone.querySelectorAll('details')
    details.forEach(detail => {
       const pre = detail.querySelector('pre')
       if (pre) {
          detail.parentNode.replaceChild(pre, detail)
       }
    })

    // Create a temporary container to hold the clone and append to body to calculate styles
    // (Though cloning computed styles is tricky if not attached. 
    // Better approach: Traverse the *original* DOM, get computed styles, apply to *clone*)
    
    const applyStyles = (source, target) => {
      if (source.nodeType !== Node.ELEMENT_NODE) return
      
      const computedStyle = window.getComputedStyle(source)
      let styleString = ''
      
      // We only care about specific properties that matter for WeChat
      // Copying ALL properties creates massive bloat and might break things
      const properties = [
        'color', 'background-color', 'font-size', 'font-family', 'font-weight',
        'line-height', 'text-align', 'border', 'border-radius', 'padding', 
        'margin', 'box-shadow', 'list-style', 'display',
        'text-decoration', 'background', 'border-left', 'border-right', 'border-top', 'border-bottom'
      ]
      
      properties.forEach(prop => {
        const val = computedStyle.getPropertyValue(prop)
        
        // Skip default values to keep HTML clean
        if (!val || val === 'initial' || val === 'none' || val === 'normal' || val === 'auto' || val === '0px' || val === 'rgba(0, 0, 0, 0)') {
           return
        }
        
        // Skip text-align: left/start (let it inherit from wrapper)
        if (prop === 'text-align' && (val === 'left' || val === 'start')) {
            return
        }

        styleString += `${prop}:${val};`
      })
      
      // Special handling for images
      if (source.tagName === 'IMG') {
        target.style.maxWidth = '100%'
        target.style.height = 'auto'
        target.style.display = 'block'
        target.style.margin = '10px auto'
        target.style.boxSizing = 'border-box'
      }

      // Special handling for lists to fix indentation in WeChat
      if (source.tagName === 'UL' || source.tagName === 'OL') {
        target.style.paddingLeft = '20px'
        target.style.listStylePosition = 'outside'
      }

      // Special handling for list items
      if (source.tagName === 'LI') {
        target.style.margin = '5px 0'
        
        // Fix for task lists (checkboxes)
        if (source.classList.contains('md-editor-task-item')) {
           // Logic moved to specific handler below
        }
      }
      
      // Special handling for input checkboxes in task lists
      // We're replacing them with Unicode symbols, so skip input elements
      if (source.tagName === 'INPUT' && source.type === 'checkbox') {
         // Skip input elements as they will be replaced with Unicode symbols
         return
      }
      
      // Handle the label/text after checkbox
      // We can't easily select "text node after checkbox" with CSS selectors in the loop
      // But if the LI has display: flex, it might break in WeChat
      // So let's try a safer approach: Block layout with inline elements
      if (source.tagName === 'LI' && source.classList.contains('md-editor-task-item')) {
         target.style.listStyleType = 'none'
         target.style.margin = '8px 0'
         target.style.paddingLeft = '0'
         target.style.overflow = 'visible' // Allow proper flow
         target.style.display = 'block' // Use block layout
         target.style.verticalAlign = 'baseline'
         target.style.lineHeight = '1.6'
         target.style.fontSize = '16px'
      }
      if (source.tagName === 'BLOCKQUOTE') {
        target.style.boxSizing = 'border-box'
        target.style.maxWidth = '100%'
      }
      
      target.setAttribute('style', styleString)
      
      // Recursively process children
      for (let i = 0; i < source.children.length; i++) {
        if (target.children[i]) {
          applyStyles(source.children[i], target.children[i])
        }
      }
      
      // Skip applying special styles for input checkboxes as they're replaced with Unicode symbols
    }
    
    applyStyles(previewEl, clone)
    
    // Wrap content in a section with max-width 100% to ensure it fits WeChat editor
    const wrapper = document.createElement('section')
    wrapper.style.width = '100%'
    wrapper.style.maxWidth = '100%'
    wrapper.style.boxSizing = 'border-box'
    wrapper.style.textAlign = 'left' // Force left alignment by default
    wrapper.style.overflowX = 'hidden' // Prevent horizontal scroll
    wrapper.style.fontFamily = "-apple-system, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif"
    
    // Move all children of clone to wrapper
    while (clone.firstChild) {
      wrapper.appendChild(clone.firstChild)
    }

    // Process task list items with a simpler, more WeChat-compatible approach
    // Replace input checkboxes with Unicode symbols to avoid WeChat filtering
    const allListItems = wrapper.querySelectorAll('li')
    const taskItems = []
    allListItems.forEach(li => {
      if (li.querySelector('input[type="checkbox"]')) {
        taskItems.push(li)
      }
    })
    
    taskItems.forEach(item => {
      // First, check and fix the parent UL element if it hasn't been processed yet
      const parentUl = item.parentElement
      if (parentUl && (parentUl.tagName === 'UL' || parentUl.tagName === 'OL') && !parentUl.hasAttribute('data-task-list-processed')) {
        // Reset UL styles for task lists to be more compatible with WeChat
        parentUl.style.listStyle = 'none'
        parentUl.style.paddingLeft = '0'
        parentUl.style.marginLeft = '0' // Ensure no extra margin
        parentUl.style.margin = '10px 0'
        parentUl.style.lineHeight = '1.6'
        parentUl.style.fontSize = '16px'
        parentUl.setAttribute('data-task-list-processed', 'true')
        
        // Remove class that might add padding
        parentUl.classList.remove('list-paddingleft-1') // Remove if present
      }
       // Find the checkbox
       const checkbox = item.querySelector('input[type="checkbox"]')
       if (!checkbox) return
       
       // Determine if checkbox is checked
       const isChecked = checkbox.checked
       
       // Collect all content after checkbox
       const childNodes = Array.from(item.childNodes)
       let hasTextContent = false
       
       // Clear the item and rebuild with proper layout
       item.innerHTML = ''
       
       // Create a span for the checkbox symbol
       const checkboxSymbol = document.createElement('span')
       checkboxSymbol.textContent = isChecked ? '☑️ ' : '☐ '
       checkboxSymbol.style.marginRight = '8px'
       checkboxSymbol.style.fontSize = '16px'
       checkboxSymbol.style.lineHeight = '1.6'
       checkboxSymbol.style.display = 'inline-block'
       checkboxSymbol.style.verticalAlign = 'middle'
       
       // Add the checkbox symbol first
       item.appendChild(checkboxSymbol)
       
       // Create a section wrapper for text content with inline-block display
       const textSection = document.createElement('section')
       textSection.style.display = 'inline-block' // Key: make section inline-block to stay on same line as symbol
       textSection.style.verticalAlign = 'middle' // Key: align with symbol
       textSection.style.margin = '0' // Reset default margin
       textSection.style.lineHeight = '1.6'
       textSection.style.fontSize = '16px'
       
       // Add all other content to the section
       childNodes.forEach(child => {
          if (child !== checkbox) {
             hasTextContent = true
             if (child.nodeType === Node.TEXT_NODE) {
                // For text nodes, add directly as text
                textSection.appendChild(document.createTextNode(child.textContent))
             } else if (child.tagName === 'SECTION') {
                // If it's a section, extract its content and add directly
                const sectionChildren = Array.from(child.childNodes)
                sectionChildren.forEach(sectionChild => {
                   if (sectionChild.nodeType === Node.TEXT_NODE) {
                      textSection.appendChild(document.createTextNode(sectionChild.textContent))
                   } else {
                      // Clone and ensure inline display
                      const clonedChild = sectionChild.cloneNode(true)
                      if (clonedChild.style) {
                         clonedChild.style.display = 'inline'
                         clonedChild.style.verticalAlign = 'middle'
                      }
                      textSection.appendChild(clonedChild)
                   }
                })
             } else if (child.tagName === 'UL' || child.tagName === 'OL') {
                // For nested lists, keep them as block elements
                child.style.marginLeft = '24px'
                child.style.marginTop = '5px'
                child.style.marginBottom = '5px'
                textSection.appendChild(child)
             } else {
                // For other element nodes, clone and append
                const clonedChild = child.cloneNode(true)
                // Ensure inline display for all elements
                if (clonedChild.style) {
                   clonedChild.style.display = 'inline'
                   clonedChild.style.verticalAlign = 'middle'
                }
                textSection.appendChild(clonedChild)
             }
          }
       })
       
       // If no text content, add a space to maintain layout
       if (!hasTextContent) {
          textSection.appendChild(document.createTextNode(' '))
       }
       
       // Add the text section to the item
       item.appendChild(textSection)
       
       // Reset LI styles to be more compatible with WeChat
       item.style.listStyle = 'none'
       item.style.padding = '0'
       item.style.margin = '8px 0'
       item.style.lineHeight = '1.6'
       item.style.fontSize = '16px'
       item.style.display = 'block'
    })

    // Process footnotes to ensure they render correctly in WeChat
    const footnotes = wrapper.querySelectorAll('.footnotes')
    footnotes.forEach(footnote => {
      // Style the footnotes section
      footnote.style.marginTop = '30px'
      footnote.style.paddingTop = '15px'
      footnote.style.borderTop = '1px solid #eee'
      footnote.style.fontSize = '14px'
      footnote.style.color = '#666'
      
      // Style the footnote title
      const title = footnote.querySelector('h2')
      if (title) {
        title.style.fontSize = '16px'
        title.style.marginBottom = '10px'
        title.style.color = '#333'
        title.style.borderBottom = '1px solid #eee'
        title.style.paddingBottom = '5px'
      }
      
      // Style individual footnote items
      const footnoteItems = footnote.querySelectorAll('ol > li')
      footnoteItems.forEach(item => {
        item.style.marginBottom = '8px'
        item.style.lineHeight = '1.6'
        
        // Style the footnote back reference
        const backRef = item.querySelector('a[href^="#fnref"]')
        if (backRef) {
          backRef.style.fontSize = '12px'
          backRef.style.marginLeft = '5px'
          backRef.style.textDecoration = 'none'
          backRef.style.color = '#999'
        }
      })
    })
    
    // Process footnote references in the text
    const footnoteRefs = wrapper.querySelectorAll('a[href^="#fn"]')
    footnoteRefs.forEach(ref => {
      ref.style.fontSize = '12px'
      ref.style.verticalAlign = 'super'
      ref.style.margin = '0 2px'
      ref.style.textDecoration = 'none'
      ref.style.color = '#07c160' // Use theme color
      ref.style.fontWeight = 'bold'
    })

    // Process code blocks to add Mac-style header
    const codeBlocks = wrapper.querySelectorAll('.md-editor-code, pre')
    codeBlocks.forEach(block => {
      // Skip if already processed or inside another processed block
      if (block.closest('.mac-window-container')) return
      
      // If it's a pre inside md-editor-code, we want the wrapper
      if (!block.parentNode || (block.tagName === 'PRE' && block.parentNode.classList.contains('md-editor-code'))) {
         return // handled by parent loop
      }

      const pre = block.tagName === 'PRE' ? block : block.querySelector('pre')
      if (!pre) return

      // Determine background color based on theme
      // Default to dark background for all themes as requested
      const isDark = true
      const bgColor = '#282c34'
      const headerColor = '#21252b'
      const textColor = '#abb2bf'
      
      // Create Mac Window Container
      // Use SECTION tag for better WeChat compatibility
      const container = document.createElement('section')
      container.className = 'mac-window-container'
      container.style.borderRadius = '6px'
      container.style.boxShadow = '0 2px 6px rgba(0,0,0,0.1)'
      container.style.overflow = 'hidden'
      container.style.margin = '15px 0'
      container.style.backgroundColor = bgColor
      container.style.display = 'block' // Ensure block display
      
      // Create Mac Header
      const header = document.createElement('div')
      header.style.background = headerColor
      header.style.padding = '8px 12px'
      header.style.display = 'flex'
      header.style.alignItems = 'center'
      header.style.borderRadius = '6px 6px 0 0'
      header.style.lineHeight = '1' // Fix height issues
      
      // Create Dots
      const createDot = (color) => {
        const dot = document.createElement('div')
        dot.style.width = '12px'
        dot.style.height = '12px'
        dot.style.borderRadius = '50%'
        dot.style.backgroundColor = color
        dot.style.marginRight = '8px'
        dot.style.display = 'inline-block' // Ensure inline-block
        return dot
      }
      
      header.appendChild(createDot('#ff5f56')) // Red
      header.appendChild(createDot('#ffbd2e')) // Yellow
      header.appendChild(createDot('#27c93f')) // Green
      
      // Assemble
      container.appendChild(header)
      
      // Clone pre to put inside container
      const preClone = pre.cloneNode(true)
      
      // Remove any background color from children to prevent red/pink lines
      const allChildren = preClone.querySelectorAll('*')
      allChildren.forEach(child => {
         child.style.backgroundColor = 'transparent'
         child.style.backgroundImage = 'none'
         
         // Force inline syntax highlighting colors
         const classList = child.classList
         if (classList.length > 0) {
            // ... (existing highlight logic) ...
            if (classList.contains('hljs-comment') || classList.contains('hljs-quote')) {
               child.style.color = '#5c6370'
               child.style.fontStyle = 'italic'
            } else if (classList.contains('hljs-doctag') || classList.contains('hljs-keyword') || classList.contains('hljs-formula')) {
               child.style.color = '#c678dd'
            } else if (classList.contains('hljs-section') || classList.contains('hljs-name') || classList.contains('hljs-selector-tag') || classList.contains('hljs-deletion') || classList.contains('hljs-subst')) {
               child.style.color = '#e06c75'
            } else if (classList.contains('hljs-literal')) {
               child.style.color = '#56b6c2'
            } else if (classList.contains('hljs-string') || classList.contains('hljs-regexp') || classList.contains('hljs-addition') || classList.contains('hljs-attribute') || classList.contains('hljs-meta')) {
               child.style.color = '#98c379'
            } else if (classList.contains('hljs-attr') || classList.contains('hljs-variable') || classList.contains('hljs-template-variable') || classList.contains('hljs-type') || classList.contains('hljs-selector-class') || classList.contains('hljs-selector-attr') || classList.contains('hljs-selector-pseudo') || classList.contains('hljs-number')) {
               child.style.color = '#d19a66'
            } else if (classList.contains('hljs-symbol') || classList.contains('hljs-bullet') || classList.contains('hljs-link') || classList.contains('hljs-selector-id') || classList.contains('hljs-title')) {
               child.style.color = '#61aeee'
            } else if (classList.contains('hljs-built_in') || classList.contains('hljs-title') && classList.contains('class_') || classList.contains('hljs-class')) {
               child.style.color = '#e6c07b'
            } else if (classList.contains('hljs-emphasis')) {
               child.style.fontStyle = 'italic'
            } else if (classList.contains('hljs-strong')) {
               child.style.fontWeight = 'bold'
            }
         }
      })
      
      // Traverse text nodes and replace spaces with NBSP to prevent WeChat from eating them
      // But we need to be careful not to introduce extra spaces at the beginning if they are just layout formatting
      const replaceSpaces = (node) => {
        if (node.nodeType === Node.TEXT_NODE) {
           // Replace space with Non-Breaking Space (0xA0)
           // But avoid replacing newlines
           let val = node.nodeValue
           // If the node contains ONLY newlines/spaces and is followed by a block, it might be structural
           // However, inside PRE, all whitespace matters.
           
           // Strategy: Replace ALL regular spaces with NBSP, but preserve newlines
           val = val.replace(/ /g, '\u00a0')
           node.nodeValue = val
        } else if (node.nodeType === Node.ELEMENT_NODE) {
           for (let i = 0; i < node.childNodes.length; i++) {
              replaceSpaces(node.childNodes[i])
           }
        }
      }
      
      // First, trim the initial and trailing newlines/spaces from the code block itself
      // We perform a deep search for the first text node and trim it
      const trimLeadingWhitespace = (node) => {
         if (node.nodeType === Node.TEXT_NODE) {
            // Found the first text node
            if (node.nodeValue.trim().length === 0) {
               // If it's pure whitespace, remove it entirely and continue search
               const parent = node.parentNode
               parent.removeChild(node)
               return false // Continue searching
            } else {
               // Trim start of this node
               node.nodeValue = node.nodeValue.replace(/^\s+/, '')
               return true // Done
            }
         } else if (node.nodeType === Node.ELEMENT_NODE) {
            // Search children
            // We need to iterate backwards-compatible loop because we might remove children
            const children = Array.from(node.childNodes)
            for (const child of children) {
               if (trimLeadingWhitespace(child)) {
                  return true // Found and trimmed
               }
            }
         }
         return false
      }
      trimLeadingWhitespace(preClone)
      
      // Also trim trailing whitespace (simpler, usually at the end of pre)
      // For trailing, innerHTML replace is usually safe enough as it doesn't affect structure as much
      // but let's be consistent and safe.
      let html = preClone.innerHTML
      html = html.replace(/\s+$/, '')
      preClone.innerHTML = html

      // Re-query children after innerHTML update
      const newChildren = preClone.querySelectorAll('*')
      newChildren.forEach(child => {
         // Re-apply styles (same logic as above)
         child.style.backgroundColor = 'transparent'
         child.style.backgroundImage = 'none'
         
         const classList = child.classList
         if (classList.length > 0) {
            if (classList.contains('hljs-comment') || classList.contains('hljs-quote')) {
               child.style.color = '#5c6370'
               child.style.fontStyle = 'italic'
            } else if (classList.contains('hljs-doctag') || classList.contains('hljs-keyword') || classList.contains('hljs-formula')) {
               child.style.color = '#c678dd'
            } else if (classList.contains('hljs-section') || classList.contains('hljs-name') || classList.contains('hljs-selector-tag') || classList.contains('hljs-deletion') || classList.contains('hljs-subst')) {
               child.style.color = '#e06c75'
            } else if (classList.contains('hljs-literal')) {
               child.style.color = '#56b6c2'
            } else if (classList.contains('hljs-string') || classList.contains('hljs-regexp') || classList.contains('hljs-addition') || classList.contains('hljs-attribute') || classList.contains('hljs-meta')) {
               child.style.color = '#98c379'
            } else if (classList.contains('hljs-attr') || classList.contains('hljs-variable') || classList.contains('hljs-template-variable') || classList.contains('hljs-type') || classList.contains('hljs-selector-class') || classList.contains('hljs-selector-attr') || classList.contains('hljs-selector-pseudo') || classList.contains('hljs-number')) {
               child.style.color = '#d19a66'
            } else if (classList.contains('hljs-symbol') || classList.contains('hljs-bullet') || classList.contains('hljs-link') || classList.contains('hljs-selector-id') || classList.contains('hljs-title')) {
               child.style.color = '#61aeee'
            } else if (classList.contains('hljs-built_in') || classList.contains('hljs-title') && classList.contains('class_') || classList.contains('hljs-class')) {
               child.style.color = '#e6c07b'
            } else if (classList.contains('hljs-emphasis')) {
               child.style.fontStyle = 'italic'
            } else if (classList.contains('hljs-strong')) {
               child.style.fontWeight = 'bold'
            }
         }
      })
      
      replaceSpaces(preClone)

      preClone.style.margin = '0'
      preClone.style.padding = '12px 16px'
      // Explicitly set background on PRE as well, in case container bg is stripped
      preClone.style.backgroundColor = bgColor 
      preClone.style.border = 'none'
      preClone.style.borderRadius = '0 0 6px 6px'
      preClone.style.color = textColor
      preClone.style.fontFamily = 'Consolas, Monaco, "Andale Mono", "Ubuntu Mono", monospace'
      preClone.style.fontSize = '14px'
      preClone.style.lineHeight = '1.6'
      // Use 'pre' instead of 'pre-wrap' to ensure spaces are respected and allow scrolling
      preClone.style.whiteSpace = 'pre' 
      preClone.style.overflowX = 'auto'
      preClone.style.display = 'block' // Ensure block display
      preClone.style.textIndent = '0' // Force reset text-indent
      preClone.style.paddingLeft = '16px' // Ensure left padding matches style

      // Ensure code element inside pre also has correct style
      const codeEl = preClone.querySelector('code')
      if (codeEl) {
        codeEl.style.whiteSpace = 'pre'
        codeEl.style.fontFamily = 'inherit'
        codeEl.style.color = 'inherit'
        codeEl.style.textIndent = '0' // Force reset text-indent
        codeEl.style.paddingLeft = '0' // Code element shouldn't have padding
        codeEl.style.marginLeft = '0'
      }
      
      container.appendChild(preClone)
      
      // Replace original block with new container
      // If block is .md-editor-code, replace it. If block is pre, replace it.
      if (block.tagName === 'PRE') {
        block.parentNode.replaceChild(container, block)
      } else {
        block.innerHTML = ''
        block.appendChild(container)
      }
    })

    // Copy to clipboard logic with fallback
    const doCopy = async () => {
      try {
        // Try Modern Clipboard API first
        const clipboardItem = new ClipboardItem({
          'text/html': new Blob([wrapper.outerHTML], { type: 'text/html' }),
          'text/plain': new Blob([wrapper.innerText], { type: 'text/plain' })
        })
        
        await navigator.clipboard.write([clipboardItem])
        MessagePlugin.success('复制成功，请直接粘贴到微信公众号后台')
      } catch (err) {
        console.warn('Clipboard API failed, trying execCommand fallback...', err)
        
        // Fallback: document.execCommand('copy')
        try {
          const container = document.createElement('div')
          // We need to clone wrapper because appending it will move it from memory to DOM
          // and if we needed to retry or use it elsewhere, it might be an issue. 
          // But here it's fine as this is the last step.
          container.appendChild(wrapper)
          
          container.style.position = 'fixed'
          container.style.left = '-9999px'
          container.style.top = '0'
          container.style.opacity = '0'
          container.style.pointerEvents = 'none'
          container.style.zIndex = '-1'
          
          document.body.appendChild(container)
          
          const selection = window.getSelection()
          const range = document.createRange()
          
          // Select the wrapper element itself
          range.selectNode(wrapper)
          
          selection.removeAllRanges()
          selection.addRange(range)
          
          const successful = document.execCommand('copy')
          
          document.body.removeChild(container)
          selection.removeAllRanges()
          
          if (successful) {
            MessagePlugin.success('复制成功，请直接粘贴到微信公众号后台')
          } else {
            throw new Error('execCommand returned false')
          }
        } catch (fallbackErr) {
          console.error('Fallback copy failed', fallbackErr)
          throw fallbackErr // Re-throw to trigger outer catch
        }
      }
    }

    await doCopy()
    
  } catch (err) {
    console.error('Copy process failed', err)
    MessagePlugin.error('复制失败，请重试')
  } finally {
    copying.value = false
  }
}

</script>

<template>
  <Dialog
    :visible="visible"
    :footer="false"
    header="微信公众号文章主题预览"
    width="90%"
    top="5vh"
    @close="handleClose"
    class="wechat-preview-dialog"
  >
    <div class="wechat-container">
      <!-- Sidebar -->
      <div class="theme-sidebar">
        <div class="theme-list">
          <div 
            v-for="theme in themes" 
            :key="theme.value"
            class="theme-item"
            :class="{ active: currentTheme === theme.value }"
            @click="currentTheme = theme.value"
          >
            <div class="color-dot" :style="{ backgroundColor: theme.color }"></div>
            <span>{{ theme.label }}</span>
          </div>
        </div>
        
        <div class="action-area">
          <Button theme="primary" block size="large" @click="copyToWeChat" :loading="copying">
            <template #icon>
              <svg viewBox="0 0 24 24" width="16" height="16" style="margin-right: 8px;">
                <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
              </svg>
            </template>
            一键复制
          </Button>
          <p class="tip">复制后直接粘贴到公众号后台编辑器</p>
        </div>
      </div>
      
      <!-- Preview Area -->
      <div class="preview-area-wrapper">
        <div class="preview-mobile-frame">
          <div class="preview-content" ref="previewRef" :class="`wechat-theme-${currentTheme}`">
            <MdPreview 
              v-if="visible"
              :modelValue="modelValue" 
              previewTheme="default" 
              codeTheme="atom" 
              :codeFoldable="false"
              editorId="wechat-preview"
            />
            
            <!-- Footer Signature (Optional) -->
            <div class="wechat-footer">
              <p>Created with EasyMD</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Dialog>
</template>

<style scoped>
.wechat-container {
  display: flex;
  height: 75vh;
  gap: 20px;
}

.theme-sidebar {
  width: 250px;
  border-right: 1px solid #eee;
  padding-right: 20px;
  display: flex;
  flex-direction: column;
}

.theme-sidebar h3 {
  margin-bottom: 20px;
  color: #333;
}

.theme-list {
  flex: 1;
  overflow-y: auto;
}

.theme-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  border-radius: 8px;
  margin-bottom: 8px;
  transition: all 0.2s;
}

.theme-item:hover {
  background-color: #f5f5f5;
}

.theme-item.active {
  background-color: #e6f7ff;
  color: #1677ff;
}

.color-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  margin-right: 12px;
}

.action-area {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.tip {
  font-size: 12px;
  color: #999;
  text-align: center;
  margin-top: 8px;
}

.preview-area-wrapper {
  flex: 1;
  background-color: #f0f2f5;
  display: flex;
  justify-content: center;
  padding: 20px;
  overflow-y: auto;
  border-radius: 8px;
}

.preview-mobile-frame {
  width: 500px; /* Increased from 375px for better visibility */
  min-height: 100%;
  background-color: #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 20px;
  overflow-x: hidden;
}

.wechat-footer {
  margin-top: 40px;
  text-align: center;
  font-size: 12px;
  color: #ccc;
  border-top: 1px dashed #eee;
  padding-top: 20px;
}
</style>

<!-- Global Styles for WeChat Themes -->
<!-- These need to be not scoped so they affect the MdPreview content -->
<style>
/* Reset some md-editor defaults to fit WeChat */
.preview-content {
  text-align: left !important;
}

.preview-content .md-editor-preview {
  padding: 0;
  font-size: 16px;
  line-height: 1.8;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, "Helvetica Neue", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei UI", "Microsoft YaHei", Arial, sans-serif;
}

.preview-content .md-editor-preview img {
  max-width: 100%;
  border-radius: 4px;
  margin: 10px 0;
  display: block;
}

/* ================== Default Theme (Green) ================== */
.wechat-theme-default .md-editor-preview h1 {
  font-size: 22px;
  font-weight: bold;
  border-bottom: 2px solid #07c160;
  padding-bottom: 10px;
  margin: 20px 0 15px;
  color: #333;
}

.wechat-theme-default .md-editor-preview h2 {
  font-size: 18px;
  font-weight: bold;
  background: rgba(7, 193, 96, 0.1);
  border-left: 5px solid #07c160;
  padding: 10px;
  margin: 20px 0 15px;
  border-radius: 0 4px 4px 0;
  color: #333;
  display: table;
  padding-right: 10px;
}

.wechat-theme-default .md-editor-preview h3 {
  font-size: 16px;
  font-weight: bold;
  color: #07c160;
  margin: 15px 0 10px;
  padding-left: 10px;
  border-left: 3px solid #07c160;
}

.wechat-theme-default .md-editor-preview blockquote {
  border-left: 4px solid #07c160;
  background-color: #f8f9fa;
  color: #666;
  padding: 15px;
  margin: 15px 0;
  border-radius: 4px;
}

.wechat-theme-default .md-editor-preview ul li::marker {
  color: #07c160;
}

.wechat-theme-default .md-editor-preview a {
  color: #07c160;
  text-decoration: none;
  border-bottom: 1px dashed #07c160;
}

/* Footnote styles for default theme */
.wechat-theme-default .md-editor-preview .footnotes {
  margin-top: 30px;
  padding-top: 15px;
  border-top: 1px solid #eee;
  font-size: 14px;
  color: #666;
}

.wechat-theme-default .md-editor-preview .footnotes h2 {
  font-size: 16px;
  margin-bottom: 10px;
  color: #333;
  border-bottom: 1px solid #eee;
  padding-bottom: 5px;
  border-left: none;
  background: none;
  padding-left: 0;
}

.wechat-theme-default .md-editor-preview .footnotes ol {
  padding-left: 20px;
  margin: 0;
}

.wechat-theme-default .md-editor-preview .footnotes li {
  margin-bottom: 8px;
  line-height: 1.6;
}

.wechat-theme-default .md-editor-preview .footnotes a {
  color: #07c160;
  text-decoration: none;
  border-bottom: none;
  font-size: 12px;
  margin-left: 5px;
}

.wechat-theme-default .md-editor-preview a[href^="#fn"] {
  font-size: 12px;
  vertical-align: super;
  margin: 0 2px;
  text-decoration: none;
  color: #07c160;
  font-weight: bold;
  border-bottom: none;
}

/* ================== Elegant Theme (Purple) ================== */
.wechat-theme-elegant .md-editor-preview {
  font-family: "Optima", "PingFang SC", serif;
}

.wechat-theme-elegant .md-editor-preview h1 {
  font-size: 24px;
  text-align: center;
  color: #722ed1;
  margin: 30px 0 20px;
  font-weight: normal;
}

.wechat-theme-elegant .md-editor-preview h2 {
  font-size: 18px;
  text-align: center;
  color: #fff;
  background-color: #722ed1;
  display: table;
  margin: 30px auto 20px;
  padding: 5px 20px;
  border-radius: 20px;
  font-weight: normal;
}

.wechat-theme-elegant .md-editor-preview h3 {
  font-size: 16px;
  color: #722ed1;
  border-bottom: 1px solid #722ed1;
  display: inline-block;
  padding-bottom: 5px;
  margin: 20px 0 10px;
}

.wechat-theme-elegant .md-editor-preview blockquote {
  border: none;
  background-color: #f9f0ff;
  color: #722ed1;
  padding: 20px;
  margin: 20px 0;
  border-radius: 8px;
  text-align: center;
  font-style: italic;
}

.wechat-theme-elegant .md-editor-preview strong {
  color: #722ed1;
}

/* ================== Tech Theme (Blue) ================== */
.wechat-theme-tech .md-editor-preview h1 {
  font-size: 22px;
  color: #165dff;
  border-bottom: 1px solid #e5e6eb;
  padding-bottom: 15px;
  margin-bottom: 20px;
}

.wechat-theme-tech .md-editor-preview h2 {
  font-size: 18px;
  color: #165dff;
  border-left: 4px solid #165dff;
  padding-left: 12px;
  margin: 25px 0 15px;
  background: linear-gradient(90deg, rgba(22, 93, 255, 0.1) 0%, rgba(255, 255, 255, 0) 100%);
  padding-top: 8px;
  padding-bottom: 8px;
  display: table;
  padding-right: 12px;
}

.wechat-theme-tech .md-editor-preview h3 {
  font-size: 16px;
  color: #165dff;
  margin: 20px 0 10px;
}

.wechat-theme-tech .md-editor-preview blockquote {
  border-left: 4px solid #165dff;
  background-color: #e8f3ff;
  color: #555;
  padding: 10px 15px;
  border-radius: 0 8px 8px 0;
}

.wechat-theme-tech .md-editor-preview code {
  color: #165dff;
  background-color: #e8f3ff;
  padding: 2px 6px;
  border-radius: 4px;
}

/* Footnote styles for tech theme */
.wechat-theme-tech .md-editor-preview .footnotes {
  margin-top: 30px;
  padding-top: 15px;
  border-top: 1px solid #e5e6eb;
  font-size: 14px;
  color: #666;
}

.wechat-theme-tech .md-editor-preview .footnotes h2 {
  font-size: 16px;
  margin-bottom: 10px;
  color: #165dff;
  border-bottom: 1px solid #e5e6eb;
  padding-bottom: 5px;
  border-left: none;
  background: none;
  padding-left: 0;
}

.wechat-theme-tech .md-editor-preview .footnotes ol {
  padding-left: 20px;
  margin: 0;
}

.wechat-theme-tech .md-editor-preview .footnotes li {
  margin-bottom: 8px;
  line-height: 1.6;
}

.wechat-theme-tech .md-editor-preview .footnotes a {
  color: #165dff;
  text-decoration: none;
  border-bottom: none;
  font-size: 12px;
  margin-left: 5px;
}

.wechat-theme-tech .md-editor-preview a[href^="#fn"] {
  font-size: 12px;
  vertical-align: super;
  margin: 0 2px;
  text-decoration: none;
  color: #165dff;
  font-weight: bold;
  border-bottom: none;
}

/* ================== Warm Theme (Orange) ================== */
.wechat-theme-warm .md-editor-preview h1 {
  font-size: 22px;
  color: #fa8c16;
  text-align: center;
  background-image: linear-gradient(to right, transparent, rgba(250, 140, 22, 0.2), transparent);
  padding: 10px 0;
  margin: 20px 0;
}

.wechat-theme-warm .md-editor-preview h2 {
  font-size: 18px;
  color: #fff;
  background-color: #fa8c16;
  display: inline-block;
  padding: 5px 15px;
  border-radius: 10px 10px 10px 0;
  margin: 20px 0 15px;
  box-shadow: 2px 2px 0 rgba(250, 140, 22, 0.3);
}

.wechat-theme-warm .md-editor-preview h3 {
  font-size: 16px;
  color: #fa8c16;
  border-bottom: 2px solid #fa8c16;
  display: inline-block;
  margin: 15px 0 10px;
}

.wechat-theme-warm .md-editor-preview blockquote {
  border: 1px dashed #fa8c16;
  background-color: #fff7e6;
  color: #8c5a2b;
  padding: 15px;
  border-radius: 8px;
  margin: 15px 0;
}

.wechat-theme-warm .md-editor-preview strong {
  color: #d46b08;
  border-bottom: 2px solid rgba(250, 140, 22, 0.3);
}

/* Footnote styles for warm theme */
.wechat-theme-warm .md-editor-preview .footnotes {
  margin-top: 30px;
  padding-top: 15px;
  border-top: 1px dashed #fa8c16;
  font-size: 14px;
  color: #666;
}

.wechat-theme-warm .md-editor-preview .footnotes h2 {
  font-size: 16px;
  margin-bottom: 10px;
  color: #fa8c16;
  border-bottom: 1px dashed #fa8c16;
  padding-bottom: 5px;
  display: block;
  padding-left: 0;
  background: none;
  border-radius: 0;
  box-shadow: none;
}

.wechat-theme-warm .md-editor-preview .footnotes ol {
  padding-left: 20px;
  margin: 0;
}

.wechat-theme-warm .md-editor-preview .footnotes li {
  margin-bottom: 8px;
  line-height: 1.6;
}

.wechat-theme-warm .md-editor-preview .footnotes a {
  color: #fa8c16;
  text-decoration: none;
  border-bottom: none;
  font-size: 12px;
  margin-left: 5px;
}

.wechat-theme-warm .md-editor-preview a[href^="#fn"] {
  font-size: 12px;
  vertical-align: super;
  margin: 0 2px;
  text-decoration: none;
  color: #fa8c16;
  font-weight: bold;
  border-bottom: none;
}

/* ================== Minimal Theme (Black & White) ================== */
.wechat-theme-minimal .md-editor-preview {
  font-family: 'Helvetica Neue', Arial, sans-serif;
  color: #333;
  line-height: 1.8;
}

.wechat-theme-minimal .md-editor-preview h1 {
  font-size: 26px;
  font-weight: 300;
  text-align: center;
  color: #000;
  padding: 20px 0;
  margin: 30px 0 20px;
  letter-spacing: 3px;
  border-top: 3px solid #000;
  border-bottom: 1px solid #000;
}

.wechat-theme-minimal .md-editor-preview h2 {
  font-size: 22px;
  font-weight: 400;
  color: #000;
  margin: 25px auto 15px;
  padding: 12px 20px;
  border-left: 5px solid #000;
  border-right: 5px solid #000;
  text-align: center;
  background: linear-gradient(90deg, rgba(0,0,0,0.05) 0%, rgba(0,0,0,0.02) 50%, rgba(0,0,0,0.05) 100%);
  outline: 1px solid #ddd;
  display: table;
}

.wechat-theme-minimal .md-editor-preview h3 {
  font-size: 18px;
  font-weight: 400;
  color: #333;
  margin: 20px 0 10px;
  padding-left: 20px;
  padding-bottom: 5px;
  border-bottom: 1px dashed #999;
  border-left: 3px solid #000;
}

.wechat-theme-minimal .md-editor-preview blockquote {
  border: none;
  background-color: #f8f8f8;
  color: #666;
  padding: 20px;
  margin: 20px 0;
  font-style: italic;
  position: relative;
}

.wechat-theme-minimal .md-editor-preview p {
  margin-bottom: 16px;
  text-align: justify;
}

.wechat-theme-minimal .md-editor-preview strong {
  font-weight: 600;
  color: #000;
}

.wechat-theme-minimal .md-editor-preview a {
  color: #000;
  text-decoration: underline;
}

/* ================== Mint Theme (Fresh Mint Green) ================== */
.wechat-theme-mint .md-editor-preview {
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
  color: #2d3748;
}

.wechat-theme-mint .md-editor-preview h1 {
  font-size: 24px;
  font-weight: 600;
  color: #38b2ac;
  text-align: center;
  margin: 25px 0 20px;
  padding: 20px 0;
  border-top: 4px solid #38b2ac;
  border-bottom: 4px solid #81e6d9;
  border-radius: 2px;
}

.wechat-theme-mint .md-editor-preview h2 {
  font-size: 20px;
  font-weight: 500;
  color: #fff;
  background: linear-gradient(135deg, #38b2ac, #4fd1c5);
  padding: 12px 25px;
  margin: 20px 0 15px;
  border-radius: 30px 5px 30px 5px;
  display: inline-block;
  box-shadow: 0 6px 15px rgba(56, 178, 172, 0.3);
  transform: rotate(-1deg);
}

.wechat-theme-mint .md-editor-preview h3 {
  font-size: 18px;
  font-weight: 500;
  color: #38b2ac;
  margin: 18px 0 10px;
  padding-left: 20px;
  border-left: 5px solid #81e6d9;
  padding-bottom: 5px;
  border-bottom: 2px dotted #81e6d9;
}

.wechat-theme-mint .md-editor-preview blockquote {
  border-left: 4px solid #38b2ac;
  background-color: #e6fffa;
  color: #234e52;
  padding: 15px 20px;
  margin: 15px 0;
  border-radius: 0 8px 8px 0;
}

.wechat-theme-mint .md-editor-preview strong {
  color: #38b2ac;
  font-weight: 600;
}

.wechat-theme-mint .md-editor-preview a {
  color: #38b2ac;
  text-decoration: none;
  border-bottom: 1px dotted #38b2ac;
}

/* ================== Sakura Theme (Romantic Pink) ================== */
.wechat-theme-sakura .md-editor-preview {
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
  color: #4a5568;
  line-height: 1.7;
}

.wechat-theme-sakura .md-editor-preview h1 {
  font-size: 26px;
  font-weight: 600;
  text-align: center;
  color: #d53f8c;
  margin: 25px 0 20px;
  padding: 25px 0;
  background: linear-gradient(135deg, rgba(246, 135, 179, 0.1) 0%, rgba(251, 182, 206, 0.1) 100%);
  border-radius: 15px;
  border: 2px dashed #fbb6ce;
}

.wechat-theme-sakura .md-editor-preview h2 {
  font-size: 20px;
  font-weight: 500;
  color: #fff;
  background: linear-gradient(45deg, #f687b3, #fbb6ce);
  padding: 12px 25px;
  margin: 20px 0 15px;
  border-radius: 25px 8px 25px 8px;
  display: inline-block;
  box-shadow: 0 5px 12px rgba(246, 135, 179, 0.4);
  transform: rotate(-0.5deg);
}

.wechat-theme-sakura .md-editor-preview h3 {
  font-size: 18px;
  font-weight: 500;
  color: #d53f8c;
  margin: 18px 0 10px;
  padding: 8px 15px;
  background: linear-gradient(90deg, rgba(246, 135, 179, 0.1) 0%, transparent 100%);
  border-left: 4px solid #f687b3;
  border-bottom: 2px dotted #fbb6ce;
}

.wechat-theme-sakura .md-editor-preview blockquote {
  border: none;
  background: linear-gradient(135deg, #fff5f7, #fed7e2);
  color: #702459;
  padding: 20px;
  margin: 20px 0;
  border-radius: 15px;
  position: relative;
  box-shadow: 0 2px 8px rgba(246, 135, 179, 0.15);
}

.wechat-theme-sakura .md-editor-preview blockquote::before {
  content: '💕';
  position: absolute;
  top: -10px;
  left: 20px;
  background: #fff;
  border-radius: 50%;
  padding: 5px;
  font-size: 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.wechat-theme-sakura .md-editor-preview strong {
  color: #d53f8c;
  font-weight: 600;
}

.wechat-theme-sakura .md-editor-preview a {
  color: #d53f8c;
  text-decoration: none;
  border-bottom: 1px solid #fbb6ce;
}

/* ================== Business Theme (Professional Deep Blue) ================== */
.wechat-theme-business .md-editor-preview {
  font-family: 'Helvetica Neue', 'PingFang SC', sans-serif;
  color: #2d3748;
  line-height: 1.6;
}

.wechat-theme-business .md-editor-preview h1 {
  font-size: 28px;
  font-weight: 600;
  color: #2c5282;
  text-align: center;
  margin: 25px 0 20px;
  position: relative;
  padding: 25px 0;
  background: linear-gradient(135deg, rgba(44, 82, 130, 0.05) 0%, rgba(66, 153, 225, 0.05) 100%);
  border-top: 3px solid #2c5282;
  border-bottom: 3px solid #2c5282;
}

.wechat-theme-business .md-editor-preview h2 {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #2c5282, #4299e1);
  padding: 12px 25px;
  margin: 20px 0 15px;
  border-left: 6px solid #2b6cb0;
  box-shadow: 0 4px 10px rgba(44, 82, 130, 0.3);
  display: table;
  padding-right: 25px;
}

.wechat-theme-business .md-editor-preview h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2c5282;
  margin: 18px 0 10px;
  padding-left: 20px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e2e8f0;
  border-left: 3px solid #4299e1;
}

.wechat-theme-business .md-editor-preview blockquote {
  border-left: 4px solid #2c5282;
  background-color: #ebf8ff;
  color: #2a4e7c;
  padding: 15px 20px;
  margin: 15px 0;
  border-radius: 0 4px 4px 0;
}

.wechat-theme-business .md-editor-preview strong {
  color: #2c5282;
  font-weight: 600;
}

.wechat-theme-business .md-editor-preview a {
  color: #2c5282;
  text-decoration: none;
  border-bottom: 1px solid #4299e1;
}

.wechat-theme-business .md-editor-preview table {
  border-collapse: collapse;
  width: 100%;
  margin: 15px 0;
  border: 1px solid #e2e8f0;
}

.wechat-theme-business .md-editor-preview th,
.wechat-theme-business .md-editor-preview td {
  border: 1px solid #e2e8f0;
  padding: 8px 12px;
  text-align: left;
}

.wechat-theme-business .md-editor-preview th {
  background-color: #2c5282;
  color: white;
  font-weight: 600;
}

/* ================== Vintage Theme (Retro Brown) ================== */
.wechat-theme-vintage .md-editor-preview {
  font-family: 'Georgia', 'Times New Roman', serif;
  color: #5d4037;
  line-height: 1.8;
}

.wechat-theme-vintage .md-editor-preview h1 {
  font-size: 28px;
  font-weight: normal;
  text-align: center;
  color: #8b5a3c;
  margin: 25px 0 20px;
  position: relative;
  padding: 25px 0;
  background: linear-gradient(90deg, transparent, rgba(139, 90, 60, 0.1), transparent);
  border: 2px solid #d7ccc8;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(139, 90, 60, 0.1);
}

.wechat-theme-vintage .md-editor-preview h2 {
  font-size: 22px;
  font-weight: normal;
  color: #8b5a3c;
  margin: 20px auto 15px;
  padding: 15px 15px;
  border-bottom: 3px double #d7ccc8;
  text-align: center;
  background: linear-gradient(90deg, transparent, rgba(215, 204, 200, 0.3), transparent);
  display: table;
}

.wechat-theme-vintage .md-editor-preview h3 {
  font-size: 20px;
  font-weight: normal;
  color: #6d4c41;
  margin: 18px 0 10px;
  font-style: italic;
  padding-left: 30px;
  padding-bottom: 8px;
  border-bottom: 1px solid #d7ccc8;
  border-left: 3px solid #8b5a3c;
}

.wechat-theme-vintage .md-editor-preview blockquote {
  border: none;
  background-color: #efebe9;
  color: #5d4037;
  padding: 20px;
  margin: 20px 0;
  border-radius: 5px;
  font-style: italic;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.wechat-theme-vintage .md-editor-preview p {
  text-indent: 2em;
  margin-bottom: 16px;
}

.wechat-theme-vintage .md-editor-preview strong {
  color: #8b5a3c;
  font-weight: 600;
}

.wechat-theme-vintage .md-editor-preview a {
  color: #8b5a3c;
  text-decoration: underline;
  font-style: italic;
}

/* Common fix for code blocks overflow */
.preview-content .md-editor-preview pre {
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* Mac-like Code Block Style for Preview */
.preview-content .md-editor-code {
  position: relative;
  margin-top: 24px !important;
  background: #282c34 !important; /* Force dark background for all themes */
  border-radius: 6px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  color: #abb2bf !important;
}

/* Header simulation */
.preview-content .md-editor-code::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 28px;
  background: #21252b !important; /* Force dark header for all themes */
  border-radius: 6px 6px 0 0;
  /* Mac Dots */
  background-image: radial-gradient(#ff5f56 50%, transparent 50%),
                    radial-gradient(#ffbd2e 50%, transparent 50%),
                    radial-gradient(#27c93f 50%, transparent 50%);
  background-size: 12px 12px;
  background-repeat: no-repeat;
  background-position: 12px center, 32px center, 52px center;
  z-index: 10;
}

/* Adjust PRE padding to account for header */
.preview-content .md-editor-code pre {
  padding-top: 40px !important;
  border-radius: 6px;
  margin: 0 !important;
  color: #abb2bf !important;
}

/* Hide original artifacts in preview */
.preview-content .md-editor-code-head,
.preview-content .md-editor-copy-button,
.preview-content .md-editor-code-lang {
  display: none !important;
}

/* 
   Force Atom One Dark Syntax Highlighting Colors 
   to ensure getComputedStyle picks up the correct dark theme colors 
   regardless of what MdPreview thinks.
*/
.preview-content .hljs-comment,
.preview-content .hljs-quote {
  color: #5c6370 !important;
  font-style: italic;
}
.preview-content .hljs-doctag,
.preview-content .hljs-keyword,
.preview-content .hljs-formula {
  color: #c678dd !important;
}
.preview-content .hljs-section,
.preview-content .hljs-name,
.preview-content .hljs-selector-tag,
.preview-content .hljs-deletion,
.preview-content .hljs-subst {
  color: #e06c75 !important;
}
.preview-content .hljs-literal {
  color: #56b6c2 !important;
}
.preview-content .hljs-string,
.preview-content .hljs-regexp,
.preview-content .hljs-addition,
.preview-content .hljs-attribute,
.preview-content .hljs-meta .hljs-string {
  color: #98c379 !important;
}
.preview-content .hljs-attr,
.preview-content .hljs-variable,
.preview-content .hljs-template-variable,
.preview-content .hljs-type,
.preview-content .hljs-selector-class,
.preview-content .hljs-selector-attr,
.preview-content .hljs-selector-pseudo,
.preview-content .hljs-number {
  color: #d19a66 !important;
}
.preview-content .hljs-symbol,
.preview-content .hljs-bullet,
.preview-content .hljs-link,
.preview-content .hljs-meta,
.preview-content .hljs-selector-id,
.preview-content .hljs-title {
  color: #61aeee !important;
}
.preview-content .hljs-built_in,
.preview-content .hljs-title.class_,
.preview-content .hljs-class .hljs-title {
  color: #e6c07b !important;
}
.preview-content .hljs-emphasis {
  font-style: italic;
}
.preview-content .hljs-strong {
  font-weight: bold;
}
.preview-content .hljs-link {
  text-decoration: underline;
}
</style>
