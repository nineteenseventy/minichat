<script setup lang="ts">
import MessageComponent from '@/components/Message.component.vue';
import type { Message } from '@/interfaces/message.interface';
import type { User } from '@/interfaces/user.interface';
import { useUserStore } from '@/stores/user.store';
import Panel from 'primevue/panel';
import Card from 'primevue/card';
import Listbox from 'primevue/listbox';
import { ref } from 'vue';
import UserComponent from '@/components/User.component.vue';

const me: User = {
  id: 'me',
  username: 'MeroFuruya',
  picture: 'https://avatars.githubusercontent.com/u/29742437?v=4',
};

const userStore = useUserStore();
userStore.setUser(me);

const message: Message = {
  author: me,
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
        <template #content><UserComponent :user="me" /></template>
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
#navbar {
  // margin: 8px;
}
</style>
