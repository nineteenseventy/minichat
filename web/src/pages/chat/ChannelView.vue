<script setup lang="ts">
import Card from 'primevue/card';
import { useRoute } from 'vue-router';
import ChannelTitleComponent from '@/components/ChannelTitleComponent.vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import type { Message } from '@/interfaces/message.interface';
import MessageComponent from '@/components/Message.component.vue';
import type { Channel } from '@/interfaces/channel.interface';
import { useApi } from '@/composables/useApi';
import Textarea from 'primevue/textarea';

const route = useRoute();
const channelId = route.params.id;
const { data: channel } = useApi(`/channels/${channelId}`).json<Channel>();

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
  <Card class="h-full">
    <template #title>
      <ChannelTitleComponent
        :channelTitle="channel?.title ?? 'loading channel name...'"
      />
    </template>
    <template #content>
      <MessageComponent :message="message" />
    </template>
    <template #footer>
      <Textarea autoResize rows="1" class="w-full"></Textarea>
    </template>
  </Card>
</template>
