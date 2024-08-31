import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.messages (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      author_id uuid not null,
      channel_id uuid not null,
      content text not null
    )
  `;
}
