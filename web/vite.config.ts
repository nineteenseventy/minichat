import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import VueRouter from 'unplugin-vue-router/vite';
import vue from '@vitejs/plugin-vue';
import vueDevTools from 'vite-plugin-vue-devtools';
import Components from 'unplugin-vue-components/vite';
import { PrimeVueResolver } from 'unplugin-vue-components/resolvers';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    VueRouter({
      /* options */
      root: fileURLToPath(new URL('./', import.meta.url)),
    }),
    vue(),
    vueDevTools(),
    Components({
      dirs: ['src/components', 'src/views'],
      resolvers: [PrimeVueResolver()],
      directoryAsNamespace: true,
    }),
  ],
  envDir: '../',
  build: {
    outDir: '../dist/web',
    emptyOutDir: true,
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
});
