<script setup lang="ts">
import Card from 'primevue/card';
import { onBeforeMount, onBeforeUnmount } from 'vue';
import UserComponent from '@/components/UserComponent.vue';
import { useUserStore } from '@/stores/userStore';
import { useTimeoutPoll } from '@vueuse/core';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import ChannelsComponent from '@/components/ChannelsComponent.vue';
import { useOnlineStatusStore } from '@/stores/onlineStatusStore';
import { useChannelStore } from '@/stores/channelStore';
import { useRouteParam } from '@/composables/useRouteParam';

const pollIntervals = [
  useTimeoutPoll(async () => await userStore.updateStore(), 60000),
  useTimeoutPoll(async () => await onlineStatusStore.updateStore(), 10000),
  useTimeoutPoll(async () => await channelStore.updateStore(), 60000),
];

onBeforeMount(() => {
  pollIntervals.forEach((pollInterval) => pollInterval.resume());
});

onBeforeUnmount(() => {
  pollIntervals.forEach((pollInterval) => pollInterval.pause());
});

const userStore = useUserStore();
const onlineStatusStore = useOnlineStatusStore();
const channelStore = useChannelStore();
const authenticatedUserId = useAuthenticatedUserStore().id;

const nestedRouteIsActive = useRouteParam('channelId');
console.log('route param:', nestedRouteIsActive.value);
</script>

<template>
  <div class="flex flex-row gap-2 p-2 h-full">
    <nav class="w-72 flex flex-col gap-2">
      <ChannelsComponent class="h-full" />
      <Card class="user-card">
        <template #content>
          <UserComponent :userId="authenticatedUserId" class="h-10" />
        </template>
      </Card>
    </nav>
    <main class="flex-1">
      <div v-if="!nestedRouteIsActive" class="flex h-full">
        <Card class="justify-self-center m-auto">
          <template #content>
            <p class="text-center">Select a channel to display its messages.</p>
          </template>
        </Card>
      </div>
      <RouterView style="align-content: center" />
    </main>
  </div>
</template>

<style scoped>
.user-card :deep(.p-card-content) {
  overflow: visible;
}
</style>
