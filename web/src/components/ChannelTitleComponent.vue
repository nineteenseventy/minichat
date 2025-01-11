<script setup lang="ts">
import { useChannelStore } from '@/stores/channelStore';
import { computed } from 'vue';

const channelStore = useChannelStore();
const props = defineProps<{
  channelId: string;
}>();

const channel = channelStore.getChannel(computed(() => props.channelId));
</script>

<template>
  <div class="flex flex-row items-center text-base font-bold color-surface-400">
    <span>{{ channel?.title }}</span>
  </div>
</template>
mes
