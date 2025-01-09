<script setup lang="ts">
import { useOnlineStatusStore } from '@/stores/onlineStatus.store';
import UserPictureComponent from './UserPicture.component.vue';

const props = defineProps<{
  userId: string;
}>();

const onlineStatusStore = useOnlineStatusStore();
const online = onlineStatusStore.getUserOnlineStatus(props.userId);
</script>

<template>
  <UserPictureComponent
    :userId="userId"
    class="outline outline-2 outline-offset-2"
    :class="online"
  />
</template>

<style scoped lang="scss">
.online {
  outline-color: green;
}
.offline {
  outline-color: red;
}
.away {
  outline-color: yellow;
}
</style>
