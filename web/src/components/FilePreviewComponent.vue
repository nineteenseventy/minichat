<script setup lang="ts">
export interface FilePreview {
  name: string;
  type: string;
  url: string;
  removeable?: boolean;
}

defineProps<FilePreview>();

const emit = defineEmits<{ remove: [] }>();
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
    </div>
    <div class="flex-1 flex items-center justify-center p-1 overflow-hidden">
      <img
        v-if="type.startsWith('image')"
        :src="url"
        class="max-w-full max-h-full aspect-square"
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
