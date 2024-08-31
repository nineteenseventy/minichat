import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels_direct (
      id uuid PRIMARY KEY,
      user_id uuid not null,
      user_id2 uuid not null
    )
  `;
}
