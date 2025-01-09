<script setup lang="ts">
const props = defineProps<{
  channelTitle: string;
}>();
</script>

<template>
  <div class="flex flex-row items-center cursor-pointer hover:underline">
    <span>{{ channelTitle }}</span>
  </div>
</template>
