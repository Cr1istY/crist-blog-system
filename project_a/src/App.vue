<script setup lang="ts">

</script>

<template>
  <n-config-provider>
    <!-- 必须包裹 message provider -->
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <!-- 所有页面内容 -->
          <router-view />
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped></style>
