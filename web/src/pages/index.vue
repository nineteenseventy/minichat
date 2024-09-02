<script setup lang="ts">
import MessageComponent from '@/components/Message.component.vue';
import type { Message } from '@/interfaces/message.interface';
import type { User } from '@/interfaces/user.interface';
import { useUserStore } from '@/stores/user.store';

const me: User = {
  id: '1234',
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
</script>

<template>
  <MessageComponent :message="message" />
</template>
