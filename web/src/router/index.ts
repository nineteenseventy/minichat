import { createAuthGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';
import { initializeAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { globalAuth0 } from '@/plugins/auth0';
import CallbackErrorView from '@/view/CallbackErrorView.vue';
import MainView from '@/views/MainView.vue';
import ChatView from '@/views/ChannelView.vue';

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
    {
      path: '/',
      name: 'home',
      component: MainView,
      children: [
        {
          path: 'channels',
          redirect: '/',
        },
        {
          path: 'channels/:channelId',
          name: 'channels',
          component: ChatView,
        },
      ],
    },
  ],
});

router.beforeEach((to) => {
  if (to.name === 'callback') return;
  return createAuthGuard()(to);
});

router.beforeEach(initializeAuthenticatedUserStore);

export default router;
