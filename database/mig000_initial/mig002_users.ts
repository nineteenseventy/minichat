import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.users (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      idp_id varchar(256) not null,
      username varchar(255) not null,
      color varchar(16),
      bio TEXT,
      picture varchar(256)
    )
  `;
}
