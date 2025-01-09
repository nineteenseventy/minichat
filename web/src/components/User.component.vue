<script setup lang="ts">
import { useUserProfileDialog } from '@/composables/useUserProfileDialog';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatus.component.vue';
import { useUserStore } from '@/stores/user.store';
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
  <div class="user" @click="openProfile(_userId)">
    <UserPictureOnlineStatusComponent
      :picture="user?.picture"
      :user-id="_userId"
    />
    <span class="username">
      {{ user?.username }}
    </span>
  </div>
</template>

<style lang="scss" scoped>
.user {
  display: flex;
  flex-direction: row;
  align-items: center;
  cursor: pointer;
  .username {
    font-weight: bold;
    margin-left: 1rem;
    font-size: 1rem;
  }
  &:hover {
    .username {
      text-decoration: underline;
    }
  }
}
</style>
