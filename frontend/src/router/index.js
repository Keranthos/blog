import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => (import('@/views/HomeView.vue'))
  },
  {
    path: '/blog',
    name: 'Blog',
    component: () => (import('@/views/BlogView.vue'))
  },
  {
    path: '/questionbox',
    name: 'QuestionBox',
    component: () => (import('@/views/QuestionboxView.vue'))
  },
  {
    path: '/moments',
    name: 'Moments',
    component: () => (import('@/views/MomentsView.vue'))
  },
  {
    path: '/moments/:id',
    name: 'MomentDetail',
    component: () => (import('@/views/MomentDetailView.vue'))
  },
  {
    path: '/search',
    name: 'Search',
    component: () => (import('@/views/SearchView.vue'))
  },
  {
    path: '/editarticle',
    name: 'editarticle',
    component: () => (import('@/views/EditArticleView.vue'))
  },
  {
    path: '/location-update',
    name: 'LocationUpdate',
    component: () => (import('@/views/LocationUpdateView.vue'))
  },
  {
    path: '/edit/:id',
    name: 'edit',
    component: () => (import('@/views/EditArticleView.vue'))
  },
  {
    path: '/blog/:id',
    name: 'BlogDetail',
    component: () => import('@/views/BlogDetailView.vue'), // ：id怎么处理？
    props: true
  },
  { path: '/fragments/books', name: 'Books', component: () => import('@/views/BooksView.vue') },
  { path: '/fragments/novels', name: 'Novels', component: () => import('@/views/NovelsView.vue') },
  { path: '/fragments/movies', name: 'Movies', component: () => import('@/views/MoviesView.vue') },
  { path: '/login-register', name: 'LoginRegister', component: () => import('@/views/LoginRegisterView.vue') },
  { path: '/user-info', name: 'UserInfo', component: () => import('@/views/LoginRegisterView.vue') },
  { path: '/profile', name: 'Profile', component: () => import('@/views/ProfileView.vue') },
  { path: '/images', name: 'ImageManage', component: () => import('@/views/ImageManageView.vue') },
  { path: '/timeline', name: 'Timeline', component: () => import('@/views/TimelineView.vue') },
  { path: '/loading-test', name: 'LoadingTest', component: () => import('@/views/LoadingTestView.vue') },
  { path: '/error', name: 'Error', component: () => import('@/views/ErrorView.vue') },
  { path: '/:catchAll(.*)', name: 'NotFound', component: () => import('@/views/NotFoundView.vue') }
]

const router = createRouter({
  history: createWebHistory('/'),
  routes
})

export default router
