import type { Plugin } from 'vue';
import { parse, populate } from 'dotenv';

export interface AssetEnvPluginOptions {
  path?: string;
}

const defaultPath = '/assets/env';
export const globalEnv: ImportMetaEnv = {} as ImportMetaEnv;
export async function loadGlobalEnv(options?: AssetEnvPluginOptions) {
  const path = options?.path ?? defaultPath;
  const response = await fetch(path);
  if (!response.ok) return;
  const txtEnv = await response.text();
  const parsed = parse(txtEnv);
  populate(globalEnv, parsed); // first load from .env file so that Vite env doesn't overwrite
  populate(globalEnv, import.meta.env); // then load from Vite env
}

// Vue Plugin Stuff

const ENV_KEY = 'env';
export const ENV_INJECTION_KEY = Symbol(ENV_KEY);

export const assetEnv: Plugin = {
  async install(app) {
    app.provide(ENV_INJECTION_KEY, globalEnv);
  },
};
