<template>
  <n-layout style="padding: 20px">
    <div class="form-container">
      <h2 class="form-title">创建新博客：</h2>

      <n-form ref="formRef" :model="post" label-placement="left" label-width="140px">
        <n-form-item path="thumbnail" label="封面">
          <n-input v-model:value="post.thumbnail" placeholder="Thumbnail URL" />
        </n-form-item>

        <n-form-item path="content" label="内容">
          <MdEditor v-model="post.content" preview-theme="vuepress" style="min-height: 800px;"/>
        </n-form-item>

        <n-form-item path="excerpt" label="摘要">
          <n-input
            type="textarea"
            v-model:value="post.excerpt"
            placeholder="Short summary for listing pages"
            :autosize="{ minRows: 2, maxRows: 4 }"
          />
        </n-form-item>

        <n-form-item path="status" label="状态">
          <n-select v-model:value="post.status" :options="statusOptions" />
        </n-form-item>

        <!-- 分类选择 -->
        <n-form-item path="category_id" label="分类">
          <n-select
            v-model:value="post.category_id"
            :options="categoryOptions"
            placeholder="Select a category"
            :loading="loadingCategories"
          />
        </n-form-item>

        <n-form-item path="tags" label="标签">
          <n-dynamic-tags v-model:value="post.tags" placeholder="Press Enter to add tags" />
        </n-form-item>

        <n-space justify="end" style="margin-top: 24px">
          <n-button type="primary" @click="submitForm">Submit Post</n-button>
        </n-space>
      </n-form>
    </div>
  </n-layout>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { MdEditor } from 'md-editor-v3'
import { useCategoryTree } from '@/composables/useCategoryTree'
import router from '@/router/index'
import 'md-editor-v3/lib/style.css'
import service from '@/utils/request'

const { categoryOptions, loading: loadingCategories, fetchCategories } = useCategoryTree()


// 表单数据
const post = reactive({
  content: '',
  excerpt: '',
  status: 'published',
  category_id: null as number | null,
  tags: [] as string[],
  thumbnail: '',
})

// 状态选项
const statusOptions = [
  { label: '公开', value: 'published' },
  { label: '归档', value: 'archived' },
  { label: '草稿', value: 'draft' },
]


// 提交表单
async function submitForm() {
  try {
    // 如果 category_id 是 null，可以传 null 或忽略（根据后端要求）
    const response = await service.post('/posts/create', post)
    if (response.status === 200) {
      router.push('/blog') // 重定向到文章列表页面
    } else {
      console.log(response.data)
    }
  } catch (error) {
    console.error('There was an error submitting your form!', error)
  }
}

// 组件挂载时加载分类
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.form-container {
  max-width: 80vw;
  margin: 0 auto;
}

.form-title {
  margin-bottom: 24px;
  font-size: 24px;
  font-weight: 600;
  color: var(--n-text-color);
}
</style>
