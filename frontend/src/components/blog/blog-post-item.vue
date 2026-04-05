<template>
  <div class="post-item" :class="{ 'is-pinned': showPin }">
    <div v-if="showPin" class="pin-badge">
      <n-tag type="error" size="small" round>置顶</n-tag>
    </div>

    <div class="post-content">
      <router-link :to="`/blog/${post.slug}`" class="post-title">
        {{ post.title }}
      </router-link>

      <div class="stats-container">
        <div class="stat-item">
          <n-icon size="16" :component="AccessibilityColor" />
          <span>{{ formatNumber(post.views) }}</span>
        </div>
        <div class="stat-item">
          <n-icon size="16" :component="HeartCircleOutline" />
          <span>{{ formatNumber(post.likes) }}</span>
        </div>
        <div class="stat-item">
          <n-icon size="16" :component="EarthSharp" />
          <span>{{ post.date }}</span>
        </div>
      </div>

      <div class="post-stats">
        <p class="post-excerpt">{{ post.excerpt }}</p>
        <div class="post-meta">
          <n-avatar
            v-if="post.user"
            size="small"
            round
            :src="post.user.avatar"
            style="margin-right: 0.8em"
          />
          <n-tag
            v-if="post.category && post.category !== ''"
            type="info"
            size="small"
            round
            @click.stop="$emit('cat-click', post.category)"
            style="cursor: pointer"
            >{{ post.category }}</n-tag
          >
          <n-avatar-group
            :options="tagOptions"
            :max="3"
            size="small"
            :style="{ display: 'inline-flex', gap: '4px', marginLeft: '8px' }"
          >
            <template #avatar="{ option }">
              <n-tooltip :delay="200">
                <template #trigger>
                  <n-avatar
                    :src="option.src"
                    :fallback="{ text: option.fallbackText }"
                    style="cursor: pointer"
                    @click.stop="$emit('tag-click', option.name)"
                  />
                </template>
                {{ option.name }}
              </n-tooltip>
            </template>

            <template #rest="{ rest, options }">
              <n-dropdown :options="dropdownOptions(options)" placement="top" trigger="hover">
                <n-avatar>+{{ rest }}</n-avatar>
              </n-dropdown>
            </template>
          </n-avatar-group>
        </div>
      </div>
    </div>

    <div class="post-thumbnail" v-if="thumbnailSrc">
      <img
        :src="thumbnailSrc"
        :alt="post.title"
        lazy
        @error="handleImageError"
        class="thumbnail-img"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { HeartCircleOutline, EarthSharp } from '@vicons/ionicons5'
import { AccessibilityColor } from '@vicons/carbon'
import type { BlogPost } from '@/types/blog'

const props = defineProps<{
  post: BlogPost
  showPin: boolean
}>()
const emit = defineEmits<{ 'tag-click': [tag: string]; 'cat-click': [cat: string] }>()

// 标签头像生成
const stringToColor = (str: string): string => {
  const colors = [
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
  if (!str || str.length === 0) return colors[0]!
  let hash = 0
  for (let i = 0; i < str.length; i++) hash = str.charCodeAt(i) + ((hash << 5) - hash)
  return colors[Math.abs(hash) % colors.length]!
}

const generateSvg = (text: string): string => {
  const char = text?.[0]?.toUpperCase() || '?'
  const bgColor = stringToColor(text || '')
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32">
    <circle cx="16" cy="16" r="16" fill="${bgColor}"/>
    <text x="16" y="21" text-anchor="middle" font-family="system-ui" font-size="16" fill="#334155" font-weight="500">${char}</text>
  </svg>`
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`
}

interface TagOption {
  name: string
  src: string
  fallbackText: string
}

const tagOptions = computed<TagOption[]>(() =>
  props.post.tags.map((tag) => ({
    name: tag,
    src: generateSvg(tag),
    fallbackText: tag.charAt(0).toUpperCase(),
  })),
)

const dropdownOptions = (opts: TagOption[]) =>
  opts.map((opt) => ({
    key: opt.name,
    label: () =>
      h('div', { style: 'display:flex;align-items:center;gap:8px' }, [
        h(NAvatar, { size: 'small', src: opt.src, fallback: { text: opt.fallbackText } }),
        h('span', opt.name),
      ]),
    onClick: () => emit('tag-click', opt.name),
  }))

// 图片处理
const defaultThumbnail =
  'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODAiIGhlaWdodD0iODAiIHZpZXdCb3g9IjAgMCA4MCA4MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjgwIiBoZWlnaHQ9IjgwIiBmaWxsPSIjZjNmNGY2Ii8+CjxwYXRoIGQ9Ik0yNSAyNUg1NVY1NUgyNVYyNVoiIHN0cm9rZT0iI2UyZThmMCIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KPHBhdGggZD0iTTM1IDQwSDQ1IiBzdHJva2U9IiNlMmU4ZjAiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIi8+CjxwYXRoIGQ9Ik0zNSA0NUg0NSIgc3Ryb2tlPSIjZTJlOGYwIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8cGF0aCBkPSJNMzUgNTBINTUiIHN0cm9rZT0iI2UyZThmMCIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiLz4KPC9zdmc+'
const thumbnailSrc = computed(() =>
  props.post.thumbnail
    ? `/api/proxy/image?url=${encodeURIComponent(props.post.thumbnail)}&quality=low`
    : defaultThumbnail,
)

const handleImageError = (e: Event) => {
  ;(e.target as HTMLImageElement).src = defaultThumbnail
}

// 工具函数
const formatNumber = (num: number) => (num >= 1000 ? `${(num / 1000).toFixed(1)}k` : num.toString())
</script>

<style scoped>
.post-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0 12px 12px;
  border-bottom: 1px solid #e5e7eb;
  border-radius: 12px;
  position: relative;
  transition: background-color 0.8s;
}
.post-item.is-pinned {
  background-color: #fff9f9;
}
.post-item:hover {
  background-color: rgba(233, 229, 226, 0.3);
}
.post-item.is-pinned:hover {
  background-color: #fff9f9;
}
.pin-badge {
  position: absolute;
  top: -6px;
  left: -6px;
  z-index: 10;
  transform: rotate(-12deg);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0);
}
.post-content {
  flex: 1;
  min-width: 0;
}
.post-title {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  text-decoration: none;
  display: block;
  margin-bottom: 6px;
  transition: color 0.2s;
}
.post-title:hover {
  color: #10b981;
}
.post-meta {
  display: flex;
  align-items: center;
  margin: 4px 0;
}
.post-excerpt {
  color: #333;
  line-height: 1.6;
  font-size: 0.95rem;
  margin: 8px 0;
}
.post-thumbnail {
  width: 176px;
  height: 99px;
  flex-shrink: 0;
  border-radius: 8px;
  overflow: hidden;
  background-color: #f8fafc;
  margin-right: 1em;
}
.thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}
.thumbnail-img:hover {
  transform: scale(1.05);
}
.stats-container {
  display: flex;
  gap: 0.8rem;
  margin-top: 8px;
}
.stat-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: #666;
  font-size: 0.85rem;
}
.stat-item:hover {
  color: #333;
}
.stat-item:hover :deep(.n-icon) {
  opacity: 1;
}
@media (max-width: 640px) {
  .post-item {
    flex-direction: column;
    align-items: start;
  }
  .post-item.is-pinned {
    padding-left: 14px;
  }
  .pin-badge {
    top: -4px;
    left: 4px;
    transform: rotate(-8deg);
  }
  .post-thumbnail {
    width: 100%;
    height: 180px;
    order: -1;
    border-radius: 10px;
  }
}
</style>
