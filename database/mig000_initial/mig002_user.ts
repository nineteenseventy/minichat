import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.users (
    )
  `;
}
