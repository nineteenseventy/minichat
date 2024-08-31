import { useEnv } from '../useEnv';
import type { Cache } from './cache.interface';
import { RedisCache } from './redis.cache';
import { useLogger } from '../useLogger';
import { InMemoryCache } from './in-memory.cache';

const cacheType = useEnv('CACHE_TYPE', 'in-memory');
const logger = useLogger('cache');
let cache: Cache;

if (cacheType === 'in-memory') {
  cache = new InMemoryCache(
    {
      type: 'in-memory',
    },
    logger,
  );
} else if (cacheType === 'redis') {
  cache = new RedisCache(
    {
      type: 'redis',
      host: useEnv('REDIS_HOST'),
      port: parseInt(useEnv('REDIS_PORT', '6379')),
      tls: useEnv('REDIS_TLS', 'false') === 'true',
    },
    logger,
  );
} else {
  throw new Error(`Invalid cache type: ${cacheType}`);
}

export const useCache = () => {
  return cache;
};
