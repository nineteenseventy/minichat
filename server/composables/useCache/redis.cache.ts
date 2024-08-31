import * as redis from 'redis';
import type {
  RedisOptions,
  Cache,
  SetOptions,
  CacheValue,
} from './cache.interface';
import type { ConsolaInstance } from 'consola';

export class RedisCache implements Cache {
  private clientConnectPromise: Promise<void>;

  constructor(
    private options: RedisOptions,
    private logger: ConsolaInstance,
  ) {
    if (this.options.type !== 'redis') {
      throw new Error('Invalid cache type');
    }
    this.logger.log('Using Redis cache');

    this.client = redis.createClient({
      socket: {
        tls: this.options.tls,
        port: this.options.port,
        host: this.options.host,
      },
    });

    this.clientConnectPromise = new Promise((resolve, reject) => {
      this.client
        .connect()
        .then(() => {
          resolve();
        })
        .catch((err: Error) => {
          reject(err);
        });
    });

    this.ensureConnected().then(() => {
      this.logger.log('Connected to Redis');
    });
  }

  private client: redis.RedisClientType;

  private async ensureConnected() {
    await this.clientConnectPromise;
  }

  async set(
    key: string,
    value: string | number,
    options?: SetOptions,
  ): Promise<CacheValue | null> {
    await this.ensureConnected();
    const redisOptions: redis.SetOptions = {
      GET: options?.get,
    };

    if (options?.keepTtl) redisOptions.KEEPTTL = true;
    else redisOptions.EX = options?.ttl;

    if (options?.ifExists) redisOptions.XX = true;
    else if (options?.ifNotExists) redisOptions.NX = true;

    return await this.client.set(key, value, redisOptions);
  }

  async get(key: string): Promise<CacheValue | null> {
    await this.ensureConnected();
    return await this.client.get(key);
  }

  async del(key: string) {
    await this.ensureConnected();
    await this.client.del(key);
  }
}
