<script setup lang="ts">
import { useUserProfileDialog } from '@/composables/useUserProfileDialog';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatus.component.vue';
import { useUserStore } from '@/stores/user.store';

const props = defineProps<{
  userId: string;
}>();

const openProfile = useUserProfileDialog();
const user = useUserStore().getUser(props.userId);
</script>

<template>
  <div class="user" @click="openProfile(userId)">
    <UserPictureOnlineStatusComponent
      :picture="user?.picture"
      :user-id="userId"
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
