import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      type varchar(64) not null,
      created_at timestamp with time zone DEFAULT now()
    )
  `;
}
