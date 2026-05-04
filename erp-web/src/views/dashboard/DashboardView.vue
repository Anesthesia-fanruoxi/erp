<template>
  <div class="dashboard-view">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-info">
        <h1 class="welcome-title">
          {{ greeting }}，{{ userInfo?.realName || '用户' }}
        </h1>
        <p class="welcome-desc">欢迎回到 ERP 管理系统，祝您工作顺利</p>
      </div>
      <div class="welcome-date">
        <div class="date-text">{{ currentDate }}</div>
        <div class="week-text">{{ currentWeek }}</div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="8">
        <div class="stat-card card-blue">
          <div class="stat-icon">
            <el-icon :size="28"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">--</div>
            <div class="stat-label">用户总数</div>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-card card-orange">
          <div class="stat-icon">
            <el-icon :size="28"><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">--</div>
            <div class="stat-label">待审核</div>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-card card-green">
          <div class="stat-icon">
            <el-icon :size="28"><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">--</div>
            <div class="stat-label">正常用户</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <div class="quick-section">
      <h3 class="section-title">快捷操作</h3>
      <el-row :gutter="16">
        <el-col :span="6">
          <div class="quick-card" @click="$router.push('/system/user')">
            <el-icon :size="24" color="#667eea"><User /></el-icon>
            <span>用户管理</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="quick-card" @click="$router.push('/system/role')">
            <el-icon :size="24" color="#764ba2"><UserFilled /></el-icon>
            <span>角色管理</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="quick-card" @click="$router.push('/system/audit')">
            <el-icon :size="24" color="#e6a23c"><Checked /></el-icon>
            <span>注册审核</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="quick-card">
            <el-icon :size="24" color="#67c23a"><DataAnalysis /></el-icon>
            <span>数据统计</span>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { User, Clock, Check, UserFilled, Checked, DataAnalysis } from '@element-plus/icons-vue'

const userStore = useUserStore()
const userInfo = computed(() => userStore.userInfo)

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const currentDate = ref('')
const currentWeek = ref('')

onMounted(() => {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  currentDate.value = `${y}年${m}月${d}日`
  const weeks = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  currentWeek.value = weeks[now.getDay()]
})
</script>

<style scoped>
.dashboard-view {
  font-family: 'Inter', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

/* 欢迎区域 */
.welcome-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: #fff;
  margin-bottom: 24px;
}

.welcome-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 8px;
  letter-spacing: 0.5px;
}

.welcome-desc {
  font-size: 14px;
  opacity: 0.85;
}

.welcome-date {
  text-align: right;
}

.date-text {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 4px;
}

.week-text {
  font-size: 20px;
  font-weight: 600;
}

/* 统计卡片 */
.stat-cards {
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px 28px;
  border-radius: 14px;
  background: #fff;
  transition: all 0.3s ease;
  cursor: default;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.card-blue .stat-icon {
  background: linear-gradient(135deg, #667eea, #764ba2);
}
.card-orange .stat-icon {
  background: linear-gradient(135deg, #f5af19, #f12711);
}
.card-green .stat-icon {
  background: linear-gradient(135deg, #11998e, #38ef7d);
}

.stat-value {
  font-size: 30px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #8c8ca1;
  margin-top: 4px;
  letter-spacing: 0.5px;
}

/* 快捷操作 */
.quick-section {
  margin-top: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  margin-bottom: 16px;
  letter-spacing: 0.5px;
}

.quick-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 24px 16px;
  background: #fff;
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.quick-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.quick-card span {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}
</style>
