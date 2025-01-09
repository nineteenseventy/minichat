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
import { useOnlineStatusStore } from '@/stores/onlineStatus.store';

onBeforeMount(() => {
  useTimeoutPoll(async () => await userStore.updateStore(), 60000, {
    immediate: true,
  });

  useTimeoutPoll(async () => await onlineStatus.updateStore(), 10000, {
    immediate: true,
  });
});

const userStore = useUserStore();
const onlineStatus = useOnlineStatusStore();
const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const messages: Message[] = [
  {
    authorId: authenticatedUserId,
    content: 'Hello, World!',
    id: '1',
    timestamp: new Date().toISOString(),
    attachments: [],
    channelId: 'Global Channel',
    read: false,
  },
  {
    authorId: '038678a5-b6f1-45dc-b1da-9f3837f4cdc8',
    content: 'Hello, World!',
    id: '2',
    timestamp: new Date().toISOString(),
    attachments: [],
    channelId: 'Global Channel',
    read: false,
  },
  {
    authorId: '2d0a4682-a7a3-4461-98bc-4b403a94f000',
    content: 'Hello, World!',
    id: '3',
    timestamp: new Date().toISOString(),
    attachments: [],
    channelId: 'Global Channel',
    read: false,
  },
];
</script>

<template>
  <div class="flex flex-row gap-2 p-2 h-full">
    <nav class="w-72 flex flex-col gap-1">
      <ChannelsComponent class="h-full" />
      <!-- <Card class="h-full">
        <template #content>
        </template>
      </Card> -->
      <Card>
        <template #content>
          <UserComponent :userId="authenticatedUserId" />
        </template>
      </Card>
    </nav>
    <main class="flex-1">
      <span class="text-red-500"> MAIN CONTENT HERE </span>
      <MessageComponent
        v-for="message in messages"
        :key="message.id"
        :message="message"
      />
    </main>
  </div>
</template>
