<script setup lang="ts">
import MessageComponent from '@/components/Message.component.vue';
import type { Message } from '@/interfaces/message.interface';
import Panel from 'primevue/panel';
import Card from 'primevue/card';
import Listbox from 'primevue/listbox';
import { ref } from 'vue';
import UserComponent from '@/components/User.component.vue';
import { useUserStore } from '@/stores/user.store';

const userStore = useUserStore();

const message: Message = {
  authorId: 'me',
  content: 'Hello, World!',
  id: '1',
  timestamp: new Date().toISOString(),
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
          ><UserComponent :userId="userStore.authenticatedUserId"
        /></template>
      </Card>
    </nav>
    <main>
      MAIN CONTENT HERE
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
