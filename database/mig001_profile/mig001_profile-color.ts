import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
  alter table minichat.users
  add column color varchar(16)
  `;
}
