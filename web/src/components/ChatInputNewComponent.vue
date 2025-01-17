<script setup lang="ts">
import { computed, ref, useTemplateRef, watch } from 'vue';
import { useMessageDraftsStore } from '@/stores/messageDraftsStore';
import { useMessageStore } from '@/stores/messageStore';
import type { NewMessage } from '@/interfaces/message.interface';
import ChatInputComponent from './ChatInputComponent.vue';
import FilesPreviewComponent from './FilesPreviewComponent.vue';
import type { FilePreview } from './FilePreviewComponent.vue';
import { useApi } from '@/composables/useApi';

const props = defineProps<{ channelId: string }>();

const messageStore = useMessageStore();
const draftsStore = useMessageDraftsStore();

const content = ref('');

watch(
  () => props.channelId,
  (channelId, beforeChannelId) => {
    if (beforeChannelId) {
      draftsStore.setMessageDraft(beforeChannelId, content.value);
    }
    content.value = draftsStore.getMessageDraft(channelId) ?? '';
  },
  { immediate: true },
);

// File handling
const fileInputElement = useTemplateRef('fileInputElement');
function openFilePicker() {
  fileInputElement.value?.click();
}

const files = ref<File[]>([]);
const filePreviews = computed<FilePreview[]>(() =>
  files.value.map(
    (file) =>
      ({
        name: file.name,
        type: file.type,
        url: URL.createObjectURL(file),
        removeable: true,
      }) satisfies FilePreview,
  ),
);
function onFileChange() {
  if (!fileInputElement.value || !fileInputElement.value.files) return;
  if (fileInputElement.value.files.length < 1) return;
  files.value.push(...Array.from(fileInputElement.value.files));
  fileInputElement.value.value = '';
}
function removeFile(file: FilePreview) {
  const index = filePreviews.value.indexOf(file);
  if (index < 0) return;
  files.value.splice(index, 1);
}

// Send message
async function onSend() {
  const contentValue = content.value.trim();
  if (contentValue.length < 1) return;

  const messageFiles = files.value.map((file) => file);
  content.value = '';
  files.value = [];
  draftsStore.clearMessageDraft(props.channelId);

  const newMessage: NewMessage = {
    content: contentValue,
    attachments: messageFiles.map((file) => ({
      filename: file.name,
      type: file.type ?? 'application/octet-stream',
    })),
  };
  await messageStore
    .sendMessage(props.channelId, newMessage)
    .then((message) => {
      if (!message?.attachments) return;
      message.attachments.forEach((attachment) => {
        const fileIndex = messageFiles.findIndex(
          (file) => file.name === attachment.filename,
        );
        if (fileIndex < 0) return;
        const file = messageFiles[fileIndex];
        useApi(`/attachment/${attachment.id}`).post(file);
        files.value.splice(fileIndex, 1);
      });
    });
}

function checkForMentions(event: Event) {
  // match "@something" with no word at start and EOL or space at end
  const allMentionsRegex = /(?<!\w)@\w+(?=\s|$)/g;
  // match "@mention" only at end of text
  const activeMentionRegex = /(?<!\w)@\w+(?=$)/g;
  const mentions = content.value.match(allMentionsRegex);
  const activeMention = [...content.value.matchAll(activeMentionRegex)][0];
  if (activeMention) {
    const matchStart = activeMention.index;
    const matchEnd = matchStart + activeMention[0].length;
    const cursorPos = (event.target as HTMLTextAreaElement).selectionStart;
    console.debug(
      `active mention "${activeMention}" found from index ${matchStart} to ${matchEnd}`,
    );
    console.log('cursor at:', cursorPos);
    console.log('cursor in mention:', matchEnd - cursorPos >= 0);
  } else {
    console.log('no active Mention');
  }
}
</script>

<template>
  <div class="p-2 bg-highlight rounded-md mb-2" v-if="filePreviews.length > 0">
    <FilesPreviewComponent :files="filePreviews" @remove="removeFile($event)" />
  </div>
  <ChatInputComponent
    v-model="content"
    enableFile
    @onSave="onSend()"
    @onFile="openFilePicker()"
    @onInput="checkForMentions"
  />
  <input
    ref="fileInputElement"
    @change="onFileChange"
    type="file"
    multiple
    class="hidden"
    aria-label="file input"
  />
</template>

<style scoped lang="scss"></style>
