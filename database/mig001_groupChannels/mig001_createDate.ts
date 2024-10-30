import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    alter table minichat.channels_group
    drop column created_at;
    `;
}
