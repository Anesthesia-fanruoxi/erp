<template>
  <el-container class="app-layout">
    <el-aside width="220px" class="sidebar">
      <SideMenu />
    </el-aside>
    <el-container>
      <el-header class="top-header-wrap">
        <TopHeader />
      </el-header>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import SideMenu from './SideMenu.vue'
import TopHeader from './TopHeader.vue'
</script>

<style scoped>
.app-layout {
  height: 100vh;
}
.sidebar {
  background: #304156;
  overflow-y: auto;
}
.top-header-wrap {
  padding: 0;
  height: 60px;
  border-bottom: 1px solid #e6e6e6;
  background: #fff;
}
.main-content {
  background: #f0f2f5;
  overflow: auto;
}
</style>
