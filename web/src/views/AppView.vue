<script setup lang="ts">
import { useApi } from '@/composables/useApi';
import { useAuth0 } from '@auth0/auth0-vue';
import { ref, type Ref } from 'vue';
const { user } = useAuth0();
const api = useApi();

const users: Ref<any> = ref([]);
api.get('/users').then((response) => {
  users.value = response.data;
});
</script>

<template>
  <div>{{ user?.name }}</div>
  <ul>
    <li v-for="user in users" :key="user.id">{{ user.name }}</li>
  </ul>
</template>
