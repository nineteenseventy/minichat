<script setup lang="ts">
import { computed, watch } from 'vue';
import MessageComponent from '@/components/MessageComponent.vue';
import { useMessageStore } from '@/stores/messageStore';
import { useChannelStore } from '@/stores/channelStore';

const messageStore = useMessageStore();
const channelStore = useChannelStore();

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
    await channelStore.setRead(channelId);
  },
  { immediate: true },
);

function onScroll(event: Event) {
  const target = event.target as HTMLElement;
  console.log(
    'scrolling',
    target.scrollTop,
    target.scrollHeight,
    target.clientHeight,
    target.scrollHeight - target.scrollTop - target.clientHeight,
  );
}
</script>

<template>
  <div
    v-on:scroll="onScroll"
    class="flex flex-col-reverse overflow-y-auto gap-4"
  >
    <MessageComponent
      v-for="messageId in messageIds"
      :key="messageId"
      :messageId="messageId"
      class="px-2 py-2"
    />

    <div v-if="!messageIds.length" class="md:justify-self-center md:m-auto">
      <p class="text-muted-color md:text-center">
        This channel does not contain any messages yet.
        <br class="hidden md:block" />
        Be the first one to say hi!
      </p>
    </div>

    <!-- <Card v-if="!messageIds.length" class="justify-self-center m-auto">
      <template #content>
      </template>
    </Card> -->
  </div>
</template>
