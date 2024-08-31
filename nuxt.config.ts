// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  modules: [
    "nuxt-auth-utils",
    "@nuxt/eslint",
    "@nuxtjs/color-mode",
    "@nuxtjs/device",
    "@pinia/nuxt",
    "@primevue/nuxt-module",
    "@vueuse/nuxt",
  ],
  alias: {
    /**
     * @fix `ERROR  Cannot find module 'pinia/dist/pinia.mjs'
     * @source https://stackoverflow.com/a/74801367/16237426
     */
    pinia: "/node_modules/@pinia/nuxt/node_modules/pinia/dist/pinia.mjs",
  },
});
