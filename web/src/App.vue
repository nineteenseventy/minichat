<script setup lang="ts">
import DynamicDialog from 'primevue/dynamicdialog';
import { RouterView } from 'vue-router';
</script>

<template>
  <DynamicDialog />
  <RouterView />
</template>

<style scoped>
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}
</style>
