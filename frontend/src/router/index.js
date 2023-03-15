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
    // name: 'Login',
    component: LoginPage
  },
  {
    path: '/Register',
    component: RegisterPage
  },
  {
    path: '/profile',
    component: ProfilePage
  },
  {
    path: '/dashboard',
    component: DashBoard,
    // meta: { requiresAuth: true },
    // children: [{
    //   path: "/dashboard/profile",
    //   component: ProfilePage,
    //   // meta: { requiresAuth: true }
    // }, ]
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