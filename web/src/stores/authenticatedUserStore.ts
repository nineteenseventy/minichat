import { useApi } from '@/composables/useApi';
import type { UserProfile } from '@/interfaces/userProfile.interface';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

let userStoreInitialized: Promise<boolean> = new Promise((resolve) =>
  resolve(false),
);

let userProfileInit: UserProfile;

export async function initializeAuthenticatedUserStore() {
  if (await userStoreInitialized) return;

  let _resolve: () => void;
  userStoreInitialized = new Promise(
    (resolve) => (_resolve = resolve.bind(null, true)),
  );

  const { data } = await useApi('/users/me/profile').json<UserProfile>();
  if (!data.value)
    throw new Error('Could not retrieve data about the logged in user!');

  data.value.picture ??= '/src/assets/images/default-user.png';
  userProfileInit = data.value;

  _resolve!();
  console.debug(
    'authenticated user store initialized with: ' + JSON.stringify(data.value),
  );
}

export const useAuthenticatedUserStore = defineStore(
  'authenticatedUser',
  () => {
    const id = ref(userProfileInit.id);
    const profile = ref<UserProfile>(userProfileInit);

    function getProfile() {
      return computed(() => profile.value);
    }

    return { id, profile, getProfile };
  },
);
