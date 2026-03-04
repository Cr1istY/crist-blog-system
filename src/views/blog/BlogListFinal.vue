<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="blog-layout">
    <BlogSidebar
      :posts="allPosts"
      :total-posts="totalPosts"
      v-model:date="selectedDate"
      v-model:cat="selectedCat"
      v-model:tag="selectedTag"
      v-model:search="searchKeyword"
      @clear-filters="clearFilters"
    />

    <main class="main-content">
      <!-- 文章列表 -->
      <BlogPostItem
        v-for="post in visiblePosts"
        :key="post.id"
        :post="post"
        :show-pin="showPinBadge(post)"
        @cat-click="handleCatClick"
        @tag-click="handleTagClick"
      />

      <!-- 哨兵元素：用于触发加载 -->
      <!-- 只有在移动端模式 (isMobile) 且还有更多数据时才显示 -->
      <div v-if="isMobile && hasMoreData" ref="loadTriggerRef" class="load-trigger">
        <n-spin v-if="isLoadingMore" size="small" description="加载中..." />
        <n-divider v-else dashed>上拉加载更多</n-divider>
      </div>

      <!-- 无更多数据提示 -->
      <div v-if="isMobile && !hasMoreData && sortedPosts.length > 0" class="no-more-text">
        <n-divider>没有更多文章了</n-divider>
      </div>

      <!-- 桌面端分页器 (仅在非移动端显示) -->
      <PaginationControls
        v-if="!isMobile && totalPages > 1"
        v-model:page="currentPage"
        v-model:page-size="pageSize"
        :total-pages="totalPages"
        @update:page-size="handlePageSizeChange"
      />

      <n-empty v-if="!loading && sortedPosts.length === 0" description="暂无匹配文章" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import BlogSidebar from '@/components/blog/BlogSideabar.vue'
import BlogPostItem from '@/components/blog/blog-post-item.vue'
import PaginationControls from '@/components/blog/PaginationControls.vue'
import { useBlogSearch } from '@/composables/useBlogSearch'
import { usePostFiltering } from '@/composables/usePostFiltering'
import { useTagRouting, useCategoryRouting } from '@/composables/useTagRouting'
import type { BlogPost, ApiPost } from '@/types/blog'

const message = useMessage()
const allPosts = ref<BlogPost[]>([])
const loading = ref(true)

// --- 核心配置 ---
const DESKTOP_PAGE_SIZE = 12
const MOBILE_PAGE_SIZE = 10 // 移动端每次加载的数量

const currentPage = ref(1)
const pageSize = ref(DESKTOP_PAGE_SIZE)
const isMobile = ref(false)
const isLoadingMore = ref(false) // 正在加载更多的状态
const loadTriggerRef = ref<HTMLElement | null>(null) // 哨兵元素引用

const selectedDate = ref<string>()
const selectedTag = ref<string>()
const selectedCat = ref<string>()
const searchKeyword = ref('')

// 路由同步
useTagRouting(selectedTag)
useCategoryRouting(selectedCat)

// 搜索与过滤
const { invertedIndex, buildIndex, search } = useBlogSearch()
const { filteredPosts } = usePostFiltering(allPosts, {
  selectedDate,
  selectedCat,
  selectedTag,
  searchKeyword,
  invertedIndex,
  searchFunction: search,
})

// 排序逻辑
const hasActiveFilter = computed(
  () => searchKeyword.value.trim() || selectedDate.value || selectedTag.value || selectedCat.value,
)

const sortedPosts = computed(() => {
  const list = [...filteredPosts.value]
  if (hasActiveFilter.value) {
    return list.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
  }
  return list.sort((a, b) => {
    if (a.is_pinned !== b.is_pinned) return a.is_pinned ? -1 : 1
    if (a.is_pinned && b.is_pinned) return a.pinned_order - b.pinned_order
    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })
})

// 【关键修改】visiblePosts：根据当前页码切片，而不是直接展示所有
// 在移动端，随着 currentPage 增加，这里会自动包含更多数据
const visiblePosts = computed(() => {
  const start = 0
  const end = currentPage.value * pageSize.value
  return sortedPosts.value.slice(start, end)
})

// 计算总页数（用于判断是否有更多数据）
const totalPages = computed(() => Math.ceil(sortedPosts.value.length / pageSize.value))
const totalPosts = computed(() => sortedPosts.value.length)

// 是否还有更多数据可加载
const hasMoreData = computed(() => currentPage.value < totalPages.value)

// 事件处理
const showPinBadge = (post: BlogPost) => !hasActiveFilter.value && post.is_pinned

const resetPagination = () => {
  currentPage.value = 1
  // 重置时不需要重新设置 pageSize，因为 watch 会处理，或者在这里显式处理
}

const handleTagClick = (tag: string) => {
  if (selectedTag.value === tag) {
    selectedTag.value = ''
  } else {
    selectedTag.value = tag
  }
  resetPagination()
}

const handleCatClick = (cat: string) => {
  if (selectedCat.value === cat) {
    selectedCat.value = ''
  } else {
    selectedCat.value = cat
  }
  resetPagination()
}

const clearFilters = () => {
  selectedDate.value = selectedTag.value = searchKeyword.value = ''
  resetPagination()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  resetPagination()
}

// --- 移动端自适应与无限滚动逻辑 ---

// 1. 检查屏幕尺寸并设置 PageSize
const checkScreenSize = () => {
  const wasMobile = isMobile.value
  isMobile.value = window.innerWidth <= 640

  if (isMobile.value) {
    pageSize.value = MOBILE_PAGE_SIZE
  } else {
    pageSize.value = DESKTOP_PAGE_SIZE
  }

  // 如果从移动端切回桌面端，或者反之，重置页码以防数据错乱
  if (wasMobile !== isMobile.value) {
    currentPage.value = 1
  }
}

// 2. 加载更多的函数
const loadMore = async () => {
  if (isLoadingMore.value || !hasMoreData.value) return

  isLoadingMore.value = true
  try {
    // 模拟网络延迟体验 (可选，实际项目中如果是纯前端过滤则不需要 await)
    // 因为你的数据是全量获取后前端过滤的，所以这里其实是瞬间完成的
    // 但为了 UI 反馈，我们保留一个微小的 tick
    await nextTick()

    currentPage.value++
  } catch (e) {
    message.error('加载失败' + e)
  } finally {
    isLoadingMore.value = false
  }
}

// 3. Intersection Observer 实现
let observer: IntersectionObserver | null = null

const setupObserver = () => {
  if (observer) observer.disconnect()

  observer = new IntersectionObserver(
    (entries) => {
      const [entry] = entries
      if (entry?.isIntersecting && hasMoreData.value && !isLoadingMore.value) {
        loadMore()
      }
    },
    {
      rootMargin: '100px', // 提前 100px 开始加载
      threshold: 0.1,
    },
  )

  if (loadTriggerRef.value) {
    observer.observe(loadTriggerRef.value)
  }
}

onMounted(async () => {
  // 初始化屏幕检测
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)

  // 设置观察器
  // 注意：需要在 DOM 渲染后设置，如果初始就是移动端且有数据
  nextTick(() => {
    setupObserver()
  })

  // 加载数据
  loading.value = true
  try {
    const res = await fetch('/api/posts/getAllPosts')
    if (!res.ok) throw new Error('API error')
    const apiPosts: ApiPost[] = await res.json()

    allPosts.value = apiPosts.map((p) => ({
      id: p.id,
      slug: p.slug,
      title: p.title,
      category: p.category || '',
      tags: Array.isArray(p.tags) ? p.tags : [],
      date: p.date || p.published_at?.split('T')[0] || '',
      excerpt: p.excerpt || '',
      views: p.views || 0,
      likes: p.likes || 0,
      thumbnail: p.thumbnail,
      is_pinned: p.is_pinned ?? false,
      pinned_order: p.pinned_order ?? 0,
    }))

    buildIndex(allPosts.value)
  } catch {
    message.error('加载文章失败')
  } finally {
    loading.value = false
    // 数据加载完后，重新检查观察器（确保哨兵元素已存在）
    nextTick(() => setupObserver())
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize)
  if (observer) observer.disconnect()
})

// 当筛选条件变化导致 sortedPosts 变化时，需要重新连接观察器
// 因为 DOM 可能会重排，哨兵元素位置变了
import { watch } from 'vue'
watch(
  [sortedPosts, isMobile],
  () => {
    nextTick(() => setupObserver())
  },
  { deep: false },
)
</script>

<style scoped>
.blog-layout {
  display: flex;
  gap: 64px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px 48px;
  min-height: 100vh;
}
.main-content {
  flex: 1;
  margin-left: 288px;
  margin-top: 24px;
}

/* 加载触发区域样式 */
.load-trigger {
  padding: 20px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60px;
}

.no-more-text {
  padding: 20px 0;
  color: #999;
  text-align: center;
  font-size: 14px;
}

@media (max-width: 640px) {
  .blog-layout {
    flex-direction: column;
    padding: 0 16px 32px;
  }
  .main-content {
    margin-left: 0;
    margin-top: 0;
  }
}
</style>
