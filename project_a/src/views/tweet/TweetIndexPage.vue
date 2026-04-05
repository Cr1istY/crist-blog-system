<script setup lang="ts">
defineOptions({
  name: 'TweetIndexPage',
})

import { ref, onMounted, getCurrentInstance } from 'vue'
import Toast from 'primevue/toast'
import TweetComposer from '@/components/tweet/TweetComposer.vue'
import TweetList from '@/components/tweet/tweet-list.vue'
import service from '@/utils/request'
import type { Tweet, User, TweetImage } from '@/types/tweet'

// --- 状态定义 ---
const tweets = ref<Tweet[]>([])
const isLoading = ref<boolean>(false) // 初始化为 false，由 loadTweets 控制
const isTweetLoading = ref<boolean>(false)
const hasMore = ref<boolean>(true) // 标记是否还有更多数据
const instance = getCurrentInstance()
const isLogin = ref<boolean>(false)
const currentUser = ref<User | null>(null)

// --- 分页配置 ---
const LIMIT = 20
let offset = 0 // 记录当前已加载的数量，作为下一次请求的 offset

// --- 工具函数 ---
const showToast = (
  severity: 'success' | 'info' | 'warn' | 'error',
  summary: string,
  detail?: string,
  life: number = 3000,
): void => {
  instance?.proxy?.$toast.add({
    severity,
    summary,
    detail,
    life,
  })
}

// --- 核心加载函数 ---
/**
 * @param isLoadMore - true 表示加载更多(追加), false 表示刷新/首次加载(重置)
 */
const loadTweets = async (isLoadMore: boolean = false): Promise<void> => {
  // 如果正在加载中，或者没有更多数据了却尝试加载更多，则直接返回
  if (isLoading.value || (isLoadMore && !hasMore.value)) {
    return
  }

  // 如果是刷新操作（非加载更多），重置状态
  if (!isLoadMore) {
    offset = 0
    tweets.value = []
    hasMore.value = true
  }

  isLoading.value = true

  try {
    // 模拟网络延迟
    if (!isLoadMore) {
      await new Promise((resolve) => setTimeout(resolve, 200))
    }

    // 构建带参数的 URL
    const url = new URL('/api/tweet/getall', window.location.origin)
    url.searchParams.append('limit', LIMIT.toString())
    url.searchParams.append('offset', offset.toString())

    const res = await fetch(url.toString())

    if (!res.ok) {
      throw new Error(`API error: ${res.status}`)
    }

    const apiTweets: Tweet[] = await res.json()

    // 映射数据 (保持原有逻辑)
    const newTweets: Tweet[] = apiTweets.map((tweet) => ({
      id: tweet.id,
      user: tweet.user,
      content: tweet.content,
      timestamp: new Date(tweet.timestamp), // 确保是 Date 对象，如果后端传的是字符串
      likes: tweet.likes,
      images: tweet.images,
    }))

    // 判断是否还有更多数据
    // 如果返回的数量少于请求的 limit，说明到底了
    if (newTweets.length < LIMIT) {
      hasMore.value = false
    }

    if (isLoadMore) {
      // 追加模式
      tweets.value = [...tweets.value, ...newTweets]
      offset += LIMIT // 更新 offset 供下一次使用
    } else {
      // 重置模式
      tweets.value = newTweets
      offset = LIMIT // 下一次从 LIMIT 开始
    }
  } catch (error) {
    console.error('加载推文失败:', error)
    showToast('error', '加载失败', '无法获取推文列表，请稍后重试')
    // 如果出错且是加载更多，可能需要回滚 offset，这里简单处理不回滚，用户可重试
    if (!isLoadMore) {
      offset = 0
    }
  } finally {
    isLoading.value = false
  }
}

// --- 用户相关 ---
const getCurrentUser = async (): Promise<User | null> => {
  try {
    const token = localStorage.getItem('access_token')
    if (!token) {
      return null
    }
    isLogin.value = true
    // 假设 service 是 axios 实例，且拦截器已处理 token
    const res = await service.get('/tweet/getCurrentUser')
    // 根据实际后端返回结构调整，假设 data 字段包含用户信息
    const userData = res.data?.data || res.data
    if (userData) {
      isLogin.value = true
      return userData as User
    }
    return null
  } catch (error) {
    console.log('获取用户信息失败:', error)
    isLogin.value = false
    return null
  }
}

// --- 业务逻辑 ---
const handleTweetSubmit = async (content: string, images: TweetImage[]): Promise<void> => {
  if (!isLogin.value || !currentUser.value) {
    showToast('error', '错误', '请先登录')
    return
  }

  isTweetLoading.value = true
  try {
    const imageUrls: string[] = images.map((img) => img.url)
    const newTweet: Tweet = {
      id: `tweet-${Date.now()}`,
      user: currentUser.value,
      content,
      timestamp: new Date(),
      likes: 0,
      images: imageUrls.length > 0 ? imageUrls : undefined,
    }

    // 插入到列表最前面
    tweets.value.unshift(newTweet)

    const imageCount = images.length
    showToast(
      'success',
      '发布成功',
      imageCount > 0 ? `已发布推文和 ${imageCount} 张图片` : '您的推文已成功发布',
    )
  } catch (error) {
    showToast('error', '发布失败', '请稍后重试')
    console.error('发布推文失败:', error)
  } finally {
    isTweetLoading.value = false
  }
}

const handleLike = (tweetId: string): void => {
  const tweet = tweets.value.find((t) => t.id === tweetId)
  if (tweet) {
    tweet.likes += 1
    // 实际项目中应调用 API 更新点赞状态
  }
}

const removeTweet = async (id: string | number): Promise<void> => {
  if (!confirm('确定要删除这条推文吗？')) {
    return
  }

  try {
    // 调用删除 API
    const res = await service.delete(`/tweet/delete/${id}`)

    // 假设后端返回 code 200 表示成功
    if (res.data?.code === 200 || res.status === 200) {
      // 从本地数组移除
      tweets.value = tweets.value.filter((t) => t.id !== id)
      showToast('success', '删除成功', '推文已删除')
    } else {
      throw new Error('Delete failed')
    }
  } catch (error) {
    console.error('删除推文失败:', error)
    showToast('error', '删除失败', '请稍后重试')
  }
}

// --- 生命周期 ---
onMounted(async (): Promise<void> => {
  // 并行执行：加载推文和获取用户信息
  await Promise.all([
    loadTweets(false), // 首次加载
    (async () => {
      const user = await getCurrentUser()
      if (user) {
        currentUser.value = user
      }
    })(),
  ])
})
</script>

<template>
  <div class="app-container">
    <Toast />

    <header class="app-header">
      <div class="header-content">
        <i class="pi pi-box logo"></i>
        <h1>Threads</h1>
      </div>
    </header>

    <main class="main-content">
      <div class="feed-container">
        <TweetComposer
          v-if="isLogin && currentUser"
          :current-user="currentUser"
          @submit="handleTweetSubmit"
          :loading="isTweetLoading"
        />

        <TweetList
          :tweets="tweets"
          :loading="isLoading && tweets.length === 0"
          :current-user="currentUser || undefined"
          @like="handleLike"
          @delete-tweet="removeTweet"
        />

        <!-- 加载更多按钮 / 加载状态指示器 -->
        <div class="load-more-container">
          <button v-if="hasMore && !isLoading" @click="loadTweets(true)" class="load-more-btn">
            加载更多推文
          </button>

          <div v-if="isLoading && tweets.length > 0" class="loading-spinner">
            <i class="pi pi-spin pi-spinner"></i>
            <span>加载中...</span>
          </div>

          <div v-if="!hasMore && tweets.length > 0" class="no-more-text">没有更多推文了</div>
        </div>
      </div>

      <aside class="sidebar">
        <div class="sidebar-card">
          <h3>当前专注</h3>
          <div class="trend-item">
            <span class="trend-category">核心目标</span>
            <span class="trend-title">研究生入学考试备考</span>
            <span class="trend-count"></span>
          </div>

          <div class="trend-item">
            <span class="trend-category">自我提升</span>
            <span class="trend-title">深度阅读计划</span>
            <span class="trend-count"></span>
          </div>

          <div class="trend-item">
            <span class="trend-category">技术探索</span>
            <span class="trend-title">cSharp mod</span>
            <span class="trend-count"></span>
          </div>
        </div>
      </aside>
    </main>
  </div>
</template>

<style>
/* 全局样式 */
:root {
  --primary-color: #1d9bf0;
  --surface-border: #eff3f4;
  --surface-hover: #f7f9f9;
  --text-color: #0f1419;
  --text-color-secondary: #536471;
  --red-500: #f91880;
  --green-500: #00ba7c;
  --blue-500: #1d9bf0;
  --orange-500: #f7931a;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  background-color: #ffffff;
  color: var(--text-color);
}
</style>

<style scoped>
.app-container {
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
}

.app-header {
  position: sticky;
  top: 0;
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--surface-border);
  z-index: 100;
}

.header-content {
  max-width: 600px;
  margin: 0 auto;
  padding: 1rem;
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.logo {
  font-size: 2rem;
  color: var(--p-orange-600);
  /* 应用动画：名称 时长 缓动函数 无限循环 */
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0); /* 初始位置和结束位置 */
  }
  50% {
    transform: translateY(-5px); /* 移动到上方 10px */
  }
}

.header-content h1 {
  font-size: 1.25rem;
  font-weight: 700;
}

.main-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 2rem;
  max-width: 900px;
  margin: 0 auto;
  padding: 1rem;
}

.feed-container {
  background-color: #ffffff;
  border-radius: 1rem;
  border: 1px solid var(--surface-border);
  display: flex;
  flex-direction: column;
}

/* 加载更多区域样式 */
.load-more-container {
  padding: 1.5rem;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  gap: 1rem;
  border-top: 1px solid var(--surface-border);
  margin-top: auto;
}

.load-more-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 2rem;
  border-radius: 9999px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 0.95rem;
}

.load-more-btn:hover {
  background-color: #1a8cd8;
}

.load-more-btn:active {
  transform: scale(0.98);
}

.loading-spinner {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-color-secondary);
  font-size: 0.9rem;
}

.no-more-text {
  color: var(--text-color-secondary);
  font-size: 0.9rem;
  font-style: italic;
}

.sidebar {
  position: sticky;
  top: 5rem;
  height: fit-content;
}

.sidebar-card {
  background-color: var(--surface-hover);
  border-radius: 1rem;
  padding: 1rem;
}

.sidebar-card h3 {
  font-size: 1.25rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.trend-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 0.75rem 0;
  cursor: pointer;
  border-bottom: 1px solid var(--surface-border);
}

.trend-item:last-child {
  border-bottom: none;
}

.trend-category {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
}

.trend-title {
  font-weight: 700;
  color: var(--text-color);
}

.trend-count {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
}

@media (max-width: 768px) {
  .main-content {
    grid-template-columns: 1fr;
  }

  .sidebar {
    display: none;
  }
}
</style>
