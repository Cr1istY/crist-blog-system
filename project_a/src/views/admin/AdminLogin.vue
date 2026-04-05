<template>
  <div class="login-container">
    <!-- 可选：添加一个极简 Logo 或标题图标 -->
    <div class="logo-placeholder">
      <img src="/favicon.ico" alt="Logo" class="logo" />
    </div>

    <n-card class="login-card" :bordered="false" size="large">
      <template #header>
        <h2 class="login-title">欢迎回来</h2>
        <p class="login-subtitle">请登录以继续</p>
      </template>

      <n-form
        ref="formRef"
        :model="formValue"
        :rules="rules"
        label-placement="left"
        class="login-form"
      >
        <n-form-item path="username" label="用户名">
          <n-input
            v-model:value="formValue.username"
            placeholder="输入用户名或邮箱"
            clearable
            size="large"
          />
        </n-form-item>

        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="formValue.password"
            type="password"
            show-password-on="click"
            placeholder="至少6位字符"
            clearable
            size="large"
          />
        </n-form-item>

        <n-space vertical :size="16" style="width: 100%">
          <n-button
            type="primary"
            block
            size="large"
            :loading="loading"
            @click="handleSubmit"
            class="login-button"
          >
            登录
          </n-button>
          <n-flex>
            <n-button text @click="handleForgotPassword" class="forgot-link"> 忘记密码 </n-button>
            <n-button text @click="handleLogout" class="admin-link"> 不登录 </n-button>
            <n-button text @click="handleRegister" class="register-link"> 注册 </n-button>
          </n-flex>
        </n-space>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
// （你的 script 部分完全保留，无需修改）
import { ref } from 'vue'
import type { FormInst, FormRules } from 'naive-ui'
import { useMessage } from 'naive-ui'
import request from '@/utils/request'
import type { AxiosRetryConfig } from '@/utils/request'
import router from '@/router/index'
import type { AxiosError } from 'axios'

interface LoginForm {
  username: string
  password: string
}

const formRef = ref<FormInst | null>(null)
const loading = ref<boolean>(false)
const message = useMessage()

const formValue = ref<LoginForm>({
  username: '',
  password: '',
})

interface LoginResponse {
  access_token: string
}

const rules: FormRules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur',
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: 'blur',
    validator(_, value) {
      if (value && value.length < 6) {
        return new Error('密码至少6位')
      }
      return true
    },
  },
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    // 表单验证
    await formRef.value.validate()

    try {
      // API 请求
      const res = await request.post<LoginResponse>('/login', formValue.value, {
        customErrorHandling: true,
      } as AxiosRetryConfig)
      const token = res.data.access_token

      if (token) {
        localStorage.setItem('access_token', token)
        message.success('登录成功')
        router.push('/admin/dashboard')
      } else {
        message.error('登录失败，请检查用户名和密码')
      }
    } catch (apiError: unknown) {
      if (apiError && typeof apiError === 'object' && 'response' in apiError) {
        // 处理 API 请求错误
        const error = apiError as AxiosError
        if (error.response) {
          // 服务器返回了错误状态码
          const status = error.response.status
          if (status === 401) {
            message.error('用户名或密码错误')
          } else {
            message.error(`服务器错误: ${status}`)
          }
        } else if (error.request) {
          // 请求已发出但没有收到响应
          message.error('网络错误，请检查网络连接')
        } else {
          // 其他错误
          message.error('登录失败，请稍后重试')
        }
      }
    }
  } catch (validationError: Error | unknown) {
    // 处理表单验证错误
    console.error('表单验证失败:', validationError)
    message.error('请检查输入内容')
  }
}

const handleForgotPassword = () => {
  message.info('请联系管理员重置密码')
}

const handleLogout = () => {
  localStorage.removeItem('access_token')
  message.success('退出登录成功')
  router.back()
}

const handleRegister = () => {
  message.info('请联系管理员注册')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4edf9 100%);
  /* 未来感柔光背景 */
}

.logo-placeholder {
  position: absolute;
  top: 15vh;
  font-size: 28px;
  font-weight: bold;
  color: var(--n-primary-color);
  opacity: 0.8;
}

.login-card {
  width: 100%;
  max-width: 420px;
  border-radius: 16px;
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.5);
  padding: 32px;
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
}

.login-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px -12px rgba(0, 0, 0, 0.12);
}

.login-title {
  font-size: 24px;
  font-weight: 600;
  text-align: center;
  margin: 0 0 8px;
  color: var(--n-text-color);
}

.login-subtitle {
  font-size: 14px;
  color: var(--n-text-color-secondary);
  text-align: center;
  margin: 0 0 24px;
}

.login-form :deep(.n-form-item) {
  margin-bottom: 20px;
}

.login-form :deep(.n-form-item-label) {
  font-weight: 500;
  color: var(--n-text-color);
}

.login-form :deep(.n-input) {
  transition: all 0.25s ease;
}

.login-form :deep(.n-input:hover),
.login-form :deep(.n-input:focus) {
  border-color: var(--n-primary-color);
  box-shadow: 0 0 0 2px rgba(var(--n-primary-color-rgb), 0.15);
}

.login-button {
  height: 48px;
  font-weight: 500;
  border-radius: 12px;
}

.forgot-link {
  align-self: center;
  font-size: 14px;
  color: var(--n-text-color-secondary);
  transition: color 0.2s ease;
}

.forgot-link:hover {
  color: var(--n-primary-color);
}

/* 响应式优化 */
@media (max-width: 480px) {
  .login-card {
    padding: 24px;
  }

  .login-title {
    font-size: 22px;
  }
}
</style>
