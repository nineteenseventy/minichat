import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.channels_public (
      id uuid PRIMARY KEY,
      name varchar(255) not null,
      description text not null,
      created_at timestamp not null default now()
    )
  `;
}
