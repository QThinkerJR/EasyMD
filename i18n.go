package main

// TextResources 定义了多语言文本资源
type TextResources struct {
	OpenFileTitle       string
	NoFileSelected      string
	SaveFileTitle       string
	NoSavePath          string
	ExportHtmlTitle     string
	MarkdownFilesFilter string
	HtmlFilesFilter     string
	AllFilesFilter      string
	HtmlLang            string
	HtmlTitle           string
}

// 语言常量
const (
	LangZhCN = "zh-CN"
	LangEnUS = "en-US"
)

// 默认语言资源 (中文)
var zhCNResources = TextResources{
	OpenFileTitle:       "打开Markdown文件",
	NoFileSelected:      "未选择文件",
	SaveFileTitle:       "保存Markdown文件",
	NoSavePath:          "未选择保存位置",
	ExportHtmlTitle:     "导出为HTML",
	MarkdownFilesFilter: "Markdown Files (*.md)",
	HtmlFilesFilter:     "HTML Files (*.html)",
	AllFilesFilter:      "All Files (*.*)",
	HtmlLang:            "zh-CN",
	HtmlTitle:           "Markdown 导出",
}

// 英文语言资源
var enUSResources = TextResources{
	OpenFileTitle:       "Open Markdown File",
	NoFileSelected:      "No file selected",
	SaveFileTitle:       "Save Markdown File",
	NoSavePath:          "No save location selected",
	ExportHtmlTitle:     "Export to HTML",
	MarkdownFilesFilter: "Markdown Files (*.md)",
	HtmlFilesFilter:     "HTML Files (*.html)",
	AllFilesFilter:      "All Files (*.*)",
	HtmlLang:            "en",
	HtmlTitle:           "Markdown Export",
}

// GetTextResources 根据语言代码返回对应的文本资源
func GetTextResources(lang string) TextResources {
	switch lang {
	case LangEnUS, "en":
		return enUSResources
	default:
		return zhCNResources
	}
}
