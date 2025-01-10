<script setup lang="ts">
import { ref } from 'vue';
import { useMessageDraftsStore } from '@/stores/messageDraftsStore';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';

const props = defineProps<{ channelId: string }>();

const draftsStore = useMessageDraftsStore();
const text = ref('');

if (props.channelId) {
  const draft = draftsStore.getMessageDraft(props.channelId);
  if (draft) text.value = draft;
}
</script>

<template>
  <div class="flex flex-row gap-3 items-center">
    <Textarea
      autoResize
      fluid
      rows="1"
      class="min-h-11 max-h-36"
      v-model="text"
    />
    <Button icon="pi pi-send" class="self-end" />
  </div>
</template>

<style scoped lang="scss"></style>
