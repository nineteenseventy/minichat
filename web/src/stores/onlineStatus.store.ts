import type {
  UserStatusResponse,
  UserStatus,
} from '@/interfaces/user.interface';
import { defineStore } from 'pinia';
import { computed, ref, shallowRef } from 'vue';
import { useApi } from '@/composables/useApi';

interface StoredUserOnlineStatus {
  id: string;
  referenceCounter: number;
  onlineStatus?: UserStatus;
}

export const useOnlineStatusStore = defineStore('onlineStatus', () => {
  const onlineStatuses = ref<StoredUserOnlineStatus[]>([]);

  function getOnlineStatus(userId: string) {
    const storedOnlineStatus = computed(() =>
      onlineStatuses.value.find((v) => v.id === userId),
    );
    if (!storedOnlineStatus.value) {
      const newStoredOnlineStatus = {
        id: userId,
        referenceCounter: 0,
      };
      onlineStatuses.value.push(newStoredOnlineStatus);

      fetchOnlineStatus(userId).then(
        (fetchedOnlineStatus) =>
          (storedOnlineStatus.value!.onlineStatus =
            fetchedOnlineStatus?.status ?? 'offline'),
      );
    }
    storedOnlineStatus.value!.referenceCounter++;
    return computed(
      () => onlineStatuses.value.find((v) => v.id === userId)?.onlineStatus,
    );
  }

  async function updateStore() {
    const activeOnlineStatuses = onlineStatuses.value.filter(
      (v) => v.referenceCounter > 0,
    );
    onlineStatuses.value = activeOnlineStatuses;

    const userIds = activeOnlineStatuses.map((v) => v.id);
    const newOnlineStatuses = await fetchOnlineStatusesAndEcho(userIds);
    for (let i = 0; i < onlineStatuses.value.length; i++) {
      const userId = onlineStatuses.value[i]?.id;
      const newOnlineStatus = newOnlineStatuses.find((v) => v.id === userId);

      if (newOnlineStatus) {
        onlineStatuses.value[i].onlineStatus = newOnlineStatus.status;
      }
    }
  }

  function unsubscribeOnlineStatus(userId: string) {
    const storedOnlineStatus = computed(() =>
      onlineStatuses.value.find((v) => v.id === userId),
    );
    if (!storedOnlineStatus.value) return;
    storedOnlineStatus.value.referenceCounter--;
  }

  async function fetchOnlineStatus(userId: string) {
    const { data } = await useApi(
      `/users/${userId}/status`,
    ).json<UserStatusResponse>();
    return data.value ?? undefined;
  }

  async function fetchOnlineStatusesAndEcho(userIds: string[]) {
    if (userIds.length === 0) {
      await echoOnly();
      return [];
    }

    const sliceSize = 100;
    const onlineStatuses: UserStatusResponse[] = [];

    for (let i = 0; i < userIds.length; i += sliceSize) {
      const slice = userIds.slice(i, i + sliceSize);

      const queryValue = slice.join(',');
      const { data } = await useApi(
        `/users/echoAndGetStatuses?ids=${queryValue}`,
      )
        .post()
        .json<UserStatusResponse[]>();

      if (data.value) onlineStatuses.push(...data.value);
    }

    return onlineStatuses;
  }

  async function echoOnly() {
    await useApi('/users/echo').post();
  }

  return {
    onlineStatuses,
    getOnlineStatus,
    updateStore,
    unsubscribeOnlineStatus,
  };
});
