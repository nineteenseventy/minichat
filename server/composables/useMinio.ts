import { Client } from 'minio';
import { useEnv } from './useEnv';
import { useLogger } from './useLogger';

const logger = useLogger('minio');
const minioclient = new Client({
  endPoint: useEnv('MINIO_ENDPOINT'),
  port: parseInt(useEnv('MINIO_PORT', '9000')),
  useSSL: useEnv('MINIO_USESSL', 'false') === 'true',
  accessKey: useEnv('MINIO_ACCESSKEY'),
  secretKey: useEnv('MINIO_SECRETKEY'),
});

minioclient
  .listBuckets()
  .then(() => {
    logger.log('Connected to Minio');
  })
  .catch((err) => {
    logger.error('Failed to connect to Minio', err);
  });

export const useMinio = () => {
  return minioclient;
};
