<template>
  <div class="layout-outer">
    <div class="page-container">
      <div>
        <!-- 头部简介区 -->
        <n-card embedded :bordered="false" class="hero-section">
          <n-space vertical align="center" size="large">
            <n-avatar
              round
              lazy
              :size="isMobile ? 80 : 100"
              src="https://avatars.githubusercontent.com/u/149236474?v=4"
            />
            <div class="intro-text">
              <n-h2 style="margin-bottom: 2px">Crist Yang</n-h2>
              <n-text type="primary">全栈学习中 · 在校学生</n-text>
              <n-p class="bio">后端开发，物联网工程，偶尔前端，顺手设计，持续写作</n-p>
            </div>
          </n-space>
        </n-card>
        <div class="content-wrapper">
          <!-- 主内容区 -->
          <n-grid :cols="isMobile ? 1 : 2" :x-gap="24" :y-gap="24" class="main-grid">
            <!-- 左侧：最新文章 + 精选文章 -->
            <n-grid-item>
              <!-- 精选文章 -->
              <n-card :bordered="false" size="medium">
                <template #header>
                  <n-button text tag="a" @click="$router.push('/blog')" class="featured-title">
                    🔥 精选文章
                  </n-button>
                </template>
                <n-list>
                  <n-list-item v-for="post in hotPosts" :key="post.slug">
                    <n-thing>
                      <template #header>
                        <router-link :to="`/blog/${post.slug}`" class="post-title">
                          {{ post.title }}
                        </router-link>
                      </template>
                      <template #description>
                        <n-space size="small">
                          <n-tag type="info" size="small">{{ post.category }}</n-tag>
                          <n-text type="secondary">{{ post.date }}</n-text>
                        </n-space>
                      </template>
                      <template #footer>
                        {{ post.excerpt }}
                      </template>
                    </n-thing>
                  </n-list-item>
                </n-list>
              </n-card>

              <!-- 最新文章 -->
              <n-card title="📝 最新文章" :bordered="false" size="medium" style="margin-top: 0px">
                <n-list>
                  <n-list-item v-for="post in latestPosts" :key="post.slug">
                    <router-link :to="`/blog/${post.slug}`" custom v-slot="{ navigate }">
                      <n-thing :title="post.title" @click="navigate" style="cursor: pointer">
                        <template #description>
                          <n-space size="small">
                            <n-tag type="info" size="small">{{ post.category }}</n-tag>
                            <n-text type="secondary">{{ post.date }}</n-text>
                          </n-space>
                        </template>
                      </n-thing>
                    </router-link>
                  </n-list-item>
                </n-list>
              </n-card>
            </n-grid-item>

            <!-- 右侧（仅桌面）：关于 + 技术栈 + 社交 -->
            <n-grid-item class="right-content">
              <!-- 关于我 -->
              <n-card :bordered="false" size="medium">
                <template #header>
                  <n-button text tag="a" @click="$router.push('/thread')" class="featured-title">
                    👨‍💻 About Me
                  </n-button>
                </template>
                <n-p>- 👨‍🎓 IoT Engineering student at CQUPT.</n-p>
                <n-p>- 🌍 Who loves open-source and make an impact.</n-p>
                <n-p>- 🛠️ Building AIoT systems with Java, Go, and Python.</n-p>
                <n-p>- 📚 Book lover & music fan - Delighting You Always.</n-p>
              </n-card>

              <!-- 技术栈 -->
              <n-card title="🛠 技术栈" :bordered="false" size="medium" style="margin-top: 0px">
                <n-space wrap :size="[12, 12]">
                  <n-tag v-for="tech in techStack" :key="tech" type="info" size="small" round>
                    <template #icon>
                      <n-icon :component="getIcon(tech)" />
                    </template>
                    {{ tech }}
                  </n-tag>
                </n-space>
                <n-space wrap :size="[12, 12]" style="margin-top: 8px">
                  <n-tag
                    v-for="tech in techStackLanguage"
                    :key="tech"
                    type="warning"
                    size="small"
                    round
                  >
                    <template #icon>
                      <n-icon :component="getIcon(tech)" />
                    </template>
                    {{ tech }}
                  </n-tag>
                </n-space>
                <n-space wrap :size="[12, 12]" style="margin-top: 8px">
                  <n-tag v-for="tech in techStackIoT" :key="tech" type="success" size="small" round>
                    <template #icon>
                      <n-icon :component="getIcon(tech)" />
                    </template>
                    {{ tech }}
                  </n-tag>
                </n-space>
              </n-card>

              <!-- 联系方式 -->
              <n-card
                title="📬 联系我"
                :bordered="false"
                size="medium"
                style="margin-top: 0px; margin-bottom: 24px"
              >
                <n-space vertical>
                  <n-button text tag="a" href="mailto:cr1st4ever@outlook.com">
                    📧 cr1st4ever@outlook.com
                  </n-button>
                  <n-button text tag="a" href="https://github.com/Cr1istY" target="_blank">
                    🐙 GitHub (@Cr1istY)
                  </n-button>
                  <n-button
                    text
                    tag="a"
                    href="https://steamcommunity.com/profiles/76561199027120002"
                    target="_blank"
                  >
                    🎮 Steam (@CristY)
                  </n-button>
                </n-space>
              </n-card>
            </n-grid-item>
          </n-grid>
        </div>
      </div>

      <footer class="site-footer">
        <div class="footer-content">
          <div class="copyright">
            <span class="icp">
              <n-button text tag="a" href="https://beian.miit.gov.cn/" target="_blank">
                渝ICP备2025056615号
              </n-button>
            </span>
            <span>© {{ new Date().getFullYear() }} Crist Yang.</span>
          </div>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'HomePage'
})
import { LogoVue, Code, HardwareChip, LogoDocker } from '@vicons/ionicons5'

import { LogoPython } from '@vicons/carbon'
import { BrandJavascript } from '@vicons/tabler'
import { Java } from '@vicons/fa'
import type { Component } from 'vue'
import { defineComponent } from 'vue'
import { Icon } from '@iconify/vue'

import { ref, onMounted } from 'vue'

const fetchPosts = async <T,>(endpoint: string): Promise<T[]> => {
  const response = await fetch(endpoint)
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  return response.json() as Promise<T[]>
}

const isMobile = ref(window.innerWidth < 768)
window.addEventListener('resize', () => {
  isMobile.value = window.innerWidth < 768
})

interface HotPost {
  slug: string
  title: string
  category: string
  date: string
  excerpt: string
}

const hotPosts = ref<HotPost[]>([])
const latestPosts = ref<HotPost[]>([])

// 模拟数据
const loadPosts = async () => {
  try {
    const [hot, latest] = await Promise.all([
      fetchPosts<HotPost>('/api/posts/hot'),
      fetchPosts<HotPost>('/api/posts/latest'),
    ])

    hotPosts.value = hot
    latestPosts.value = latest
  } catch (error) {
    console.error('Failed to load posts:', error)
    // 可以在这里添加错误处理逻辑，比如显示错误消息
  }
}

const GoIcon = defineComponent({
  render() {
    return h(Icon, { icon: 'simple-icons:jaeger' })
  },
})

const TsIcon = defineComponent({
  render() {
    return h(Icon, { icon: 'simple-icons:typescript' })
  },
})

const CppIcon = defineComponent({
  render() {
    return h(Icon, { icon: 'simple-icons:cplusplus' })
  },
})

const EspIcon = defineComponent({
  render() {
    return h(Icon, { icon: 'simple-icons:espressif' })
  },
})

const MqttIcon = defineComponent({
  render() {
    return h(Icon, { icon: 'simple-icons:mqtt' })
  },
})

const StmIcon = defineComponent({
  render() {
    return h(Icon, { icon: 'simple-icons:stmicroelectronics' })
  },
})

const techStack: TechType[] = ['Vue 3', 'TypeScript', 'JavaScript', 'Docker']
const techStackLanguage: TechType[] = ['Java', 'Go', 'Python', 'C/C++']
const techStackIoT: TechType[] = ['Arduino', 'ESP32', 'Mqtt', 'STM32']

type TechType =
  | 'Vue 3'
  | 'TypeScript'
  | 'JavaScript'
  | 'Docker'
  | 'Java'
  | 'Go'
  | 'Python'
  | 'C/C++'
  | 'Arduino'
  | 'ESP32'
  | 'Mqtt'
  | 'STM32'

const iconMap: Record<TechType, Component> = {
  // Web技术栈
  'Vue 3': LogoVue,
  TypeScript: TsIcon,
  JavaScript: BrandJavascript,
  Docker: LogoDocker,

  // 编程语言
  Java: Java,
  Go: GoIcon,
  Python: LogoPython,
  'C/C++': CppIcon,

  // IoT相关
  Arduino: HardwareChip,
  ESP32: EspIcon,
  Mqtt: MqttIcon,
  STM32: StmIcon,
}

const getIcon = (tech: TechType) => {
  return iconMap[tech] ?? Code
}

onMounted(() => {
  loadPosts()
})
</script>

<style scoped>
.layout-outer {
  overflow-x: hidden; /* 阻止任何子元素溢出产生滚动条 */
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 0 16px;
  box-sizing: border-box;
}

.content-wrapper {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  box-sizing: border-box;
}

.hero-section {
  padding: 32px 24px;
  text-align: center;
  width: 100vw;
  margin-left: calc(-50vw + 50%);
}

.main-grid {
  margin-top: 24px;
  margin-left: 0;
}

/* 深度选择器修改 n-card 的样式 */
.hero-section > :deep(.n-card) {
  border-radius: 0;
  max-width: 1200px;
  margin: 0 auto;
  background: transparent;
  box-shadow: none;
}

.intro-text {
  max-width: 600px;
  margin: 0 auto;
}

.bio {
  margin-top: 8px;
  color: var(--n-text-color);
  line-height: 1.6;
}

.right-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.right-content :deep(.n-card) {
  flex: 1;
  margin-top: 24px;
}

.right-content :deep(.n-card:first-child) {
  margin-top: 0px;
}

.featured-title {
  color: inherit;
  text-decoration: none;
  font-size: 18px;
  font-weight: 500;
  cursor: pointer;
}

.featured-title:hover {
  color: #f0a020;
}

.post-title {
  color: inherit;
  text-decoration: none;
  cursor: pointer;
}

.post-title:hover {
  color: #18a058;
  text-decoration: underline;
}

.site-footer {
  background-color: #fff;
  padding: 20px 0;
  margin-top: auto;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.copyright {
  display: flex;
  flex-direction: row;
  gap: 12px;
}

.icp {
  display: flex;
  flex-direction: row;
  color: #666;
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 0;
    margin: 0;
  }

  .main-grid {
    margin: 0;
  }

  .footer-content {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }

  .copyright {
    align-items: center;
  }

  .hero-section {
    padding: 32px 16px;
  }

  /* 移动端禁用缩放 */
  .page-container {
    transform: none;
    width: 100%;
  }
}
</style>
