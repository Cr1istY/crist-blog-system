<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <aside class="sidebar">
    <n-h3 class="sidebar-title">📝 所有文章</n-h3>
    <n-text type="secondary">共 {{ totalPosts }} 篇文章 · 按发布时间倒序排列</n-text>

    <div class="search-box">
      <n-input
        v-model:value="localSearch"
        placeholder="find the thoughts..."
        clearable
        size="small"
        @update:value="onSearchInput"
        @focus="showSuggestions = true"
        @blur="hideSuggestionsIfNotHovered"
      />
      <div
        v-if="showSuggestions && suggestions.length"
        class="search-suggestions"
        @mouseenter="isMouseOverSuggestions = true"
        @mouseleave="isMouseOverSuggestions = false"
      >
        <div
          v-for="(s, i) in suggestions"
          :key="i"
          class="suggestion-item"
          @click="selectSuggestion(s)"
        >
          {{ s }}
        </div>
      </div>
    </div>

    <div class="calendar-heatmap">
      <n-h3 class="sidebar-title">📆 发文日历</n-h3>
      <HeatmapCalendar
        :posts="posts"
        :selected-date="modelDate"
        @date-click="$emit('update:date', $event === modelDate ? undefined : $event)"
      />
    </div>

    <div class="tag-cloud">
      <div v-if="!isTargetSize">
        <n-h3 class="sidebar-title">🧩 分类</n-h3>
        <n-space wrap>
          <n-tag
            v-for="cat in topCategories"
            :key="cat.name"
            size="small"
            round
            :bordered="true"
            :color="getCatColor(cat.name)"
            style="cursor: pointer"
            @click="$emit('update:cat', cat.name === modelCat ? undefined : cat.name)"
          >
            {{ cat.name }}
          </n-tag>
        </n-space>
      </div>
      <n-h3 class="sidebar-title">🗂️ 标签</n-h3>
      <n-space wrap>
        <n-tag
          v-for="tag in topTags"
          :key="tag.name"
          size="small"
          round
          :bordered="true"
          :color="getTagColor(tag.name)"
          style="cursor: pointer"
          @click="$emit('update:tag', tag.name === modelTag ? undefined : tag.name)"
        >
          {{ tag.name }} ({{ tag.count }})
        </n-tag>
      </n-space>
    </div>

    <div class="footer-links">
      <n-h3 class="sidebar-title">😊 更多信息</n-h3>
      <n-space vertical>
        <n-button text @click="$router.push('/')">🏠 foreveryang</n-button>
        <n-button
          text
          tag="a"
          href="https://github.com/Cr1istY/foreveryangDot-frontend"
          target="_blank"
        >
          💻 source code
        </n-button>
        <n-button text tag="a" href="https://beian.miit.gov.cn/" target="_blank">
          渝ICP备2025056615号
        </n-button>
      </n-space>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import HeatmapCalendar from '../HeatmapCalendarComponent.vue'
import type { BlogPost } from '@/types/blog'

interface Props {
  posts: BlogPost[]
  totalPosts: number
  modelValue?: string // 兼容旧API
  date?: string
  tag?: string
  cat?: string
  search?: string
}

interface Emits {
  'update:date': [date?: string]
  'update:tag': [tag?: string]
  'update:cat': [cat?: string]
  'update:search': [keyword: string]
  'clear-filters': []
}

const props = withDefaults(defineProps<Props>(), {})
const emit = defineEmits<Emits>()

// 双向绑定代理
const modelDate = computed<string>({
  get: () => props.date ?? '',
  set: (val) => emit('update:date', val || undefined),
})
const modelTag = computed<string>({
  get: () => props.tag ?? '',
  set: (val) => emit('update:tag', val || undefined),
})
const modelCat = computed<string>({
  get: () => props.cat ?? '',
  set: (val) => emit('update:cat', val || undefined),
})

const localSearch = ref<string>(props.search || '')
const showSuggestions = ref<boolean>(false)
const isMouseOverSuggestions = ref<boolean>(false)
const suggestions = ref<string[]>([])
const keywordCache = ref<Set<string>>(new Set())

// --- 搜索建议相关逻辑 ---

// 从 posts 中提取所有唯一的标题作为建议关键词
const updateKeywordCache = (): void => {
  const newCache = new Set<string>()
  props.posts.forEach((post) => {
    if (post.title) {
      newCache.add(post.title)
    }
  })
  keywordCache.value = newCache
}

// 计算搜索建议列表
const computeSuggestions = (query: string): void => {
  if (!query.trim()) {
    suggestions.value = []
    return
  }
  const lowerQuery = query.toLowerCase()
  suggestions.value = [...keywordCache.value]
    .filter((kw) => kw.toLowerCase().includes(lowerQuery))
    .slice(0, 5)
}

// 定义防抖函数类型
type DebounceFunction<T extends unknown[]> = (...args: T) => void

// 实现防抖函数
const debounce = <T extends unknown[]>(
  fn: DebounceFunction<T>,
  delay: number,
): DebounceFunction<T> => {
  let timer: number | null = null
  return (...args: T) => {
    if (timer) clearTimeout(timer)
    timer = window.setTimeout(() => fn(...args), delay) as unknown as number // 类型转换
  }
}

// 防抖版本的建议更新函数
const updateSuggestionsDebounced = debounce((query: string) => {
  computeSuggestions(query)
}, 250)

// 输入框值更新时的处理函数
const onSearchInput = (val: string): void => {
  // 立即发出搜索事件，让父组件过滤数据
  emit('update:search', val)
  // 防抖更新本地建议列表
  updateSuggestionsDebounced(val)
}

// 处理点击建议项
const selectSuggestion = (text: string): void => {
  localSearch.value = text
  // 发出搜索事件
  emit('update:search', text)
  // 隐藏建议列表
  showSuggestions.value = false
  // 清除鼠标悬停状态
  isMouseOverSuggestions.value = false
}

// 处理输入框失去焦点时的逻辑
const hideSuggestionsIfNotHovered = (): void => {
  // 使用 nextTick 确保 click 事件处理完毕
  nextTick(() => {
    if (!isMouseOverSuggestions.value) {
      showSuggestions.value = false
    }
  })
}

// --- 标签云相关逻辑 ---

// 计算热门标签
const topTags = computed(() => {
  const map = new Map<string, number>()
  props.posts.forEach((p) => {
    p.tags.forEach((t) => {
      const currentCount = map.get(t) || 0
      map.set(t, currentCount + 1)
    })
  })
  return Array.from(map.entries())
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 13)
})

const topCategories = computed(() => {
  const map = new Map<string, number>()
  props.posts.forEach((p) => {
    const currentCount = map.get(p.category) || 0
    map.set(p.category, currentCount + 1)
  })
  return Array.from(map.entries())
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
})

// 获取标签颜色
const getTagColor = (tag: string) => {
  return modelTag.value === tag
    ? { color: '#e6f7ff', textColor: '#1890ff' }
    : { color: '#f0f9ff', textColor: '#007bff' }
}

const getCatColor = (cat: string) => {
  return modelCat.value === cat
    ? { color: '#e6f7ff', textColor: '#1890ff' }
    : { color: '#f0f9ff', textColor: '#007bff' }
}

// --- 监听器 ---

// 监听 props.posts 的变化，更新关键词缓存和建议列表
watch(
  () => props.posts,
  () => {
    updateKeywordCache()
    // 如果当前有搜索词，重新计算建议
    if (localSearch.value) {
      computeSuggestions(localSearch.value)
    } else {
      // 如果没有搜索词，清空建议
      suggestions.value = []
    }
  },
  { immediate: true }, // 组件挂载时立即执行一次，初始化缓存
)

const isTargetSize = ref(false)
const windowHeight = ref(window.innerHeight)
const windowWidth = ref(window.innerWidth)

// 节流函数
const throttle = <T extends (...args: never[]) => void>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timer: number | null = null
  return (...args: Parameters<T>) => {
    if (timer) return
    timer = window.setTimeout(() => {
      fn(...args)
      timer = null
    }, delay)
  }
}

// 更新窗口尺寸的函数
const updateWindowSize = () => {
  windowHeight.value = window.innerHeight
  windowWidth.value = window.innerWidth
}

// 使用节流优化后的处理函数
const handleResize = throttle(updateWindowSize, 200) // 100ms内只执行一次

watch([windowHeight, windowWidth], ([height, width]) => {
  isTargetSize.value = height < 970 && width > 670
})

// 使用节流优化后的事件监听
window.addEventListener('resize', handleResize)


// 监听外部 search prop 的变化，同步到 localSearch
watch(
  () => props.search,
  (newVal) => {
    localSearch.value = newVal || ''
    // 如果外部清空了搜索，也应清空建议
    if (!newVal) {
      suggestions.value = []
    }
  },
)

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.sidebar {
  width: 260px;
  position: fixed;
  top: 24px;
  height: calc(100vh - 88px);
  overflow-y: auto;
  padding-right: 16px;
}
.sidebar-title {
  font-size: 16px;
  font-weight: 600;
  margin: 18px 0 12px;
  color: #333;
}
.search-box {
  position: relative;
  margin: 16px 0;
}
.search-suggestions {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-height: 200px;
  overflow-y: auto;
  margin-top: 4px;
  z-index: 1000;
}
.suggestion-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 14px;
  color: #374151;
}
.suggestion-item:hover {
  background-color: #f9fafb;
}
.tag-cloud,
.footer-links {
  margin-top: 24px;
}
@media (max-width: 640px) {
  .sidebar {
    width: 100%;
    position: static;
    height: auto;
    margin-bottom: 24px;
  }
}
</style>
