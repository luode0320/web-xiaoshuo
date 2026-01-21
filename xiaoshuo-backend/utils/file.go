package utils

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"xiaoshuo-backend/models"
)

// ParseChapterFromTXT 解析TXT文件的章节
func ParseChapterFromTXT(filepath string) ([]models.Chapter, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 定义章节标题的正则表达式
	// 支持中文数字和阿拉伯数字的章节格式
	chapterRegex := regexp.MustCompile(`^(第[一二三四五六七八九十百千万零\d]+[章节回部卷].*|Chapter\s+\d+|Prologue|Epilogue|引子|序|尾声|后记|\d+\..*|.*卷.*|\d+)$`)

	var chapters []models.Chapter
	var currentChapter *models.Chapter
	chapterIndex := 1

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// 检查是否是章节标题（长度不能太长，一般章节标题不会超过50个字符）
		if chapterRegex.MatchString(trimmedLine) && len(trimmedLine) < 50 {
			// 如果已经有当前章节，先保存它
			if currentChapter != nil && currentChapter.Content != "" {
				currentChapter.WordCount = calculateWordCount(currentChapter.Content)
				chapters = append(chapters, *currentChapter)
			}

			// 创建新章节
			currentChapter = &models.Chapter{
				Title:    trimmedLine,
				Content:  "",
				Position: chapterIndex,
			}
			chapterIndex++
		} else {
			// 如果当前有章节，将当前行添加到章节内容中
			if currentChapter != nil {
				if currentChapter.Content != "" {
					currentChapter.Content += "\n"
				}
				currentChapter.Content += line
			}
		}
	}

	// 保存最后一个章节
	if currentChapter != nil && currentChapter.Content != "" {
		currentChapter.WordCount = calculateWordCount(currentChapter.Content)
		chapters = append(chapters, *currentChapter)
	}

	// 如果没有找到任何章节，将整个文件作为一个章节
	if len(chapters) == 0 {
		content, err := ReadFileContent(filepath)
		if err != nil {
			return nil, err
		}
		chapters = []models.Chapter{
			{
				Title:     "第一章",
				Content:   content,
				Position:  1,
				WordCount: calculateWordCount(content),
			},
		}
	}

	return chapters, nil
}

// ParseChapterFromEPUB 解析EPUB文件的章节
func ParseChapterFromEPUB(filepath string) ([]models.Chapter, error) {
	// 由于EPUB解析库的复杂性，我们返回一个占位章节
	// 在实际部署中，需要使用EPUB解析库来提取章节
	content, err := ReadFileContent(filepath)
	if err != nil {
		return nil, err
	}

	// 对于EPUB文件，我们暂时返回一个包含整个内容的章节
	chapters := []models.Chapter{
		{
			Title:     "第一章",
			Content:   removeHTMLTags(content),
			Position:  1,
			WordCount: calculateWordCount(removeHTMLTags(content)),
		},
	}

	return chapters, nil
}

// removeHTMLTags 移除HTML标签
func removeHTMLTags(htmlContent string) string {
	// 使用正则表达式简单移除HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(htmlContent, "")
}

// calculateWordCount 计算字数
func calculateWordCount(content string) int {
	// 移除空白字符后计算长度
	cleaned := strings.ReplaceAll(content, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")
	cleaned = strings.TrimSpace(cleaned)
	return len([]rune(cleaned)) // 使用rune来正确处理中文字符
}

// IsEPUBFile 检查文件是否为EPUB格式
func IsEPUBFile(filepath string) bool {
	return strings.ToLower(filepath[len(filepath)-5:]) == ".epub"
}

// ReadFileContent 读取文件内容
func ReadFileContent(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// DeleteFile 删除指定路径的文件
func DeleteFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// 文件不存在，无需删除
		return nil
	}
	
	return os.Remove(filepath)
}