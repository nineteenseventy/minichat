import { createAuthGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/callback',
      name: 'callback',
      component: () => import('../views/CallbackView.vue'),
    },
    {
      path: '',
      beforeEnter: createAuthGuard(),
      children: [
        {
          path: '/',
          name: 'about',
          component: () => import('../views/AppView.vue'),
        },
      ],
    },
  ],
});

export default router;
