import { createFetch, useFetch } from '@vueuse/core';
import auth0 from '@/auth0';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type BetterUseFetch = typeof useFetch<any>;

export const useApi: BetterUseFetch = createFetch({
  baseUrl: import.meta.env.VITE_API_URL,
  options: {
    async beforeFetch({ options }) {
      const token = await auth0.getAccessTokenSilently();
      // if (!options.headers) options.headers = {};
      (<Record<string, string>>options.headers).Authorization =
        `Bearer ${token}`;
      return { options };
    },
  },
});
