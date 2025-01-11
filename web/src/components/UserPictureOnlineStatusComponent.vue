<script setup lang="ts">
import { useOnlineStatusStore } from '@/stores/onlineStatusStore';
import UserPictureComponent from './UserPictureComponent.vue';
import { onBeforeUnmount } from 'vue';

const props = defineProps<{
  userId: string;
}>();

onBeforeUnmount(() => {
  onlineStatusStore.unsubscribeOnlineStatus(_userId);
});

const onlineStatusStore = useOnlineStatusStore();
const _userId = props.userId;
const onlineStatus = onlineStatusStore.getOnlineStatus(_userId);
</script>

<template>
  <UserPictureComponent
    :userId="_userId"
    class="outline outline-2 outline-offset--1 p-1"
    :class="onlineStatus"
  />
</template>

<style scoped>
.outline-offset--1 {
  outline-offset: -2px;
}
.online {
  outline-color: green;
}
.offline {
  outline-color: gray;
}
.away {
  outline-color: yellow;
}
</style>
