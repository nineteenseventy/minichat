<script setup lang="ts">
import MessageComponent from '@/components/Message.component.vue';
import type { Message } from '@/interfaces/message.interface';
import type { User } from '@/interfaces/user.interface';
import { useUserStore } from '@/stores/user.store';

const me: User = {
  id: '28585f72-0eb6-49ab-97ec-b7453a0f526b',
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
