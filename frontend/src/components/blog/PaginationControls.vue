<template>
  <div class="pagination-wrapper">
    <n-pagination
      v-model:page="localPage"
      v-model:page-size="localPageSize"
      :page-count="totalPages"
      show-size-picker
      :page-sizes="[12, 20]"
      class="pagination"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  pageSize: number
  totalPages: number
}>()

const emit = defineEmits<{
  'update:page': [page: number]
  'update:page-size': [size: number]
}>()

const localPage = computed({
  get: () => props.page,
  set: (val) => emit('update:page', val)
})

const localPageSize = computed({
  get: () => props.pageSize,
  set: (val) => emit('update:page-size', val)
})
</script>

<style scoped>
.pagination-wrapper { display: flex; justify-content: center; margin-top: 32px; min-height: 40px }
:deep(.n-pagination) { height: 40px }
</style>
