import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/Profile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/novel/:id',
    name: 'NovelDetail',
    component: () => import('@/views/novel/Detail.vue')
  },
  {
    path: '/read/:id',
    name: 'Reader',
    component: () => import('@/views/novel/Reader.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/upload',
    name: 'Upload',
    component: () => import('@/views/novel/Upload.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/category',
    name: 'Category',
    component: () => import('@/views/category/List.vue')
  },
  {
    path: '/ranking',
    name: 'Ranking',
    component: () => import('@/views/ranking/List.vue')
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/views/search/List.vue')
  },
  {
    path: '/admin/review',
    name: 'AdminReview',
    component: () => import('@/views/admin/Review.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/About.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next('/')
  } else {
    next()
  }
})

export default router