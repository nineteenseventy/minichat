import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.private_channel_members (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      channel_id uuid not null,
      user_id uuid not null
    )
  `;
}
