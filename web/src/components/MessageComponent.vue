<script setup lang="ts">
import UserComponent from './UserComponent.vue';
import { useRelativeFormattedDate } from '@/composables/useFormattedDate';
import { computed } from 'vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { parseDate } from '@/utils/date/parseDate';
import { useMessageStore } from '@/stores/messageStore';

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

function editMessage() {
  throw new Error('Not implemented');
}
function deleteMessage() {
  if (!message.value?.id) return;
  messageStore.deleteMessage(message.value.id);
}

const messageContent = computed(() => message.value?.content.split('\n'));
</script>

<template>
  <div class="px-4 py-2 hover:bg-white hover:bg-opacity-5 rounded-content">
    <div class="flex flex-row pb-1">
      <UserComponent
        v-if="message"
        :userId="message.authorId"
        class="w-full h-7"
      />
      <div class="flex gap-2">
        <span
          v-if="isMyMessage"
          class="cursor-pointer hover:underline"
          @click="editMessage()"
        >
          Edit
        </span>
        <span
          v-if="isMyMessage"
          class="cursor-pointer hover:underline"
          @click="deleteMessage()"
        >
          Delete
        </span>
      </div>
    </div>
    <span class="flex flex-col gap-1">
      <span
        v-for="(message, index) in messageContent"
        :key="index"
        class="break-words min-h-4"
      >
        {{ message }}
      </span>
    </span>
    <span class="text-xs/3">
      {{ timestamp }}
    </span>
  </div>
</template>
