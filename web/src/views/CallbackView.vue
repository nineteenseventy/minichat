<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';
import { useRouter } from 'vue-router';

const auth0 = useAuth0();
const router = useRouter();

auth0.handleRedirectCallback().then(() => {
  router.push('/');
});
</script>

<template>
  <div></div>
</template>
