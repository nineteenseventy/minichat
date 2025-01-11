<script setup lang="ts">
import UserComponent from './UserComponent.vue';
import { useRelativeFormattedDate } from '@/composables/useFormattedDate';
import { computed, ref, watch } from 'vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { parseDate } from '@/utils/date/parseDate';
import { useMessageStore } from '@/stores/messageStore';
import ChatInputComponent from './ChatInputComponent.vue';
import type { NewMessage } from '@/interfaces/message.interface';

const messageStore = useMessageStore();

const props = defineProps<{
  messageId: string;
}>();

const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;
const message = messageStore.getMessage(computed(() => props.messageId));

const timestamp = computed(() => {
  if (!message.value?.timestamp) return '';
  const date = parseDate(message.value.timestamp);
  return useRelativeFormattedDate(date);
});

const isMyMessage = computed(() => {
  return message.value?.authorId === authenticatedUserId;
});

const mode = ref<'view' | 'edit'>('view');

function deleteMessage() {
  if (!message.value?.id) return;
  messageStore.deleteMessage(message.value.id);
}

const content = ref('');
watch(
  () => message.value?.content,
  (contentValue) => {
    content.value = contentValue ?? '';
  },
  { immediate: true },
);

async function onAfterEdit() {
  if (!message.value?.id) return;
  const updatedMessage: NewMessage = {
    content: content.value,
  };

  await messageStore.updateMessage(message.value.id, updatedMessage);

  mode.value = 'view';
}

const messageRows = computed(() => content.value.split('\n'));
</script>

<template>
  <div class="hover:bg-white hover:bg-opacity-5 rounded-content message">
    <div class="flex flex-row pb-1">
      <UserComponent
        v-if="message"
        :userId="message.authorId"
        class="w-full h-7"
      />
      <div class="flex gap-2 controls invisible">
        <span
          v-if="isMyMessage"
          class="cursor-pointer hover:underline pi pi-pencil"
          @click="mode = 'edit'"
        >
        </span>
        <span
          v-if="isMyMessage"
          class="cursor-pointer hover:underline pi pi-trash"
          @click="deleteMessage()"
        >
        </span>
      </div>
    </div>
    <span class="flex flex-col gap-1" v-if="mode === 'view'">
      <span
        v-for="(message, index) in messageRows"
        :key="index"
        class="break-words min-h-4"
      >
        {{ message }}
      </span>
    </span>
    <ChatInputComponent
      v-if="mode === 'edit'"
      v-model="content"
      @onSave="onAfterEdit()"
    />
    <span class="text-xs/3">
      {{ timestamp }}
    </span>
  </div>
</template>

<style scoped lang="scss">
.message:hover .controls {
  visibility: inherit;
}
</style>
