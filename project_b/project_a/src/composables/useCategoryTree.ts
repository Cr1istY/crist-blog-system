import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import service from '@/utils/request'
import { normalizeCategories, buildTreeOptions, type Category } from '@/utils/category'

export function useCategoryTree(options?: {
  apiPath?: string
  onError?: (error: unknown) => void
}) {
  const categoryOptions = ref<{ label: string; value: string }[]>([])
  const loading = ref(false)
  const message = useMessage()
  const apiPath = options?.apiPath || '/category/getAll'

  const fetchCategories = async () => {
    loading.value = true
    try {
      const token = localStorage.getItem('access_token')
      const response = await service.get<Category[]>(apiPath, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      })

      // 规范化 + 构建树形选项（不修改原始响应）
      const normalized = normalizeCategories(response.data)
      categoryOptions.value = buildTreeOptions(normalized)
    } catch (error) {
      console.error('Failed to load categories:', error)
      options?.onError?.(error) // 优先使用自定义错误处理
      message.error('加载分类列表失败')
    } finally {
      loading.value = false
    }
  }

  return {
    categoryOptions,
    loading,
    fetchCategories,
    // 暴露工具函数供特殊场景使用
    buildTreeOptions,
    normalizeCategories,
  }
}
