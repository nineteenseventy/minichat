<script setup lang="ts">
import { ref } from 'vue';
import { useMessageDraftsStore } from '@/stores/messageDraftsStore';
import { onBeforeRouteUpdate } from 'vue-router';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';

const props = defineProps<{ channelId: string }>();

const draftsStore = useMessageDraftsStore();

const currentDraft = draftsStore.getMessageDraft(props.channelId);
const text = ref(currentDraft);

// TODO delete draft from store when a message is posted in the channel
onBeforeRouteUpdate(async (to) => {
  console.log(text.value);
  if (text.value) {
    console.log('draft saved');
    draftsStore.setMessageDraft(props.channelId, text.value);
  }
  const fetchedDraft = draftsStore.getMessageDraft(
    to.params.channelId as string,
  );
  console.log(fetchedDraft);
  text.value = fetchedDraft;
  console.log('draft fetched');
});
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
