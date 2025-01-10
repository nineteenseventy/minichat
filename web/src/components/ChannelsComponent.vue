<script setup lang="ts">
import Menu from 'primevue/menu';
import type { MenuItem } from 'primevue/menuitem';
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import type { MenuItemCommandEvent } from 'primevue/menuitem';
import { unpackRouterParam } from '@/utils/router';
import { useChannelStore } from '@/stores/channelStore';
import type { Channel } from '@/interfaces/channel.interface';

const router = useRouter();
const route = useRoute();
const channelStore = useChannelStore();

const selectedChannelId = ref<string | undefined>();

watch(
  () => unpackRouterParam(route.params.channelId),
  (newChannelId, oldChannelId) => {
    if (newChannelId === oldChannelId) return;
    selectedChannelId.value = newChannelId;
  },
  { immediate: true },
);

function mapUrl(channelId: string): string {
  return `/channels/${channelId}`;
}

function clickCommandFactory(
  channelId: string,
): (e: MenuItemCommandEvent) => void {
  return (e: MenuItemCommandEvent) => {
    e.originalEvent.preventDefault();
    router.push(mapUrl(channelId));
  };
}

const menuItems = computed(() => {
  return mapMenuItems(channelStore.channels);
});

function mapMenuItems(channels: Channel[]): MenuItem[] {
  return [
    {
      label: 'Public Channels',
      items: channels.filter((c) => c.type === 'public').map(mapChannel),
    },
    {
      separator: true,
    },
    {
      label: 'Private Channels',
      items: channels
        .filter((c) => ['direct', 'group'].includes(c.type))
        .map(mapChannel),
    },
  ];
}

function mapChannel(channel: Channel): MenuItem {
  return {
    label: channel.title,
    url: mapUrl(channel.id),
    badge: channel.unreadCount ? channel.unreadCount : null,
    key: channel.id,
    class: { 'font-bold': selectedChannelId.value === channel.id },
    command: clickCommandFactory(channel.id),
  };
}
</script>

<template>
  <Menu :model="menuItems" :key="selectedChannelId" />
</template>

<style lang="css">
.menu:not(:only-child) {
  width: 100%;
}
</style>
