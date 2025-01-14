let env: Record<string, string> | undefined = undefined;
export async function loadEnv() {
  if (env) return env;
  const response = await fetch('/assets/.env');
  env = await response.json();
  return env;
}
