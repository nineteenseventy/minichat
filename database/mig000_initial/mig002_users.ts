import type postgres from 'postgres';

export default function (sql: postgres.Sql) {
  return sql`
    create table minichat.users (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      profile_picture varchar(255),
      idp_id varchar(255) not null
    )
  `;
}
