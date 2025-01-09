import { definePreset } from '@primevue/themes';
import Aura from '@primevue/themes/aura';
import { card } from './components/card';

export default definePreset(Aura, {
  components: {
    card,
  },
  css: () => ``,
});
