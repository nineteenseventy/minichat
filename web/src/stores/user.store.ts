import { ref, computed, type Ref } from 'vue';
import { defineStore } from 'pinia';
import { useApi } from '@/composables/useApi';
import type { User } from '@/interfaces/user.interface';

interface StoredUser {
  id: string;
  referenceCounter: number;
  user?: User;
}

const userStoreInitialized = ref(false);
const initializationStarted = ref(false);

export async function initializeUserStore() {
  const store = useUserStore();
  if (!userStoreInitialized.value && !initializationStarted.value) {
    initializationStarted.value = true;
    console.debug('initialize() called');
    const { data } = await useApi('/users/me').json<User>();
    if (!data.value) throw new Error('Something went wrong!');
    store.authenticatedUserId = data.value.id;
    userStoreInitialized.value = true;
    console.info('user store initialized with: ' + JSON.stringify(data.value));
  }
}

export const useUserStore = defineStore('user', () => {
  const authenticatedUserId = ref<string>('');
  const users = ref<StoredUser[]>([]);
  function getUser(userId: string): Ref<User | undefined> {
    const storedUser = computed(() => users.value.find((v) => v.id === userId));
    if (!storedUser.value) {
      const newStoredUser = {
        id: userId,
        referenceCounter: 0,
      };
      users.value.push(newStoredUser);
      fetchUser(userId).then(
        (fetchedUser) => (storedUser.value!.user = fetchedUser),
      );
    }
    storedUser.value!.referenceCounter++;
    return computed(() => users.value.find((v) => v.id === userId)?.user);
  }

  async function updateStore() {
    console.log('updateStore() called');
    const activeUsers = users.value.filter((v) => v.referenceCounter > 0);
    users.value = activeUsers;
    for (let i = 0; i < activeUsers.length; i++) {
      users.value[i].user = await fetchUser(users.value[i].id);
    }
  }

  function unsubscribeUser(userId: string) {
    const storedUser = computed(() => users.value.find((v) => v.id === userId));
    if (!storedUser.value) return;
    storedUser.value.referenceCounter--;
  }

  async function fetchUser(userId: string) {
    const { data } = await useApi(`/users/${userId}`).json<User>();
    return data.value ?? undefined;
  }

  return { authenticatedUserId, users, getUser, updateStore, unsubscribeUser };
});
