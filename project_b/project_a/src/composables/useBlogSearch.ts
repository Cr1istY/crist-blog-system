import type { BlogPost } from '@/types/blog'
import { ref } from 'vue'

export function useBlogSearch() {
  const invertedIndex = ref<Record<string, number[]>>({})

  const buildIndex = (posts: BlogPost[]) => {
    const index: Record<string, number[]> = {}
    posts.forEach((post) => {
      const terms = new Set<string>()

      const titleLower = post.title.toLowerCase()
      // 添加完整标题
      terms.add(titleLower)
      // 添加标题分词
      titleLower.split(/\s+/).forEach((t) => terms.add(t))

      // 添加简介
      const excerptLower = post.excerpt.toLowerCase()
      terms.add(excerptLower)
      excerptLower.split(/\s+/).forEach((t) => terms.add(t))

      // 添加标签
      post.tags.forEach((originalTag) => {
        const tagLower = originalTag.toLowerCase()
        // 添加完整标签名
        terms.add(tagLower)
        // 添加标签分词 (如果标签内有空格或特殊分隔符)
        tagLower.split(/\s+/).forEach((t) => terms.add(t))
      })

      terms.forEach((term) => {
        if (!index[term]) index[term] = []
        if (!index[term].includes(post.id)) index[term].push(post.id)
      })
    })
    invertedIndex.value = index
  }

  // 保持 search 函数不变
  const search = (query: string): number[] => {
    const terms = query
      .toLowerCase()
      .trim()
      .split(/\s+/)
      .filter((t) => t)
    if (terms.length === 0) return []

    let result = new Set<number>()
    terms.forEach((term) => {
      const matches = Object.keys(invertedIndex.value)
        .filter((key) => key.includes(term))
        .flatMap((key) => invertedIndex.value[key] || [])
      if (result.size === 0) {
        matches.forEach((id) => result.add(id))
      } else {
        const temp = new Set(matches.filter((id) => result.has(id)))
        result = temp
      }
    })
    return Array.from(result)
  }

  return { invertedIndex, buildIndex, search }
}
