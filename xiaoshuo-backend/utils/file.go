package utils

import (
	"archive/zip"
	"io"
	"os"
	"regexp"
	"strings"
)

// ReadFileContent 读取文件内容
func ReadFileContent(filepath string) (string, error) {
	// 检查文件扩展名以确定处理方式
	if strings.HasSuffix(strings.ToLower(filepath), ".epub") {
		return readEpubContent(filepath)
	}
	
	// 对于.txt文件，直接读取
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// readEpubContent 读取EPUB文件内容
func readEpubContent(filepath string) (string, error) {
	// 打开EPUB文件（EPUB本质上是一个ZIP文件）
	reader, err := zip.OpenReader(filepath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var fullText strings.Builder
	
	// 遍历EPUB文件中的所有文件
	for _, file := range reader.File {
		// 检查文件是否是HTML或XHTML文件（通常是内容文件）
		if isContentFile(file.Name) {
			f, err := file.Open()
			if err != nil {
				continue // 跳过无法打开的文件
			}
			
			// 读取文件内容
			contentBytes, err := io.ReadAll(f)
			f.Close() // 立即关闭文件
			
			if err != nil {
				continue // 跳过读取失败的文件
			}
			
			// 移除HTML标签，只保留文本内容
			textContent := removeHTMLTags(string(contentBytes))
			
			// 将内容添加到结果中
			fullText.WriteString("\n")
			fullText.WriteString(textContent)
		}
	}

	return fullText.String(), nil
}

// isContentFile 检查文件是否是EPUB内容文件
func isContentFile(filename string) bool {
	// 检查是否是HTML或XHTML文件，但排除导航文件
	lowerName := strings.ToLower(filename)
	if strings.HasSuffix(lowerName, ".html") || 
	   strings.HasSuffix(lowerName, ".xhtml") ||
	   strings.HasSuffix(lowerName, ".htm") {
		// 排除一些常见的非内容文件
		if strings.Contains(lowerName, "toc") || 
		   strings.Contains(lowerName, "nav") ||
		   strings.Contains(lowerName, "cover") ||
		   strings.Contains(lowerName, "titlepage") {
			return false
		}
		return true
	}
	return false
}

// removeHTMLTags 移除HTML标签，只保留文本内容
func removeHTMLTags(html string) string {
	// 使用正则表达式移除HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")
	
	// 清理多余的空白字符
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// 去除首尾空白
	return strings.TrimSpace(text)
}