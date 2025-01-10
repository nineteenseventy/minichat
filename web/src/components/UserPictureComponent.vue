<script setup lang="ts">
import defaultUser from '@/assets/images/default-user.png';
import { useUserStore } from '@/stores/userStore';
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
  <div class="flex rounded-full overflow-hidden w-12 h-12">
    <img
      :src="user?.picture ?? defaultUser"
      alt="Profile picture"
      class="w-12 h-auto"
    />
  </div>
</template>
