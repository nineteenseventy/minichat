<script setup lang="ts">
import FilePreviewComponent, {
  type FilePreview,
} from './FilePreviewComponent.vue';

defineProps<{ files: FilePreview[] }>();
const emit = defineEmits<{ remove: [FilePreview] }>();
</script>

<template>
  <div class="flex flex-row gap-3">
    <FilePreviewComponent
      v-for="(file, i) in files"
      v-bind="file"
      :key="i"
      @remove="emit('remove', file)"
    />
  </div>
</template>
