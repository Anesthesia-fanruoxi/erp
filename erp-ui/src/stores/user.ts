import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCurrentUser } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)

  async function fetchUserInfo() {
    const res = await getCurrentUser()
    userInfo.value = res.data
  }

  function clearUser() {
    userInfo.value = null
  }

  return { userInfo, fetchUserInfo, clearUser }
})
