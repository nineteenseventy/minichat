import { ref, computed, type Ref } from 'vue';
import { defineStore } from 'pinia';
import { useApi } from '@/composables/useApi';
import type { User } from '@/interfaces/user.interface';

interface StoredUser {
  id: string;
  referenceCounter: number;
  user?: User;
}

export const useUserStore = defineStore('user', () => {
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
    const activeUsers = users.value.filter((v) => v.referenceCounter > 0);
    users.value = activeUsers;
    for (let i = 0; i < activeUsers.length; i++) {
      users.value[i].user = await fetchUser(users.value[i].id);
    }
  }

  function unsubscribeUser(userId: string) {
    const storedUser = computed(() => users.value.find((v) => v.id === userId));
    if (!storedUser.value) return;
    if (storedUser.value.referenceCounter === 0) {
      console.warn('unsubscribeUser was called to often for user', userId);
      return;
    }
    storedUser.value.referenceCounter--;
  }

  async function updateUser(userId: string, user: Partial<User>) {
    const storedUserIndex = users.value.findIndex((v) => v.id === userId);
    if (storedUserIndex === -1) return;
    if (!users.value[storedUserIndex].user) {
      users.value[storedUserIndex].user = {
        ...(user as User),
        id: userId,
      };
      return;
    }
    users.value[storedUserIndex].user = {
      ...users.value[storedUserIndex].user,
      ...user,
      id: userId,
    };
  }

  async function fetchUser(userId: string) {
    const { data } = await useApi(`/users/${userId}`).json<User>();
    return data.value ?? undefined;
  }

  return { users, getUser, updateStore, unsubscribeUser, updateUser };
});
