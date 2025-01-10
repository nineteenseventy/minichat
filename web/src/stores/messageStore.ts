import { defineStore } from 'pinia';
import { ref } from 'vue';

import type { Message } from '@/interfaces/message.interface';

const testMessages: Message[] = [
  {
    authorId: '2d0a4682-a7a3-4461-98bc-4b403a94f000',
    content: 'Hello, World!',
    id: '1',
    timestamp: new Date().toISOString(),
    attachments: [],
    channelId: 'Global Channel',
    read: false,
  },
  {
    authorId: '038678a5-b6f1-45dc-b1da-9f3837f4cdc8',
    content: 'Hello, World!',
    id: '2',
    timestamp: new Date().toISOString(),
    attachments: [],
    channelId: 'Global Channel',
    read: false,
  },
  {
    authorId: '2d0a4682-a7a3-4461-98bc-4b403a94f000',
    content: 'Hello, World!',
    id: '3',
    timestamp: new Date().toISOString(),
    attachments: [],
    channelId: 'Global Channel',
    read: false,
  },
];

export const useMessageStore = defineStore('messageDrafts', () => {
  const messages = ref<Message[]>([]);
  const PAGE_SIZE = 20;
  function storeMessages(messagesToStore: Message[] | Message) {
    if (Array.isArray(messagesToStore)) {
      messages.value.push(...messagesToStore);
    } else {
      messages.value.push(messagesToStore);
    }
  }
  async function getMessages(
    channelId: string,
    fromTimestamp?: Date | string,
    reverse?: boolean,
  ): Promise<Message[]> {
    const channelMessages = messages.value.filter(
      (v) => v.channelId === channelId,
    );
    if (
      channelMessages.length === 0 ||
      (fromTimestamp && _needsFetching(channelMessages, fromTimestamp, reverse))
    ) {
      await _fetchMessages(channelId, fromTimestamp, reverse);
      return messages.value;
    }
    return [];

    function _needsFetching(
      channelMessages: Message[],
      fromtimestamp: Date | string,
      reverse?: boolean,
    ): boolean {
      if (channelMessages.length === 0) return true;
      const timestamp = new Date(fromtimestamp);
      if (reverse) {
        return new Date(channelMessages[0].timestamp) > timestamp;
      } else {
        return (
          new Date(channelMessages[channelMessages.length - 1].timestamp) <
          timestamp
        );
      }
    }
    async function _fetchMessages(
      channelId: string,
      fromTimestamp: string,
      reverse: boolean,
    ) {
      const fetchedMessages = testMessages.find(
        (v) => v.channelId == channelId,
      );
      if (reverse) {
        fetchedMessages.reverse();
      }
      storeMessages(fetchedMessages);
    }
  }
  return { messages, getMessages, storeMessages };
});
