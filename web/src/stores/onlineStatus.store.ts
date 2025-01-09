import type { OnlineStatus } from '@/interfaces/onlineStatus.interface';
import { defineStore } from 'pinia';
import { computed, shallowRef, triggerRef } from 'vue';
import { useAuthenticatedUserStore } from './authenticatedUser.store';

type UserOnlineStatusMap = Record<string, OnlineStatus>;

export const useOnlineStatusStore = defineStore('onlineStatus', () => {
  const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

  const userOnlineStatusMap = shallowRef<UserOnlineStatusMap>({});

  function setUserOnlineStatus(userId: string, status: OnlineStatus) {
    if (userId === authenticatedUserId) {
      userId = 'me';
    }
    userOnlineStatusMap.value[userId] = status;
    triggerRef(userOnlineStatusMap);
  }

  function getUserOnlineStatus(userId: string) {
    if (userId === authenticatedUserId) {
      userId = 'me';
    }
    return computed(() => userOnlineStatusMap.value[userId] ?? 'offline');
  }

  function clearUserOnlineStatus(userId: string) {
    if (userId === authenticatedUserId) {
      userId = 'me';
      console.warn('Clearing own online status');
    }
    delete userOnlineStatusMap.value[userId];
    triggerRef(userOnlineStatusMap);
  }

  return {
    userOnlineStatusMap,
    setUserOnlineStatus,
    getUserOnlineStatus,
    clearUserOnlineStatus,
  };
});
