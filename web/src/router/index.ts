import { createAuthGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';
import { initializeAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { globalAuth0 } from '@/plugins/auth0';
import CallbackErrorView from '@/views/CallbackErrorView.vue';
import MainView from '@/views/MainView.vue';
import ChatView from '@/views/ChannelView.vue';
import UserSettingsComponent from '@/components/UserSettingsComponent.vue';
import { globalEnv } from '@/plugins/assetEnvPlugin';

const router = createRouter({
  history: createWebHistory(globalEnv.BASE_URL),
  routes: [
    {
      path: '/callback',
      name: 'callback',
      component: CallbackErrorView,
      beforeEnter: async (to) => {
        let appState;
        try {
          appState = await globalAuth0.handleRedirectCallback(to.fullPath);
        } catch (error) {
          console.error('Error handling redirect callback:', error);
          return;
        }
        return { path: appState.appState?.target ?? '/' };
      },
    },
    {
      path: '/settings/profile',
      name: 'profileSettings',
      component: UserSettingsComponent,
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

router.beforeEach((to) => {
  if (to.name === 'callback') return;
  return initializeAuthenticatedUserStore();
});

export default router;
