import { defineStore } from 'pinia';
import type {
  GetMessagesQuery,
  Message,
  NewMessage,
  UpdateMessage,
} from '@/interfaces/message.interface';
import { computed, ref, type Ref } from 'vue';
import { useApi } from '@/composables/useApi';
import type { Result } from '@/interfaces/util.interface';
import { parseDate } from '@/utils/date/parseDate';
import { params } from '@/utils/url';

interface ClearMessagesOptions {
  channelId?: string;
  messageId?: string;
  before?: Date;
  after?: Date;
}

export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([]);

  function storeMessage(message: Message) {
    const existingMessageIndex = messages.value.findIndex(
      (v) => v.id === message.id,
    );
    if (existingMessageIndex !== -1)
      messages.value[existingMessageIndex] = message;
    else messages.value.push(message);
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
    data.value.data.forEach(storeMessage);
    sortMessages();
  }

  function sortMessages() {
    messages.value.sort(
      (a, b) =>
        parseDate(b.timestamp).getTime() - parseDate(a.timestamp).getTime(),
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

    const { channelId, before, after, messageId } = options;
    if (channelId) filters.push((v) => v.channelId === channelId);
    if (before) filters.push((v) => new Date(v.timestamp) < before);
    if (after) filters.push((v) => new Date(v.timestamp) > after);
    if (messageId) filters.push((v) => v.id === messageId);

    messages.value = messages.value.filter((v) => filters.every((f) => !f(v)));
  }

  async function sendMessage(
    channelId: string,
    newMessage: NewMessage,
  ): Promise<Message | undefined> {
    const request = useApi(`/messages/${channelId}`).post(newMessage);
    const { data } = await request.json<Message>();
    if (!data.value) return;
    storeMessage(data.value);
    sortMessages();
    return data.value;
  }

  async function updateMessage(messageId: string, newMessage: UpdateMessage) {
    const request = useApi(`/messages/${messageId}`).patch(newMessage);
    const { data } = await request.json<Message>();
    if (!data.value) return;
    storeMessage(data.value);
  }

  async function deleteMessage(messageId: string) {
    const request = useApi(`/messages/${messageId}`).delete();
    const { data } = await request.json<Message>();
    if (!data.value) return;
    clearMessages({ messageId });
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
    updateMessage,
    deleteMessage,
  };
});
