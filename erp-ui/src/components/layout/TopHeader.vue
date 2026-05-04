<template>
  <header class="top-header">
    <div class="left">
      <!-- breadcrumb placeholder -->
    </div>
    <div class="right">
      <el-dropdown @command="handleCommand">
        <span class="user-dropdown">
          <el-icon><UserFilled /></el-icon>
          {{ userInfo?.realName || '用户' }}
          <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人信息</el-dropdown-item>
            <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
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
import { computed, ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { logout, getProfile, updateProfile } from '@/api/auth'
import { ElMessage } from 'element-plus'
import { UserFilled, ArrowDown } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()
const userInfo = computed(() => userStore.userInfo)

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
      router.push('/login')
      ElMessage.success('已退出登录')
    }
  }
}
</script>

<style scoped>
.top-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: #fff;
}
.user-dropdown {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: #303133;
  font-size: 14px;
}
</style>
