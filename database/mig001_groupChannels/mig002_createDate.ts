import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    alter table minichat.channels
    add column created_at timestamp not null default now();
  `;
}
