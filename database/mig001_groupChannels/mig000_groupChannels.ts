import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    alter table minichat.channels_private
    rename to channels_group
  `;
}
