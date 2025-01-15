<script setup lang="ts">
import { computed, ref, watch, watchEffect, type Ref } from 'vue';
import MessageComponent from '@/components/MessageComponent.vue';
import { useMessageStore } from '@/stores/messageStore';
import { useChannelStore } from '@/stores/channelStore';
import Card from 'primevue/card';
import { useTimeout, useTimeoutFn, useTimeoutPoll } from '@vueuse/core';

const messageStore = useMessageStore();
const channelStore = useChannelStore();

const props = defineProps<{
  channelId: string;
}>();

const messages = messageStore.getMinimalMessages(
  computed(() => props.channelId),
);

let startDate: Date;
let endDate: Date;
const reversedScrollContainer: Ref<HTMLElement> =
  ref<HTMLElement>() as Ref<HTMLElement>;

const poll = useTimeoutPoll(async () => {
  if (reversedScrollContainer.value.scrollTop >= -50) {
    await messageStore.loadMessages(props.channelId, undefined, endDate, 50);
  } else {
    await messageStore.loadMessages(props.channelId, startDate, undefined, 10);
  }
}, 1000);

function updateDates() {
  if (messages.value.length) {
    startDate = messages.value[0]?.timestamp || new Date();
    endDate =
      messages.value[messages.value.length - 1]?.timestamp || new Date();
  }
}

watch(
  () => props.channelId,
  async (channelId, beforeChannelId) => {
    poll.pause();
    if (beforeChannelId)
      messageStore.clearMessages({ channelId: beforeChannelId });
    await messageStore.loadMessages(channelId);
    await channelStore.setRead(channelId);
    updateDates();
    poll.resume();
  },
  { immediate: true },
);

let loading = false;
async function loadMessages(type: 'before' | 'after', date: Date) {
  if (loading) return;
  console.log('loading', type, date);
  loading = true;
  if (type === 'before')
    await messageStore.loadMessages(props.channelId, date, undefined, 10);
  else if (type === 'after')
    await messageStore.loadMessages(props.channelId, undefined, date, 10);
  updateDates();
  loading = false;
}

const startLoadingPixels = 50;

let scrollTimeout: ReturnType<typeof useTimeoutFn> | null = null;
function onScroll() {
  poll.pause();
  if (scrollTimeout) scrollTimeout.stop();
  scrollTimeout = useTimeoutFn(poll.resume, 1000);

  if (reversedScrollContainer.value.scrollTop <= startLoadingPixels) {
    loadMessages('before', startDate);
  } else if (
    reversedScrollContainer.value.scrollTop >=
    reversedScrollContainer.value.scrollHeight -
      reversedScrollContainer.value.clientHeight -
      startLoadingPixels
  ) {
    loadMessages('after', endDate);
  }
}
</script>

<template>
  <div
    ref="reversedScrollContainer"
    v-on:scroll="onScroll"
    class="flex flex-col-reverse overflow-y-auto gap-4"
  >
    <MessageComponent
      v-for="message in messages"
      :key="message.id"
      :messageId="message.id"
      class="px-2 py-2"
    />
    <Card v-if="!messages.length" class="justify-self-center m-auto">
      <template #content>
        <p class="text-center">
          This channel does not contain any messages yet.<br />
          Be the first one to say hi!
        </p>
      </template>
    </Card>
  </div>
</template>
