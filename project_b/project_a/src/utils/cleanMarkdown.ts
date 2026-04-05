function markdownToPlainText(md: string): string {
  // 1. 移除代码块（避免干扰）
  let text = md.replace(/```[\s\S]*?```/g, '')
  text = text.replace(/`[^`]*`/g, '') // 行内代码

  // 2. 移除链接 [text](url) → text
  text = text.replace(/\[([^\]]+)\]\([^)]*\)/g, '$1')

  // 3. 移除图片 ![alt](src) → (空)
  text = text.replace(/!\[[^\]]*\]\([^)]*\)/g, '')

  // 4. 移除强调/加粗等标记：**bold**、__bold__、*italic*、_italic_
  text = text.replace(/(\*\*|__)(.*?)\1/g, '$2')
  text = text.replace(/(\*|_)(.*?)\1/g, '$2')

  // 5. 移除标题标记 #、列表符号 - * +、引用 >
  text = text.replace(/^(\s*[-*+]\s+|\s*\d+\.\s+|\s*>+\s*|#*)/gm, '')

  // 6. 合并多个空白符（包括换行）为单个空格
  text = text.replace(/\s+/g, ' ')

  // 7. 去除首尾空格
  return text.trim()
}

export default markdownToPlainText
