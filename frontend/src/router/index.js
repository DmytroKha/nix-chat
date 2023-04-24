// import Vue from 'vue'
import {
  createRouter,
  // createWebHashHistory,
  createWebHistory
} from 'vue-router'

import LoginPage from "../pages/Login.vue";
import RegisterPage from "../pages/Register.vue";
import ProfilePage from "../pages/Profile.vue";
import DashBoard from "../pages/DashBoard.vue";

// Vue.use(VueRouter)

const routes = [{
    path: '/',
    name: 'Login',
    component: LoginPage
  },
  {
    path: '/Register',
    name: 'Register',
    component: RegisterPage
  },
  {
    path: '/profile',
    name: 'profile',
    component: ProfilePage
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashBoard,
  },
]

const router = createRouter({
  // history: createWebHashHistory(),
  history: createWebHistory(),
  routes
})

// router.beforeEach((to, from, next) => {
//   const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
//   const isAuthenticated = firebase.auth().currentUser;
//   console.log("isauthenticated", isAuthenticated);
//   if (requiresAuth && !isAuthenticated) {
//       next("/login");
//   } else {
//       next();
//   }
// });

export default router