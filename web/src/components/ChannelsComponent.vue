<script setup lang="ts">
import Menu from 'primevue/menu';
import type { MenuItem } from 'primevue/menuitem';
import { ref } from 'vue';
import type {
  ChannelsResponse,
  PublicChannel,
  DirectChannel,
  GroupChannel,
} from '@/interfaces/channel.interface';
import { useTimeoutPoll } from '@vueuse/core';
import { useApi } from '@/composables/useApi';
import { useRouter } from 'vue-router';
import { MenuItemCommandEvent } from 'primevue/menuitem';

const router = useRouter();

const menuItems = ref<MenuItem[]>([]);
const selectedChannelId = ref<string | null>(null);

useTimeoutPoll(
  async () => {
    const response = await useApi('/channels').json<ChannelsResponse>();
    menuItems.value = mapMenuItems(response.data.value!);
  },
  60000,
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

function mapMenuItems(data: ChannelsResponse): MenuItem[] {
  return [
    {
      label: 'Public Channels',
      items: data.public.map(mapPublicChannel),
    },
    {
      separator: true,
    },
    {
      label: 'Private Channels',
      items: data.private.map((c) =>
        c.type === 'group' ? mapGroupChannel(c) : mapDirectChannel(c),
      ),
    },
  ];
}

function mapPublicChannel(channel: PublicChannel): MenuItem {
  return {
    label: channel.title,
    url: mapUrl(channel.id),
    badge: channel.unreadCount ? channel.unreadCount : null,
    key: channel.id,
    class: { 'font-bold': selectedChannelId.value === channel.id },
    command: clickCommandFactory(channel.id),
  };
}

function mapDirectChannel(channel: DirectChannel): MenuItem {
  return {
    label: channel.title,
    url: mapUrl(channel.id),
    badge: channel.unreadCount ? channel.unreadCount : null,
    key: channel.id,
    class: { 'font-bold': selectedChannelId.value === channel.id },
    command: clickCommandFactory(channel.id),
  };
}

function mapGroupChannel(channel: GroupChannel): MenuItem {
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
  <Menu :model="menuItems" />
</template>

<style lang="css">
.menu:not(:only-child) {
  width: 100%;
}
</style>
