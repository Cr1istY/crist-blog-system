<script setup lang="ts">
defineOptions({
  name: 'TweetComposer',
})

import { ref, computed, watch, onBeforeUnmount } from 'vue'
import {
  ImageOutline,
  CloseOutline,
  WarningOutline,
  CheckmarkCircleOutline,
} from '@vicons/ionicons5'
import type { TweetComposerProps, UploadedImage, TweetImage } from '@/types/tweet'
import { uploadImage } from '@/services/upload'
import service from '@/utils/request'

// 定义 Props
const props = withDefaults(defineProps<TweetComposerProps>(), {
  placeholder: '有什么新鲜事？',
  submitLabel: '发布',
  maxLength: 280,
  maxImages: 9,
  allowedImageTypes: () => ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  maxImageSize: 5,
  currentUser: undefined,
})

// 定义 Emits
const emit = defineEmits<{
  submit: [content: string, images: TweetImage[]]
}>()

// 状态管理
const content = ref<string>('')
const isSubmitting = ref<boolean>(false)
const images = ref<UploadedImage[]>([])
const fileInputRef = ref<HTMLInputElement | null>(null)

// 计算属性
const charCount = computed<number>(() => content.value.length)
const isOverLimit = computed<boolean>(() => charCount.value > props.maxLength)

const canSubmit = computed<boolean>(
  () =>
    content.value.trim().length > 0 &&
    charCount.value <= props.maxLength &&
    !isSubmitting.value &&
    images.value.every((img) => !img.uploading && !img.error),
)

const canAddMoreImages = computed<boolean>(() => images.value.length < props.maxImages)
const remainingImageCount = computed<number>(() => props.maxImages - images.value.length)

// 获取当前图片数量对应的布局类名后缀 (用于匹配展示页的 CSS)
const imageLayoutClass = computed<string>(() => {
  const count = images.value.length
  if (count === 0) return ''
  return `image-count-${Math.min(count, 9)}`
})

// 判断是否为单图模式 (用于特殊处理 object-fit)
const isSingleImageMode = computed<boolean>(() => images.value.length === 1)

// 工具函数
const generateId = (): string => `img-${Date.now()}-${Math.random().toString(36).slice(2)}`

const validateImage = (file: File): string | null => {
  if (!props.allowedImageTypes.includes(file.type)) {
    return '不支持的图片格式'
  }
  const maxSizeBytes = props.maxImageSize * 1024 * 1024
  if (file.size > maxSizeBytes) {
    return `图片大小不能超过 ${props.maxImageSize}MB`
  }
  return null
}

const createPreviewUrl = (file: File): string => {
  return URL.createObjectURL(file)
}

// 文件选择处理
const handleFileSelect = (event: Event): void => {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (!files || files.length === 0) return

  const remainingSlots = props.maxImages - images.value.length
  const filesToProcess = Array.from(files).slice(0, remainingSlots)

  filesToProcess.forEach((file: File) => {
    const error = validateImage(file)
    const uploadedImage: UploadedImage = {
      id: generateId(),
      file,
      previewUrl: createPreviewUrl(file),
      uploading: false,
      progress: 0,
      error: error || undefined,
    }
    images.value.push(uploadedImage)
  })

  // 重置 input 以便重复选择相同文件
  target.value = ''
}

const triggerFileSelect = (): void => {
  fileInputRef.value?.click()
}

const removeImage = (imageId: string): void => {
  const image = images.value.find((img) => img.id === imageId)
  if (image) {
    URL.revokeObjectURL(image.previewUrl)
  }
  images.value = images.value.filter((img) => img.id !== imageId)
}

// 上传逻辑
const uploadImageToServer = async (image: UploadedImage): Promise<void> => {
  image.uploading = true
  image.progress = 0
  try {
    const response = await uploadImage(image.file, {
      onProgress: (progress) => {
        image.progress = progress.percentage
      },
    })
    image.uploading = false
    image.progress = 100
    image.id = response.id
    image.uploadedUrl = response.url
  } catch (error) {
    image.uploading = false
    image.error = '上传失败，请重试'
    console.error('上传错误:', error)
  }
}

const uploadAllImages = async (): Promise<boolean> => {
  const pendingImages = images.value.filter((img) => !img.uploading && !img.error)
  if (pendingImages.length === 0) return true

  const uploadPromises = pendingImages.map((img) => uploadImageToServer(img))
  await Promise.all(uploadPromises)

  return images.value.every((img) => !img.error)
}

// 提交逻辑
const handleSubmit = async (): Promise<void> => {
  if (!canSubmit.value) return

  isSubmitting.value = true
  try {
    const uploadSuccess = await uploadAllImages()
    if (!uploadSuccess) return

    const tweetImages: TweetImage[] = images.value.map((img) => ({
      id: img.id,
      url: img.uploadedUrl || img.previewUrl,
    }))

    const imageIds = tweetImages.map((img) => img.id)
    const tweet = {
      content: content.value.trim(),
      image_ids: imageIds,
    }

    // 调用接口
    console.log('提交内容:', tweet)
    const response = await service.post('/tweet/create', tweet)
    console.log('提交成功:', response)

    emit('submit', content.value.trim(), tweetImages)

    // 清理资源
    images.value.forEach((img) => {
      URL.revokeObjectURL(img.previewUrl)
    })
    content.value = ''
    images.value = []
  } catch (error) {
    console.error('提交失败:', error)
  } finally {
    isSubmitting.value = false
  }
}

const handleKeydown = (event: KeyboardEvent): void => {
  if (event.ctrlKey && event.key === 'Enter') {
    event.preventDefault()
    handleSubmit()
  }
}

watch(content, (newValue: string): void => {
  if (newValue.length > props.maxLength + 10) {
    content.value = newValue.slice(0, props.maxLength)
  }
})

onBeforeUnmount((): void => {
  images.value.forEach((img) => {
    URL.revokeObjectURL(img.previewUrl)
  })
})
</script>

<template>
  <div class="tweet-composer">
    <!-- 隐藏的文件输入框 -->
    <input
      ref="fileInputRef"
      type="file"
      accept="image/*"
      multiple
      class="hidden-input"
      @change="handleFileSelect"
    />

    <div class="composer-header">
      <NAvatar
        :src="currentUser?.avatar"
        size="large"
        round
        aria-label="用户头像"
        :fallback-src="'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg'"
      />

      <div class="composer-input-container">
        <!-- 文本输入 -->
        <NInput
          v-model:value="content"
          type="textarea"
          :placeholder="placeholder"
          :class="['composer-input', { 'over-limit': isOverLimit }]"
          :maxlength="maxLength + 10"
          :rows="3"
          :autosize="{ minRows: 3, maxRows: 6 }"
          @keydown="handleKeydown"
          show-count
        />

        <!-- 图片预览区域 (核心修改部分) -->
        <div v-if="images.length > 0" class="image-preview-wrapper">
          <NImageGroup>
            <!--
              动态绑定类名以匹配展示页的 CSS (.image-count-1, .image-count-2 等)
              这确保了所见即所得
            -->
            <div class="tweet-images" :class="imageLayoutClass">
              <div v-for="image in images" :key="image.id" class="image-item-container">
                <div class="image-inner-wrapper">
                  <NImage
                    :src="image.previewUrl"
                    alt="预览图片"
                    class="tweet-image"
                    :style="{
                      objectFit: isSingleImageMode ? 'contain' : 'cover',
                      objectPosition: 'center',
                    }"
                    :img-props="{ class: 'custom-img-inner' }"
                  >
                    <template #placeholder>
                      <div class="image-loading-placeholder">
                        <NProgress
                          type="circle"
                          :percentage="image.progress"
                          :size="40"
                          :show-indicator="false"
                        />
                      </div>
                    </template>
                  </NImage>

                  <!-- 上传进度覆盖层 -->
                  <div v-if="image.uploading" class="upload-progress-overlay">
                    <NProgress
                      :percentage="image.progress"
                      type="line"
                      :height="4"
                      :show-indicator="true"
                      indicator-placement="inside"
                      class="image-progress-bar"
                    />
                  </div>

                  <!-- 错误覆盖层 -->
                  <div v-if="image.error" class="image-error-overlay">
                    <NIcon :component="WarningOutline" size="24" color="#fff" />
                    <span class="error-text">{{ image.error }}</span>
                  </div>

                  <!-- 删除按钮 -->
                  <NButton
                    type="default"
                    circle
                    size="small"
                    class="remove-image-btn"
                    @click="removeImage(image.id)"
                    :disabled="image.uploading"
                    aria-label="删除图片"
                  >
                    <template #icon>
                      <NIcon :component="CloseOutline" />
                    </template>
                  </NButton>

                  <!-- 上传成功标识 -->
                  <div
                    v-if="!image.uploading && !image.error && image.uploadedUrl"
                    class="upload-success-badge"
                  >
                    <NIcon :component="CheckmarkCircleOutline" size="18" color="#18a058" />
                  </div>
                </div>
              </div>
            </div>
          </NImageGroup>
        </div>

        <!-- 底部操作栏 -->
        <div class="composer-footer">
          <div class="composer-actions">
            <NButton
              type="default"
              :disabled="!canAddMoreImages"
              @click="triggerFileSelect"
              class="add-image-btn"
              dashed
            >
              <template #icon>
                <NIcon :component="ImageOutline" />
              </template>
              <span v-if="canAddMoreImages" class="add-image-text">
                图片 ({{ remainingImageCount }})
              </span>
            </NButton>
          </div>

          <div class="composer-submit">
            <NTag
              :type="charCount > maxLength * 0.9 ? (isOverLimit ? 'error' : 'warning') : 'default'"
              :bordered="false"
              class="char-count-tag"
              size="small"
            >
              {{ charCount }}/{{ maxLength }}
            </NTag>

            <NButton
              type="primary"
              :disabled="!canSubmit"
              :loading="isSubmitting"
              class="submit-button"
              @click="handleSubmit"
            >
              {{ submitLabel }}
            </NButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tweet-composer {
  padding: 1rem;
  border-bottom: 8px solid var(--n-color-hover);
  background-color: var(--n-color);
}

.composer-header {
  display: flex;
  gap: 12px;
}

.composer-input-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.composer-input {
  font-size: 15px;
  resize: none;
  border: none;
  box-shadow: none;
  padding-left: 0;
  padding-right: 0;
}

.composer-input:focus {
  box-shadow: none;
}

.composer-input.over-limit :deep(.n-input__textarea-el) {
  color: var(--n-error-text-color, #d03050);
}

/* 隐藏文件输入 */
.hidden-input {
  display: none;
}

/* --- 图片预览容器 (完全复刻展示页逻辑) --- */
.image-preview-wrapper {
  margin-top: 12px;
  width: 100%;
}

.tweet-images {
  display: grid;
  gap: 6px;
  width: 100%;
  box-sizing: border-box;
}

.image-item-container {
  /* 占位容器，确保 Grid 布局正常 */
  width: 100%;
  height: 100%;
  position: relative;
}

.image-inner-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 4px;
  overflow: hidden;
  background-color: #f0f0f0;
}

.tweet-image {
  width: 100%;
  height: 100%;
  display: block;
  cursor: zoom-in;
  border-radius: 4px;
}

:deep(.custom-img-inner) {
  width: 100%;
  height: 100%;
  /* object-fit 由父组件 NImage 的 style 动态控制 */
}

.image-loading-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: #f5f5f5;
}

/* --- 布局场景 1：单图 --- */
/* 不强制正方形，显示原图比例，限制最大尺寸 */
.image-count-1 {
  display: block;
  max-width: 70%;
  max-height: 400px;
}

.image-count-1 .image-item-container {
  /* 单图不需要 aspect-ratio，让内容撑开 */
  aspect-ratio: auto;
  height: auto;
  max-height: 400px;
}

/* --- 布局场景 2：双图 --- */
.image-count-2 {
  grid-template-columns: repeat(2, 1fr);
}

.image-count-2 .image-item-container {
  aspect-ratio: 1 / 1;
}

/* --- 布局场景 3-9：九宫格 --- */
.image-count-3,
.image-count-4,
.image-count-5,
.image-count-6,
.image-count-7,
.image-count-8,
.image-count-9 {
  grid-template-columns: repeat(3, 1fr);
}

.image-count-3 .image-item-container,
.image-count-4 .image-item-container,
.image-count-5 .image-item-container,
.image-count-6 .image-item-container,
.image-count-7 .image-item-container,
.image-count-8 .image-item-container,
.image-count-9 .image-item-container {
  aspect-ratio: 1 / 1;
}

/* --- 覆盖层样式 --- */
.upload-progress-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2;
  padding: 10px;
}

.image-progress-bar {
  width: 100%;
}

.image-error-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(208, 48, 80, 0.9);
  color: white;
  gap: 8px;
  z-index: 2;
}

.error-text {
  font-size: 12px;
  text-align: center;
  padding: 0 4px;
}

.remove-image-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  border: none;
  transition: all 0.2s;
}

.remove-image-btn:hover:not(:disabled) {
  background: rgba(208, 48, 80, 0.9);
  transform: scale(1.1);
}

.remove-image-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.upload-success-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 10;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  padding: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* --- 底部操作栏 --- */
.composer-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  padding-top: 12px;
}

.composer-actions {
  display: flex;
  gap: 8px;
}

.add-image-btn {
  border-radius: 9999px;
}

.add-image-text {
  margin-left: 4px;
  font-size: 13px;
}

.composer-submit {
  display: flex;
  align-items: center;
  gap: 12px;
}

.char-count-tag {
  font-family: monospace;
}

.submit-button {
  border-radius: 9999px;
  padding-left: 24px;
  padding-right: 24px;
  font-weight: 600;
}

/* 移动端适配 */
@media (max-width: 480px) {
  .image-count-1 {
    max-width: 85%;
  }

  .composer-footer {
    flex-wrap: wrap;
    gap: 12px;
  }

  .composer-actions {
    order: 2;
    width: 100%;
  }

  .composer-submit {
    order: 1;
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
