<script setup lang="ts">
import { useApi } from '@/composables/useApi';
import { useAuth0 } from '@auth0/auth0-vue';
const { user } = useAuth0();
const { data } = useApi('/users');
</script>

<template>
  <div>{{ user?.name }}</div>
  <ul>
    <li v-for="user in data?.user" :key="user.id">{{ user.name }}</li>
  </ul>
</template>
