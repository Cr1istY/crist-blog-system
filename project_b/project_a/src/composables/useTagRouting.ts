import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Ref } from 'vue'

export function useTagRouting(selectedTag: Ref<string | undefined>) {
  const route = useRoute()
  const router = useRouter()

  // 初始化：从路由读取tag
  if (route.query.tag) {
    const tag = route.query.tag
    selectedTag.value = Array.isArray(tag)
      ? (tag[0] ?? undefined)
      : (tag ?? undefined)
  }

  // 监听selectedTag变化，同步到路由
  watch(selectedTag, (newTag) => {
    const newQuery = { ...route.query }
    if (newTag) {
      newQuery.tag = newTag
    } else {
      delete newQuery.tag
    }
    router.push({ query: newQuery })
  }, { immediate: false })

  // 监听路由变化（前进/后退）
  watch(() => route.query.tag, (newVal) => {
    if (newVal === selectedTag.value) return
    selectedTag.value = Array.isArray(newVal) ? (newVal[0] ?? undefined) : (newVal || undefined)
  })

  return { route, router }
}

export function useCategoryRouting(selectedCategory: Ref<string | undefined>) {
  const route = useRoute()
  const router = useRouter()

  // 初始化：从路由读取category
  if (route.query.cat) {
    const cat = route.query.cat
    selectedCategory.value = Array.isArray(cat)
      ? (cat[0] ?? undefined)
      : (cat ?? undefined)
  }

  // 监听selectedCategory变化，同步到路由
  watch(selectedCategory, (newCat) => {
    const newQuery = { ...route.query }
    if (newCat) {
      newQuery.cat = newCat
    } else {
      delete newQuery.cat
    }
    router.push({ query: newQuery })
  }, { immediate: false })

  // 监听路由变化（前进/后退）
  watch(() => route.query.cat, (newVal) => {
    if (newVal === selectedCategory.value) return
    selectedCategory.value = Array.isArray(newVal) ? (newVal[0] ?? undefined) : (newVal || undefined)
  })

  return { route, router }
}
