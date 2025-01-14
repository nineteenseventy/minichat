<script setup lang="ts">
import Card from 'primevue/card';
import ChannelTitleComponent from '@/components/ChannelTitleComponent.vue';
import ChatInputComponent from '@/components/ChatInputNewComponent.vue';
import { useRouteParam } from '@/composables/useRouteParam';
import MessagesComponent from '@/components/MessagesComponent.vue';

const channelId = useRouteParam('channelId');
</script>

<template>
  <Card v-if="channelId" class="h-full">
    <template #title>
      <ChannelTitleComponent :channelId="channelId" />
    </template>
    <template #content>
      <MessagesComponent :channelId="channelId" class="h-full max-h-full" />
    </template>
    <template #footer>
      <ChatInputComponent :channelId="channelId" />
    </template>
  </Card>
</template>
