export const useEnv = (key: string, defaultValue?: string) => {
  const value = process.env[key];
  if (value === undefined) {
    if (defaultValue === undefined) {
      throw new Error(`Missing environment variable: ${key}`);
    }
    return defaultValue;
  }
  return value;
};
