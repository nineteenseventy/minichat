import { useApi } from '@/composables/useApi';
import type { User } from '@auth0/auth0-vue';
import { defineStore } from 'pinia';
import { ref } from 'vue';

let userStoreInitialized: Promise<boolean> = new Promise((resolve) =>
  resolve(false),
);

export async function initializeAuthenticatedUserStore() {
  if (await userStoreInitialized) return;

  let _resolve: () => void;
  userStoreInitialized = new Promise(
    (resolve) => (_resolve = resolve.bind(null, true)),
  );

  const store = useAuthenticatedUserStore();
  const { data } = await useApi('/users/me').json<User>();
  if (!data.value) throw new Error('could');
  store.authenticatedUserId = data.value.id;

  _resolve!();
  console.info('user store initialized with: ' + JSON.stringify(data.value));
}

export const useAuthenticatedUserStore = defineStore(
  'authenticatedUser',
  () => {
    const authenticatedUserId = ref<string>('undefined');

    return { authenticatedUserId };
  },
);
