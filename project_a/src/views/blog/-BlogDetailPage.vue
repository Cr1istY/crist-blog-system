<!-- eslint-disable vue/multi-word-component-names -->
<!-- this page has been repalced by BlogDetailPageWithSlug.vue-->
<template>
  <div class="blog-detail-wrapper">
    <!-- 装饰性背景：左右淡色渐变，仅视觉增强 -->
    <div class="decorative-bg"></div>

    <div class="blog-detail-page">
      <n-spin :show="loading" description="加载中..." size="large">
        <!-- 错误提示 -->
        <n-alert v-if="error" type="error" show-icon style="margin-bottom: 24px">
          {{ error }}
          <template #action>
            <n-button text @click="retry">重试</n-button>
          </template>
        </n-alert>

        <!-- 文章内容 -->
        <n-card v-else-if="post" :bordered="false" class="post-card" embedded>
          <!-- 图片区 -->
          <div class="post-image">
            <n-image
              :src="processImageUrl(post.thumbnail)"
              :alt="post.title"
              object-fit="cover"
              preview-disabled
              width="100%"
              style="border-radius: 8px; margin-bottom: 24px"
            />
          </div>

          <!-- 标题区 -->
          <div class="post-header">
            <h1 class="post-title">{{ post.title }}</h1>
            <div class="post-meta">
              <n-tag type="info" size="small" round>{{ post.category }}</n-tag>
              <n-text type="secondary">{{ post.date }}</n-text>
              <div class="post-stats">
                <span class="stat-item">
                  <n-icon size="16" :component="EyeOutline" />
                  {{ formatNumber(post.views) }}
                </span>
                <span class="stat-item">
                  <n-icon size="16" :component="HeartOutline" />
                  {{ formatNumber(post.likes) }}
                </span>
              </div>
            </div>

            <!-- 标签云 -->
            <div v-if="post.tags && post.tags.length > 0" class="post-tags">
              <n-tag
                v-for="tag in post.tags"
                :key="tag"
                size="small"
                round
                type="success"
                @click="goToTag(tag)"
                style="cursor: pointer; margin-right: 8px; margin-top: 8px"
              >
                {{ tag }}
              </n-tag>
            </div>
          </div>

          <!-- Markdown 正文 -->
          <div class="markdown-container">
            <MdPreview
              v-model="post.content"
              :preview-only="true"
              :theme="isDark ? 'dark' : 'light'"
              language="zh-CN"
              :toolbars="[]"
              :editor-id="`post-${post.id}`"
              class="markdown-preview"
            />
          </div>

          <!-- 底部操作栏 -->
          <div class="post-footer">
            <n-space>
              <n-button
                type="tertiary"
                size="small"
                @click="router.push('/blog')"
                style="font-weight: 500"
              >
                ← 返回
              </n-button>
              <n-button type="tertiary" size="small" @click="sharePost" style="font-weight: 500">
                🤝 分享
              </n-button>
              <n-button
                :type="like_flag ? 'warning' : 'tertiary'"
                size="small"
                @click="likePost"
                style="font-weight: 500"
              >
                哎哟不错哟👍
              </n-button>
              <router-link v-if="update_flag" :to="`/admin/update/${post.id}`">
                <n-button type="tertiary" size="small" style="font-weight: 500"> 修改 </n-button>
              </router-link>
              <n-button
                v-if="update_flag"
                type="tertiary"
                size="small"
                @click="deletePost"
                style="font-weight: 500"
              >
                删除
              </n-button>
            </n-space>
          </div>
        </n-card>

        <!-- 返回顶部按钮 -->
        <n-back-top :right="24" :bottom="80" />
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { EyeOutline, HeartOutline } from '@vicons/ionicons5'

// 类型定义
interface BlogPost {
  id: number
  title: string
  category: string
  date: string
  content: string
  tags: string[]
  views: number
  likes: number
  excerpt: string
  meta_title?: string
  meta_description?: string
  thumbnail?: string
}

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const error = ref<string | null>(null)
const post = ref<BlogPost | null>(null)
const isDark = ref(false)

const update_flag = ref(false)
const like_flag = ref(false)
const token = localStorage.getItem('access_token')
if (token) {
  update_flag.value = true
}

// 格式化数字（如 1200 → 1.2k）
const formatNumber = (num: number): string => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 图片处理
const processImageUrl = (url?: string): string => {
  if (url) {
    return `/api/proxy/image?url=${encodeURIComponent(url)}`
  }
  return ""
}

// 获取文章
const fetchPost = async (id: string): Promise<void> => {
  loading.value = true
  error.value = null
  try {
    const response = await fetch(`/api/posts/get/${id}`)
    if (!response.ok) throw new Error('文章不存在或已删除')
    const data: BlogPost = await response.json()
    // 模拟数据放大
    data.views = data.views
    data.likes = data.likes
    post.value = data

    // 设置 SEO
    document.title = data.meta_title || `${data.title} - foreveryang`
    let metaDesc = document.querySelector<HTMLMetaElement>('meta[name="description"]')
    if (metaDesc) {
      metaDesc.content = data.meta_description || data.excerpt || ''
    } else {
      metaDesc = document.createElement('meta')
      metaDesc.name = 'description'
      metaDesc.content = data.meta_description || data.excerpt || ''
      document.head.appendChild(metaDesc)
    }

    // 添加views-出错不做处理
    await fetch(`/api/posts/addViews/${id}`)
  } catch (err) {
    const msg = err instanceof Error ? err.message : '未知错误'
    error.value = msg
    message.error(msg)
  } finally {
    loading.value = false
  }
}

// 重试
const retry = (): void => {
  const id = route.params.id as string | undefined
  if (id) {
    fetchPost(id)
  }
}

// 跳转到标签页
const goToTag = (tag: string): void => {
  router.push({ path: '/blog', query: { tag } })
}

// 分享文章
const sharePost = async (): Promise<void> => {
  const url = window.location.href
  try {
    if (navigator.share) {
      await navigator.share({
        title: post.value?.title || 'foreveryang',
        url,
      })
    } else {
      await navigator.clipboard.writeText(url)
      message.success('链接已复制到剪贴板')
    }
  } catch (err) {
    // 忽略用户取消分享（AbortError）
    if (err instanceof Error && err.name !== 'AbortError') {
      try {
        await navigator.clipboard.writeText(url)
        message.success('链接已复制到剪贴板')
      } catch {
        message.error('复制失败，请手动复制链接')
      }
    }
  }
}

const likePost = async (): Promise<void> => {
  if (!post.value?.id) return
  if (like_flag.value) return
  try {
    await fetch(`/api/posts/addLikes/${post.value?.id}`)
  } finally {
    like_flag.value = true
  }
}

const deletePost = async (): Promise<void> => {
  if (!post.value?.id) return
  dialog.warning({
    title: '删除文章',
    content: '确定要删除这篇文章吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        loading.value = true
        const response = await fetch(`/api/posts/delete/${post.value?.id}`, {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        })

        if (!response.ok) {
          throw new Error('删除失败')
        }
        message.success('文章已删除')
        router.push('/blog')
      } catch (error) {
        const msg = error instanceof Error ? error.message : '删除失败'
        message.error(msg)
      } finally {
        loading.value = false
      }
    },
  })
}

// 监听路由变化
watch(
  () => route.params.id,
  (newId) => {
    if (typeof newId === 'string') {
      fetchPost(newId)
    }
  },
  { immediate: true },
)

// 初始化主题（可选）
onMounted(() => {
  // isDark.value = localStorage.getItem('theme') === 'dark'
})
</script>

<style scoped>
.blog-detail-wrapper {
  position: relative;
  min-height: 100vh;
}

/* 装饰背景：左右淡色渐变 */
.decorative-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: -1;
  background: linear-gradient(
    to right,
    rgba(100, 100, 255, 0.03) 0%,
    transparent 15%,
    transparent 85%,
    rgba(100, 100, 255, 0.03) 100%
  );
}

.blog-detail-page {
  max-width: 1080px;
  margin: 0 auto;
  padding: 24px 24px;
  min-height: calc(100vh - 120px);
}

.post-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  background: var(--n-color);
}

.post-header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--n-border-color);
}

.post-title {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.4;
  margin: 0 0 16px;
  color: var(--n-title-text-color);
}

.post-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.post-stats {
  display: flex;
  gap: 16px;
  margin-left: auto;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--n-text-color);
  font-size: 14px;
}

.post-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.markdown-container {
  margin: 24px 0;
}

/* 覆盖 md-editor-v3 样式 */
.markdown-preview :deep(.md-editor-preview-wrapper) {
  padding: 0 !important;
  background: transparent !important;
  border: none !important;
}

.markdown-preview :deep(.markdown-body) {
  font-size: 16px;
  line-height: 1.8;
  color: var(--n-text-color);
  --md-primary-color: var(--n-primary-color);
  max-width: 800px;
  margin: 0 auto;
}

.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3) {
  margin-top: 1.5em;
  margin-bottom: 0.8em;
  font-weight: 600;
}

.markdown-preview :deep(p) {
  margin: 1em 0;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 16px 0;
}

.post-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--n-border-color);
}

.post-image {
  width: 100%;
  margin-bottom: 24px;
}

/* 响应式 */
@media (max-width: 768px) {
  .post-title {
    font-size: 24px;
  }

  .post-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .post-stats {
    margin-left: 0;
  }

  .stat-item {
    font-size: 13px;
  }

  .blog-detail-page {
    padding: 16px;
  }
}
</style>
