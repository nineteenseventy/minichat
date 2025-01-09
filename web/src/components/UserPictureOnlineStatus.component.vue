<script setup lang="ts">
import { useOnlineStatusStore } from '@/stores/onlineStatus.store';
import UserPictureComponent from './UserPicture.component.vue';
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
    class="outline outline-2 outline-offset-2"
    :class="onlineStatus"
  />
</template>

<style scoped lang="scss">
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
