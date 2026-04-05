<template>
  <div class="not-found-page">
    <n-result
      status="404"
      title="页面走丢了 🥺"
      description="抱歉，你访问的页面不存在。可能是被猫叼走了，或者它从未存在过。"
    >
      <template #footer>
        <n-button
          type="primary"
          @click="goHome"
          size="large"
          ghost
        >
          回到首页
        </n-button>
      </template>
    </n-result>
  </div>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goHome = () => {
  router.push('/');
};
</script>

<style scoped>
.not-found-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4edf9 100%);
}
</style>
