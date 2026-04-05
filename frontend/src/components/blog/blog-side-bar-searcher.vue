<!-- src/components/SearchBar.vue -->
<template>
  <div class="search-container">
    <div class="search-wrapper" :class="{ 'is-focused': showSuggestions }">
      <n-input
        ref="inputRef"
        v-model:value="localSearch"
        placeholder="搜索文章标题..."
        clearable
        size="small"
        round
        @update:value="onSearchInput"
        @focus="handleFocus"
        @blur="handleBlur"
        class="custom-search-input"
      >
        <template #prefix>
          <span class="search-icon">🔍</span>
        </template>
      </n-input>

      <transition name="slide-fade">
        <div
          v-if="showSuggestions && suggestions.length"
          class="search-suggestions"
          :class="{ 'mobile-mode': isMobile }"
          @mouseenter="isMouseOverSuggestions = true"
          @mouseleave="isMouseOverSuggestions = false"
        >
          <div class="suggestions-header" :class="{ visible: isMobile || showSuggestions }">
            <span class="header-title">
              {{ localSearch.trim() ? '🔍 搜索结果' : '🔥 热门推荐' }}
            </span>
            <n-button v-if="isMobile" text size="tiny" @click="closeSuggestions"> 关闭 </n-button>
          </div>

          <div
            v-for="(s, i) in suggestions"
            :key="i"
            class="suggestion-item"
            @click="selectSuggestion(s)"
            role="button"
            tabindex="0"
            @keydown.enter="selectSuggestion(s)"
          >
            <span v-if="!localSearch.trim()" class="item-icon emoji-icon">✨</span>
            <n-icon v-else class="item-icon"><ArrowForwardCircleOutline /></n-icon>

            <span class="item-text" v-html="highlightMatch(s, localSearch)"></span>
          </div>

          <div v-if="suggestions.length >= 5" class="suggestions-footer">
            <n-text depth="3" style="font-size: 12px">
              {{ localSearch.trim() ? '仅显示前 5 条结果' : '更多请输入关键词' }}
            </n-text>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { ArrowForwardCircleOutline } from '@vicons/ionicons5'
import type { BlogPost } from '@/types/blog'
import { type InputInst } from 'naive-ui'

interface Props {
  modelValue?: string
  posts: BlogPost[]
}

interface Emits {
  'update:modelValue': [keyword: string]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
})
const emit = defineEmits<Emits>()

const inputRef = ref<InputInst | null>(null)
const localSearch = ref<string>(props.modelValue)
const showSuggestions = ref<boolean>(false)
const isMouseOverSuggestions = ref<boolean>(false)
const suggestions = ref<string[]>([])
const keywordCache = ref<Set<string>>(new Set())
const isMobile = ref<boolean>(false)

// 检测是否为移动设备
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

// --- 搜索建议相关逻辑 ---

const updateKeywordCache = (): void => {
  const newCache = new Set<string>()
  props.posts.forEach((post) => {
    if (post.title) {
      newCache.add(post.title)
    }
  })
  keywordCache.value = newCache

  // 如果当前是聚焦状态且没有输入内容，缓存更新后重新显示推荐
  if (showSuggestions.value && !localSearch.value.trim()) {
    computeSuggestions('')
  }
}

const computeSuggestions = (query: string): void => {
  const allKeywords = Array.from(keywordCache.value)

  // 情况 A: 有搜索词 -> 进行过滤匹配
  if (query.trim()) {
    const lowerQuery = query.toLowerCase()
    suggestions.value = allKeywords
      .filter((kw) => kw.toLowerCase().includes(lowerQuery))
      .slice(0, 5)
  }
  // 情况 B: 无搜索词 (点击空输入框) -> 显示 5 个热门推荐
  else {
    const allKeywords = Array.from(keywordCache.value)
    // 复制数组并打乱顺序
    const shuffled = [...allKeywords].sort(() => Math.random() - 0.5)
    // 取前 5 个
    suggestions.value = shuffled.slice(0, 5)
  }
}

// 高亮匹配文字
const highlightMatch = (text: string, query: string) => {
  if (!query.trim()) return text // 无搜索词时不高亮
  const regex = new RegExp(`(${query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

type DebounceFunction<T extends unknown[]> = (...args: T) => void

const debounce = <T extends unknown[]>(
  fn: DebounceFunction<T>,
  delay: number,
): DebounceFunction<T> => {
  let timer: number | null = null
  return (...args: T) => {
    if (timer) clearTimeout(timer)
    timer = window.setTimeout(() => fn(...args), delay) as unknown as number
  }
}

const updateSuggestionsDebounced = debounce((query: string) => {
  computeSuggestions(query)
}, 250)

const onSearchInput = (val: string): void => {
  emit('update:modelValue', val)
  updateSuggestionsDebounced(val)
}

const selectSuggestion = (text: string): void => {
  localSearch.value = text
  emit('update:modelValue', text)
  closeSuggestions()

  // 移动端选中后让输入框失去焦点，收起键盘
  if (isMobile.value && inputRef.value) {
    inputRef.value.blur()
  }
}

const handleFocus = () => {
  showSuggestions.value = true
  // 聚焦时立即计算建议（无论是空还是已有文字）
  computeSuggestions(localSearch.value)
}

const handleBlur = () => {
  nextTick(() => {
    if (!isMouseOverSuggestions.value) {
      closeSuggestions()
    }
  })
}

const closeSuggestions = () => {
  showSuggestions.value = false
  isMouseOverSuggestions.value = false
}

// --- 监听器 ---

watch(
  () => props.posts,
  () => {
    updateKeywordCache()
    // 如果当前有内容，重新计算；如果是空的且聚焦，也会在上面 updateKeywordCache 里处理
    if (localSearch.value.trim()) {
      computeSuggestions(localSearch.value)
    }
  },
  { immediate: true },
)

watch(
  () => props.modelValue,
  (newVal) => {
    localSearch.value = newVal || ''
    if (!newVal && showSuggestions.value) {
      // 如果外部清空了搜索，且面板开着，显示推荐
      computeSuggestions('')
    } else if (!newVal) {
      // 如果没聚焦，则清空
      if (!showSuggestions.value) suggestions.value = []
    }
  },
)

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.search-container {
  position: relative;
  margin: 16px 0;
  width: 100%;
}

.search-wrapper {
  position: relative;
  width: 100%;
}

.custom-search-input :deep(.n-input__wrapper) {
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  background-color: #f9fafb;
}

.custom-search-input :deep(.n-input__wrapper:hover) {
  background-color: #f3f4f6;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.custom-search-input :deep(.n-input__wrapper.n-input__wrapper--focused) {
  background-color: #fff;
  box-shadow: 0 4px 16px rgba(59, 130, 246, 0.25);
  border-color: #3b82f6;
}

.search-icon {
  font-size: 14px;
  opacity: 0.6;
}

/* 建议列表通用样式 */
.search-suggestions {
  background: white;
  border-radius: 12px;
  box-shadow:
    0 10px 25px -5px rgba(0, 0, 0, 0.1),
    0 8px 10px -6px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  overflow: hidden;
  border: 1px solid #e5e7eb;
}

/* 桌面端：绝对定位 */
@media (min-width: 769px) {
  .search-suggestions {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 8px;
    max-height: 240px;
    overflow-y: auto;
  }
}

/* 移动端：流式布局，全宽显示 */
@media (max-width: 768px) {
  .search-suggestions.mobile-mode {
    position: absolute;
    margin-top: 8px;
    width: 100%;
    max-height: none;
    border: 1px solid #e5e7eb;
  }
}

.suggestions-header {
  display: none; /* 默认隐藏 */
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background-color: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
  font-weight: 600;
  font-size: 14px;
  color: #374151;
}

/* 强制显示头部的条件：移动端 或 (桌面端且有内容/推荐时想显示也可以打开这个) */
/* 这里为了保持桌面端简洁，默认桌面端不显示 header，除非你希望一直显示 */
@media (max-width: 768px) {
  .suggestions-header {
    display: flex;
  }
}
/* 如果你希望桌面端也显示这个标题栏，可以把上面的 display:none 去掉，或者添加如下逻辑 */
/* .suggestions-header.visible { display: flex; } */

.header-title {
  display: flex;
  align-items: center;
  gap: 6px;
}

.suggestion-item {
  display: flex;
  align-items: center;
  padding: 10px 14px;
  cursor: pointer;
  font-size: 14px;
  color: #374151;
  transition: background-color 0.2s ease;
  border-bottom: 1px solid #f3f4f6;
}

.suggestion-item:last-child {
  border-bottom: none;
}

.suggestion-item:hover,
.suggestion-item:focus {
  background-color: #eff6ff;
}

.item-icon {
  margin-right: 10px;
  font-size: 16px;
  opacity: 0.7;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
}

.emoji-icon {
  font-style: normal;
}

.item-text {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-text :deep(mark) {
  background-color: #fef3c7;
  color: #92400e;
  padding: 1px 2px;
  border-radius: 2px;
  font-weight: 600;
}

.suggestions-footer {
  padding: 8px 12px;
  text-align: center;
  background-color: #f9fafb;
  border-top: 1px solid #e5e7eb;
}

/* 过渡动画 */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.25s ease-out;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 移动端特定优化 */
@media (max-width: 768px) {
  .custom-search-input :deep(.n-input__wrapper) {
    height: 44px; /* 更大的点击区域 */
  }

  .custom-search-input :deep(.n-input__input-el) {
    font-size: 16px; /* 防止 iOS 自动缩放 */
  }

  .suggestion-item {
    padding: 12px 16px; /* 更大的触摸区域 */
    font-size: 15px;
  }
}
</style>
