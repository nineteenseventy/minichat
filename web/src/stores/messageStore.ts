import { defineStore } from 'pinia';
import type {
  GetMessagesQuery,
  Message,
  NewMessage,
} from '@/interfaces/message.interface';
import { computed, ref, type Ref } from 'vue';
import { useApi } from '@/composables/useApi';
import type { Result } from '@/interfaces/util.interface';
import { parseDate } from '@/utils/date/parseDate';
import { params } from '@/utils/url';

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

  function getMessage(messageId: Ref<string>) {
    return computed(() => messages.value.find((v) => v.id === messageId.value));
  }

  async function loadMessages(
    channelId: string,
    before?: Date,
    after?: Date,
    count = 50,
  ) {
    const queryData: GetMessagesQuery = {
      count,
      before: before?.toISOString(),
      after: after?.toISOString(),
    };
    const request = useApi(`/messages/${channelId}?${params(queryData)}`);
    const { data } = await request.json<Result<Message[]>>();
    if (!data.value) return;
    const newMessages = data.value.data ?? [];
    messages.value = newMessages.concat(messages.value);
    sortMessages();
  }

  function sortMessages() {
    messages.value.sort(
      (a, b) =>
        parseDate(a.timestamp).getTime() - parseDate(b.timestamp).getTime(),
    );
  }

  function getMessages(channelId: Ref<string>) {
    return computed(() =>
      messages.value.filter((v) => v.channelId === channelId.value),
    );
  }

  function getMessageIds(channelId: Ref<string>) {
    return computed(() =>
      messages.value
        .filter((v) => v.channelId === channelId.value)
        .map((v) => v.id),
    );
  }

  function clearMessages(options: ClearMessagesOptions) {
    const filters: ((v: Message) => boolean)[] = [];

    const { channelId, before, after } = options;
    if (channelId) filters.push((v) => v.channelId === channelId);
    if (before) filters.push((v) => new Date(v.timestamp) < before);
    if (after) filters.push((v) => new Date(v.timestamp) > after);

    messages.value = messages.value.filter((v) => filters.every((f) => f(v)));
  }

  async function sendMessage(channelId: string, newMessage: NewMessage) {
    const request = useApi(`/messages/${channelId}`, {
      method: 'POST',
      body: JSON.stringify(newMessage),
    });
    const { data } = await request.json<Message>();
    if (!data.value) return;
    storeMessage(data.value);
  }

  return {
    messages,
    storeMessage,
    getMessage,
    getMessages,
    clearMessages,
    loadMessages,
    getMessageIds,
    sendMessage,
  };
});
