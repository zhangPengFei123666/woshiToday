import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/store/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/components/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台', icon: 'Odometer' }
      },
      {
        path: 'group',
        name: 'Group',
        component: () => import('@/views/group/index.vue'),
        meta: { title: '任务组管理', icon: 'Folder' }
      },
      {
        path: 'task',
        name: 'Task',
        component: () => import('@/views/task/index.vue'),
        meta: { title: '任务管理', icon: 'List' }
      },
      {
        path: 'instance',
        name: 'Instance',
        component: () => import('@/views/instance/index.vue'),
        meta: { title: '执行记录', icon: 'Document' }
      },
      {
        path: 'executor',
        name: 'Executor',
        component: () => import('@/views/executor/index.vue'),
        meta: { title: '执行器管理', icon: 'Monitor' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title || '分布式任务调度系统'} - Scheduler`
  
  const userStore = useUserStore()
  const token = userStore.token
  
  if (to.meta.requiresAuth === false) {
    // 不需要认证的页面
    if (token && to.path === '/login') {
      next('/dashboard')
    } else {
      next()
    }
  } else {
    // 需要认证的页面
    if (token) {
      next()
    } else {
      next('/login')
    }
  }
})

export default router

