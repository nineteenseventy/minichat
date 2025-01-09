type DesignTokenFunction = (key: string) => string;

export interface CssFunctionProps {
  dt: DesignTokenFunction;
}

type CssFunction = (props: CssFunctionProps) => string;

interface ComponentTypeExtras {
  css?: string | CssFunction;
}

export type ComponentType<T> = T & ComponentTypeExtras;
