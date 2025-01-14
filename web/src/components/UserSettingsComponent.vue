<script setup lang="ts">
import { reactive, computed } from 'vue';
import {
  Button,
  ColorPicker,
  FileUpload,
  InputText,
  Panel,
  Textarea,
  useConfirm,
} from 'primevue';
import InputMask from 'primevue/inputmask';
import { useAuthenticatedUserStore } from '@/stores/authenticatedUserStore';
import { splitAndCapitalizeCamelCase } from '@/utils/strings/splitAndCapitalizeCamelCase';

interface UserSettingsPictureField {
  picture?: string;
}
interface UserSettingsField {
  type: 'InputText' | 'TextArea' | 'FileUpload' | 'ColorPicker';
  vModel?: string;
  placeholder?: string | undefined;
}

type UserSettings = Record<
  string,
  UserSettingsField & UserSettingsPictureField
>;

const user = useAuthenticatedUserStore().getProfile();

const fields = reactive<UserSettings>({
  userName: {
    type: 'InputText',
    vModel: '',
    placeholder: computed(() => user.value.username).value,
  },
  bio: {
    type: 'TextArea',
    vModel: '',
    placeholder: computed(() => user.value.bio).value,
  },
  profilePicture: {
    type: 'FileUpload',
    picture: user.value.picture,
  },
  color: {
    type: 'ColorPicker',
    vModel: user.value.color,
  },
});

function handleKeydown(event: KeyboardEvent) {
  const regex = /^[a-fA-F0-9]$/;
  const allowedKeys = [
    'ArrowLeft',
    'ArrowRight',
    'ArrowUp',
    'ArrowDown',
    'End',
    'Home',
    'Backspace',
    'Delete',
    'Tab',
    'Control',
    'Shift',
    'Alt',
    'Meta',
  ];
  console.log(event.key);
  if (
    !regex.test(event.key) &&
    !(event.ctrlKey || event.metaKey) &&
    !allowedKeys.includes(event.key)
  ) {
    event.preventDefault();
  }
}

const confirm = useConfirm();
function confirmSaveSettings(event: Event) {
  console.log(event.target);
  confirm.require({
    target: event.target as HTMLButtonElement,
    message: 'Do you want to save your settings?',
    header: 'Confirmation',
    icon: 'pi pi-info-circle',
    rejectProps: {
      label: 'Cancel',
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: 'Save',
      severity: 'primary',
    },
    accept: () => {
      console.log('accepted');
    },
    reject: () => {
      console.log('rejected');
    },
  });
}
</script>

<template>
  <Panel style="null">
    <template #header><span class="text-xl font-bold">Settings</span></template>
    <form id="settings-form">
      <template v-for="(item, key) in fields" :key="key">
        <label :for="key">{{ splitAndCapitalizeCamelCase(key) }}:</label>
        <InputText
          v-if="item.type === 'InputText'"
          v-model.trim="item.vModel"
          :placeholder="item.placeholder"
        />
        <Textarea
          v-else-if="item.type === 'TextArea'"
          v-model.trim="item.vModel"
          :placeholder="item.placeholder"
          rows="5"
          class="min-h-11 max-h-36"
        />
        <template v-else-if="item.type === 'FileUpload'">
          <FileUpload
            mode="advanced"
            accept="image/*"
            :showUploadButton="false"
            :maxFileSize="4000000"
          >
            <template #empty>
              <img :src="item.picture" class="file-upload-picture" />
              <span>Drag and drop your new image here (<i>max. 4MB</i>).</span>
            </template>
          </FileUpload>
        </template>
        <div v-else-if="item.type === 'ColorPicker'" class="color-picker">
          <ColorPicker
            v-model="item.vModel"
            :inputId="key"
            format="hex"
            inline
          />
          <InputMask
            v-model="item.vModel"
            placeholder="#FF0000"
            mask="#******"
            @keydown="handleKeydown($event as KeyboardEvent)"
            ref="mask"
          />
        </div>
      </template>
      <Button
        label="Change Your Password"
        severity="info"
        :to="'/not-implemented'"
        :as="'router-link'"
      />
    </form>
    <template #footer>
      <Button
        @click.prevent="confirmSaveSettings"
        type="submit"
        form="settings-form"
        >Save Settings</Button
      >
    </template>
  </Panel>
</template>

<style scoped lang="scss">
.p-panel {
  height: 100%;
  align-content: normal !important;
}
#settings-form {
  display: grid;
  grid-template-columns: auto 2fr;
  gap: 1rem;
  label {
    grid-column-start: 1;
    grid-column-end: 1;
    padding-right: 1rem;
    display: flex;
    align-items: center;
  }
  .file-upload-picture {
    display: inline;
    margin-right: 1rem;
    width: 4rem;
  }
  .color-picker {
    display: flex;
    flex-direction: column;
    Button {
      margin-top: auto;
      margin-bottom: auto;
      display: inline;
    }
  }
}
:deep(.p-panel-footer) {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}
</style>
