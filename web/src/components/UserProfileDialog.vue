<script setup lang="ts">
import { useApi } from '@/composables/useApi';
import type { UserProfile } from '@/interfaces/userProfile.interface';
import { computed, inject } from 'vue';
import type { DynamicDialogInstance } from 'primevue/dynamicdialogoptions';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import Button from 'primevue/button';
import { useRouter } from 'vue-router';
import UserPictureOnlineStatusComponent from './UserPictureOnlineStatusComponent.vue';
import SpinnerComponent from './SpinnerComponent.vue';
import { useChannelStore } from '@/stores/channelStore';

interface DialogRef {
  value: DynamicDialogInstance;
}
const dialogRef = inject<DialogRef>('dialogRef');
const router = useRouter();
const api = useApi;

const channelStore = useChannelStore();
const authenticatedUserId = useAuthenticatedUserStore().id;

const user = (dialogRef?.value.data.user as string) ?? authenticatedUserId;
const { data, error, isFetching } = api(`/users/${user}/profile`, {
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

const messageUser = async () => {
  const channel = await channelStore.getDirectChannel(user);
  if (channel.id) {
    router.push(`/channels/${channel.id}`);
    close();
  } else {
    console.warn(
      'something went wrong while trying to find/create channel for this user',
    );
  }
};
</script>

<template>
  <div class="flex flex-col mt-2">
    <div
      class="absolute inset-0 flex justify-center items-center backdrop-blur bg-black bg-opacity-50"
      v-if="isFetching"
    >
      <SpinnerComponent />
    </div>
    <UserPictureOnlineStatusComponent :userId="user" class="h-10 w-10" />
    <span class="font-bold text-2xl mix-blend-difference">{{
      data?.username
    }}</span>
    <span v-if="!!error" class="pt-2 text-orange-500"
      >Profile Could not be retrieved!</span
    >
    <div
      class="mt-4 flex flex-col bg-black bg-opacity-30 rounded-lg px-2 py-4 min-h-20"
    >
      <span class="font-bold mb-2">Bio</span>
      <div class="flex flex-col">
        <span v-for="line in bio" :key="line">{{ line }}</span>
      </div>
    </div>
    <div class="flex gap-4 mt-4">
      <Button @click="close()">Close</Button>
      <Button v-if="isMe" @click="editMyProfile()">Edit</Button>
      <Button v-if="!isMe" @click="messageUser()">Message</Button>
    </div>
  </div>
</template>
