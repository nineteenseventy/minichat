<script setup lang="ts">
import { useApi } from '@/composables/useApi';
import type { UserProfile } from '@/interfaces/userProfile.interface';
import { computed, inject } from 'vue';
import type { DynamicDialogInstance } from 'primevue/dynamicdialogoptions';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUser.store';
import Button from 'primevue/button';
import { useRouter } from 'vue-router';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatusComponent.vue';
import SpinnerComponent from './SpinnerComponent.vue';

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
    <UserPictureOnlineStatusComponent :userId="user" />
    <span class="font-bold text-2xl mix-blend-difference">{{
      data?.username
    }}</span>
    <span v-if="!!error" class="pt-2 text-orange-500"
      >Profile Could not be retrieved!</span
    >
    <div
      class="mt-4 flex flex-col bg-transparent bg-opacity-30 rounded-lg px-2 py-4 min-h-20"
    >
      <span class="font-bold mb-2">Bio</span>
      <div class="flex flex-col">
        <span v-for="line in bio" :key="line">{{ line }}</span>
      </div>
    </div>
    <div class="flex gap-4 mt-4">
      <Button @click="close()">Close</Button>
      <Button v-if="isMe" @click="editMyProfile()">Edit</Button>
    </div>
  </div>
</template>
