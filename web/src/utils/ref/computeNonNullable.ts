import { computed, type ComputedRef, type Ref } from 'vue';

export type NonNullable<T> = T extends null | undefined ? never : T;

function IsNonNullable<T>(value: T): value is NonNullable<T> {
  return value !== null && value !== undefined;
}

export function computeNonNullable<T>(
  value: Ref<T>,
): ComputedRef<NonNullable<T>> {
  return computed(() => {
    const v = value.value;
    if (IsNonNullable(v)) {
      throw new Error('Value is null or undefined');
    }
    return v as NonNullable<T>;
  });
}
