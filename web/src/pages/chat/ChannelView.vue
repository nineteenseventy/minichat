<script setup lang="ts">
import Card from 'primevue/card';
import { useRoute } from 'vue-router';
import ChannelTitleComponent from '@/components/ChannelTitleComponent.vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import type { Message } from '@/interfaces/message.interface';
import MessageComponent from '@/components/MessageComponent.vue';
import type { Channel } from '@/interfaces/channel.interface';
import { useApi } from '@/composables/useApi';
import ChatInputComponent from '@/components/ChatInputComponent.vue';
import { useMessageStore } from '@/stores/messageStore';

const route = useRoute();
const channelId = route.params.id as string;
const { data: channel } = useApi(`/channels/${channelId}`).json<Channel>();

const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const messageStore = useMessageStore();
const messages = messageStore.getMessages(channelId);
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
      <ChatInputComponent :channelId="channelId" />
    </template>
  </Card>
</template>
