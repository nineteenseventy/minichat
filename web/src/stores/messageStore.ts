import { defineStore } from 'pinia';
import type { Message } from '@/interfaces/message.interface';
import { ref } from 'vue';

const testMessages: Message[] = [
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

interface ClearMessagesOptions {
  channelId?: string;
  before?: Date;
  after?: Date;
}

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>(testMessages);

  function storeMessage(message: Message) {
    messages.value.push(message);
  }

  function getMessage(channelId: string) {
    return messages.value.filter((v) => v.channelId === channelId);
  }

  function getMessages(channelId?: string) {
    if (channelId) {
      // return messages.value.filter((v) => v.channelId === channelId);
      return messages.value.map((v) => ({ ...v, channelId }));
    }
    return messages.value;
  }

  function clearMessages(options: ClearMessagesOptions) {
    const filters: ((v: Message) => boolean)[] = [];

    const { channelId, before, after } = options;
    if (channelId) filters.push((v) => v.channelId === channelId);
    if (before) filters.push((v) => new Date(v.timestamp) < before);
    if (after) filters.push((v) => new Date(v.timestamp) > after);

    messages.value = messages.value.filter((v) => filters.every((f) => f(v)));
  }

  return { messages, storeMessage, getMessage, getMessages, clearMessages };
});
