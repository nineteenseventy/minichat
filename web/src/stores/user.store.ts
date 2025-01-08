import { ref, computed, type Ref, reactive } from 'vue';
import { defineStore } from 'pinia';
import { useApi } from '@/composables/useApi';
import type { UserProfile } from '@/interfaces/userProfile.interface'; // TODO: create seperate interface for stored user
import type { User } from '@/interfaces/user.interface';

interface StoredUser {
  id: string;
  referenceCounter: number;
  user?: User;
}

const userStoreInitialized = ref(false);
const initializationStarted = ref(false);

async function initialize() {
  return (await useApi('/users/me').json<User>()).data.value?.id;
}

export const useUserStore = () => {
  const store = userStore();
  if (!userStoreInitialized.value && !initializationStarted.value) {
    initializationStarted.value = true;
    console.debug('initialize() called');
    void initialize().then((data) => {
      if (!data) throw new Error('Something went wrong!');
      store.authenticatedUserId = data;
      userStoreInitialized.value = true;
      console.info('user store initialized with: ' + JSON.stringify(data));
    });
  }
  return store;
};

const userStore = defineStore('user', () => {
  const authenticatedUserId = ref<string>(
    '2d0a4682-a7a3-4461-98bc-4b403a94f000',
  );
  const users = ref<StoredUser[]>([]);
  function getUser(userId: string): Ref<User | undefined> {
    let storedUser = users.value.find((v) => v.id === userId);
    if (!storedUser) {
      storedUser = {
        id: userId,
        referenceCounter: 0,
      };
      updateUser(storedUser);
      users.value.push(storedUser);
    }
    storedUser.referenceCounter += 1;
    return computed(() => users.value.find((v) => v.id === userId)?.user);
  }

  async function updateStore() {
    const activeUsers = users.value.filter((v) => v.referenceCounter > 0);
    users.value = activeUsers;
    for (const storedUser of activeUsers) {
      await updateUser(storedUser);
    }
  }

  async function updateUser(storedUser: StoredUser) {
    const { data } = await useApi(`/users/${storedUser.id}`).json<User>();
    if (data.value) {
      storedUser.user = data.value;
    }
  }

  return { authenticatedUserId, users, getUser, updateStore };
});
