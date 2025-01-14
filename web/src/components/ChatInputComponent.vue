<script setup lang="ts">
import Button from 'primevue/button';
import Textarea from 'primevue/textarea';

defineProps<{
  enableCancel?: boolean;
}>();

const model = defineModel();

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
      icon="pi pi-send"
      class="self-end min-h-11 min-w-11"
      @click="emit('onSave')"
    />
  </div>
</template>

<style scoped lang="scss"></style>
