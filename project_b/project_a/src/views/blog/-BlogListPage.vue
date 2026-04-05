<template>
  <div class="blog-layout">
    <!-- 左侧边栏 -->
    <aside class="sidebar">
      <n-h3 class="sidebar-title" style="padding-top: 18px">📝 所有文章</n-h3>
      <n-text type="secondary">共 {{ totalPosts }} 篇文章 · 按发布时间倒序排列</n-text>

      <!-- 🔍 智能搜索框 -->
      <div class="search-box">
        <n-input
          v-model:value="searchKeyword"
          placeholder="find the thoughts..."
          clearable
          size="small"
          @update:value="onSearch"
          @focus="handleFocus"
          @blur="handleBlur"
        />
        <!-- 搜索建议下拉 -->
        <div v-if="showSuggestions && searchSuggestions.length > 0" class="search-suggestions">
          <div
            v-for="(suggestion, index) in searchSuggestions"
            :key="index"
            class="suggestion-item"
            @click="selectSuggestion(suggestion)"
          >
            {{ suggestion }}
          </div>
        </div>
      </div>

      <div class="calendar-heatmap">
        <n-h3 class="sidebar-title">📆 发文日历</n-h3>
        <HeatmapCalendar
          :posts="allPosts"
          @date-click="onDateSelect"
          :selected-date="selectedDate"
        />
      </div>

      <div class="tag-cloud" style="margin-top: 24px">
        <n-h3 class="sidebar-title">🔖 所有标签</n-h3>
        <n-tag
          v-for="tag in tagCounts"
          :key="tag.name"
          size="small"
          round
          :bordered="true"
          :color="getTagColor(tag.name)"
          style="margin: 4px; cursor: pointer"
          @click="onTagSelect(tag.name)"
        >
          {{ tag.name }} ({{ tag.count }})
        </n-tag>
      </div>
      <div class="tag-cloud" style="margin-top: 24px">
        <n-h3 class="sidebar-title">😊 更多信息</n-h3>
        <n-space vertical>
          <n-button text tag="a" @click="$router.push('/')"> 🏠 foreveryang </n-button>
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

    <!-- 右侧主内容 -->
    <main class="main-content">
      <div v-for="post in paginatedPosts" :key="post.id" class="post-item-layout">
        <div class="post-text">
          <router-link :to="`/blog/${post.slug}`" class="post-title">
            {{ post.title }}
          </router-link>
          <div class="post-meta">
            <n-text type="secondary" class="post-date">{{ post.date }}</n-text>
            <!-- 使用 NAvatarGroup 展示标签头像 -->
            <n-avatar-group
              :options="getPostTagOptions(post)"
              :max="3"
              size="small"
              :style="{ display: 'inline-flex', gap: '4px' }"
            >
              <!-- 自定义每个 avatar -->
              <template #avatar="{ option }">
                <n-tooltip placement="top">
                  <template #trigger>
                    <n-avatar
                      :src="option.src"
                      :fallback="{ text: option.name?.[0]?.toUpperCase() || '?' }"
                      style="cursor: pointer"
                      @click.stop="onTagSelect(option.name)"
                    />
                  </template>
                  {{ option.name }}
                </n-tooltip>
              </template>

              <template #rest="{ rest, options: restOptions }">
                <n-dropdown
                  :options="createDropdownOptions(restOptions)"
                  placement="top"
                  trigger="hover"
                  :style="{ cursor: 'pointer' }"
                >
                  <n-avatar>+{{ rest }}</n-avatar>
                </n-dropdown>
              </template>
            </n-avatar-group>
          </div>
          <div class="post-stats">
            <p class="post-excerpt">{{ post.excerpt }}</p>
            <div class="stats-container">
              <div class="stat-item">
                <n-icon size="16" :component="EyeOutline" />
                <span>{{ formatNumber(post.views) }}</span>
              </div>
              <div class="stat-item">
                <n-icon size="16" :component="HeartOutline" />
                <span>{{ formatNumber(post.likes) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="post-thumbnail">
          <img
            :src="processImageUrl(post.thumbnail)"
            :alt="post.title || '文章缩略图'"
            loading="lazy"
            @error="handleImageError"
            @loadstart="handleLoadStart"
            :data-post-id="post.id"
            class="thumbnail-img"
          />
        </div>
      </div>

      <div class="pagination-wrapper">
        <n-pagination
          v-show="totalPages > 1"
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-size="pageSize"
          show-size-picker
          :page-sizes="[12, 20]"
          @update:page-size="handlePageSizeChange"
          class="pagination"
        />
        <div v-show="totalPages <= 1" class="pagination-placeholder"></div>
      </div>

      <n-empty
        v-if="filteredPosts.length === 0 && !loading"
        description="暂无匹配文章"
        style="margin-top: 64px"
      />
    </main>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, watch, onUnmounted, h } from 'vue'
import { useMessage, NAvatar } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import HeatmapCalendar from '@/components/HeatmapCalendarComponent.vue'
import { EyeOutline, HeartOutline } from '@vicons/ionicons5'

// 默认占位图（80x80 SVG）
const defaultThumbnail =
  'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODAiIGhlaWdodD0iODAiIHZpZXdCb3g9IjAgMCA4MCA4MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjgwIiBoZWlnaHQ9IjgwIiBmaWxsPSIjZjNmNGY2Ii8+CjxwYXRoIGQ9Ik0yNSAyNUg1NVY1NUgyNVYyNVoiIHN0cm9rZT0iI2UyZThmMCIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KPHBhdGggZD0iTTM1IDQwSDQ1IiBzdHJva2U9IiNlMmU4ZjAiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIi8+CjxwYXRoIGQ9Ik0zNSA0NUg0NSIgc3Ryb2tlPSIjZTJlOGYwIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8cGF0aCBkPSJNMzUgNTBINTUiIHN0cm9rZT0iI2UyZThmMCIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiLz4KPC9zdmc+Cg=='

interface BlogPost {
  id: number
  slug: string
  title: string
  tags: string[]
  date: string
  excerpt: string
  views: number
  likes: number
  thumbnail?: string
}

const allPosts = ref<BlogPost[]>([])
const loading = ref<boolean>(true)
const currentPage = ref<number>(1)
const pageSize = ref<number>(12)
const selectedDate = ref<string | undefined>(undefined)
const selectedTag = ref<string | undefined>(undefined)
const searchKeyword = ref<string>('')
const searchSuggestions = ref<string[]>([])
const showSuggestions = ref<boolean>(false)
const blurTimer = ref<number | null>(null)

// ======================
// 🔍 全文搜索：动态构建倒排索引（基于标题 + 标签）
// ======================
let mockInvertedIndex: Record<string, number[]> = {}

const buildMockInvertedIndex = (posts: BlogPost[]): void => {
  const index: Record<string, number[]> = {}
  posts.forEach((post) => {
    const terms = new Set<string>()
    post.title
      .toLowerCase()
      .split(/\s+/)
      .forEach((t) => terms.add(t))
    post.tags.forEach((tag) => terms.add(tag.toLowerCase()))
    terms.forEach((term) => {
      if (!index[term]) index[term] = []
      if (!index[term].includes(post.id)) {
        index[term].push(post.id)
      }
    })
  })
  mockInvertedIndex = index
}

const fullTextSearch = (query: string): number[] => {
  const terms = query
    .toLowerCase()
    .trim()
    .split(/\s+/)
    .filter((t) => t.length > 0)
  if (terms.length === 0) return []

  let resultIds: Set<number> | null = null
  for (const term of terms) {
    const ids = Object.keys(mockInvertedIndex)
      .filter((key) => key.includes(term))
      .flatMap((key) => mockInvertedIndex[key] || [])
      .filter((id): id is number => typeof id === 'number')
    if (resultIds === null) {
      resultIds = new Set(ids)
    } else {
      resultIds = new Set(ids.filter((id) => resultIds!.has(id)))
    }
  }
  return resultIds ? Array.from(resultIds) : []
}

const loadAllPosts = async () => {
  try {
    loading.value = true
    const response = await fetch('/api/posts/getAllPosts')
    if (!response.ok) {
      throw new Error('Network response was not ok')
    }
    const apiPosts: BlogPost[] = await response.json()

    const blogPosts: BlogPost[] = apiPosts.map((post) => ({
      id: post.id,
      slug: post.slug,
      title: post.title,
      tags: post.tags,
      date: post.date,
      excerpt: post.excerpt,
      views: post.views,
      likes: post.likes,
      thumbnail: post.thumbnail,
    }))

    allPosts.value = blogPosts.sort(
      (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
    )
    buildMockInvertedIndex(allPosts.value)
  } catch (error) {
    useMessage().error('加载文章失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAllPosts()
})

onUnmounted(() => {
  if (blurTimer.value !== null) {
    window.clearTimeout(blurTimer.value)
    blurTimer.value = null
  }
})

// ======================
// 🖼️ 标签头像生成（带彩色背景 SVG）
// ======================

type ColorHex = `#${string}`

const stringToColor = (str: string | undefined): ColorHex => {
  const colors: ColorHex[] = [
    '#fce7f3',
    '#e0f2fe',
    '#dcfce7',
    '#fff7ed',
    '#ede9fe',
    '#ffe4e6',
    '#ffedd5',
    '#dcf5e9',
    '#e6f7ff',
    '#f0f9ff',
  ]
  if (!str) {
    return colors[0]! // 使用非空断言确保返回值
  }
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]! // 使用非空断言确保返回值
}

const PlaceholderSvg = (text: string): string => {
  const char = text?.[0]?.toUpperCase() || '?'
  const bgColor = stringToColor(text)
  const textColor = '#334155'
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32">
      <circle cx="16" cy="16" r="16" fill="${bgColor}"/>
      <text x="16" y="21" text-anchor="middle" font-family="system-ui, sans-serif" font-size="16" fill="${textColor}" font-weight="500">${char}</text>
    </svg>
  `
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`
}

const tagIconMap: Record<string, string> = {}

const getPostTagOptions = (post: BlogPost) => {
  return post.tags.map((tag) => ({
    name: tag,
    src: tagIconMap[tag] || PlaceholderSvg(tag),
    fallbackText: tag?.[0]?.toUpperCase() || '?',
  }))
}

const createDropdownOptions = (
  restOptions: { name: string; src: string; fallbackText: string }[],
) => {
  return restOptions.map((opt) => ({
    key: opt.name,
    label: () =>
      h('div', { style: { display: 'flex', alignItems: 'center', gap: '8px' } }, [
        h(NAvatar, {
          size: 'small',
          src: opt.src,
          fallback: { text: opt.fallbackText },
        }),
        h('span', opt.name),
      ]),
    onClick: () => onTagSelect(opt.name),
  }))
}

// ======================
// 🔎 搜索建议 & 防抖
// ======================

const getAllKeywords = (): string[] => {
  const keywords = new Set<string>()
  allPosts.value.forEach((post) => {
    keywords.add(post.title)
    post.tags.forEach((tag) => keywords.add(tag))
  })
  return Array.from(keywords)
}

const computeSuggestions = (query: string): void => {
  if (!query.trim()) {
    searchSuggestions.value = []
    showSuggestions.value = false
    return
  }
  const lowerQuery = query.toLowerCase()
  const allKeywords = getAllKeywords()
  const matches = allKeywords.filter((kw) => kw.toLowerCase().includes(lowerQuery)).slice(0, 5)
  searchSuggestions.value = matches
  showSuggestions.value = matches.length > 0
}

const debounce = <T extends (...args: string[]) => void>(
  func: T,
  delay: number,
): ((...args: Parameters<T>) => void) => {
  let timeoutId: number | null = null
  return (...args: Parameters<T>) => {
    if (timeoutId !== null) {
      window.clearTimeout(timeoutId)
    }
    timeoutId = window.setTimeout(() => {
      func(...args)
    }, delay)
  }
}

const debouncedCompute = debounce((val: string) => {
  computeSuggestions(val)
}, 200)

watch(searchKeyword, (newVal) => {
  debouncedCompute(newVal)
})

const handleFocus = (): void => {
  if (searchKeyword.value.trim()) {
    showSuggestions.value = true
  }
}

const handleBlur = (): void => {
  if (blurTimer.value !== null) {
    window.clearTimeout(blurTimer.value)
  }
  blurTimer.value = window.setTimeout(() => {
    showSuggestions.value = false
    blurTimer.value = null
  }, 200)
}

const selectSuggestion = (text: string): void => {
  searchKeyword.value = text
  showSuggestions.value = false
  onSearch()
}

const onSearch = (): void => {
  currentPage.value = 1
}

// ======================
// ✅ 路由参数同步：tag
// ======================

const route = useRoute()
const router = useRouter()

// 安全提取 tag 参数
const extractTagFromQuery = (): string | undefined => {
  const tag = route.query.tag
  if (Array.isArray(tag)) {
    return tag[0] ?? undefined
  }
  return tag ?? undefined
}

// 监听 URL 中的 tag 参数变化（包括初始加载）
watch(
  () => route.query.tag,
  () => {
    selectedTag.value = extractTagFromQuery()
    currentPage.value = 1
  },
  { immediate: true },
)

// 用户点击标签时，同步到 URL
const onTagSelect = (tag: string): void => {
  const newTag = selectedTag.value === tag ? undefined : tag
  selectedTag.value = newTag
  currentPage.value = 1

  // 构建新查询对象
  const newQuery = { ...route.query }
  if (newTag) {
    newQuery.tag = newTag
  } else {
    delete newQuery.tag
  }

  // 无刷新更新 URL
  router.push({ query: newQuery })
}

// ======================
// ✅ 筛选逻辑
// ======================

const filteredPosts = computed(() => {
  let candidates = [...allPosts.value]

  if (searchKeyword.value.trim()) {
    const ids = fullTextSearch(searchKeyword.value)
    const matchedIds = new Set(ids)
    candidates = candidates.filter((post) => matchedIds.has(post.id))
  }

  if (selectedDate.value) {
    candidates = candidates.filter((post) => post.date === selectedDate.value)
  }

  if (selectedTag.value) {
    candidates = candidates.filter((post) => post.tags.includes(selectedTag.value!))
  }

  return candidates
})

const paginatedPosts = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredPosts.value.slice(start, start + pageSize.value)
})

const totalPages = computed(() => Math.ceil(filteredPosts.value.length / pageSize.value))
const totalPosts = computed(() => filteredPosts.value.length)

const handlePageSizeChange = (size: number): void => {
  pageSize.value = size
  currentPage.value = 1
}

const tagCounts = computed(() => {
  const map: Record<string, number> = {}
  allPosts.value.forEach((post) => {
    post.tags.forEach((tag) => {
      map[tag] = (map[tag] || 0) + 1
    })
  })
  return Object.entries(map)
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 20)
})

const getTagColor = (tag: string) => {
  if (selectedTag.value === tag) {
    return { color: '#e6f7ff', textColor: '#1890ff' }
  }
  return { color: '#f0f9ff', textColor: '#007bff' }
}

const onDateSelect = (date: string): void => {
  selectedDate.value = selectedDate.value === date ? undefined : date
  currentPage.value = 1
}

// ======================
// 🖼️ 图片处理
// ======================

const handleLoadStart = (e: Event) => {
  const img = e.target as HTMLImageElement
  const postId = Number(img.dataset.postId)
  const post = allPosts.value.find((p) => p.id === postId)
  console.log('图片开始加载:', {
    src: img.src,
    hasThumbnail: !!post?.thumbnail,
    thumbnailUrl: post?.thumbnail,
  })
}

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  const postId = Number(img.dataset.postId)

  if (e instanceof ErrorEvent) {
    console.error('图片加载失败:', {
      src: img.src,
      error: e.message,
      postId,
    })
  }
  img.src = defaultThumbnail
}

const processImageUrl = (url?: string): string => {
  if (url) {
    return `/api/proxy/image?url=${encodeURIComponent(url)}&quality=low`
  }
  return ''
}


const formatNumber = (num: number): string => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}
</script>

<style scoped>
.sidebar {
  width: 250px;
  flex-shrink: 0;
  position: fixed;
  top: 24px; /* 与页面顶部边距保持一致 */
  height: calc(100vh - 88px); /* 减去顶部边距和底部间距 */
  overflow-y: auto;
}

.main-content {
  flex: 1;
  min-width: 0;
  margin-left: 288px; /* 侧边栏宽度 + 间距 */
  margin-top: 24px;
}

.blog-layout {
  display: flex;
  gap: 64px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px 48px;
  min-height: 100vh;
}

.sidebar-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 12px;
  color: #333;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.post-item-layout {
  display: flex;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid #e5e7eb;
}

.post-item-layout:last-child {
  border-bottom: none;
}

.post-text {
  flex: 1;
  min-width: 0;
}

.post-title {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  text-decoration: none;
  transition: color 0.2s;
  display: block;
  margin-bottom: 6px;
}

.post-title:hover {
  color: #10b981;
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  margin: 0px 0 2px;
}

.post-date {
  margin-right: 8px;
}

.post-stats {
  margin: 0;
}

.post-excerpt {
  color: #333;
  line-height: 1.6;
  font-size: 0.95rem;
  margin-top: 0;
  margin-bottom: 0;
}

.post-thumbnail {
  width: 160px;
  height: 90px;
  flex-shrink: 0;
}

.post-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  background-color: #f8fafc;
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
  z-index: 1000;
  max-height: 200px;
  overflow-y: auto;
  margin-top: 4px;
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

@media (max-width: 640px) {
  .blog-layout {
    flex-direction: column;
    padding: 0 16px 32px;
  }

  .sidebar {
    width: 100%;
    height: 100%;
    display: block;
    position: static;
    left: auto;
    margin-bottom: 32px;
    margin-bottom: 0;
  }

  .main-content {
    margin-left: 0;
    margin-top: 0;
  }

  .post-item-layout {
    flex-direction: column;
    gap: 12px;
  }

  .post-thumbnail {
    width: 100%;
    height: 180px;
    order: -1;
  }

  .post-thumbnail img {
    border-radius: 10px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  min-height: 40px; /* 或根据实际调整 */
}

.pagination,
.pagination-placeholder {
  height: 40px; /* 统一高度 */
  display: flex;
  align-items: center;
  justify-content: center;
}

.post-thumbnail {
  width: 160px;
  height: 90px;
  flex-shrink: 0;
  overflow: hidden;
  background-color: #f8fafc;
  border-radius: 8px;
}

.thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  display: block;
  transition: transform 0.3s ease;
}

.thumbnail-img:hover {
  transform: scale(1.05);
}

.stats-container {
  display: flex;
  gap: 1.5rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: #666;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.stat-item:hover {
  color: #333;
}

.stat-item span {
  font-weight: 500;
}

/* 极简风格：移除所有装饰效果 */
.stat-item {
  background: transparent;
  border: none;
  padding: 0;
}

/* 图标样式微调 */
:deep(.n-icon) {
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.stat-item:hover :deep(.n-icon) {
  opacity: 1;
}
</style>
