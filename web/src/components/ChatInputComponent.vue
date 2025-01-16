<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';
import { Button, Textarea, Popover } from 'primevue';
import { useApi } from '@/composables/useApi';
import type { Message } from '@/interfaces/message.interface';

defineProps<{
  enableCancel?: boolean;
}>();

const model = defineModel<string>();

const emit = defineEmits<{
  onSave: [];
  onCancel: [];
}>();

async function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    emit('onSave');
  }
}

const fileUploadPopOver = useTemplateRef('fileUploadPopOver');
function showFileAttachmentPopover(event: Event) {
  fileUploadPopOver.value?.toggle(event);
}

const fileInputElement = useTemplateRef('fileInputElement');
async function submitImage() {
  if (!fileInputElement.value || !fileInputElement.value.files) return;
  const file = fileInputElement.value.files[0];
  if (!file) return;
  const { error } = await useApi(`/messages/${id}`).post(file).json<Message>();
  if (error) console.log('Error posting image: ', error);
  fileUploadPopOver.value?.hide();
}

const newImageSrc = ref('');
function onFileChange() {
  if (!fileInputElement.value || !fileInputElement.value.files) return;
  const file = fileInputElement.value.files[0];
  if (!file) return;
  newImageSrc.value = URL.createObjectURL(file);
}
</script>

<template>
  <div class="flex flex-row gap-3">
    <Textarea
      autoResize
      fluid
      rows="1"
      class="min-h-11 max-h-36"
      v-model="model"
      @keydown="onKeydown($event)"
    />
    <Button
      v-if="enableCancel"
      icon="pi pi-times"
      class="self-end min-h-11 min-w-11"
      severity="danger"
      @click="emit('onCancel')"
    />
    <Button
      icon="pi pi-paperclip"
      class="min-h-11 min-w-11"
      severity="help"
      @click="showFileAttachmentPopover"
    />
    <Button
      icon="pi pi-send"
      class="self-end min-h-11 min-w-11"
      @click="emit('onSave')"
    />
    <Popover ref="fileUploadPopOver">
      <div class="flex flex-col items-center gap-4">
        <input
          ref="fileInputElement"
          @change="onFileChange"
          type="file"
          accept="image/*"
          class="w-full"
        />

        <Button @click="submitImage" class="self-end">Submit</Button>
      </div>
    </Popover>
  </div>
</template>

<style scoped lang="scss"></style>
