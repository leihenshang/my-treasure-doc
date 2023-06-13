import { createWebHashHistory, createRouter } from 'vue-router';
import { useUserinfoStore } from "./stores/user/userinfo.js";


const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/LogIn' },
    { path: '/LogIn', name: 'LogIn', component: () => import('./views/LogIn.vue') },
    {
      path: '/HomePage', name: 'HomePage', component: () => import('./views/home/HomePage.vue'), redirect: { name: 'Collection' },
      children: [
        { path: '/Collection', name: 'Collection', component: () => import('./views/home/Collection.vue') },
        { path: '/Work', name: 'Work', component: () => import('./views/home/notes/Work.vue') },
        { path: '/Life', name: 'Life', component: () => import('./views/home/notes/Life.vue') },
        { path: '/Experience', name: 'Experience', component: () => import('./views/home/notes/Experience.vue') },
        { path: '/Plan', name: 'Plan', component: () => import('./views/home/Plan.vue') },
        { path: '/Diary', name: 'Diary', component: () => import('./views/home/Diary.vue') },
      ],
    }
  ]
})

router.beforeEach(async (to, from) => {
  let isAuthenticated = false
  const storeUserinfo = useUserinfoStore()
  if (storeUserinfo.userId > 0) {
    isAuthenticated = true
  }

  if (
    // 检查用户是否已登录
    !isAuthenticated &&
    // ❗️ 避免无限重定向
    to.name !== 'LogIn'
  ) {
    // 将用户重定向到登录页面
    return { name: 'LogIn' }
  }
})

export { router }