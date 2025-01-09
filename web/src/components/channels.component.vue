<script setup lang="ts">
import Menu from 'primevue/menu';
import type { MenuItem } from 'primevue/menuitem';
import { ref } from 'vue';

const selectedElement = ref();

const menuItems = ref<MenuItem[]>([]);
</script>

<template>
  <Menu :model="menuItems" />
</template>

<style lang="css">
.menu:not(:only-child) {
  width: 100%;
}
</style>
