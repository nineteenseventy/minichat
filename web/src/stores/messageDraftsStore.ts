import { defineStore } from 'pinia';
import { ref } from 'vue';

interface StoredMessageDraft {
  channelId: string;
  content: string;
}

export const useMessageDraftsStore = defineStore('messageDrafts', () => {
  const messageDrafts = ref<StoredMessageDraft[]>([]);
  function getMessageDraft(channelId: string) {
    return messageDrafts.value.find((v) => v.channelId === channelId)?.content;
  }
  function setMessageDraft(channelId: string, content?: string) {
    if (content) {
      messageDrafts.value.push({
        channelId: channelId,
        content: content,
      });
    } else {
      messageDrafts.value = messageDrafts.value.filter(
        (v) => v.channelId !== channelId,
      );
    }
  }

  function clearMessageDraft(channelId: string) {
    setMessageDraft(channelId, undefined);
  }
  return { messageDrafts, getMessageDraft, setMessageDraft, clearMessageDraft };
});
