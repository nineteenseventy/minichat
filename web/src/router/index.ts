import { createAuthGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';
import { routes } from 'vue-router/auto-routes';
import auth0 from '../auth0';
import { initializeAuthenticatedUserStore } from '@/stores/authenticatedUser.store';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/callback',
      name: 'callback',
      redirect: '/',
      beforeEnter: async (to) => {
        const appState = await auth0.handleRedirectCallback(to.fullPath);
        return { path: appState.appState?.target };
      },
    },
    ...routes,
  ],
});

router.beforeEach((to) => {
  if (to.name === 'callback') return;
  return createAuthGuard()(to);
});

router.beforeEach(initializeAuthenticatedUserStore);

export default router;
