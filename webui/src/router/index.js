import { createRouter, createWebHistory } from 'vue-router'
import PublicLayout from '@/layouts/PublicLayout.vue'
import PrivateLayout from '@/layouts/PrivateLayout.vue'

const routes = [
  {
    path: '/',
    component: PublicLayout,
    children: [
      { path: '', name: 'login', component: () => import('@/views/Login.vue') },
    ],
  },
  {
    path: '/app',
    component: PrivateLayout,
    children: [
      { path: '',                  name: 'dashboard',       component: () => import('@/views/Dashboard.vue') },
      { path: 'profile',           name: 'my-profile',      component: () => import('@/views/MyProfile.vue') },
      { path: 'users/:id',         name: 'user-profile',    component: () => import('@/views/UserProfile.vue') },
      { path: 'conversations',     name: 'my-conversations', component: () => import('@/views/MyConversations.vue') },
      { path: 'conversations/:id', name: 'conversation',     component: () => import('@/views/Conversation.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, from, next) => {
  const loggedIn = !!localStorage.getItem('userId')
  if (to.path.startsWith('/app') && !loggedIn) {
    return next({ name: 'login' })
  }
  if (to.name === 'login' && loggedIn) {
    return next({ name: 'dashboard' })
  }
  next()
})

export default router