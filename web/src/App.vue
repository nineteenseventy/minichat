<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup';
import DynamicDialog from 'primevue/dynamicdialog';
</script>

<template>
  <DynamicDialog />
  <ConfirmPopup />
  <div class="flex p-2 h-full w-full">
    <RouterView />
  </div>
</template>
