<script setup lang="ts">
import { useApi } from '@/composables/useApi';
import { useAuth0 } from '@auth0/auth0-vue';
const { user } = useAuth0();

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const { data } = useApi('/users').json();
</script>

<template>
  <div>{{ user?.email }}</div>
  <ul>
    <li v-for="user in data.data" :key="user.id">{{ user.name }}</li>
  </ul>
</template>
