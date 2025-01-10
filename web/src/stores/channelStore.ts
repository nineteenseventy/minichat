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

  async function fetchChannels() {
    const { data } = await useApi('/channels').json<Result<Channel[]>>();
    return data.value?.data ?? undefined;
  }

  async function fetchChannel(channelId: string) {
    const { data } = await useApi(`/channels/${channelId}`).json<Channel>();
    return data.value ?? undefined;
  }

  return { channels, getChannel, updateStore };
});
