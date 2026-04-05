import { createRouter, createWebHistory, type RouteLocationNormalized } from 'vue-router'
import service from '@/utils/request'
// 注意：这里不再导入具体的 Vue 组件文件

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      // 修改为动态导入
      component: () => import('@/views/blog/HomePage.vue'),
      meta: { title: 'Home' },
    },
    {
      path: '/blog',
      name: 'blog',
      component: () => import('@/views/blog/BlogListFinal.vue'),
      meta: { title: 'Blog' },
    },
    {
      path: '/thread',
      name: 'thread',
      component: () => import('@/views/tweet/TweetIndexPage.vue'),
      meta: { title: 'Thread' },
    },
    {
      path: '/blog/:slug',
      name: 'PostDetailBySlug',
      component: () => import('@/views/blog/BlogDetailPageWithSlug.vue'),
    },

    // --- Admin 管理路由 ---
    // 这些通常体积较大，拆分后效果最明显
    {
      path: '/admin',
      name: 'login',
      component: () => import('@/views/admin/AdminLogin.vue'),
      meta: { title: 'Admin Login' },
    },
    {
      path: '/admin/post/create',
      name: 'create',
      component: () => import('@/views/post/AdminCreatePost.vue'),
      meta: { title: 'Create Post' },
    },
    {
      path: '/admin/update/:id',
      name: 'update',
      component: () => import('@/views/post/AdminUpdatePost.vue'),
      meta: { title: 'Update Post' },
    },
    {
      path: '/admin/category/create',
      name: 'CreateCategory',
      component: () => import('@/views/category/CreateCategory.vue'),
      meta: { title: 'Create Category' },
    },
    {
      path: '/admin/dashboard',
      name: 'AdminIndex',
      component: () => import('@/views/admin/AdminIndexPage.vue'),
      meta: { title: 'Admin Index' },
    },
    {
      path: '/admin/changeUserInfo',
      name: 'ChangeUserInfo',
      component: () => import('@/views/admin/user-change-info.vue'),
      meta: { title: '修改用户信息' },
    },

    // 404
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundPage.vue'),
      meta: { title: '页面走丢了' },
    },
  ],
})

router.beforeEach(async (to: RouteLocationNormalized, from, next) => {
  const token = localStorage.getItem('access_token')

  // 逻辑保持不变
  const needAuth =
    to.matched.some((record) => record.meta.requiresAuth) ||
    (to.path.startsWith('/admin') && to.path !== '/admin')

  if (needAuth) {
    if (token) {
      if (from.path === '/admin') {
        next()
        return
      }
      try {
        // 建议加上 try-catch 防止刷新 token 失败导致路由卡死
        const newToken = await service.post('/auth/refresh')
        if (newToken && newToken.data?.access_token) {
          localStorage.setItem('access_token', newToken.data.access_token)
          console.log('刷新token成功')
        }
        next()
      } catch (error) {
        console.error('Token refresh failed', error)
        localStorage.removeItem('access_token')
        next({
          path: '/admin',
          query: { redirect: to.fullPath },
        })
      }
    } else {
      next({
        path: '/admin',
        query: { redirect: to.fullPath },
      })
    }
  } else {
    next()
  }
})

export default router
