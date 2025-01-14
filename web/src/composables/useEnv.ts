import { ENV_INJECTION_KEY } from '@/plugins/assetEnvPlugin';
import { inject } from 'vue';

export const useEnv = () => {
  return inject(ENV_INJECTION_KEY) as ImportMetaEnv;
};
