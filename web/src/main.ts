import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import PrimeVue from 'primevue/config';
import Theme from './theme/theme';
import '../index.css';
import DialogService from 'primevue/dialogservice';

import App from './App.vue';
import router from './router';
import auth0 from './auth0';
import { assetEnvLoader } from './plugins/AssetEnvPlugin';

const app = createApp(App);
const pinia = createPinia();

app.use(assetEnvLoader, {});
app.use(router);
app.use(pinia);
app.use(auth0);
app.use(PrimeVue, {
  theme: {
    preset: Theme,
    options: {
      darkModeSelector: 'system',
      cssLayer: {
        name: 'primevue',
        order: 'tailwind-base, primevue, tailwind-utilities',
      },
    },
  },
});
app.use(DialogService);

app.mount('app-root');
