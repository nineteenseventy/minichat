import { ref, type Ref } from 'vue';
import { onBeforeRouteUpdate, useRoute } from 'vue-router';
import { unpackRouterParam } from '@/utils/router';

export function useRouteParam(key: string): Ref<string | undefined> {
  const route = useRoute();
  const current = unpackRouterParam(route.params[key]);
  const value = ref<string | undefined>(current);
  onBeforeRouteUpdate((to) => {
    value.value = unpackRouterParam(to.params[key]);
  });
  return value;
}
