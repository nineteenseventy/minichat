import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.attachments (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      message_id uuid not null,
      filename varchar(256) not null,
      type varchar(64) not null,
      url varchar(256) not null
    )
  `;
}
