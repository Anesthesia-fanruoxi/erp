<template>
  <div class="side-menu">
    <!-- Logo区域 -->
    <div class="logo-section">
      <div class="logo-icon">
        <el-icon :size="24"><Monitor /></el-icon>
      </div>
      <transition name="fade">
        <span class="logo-text">ERP System</span>
      </transition>
    </div>

    <!-- 菜单区域 -->
    <div class="menu-wrapper">
      <el-menu
        :default-active="$route.path"
        :background-color="'transparent'"
        text-color="rgba(255, 255, 255, 0.65)"
        active-text-color="#fff"
        router
        :unique-opened="true"
        class="custom-menu"
      >
        <el-menu-item index="/dashboard" class="menu-item-custom">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>

        <template v-for="menu in menus" :key="menu.id">
          <el-sub-menu
            v-if="menu.children && menu.children.length"
            :index="menu.path || String(menu.id)"
            class="sub-menu-custom"
          >
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
          <el-menu-item v-else :index="menu.path" class="menu-item-custom">
            <el-icon><component :is="menu.icon || 'Menu'" /></el-icon>
            <span>{{ menu.name }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </div>

    <!-- 底部版本信息 -->
    <div class="sidebar-footer">
      <span>v0.1.0</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePermissionStore } from '@/stores/permission'
import { House, Monitor } from '@element-plus/icons-vue'
import MenuItemRecursive from './MenuItemRecursive.vue'

const permissionStore = usePermissionStore()
const menus = computed(() => permissionStore.menus)
</script>

<style scoped>
.side-menu {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Logo区域 */
.logo-section {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.logo-text {
  font-size: 17px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1.5px;
  white-space: nowrap;
}

/* 菜单容器 */
.menu-wrapper {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px;
}

/* 滚动条样式 */
.menu-wrapper::-webkit-scrollbar {
  width: 4px;
}
.menu-wrapper::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 4px;
}
.menu-wrapper::-webkit-scrollbar-track {
  background: transparent;
}

/* 自定义菜单整体 */
.custom-menu {
  border-right: none !important;
}

/* 一级菜单项 */
.custom-menu :deep(.el-menu-item) {
  height: 44px;
  line-height: 44px;
  margin: 2px 0;
  border-radius: 10px;
  padding-left: 16px !important;
  font-size: 14px;
  letter-spacing: 0.5px;
  transition: all 0.25s ease;
  color: rgba(255, 255, 255, 0.65);
}

.custom-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.08) !important;
  color: #fff;
}

.custom-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.5), rgba(118, 75, 162, 0.4)) !important;
  color: #fff;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.25);
}

.custom-menu :deep(.el-menu-item .el-icon) {
  font-size: 18px;
  margin-right: 8px;
}

/* 子菜单标题 */
.custom-menu :deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  margin: 2px 0;
  border-radius: 10px;
  padding-left: 16px !important;
  font-size: 14px;
  letter-spacing: 0.5px;
  transition: all 0.25s ease;
  color: rgba(255, 255, 255, 0.65);
}

.custom-menu :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.08) !important;
  color: #fff;
}

.custom-menu :deep(.el-sub-menu__title .el-icon) {
  font-size: 18px;
  margin-right: 8px;
}

/* 子菜单箭头 */
.custom-menu :deep(.el-sub-menu__icon-arrow) {
  color: rgba(255, 255, 255, 0.35);
  font-size: 12px;
}

/* 展开的子菜单 */
.custom-menu :deep(.el-sub-menu.is-opened > .el-sub-menu__title) {
  color: #fff;
}

/* 子菜单内部项 */
.custom-menu :deep(.el-menu--inline .el-menu-item) {
  height: 40px;
  line-height: 40px;
  padding-left: 48px !important;
  font-size: 13px;
  border-radius: 8px;
  margin: 1px 4px;
}

/* 底部版本信息 */
.sidebar-footer {
  padding: 12px 0;
  text-align: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.2);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
  letter-spacing: 1px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
