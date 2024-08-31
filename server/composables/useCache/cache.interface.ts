export interface BaseOptions {
  type: string;
}

export interface MemoryOptions extends BaseOptions {
  type: 'in-memory';
}

export interface RedisOptions extends BaseOptions {
  type: 'redis';
  host: string;
  port: number;
  tls: boolean;
}

export type CacheOptions = MemoryOptions | RedisOptions;
export type CacheValue = string | number;

export interface SetOptions {
  get?: true;
  ifNotExists?: true;
  ifExists?: true;
  ttl?: number;
  keepTtl?: true;
}

export interface Cache {
  set(
    key: string,
    value: CacheValue,
    options?: SetOptions,
  ): Promise<CacheValue | null>;
  get(key: string): Promise<CacheValue | null>;
  del(key: string): Promise<void>;
}
