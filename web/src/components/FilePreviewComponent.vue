<script setup lang="ts">
import { saveAs } from 'file-saver-es';

export interface FilePreview {
  name: string;
  type: string;
  url: string;
  removeable?: boolean;
}

const props = defineProps<
  FilePreview & {
    enableDownload?: boolean;
    enableLargePreview?: boolean;
  }
>();

const emit = defineEmits<{ remove: [] }>();

function onLargePreview() {
  if (!props.enableLargePreview) return;
  window.open(props.url, '_blank');
}
</script>

<template>
  <div
    class="flex w-32 h-32 border border-gray-200 rounded-lg flex-col items-center"
  >
    <div class="flex flex-row justify-between w-full p-1">
      <span class="text-sm font-semibold truncate">{{ name }}</span>
      <i
        v-if="removeable"
        @click="emit('remove')"
        class="pi pi-times cursor-pointer"
      ></i>
      <i
        v-if="enableDownload"
        @click="saveAs(url, name)"
        class="pi pi-download cursor-pointer"
      ></i>
    </div>
    <div
      class="flex-1 flex items-center justify-center p-1 overflow-hidden"
      :class="{ 'cursor-pointer': enableLargePreview }"
      @click="onLargePreview()"
    >
      <img
        v-if="type.startsWith('image')"
        :src="url"
        class="max-w-full max-h-full aspect-square rounded-md pi pi-image"
        alt="file preview"
      />
      <div v-else class="flex flex-col gap-1">
        <i class="pi pi-file text-5xl"></i>
        <p class="text-xs text-gray-500 truncate">{{ type }}</p>
      </div>
    </div>
    <!-- <div v-else class="flex flex-col gap-1">
      <p class="text-xs text-gray-500 truncate">{{ type }}</p>
    </div> -->
  </div>
</template>

<style scoped lang="scss">
img {
  position: relative;
  &::before {
    content: '';
    display: block;
    padding-top: 100%;
    background-color: rgba($color: white, $alpha: 0.1);
  }
  &::after {
    content: '\e972';
    font-size: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    z-index: 2;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }
}
</style>
