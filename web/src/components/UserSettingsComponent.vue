<script setup lang="ts">
import { computed, ref, useTemplateRef, watchEffect } from 'vue';
import {
  Button,
  ColorPicker,
  Dialog,
  InputText,
  Card,
  Textarea,
} from 'primevue';
import { Form, type FormSubmitEvent } from '@primevue/forms';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { useApi } from '@/composables/useApi';
import type { UpdateUserProfilePayload } from '@/interfaces/userProfile.interface';
import UserPictureComponent from './UserPictureComponent.vue';

const authenticatedUser = useAuthenticatedUserStore().authenticatedUserId;

const user = useApi(`/users/${authenticatedUser}/profile`).json();
const initialValues = computed(() => user.data.value);

async function submit(event: FormSubmitEvent) {
  const payload: UpdateUserProfilePayload = {};
  console.log(event.states);
  if (event.states.username.dirty) {
    payload.username = event.states.username.value;
  } else {
    payload.username = initialValues.value.username;
  }

  if (event.states.bio.dirty) {
    payload.bio = event.states.bio.value;
  } else {
    payload.bio = initialValues.value.bio;
  }

  const newUser = await useApi('/users/me/profile').put(payload).json();
  console.log(newUser.data.value);
}

const dialogVisible = ref(false);
const fileInputElement = useTemplateRef('fileInputElement');
async function submitImage() {
  if (!fileInputElement.value || !fileInputElement.value.files) return;
  const file = fileInputElement.value.files[0];
  if (!file) return;
  const newUser = await useApi('/users/me/picture').post(file).json();
  if (!newUser.data) return;
  console.log(newUser.data.value);
}
</script>

<template>
  <Card class="h-full">
    <template #title><span class="text-xl font-bold">Settings</span></template>
    <template #content>
      <Form
        :initialValues="initialValues"
        @submit="submit"
        class="grid grid-cols-[auto,2fr] gap-4 items-center"
        id="settings-form"
      >
        <label for="picture">Profile Picture:</label>
        <UserPictureComponent
          :userId="authenticatedUser"
          @click="dialogVisible = true"
          class="h-20 w-20"
        />
        <label for="username">Username:</label>
        <InputText name="username" placeholder="Username" />
        <label for="bio">Bio:</label>
        <Textarea
          name="bio"
          placeholder="Bio"
          rows="5"
          class="min-h-11 max-h-36"
        />
      </Form>
      <Dialog header="Upload a new file" v-model:visible="dialogVisible" modal>
        <input ref="fileInputElement" type="file" accept="image/*" />
        <Button @click="submitImage">Submit</Button>
      </Dialog>
    </template>
    <template #footer>
      <div class="flex w-full justify-between">
        <Button
          label="Change Your Password"
          severity="info"
          :to="'/not-implemented'"
          :as="'router-link'"
        />
        <Button type="submit" form="settings-form">Save Settings</Button>
      </div>
    </template>
  </Card>
</template>
