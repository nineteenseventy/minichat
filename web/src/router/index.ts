import { createAuthGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';
import { routes } from 'vue-router/auto-routes';
import { initializeAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import { globalAuth0 } from '@/plugins/auth0';
import CallbackErrorView from '@/view/CallbackErrorView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/callback',
      name: 'callback',
      component: CallbackErrorView,
      beforeEnter: async (to) => {
        const appState = await globalAuth0.handleRedirectCallback(to.fullPath);
        console.log(globalAuth0.error);
        if (globalAuth0.error) return { params: { error: globalAuth0.error } };
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
