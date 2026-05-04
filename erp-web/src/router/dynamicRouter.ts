import type { RouteRecordRaw } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'

// 路径到组件的映射表
const viewMap: Record<string, () => Promise<any>> = {
  '/dashboard': () => import('@/views/dashboard/DashboardView.vue'),
  '/system/user': () => import('@/views/system/UserManage.vue'),
  '/system/role': () => import('@/views/system/RoleManage.vue'),
  '/system/menu': () => import('@/views/system/MenuManage.vue'),
  '/system/audit': () => import('@/views/system/AuditManage.vue'),
  '/audit/log': () => import('@/views/audit/AuditLogView.vue'),
  // 业务管理
  '/business/sale': () => import('@/views/business/SaleView.vue'),
  '/business/purchase': () => import('@/views/business/PurchaseView.vue'),
  '/business/retail': () => import('@/views/business/RetailView.vue'),
  '/business/stock': () => import('@/views/business/StockView.vue'),
  '/business/contract': () => import('@/views/business/ContractView.vue'),  // 流程管理
  '/workflow/purchase': () => import('@/views/workflow/PurchaseFlowView.vue'),
  '/workflow/sale': () => import('@/views/workflow/SaleFlowView.vue'),
}

// 根据后端菜单树生成路由
export function generateRoutes(menus: MenuItem[]): RouteRecordRaw[] {
  const layoutRoute: RouteRecordRaw = {
    path: '/',
    component: AppLayout,
    children: []
  }

  function buildRoutes(items: MenuItem[], routes: RouteRecordRaw[]) {
    for (const item of items) {
      if (item.children && item.children.length > 0) {
        buildRoutes(item.children, routes)
      } else if (item.path) {
        const component = viewMap[item.path]
        if (component != null) {
          routes.push({
            path: item.path,
            name: item.name,
            component,
            meta: { title: item.name, icon: item.icon }
          })
        }
      }
    }
  }

  buildRoutes(menus, layoutRoute.children as RouteRecordRaw[])

  // 添加 dashboard 作为默认页
  const hasDashboard = layoutRoute.children!.some(r => r.path === '/dashboard')
  if (!hasDashboard) {
    layoutRoute.children!.unshift({
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/DashboardView.vue'),
      meta: { title: '首页', icon: 'House' }
    })
  }

  return [layoutRoute]
}
