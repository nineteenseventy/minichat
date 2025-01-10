<script setup lang="ts">
import { computed, defineProps, effect, watch } from 'vue';
import MessageComponent from '@/components/MessageComponent.vue';
import { useMessageStore } from '@/stores/messageStore';

const messageStore = useMessageStore();

const props = defineProps<{
  channelId: string;
}>();

const messageIds = messageStore.getMessageIds(computed(() => props.channelId));

watch(
  () => props.channelId,
  async (channelId, beforeChannelId) => {
    if (beforeChannelId)
      messageStore.clearMessages({ channelId: beforeChannelId });
    await messageStore.loadMessages(channelId);
  },
  { immediate: true },
);

function onScroll(event: Event) {
  const target = event.target as HTMLElement;
  console.log('scrolling', target.scrollTop, target.scrollHeight);
}
</script>

<template>
  <div v-on:scroll="onScroll" class="flex flex-col-reverse">
    <MessageComponent
      v-for="messageId in messageIds"
      :key="messageId"
      :messageId="messageId"
    />
  </div>
</template>
