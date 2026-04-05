import { computed, type Ref } from 'vue'
import type { BlogPost } from '@/types/blog'

// 定义传入 filters 的接口，增加 searchFunction
interface Filters {
  selectedDate: Ref<string | undefined>
  selectedCat: Ref<string | undefined>
  selectedTag: Ref<string | undefined>
  searchKeyword: Ref<string>
  invertedIndex: Ref<Record<string, number[]>>
  searchFunction: (query: string) => number[]
}

export function usePostFiltering(
  allPosts: Ref<BlogPost[]>,
  filters: Filters,
) {
  const { selectedDate, selectedCat, selectedTag, searchKeyword, searchFunction /*, invertedIndex */ } = filters

  const filteredPosts = computed(() => {
    let candidates = [...allPosts.value]

    // 搜索筛选
    if (searchKeyword.value.trim()) {
      const ids = searchFunction(searchKeyword.value)
      const idSet = new Set(ids)
      candidates = candidates.filter((post) => idSet.has(post.id))
    }

    // 日期筛选
    if (selectedDate.value) {
      candidates = candidates.filter((post) => post.date === selectedDate.value)
    }

    if (selectedCat.value) {
      candidates = candidates.filter((post) => post.category === selectedCat.value!)
    }

    // 标签筛选
    if (selectedTag.value) {
      candidates = candidates.filter((post) => post.tags.includes(selectedTag.value!))
    }

    return candidates
  })

  return { filteredPosts }
}
