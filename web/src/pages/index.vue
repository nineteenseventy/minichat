<script setup lang="ts">
import MessageComponent from '@/components/Message.component.vue';
import type { Message } from '@/interfaces/message.interface';
import Panel from 'primevue/panel';
import Card from 'primevue/card';
import Listbox from 'primevue/listbox';
import { onBeforeMount, ref } from 'vue';
import UserComponent from '@/components/User.component.vue';
import { useUserStore } from '@/stores/user.store';
import { useTimeoutPoll } from '@vueuse/core';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';

onBeforeMount(() => {
  useTimeoutPoll(async () => await userStore.updateStore(), 60000, {
    immediate: true,
  });
});

const userStore = useUserStore();
const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const message: Message = {
  authorId: authenticatedUserId,
  content: 'Hello, World!',
  id: '1',
  timestamp: new Date().toISOString(),
  attachments: [],
  channelId: 'Global Channel',
  read: false,
};

const selectedElement = ref();
const publicChannels = ref(['Global Channel', 'BBS1', 'Germany']);
const privateAndDirectChannels = ref([
  'John Doe',
  'Vue Forum',
  'Black Ops 3',
  'Animals',
]);
</script>

<template>
  <div id="index">
    <nav id="navbar">
      <Panel header="Public Channels" toggleable>
        <p><Listbox :options="publicChannels"></Listbox></p>
      </Panel>
      <Panel header="Group and Private Chats" toggleable>
        <p><Listbox :options="privateAndDirectChannels"></Listbox></p>
      </Panel>
      <Card>
        <template #content
          ><UserComponent :userId="authenticatedUserId"
        /></template>
      </Card>
    </nav>
    <main>
      <span class="text-red-500"> MAIN CONTENT HERE </span>
      <MessageComponent :message="message" />
    </main>
  </div>
</template>
<style lang="scss" scoped>
#index {
  display: flex;
  flex-direction: row;
}

nav {
  width: 300px;
}

main {
  flex: 1;
}
</style>
