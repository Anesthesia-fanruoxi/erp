import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'
import { useUserStore } from '@/stores/user'
import { generateRoutes } from './dynamicRouter'

// 静态路由(不需要权限)
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/RegisterView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const permissionStore = usePermissionStore()
  const userStore = useUserStore()

  // 已登录访问登录/注册页，重定向到首页
  if ((to.path === '/login' || to.path === '/register') && authStore.token) {
    next('/dashboard')
    return
  }

  // 不需要认证的页面直接放行
  if (to.meta.requiresAuth === false) {
    next()
    return
  }

  // 未登录跳转登录页
  if (!authStore.token) {
    next('/login')
    return
  }

  // 已登录但未加载动态菜单
  if (permissionStore.menus.length === 0) {
    try {
      await permissionStore.fetchMenus()
      const dynamicRoutes = generateRoutes(permissionStore.menus)
      dynamicRoutes.forEach(route => router.addRoute(route))
      await userStore.fetchUserInfo()
      next({ ...to, replace: true })
    } catch {
      // 加载失败清除登录状态
      authStore.clearAuth()
      next('/login')
    }
    return
  }

  next()
})

export default router
