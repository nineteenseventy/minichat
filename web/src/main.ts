import './assets/main.css';
import 'primeicons/primeicons.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import PrimeVue from 'primevue/config';
import Theme from './theme/theme';
import '../index.css';
import DialogService from 'primevue/dialogservice';
import ConfirmationService from 'primevue/confirmationservice';

import App from './App.vue';
import router from './router';
import { auth0 } from './plugins/auth0';
import Ripple from 'primevue/ripple';
import { assetEnv, loadGlobalEnv } from './plugins/assetEnvPlugin';

async function main() {
  // throw new Error('Not implemented');
  await loadGlobalEnv();

  const app = createApp(App);
  const pinia = createPinia();

  app.use(assetEnv);
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
  app.use(ConfirmationService);
  app.directive('ripple', Ripple);

  app.mount('app-root');
}

main().catch(console.error);
