<script setup lang="ts">
import { useUserProfileDialog } from '@/composables/useUserProfileDialog';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatusComponent.vue';
import { useUserStore } from '@/stores/userStore';
import { onBeforeUnmount } from 'vue';

const props = defineProps<{
  userId: string;
}>();

onBeforeUnmount(() => userStore.unsubscribeUser(_userId));

const openProfile = useUserProfileDialog();
const userStore = useUserStore();
const _userId = props.userId;
const user = userStore.getUser(_userId);
</script>

<template>
  <div class="flex flex-row items-center">
    <UserPictureOnlineStatusComponent :userId="_userId" class="h-full" />
    <span
      @click="openProfile(_userId)"
      class="font-bold ml-2 cursor-pointer hover:underline"
    >
      {{ user?.username }}
    </span>
  </div>
</template>
