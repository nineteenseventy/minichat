import { createFetch } from '@vueuse/core';
import { globalAuth0 } from '@/plugins/auth0';

export const useApi = createFetch({
  baseUrl: import.meta.env.VITE_API_URL,
  options: {
    async beforeFetch({ options }) {
      const token = await globalAuth0.getAccessTokenSilently();
      // if (!options.headers) options.headers = {};
      (<Record<string, string>>options.headers).Authorization =
        `Bearer ${token}`;
      return { options };
    },
  },
});
