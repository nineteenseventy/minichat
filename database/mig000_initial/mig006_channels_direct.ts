import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels_direct (
      id uuid PRIMARY KEY,
      user1_id uuid not null,
      user2_id uuid not null
    )
  `;
}
