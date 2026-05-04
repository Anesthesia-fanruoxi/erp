import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMenus } from '@/api/auth'

export const usePermissionStore = defineStore('permission', () => {
  const menus = ref<MenuItem[]>([])
  const permissions = ref<string[]>([])

  async function fetchMenus() {
    const res = await getMenus()
    menus.value = res.data.menus || []
    permissions.value = res.data.permissions || []
  }

  function clearPermission() {
    menus.value = []
    permissions.value = []
  }

  return { menus, permissions, fetchMenus, clearPermission }
})
