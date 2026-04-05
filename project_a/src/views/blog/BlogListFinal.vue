<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="blog-layout">
    <BlogSidebar
      v-if="!isMobile"
      :posts="allPosts"
      :total-posts="sortedPosts.length"
      v-model:date="selectedDate"
      v-model:cat="selectedCat"
      v-model:tag="selectedTag"
      v-model:search="searchKeyword"
      @clear-filters="clearFilters"
    />

    <BlogSideBarSearcher v-if="isMobile" v-model="searchKeyword" :posts="allPosts" />

    <main class="main-content">
      <!--
        核心区别：
        移动端 (isMobile=true): 渲染 visiblePosts (累积列表)
        桌面端 (isMobile=false): 渲染 paginatedPosts (当前页切片)
      -->
      <BlogPostItem
        v-for="post in displayPosts"
        :key="post.id"
        :post="post"
        :show-pin="showPinBadge(post)"
        @cat-click="handleCatClick"
        @tag-click="handleTagClick"
      />

      <!-- === 移动端区域：无限滚动 === -->
      <template v-if="isMobile">
        <!-- 哨兵元素：用于触发加载 -->
        <div v-if="hasMoreData" ref="loadTriggerRef" class="load-trigger">
          <n-spin v-if="isLoadingMore" size="small" description="加载中..." />
          <n-divider v-else dashed>上拉加载更多</n-divider>
        </div>

        <!-- 无更多数据提示 -->
        <div v-if="!hasMoreData && sortedPosts.length > 0" class="no-more-text">
          <n-divider>没有更多文章了</n-divider>
        </div>
      </template>

      <!-- === 桌面端区域：分页控件 === -->
      <template v-else>
        <PaginationControls
          v-if="totalPages > 1"
          v-model:page="currentPage"
          v-model:page-size="pageSize"
          :total-pages="totalPages"
          @update:page-size="handlePageSizeChange"
        />
      </template>

      <!-- 空状态 -->
      <n-empty v-if="!loading && sortedPosts.length === 0" description="暂无匹配文章" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useMessage } from 'naive-ui'
import BlogSidebar from '@/components/blog/BlogSideabar.vue'
import BlogPostItem from '@/components/blog/blog-post-item.vue'
import PaginationControls from '@/components/blog/PaginationControls.vue'
import BlogSideBarSearcher from '@/components/blog/blog-side-bar-searcher.vue'
import { useBlogSearch } from '@/composables/useBlogSearch'
import { usePostFiltering } from '@/composables/usePostFiltering'
import { useTagRouting, useCategoryRouting } from '@/composables/useTagRouting'
import type { BlogPost, ApiPost } from '@/types/blog'

const message = useMessage()
const allPosts = ref<BlogPost[]>([])
const loading = ref(true)

// --- 配置 ---
const DESKTOP_PAGE_SIZE = 12
const MOBILE_PAGE_SIZE = 10

// --- 状态 ---
const currentPage = ref(1)
const pageSize = ref(DESKTOP_PAGE_SIZE)
const isMobile = ref(false) // 设备状态
const isLoadingMore = ref(false)
const loadTriggerRef = ref<HTMLElement | null>(null)

// 筛选状态
const selectedDate = ref<string>()
const selectedTag = ref<string>()
const selectedCat = ref<string>()
const searchKeyword = ref('')

// 路由同步
useTagRouting(selectedTag)
useCategoryRouting(selectedCat)

// 搜索与过滤逻辑
const { invertedIndex, buildIndex, search } = useBlogSearch()
const { filteredPosts } = usePostFiltering(allPosts, {
  selectedDate,
  selectedCat,
  selectedTag,
  searchKeyword,
  invertedIndex,
  searchFunction: search,
})

// --- 排序逻辑 ---
const hasActiveFilter = computed(
  () =>
    !!searchKeyword.value.trim() ||
    !!selectedDate.value ||
    !!selectedTag.value ||
    !!selectedCat.value,
)

const sortedPosts = computed(() => {
  const list = [...filteredPosts.value]
  if (hasActiveFilter.value) {
    return list.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
  }
  // 置顶逻辑
  return list.sort((a, b) => {
    if (a.is_pinned !== b.is_pinned) return a.is_pinned ? -1 : 1
    if (a.is_pinned && b.is_pinned) return a.pinned_order - b.pinned_order
    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })
})

// --- 核心显示逻辑区分 ---

// 1. 移动端模式：返回从第1条到当前页末的所有数据 (累积)
const visiblePosts = computed(() => {
  if (!isMobile.value) return []
  const end = currentPage.value * pageSize.value
  return sortedPosts.value.slice(0, end)
})

// 2. 桌面端模式：仅返回当前页的数据 (切片)
const paginatedPosts = computed(() => {
  if (isMobile.value) return []
  const start = (currentPage.value - 1) * pageSize.value
  return sortedPosts.value.slice(start, start + pageSize.value)
})

// 根据设备状态决定渲染哪个列表
const displayPosts = computed(() => {
  return isMobile.value ? visiblePosts.value : paginatedPosts.value
})

// 分页计算
const totalPages = computed(() => Math.ceil(sortedPosts.value.length / pageSize.value))
const hasMoreData = computed(() => currentPage.value < totalPages.value)

// --- 事件处理 ---

const showPinBadge = (post: BlogPost) => !hasActiveFilter.value && post.is_pinned

const resetPagination = () => {
  currentPage.value = 1
  // 如果是移动端，重置后需要重新触发观察器检查（因为列表长度变了）
  if (isMobile.value) {
    nextTick(() => setupObserver())
  }
}

const handleTagClick = (tag: string) => {
  selectedTag.value = selectedTag.value === tag ? '' : tag
  resetPagination()
}

const handleCatClick = (cat: string) => {
  selectedCat.value = selectedCat.value === cat ? '' : cat
  resetPagination()
}

const clearFilters = () => {
  selectedDate.value = undefined
  selectedTag.value = undefined
  searchKeyword.value = ''
  resetPagination()
}

const handlePageSizeChange = (size: number) => {
  pageSize.value = size
  resetPagination()
}

// --- 设备检测与自适应 ---

const checkScreenSize = () => {
  const wasMobile = isMobile.value
  // 断点可根据需求调整，这里沿用之前的 640px 或 768px
  isMobile.value = window.innerWidth <= 768

  // 切换设备时更新每页数量
  pageSize.value = isMobile.value ? MOBILE_PAGE_SIZE : DESKTOP_PAGE_SIZE

  // 【关键】设备类型切换时，必须重置页码到第1页
  // 否则：桌面看第5页 -> 变手机 -> 此时 currentPage=5，但手机每页少，可能导致索引越界或显示空白
  if (wasMobile !== isMobile.value) {
    currentPage.value = 1
    // 如果变成了移动端，需要重新设置观察器
    if (isMobile.value) {
      nextTick(() => setupObserver())
    }
  }
}

// --- 无限滚动逻辑 (仅移动端有效) ---

const loadMore = async () => {
  // 双重检查：确保只在移动端且有数据时加载
  if (!isMobile.value || isLoadingMore.value || !hasMoreData.value) return

  isLoadingMore.value = true
  try {
    await nextTick()
    currentPage.value++
  } catch (e) {
    message.error('加载失败' + e)
  } finally {
    isLoadingMore.value = false
    // 加载完后，DOM更新，观察器会自动监测新的哨兵位置（如果还有下一页）
  }
}

let observer: IntersectionObserver | null = null

const setupObserver = () => {
  // 清理旧的观察器
  if (observer) {
    observer.disconnect()
    observer = null
  }

  // 只有在移动端、有触发元素、且还有更多数据时才建立观察
  if (!isMobile.value || !loadTriggerRef.value || !hasMoreData.value) {
    return
  }

  observer = new IntersectionObserver(
    (entries) => {
      const [entry] = entries
      if (entry?.isIntersecting) {
        loadMore()
      }
    },
    {
      rootMargin: '100px', // 提前100px加载
      threshold: 0.1,
    },
  )

  observer.observe(loadTriggerRef.value)
}

// --- 生命周期 ---

onMounted(async () => {
  // 初始化检测
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)

  // 初始数据加载
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
      user: p.user ?? {},
    }))

    buildIndex(allPosts.value)
  } catch {
    message.error('加载文章失败')
  } finally {
    loading.value = false
    // 数据加载完成后，如果是移动端，设置观察器
    nextTick(() => {
      if (isMobile.value) setupObserver()
    })
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', checkScreenSize)
  if (observer) observer.disconnect()
})

// --- 监听变化 ---

// 当筛选条件变化 (sortedPosts 变化) 或 设备类型变化 (isMobile 变化) 时，重新设置观察器
watch(
  [sortedPosts, isMobile],
  () => {
    if (isMobile.value) {
      nextTick(() => setupObserver())
    }
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

/* 移动端加载触发器样式 */
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

/* 响应式布局调整 */
@media (max-width: 768px) {
  .blog-layout {
    flex-direction: column;
    padding: 32px 16px 32px;
  }
  .main-content {
    margin-left: 0;
    margin-top: 0;
  }
}
</style>
