<script setup lang="ts">
import { NCard, NAvatar, NDivider, NButton, NImage, NIcon, NSpin, NImageGroup } from 'naive-ui'
import { TrashOutline, CheckmarkCircle, PinOutline, TimeOutline } from '@vicons/ionicons5'
import type { TweetListProps } from '@/types/tweet'

withDefaults(defineProps<TweetListProps>(), {
  loading: false,
  tweets: () => [],
})

const emit = defineEmits<{
  (e: 'delete-tweet', id: string | number): void
}>()

const formatTimestamp = (date: Date): string => {
  const now = new Date()
  const diff = now.getTime() - new Date(date).getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`

  return new Date(date).toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
  })
}

const handleDelete = (id: string | number): void => {
  emit('delete-tweet', id)
}
</script>

<template>
  <div class="page-container">
    <div class="tweet-list">
      <!-- 加载状态 -->
      <div v-if="loading" class="status-state">
        <n-spin size="large" />
        <span class="status-text">加载中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="tweets.length === 0" class="status-state">
        <n-icon :component="PinOutline" size="48" color="var(--n-color-text-disabled)" />
        <span class="status-text">暂无推文</span>
      </div>

      <!-- 列表 -->
      <!-- 使用 NImageGroup 包裹，实现组内图片预览联动 -->
      <n-image-group>
        <n-card
          v-for="tweet in tweets"
          :key="tweet.id"
          class="tweet-card"
          hoverable
          content-style="padding: 0;"
          :bordered="false"
        >
          <template #header>
            <div class="tweet-header">
              <n-avatar
                :src="tweet.user.avatar"
                round
                size="large"
                :fallback-src="'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg'"
              />
              <div class="tweet-user-info">
                <div class="user-name-row">
                  <span class="display-name">{{ tweet.user.displayName }}</span>
                  <n-icon
                    v-if="tweet.user.verified"
                    :component="CheckmarkCircle"
                    class="verified-badge"
                    color="#2080f0"
                  />
                </div>
                <span class="username">@{{ tweet.user.username }}</span>
              </div>
              <div class="tweet-timestamp">
                <n-icon :component="TimeOutline" style="margin-right: 4px; vertical-align: middle" />
                {{ formatTimestamp(tweet.timestamp) }}
              </div>
            </div>
          </template>

          <template #default>
            <div class="tweet-body">
              <p class="tweet-content">{{ tweet.content }}</p>

              <!-- 图片区域 -->
              <div
                v-if="tweet.images && tweet.images.length > 0"
                class="tweet-images"
                :class="`image-count-${Math.min(tweet.images.length, 9)}`"
              >
                <n-image
                  v-for="(image, index) in tweet.images"
                  :key="index"
                  :src="image"
                  class="tweet-image"
                  object-fit="cover"
                  preview-toolbar-placement="bottom-center"
                  :img-props="{ class: 'custom-img-inner' }"
                />
              </div>
            </div>

            <n-divider style="margin: 0" />

            <div class="tweet-actions">
              <n-button
                v-if="currentUser && tweet.user.id === currentUser.id"
                quaternary
                circle
                type="error"
                @click="handleDelete(tweet.id)"
              >
                <template #icon>
                  <n-icon :component="TrashOutline" />
                </template>
              </n-button>
            </div>
          </template>
        </n-card>
      </n-image-group>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  background-color: var(--n-color);
  min-height: 100vh;
  padding-bottom: 20px;
}

.tweet-list {
  display: flex;
  flex-direction: column;
}

.status-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: var(--n-color-text-disabled);
  gap: 1rem;
}

.tweet-card {
  border-radius: 0;
  /* 朋友圈风格通常没有明显的卡片边框，只有底部分割线 */
  border-bottom: 8px solid var(--n-color-hover);
  box-shadow: none;
  transition: background-color 0.2s ease;
  width: 100%;
}

/* 移除 header 的默认 padding，自定义 */
:deep(.n-card__header) {
  padding: 12px 16px;
  border-bottom: none;
}

.tweet-header {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.tweet-user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.user-name-row {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.display-name {
  font-weight: 600;
  color: var(--n-color-text);
  font-size: 15px;
}

.verified-badge {
  font-size: 14px;
}

.username {
  color: var(--n-color-text-3);
  font-size: 13px;
  margin-top: 2px;
}

.tweet-timestamp {
  color: var(--n-color-text-3);
  font-size: 12px;
  display: flex;
  align-items: center;
  white-space: nowrap;
  margin-left: auto; /* 时间靠右 */
}

.tweet-body {
  padding: 0 16px 12px;
}

.tweet-content {
  margin: 0;
  font-size: 15px;
  line-height: 1.6;
  color: var(--n-color-text);
  white-space: pre-wrap;
  word-break: break-word;
}

/* --- 朋友圈风格图片布局核心 --- */
.tweet-images {
  display: grid;
  gap: 6px; /* 图片间距 */
  margin-top: 12px;
  width: 100%;
  box-sizing: border-box;
}

.tweet-image {
  border-radius: 4px;
  overflow: hidden;
  background-color: #f0f0f0;
  cursor: zoom-in;
  /* 关键：让图片填满网格单元 */
  width: 100%;
  height: 100%;
  display: block;
}

/* 内部 img 标签属性，确保 cover 生效 */
:deep(.custom-img-inner) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* --- 场景 1：单图 (1) --- */
/* 朋友圈单图通常不强制正方形，而是显示原图比例，但限制最大尺寸 */
.image-count-1 {
  display: block; /* 退出 grid，使用 block 流 */
  max-width: 70%; /* 限制最大宽度，避免太长 */
  max-height: 400px; /* 限制最大高度 */
}

.image-count-1 .tweet-image {
  /* 单图不强制 aspect-ratio，让图片自然撑开，但受限于父级 max-width/height */
  aspect-ratio: auto;
  /* 如果希望单图也是正方形裁剪，取消下面这行的注释并删除上面的 max-width/max-height 逻辑 */
  /* aspect-ratio: 1/1; */
}

/* --- 场景 2：双图 (2) --- */
/* 两列等分，高度自适应或固定比例，朋友圈通常是 1:1 或 16:9，这里采用 1:1 保证整齐 */
.image-count-2 {
  grid-template-columns: repeat(2, 1fr);
}

.image-count-2 .tweet-image {
  aspect-ratio: 1 / 1;
}

/* --- 场景 3：三图及以上 (3-9) --- */
/* 经典的 3 列九宫格，所有图片强制 1:1 */
.image-count-3,
.image-count-4,
.image-count-5,
.image-count-6,
.image-count-7,
.image-count-8,
.image-count-9 {
  grid-template-columns: repeat(3, 1fr);
}

.image-count-3 .tweet-image,
.image-count-4 .tweet-image,
.image-count-5 .tweet-image,
.image-count-6 .tweet-image,
.image-count-7 .tweet-image,
.image-count-8 .tweet-image,
.image-count-9 .tweet-image {
  aspect-ratio: 1 / 1;
}

/* 特殊处理：如果是 4 张图，有些设计喜欢 2x2，这里保持 3 列布局（第一行3张，第二行1张）
   如果你更喜欢 2x2，可以单独给 .image-count-4 设置 grid-template-columns: repeat(2, 1fr);
   但在微信朋友圈中，4张图通常是第一行3张，第二行1张（左对齐），或者全部 1:1 排列。
   此处保持通用 3 列流式布局。
*/

.tweet-actions {
  display: flex;
  justify-content: flex-end;
  padding: 8px 16px 12px;
}

/* 移动端适配 */
@media (max-width: 480px) {
  .page-container {
    max-width: 100%;
  }

  .image-count-1 {
    max-width: 85%; /* 手机上单图更宽 */
  }

  /* 手机上如果图片太多，3列可能太小，可根据需求调整为2列，但朋友圈通常保持3列 */
  /* 如果希望手机上多图变2列，解开下面注释 */
  /*
  .image-count-3, .image-count-4, .image-count-5,
  .image-count-6, .image-count-7, .image-count-8, .image-count-9 {
    grid-template-columns: repeat(2, 1fr);
  }
  */
}
</style>
