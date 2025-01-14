import type { Plugin } from 'vue';
import { parse, populate } from 'dotenv';

export interface AssetEnvPluginOptions {
  path?: string;
}

// export type AssetEnvPluginOptions = _AssetEnvPluginOptions | undefined;

const defaultPath = '/assets/.env';
const processEnv = import.meta.env;

export const assetEnvLoader: Plugin<AssetEnvPluginOptions> = {
  async install(_, options = {}) {
    const path = options?.path ?? defaultPath;
    const response = await fetch(path);
    if (!response.ok) return;
    const env = await response.text();
    const parsed = parse(env);
    populate(parsed, processEnv);
  },
};
