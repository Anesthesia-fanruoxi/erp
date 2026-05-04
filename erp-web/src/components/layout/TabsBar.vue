<template>
  <div class="tabs-bar">
    <el-tabs
      v-model="tabsStore.activeTab"
      type="card"
      class="tabs-nav"
      @tab-click="handleTabClick"
      @tab-remove="handleTabRemove"
    >
      <el-tab-pane
        v-for="tab in tabsStore.tabs"
        :key="tab.path"
        :label="tab.title"
        :name="tab.path"
        :closable="tab.path !== '/dashboard'"
      >
        <template #label>
          <span class="tab-label">
            <el-icon v-if="tab.icon" class="tab-icon"><component :is="tab.icon" /></el-icon>
            {{ tab.title }}
          </span>
        </template>
      </el-tab-pane>
    </el-tabs>

    <!-- 右侧操作 -->
    <el-dropdown trigger="click" @command="handleCommand">
      <el-button size="small" class="tabs-action-btn">
        <el-icon><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="closeOthers">关闭其他</el-dropdown-item>
          <el-dropdown-item command="closeAll">关闭全部</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { useTabsStore } from '@/stores/tabs'
import { useRouter } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import type { TabsPaneContext } from 'element-plus'

const tabsStore = useTabsStore()
const router = useRouter()

function handleTabClick(pane: TabsPaneContext) {
  const path = pane.paneName as string
  if (router.currentRoute.value.path !== path) {
    router.push(path)
  }
}

function handleTabRemove(path: string) {
  tabsStore.removeTab(path)
  // 如果关闭后 activeTab 变了，跳转
  if (router.currentRoute.value.path === path) {
    router.push(tabsStore.activeTab)
  }
}

function handleCommand(cmd: string) {
  if (cmd === 'closeOthers') {
    const current = tabsStore.activeTab
    tabsStore.tabs.splice(1) // 保留首页
    const dashboard = { path: '/dashboard', title: '首页', icon: 'House' }
    tabsStore.tabs.splice(0, 1, dashboard)
    if (current !== '/dashboard') {
      tabsStore.addTab({
        path: current,
        title: tabsStore.tabs.find(t => t.path === current)?.title || current
      })
    }
    router.push(tabsStore.activeTab)
  } else if (cmd === 'closeAll') {
    tabsStore.clearTabs()
    router.push('/dashboard')
  }
}
</script>

<style scoped>
.tabs-bar {
  display: flex;
  align-items: center;
  background: #fff;
  border-bottom: 2px solid #e8eaf0;
  padding: 0 8px 0 0;
  height: 40px;
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
}

.tabs-nav {
  flex: 1;
  overflow: hidden;
}

/* 覆盖 el-tabs card 样式 */
.tabs-nav :deep(.el-tabs__header) {
  margin: 0;
  border-bottom: none;
}

.tabs-nav :deep(.el-tabs__nav-wrap) {
  margin-bottom: 0;
}

.tabs-nav :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.tabs-nav :deep(.el-tabs__nav) {
  border: none;
}

.tabs-nav :deep(.el-tabs__item) {
  height: 36px;
  line-height: 36px;
  border: 1px solid #d0d3e0 !important;
  border-bottom: none !important;
  border-radius: 6px 6px 0 0;
  font-size: 13px;
  color: #303133;
  padding: 0 14px;
  margin: 4px 2px 0;
  transition: all 0.2s;
  background: #eef0f6;
}

.tabs-nav :deep(.el-tabs__item:hover) {
  color: #667eea;
  background: #e4e8ff;
  border-color: #b0b8f0 !important;
}

.tabs-nav :deep(.el-tabs__item.is-active) {
  color: #fff;
  background: linear-gradient(135deg, #667eea, #764ba2);
  font-weight: 600;
  border-color: #667eea !important;
  box-shadow: 0 -2px 8px rgba(102, 126, 234, 0.3);
}

.tabs-nav :deep(.el-tabs__item.is-active .is-icon-close) {
  color: rgba(255, 255, 255, 0.8);
}

.tabs-nav :deep(.el-tabs__item.is-active .is-icon-close:hover) {
  color: #fff;
  background: rgba(255, 255, 255, 0.2);
}

.tabs-nav :deep(.el-tabs__item .is-icon-close) {
  font-size: 12px;
  margin-left: 4px;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 5px;
}

.tab-icon {
  font-size: 13px;
}

.tabs-action-btn {
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  padding: 4px 8px;
  margin-left: 6px;
  flex-shrink: 0;
  background: #f8f9fc;
  color: #606266;
}

.tabs-action-btn:hover {
  border-color: #667eea;
  color: #667eea;
  background: #eef1ff;
}
</style>
