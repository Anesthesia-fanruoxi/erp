<template>
  <header class="top-header">
    <div class="header-left">
      <div class="breadcrumb-area">
        <el-icon :size="18" color="#8c8ca1"><HomeFilled /></el-icon>
        <span class="current-page">{{ currentRouteName }}</span>
      </div>
    </div>

    <div class="header-right">
      <!-- 当前时间 -->
      <div class="header-time">
        <el-icon :size="14"><Clock /></el-icon>
        <span>{{ currentTime }}</span>
      </div>

      <div class="header-divider"></div>

      <!-- 用户下拉 -->
      <el-dropdown @command="handleCommand" trigger="click">
        <div class="user-dropdown">
          <div class="user-avatar">
            <el-icon :size="16"><UserFilled /></el-icon>
          </div>
          <span class="user-name">{{ userInfo?.realName || '用户' }}</span>
          <el-icon :size="12" class="arrow-icon"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu class="user-dropdown-menu">
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>个人信息
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!-- 个人信息弹窗 -->
    <el-dialog v-model="profileVisible" title="个人信息" width="480px" @closed="resetProfileForm">
      <el-form
        ref="profileFormRef"
        :model="profileForm"
        :rules="profileRules"
        label-width="80px"
        v-loading="profileLoading"
      >
        <el-form-item label="用户名">
          <el-input :model-value="profileForm.userName" disabled />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="profileForm.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机" prop="phone">
          <el-input v-model="profileForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input
            v-model="profileForm.password"
            type="password"
            show-password
            placeholder="不修改请留空"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="profileVisible = false">取消</el-button>
        <el-button type="primary" :loading="profileSubmitting" @click="handleProfileSubmit">保存</el-button>
      </template>
    </el-dialog>
  </header>
</template>

<script setup lang="ts">
import { computed, ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { useTabsStore } from '@/stores/tabs'
import { logout } from '@/api/auth'
import { getProfile, updateProfile } from '@/api/auth'
import { ElMessage } from 'element-plus'
import { UserFilled, ArrowDown, HomeFilled, Clock, User, SwitchButton } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()
const tabsStore = useTabsStore()
const userInfo = computed(() => userStore.userInfo)

const currentRouteName = computed(() => {
  return (route.meta?.title as string) || '首页'
})

// 实时时间
const currentTime = ref('')
let timer: number | null = null

function updateTime() {
  const now = new Date()
  const h = String(now.getHours()).padStart(2, '0')
  const m = String(now.getMinutes()).padStart(2, '0')
  const s = String(now.getSeconds()).padStart(2, '0')
  currentTime.value = `${h}:${m}:${s}`
}

onMounted(() => {
  updateTime()
  timer = window.setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

// ---- 个人信息弹窗 ----
const profileVisible = ref(false)
const profileLoading = ref(false)
const profileSubmitting = ref(false)
const profileFormRef = ref<FormInstance>()

const profileForm = reactive({
  userName: '',
  realName: '',
  email: '',
  phone: '',
  password: ''
})

const profileRules: FormRules = {
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }]
}

async function openProfile() {
  profileVisible.value = true
  profileLoading.value = true
  try {
    const res = await getProfile()
    const data = res.data
    profileForm.userName = data.userName || ''
    profileForm.realName = data.realName || ''
    profileForm.email = data.email || ''
    profileForm.phone = data.phone || ''
    profileForm.password = ''
  } catch (err: any) {
    ElMessage.error(err.message || '获取个人信息失败')
  } finally {
    profileLoading.value = false
  }
}

async function handleProfileSubmit() {
  if (!profileFormRef.value) return
  await profileFormRef.value.validate(async (valid) => {
    if (!valid) return
    profileSubmitting.value = true
    try {
      const payload: any = {
        realName: profileForm.realName,
        email: profileForm.email,
        phone: profileForm.phone
      }
      if (profileForm.password) {
        payload.password = profileForm.password
      }
      await updateProfile(payload)
      ElMessage.success('个人信息更新成功')
      profileVisible.value = false
      // 刷新用户信息
      userStore.fetchUserInfo()
    } catch (err: any) {
      ElMessage.error(err.message || '更新失败')
    } finally {
      profileSubmitting.value = false
    }
  })
}

function resetProfileForm() {
  profileForm.userName = ''
  profileForm.realName = ''
  profileForm.email = ''
  profileForm.phone = ''
  profileForm.password = ''
  profileFormRef.value?.resetFields()
}

async function handleCommand(command: string) {
  if (command === 'profile') {
    openProfile()
  } else if (command === 'logout') {
    try {
      await logout()
    } finally {
      authStore.clearAuth()
      userStore.clearUser()
      permissionStore.clearPermission()
      tabsStore.clearTabs()
      router.push('/login')
      ElMessage.success('已退出登录')
    }
  }
}
</script>

<style scoped>
.top-header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: #fff;
  font-family: 'Inter', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.header-left {
  display: flex;
  align-items: center;
}

.breadcrumb-area {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-page {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a2e;
  letter-spacing: 0.5px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-time {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #8c8ca1;
  font-variant-numeric: tabular-nums;
}

.header-divider {
  width: 1px;
  height: 20px;
  background: #e8e8ef;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 10px;
  transition: all 0.25s ease;
}

.user-dropdown:hover {
  background: #f5f5fa;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  letter-spacing: 0.5px;
}

.arrow-icon {
  color: #a0a3bd;
  transition: transform 0.3s;
}
</style>
