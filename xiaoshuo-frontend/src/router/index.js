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
    meta: { requiresAuth: true },
    children: [
      {
        path: 'basic',
        name: 'BasicInfo',
        component: () => import('@/views/user/BasicInfo.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'uploads',
        name: 'Uploads',
        component: () => import('@/views/user/Uploads.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'comments',
        name: 'Comments',
        component: () => import('@/views/user/Comments.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'ratings',
        name: 'Ratings',
        component: () => import('@/views/user/Ratings.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'social',
        name: 'SocialHistoryPage',
        component: () => import('@/views/user/SocialHistory.vue'),
        meta: { requiresAuth: true }
      }
    ]
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
    path: '/admin/standard',
    name: 'AdminStandard',
    component: () => import('@/views/admin/Standard.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/monitor',
    name: 'AdminMonitor',
    component: () => import('@/views/admin/Monitor.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/novel/:id/social-history',
    name: 'SocialHistory',
    component: () => import('@/views/novel/SocialHistory.vue'),
    meta: { requiresAuth: true }
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