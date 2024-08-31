import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels_private (
      id uuid PRIMARY KEY,
      title varchar(255) not null
    )
  `;
}