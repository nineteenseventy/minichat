<script setup lang="ts">
import Card from 'primevue/card';
import ChannelTitleComponent from '@/components/ChannelTitleComponent.vue';
import ChatInputComponent from '@/components/ChatInputComponent.vue';
import { useRouteParam } from '@/composables/useRouteParam';
import MessagesComponent from '@/components/MessagesComponent.vue';

const channelId = useRouteParam('channelId');
</script>

<template>
  <Card class="h-full">
    <template #title>
      <ChannelTitleComponent v-if="channelId" :channelId="channelId" />
    </template>
    <template #content>
      <MessagesComponent v-if="channelId" :channelId="channelId" />
    </template>
    <template #footer>
      <ChatInputComponent v-if="channelId" :channelId="channelId" />
    </template>
  </Card>
</template>
