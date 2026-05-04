<template>
  <el-sub-menu v-if="item.children && item.children.length" :index="item.path || String(item.id)">
    <template #title>
      <el-icon><component :is="item.icon || 'Menu'" /></el-icon>
      <span>{{ item.name }}</span>
    </template>
    <menu-item-recursive
      v-for="child in item.children"
      :key="child.id"
      :item="child"
    />
  </el-sub-menu>
  <el-menu-item v-else :index="item.path">
    <el-icon><component :is="item.icon || 'Menu'" /></el-icon>
    <span>{{ item.name }}</span>
  </el-menu-item>
</template>

<script setup lang="ts">
defineProps<{
  item: MenuItem
}>()
</script>
