import { useMinio } from './useMinio';

export const useMinioBucket = async (bucketName: string) => {
  const client = useMinio();
  if (!(await client.bucketExists(bucketName))) {
    await client.makeBucket(bucketName);
  }
  return client;
};
