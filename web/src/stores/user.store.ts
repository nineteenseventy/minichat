import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import type { UserProfile } from '@/interfaces/userProfile.interface';

export const useUserStore = defineStore('user', () => {
  const user = ref<UserProfile | undefined>(undefined);
  function setUser(newUser: UserProfile) {
    user.value = newUser;
  }

  const id = computed(() => user.value?.id ?? undefined);
  const username = computed(() => user.value?.username ?? undefined);
  const bio = computed(() => user.value?.bio ?? undefined);
  const profilePicture = computed(() => user.value?.picture ?? undefined);

  return { user, setUser, id, username, bio, profilePicture };
});
