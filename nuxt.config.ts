// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  devtools: { enabled: true },
  modules: [
    'nuxt-auth-utils',
    '@nuxt/eslint',
    '@nuxtjs/color-mode',
    '@nuxtjs/device',
    '@pinia/nuxt',
    '@primevue/nuxt-module',
    '@vueuse/nuxt',
  ],
})