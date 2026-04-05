<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="blog-detail-wrapper">
    <div class="blog-detail-page">
      <n-spin :show="loading" description="加载中..." size="large" style="margin-top: 24px">
        <!-- 错误提示 -->
        <n-alert v-if="error" type="error" show-icon style="margin-top: 24px">
          {{ error }}
          <template #action>
            <n-button text @click="retry">重试</n-button>
          </template>
        </n-alert>

        <!-- 文章内容 -->
        <n-card v-else-if="post" :bordered="false" class="post-card">
          <!-- 图片区 -->
          <div class="post-image">
            <n-image
              :src="processImageUrl(post.thumbnail)"
              :alt="post.title"
              object-fit="cover"
              preview-disabled
              width="100%"
              style="border-radius: 8px"
            />
          </div>

          <!-- 标题区 -->
          <div class="post-header">
            <h1 class="post-title">{{ post.title }}</h1>

            <div class="post-description-bootom">
              <!-- 标签云 -->
              <div v-if="post.tags && post.tags.length > 0" class="post-tags">
                <n-tag
                  v-for="tag in post.tags"
                  :key="tag"
                  size="small"
                  round
                  type="success"
                  @click="goToTag(tag)"
                  class="tag-item"
                >
                  {{ tag }}
                </n-tag>
              </div>
              <!--description-->
              <div class="post-excerpt" v-if="post.excerpt && post.excerpt.length > 0">
                <n-highlight
                  :text="text"
                  :patterns="patterns"
                  :highlight-style="{
                    padding: '0 6px',
                    borderRadius: themeVars.borderRadius,
                    display: 'inline-block',
                    color: themeVars.baseColor,
                    background: themeVars.primaryColor,
                    transition: `all .3s ${themeVars.cubicBezierEaseInOut}`,
                  }"
                />
              </div>
            </div>
            <div class="post-meta">
              <n-tag
                type="info"
                size="small"
                round
                style="cursor: pointer"
                @click="goToCat(post.category)"
                >{{ post.category }}</n-tag
              >
              <div class="post-stats">
                <n-avatar v-if="post.user_avatar" round size="small" :src="post.user_avatar" />
                <span class="stat-item">
                  <n-icon size="16" :component="EyeOutline" />
                  {{ formatNumber(post.views) }}
                </span>
                <span class="stat-item">
                  <n-icon size="16" :component="HeartOutline" />
                  {{ formatNumber(post.likes) }}
                </span>
                <span class="stat-item">
                  <n-icon size="16" :component="ChatbubbleOutline" />
                  Edited on {{ post.date }}
                </span>
              </div>
            </div>
          </div>

          <!-- Markdown 正文 -->
          <div class="markdown-container">
            <MdPreview
              v-model="post.content"
              :preview-only="true"
              preview-theme="vuepress"
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
              <n-button
                v-if="update_flag && !isPinned"
                type="tertiary"
                size="small"
                @click="pinPost"
                style="font-weight: 500"
              >
                置顶
              </n-button>
              <n-button
                v-if="update_flag && isPinned"
                type="tertiary"
                size="small"
                @click="unpinPost"
                style="font-weight: 500"
              >
                取消置顶
              </n-button>
            </n-space>
          </div>
        </n-card>
      </n-spin>
    </div>
  </div>
  <!-- 返回顶部按钮 -->
  <n-back-top :right="24" :bottom="80" />
  <n-button
    class="sidebar-toggle-float-btn"
    strong
    secondary
    circle
    @click="toggleDrawer"
    title="切换侧边栏"
  >
    <template #icon>
      <n-icon :component="MenuOutline" />
    </template>
  </n-button>
  <n-drawer v-model:show="active" :width="502" :placement="placement">
    <n-drawer-content title="目录" closable>
      <MdCatalog
        v-if="post"
        :editor-id="`post-${post.id}`"
        :scrollElement="scrollElement"
        :theme="isDark ? 'dark' : 'light'"
      />
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage, useDialog, useThemeVars } from 'naive-ui'
import type { DrawerPlacement } from 'naive-ui'
import { MdPreview, MdCatalog } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import { EyeOutline, HeartOutline, ChatbubbleOutline, MenuOutline } from '@vicons/ionicons5'
import service from '@/utils/request'

// 类型定义
interface BlogPost {
  id: number
  slug: string
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
  is_pinned?: boolean
  user_avatar?: string
}

const route = useRoute()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()

const loading = ref(false)
const error = ref<string | null>(null)
const post = ref<BlogPost | null>(null)
const isDark = ref(false)

const isPinned = ref(false)

const update_flag = ref(false)
const like_flag = ref(false)
const token = localStorage.getItem('access_token')

const scrollElement = document.documentElement

if (token) {
  update_flag.value = true
}

// 处理高亮
const text = ref('')
const patterns = ref([''])
const themeVars = useThemeVars()

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
    return `/api/proxy/image?url=${encodeURIComponent(url)}&quality=medium`
  }
  return ''
}

// 获取文章
const fetchPost = async (slug: string): Promise<void> => {
  loading.value = true
  error.value = null
  try {
    const response = await fetch(`/api/posts/getBySlug/${slug}`)
    if (!response.ok) throw new Error('文章不存在或已删除')
    const data: BlogPost = await response.json()
    post.value = data
    // 高亮处理
    text.value = post.value.excerpt ?? ''
    patterns.value = post.value.tags ?? []
    isPinned.value = data.is_pinned || false
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
    ;(async () => {
      try {
        await fetch(`/api/posts/addViews/${post.value?.id}`)
      } catch (err) {
        // 静默失败，或者只打印日志，不提示用户，不影响主流程
        console.warn('增加浏览量失败:', err)
      }
    })()
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
const goToCat = (cat: string): void => {
  cat = cat.split('/').pop() ?? cat
  router.push({ path: '/blog', query: { cat } })
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
        if (!post.value) {
          message.error('文章不存在')
          return
        }
        await service.delete(`/posts/delete/${post.value.id}`)
        message.success('文章已删除')
        router.push('/blog')
      } catch (error) {
        const msg = error instanceof Error ? error.message : '删除失败'
        message.error(msg)
        if (msg.includes('401')) {
          router.push('/admin')
        }
      } finally {
        loading.value = false
      }
    },
  })
}

const pinPost = async (): Promise<void> => {
  if (!post.value?.id) return
  dialog.warning({
    title: '置顶文章',
    content: '确定要置顶这篇文章吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        loading.value = true
        if (!post.value) {
          message.error('文章不存在')
          return
        }
        await service.put(`/posts/pin/${post.value.id}`)
        message.success('文章已删除')
        router.push('/blog')
      } catch (error) {
        const msg = error instanceof Error ? error.message : '删除失败'
        message.error(msg)
        if (msg.includes('401')) {
          router.push('/admin')
        }
      } finally {
        loading.value = false
      }
    },
  })
}

const unpinPost = async (): Promise<void> => {
  if (!post.value?.id) return
  dialog.warning({
    title: '取消置顶文章',
    content: '确定要取消置顶这篇文章吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        loading.value = true
        if (!post.value) {
          message.error('文章不存在')
          return
        }
        await service.put(`/posts/unpin/${post.value.id}`)
        message.success('文章已取消置顶')
        router.push('/blog')
      } catch (error) {
        const msg = error instanceof Error ? error.message : '取消置顶失败'
        message.error(msg)
        if (msg.includes('401')) {
          router.push('/admin')
        }
      } finally {
        loading.value = false
      }
    },
  })
}

const active = ref(false)
const placement = ref<DrawerPlacement>('right')
function toggleDrawer() {
  active.value = true
}

// 监听路由变化
watch(
  () => route.params.slug,
  (newSlug) => {
    if (typeof newSlug === 'string') {
      fetchPost(newSlug)
    }
  },
  { immediate: true },
)

// 初始化主题（可选）
onMounted(() => {
  isDark.value = localStorage.getItem('theme') === 'dark'
})
</script>

<style scoped>
.blog-detail-wrapper {
  position: relative;
  min-height: 100vh;
}

.blog-detail-page {
  max-width: 1080px;
  margin: 0 auto;
  min-height: calc(100vh - 120px);
}

.post-card {
  border-radius: 24px;
  background: var(--n-color);
}

.post-header {
  margin-top: 1cm;
  padding-bottom: 1em;
  border-bottom: 1px solid var(--n-border-color);
}

.post-title {
  display: none;
  font-size: 28px;
  font-weight: 700;
  line-height: 1.4;
  margin: 16px 0;
  color: var(--n-title-text-color);
}

.post-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 1em;
}

.post-stats {
  display: flex;
  gap: 16px;
  margin-left: auto;
  margin-right: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--n-text-color);
  font-size: 14px;
}

.post-description-bootom {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.post-tags {
  gap: 8px;
}

.tag-item {
  cursor: pointer;
  margin-right: 8px;
  transition:
    background-color 0.5s cubic-bezier(0.4, 0, 0.2, 1),
    color 0.5s cubic-bezier(0.4, 0, 0.2, 1),
    transform 0.5s cubic-bezier(0.4, 0, 0.2, 1),
    box-shadow 0.5s cubic-bezier(0.4, 0, 0.2, 1),
    border-color 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateY(0);
}

.tag-item:hover {
  transform: translateY(-2px);
  background-color: #059669 !important;
  color: #ffffff !important;
  border-color: #059669 !important;
  box-shadow:
    0 10px 15px -3px rgba(5, 150, 105, 0.4),
    0 4px 6px -2px rgba(5, 150, 105, 0.2);
}

.tag-item:active {
  transform: translateY(-1px) scale(0.98);
  box-shadow: 0 4px 6px -1px rgba(5, 150, 105, 0.3);
  transition-duration: 0.1s;
}

.post-excerpt {
  margin-left: auto;
  flex-shrink: 0;
}

.post-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--n-border-color);
}

.post-image {
  max-height: 50vh;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 8px;
  overflow: hidden;
}

.post-image :deep(img),
.post-image :deep(.n-image__img) {
  max-width: 100%;
  height: auto;
  transition: transform 2s ease;
}

.post-image:hover :deep(img) {
  border-radius: 8px;
  transform: scale(1.04);
}

.sidebar-toggle-float-btn {
  position: fixed;
  height: 44px;
  width: 44px;
  right: 24px;
  bottom: 140px; /* 位于返回顶部按钮 (bottom: 80) 之上 */
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.sidebar-toggle-float-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}
/* 响应式 */
@media (max-width: 768px) {
  .post-header {
    margin-bottom: 0;
    padding-bottom: 0;
    border-bottom: 1px solid var(--n-border-color);
  }

  .post-title {
    font-size: 24px;
  }

  .post-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 2em;
    padding-bottom: 1em;
  }

  .post-stats {
    margin-left: 0;
  }

  .post-excerpt {
    margin-top: 2em;
    margin-left: auto;
    margin-bottom: 1em;
    flex-shrink: 1;
    word-break: break-word;
    overflow-wrap: break-word;
    display: block;
  }

  .stat-item {
    font-size: 13px;
  }

  .blog-detail-page {
    padding: 16px;
  }
}
</style>
