<script setup lang="ts">
import { globalAuth0 } from '@/plugins/auth0';
import Card from 'primevue/card';
import { ref, watch } from 'vue';

const { error } = globalAuth0;

const errors = ref([error.value]);

watch(
  () => error.value,
  (value) => {
    errors.value.push(value);
  },
);
</script>
<template>
  <Card class="w-full h-full m-4">
    <template #title> Callback Error </template>
    <template #content>
      <div v-for="error in errors" :key="error" class="text-red-500 mb-4">
        {{ error }}
      </div>
    </template>
  </Card>
</template>
