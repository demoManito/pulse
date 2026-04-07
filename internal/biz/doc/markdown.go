package doc

import (
	"crypto/sha256"
	"fmt"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// HTMLToMarkdown 将 HTML 内容转换为 Markdown。
func HTMLToMarkdown(html string) (string, error) {
	md, err := htmltomarkdown.ConvertString(html)
	if err != nil {
		return "", fmt.Errorf("html to markdown: %w", err)
	}
	return strings.TrimSpace(md), nil
}

// ContentHash 计算内容的 SHA256 哈希值。
func ContentHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}
