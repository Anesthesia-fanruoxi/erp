<template>
  <div class="dashboard-view">
    <h1>欢迎使用 ERP 系统</h1>
    <p class="welcome-text">您好，{{ userInfo?.realName || '用户' }}，祝您工作愉快！</p>

    <el-row :gutter="20" style="margin-top: 24px">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon :size="20" color="#409EFF"><User /></el-icon>
              <span>用户总数</span>
            </div>
          </template>
          <div class="card-value">--</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon :size="20" color="#E6A23C"><Clock /></el-icon>
              <span>待审核</span>
            </div>
          </template>
          <div class="card-value">--</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon :size="20" color="#67C23A"><Check /></el-icon>
              <span>正常用户</span>
            </div>
          </template>
          <div class="card-value">--</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { User, Clock, Check } from '@element-plus/icons-vue'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)
</script>

<style scoped>
.dashboard-view {
  padding: 20px;
}
.welcome-text {
  color: #606266;
  font-size: 16px;
  margin-top: 8px;
}
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
.card-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  text-align: center;
}
</style>
