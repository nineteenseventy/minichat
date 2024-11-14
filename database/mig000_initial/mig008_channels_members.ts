import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels_members (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      channel_id uuid not null,
      user_id uuid not null,
      last_read_message_timestamp timestamp with time zone not null default now(),
      joined_at_timestamp timestamp with time zone not null default now()
    )
  `;
}
