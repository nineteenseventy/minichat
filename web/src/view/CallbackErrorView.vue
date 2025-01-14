<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';
import Card from 'primevue/card';

const { error } = useAuth0();
console.log(error.value);
</script>
<template>
  <Card>
    <template #title> Callback </template>
    <template #content> {{ error }} </template>
  </Card>
</template>
