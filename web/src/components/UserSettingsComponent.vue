<script setup lang="ts">
import { computed, ref, useTemplateRef } from 'vue';
import { Button, Dialog, InputText, Card, Textarea } from 'primevue';
import { Form, type FormSubmitEvent } from '@primevue/forms';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { useApi } from '@/composables/useApi';
import type { UpdateUserProfilePayload } from '@/interfaces/userProfile.interface';
import UserPictureComponent from './UserPictureComponent.vue';
import { useUserStore } from '@/stores/userStore';
import type { User } from '@/interfaces/user.interface';
import { useAuth0 } from '@auth0/auth0-vue';

const authenticatedUser = useAuthenticatedUserStore().authenticatedUserId;
const userStore = useUserStore();
const auth0 = useAuth0();

const userProfileResponse = useApi(
  `/users/${authenticatedUser}/profile`,
).json();
const userProfile = computed(() => userProfileResponse.data.value);

async function submit(event: FormSubmitEvent) {
  const payload: UpdateUserProfilePayload = {};
  if (event.states.username.dirty) {
    payload.username = event.states.username.value;
  } else {
    payload.username = userProfile.value.username;
  }

  if (event.states.bio.dirty) {
    payload.bio = event.states.bio.value;
  } else {
    payload.bio = userProfile.value.bio;
  }

  const newUser = await useApi('/users/me/profile').put(payload).json();
  userStore.updateUser(authenticatedUser, newUser.data.value);
}

const dialogVisible = ref(false);
const fileInputElement = useTemplateRef('fileInputElement');
const newImageSrc = ref('');
function onFileChange() {
  if (!fileInputElement.value || !fileInputElement.value.files) return;
  const file = fileInputElement.value.files[0];
  if (!file) return;
  newImageSrc.value = URL.createObjectURL(file);
}
async function submitImage() {
  if (!fileInputElement.value || !fileInputElement.value.files) return;
  const file = fileInputElement.value.files[0];
  if (!file) return;
  dialogVisible.value = false;
  const newUser = await useApi('/users/me/picture').post(file).json<User>();
  if (!newUser.data.value) return;
  userProfile.value.picture = newUser.data.value.picture;
  userStore.updateUser(authenticatedUser, newUser.data.value);
}
</script>

<template>
  <Card class="h-full w-full">
    <template #title><span class="text-xl font-bold">Settings</span></template>
    <template #content>
      <Form
        :initialValues="userProfile"
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
        <div class="flex flex-row items-center gap-4">
          <img
            :src="newImageSrc"
            :class="{ invisible: newImageSrc === '' }"
            class="w-24 h-24 rounded-full"
          />
          <input
            ref="fileInputElement"
            @change="onFileChange"
            type="file"
            accept="image/*"
            class="w-full"
          />
        </div>
        <template #footer>
          <Button @click="submitImage">Submit</Button>
        </template>
      </Dialog>
    </template>
    <template #footer>
      <div class="flex w-full justify-between">
        <div class="flex gap-2">
          <Button
            label="Change Your Password"
            severity="info"
            :to="'/not-implemented'"
            :as="'router-link'"
          />
          <Button @click="auth0.logout()" severity="danger">Logout</Button>
        </div>
        <div class="flex gap-2">
          <Button to="/" as="router-link" severity="secondary">Close</Button>
          <Button type="submit" form="settings-form">Save Settings</Button>
        </div>
      </div>
    </template>
  </Card>
</template>
