<script setup lang="ts">
import UserComponent from './UserComponent.vue';
import { useRelativeFormattedDate } from '@/composables/useFormattedDate';
import { computed, ref, watch } from 'vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { parseDate } from '@/utils/date/parseDate';
import { useMessageStore } from '@/stores/messageStore';
import ChatInputComponent from './ChatInputComponent.vue';
import type { NewMessage } from '@/interfaces/message.interface';
import { useMessageRenderer } from '@/composables/useMessageRenderer';
import { useConfirm } from 'primevue/useconfirm';

const props = defineProps<{
  messageId: string;
}>();

const messageRenderer = useMessageRenderer();

const messageStore = useMessageStore();
const message = messageStore.getMessage(computed(() => props.messageId));

const timestamp = computed(() => {
  if (!message.value?.timestamp) return '';
  const date = parseDate(message.value.timestamp);
  return useRelativeFormattedDate(date);
});

const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const isMyMessage = computed(() => {
  return message.value?.authorId === authenticatedUserId;
});

const mode = ref<'view' | 'edit'>('view');

function deleteMessage() {
  if (!message.value?.id) return;
  messageStore.deleteMessage(message.value.id);
}

const confirm = useConfirm();
const confirmDeleteMessage = (event: Event) => {
  // (event.target as HTMLSpanElement).focus();
  confirm.require({
    target: event.target as HTMLSpanElement,
    message: 'Are your sure you want to delete this message?',
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: 'Cancel',
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: 'Delete',
      severity: 'danger',
    },
    accept: deleteMessage,
  });
};

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
</script>

<template>
  <div class="hover:bg-highlight rounded-content message">
    <div class="flex flex-row pb-1">
      <UserComponent
        v-if="message"
        :userId="message.authorId"
        class="w-full h-7"
      />
      <div class="flex gap-2 controls invisible" v-if="mode === 'view'">
        <span
          v-if="isMyMessage"
          class="cursor-pointer hover:underline pi pi-pencil"
          @click="mode = 'edit'"
        >
        </span>
        <span
          v-if="isMyMessage"
          class="cursor-pointer hover:underline pi pi-trash"
          @click="confirmDeleteMessage"
        >
        </span>
      </div>
    </div>
    <span class="flex flex-col gap-1" v-if="mode === 'view'">
      <span v-html="messageRenderer(message?.content)" />
    </span>
    <ChatInputComponent
      v-if="mode === 'edit'"
      v-model="content"
      :enable-cancel="true"
      @onSave="onAfterEdit()"
      @onCancel="mode = 'view'"
    />
    <span class="text-xs/3" :v-tooltip="'a'">
      {{ timestamp }}
    </span>
  </div>
</template>

<style scoped lang="scss">
.message:hover .controls {
  visibility: inherit;
}
</style>
