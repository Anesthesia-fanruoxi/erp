import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface TabItem {
  path: string
  title: string
  icon?: string
}

export const useTabsStore = defineStore('tabs', () => {
  const tabs = ref<TabItem[]>([
    { path: '/dashboard', title: '首页', icon: 'House' }
  ])
  const activeTab = ref('/dashboard')

  function addTab(tab: TabItem) {
    const exists = tabs.value.find(t => t.path === tab.path)
    if (!exists) {
      tabs.value.push(tab)
    }
    activeTab.value = tab.path
  }

  function removeTab(path: string) {
    const idx = tabs.value.findIndex(t => t.path === path)
    if (idx === -1) return

    tabs.value.splice(idx, 1)

    // 如果关闭的是当前激活的标签，跳到相邻标签
    if (activeTab.value === path) {
      const next = tabs.value[idx] ?? tabs.value[idx - 1]
      if (next) activeTab.value = next.path
    }
  }

  function clearTabs() {
    tabs.value = [{ path: '/dashboard', title: '首页', icon: 'House' }]
    activeTab.value = '/dashboard'
  }

  return { tabs, activeTab, addTab, removeTab, clearTabs }
})
