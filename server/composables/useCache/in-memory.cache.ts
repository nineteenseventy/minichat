import type {
  Cache,
  CacheValue,
  MemoryOptions,
  SetOptions,
} from './cache.interface';
import type { ConsolaInstance } from 'consola';

type CacheEntry = {
  value: string | number;
  expiresAt?: number;
};

export class InMemoryCache implements Cache {
  private cache: Map<string, CacheEntry> = new Map();

  constructor(
    private options: MemoryOptions,
    private logger: ConsolaInstance,
  ) {
    if (this.options.type !== 'in-memory') {
      throw new Error('Invalid cache type');
    }
    this.logger.log('Using in-memory cache');
  }

  async set(
    key: string,
    value: CacheValue,
    options?: SetOptions,
  ): Promise<CacheValue | null> {
    let result: CacheEntry | undefined = undefined;
    if (options?.get) {
      result = this.cache.get(key);
    }
    if (options?.ifExists && !this.cache.has(key)) {
      return result?.value ?? null;
    }
    if (options?.ifNotExists && this.cache.has(key)) {
      return result?.value ?? null;
    }

    let expiresAt: number | undefined = undefined;

    if (options?.ttl) {
      expiresAt = Date.now() + options.ttl;
    } else if (options?.keepTtl && result?.expiresAt) {
      expiresAt = result.expiresAt;
    }

    this.cache.set(key, { value, expiresAt });
    return result?.value ?? null;
  }

  async get(key: string): Promise<CacheValue | null> {
    const entry = this.cache.get(key);
    if (!entry) {
      return null;
    }

    if (entry.expiresAt && entry.expiresAt < Date.now()) {
      this.cache.delete(key);
      return null;
    }

    return entry.value;
  }

  async del(key: string): Promise<void> {
    this.cache.delete(key);
  }
}
