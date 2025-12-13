package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gohugoio/hugo-goldmark-extensions/extras"
	mathjax "github.com/litao91/goldmark-mathjax"
	admonitions "github.com/stefanfritsch/goldmark-admonitions"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const AppVersion = "1.0.2"

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppVersion returns the application version
func (a *App) GetAppVersion() string {
	return AppVersion
}

// OpenFileResult represents the result of opening a file
type OpenFileResult struct {
	Content string `json:"content"`
	Path    string `json:"path"`
}

// OpenFile opens a file dialog and returns the selected file content and path
func (a *App) OpenFile() (OpenFileResult, error) {
	defaultPath, _ := os.Getwd()

	dialogOptions := runtime.OpenDialogOptions{
		DefaultDirectory: defaultPath,
		DefaultFilename:  "",
		Title:            "打开Markdown文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files (*.md)",
				Pattern:     "*.md",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	selectedFile, err := runtime.OpenFileDialog(a.ctx, dialogOptions)
	if err != nil {
		return OpenFileResult{}, err
	}

	if selectedFile == "" {
		return OpenFileResult{}, fmt.Errorf("未选择文件")
	}

	content, err := os.ReadFile(selectedFile)
	if err != nil {
		return OpenFileResult{}, err
	}

	result := string(content)
	fmt.Printf("成功读取文件: %s, 内容长度: %d\n", selectedFile, len(result))

	return OpenFileResult{
		Content: result,
		Path:    selectedFile,
	}, nil
}

// SaveFile saves content to the specified file path
func (a *App) SaveFile(content, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// SaveFileAs opens a save dialog and saves content to the selected path
func (a *App) SaveFileAs(content string) (string, error) {
	defaultPath, _ := os.Getwd()

	dialogOptions := runtime.SaveDialogOptions{
		DefaultDirectory: defaultPath,
		DefaultFilename:  "untitled.md",
		Title:            "保存Markdown文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files (*.md)",
				Pattern:     "*.md",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	selectedFile, err := runtime.SaveFileDialog(a.ctx, dialogOptions)
	if err != nil {
		return "", err
	}

	if selectedFile == "" {
		return "", fmt.Errorf("未选择保存位置")
	}

	// Ensure .md extension
	if !strings.HasSuffix(strings.ToLower(selectedFile), ".md") {
		selectedFile += ".md"
	}

	err = os.WriteFile(selectedFile, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return selectedFile, nil
}

// ExportToHTML converts markdown content to HTML and saves it
func (a *App) ExportToHTML(content, path string) error {
	// 使用 goldmark 将 Markdown 转换为 HTML
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Linkify,
			extension.Table,         // 表格支持
			extension.TaskList,      // 任务列表
			extension.Footnote,      // 角标支持
			mathjax.MathJax,         // 数学公式支持
			&admonitions.Extender{}, // Admonitions (!!! note) 支持
			extras.New(extras.Config{ // 上标下标支持
				Subscript:   extras.SubscriptConfig{Enable: true},
				Superscript: extras.SuperscriptConfig{Enable: true},
				Delete:      extras.DeleteConfig{Enable: true},
			}),
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"), // 使用 GitHub 风格
				highlighting.WithFormatOptions(
				// html.WithLineNumbers(true), // 可选：显示行号
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // 自动生成标题 ID
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // 硬换行
			html.WithXHTML(),     // XHTML 样式
			html.WithUnsafe(),    // 允许原生 HTML
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return err
	}

	// HTML 模板
	fullHTML := buildHTMLTemplate(buf.String())
	return os.WriteFile(path, []byte(fullHTML), 0644)
}

// buildHTMLTemplate 构建完整的 HTML 模板
func buildHTMLTemplate(content string) string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown Export</title>
    
    <!-- Mermaid -->
    <script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
    
    <!-- ECharts -->
    <script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
    
    <!-- KaTeX -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css">
    <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js"></script>
    <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js"></script>
    
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif; line-height: 1.6; margin: 0; padding: 0; color: #24292e; background-color: #fff; }
        .markdown-body { box-sizing: border-box; min-width: 200px; max-width: 980px; margin: 0 auto; padding: 45px; }
        h1, h2, h3, h4, h5, h6 { margin-top: 24px; margin-bottom: 16px; font-weight: 600; line-height: 1.25; }
        h1 { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h2 { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        
        /* Code Highlight */
        pre { background-color: #f6f8fa; border-radius: 6px; padding: 16px; overflow: auto; }
        code { font-family: 'SFMono-Regular', Consolas, monospace; font-size: 85%; }
        /* Inline code */
        :not(pre) > code { background-color: rgba(27, 31, 35, 0.05); border-radius: 3px; padding: 0.2em 0.4em; }
        
        blockquote { margin: 0; padding: 0 1em; color: #6a737d; border-left: 0.25em solid #dfe2e5; }
        table { border-collapse: collapse; width: 100%; margin-bottom: 16px; }
        table th { font-weight: 600; background-color: #f6f8fa; }
        table th, table td { padding: 6px 13px; border: 1px solid #dfe2e5; }
        .mermaid { text-align: center; margin: 16px 0; min-height: 100px; }
        .echarts-container { width: 100%; height: 400px; margin: 16px 0; }
        
        /* Admonitions */
        .admonition {
            padding: 15px;
            margin-bottom: 20px;
            border: 1px solid transparent;
            border-radius: 4px;
            border-left-width: 5px;
            border-left-style: solid;
        }
        .admonition .adm-title {
            margin: 0;
            margin-bottom: 10px;
            font-weight: bold;
            text-transform: uppercase;
        }
        /* Blue/Info */
        .admonition.adm-note, .admonition.adm-abstract, .admonition.adm-info, .admonition.adm-example, .admonition.adm-quote {
            border-color: #eee; border-left-color: #428bca; background-color: #f6faff;
        }
        /* Green/Success */
        .admonition.adm-success, .admonition.adm-tip, .admonition.adm-hint {
            border-color: #eee; border-left-color: #5cb85c; background-color: #f4fcf4;
        }
        /* Yellow/Warning */
        .admonition.adm-warning, .admonition.adm-caution, .admonition.adm-attention, .admonition.adm-question {
            border-color: #eee; border-left-color: #f0ad4e; background-color: #fcf8f2;
        }
        /* Red/Error */
        .admonition.adm-danger, .admonition.adm-error, .admonition.adm-bug, .admonition.adm-failure {
            border-color: #eee; border-left-color: #d9534f; background-color: #fdf7f7;
        }
        
        /* Footnotes */
        .footnotes {
            margin-top: 30px;
            border-top: 1px solid #eaecef;
            padding-top: 20px;
            font-size: 0.85em;
            color: #6a737d;
        }
    </style>
</head>
<body>
    <div class="markdown-body">` + content + `</div>
    <script>
        window.addEventListener('DOMContentLoaded', function() {
            // Process Mermaid
            document.querySelectorAll('pre code.language-mermaid').forEach(el => {
                const pre = el.parentElement;
                const div = document.createElement('div');
                div.className = 'mermaid';
                div.textContent = el.textContent;
                pre.parentNode.replaceChild(div, pre);
            });
            mermaid.initialize({ startOnLoad: false, theme: 'default', securityLevel: 'loose' });
            mermaid.run();

            // Process ECharts
            document.querySelectorAll('pre code.language-echarts').forEach((el, index) => {
                try {
                    const json = el.textContent;
                    // Use new Function to allow JS object syntax (keys without quotes, functions, etc.)
                    const option = new Function('return ' + json)();
                    const pre = el.parentElement;
                    const div = document.createElement('div');
                    div.className = 'echarts-container';
                    div.id = 'echarts-chart-' + index;
                    pre.parentNode.replaceChild(div, pre);
                    
                    const myChart = echarts.init(div);
                    myChart.setOption(option);
                    window.addEventListener('resize', () => myChart.resize());
                } catch (e) {
                    console.error('ECharts error', e);
                    const pre = el.parentElement;
                    pre.innerHTML = '<p style="color:red">ECharts Error: ' + e.message + '</p>';
                }
            });

            // Math
            if (typeof renderMathInElement !== 'undefined') {
                renderMathInElement(document.body, {
                    delimiters: [
                        {left: '$$', right: '$$', display: true},
                        {left: '$', right: '$', display: false},
                        {left: '\\(', right: '\\)', display: false},
                        {left: '\\[', right: '\\]', display: true}
                    ],
                    throwOnError: false
                });
            }
        });
    </script>
</body>
</html>`
}

// ExportToHTMLAs opens a save dialog and exports content as HTML
func (a *App) ExportToHTMLAs(content string) (string, error) {
	defaultPath, _ := os.Getwd()

	dialogOptions := runtime.SaveDialogOptions{
		DefaultDirectory: defaultPath,
		DefaultFilename:  "export.html",
		Title:            "导出为HTML",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "HTML Files (*.html)",
				Pattern:     "*.html",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	}

	selectedFile, err := runtime.SaveFileDialog(a.ctx, dialogOptions)
	if err != nil {
		return "", err
	}

	if selectedFile == "" {
		return "", fmt.Errorf("未选择保存位置")
	}

	// Ensure .html extension
	if !strings.HasSuffix(strings.ToLower(selectedFile), ".html") {
		selectedFile += ".html"
	}

	err = a.ExportToHTML(content, selectedFile)
	if err != nil {
		return "", err
	}

	return selectedFile, nil
}
