<script setup lang="ts">
import Button from 'primevue/button';
import Textarea from 'primevue/textarea';

const model = defineModel();

const emit = defineEmits<{
  onSave: [];
}>();

async function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    emit('onSave');
  }
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
    <Button icon="pi pi-send" class="self-end" @click="emit('onSave')" />
  </div>
</template>

<style scoped lang="scss"></style>
