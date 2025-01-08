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

export const useUserStore = defineStore('user', () => {
  const users = ref<StoredUser[]>([]);
  function getUser(userId: string): Ref<User | undefined> {
    let storedUser = users.value.find((v) => v.id === userId);
    if (!storedUser) {
      storedUser = {
        id: userId,
        referenceCounter: 0,
      };
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

  return { users, getUser, updateStore };
});
