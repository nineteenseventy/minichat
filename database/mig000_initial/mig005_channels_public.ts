import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels_public (
      id uuid PRIMARY KEY,
      title varchar(256) not null,
      description text
    )
  `;
}
