package utils

import (
	"strings"
)

func ExtractPostTitle(content string) string {
	index := strings.Index(content, "#")
	if index == -1 {
		return ""
	}
	afterHash := content[index+1:]
	lines := strings.Split(afterHash, "\n")
	title := strings.TrimSpace(lines[0])
	return title
}
