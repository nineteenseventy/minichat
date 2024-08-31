// import { useCache } from '../composables/useCache';
// import { useMinio } from '../composables/useMinio';
// import { useSql } from '../composables/useSql';

export default defineEventHandler(async () => {
  // const minio = useMinio();
  // minio.listBuckets().then((buckets) => {});

  // const sql = useSql();
  // await sql`SELECT 1;`;

  // const cache = useCache();
  // await cache.set('health', 'OK');
  // const result = await cache.get('health');
  // console.log('Health check:', result);

  return 1;
});
