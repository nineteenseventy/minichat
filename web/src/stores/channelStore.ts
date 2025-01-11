import { useApi } from '@/composables/useApi';
import type { Channel } from '@/interfaces/channel.interface';
import type { Result } from '@/interfaces/util.interface';
import { defineStore } from 'pinia';
import { computed, ref, type Ref } from 'vue';

export const useChannelStore = defineStore('channel', () => {
  const channels = ref<Channel[]>([]);

  function getChannel(channelId: Ref<string>): Ref<Channel | undefined> {
    const channel = computed(() =>
      channels.value.find((v) => v.id === channelId.value),
    );
    if (!channel.value) {
      fetchChannel(channelId.value).then((fetchedChannel) => {
        if (!fetchedChannel) return;
        if (channels.value.find((v) => v.id === fetchedChannel.id)) return;
        channels.value.push(fetchedChannel);
      });
    }
    return channel;
  }

  async function updateStore() {
    const newChannels = await fetchChannels();
    if (!newChannels) return;
    channels.value = newChannels;
  }

  function storeChannel(channel: Channel) {
    const index = channels.value.findIndex((v) => v.id === channel.id);
    if (index === -1) {
      channels.value.push(channel);
    } else {
      channels.value[index] = channel;
    }
  }

  async function fetchChannels() {
    const { data } = await useApi('/channels').json<Result<Channel[]>>();
    return data.value?.data ?? undefined;
  }

  async function fetchChannel(channelId: string) {
    const { data } = await useApi(`/channels/${channelId}`).json<Channel>();
    return data.value ?? undefined;
  }

  async function getDirectChannel(userId: string) {
    const request = useApi(`/users/${userId}/channel`);
    const { data: channel } = await request.json<Channel>();
    if (!channel.value) throw new Error('Channel could not be created!');
    storeChannel(channel.value);
    return channel.value;
  }

  async function setRead(channelId: string) {
    await useApi(`/channels/${channelId}/setRead`).patch();
    const channelIndex = channels.value.findIndex((v) => v.id === channelId);
    if (channelIndex !== -1) channels.value[channelIndex].unreadCount = 0;
  }

  return {
    channels,
    getChannel,
    updateStore,
    storeChannel,
    getDirectChannel,
    setRead,
  };
});
