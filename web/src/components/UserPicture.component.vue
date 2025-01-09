<script setup lang="ts">
import defaultUser from '@/assets/images/default-user.png';
import { useUserStore } from '@/stores/user.store';
import { onBeforeUnmount } from 'vue';

const props = defineProps<{
  userId: string;
}>();

onBeforeUnmount(() => userStore.unsubscribeUser(_userId));

const userStore = useUserStore();
const _userId = props.userId;
const user = userStore.getUser(_userId);
</script>

<template>
  <div class="picture">
    <img :src="user?.picture ?? defaultUser" alt="Profile picture" />
  </div>
</template>

<style scoped>
.picture {
  display: flex;
  border-radius: 50%;
  overflow: hidden;
  width: 3rem;
  height: 3rem;
  img {
    width: 3rem;
    height: auto;
  }
}
</style>
