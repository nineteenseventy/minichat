import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import DialogService from 'primevue/dialogservice';

import App from './App.vue';
import router from './router';
import auth0 from './auth0';

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(auth0);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: '.p-darkmode',
    },
  },
});
app.use(DialogService);

app.mount('#app');
