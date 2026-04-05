<template>
  <n-layout>
    <n-layout-header style="padding: 20px; background-color: #f8f8f9">
      <div style="font-size: 24px; font-weight: bold">管理员面板</div>
    </n-layout-header>
    <n-layout-content content-style="padding: 20px;">
      <n-space vertical>
        <n-card title="文章管理" :bordered="false">
          <n-button type="primary" tag="a" :href="createUrl" target="_self">创建新文章</n-button>
          <!-- TODO 管理所有文章页面：可管理已删除文章，或者草稿文章等 -->
          <n-button type="primary" tag="a" :href="createUrl" target="_self">管理所有文章</n-button>
        </n-card>
        <n-card title="分类管理" :bordered="false">
          <n-button type="warning" tag="a" :href="categoryUrl" target="_self">创建分类</n-button>
          <!-- TODO 管理分类页面：可删除分类 -->
          <n-button type="warning" tag="a" :href="categoryUrl" target="_self">管理分类</n-button>
        </n-card>
        <n-card title="快速跳转" :bordered="false">
          <n-button type="primary" @click="$router.push('/')">回到首页</n-button>
          <n-button type="primary" @click="$router.push('/blog')">回到博客页面</n-button>
          <n-button type="primary" @click="$router.push('/admin/changeUserInfo')">更改个人信息</n-button>
        </n-card>
      </n-space>
    </n-layout-content>
  </n-layout>
</template>

<script lang="ts" setup>
defineOptions({
  name: 'AdminIndexPage'
})
import { computed } from 'vue'
import { NLayout, NLayoutHeader, NLayoutContent, NSpace, NCard, NButton } from 'naive-ui'

// 计算属性用于动态生成链接
const createUrl = computed(() => `/admin/post/create`)
const categoryUrl = computed(() => `/admin/category/create`)
</script>

<style scoped>
.n-button {
  margin-right: 2vw;
}
</style>
