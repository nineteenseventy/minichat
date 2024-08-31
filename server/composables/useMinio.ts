import { Client } from 'minio';
import { useEnv } from './useEnv';

const minioclient = new Client({
  endPoint: useEnv('MINIO_ENDPOINT'),
  port: parseInt(useEnv('MINIO_PORT', '9000')),
  useSSL: useEnv('MINIO_USESSL', 'false') === 'true',
  accessKey: useEnv('MINIO_ACCESSKEY'),
  secretKey: useEnv('MINIO_SECRETKEY'),
});

export const useMinio = () => {
  return minioclient;
};
