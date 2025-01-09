<script setup lang="ts">
import { useApi } from '@/composables/useApi';
import type { UserProfile } from '@/interfaces/userProfile.interface';
import { computed, inject } from 'vue';
import type { DynamicDialogInstance } from 'primevue/dynamicdialogoptions';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import Button from 'primevue/button';
import { useRouter } from 'vue-router';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatus.component.vue';
import SpinnerComponent from './Spinner.component.vue';

interface DialogRef {
  value: DynamicDialogInstance;
}
const dialogRef = inject<DialogRef>('dialogRef');
const router = useRouter();

const authenticatedUserId = useAuthenticatedUserStore().authenticatedUserId;

const user = dialogRef?.value.data.user ?? authenticatedUserId;
const { data, error, isFetching } = useApi(`/users/${user}/profile`, {
  afterFetch(ctx) {
    if (dialogRef) {
      dialogRef.value.options.props!.style.backgroundColor = ctx.data.color;
    }
    return ctx;
  },
}).json<UserProfile>();

const isMe = computed(() => data?.value?.id === authenticatedUserId);
const close = () => dialogRef?.value.close();
const editMyProfile = () => {
  dialogRef?.value.close();
  router.push('/settings/profile');
};

const bio = computed(() => {
  return data?.value?.bio?.split('\n') ?? [];
});
</script>

<template>
  <div class="flex flex-col mt-2">
    <div
      class="absolute inset-0 flex justify-center items-center backdrop-blur bg-black bg-opacity-50"
      v-if="isFetching"
    >
      <SpinnerComponent />
    </div>
    <UserPictureOnlineStatusComponent :picture="data?.picture" :userId="user" />
    <span class="username">{{ data?.username }}</span>
    <span v-if="!!error" class="profile-error"
      >Profile Could not be retrieved!</span
    >
    <div class="bio">
      <span class="title">Bio</span>
      <div class="content">
        <span v-for="line in bio" :key="line">{{ line }}</span>
      </div>
    </div>
    <div class="controls">
      <Button @click="close()">Close</Button>
      <Button v-if="isMe" @click="editMyProfile()">Edit</Button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.5);
}

.user-profile {
  display: flex;
  flex-direction: column;
  padding: 1rem;
}

.username {
  font-weight: bold;
  font-size: 1.5rem;
  mix-blend-mode: difference;
}

.bio {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  background-color: rgba(0, 0, 0, 0.3);
  border-radius: 0.5rem;
  padding: 0.5rem 1rem;
  min-height: 5rem;
  font-size: 1rem;
  .title {
    font-weight: bold;
    margin-bottom: 0.5rem;
  }
  .content {
    display: flex;
    flex-direction: column;
  }
}

.controls {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.profile-error {
  padding-top: 0.5rem;
  color: orange;
}
</style>
