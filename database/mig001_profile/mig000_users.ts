import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
  alter table minichat.users
  add column bio TEXT,
  add column username varchar(255) not null
  `;
}
