<script setup lang="ts">
import type { Message } from '@/interfaces/message.interface';
import Panel from 'primevue/panel';
import UserComponent from './User.component.vue';
import { useRelativeFormattedDate } from '@/composables/useFormattedDate';
import { computed } from 'vue';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';

const props = defineProps<{
  /**
   * The message content.
   */
  message: Message;
}>();

const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const timestamp = computed(() => {
  if (!props.message.timestamp) return '';
  const date = new Date(props.message.timestamp);
  return useRelativeFormattedDate(date);
});

const isMyMessage = computed(() => {
  return props.message.authorId === authenticatedUserId;
});

function editMessage() {
  console.log('Edit message');
}
function deleteMessage() {
  console.log('Delete message');
}
</script>

<template>
  <Panel>
    <template #header>
      <UserComponent :userId="message.authorId" class="w-full" />
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
      {{ message.content }}
    </span>
    <template #footer>
      <span class="text-xs">
        {{ timestamp }}
      </span>
    </template>
  </Panel>
</template>
