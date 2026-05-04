<template>
  <div class="side-menu">
    <div class="logo">ERP System</div>
    <el-menu
      :default-active="$route.path"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
      router
    >
      <el-menu-item index="/dashboard">
        <el-icon><House /></el-icon>
        <span>首页</span>
      </el-menu-item>
      <template v-for="menu in menus" :key="menu.id">
        <el-sub-menu v-if="menu.children && menu.children.length" :index="menu.path || String(menu.id)">
          <template #title>
            <el-icon><component :is="menu.icon || 'Menu'" /></el-icon>
            <span>{{ menu.name }}</span>
          </template>
          <menu-item-recursive
            v-for="child in menu.children"
            :key="child.id"
            :item="child"
          />
        </el-sub-menu>
        <el-menu-item v-else :index="menu.path">
          <el-icon><component :is="menu.icon || 'Menu'" /></el-icon>
          <span>{{ menu.name }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePermissionStore } from '@/stores/permission'
import { House } from '@element-plus/icons-vue'
import MenuItemRecursive from './MenuItemRecursive.vue'

const permissionStore = usePermissionStore()
const menus = computed(() => permissionStore.menus)
</script>

<style scoped>
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}
.el-menu {
  border-right: none;
}
</style>
