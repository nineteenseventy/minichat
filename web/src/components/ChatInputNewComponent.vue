<script setup lang="ts">
import { ref, watch } from 'vue';
import { useMessageDraftsStore } from '@/stores/messageDraftsStore';
import { useMessageStore } from '@/stores/messageStore';
import type { NewMessage } from '@/interfaces/message.interface';
import ChatInputComponent from './ChatInputComponent.vue';

const messageStore = useMessageStore();
const draftsStore = useMessageDraftsStore();

const props = defineProps<{ channelId: string }>();

const content = ref('');

watch(
  () => props.channelId,
  (channelId, beforeChannelId) => {
    if (beforeChannelId) {
      draftsStore.setMessageDraft(beforeChannelId, content.value);
    }
    content.value = draftsStore.getMessageDraft(channelId) ?? '';
  },
  { immediate: true },
);

if (props.channelId) {
  const draft = draftsStore.getMessageDraft(props.channelId);
  if (draft) content.value = draft;
}

async function onSend() {
  const contentValue = content.value.trim();
  if (contentValue.length < 1) return;

  content.value = '';
  draftsStore.clearMessageDraft(props.channelId);

  const newMessage: NewMessage = {
    content: contentValue,
  };
  await messageStore.sendMessage(props.channelId, newMessage);
}
</script>

<template>
  <ChatInputComponent v-model="content" @onSave="onSend()" />
</template>

<style scoped lang="scss"></style>
