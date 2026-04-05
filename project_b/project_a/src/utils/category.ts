// 定义常量
export const UUID_NIL = '00000000-0000-0000-0000-000000000000'

// 分类接口（与后端一致）
export interface Category {
  id: string
  name: string
  parent_id?: string | null
}

/**
 * 规范化分类数据：将 UUID_NIL 转为 null（不修改原数据）
 */
export function normalizeCategories(categories: Category[]): Category[] {
  return categories.map(cat => ({
    ...cat,
    parent_id: cat.parent_id === UUID_NIL ? null : cat.parent_id
  }))
}

/**
 * 构建带缩进的树形下拉选项（纯函数，无副作用）
 */
export function buildTreeOptions(
  categories: Category[],
  parentId: string | null = null,
  depth = 0
): { label: string; value: string }[] {
  return categories
    .filter(cat => cat.parent_id === parentId)
    .flatMap(cat => {
      const indent = '　'.repeat(depth) // 全角空格缩进
      const current = [{
        label: depth > 0 ? `${indent}└─ ${cat.name}` : cat.name,
        value: cat.id
      }]
      const children = buildTreeOptions(categories, cat.id, depth + 1)
      return [...current, ...children]
    })
}
