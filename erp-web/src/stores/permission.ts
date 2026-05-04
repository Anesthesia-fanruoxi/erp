import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMenus } from '@/api/auth'

export const usePermissionStore = defineStore('permission', () => {
  const menus = ref<MenuItem[]>([])
  const permissions = ref<string[]>([])
  const menuLoaded = ref(false)

  async function fetchMenus() {
    const res = await getMenus()
    menus.value = res.data.menus || []
    permissions.value = res.data.permissions || []
    menuLoaded.value = true
  }

  function clearPermission() {
    menus.value = []
    permissions.value = []
    menuLoaded.value = false
  }

  return { menus, permissions, menuLoaded, fetchMenus, clearPermission }
})
