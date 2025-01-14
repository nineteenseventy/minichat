export function unpackRouterParam(
  param: undefined | string | string[],
): string | undefined {
  if (Array.isArray(param)) {
    return param[0];
  }
  return param;
}
