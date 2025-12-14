# EasyMD

EasyMD 是一个基于 [Wails](https://wails.io/) + [Vue 3](https://vuejs.org/) 的现代化 Markdown 桌面编辑器。支持实时预览、流程图、甘特图、数学公式、代码高亮和 HTML、PDF 导出等丰富功能。

<img width="1480" height="990" alt="image" src="https://github.com/user-attachments/assets/36fbf767-2766-446c-8319-a027cbd9ac02" />

## ✨ 功能特性

- **实时预览**: 所见即所得的编辑体验。
- **丰富扩展**: 支持表格、任务列表、脚注、上标/下标等 Markdown 扩展语法。
- **图表支持**: 内置 Mermaid 支持，轻松绘制流程图、甘特图、时序图等。
- **数学公式**: 支持 MathJax，完美渲染数学公式。
- **代码高亮**: 支持多种编程语言的代码高亮显示。
- **文件管理**: 支持打开、保存、另存为本地文件。
- **导出功能**: 支持将文档导出为HTML和PDF。

## 🛠️ 技术栈

- **后端**: Go, Wails v2
- **前端**: Vue 3, Vite
- **UI 组件**: TDesign Vue Next
- **编辑器内核**: md-editor-v3
- **Markdown 解析 (Go)**: goldmark (及其扩展)

## 🚀 快速开始

### 前置要求

在开始之前，请确保您的开发环境已安装以下工具：

- [Go](https://go.dev/) (推荐 1.23 或更高版本)
- [Node.js](https://nodejs.org/) (npm)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### 安装 Wails

如果您还没有安装 Wails，请运行以下命令：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 获取代码

```bash
git clone https://github.com/QThinkerJR/EasyMD.git
cd EasyMD
```

### 开发模式

在项目根目录下运行以下命令启动开发模式：

```bash
wails dev
```

此命令将启动一个 Vite 开发服务器，支持前端热重载。如果需要在浏览器中调试，可以访问控制台输出的本地地址（通常是 http://localhost:34115）。

### 构建发布

要构建生产环境的可执行文件，请运行：

```bash
wails build
```

构建产物将位于 `build/bin` 目录下。

## 📦 依赖管理

### 后端依赖 (Go)

主要依赖包括：
- `github.com/wailsapp/wails/v2`: Wails 框架
- `github.com/yuin/goldmark`: Markdown 解析器
- `github.com/litao91/goldmark-mathjax`: MathJax 支持
- `github.com/stefanfritsch/goldmark-admonitions`: 提示框支持
- `github.com/gohugoio/hugo-goldmark-extensions`: 其他扩展

### 前端依赖 (Vue)

主要依赖包括：
- `vue`: 前端框架
- `md-editor-v3`: Markdown 编辑器组件
- `tdesign-vue-next`: UI 组件库
- `@vavt/v3-extension`: Vue 3 扩展工具

## 📄 开源协议

本项目采用 [MIT License](LICENSE) 开源协议。
