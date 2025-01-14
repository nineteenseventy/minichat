function stringify(val: unknown, nested = false): string | undefined {
  if (val === undefined || val === null) {
    return undefined;
  }
  if (Array.isArray(val) && !nested) {
    return val.map((v) => stringify(v, true)).join(',');
  }
  if (val instanceof Date) {
    return val.toISOString();
  }
  if (val === null || val === undefined || typeof val === 'string') {
    return val;
  }
  if (val.toString && typeof val.toString === 'function') {
    return val.toString();
  }
  try {
    return JSON.stringify(val);
  } catch (e) {
    throw new Error(`Failed to stringify value: ${val}`);
  }
}

export function params(data: Record<string, unknown> | object): string {
  const _data: Record<string, string> = {};
  for (const key in data) {
    const val = data[key as keyof typeof data];
    const str = stringify(val);
    if (str !== undefined) {
      _data[key] = str;
    }
  }
  return new URLSearchParams(_data).toString();
}
