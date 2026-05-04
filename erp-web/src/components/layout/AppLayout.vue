<template>
  <el-container class="app-layout">
    <el-aside width="230px" class="sidebar">
      <SideMenu />
    </el-aside>
    <el-container class="main-container">
      <el-header class="top-header-wrap">
        <TopHeader />
      </el-header>

      <!-- 标签栏 -->
      <TabsBar />

      <el-main class="main-content">
        <div class="content-wrapper">
          <router-view v-slot="{ Component, route }">
            <keep-alive :include="cachedViews">
              <component :is="Component" :key="route.path" />
            </keep-alive>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SideMenu from './SideMenu.vue'
import TopHeader from './TopHeader.vue'
import TabsBar from './TabsBar.vue'
import { useTabsStore } from '@/stores/tabs'

const tabsStore = useTabsStore()

// keep-alive 缓存已打开的标签页对应的组件名
const cachedViews = computed(() => tabsStore.tabs.map(t => t.path))
</script>

<style scoped>
.app-layout {
  height: 100vh;
  font-family: 'Inter', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.sidebar {
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  overflow: hidden;
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.15);
  z-index: 10;
}

.main-container {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: 100%;
}

.top-header-wrap {
  padding: 0;
  height: 64px;
  flex-shrink: 0;
  border-bottom: 1px solid #f0f0f5;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  z-index: 5;
}

.main-content {
  background: #f5f6fa;
  overflow: auto;
  padding: 0;
  flex: 1;
  height: 0;
}

.content-wrapper {
  padding: 20px;
  min-height: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.content-wrapper > * {
  flex-shrink: 0;
}

.content-wrapper > :last-child {
  flex: 1;
}
</style>
