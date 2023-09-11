import { createWebHashHistory, createRouter } from 'vue-router';
import { useUserinfoStore } from "./stores/user/userinfo";


const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/LogIn' },
    { path: '/LogIn', name: 'LogIn', component: () => import('./views/LogIn.vue') },
    {
      path: '/HomePage', name: 'HomePage',
      component: () => import('./views/home/HomePage.vue'),
      redirect: { name: 'Collection' },
      children: [
        { path: '/Collection', name: 'Collection', component: () => import('./views/home/Collection.vue') },
        { path: '/Plan', name: 'Plan', component: () => import('./views/home/Plan.vue') },
        { path: '/Write', name: 'Write', component: () => import('./views/home/Write.vue') },
        { path: '/New', name: 'New', component: () => import('./views/home/Note.vue') },
        { path: '/Edit', name: 'Edit', component: () => import('./views/home/Edit.vue') },
      ],
    }
  ]
})

router.beforeEach(async (to, from) => {
  let isAuthenticated = false


  type TypeUserInfo = {
    id: number
  }

  const storeUserinfo = useUserinfoStore()
  if (storeUserinfo.userId > 0) {
    isAuthenticated = true

    // use local storage user info
  } else {
    let localUserInfo: string | null = localStorage.getItem('userInfo')
    let userInfo: TypeUserInfo = { id: 0 }
    if (localUserInfo) {
      userInfo = JSON.parse(localUserInfo)
      if (userInfo && userInfo.id > 0) {
        storeUserinfo.updateUserinfo(userInfo)
        isAuthenticated = true
      }
    }
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