<script setup lang="ts">
import type { User } from '@/interfaces/user.interface';
import { useUserProfile } from '@/composables/useUserProfile';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatus.component.vue';

defineProps<{
  user: User;
}>();

const openProfile = useUserProfile();
</script>

<template>
  <div class="user" @click="openProfile(user.id)">
    <UserPictureOnlineStatusComponent
      :picture="user.picture"
      :user-id="user.id"
    />
    <span class="username">
      {{ user.username }}
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
