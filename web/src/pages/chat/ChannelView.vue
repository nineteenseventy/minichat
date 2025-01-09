<script setup lang="ts">
import Card from 'primevue/card';
import { useRoute } from 'vue-router';
import ChannelTitleComponent from '@/components/ChannelTitleComponent.vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import type { Message } from '@/interfaces/message.interface';
import MessageComponent from '@/components/MessageComponent.vue';
import type { Channel } from '@/interfaces/channel.interface';
import { useApi } from '@/composables/useApi';
import Textarea from 'primevue/textarea';

const route = useRoute();
const channelId = route.params.id;
const { data: channel } = useApi(`/channels/${channelId}`).json<Channel>();

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
  <Card class="h-full">
    <template #title>
      <ChannelTitleComponent
        :channelTitle="channel?.title ?? 'loading channel name...'"
      />
    </template>
    <template #content>
      <MessageComponent
        v-for="message in messages"
        :key="message.id"
        :message="message"
      />
    </template>
    <template #footer>
      <Textarea
        autoResize
        rows="1"
        class="w-full"
        style="min-height: 42px"
      ></Textarea>
    </template>
  </Card>
</template>
