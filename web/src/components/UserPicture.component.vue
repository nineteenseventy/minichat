<script setup lang="ts">
import defaultUser from '@/assets/images/default-user.png';

defineProps<{
  picture?: string;
}>();
</script>

<template>
  <div class="picture">
    <img :src="picture ?? defaultUser" alt="Profile picture" />
  </div>
</template>

<style scoped>
.picture {
  display: flex;
  border-radius: 50%;
  overflow: hidden;
  width: 3rem;
  height: 3rem;
  img {
    width: 3rem;
    height: auto;
  }
}
</style>
