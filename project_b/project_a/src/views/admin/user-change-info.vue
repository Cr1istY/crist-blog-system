<template>
  <div class="user-profile-container">
    <n-card title="修改个人信息" :bordered="false" style="max-width: 600px; margin: 0 auto">
      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="top"
        require-mark-placement="right-hanging"
      >
        <!-- 头像 URL -->
        <n-form-item label="头像链接" path="avatar">
          <n-input
            v-model:value="formData.avatar"
            type="text"
            placeholder="请输入头像图片的 URL 地址"
            clearable
          />
          <div v-if="formData.avatar" class="avatar-preview">
            <n-avatar :src="formData.avatar" size="large" />
          </div>
        </n-form-item>

        <!-- 昵称 -->
        <n-form-item label="昵称" path="nickname">
          <n-input
            v-model:value="formData.nickname"
            type="text"
            placeholder="请输入您的昵称"
            clearable
          />
        </n-form-item>

        <!-- 邮箱 (使用 text 类型 + 自定义正则验证) -->
        <n-form-item label="电子邮箱" path="email">
          <n-input
            v-model:value="formData.email"
            type="text"
            placeholder="example@email.com"
            clearable
          />
        </n-form-item>

        <!-- 个人简介 -->
        <n-form-item label="个人简介" path="bio">
          <n-input
            v-model:value="formData.bio"
            type="textarea"
            placeholder="介绍一下你自己..."
            :rows="4"
            maxlength="200"
            show-count
          />
        </n-form-item>

        <!-- 提交按钮 -->
        <n-form-item>
          <n-space justify="end">
            <n-button @click="handleReset">重置</n-button>
            <n-button type="primary" :loading="loading" @click="handleSubmit"> 保存修改 </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, NAvatar } from 'naive-ui'
// 修复：使用 import type 解决 verbatimModuleSyntax 错误
import type { FormRules, FormInst, FormValidationError } from 'naive-ui'
import type { User } from '@/types/tweet'
import service from '@/utils/request'

// --- 类型定义 ---

interface UserProfile {
  nickname: string
  email: string
  avatar: string
  bio: string
}

// --- 状态管理 ---

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref<boolean>(false)

// 表单数据
const formData = reactive<UserProfile>({
  nickname: '',
  email: '',
  avatar: '',
  bio: '',
})

// 表单验证规则
const rules: FormRules = {
  nickname: {
    required: true,
    message: '请输入昵称',
    trigger: ['blur', 'input'],
  },
  email: {
    required: true,
    message: '请输入电子邮箱',
    trigger: ['blur', 'input'],
    validator: (_rule: unknown, value: string) => {
      if (!value) {
        return new Error('请输入电子邮箱')
      }
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(value)) {
        return new Error('请输入有效的电子邮箱地址')
      }
      return true
    },
  },
  avatar: {
    validator: (_rule: unknown, value: string) => {
      if (value && !value.startsWith('http')) {
        return new Error('头像必须是有效的 HTTP/HTTPS 链接')
      }
      return true
    },
    trigger: 'blur',
  },
}

// --- 方法 ---

/**
 * 获取当前用户信息
 */
const fetchUserInfo = async () => {
  try {
    const token = localStorage.getItem('access_token')
    if (!token) {
      return null
    }
    // 假设 service 是 axios 实例，且拦截器已处理 token
    const res = await service.get('/tweet/getCurrentUser')
    // 根据实际后端返回结构调整，假设 data 字段包含用户信息
    const userData = res.data?.data || res.data
    if (userData) {
      return userData as User
    }
    return null
  } catch (error) {
    console.log('获取用户信息失败:', error)
    return null
  }
}

/**
 * 提交表单
 */
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (errors: FormValidationError[] | undefined) => {
    if (errors) {
      message.error('请检查表单填写是否正确')
      return
    }

    loading.value = true
    try {
      const payload = {
        nickname: formData.nickname,
        email: formData.email,
        avatar: formData.avatar,
        bio: formData.bio,
      }

      const res = await service.post('/changeUserInfo', payload)
      if (res.data.code !== 200) {
        message.error(res.data.message)
        return
      }
      message.success('用户信息更新成功')
    } catch (error: unknown) {
      let errMsg = '更新失败，请稍后重试'

      if (typeof error === 'object' && error !== null && 'message' in error) {
        errMsg = (error as { message: string }).message
      } else if (typeof error === 'string') {
        errMsg = error
      }

      console.error(error)
      message.error(errMsg)
    } finally {
      loading.value = false
    }
  })
}

/**
 * 重置表单
 */
const handleReset = () => {
  if (!formRef.value) return
  formRef.value.restoreValidation()
  fetchUserInfo()
}

// 生命周期
onMounted(async() => {
  const user = await fetchUserInfo()
  if (user !== null) {
    formData.nickname = user.displayName
    formData.email = user.email
    formData.avatar = user.avatar
    formData.bio = user.bio
  }
})
</script>

<style scoped>
.user-profile-container {
  padding: 24px;
  background-color: #f5f7f9;
  min-height: 100vh;
}

.avatar-preview {
  margin-top: 12px;
}
</style>
