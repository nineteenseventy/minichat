<script setup lang="ts">
import Card from 'primevue/card';
import { onBeforeMount } from 'vue';
import UserComponent from '@/components/UserComponent.vue';
import { useUserStore } from '@/stores/userStore';
import { useTimeoutPoll } from '@vueuse/core';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import ChannelsComponent from '@/components/ChannelsComponent.vue';
import { useOnlineStatusStore } from '@/stores/onlineStatusStore';
import { useChannelStore } from '@/stores/channelStore';

onBeforeMount(() => {
  useTimeoutPoll(async () => await userStore.updateStore(), 60000, {
    immediate: true,
  });

  useTimeoutPoll(async () => await onlineStatusStore.updateStore(), 10000, {
    immediate: true,
  });

  useTimeoutPoll(async () => await channelStore.updateStore(), 60000, {
    immediate: true,
  });
});

const userStore = useUserStore();
const onlineStatusStore = useOnlineStatusStore();
const channelStore = useChannelStore();
const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;
</script>

<template>
  <div class="flex flex-row gap-2 p-2 h-full">
    <nav class="w-72 flex flex-col gap-2">
      <ChannelsComponent class="h-full" />
      <Card>
        <template #content>
          <UserComponent :userId="authenticatedUserId" />
        </template>
      </Card>
    </nav>
    <main class="flex-1">
      <RouterView />
    </main>
  </div>
</template>
