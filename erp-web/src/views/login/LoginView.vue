<template>
  <div class="login-view">
    <!-- 动态背景粒子 -->
    <div class="bg-particles">
      <div v-for="n in 6" :key="n" :class="'particle p' + n"></div>
    </div>

    <div class="login-container">
      <!-- 左侧品牌区域 -->
      <div class="brand-panel">
        <div class="brand-content">
          <div class="brand-icon">
            <el-icon :size="48"><Monitor /></el-icon>
          </div>
          <h1 class="brand-title">ERP System</h1>
          <p class="brand-subtitle">企业资源管理平台</p>
          <div class="brand-features">
            <div class="feature-item">
              <el-icon :size="18"><Check /></el-icon>
              <span>高效的企业资源管理</span>
            </div>
            <div class="feature-item">
              <el-icon :size="18"><Check /></el-icon>
              <span>安全的权限控制体系</span>
            </div>
            <div class="feature-item">
              <el-icon :size="18"><Check /></el-icon>
              <span>灵活的角色配置管理</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录表单 -->
      <div class="form-panel">
        <div class="form-content">
          <div class="form-header">
            <h2>欢迎回来</h2>
            <p>请登录您的账号</p>
          </div>

          <el-form ref="formRef" :model="form" :rules="rules" label-width="0" size="large">
            <el-form-item prop="userName">
              <el-input
                v-model="form.userName"
                placeholder="请输入用户名"
                :prefix-icon="User"
                class="custom-input"
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="form.password"
                type="password"
                placeholder="请输入密码"
                :prefix-icon="Lock"
                show-password
                class="custom-input"
                @keyup.enter="handleLogin"
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                :loading="loading"
                class="login-btn"
                @click="handleLogin"
              >
                <span v-if="!loading">登 录</span>
                <span v-else>登录中...</span>
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { login } from '@/api/auth'
import { ElMessage } from 'element-plus'
import { User, Lock, Monitor, Check } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const formRef = ref<FormInstance>()
const form = ref<LoginForm>({
  userName: '',
  password: ''
})

const rules: FormRules = {
  userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不少于6位', trigger: 'blur' }
  ]
}

async function handleLogin() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const res = await login(form.value)
      authStore.setToken(res.data.token)
      ElMessage.success('登录成功')
      router.push('/dashboard')
    } catch (err: any) {
      ElMessage.error(err.message || '登录失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0c29 0%, #1a1a3e 40%, #24243e 100%);
  position: relative;
  overflow: hidden;
  font-family: 'Inter', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

/* 背景粒子动画 */
.bg-particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.particle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.15;
  animation: float 15s infinite ease-in-out;
}

.p1 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, #667eea, transparent);
  top: -5%;
  right: -5%;
  animation-delay: 0s;
}
.p2 {
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, #764ba2, transparent);
  bottom: 10%;
  left: -3%;
  animation-delay: -3s;
}
.p3 {
  width: 150px;
  height: 150px;
  background: radial-gradient(circle, #4facfe, transparent);
  top: 40%;
  left: 30%;
  animation-delay: -6s;
}
.p4 {
  width: 250px;
  height: 250px;
  background: radial-gradient(circle, #00f2fe, transparent);
  bottom: -8%;
  right: 20%;
  animation-delay: -9s;
}
.p5 {
  width: 120px;
  height: 120px;
  background: radial-gradient(circle, #a18cd1, transparent);
  top: 15%;
  left: 10%;
  animation-delay: -12s;
}
.p6 {
  width: 180px;
  height: 180px;
  background: radial-gradient(circle, #667eea, transparent);
  top: 60%;
  right: 10%;
  animation-delay: -4s;
}

@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  33% { transform: translateY(-20px) scale(1.05); }
  66% { transform: translateY(15px) scale(0.95); }
}

/* 登录容器 */
.login-container {
  display: flex;
  width: 880px;
  min-height: 520px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(10px);
  z-index: 1;
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 左侧品牌面板 */
.brand-panel {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  position: relative;
  overflow: hidden;
}

.brand-panel::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1), transparent);
  border-radius: 50%;
}

.brand-content {
  position: relative;
  z-index: 1;
  text-align: center;
  color: #fff;
}

.brand-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.brand-title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 2px;
  margin-bottom: 8px;
}

.brand-subtitle {
  font-size: 14px;
  opacity: 0.85;
  margin-bottom: 36px;
  letter-spacing: 4px;
}

.brand-features {
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  opacity: 0.9;
  letter-spacing: 0.5px;
}

/* 右侧表单面板 */
.form-panel {
  flex: 1;
  background: rgba(255, 255, 255, 0.98);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 44px;
}

.form-content {
  width: 100%;
  max-width: 340px;
}

.form-header {
  margin-bottom: 36px;
}

.form-header h2 {
  font-size: 26px;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.form-header p {
  font-size: 14px;
  color: #8c8ca1;
}

/* 自定义输入框 */
.custom-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  padding: 4px 16px;
  box-shadow: 0 0 0 1px #e4e7ed inset;
  transition: all 0.3s ease;
  background: #f8f9fc;
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #b4bccc inset;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1.5px #667eea inset;
  background: #fff;
}

.custom-input :deep(.el-input__prefix .el-icon) {
  color: #a0a3bd;
  font-size: 18px;
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 46px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 4px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.login-btn:active {
  transform: translateY(0);
}

/* 覆盖 el-form-item 间距 */
:deep(.el-form-item) {
  margin-bottom: 22px;
}
</style>
