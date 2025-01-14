import {
  createAuth0,
  type Auth0VueClientOptions,
  type Auth0VueClient,
  type Auth0PluginOptions,
} from '@auth0/auth0-vue';
import { type Plugin } from 'vue';
import { globalEnv } from './assetEnvPlugin';

export let globalAuth0: Auth0VueClient;

export const auth0: Plugin = {
  install(app) {
    const redirectUri = new URL('/callback', globalEnv.VITE_EXTERNAL_URL).href;
    console.log(globalEnv.VITE_EXTERNAL_URL);
    const options: Auth0VueClientOptions = {
      domain: globalEnv.VITE_AUTH0_DOMAIN,
      clientId: globalEnv.VITE_AUTH0_CLIENT_ID,
      authorizationParams: {
        redirect_uri: redirectUri,
        audience: globalEnv.VITE_AUTH0_AUDIENCE,
      },
    };

    const clientOptions: Auth0PluginOptions = {
      errorPath: '/callback',
      skipRedirectCallback: true,
    };

    const auth0 = createAuth0(options, clientOptions);
    globalAuth0 = auth0;
    app.use(auth0);
  },
};
