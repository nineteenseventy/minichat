import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import '../index.css';
import DialogService from 'primevue/dialogservice';

import App from './App.vue';
import router from './router';
import auth0 from './auth0';

const app = createApp(App);
const pinia = createPinia();

app.use(router);
app.use(pinia);
app.use(auth0);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: '.p-darkmode',
      cssLayer: {
        name: 'primevue',
        order: 'tailwind-base, primevue, tailwind-utilities',
      },
    },
  },
});
app.use(DialogService);

app.mount('#app');
