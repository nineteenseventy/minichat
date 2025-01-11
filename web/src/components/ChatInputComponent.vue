<script setup lang="ts">
import { ref } from 'vue';
import { useMessageDraftsStore } from '@/stores/messageDraftsStore';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';
import { useMessageStore } from '@/stores/messageStore';
import type { NewMessage } from '@/interfaces/message.interface';

const messageStore = useMessageStore();
const draftsStore = useMessageDraftsStore();

const props = defineProps<{ channelId: string }>();

const content = ref('');

if (props.channelId) {
  const draft = draftsStore.getMessageDraft(props.channelId);
  if (draft) content.value = draft;
}

async function onSend() {
  const contentValue = content.value.trim();
  if (contentValue.length < 1) return;

  const newMessage: NewMessage = {
    content: contentValue,
  };

  await messageStore.sendMessage(props.channelId, newMessage);
  content.value = '';
  draftsStore.clearMessageDraft(props.channelId);
}

async function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    await onSend();
  }
}
</script>

<template>
  <div class="flex flex-row gap-3 items-center">
    <Textarea
      autoResize
      fluid
      rows="1"
      class="min-h-11 max-h-36"
      v-model="content"
      @keydown="onKeydown($event)"
    />
    <Button icon="pi pi-send" class="self-end" @click="onSend()" />
  </div>
</template>

<style scoped lang="scss"></style>
