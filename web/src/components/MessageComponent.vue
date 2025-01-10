<script setup lang="ts">
import Panel from 'primevue/panel';
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
  throw new Error('Not implemented');
}
</script>

<template>
  <Panel>
    <template #header>
      <UserComponent v-if="message" :userId="message.authorId" class="w-full" />
      <div class="flex justify-end gap-2">
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
    </template>
    <span class="content">
      {{ message?.content }}
    </span>
    <template #footer>
      <span class="text-xs">
        {{ timestamp }}
      </span>
    </template>
  </Panel>
</template>
