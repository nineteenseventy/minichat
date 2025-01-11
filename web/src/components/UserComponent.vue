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
  <div
    class="flex flex-row items-center cursor-pointer hover:underline"
    @click="openProfile(_userId)"
  >
    <UserPictureOnlineStatusComponent :userId="_userId" class="h-full" />
    <span class="font-bold ml-2 hover:underline">
      {{ user?.username }}
    </span>
  </div>
</template>
