import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.message_mentions (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      message_id uuid not null,
      user_id uuid not null
    )
  `;
}
