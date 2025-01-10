<script setup lang="ts">
import Card from 'primevue/card';
import ChannelTitleComponent from '@/components/ChannelTitleComponent.vue';
import type { Message } from '@/interfaces/message.interface';
import MessageComponent from '@/components/MessageComponent.vue';
import ChatInputComponent from '@/components/ChatInputComponent.vue';
import { useMessageStore } from '@/stores/messageStore';
import { effect } from 'vue';
import { useRouteParam } from '@/composables/useRouteParam';

const messageStore = useMessageStore();

const channelId = useRouteParam('channelId');

let messages: Message[] = [];

effect(() => {
  if (channelId.value) {
    messages = messageStore.getMessages(channelId.value);
  }
});
</script>

<template>
  <Card class="h-full">
    <template #title>
      <ChannelTitleComponent :channelId="channelId!" />
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
