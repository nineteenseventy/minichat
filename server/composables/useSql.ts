import postgres from 'postgres';
import { useLogger } from './useLogger';

const logger = useLogger('sql');
const sql = postgres({
  onnotice: onNotice,
  debug: debug,
  onclose: onClose,
});

sql`SELECT 1;`.then(() => {
  logger.log('Connected to PostgreSQL');
});

export const useSql = () => {
  return sql;
};

function debug(
  connection: number,
  query: string,
  parameters: unknown[],
  paramTypes: unknown[],
) {
  const additional = [`connection: ${connection}`];
  if (parameters.length > 0)
    additional.push(`parameters: ${parameters.join(', ')}`);
  if (paramTypes.length > 0)
    additional.push(`paramTypes: ${paramTypes.join(', ')}`);

  logger.debug({
    message: query.trim(),
    additional,
  });
}

function onClose(connectionId: number) {
  logger.info(`Connection ${connectionId} closed`);
}

function onNotice(notice: postgres.Notice) {
  if (notice.message) {
    logger.info({
      message: notice.message,
      hint: notice.hint,
      arguments: notice.arguments,
    });
  } else {
    logger.warn(notice);
  }
}
