import type { CardDesignTokens } from '@primevue/themes/types/card';
import type { ComponentType } from '../util';

export const card = {
  root: {
    borderRadius: '{content.border.radius}',
  },
  css: ({ dt }) => `
    .p-card {
    border: 1px solid ${dt('content.border.color')};
  }
    .p-card-body {
      max-width: 100%;
      max-height: 100%;
      height: 100%;
      width: 100%;
    }
    .p-card-content {
      flex: 1;
      overflow: hidden;
    }
  `,
} satisfies ComponentType<CardDesignTokens>;
