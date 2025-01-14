import { createFetch } from '@vueuse/core';
import { globalEnv } from '@/plugins/assetEnvPlugin';
import { globalAuth0 } from '@/plugins/auth0';

console.log('use Api gets imported');
export const useApi = createFetch({
  baseUrl: globalEnv.VITE_API_URL,
  options: {
    async beforeFetch({ options }) {
      const token = await globalAuth0.getAccessTokenSilently();
      (<Record<string, string>>options.headers).Authorization =
        `Bearer ${token}`;
      return { options };
    },
  },
});
