<template>
  <div class="placeholder-view">
    <el-empty description="销售管理 · 功能开发中">
      <template #image>
        <el-icon :size="80" color="#667eea"><TrendCharts /></el-icon>
      </template>
    </el-empty>
  </div>
</template>

<script setup lang="ts">
import { TrendCharts } from '@element-plus/icons-vue'
</script>

<style scoped>
.placeholder-view {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
}
</style>
