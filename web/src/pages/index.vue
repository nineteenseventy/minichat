<script setup lang="ts">
import MessageComponent from '@/components/Message.component.vue';
import type { Message } from '@/interfaces/message.interface';
import Card from 'primevue/card';
import { onBeforeMount } from 'vue';
import UserComponent from '@/components/User.component.vue';
import { useUserStore } from '@/stores/user.store';
import { useTimeoutPoll } from '@vueuse/core';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import ChannelsComponent from '@/components/channels.component.vue';

onBeforeMount(() => {
  useTimeoutPoll(async () => await userStore.updateStore(), 60000, {
    immediate: true,
  });
});

const userStore = useUserStore();
const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const message: Message = {
  authorId: authenticatedUserId,
  content: 'Hello, World!',
  id: '1',
  timestamp: new Date().toISOString(),
  attachments: [],
  channelId: 'Global Channel',
  read: false,
};
</script>

<template>
  <div class="flex flex-row gap-2 p-2 h-full">
    <nav class="w-72 flex flex-col gap-1">
      <Card class="h-full">
        <template #content>
          <ChannelsComponent class="h-full" />
        </template>
      </Card>
      <Card>
        <template #content>
          <UserComponent :userId="authenticatedUserId" />
        </template>
      </Card>
    </nav>
    <main class="flex-1">
      <span class="text-red-500"> MAIN CONTENT HERE </span>
      <MessageComponent :message="message" />
    </main>
  </div>
</template>
